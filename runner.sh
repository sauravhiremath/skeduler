#!/bin/bash
go build -o build/ .
echo "[*] Build completed"
./build/skeduler || echo "[x] SERVER INITIALISATION FAILED"
echo "[*] Server closed successfully"