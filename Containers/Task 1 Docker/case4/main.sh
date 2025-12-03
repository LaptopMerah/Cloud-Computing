#!/bin/bash

docker network create case4-network 2>/dev/null || echo "Network already exists"

bash run_database.sh
bash run_scrapper.sh
bash run_webserver.sh

echo "All containers started successfully!"
