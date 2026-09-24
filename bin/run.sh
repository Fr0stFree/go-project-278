#!/usr/bin/env bash
set -euo pipefail

echo "[run.sh] Starting service"

echo "[run.sh] Running database migrations"
goose -dir /app/db/migrations postgres "${DATABASE_URL}" up

echo "[run.sh] Starting backend and Caddy"
exec npx \
    --no-install \
    concurrently \
    --kill-others \
    "/app/bin/shortener" \
    "caddy run --config /etc/caddy/Caddyfile --adapter caddyfile"
