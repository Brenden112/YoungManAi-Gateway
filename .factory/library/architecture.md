# Architecture Reference — B2B Multi-Provider AI API Gateway

## Base Project

new-api (Go 1.22+, Gin, GORM v2). Module path: `github.com/QuantumNous/new-api`.

## Layered Call Chain for POST /v1/chat/completions

```
HTTP Request
  └─ router/relay-router.go          SetRelayRouter()
       └─ middleware/auth.go          TokenAuth()          — validate Bearer sk-xxx, set user/token context
       └─ middleware/distributor.go   Distribute()         — select Channel by model+group, set channel context
       └─ controller/relay.go         Relay()              — build RelayInfo, call relayHandler()
            └─ relay/compatible_handler.go  TextHelper()   — parse request, pre-consume quota
                 └─ relay/channel/adapter.go  Adaptor      — ConvertOpenAIRequest / DoRequest / DoResponse
                 └─ service/billing_session.go  BillingSession.Settle()
                 └─ model/log.go      RecordLog()
```

## Key Models

| File | Struct | Purpose |
|------|--------|---------|
| model/user.go | User | user account, Role, Quota, Group |
| model/token.go | Token | user API key (sk-xxx), hashed, RemainQuota |
| model/channel.go | Channel | upstream provider config, Type int, Key (encrypted) |
| model/log.go | Log | usage log per request |

## Channel Type System

Defined in `constant/channel.go`. Integer constants 0–57 (ChannelTypeDummy is sentinel).
New types must be added before ChannelTypeDummy and registered in:
- `constant/channel.go` ChannelTypeNames map
- `constant/channel.go` ChannelBaseURLs slice
- `relay/channel/adapter.go` factory switch (GetAdaptor)
- `relay/relay_adaptor.go` GetAdaptor()

## Provider Adapter Interface

```go
// relay/channel/adapter.go
type Adaptor interface {
    Init(info *relaycommon.RelayInfo)
    GetRequestURL(info *relaycommon.RelayInfo) (string, error)
    SetupRequestHeader(c *gin.Context, req *http.Header, info *relaycommon.RelayInfo) error
    ConvertOpenAIRequest(...) (any, error)
    DoRequest(...) (any, error)
    DoResponse(...) (usage any, err *types.NewAPIError)
    GetModelList() []string
    GetChannelName() string
    // ... other convert methods
}
```

## Provider Type Extension Plan

Add `ProviderType` string field to `model/Channel`:
- `official_cloud`   — direct upstream (OpenAI, Anthropic, etc.)
- `aggregator`       — multi-provider aggregator (OpenRouter, etc.)
- `authorized_proxy` — licensed reseller
- `experimental_proxy` — internal/experimental (kiro-gateway, etc.)

`experimental_proxy` defaults:
- risk_level = high
- available_scope = internal_only
- visibility = internal_only
- manual_enable_required = true
- enabled = false

## Access Control Extension Plan

Add `IsInternal bool` to `model/User` (or use Group="internal").
Enforce in `middleware/distributor.go`: if selected channel has ProviderType=experimental_proxy,
reject unless user IsInternal AND channel explicitly enabled.

## Billing Flow

1. `service/pre_consume_quota.go` — pre-deduct estimated quota
2. Channel adapter executes upstream request
3. `service/billing_session.go` BillingSession.Settle(actualQuota) — adjust delta
4. `model/log.go` RecordLog() — write usage_log

## Database Compatibility

All three DBs must work: SQLite, MySQL ≥5.7.8, PostgreSQL ≥9.6.
Use GORM abstractions. For reserved words use `commonGroupCol`/`commonKeyCol` from model/main.go.
Boolean literals: `commonTrueVal`/`commonFalseVal`.

## Frontend

- Default: `web/default/` — React 19, Rsbuild, Base UI, Tailwind, Bun
- Classic: `web/classic/` — React 18, Vite, Semi Design
- Admin pages live under `web/default/src/pages/`
- i18n: `web/default/src/i18n/locales/{lang}.json`

## Security Invariants

- Upstream channel keys stored encrypted (AES via common/crypto.go)
- User API keys stored as SHA-256 hash prefix + random suffix
- No full prompt/response stored by default
- experimental_proxy never visible to non-internal users
- No automatic fallback from official_cloud → experimental_proxy
