#!/usr/bin/env bash
set -euo pipefail

# Runs the pre-release migration schema smoke. With no SQL_DSN, it uses a
# temporary SQLite database. Set SQL_DSN to a MySQL or PostgreSQL DSN for
# service-backed CI runs.

export SESSION_SECRET="${SESSION_SECRET:-migration-fixture-session-secret}"
export CRYPTO_SECRET="${CRYPTO_SECRET:-migration-fixture-crypto-secret}"
export REDIS_CONN_STRING="${REDIS_CONN_STRING:-}"
export GOCACHE="${GOCACHE:-/tmp/go-build-cache}"
export GOPATH="${GOPATH:-/tmp/go-path}"

GO_BIN="${GO_BIN:-/tmp/go2510/bin/go}"
if [ ! -x "$GO_BIN" ]; then
  GO_BIN="go"
fi

"$GO_BIN" test ./tests/migration -count=1
