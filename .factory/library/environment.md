# Environment Reference

## Runtime Requirements

| Component | Version |
|-----------|---------|
| Go | 1.22+ |
| Node.js | 18+ (for frontend) |
| Bun | latest (preferred package manager) |
| Docker | 20+ (optional, for compose) |

## Database Options (pick one)

| DB | Min Version | DSN env var |
|----|-------------|-------------|
| SQLite | 3.x | `SQL_DSN` empty or file path |
| MySQL | 5.7.8+ | `SQL_DSN=user:pass@tcp(host:3306)/dbname` |
| PostgreSQL | 9.6+ | `SQL_DSN=host=... user=... dbname=...` |

## Key Environment Variables

```bash
# Core
SQL_DSN=                        # empty = SQLite ./data/new-api.db
REDIS_CONN_STRING=              # optional, redis://...
SESSION_SECRET=                 # required in production
INITIAL_ROOT_TOKEN=             # optional, set root API key on first boot

# Security
CRYPTO_SECRET=                  # AES key for encrypting channel keys

# Feature flags
LOG_CONSUME_ENABLED=true        # write consume logs
STORE_FULL_TEXT_ENABLED=false   # default: do NOT store prompt/response

# Experimental proxy
EXPERIMENTAL_PROXY_ENABLED=false  # global kill-switch
```

## Build Commands

```bash
# Backend
go build -o new-api .

# Frontend (default theme)
cd web/default && bun install && bun run build

# Full build via makefile
make build
```

## Dev Server

```bash
# Backend (hot reload via air)
air

# Frontend dev server
cd web/default && bun run dev
```

## Test Commands

```bash
# Unit tests
go test ./...

# Specific package
go test ./service/... -v

# Frontend
cd web/default && bun run test
```

## Docker

```bash
# Dev compose (SQLite + app)
docker compose -f docker-compose.dev.yml up

# Production compose
docker compose up -d
```

## Ports

| Service | Default Port |
|---------|-------------|
| API Gateway | 3000 |
| Frontend dev | 3001 |
