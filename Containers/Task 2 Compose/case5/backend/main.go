package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type shortenRequest struct {
	URL  string `json:"url"`
	Code string `json:"code"`
}

type shortenResponse struct {
	Code string `json:"code"`
	URL  string `json:"url"`
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func openDBWithRetry() (*sql.DB, error) {
	host := getenv("DB_HOST", "database")
	port := getenv("DB_PORT", "3306")
	user := getenv("DB_USER", "case5user")
	pass := getenv("DB_PASS", "case5pass")
	name := getenv("DB_NAME", "shorten_link")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true", user, pass, host, port, name)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	maxRetries := 30
	if v := os.Getenv("DB_MAX_RETRIES"); v != "" {
		fmt.Sscanf(v, "%d", &maxRetries)
	}
	retryDelay := 2 * time.Second
	if v := os.Getenv("DB_RETRY_DELAY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			retryDelay = d
		}
	}
	for i := 0; i < maxRetries; i++ {
		if err := db.Ping(); err == nil {
			// Ensure table exists (idempotent)
			_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS links (
				id INT NOT NULL AUTO_INCREMENT,
				code VARCHAR(64) NOT NULL,
				url TEXT NOT NULL,
				created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
				PRIMARY KEY (id),
				UNIQUE KEY uniq_code (code)
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`)
			return db, nil
		}
		log.Printf("waiting for database... (%d/%d)", i+1, maxRetries)
		time.Sleep(retryDelay)
	}
	return db, fmt.Errorf("database not reachable after retries")
}

func main() {
	db, err := openDBWithRetry()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	mux.HandleFunc("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req shortenRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		url := strings.TrimSpace(req.URL)
		if url == "" {
			http.Error(w, "url required", http.StatusBadRequest)
			return
		}
		code := strings.TrimSpace(req.Code)
		if code == "" {
			http.Error(w, "code required", http.StatusBadRequest)
			return
		}
		if !isValidCode(code) {
			http.Error(w, "invalid code", http.StatusBadRequest)
			return
		}
		if _, err := db.Exec("INSERT INTO links (code, url) VALUES (?, ?)", code, url); err != nil {
			if isDuplicateError(err) {
				http.Error(w, "code already exists", http.StatusConflict)
				return
			}
			log.Printf("failed to store code: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shortenResponse{Code: code, URL: url})
	})

	// Optional: fetch original by code
	mux.HandleFunc("/api/lookup/", func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.Path, "/api/lookup/")
		if code == "" {
			http.Error(w, "code required", http.StatusBadRequest)
			return
		}
		var url string
		err := db.QueryRow("SELECT url FROM links WHERE code = ?", code).Scan(&url)
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		} else if err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(shortenResponse{Code: code, URL: url})
	})

	mux.HandleFunc("/api/links", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		rows, err := db.Query("SELECT code, url FROM links ORDER BY created_at DESC LIMIT 50")
		if err != nil {
			log.Printf("failed to list links: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var results []shortenResponse
		for rows.Next() {
			var item shortenResponse
			if err := rows.Scan(&item.Code, &item.URL); err != nil {
				log.Printf("failed to scan row: %v", err)
				http.Error(w, "database error", http.StatusInternalServerError)
				return
			}
			results = append(results, item)
		}
		if err := rows.Err(); err != nil {
			log.Printf("rows error: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		path := strings.Trim(r.URL.Path, "/")
		if path == "" {
			http.NotFound(w, r)
			return
		}
		if strings.Contains(path, "/") {
			http.NotFound(w, r)
			return
		}
		if !isValidCode(path) {
			http.NotFound(w, r)
			return
		}
		var dest string
		err := db.QueryRow("SELECT url FROM links WHERE code = ?", path).Scan(&dest)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		} else if err != nil {
			log.Printf("lookup failed: %v", err)
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, dest, http.StatusFound)
	})

	addr := ":8080"
	log.Printf("backend listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func isValidCode(code string) bool {
	if len(code) < 3 || len(code) > 32 {
		return false
	}
	for _, r := range code {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}
