#!/bin/bash
export MONGO_DBNAME=example_db
export MONGO_PASSWORD=example_password

now=$(date '+%d/%m/%Y %H:%M:%S')

go build -o build/ .
echo "$now [*] Build completed"

go test ./...

./build/skeduler || sudo kill -9 $(sudo lsof -t -i:8080) || echo "[x] SERVER INITIALISATION FAILED"
echo "[*] Server closed successfully"