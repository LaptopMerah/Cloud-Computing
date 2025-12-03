import requests
import mysql.connector
import time

API_URL = "https://api.chucknorris.io/jokes/random"

DB_CONFIG = {
    'host': 'database-mysql',
    'user': 'case4',
    'password': 'passdb_case4',
    'database': 'db_case4'
}

def wait_for_database():
    print("Waiting for database to be ready...")
    while True:
        try:
            conn = mysql.connector.connect(**DB_CONFIG)
            conn.close()
            print("Database is ready!")
            return
        except mysql.connector.Error:
            print("Database not ready, retrying in 5 seconds...")
            time.sleep(5)

def create_table():
    conn = mysql.connector.connect(**DB_CONFIG)
    cursor = conn.cursor()
    cursor.execute('''
        CREATE TABLE IF NOT EXISTS jokes (
            id INT AUTO_INCREMENT PRIMARY KEY,
            joke_id VARCHAR(50) UNIQUE,
            joke_text TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    ''')
    conn.commit()
    cursor.close()
    conn.close()

def fetch_and_store_joke():
    try:
        response = requests.get(API_URL)
        data = response.json()
        
        conn = mysql.connector.connect(**DB_CONFIG)
        cursor = conn.cursor()
        
        cursor.execute(
            "INSERT IGNORE INTO jokes (joke_id, joke_text) VALUES (%s, %s)",
            (data['id'], data['value'])
        )
        
        conn.commit()
        print(f"Stored joke: {data['id']}")
        
        cursor.close()
        conn.close()
    except Exception as e:
        print(f"Error: {e}")

if __name__ == "__main__":
    wait_for_database()
    create_table()
    
    print("Starting scrapper...")
    while True:
        fetch_and_store_joke()
        time.sleep(10)
