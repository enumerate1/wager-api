#!/bin/bash

echo "start wager-api, database ..."
docker compose -f  $(pwd)/docker-compose.yml up --build --abort-on-container-exit 