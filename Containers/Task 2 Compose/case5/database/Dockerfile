FROM mysql:8.0

# Environment variables will be overridden by docker-compose
ENV MYSQL_ROOT_PASSWORD=case5rootpass \
    MYSQL_DATABASE=shorten_link \
    MYSQL_USER=case5user \
    MYSQL_PASSWORD=case5pass

# Initialize schema and tables on first run
COPY init.sql /docker-entrypoint-initdb.d/01-init.sql

# Expose default MySQL port
EXPOSE 3306
