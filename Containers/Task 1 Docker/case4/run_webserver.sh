#!/bin/bash

docker rm -f webserver

docker run \
    -dit \
    --name webserver \
    --network case4-network \
    -e DB_HOST=database-mysql \
    -e DB_USER=case4 \
    -e DB_PASSWORD=passdb_case4 \
    -e DB_NAME=db_case4 \
    -v "$(pwd)/web:/var/www/html" \
    --publish 5555:80 \
    php:8-apache \
    bash -c "docker-php-ext-install mysqli && apache2-foreground"

echo "Webserver started!"
