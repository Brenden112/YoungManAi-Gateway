#!/usr/bin/env bash
set -euo pipefail

# CI wrapper for cross-database migration checks. With no SQL_DSN it uses the
# project's SQLite migration smoke. In CI jobs, set SQL_DSN for MySQL or
# PostgreSQL service-backed runs. The underlying test uses GORM migration paths.

if ! command -v go >/dev/null 2>&1; then
  echo "BLOCKED blocked_external_dependency_cross_db_runtime: go binary not found"
  exit 2
fi

bash scripts/check-migrations.sh
