#!/bin/bash
export MONGO_DBNAME=example_db
export MONGO_PASSWORD=example_password

go build -o build/ .
echo "[*] Build completed"

./build/skeduler || sudo kill -9 $(sudo lsof -t -i:8080) || echo "[x] SERVER INITIALISATION FAILED"
echo "[*] Server closed successfully"