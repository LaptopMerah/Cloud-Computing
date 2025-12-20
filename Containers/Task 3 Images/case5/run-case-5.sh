#!/bin/bash

docker rm -f case5-jupyter

docker run -d \
    --name case5-jupyter \
    -p 8888:8888 \
    -v "./jupyter/notebooks:/home/jupyter/notebooks" \
    case5-jupyter:latest



