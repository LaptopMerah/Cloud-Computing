#!/bin/bash

docker rm -f database-mysql

docker run \
    -dit \
    --name database-mysql \
    --network case4-network \
    -v "$(pwd)/DB-data:/var/lib/mysql" \
    -e MYSQL_DATABASE=db_case4 \
    -e MYSQL_PASSWORD=passdb_case4 \
    -e MYSQL_ROOT_PASSWORD=mydb6789tyui \
    -e MYSQL_ROOT_HOST=% \
    -e MYSQL_USER=case4 \
    mysql:8.0-debian

echo "Database MySQL started!"
