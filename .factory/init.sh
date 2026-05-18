#!/usr/bin/env bash
set -euo pipefail

echo "[init] Setting up B2B AI Gateway dev environment..."

# 1. Go dependencies
echo "[init] Downloading Go modules..."
go mod download

# 2. Frontend dependencies
if command -v bun &>/dev/null; then
  echo "[init] Installing frontend deps with bun..."
  (cd web/default && bun install)
else
  echo "[WARN] bun not found, skipping frontend install"
fi

# 3. Create data directory for SQLite
mkdir -p data

# 4. Copy env example if .env doesn't exist
if [ ! -f .env ]; then
  cp .env.example .env
  echo "[init] Created .env from .env.example — review and update secrets before running"
fi

# 5. Verify Go build
echo "[init] Verifying Go build..."
go build -o /tmp/new-api-check . && rm /tmp/new-api-check
echo "[init] Go build OK"

# 6. Run unit tests
echo "[init] Running unit tests..."
go test ./... -count=1 -timeout 60s

echo "[init] Done. Run 'go run .' to start the backend."
