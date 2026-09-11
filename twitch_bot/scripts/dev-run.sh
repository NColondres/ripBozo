#!/usr/bin/env bash

set -euo pipefail

URL="https://subgroup-navigator-bondless.ngrok-free.dev"
PORT=8080

if [ -z "$URL" ]; then
    echo "No URL set in script. Please update script"
    exit 1
fi

if ! ngrok --version 2>&1 > /dev/null; then
    echo "localtunnel is not installed"
    exit 1
fi

ngrok http 8080 --url "$URL" > /dev/null 2>&1 &

NGROK_PID=$!

cleanup() {
  echo "Cleaning up..."
  kill $NGROK_PID 2>/dev/null || true
  kill $GO_PID 2>/dev/null || true
}

trap cleanup EXIT


export TWITCH_CALLBACK_URL="$URL"
export TWITCH_REDIRECT_URL="http://localhost:8080"

echo "ngrok public url: $URL"


./bin/ripBozo &

GO_PID=$!
wait $GO_PID
