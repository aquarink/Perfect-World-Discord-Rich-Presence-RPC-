#!/usr/bin/env bash
set -e

# Load local production config if present (ignored by git)
if [ -f "prod.env" ]; then
  # shellcheck source=/dev/null
  source prod.env
fi

LDFLAGS="-s -w -H windowsgui"
if [ -n "$CLIENT_ID" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildClientID=$CLIENT_ID'"
fi
if [ -n "$DETAILS" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildDetails=$DETAILS'"
fi
if [ -n "$STATE" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildState=$STATE'"
fi
if [ -n "$LARGE_IMAGE" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildLargeImage=$LARGE_IMAGE'"
fi
if [ -n "$LARGE_TEXT" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildLargeText=$LARGE_TEXT'"
fi
if [ -n "$SMALL_IMAGE" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildSmallImage=$SMALL_IMAGE'"
fi
if [ -n "$SMALL_TEXT" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildSmallText=$SMALL_TEXT'"
fi
if [ -n "$WEB_URL" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildWebUrl=$WEB_URL'"
fi
if [ -n "$DISCORD_URL" ]; then
  LDFLAGS="$LDFLAGS -X 'main.BuildDiscordUrl=$DISCORD_URL'"
fi

echo "=== Building Perfect World Discord RPC Starter ==="

mkdir -p bin

echo "[1/2] Compiling Windows 32-bit (x86) -> bin/RealmOfChaos.exe"
GOOS=windows GOARCH=386 go build -ldflags="$LDFLAGS" -o bin/RealmOfChaos.exe .

echo "[2/2] Compiling Windows 64-bit (x64) -> bin/RealmOfChaos_x64.exe"
GOOS=windows GOARCH=amd64 go build -ldflags="$LDFLAGS" -o bin/RealmOfChaos_x64.exe .

cp config.json bin/config.json

echo "=== Build Complete! Executables are located in bin/ ==="
ls -lh bin/
