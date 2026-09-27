#!/usr/bin/env bash
set -e

echo "=== Building Perfect World Discord RPC Starter ==="

mkdir -p bin

echo "[1/2] Compiling Windows 32-bit (x86) -> bin/RealmOfChaos.exe"
GOOS=windows GOARCH=386 go build -ldflags="-s -w -H windowsgui" -o bin/RealmOfChaos.exe .

echo "[2/2] Compiling Windows 64-bit (x64) -> bin/RealmOfChaos_x64.exe"
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o bin/RealmOfChaos_x64.exe .

cp config.json bin/config.json

echo "=== Build Complete! Executables are located in bin/ ==="
ls -lh bin/
