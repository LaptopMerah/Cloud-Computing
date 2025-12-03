#!/bin/bash

docker rm -f scrapper-chucknorris

# Run the scrapper container directly with Python image
docker run \
    -dit \
    --name scrapper-chucknorris \
    --network case4-network \
    -e DB_HOST=database-mysql \
    -e DB_USER=case4 \
    -e DB_PASSWORD=passdb_case4 \
    -e DB_NAME=db_case4 \
    -e API_URL=https://api.chucknorris.io/jokes/random \
    -e INTERVAL=10 \
    -v "$(pwd)/scripts:/app" \
    python:3.11-slim \
    bash -c "pip install requests mysql-connector-python && python -u /app/scrapper.py"

echo "Scrapper started!"
