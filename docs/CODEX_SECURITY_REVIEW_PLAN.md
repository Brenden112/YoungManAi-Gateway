# Codex Security Review Plan

## Objective

Assess whether M0-M16 preserve the MVP security invariants:

- No real provider key exposure.
- User API keys are not stored or returned in plaintext.
- experimental_proxy is hidden and blocked for normal users.
- official_cloud failures never fallback into experimental_proxy.
- Billing cannot be bypassed by error, retry, or quota edge cases.
- Full prompts/responses are not logged by default.
- Admin-only APIs cannot be reached by normal users.
- Migrations remain safe across SQLite, MySQL, and PostgreSQL.

## Threat Model

Actors:

- Anonymous caller with no token.
- Normal user with valid token.
- Internal user with valid token but `AllowExperimental=false`.
- Internal user with `AllowExperimental=true`.
- Admin/root user.
- Compromised or misconfigured channel/provider account.
- Local operator running docker compose with default credentials.

Assets:

- User API keys.
- Upstream provider keys.
- User quota and billing records.
- usage_log contents.
- Admin channel/token/balance controls.
- experimental_proxy endpoints.

## Security Review Areas

### Key And Secret Handling

Audit files:

- `model/token.go`
- `controller/token.go`
- `middleware/auth.go`
- `model/provider_account.go`
- `common/crypto.go`
- `controller/channel.go`

Checks:

- Token creation computes `KeyHash` and `KeyPrefix`.
- API responses mask keys except explicitly owner-authorized full-key retrieval paths.
- DB scan verifies no new token row stores full `sk-` in plaintext.
- ProviderAccount stores encrypted upstream key and does not log raw key.
- `CRYPTO_SECRET` behavior is safe when unset or invalid.
- `common/json.go` wrapper is used for JSON marshal/unmarshal in business code.

Commands:

```bash
rg -n 'KeyHash|KeyPrefix|MaskTokenKey|GetMaskedKey|GetTokenByKey|sk-' model controller middleware service
rg -n 'EncryptedKey|Encrypt|Decrypt|CRYPTO_SECRET|ProviderAccount' common model controller
rg -n 'json\\.Marshal|json\\.Unmarshal|json\\.NewDecoder|json\\.NewEncoder|encoding/json' --glob '*.go'
```

High-risk questions:

- Does backward compatibility keep plaintext `Token.Key` persisted?
- Can admin token list or search leak full `sk-` values?
- Can ProviderAccount encrypted key be overwritten or decrypted into logs?

### experimental_proxy Isolation

Audit files:

- `constant/channel.go`
- `constant/context_key.go`
- `model/channel.go`
- `model/channel_cache.go`
- `model/ability.go`
- `service/internal_user.go`
- `service/channel_select.go`
- `middleware/distributor.go`
- `controller/relay.go`
- `controller/model.go`
- `controller/channel.go`

Checks:

- experimental_proxy defaults disabled, high-risk, internal-only, manual-enable-required.
- Candidate selection excludes experimental_proxy unless explicitly allowed.
- Allow requires both internal user and token `AllowExperimental`.
- Non-internal users cannot see experimental models in `/v1/models`.
- Non-admin channel APIs hide or forbid experimental_proxy data as specified.
- Retry logic does not cross into experimental_proxy on official/aggregator failure.
- `ContextKeyIsExperimentalProxy` and `ContextKeyChannelProviderType` are set only after valid channel selection.

Commands:

```bash
rg -n 'experimental_proxy|ProviderTypeExperimentalProxy|AllowExperimental|ContextKeyIsExperimentalProxy|ContextKeyTokenAllowExperimental' constant model service middleware controller relay
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./model/... ./service/... -run 'Experimental|Internal|ChannelCache|ProviderType' -count=1
```

High-risk questions:

- Can explicit channel selection bypass `Distribute` experimental guard?
- Can retry call sites omit `AllowExperimental`?
- Does channel cache preserve stale enabled experimental channels after admin disable?

### Billing And Quota Integrity

Audit files:

- `service/pre_consume_quota.go`
- `service/billing_session.go`
- `service/quota.go`
- `service/text_quota.go`
- `service/tiered_settle.go`
- `controller/relay.go`
- `model/log.go`
- `model/token.go`
- `model/user.go`

Checks:

- Zero balance returns 402/429 before upstream call.
- Pre-consumed quota is returned on failed request.
- Successful settle updates user and token usage exactly once.
- Error logs do not trigger usage stats or quota deduction.
- Tiered billing snapshots cannot be modified between pre-consume and settle.
- Concurrent settle/refund paths are idempotent or guarded.

Commands:

```bash
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./service/... -run 'PreConsume|Quota|Settle|Tiered|Refund|CAS' -count=1
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./model/... -run 'Log|TokenUsage|Regression' -count=1
rg -n 'PreConsumeQuota|ReturnPreConsumedQuota|PostConsumeQuota|DecreaseUserQuota|IncreaseUserQuota|GetTokenUsageStats|LogTypeError|LogTypeConsume' service model controller
```

High-risk questions:

- Can a panic or streaming disconnect avoid quota settlement?
- Can failed upstream calls deduct final quota?
- Can manual top-up be called by non-admin users?

### Logging Privacy

Audit files:

- `common/constants.go`
- `model/log.go`
- `controller/log.go`
- `controller/relay.go`
- `service/log_info_generate.go`
- `types/error.go`

Checks:

- `StoreFullTextEnabled` defaults false.
- `RecordConsumeLog` clears `Content` when full text storage is disabled.
- `Other` and error strings do not carry full prompt/completion.
- Admin log query does not expose prompt_text/completion_text fields.
- Non-admin log query scope cannot reveal other users' data.
- `MaskSensitiveErrorWithStatusCode` is used for error logs.

Commands:

```bash
rg -n 'StoreFullTextEnabled|Content|Other|RecordConsumeLog|RecordErrorLog|MaskSensitiveErrorWithStatusCode|prompt|completion' common model controller service types
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./model/... -run 'Log|StoreFullText|Privacy' -count=1
```

High-risk questions:

- Are prompts/responses stored in `Log.Content`, `Other`, or nested JSON under a different key?
- Does streaming path log partial completions?
- Does error masking cover provider raw response bodies?

### Admin API Authorization

Audit files:

- `router/api-router.go`
- `middleware/auth.go`
- `controller/channel.go`
- `controller/token.go`
- `controller/log.go`
- `controller/user.go`

Checks:

- Provider summary, disable experimental, admin token list, log filters, and top-up routes are under admin/root auth.
- Normal users cannot access admin data by guessing routes.
- Token list always masks key material.
- Log list supports filters but remains scoped by role.
- Manual recharge writes `LogTypeManage`.

Commands:

```bash
rg -n 'AdminAuth|RootAuth|disable-experimental|provider-summary|admin/token|GetAdminAllTokens|add_quota|topup|GetAllLogs' router controller middleware
```

Required supplemental tests:

- Normal user requests to each admin route return 401/403.
- Admin route responses are checked for secret leakage.
- Top-up audit log is created and tied to admin actor.

### Database And Migration Security

Audit files:

- `model/main.go`
- all new model files under `model/`

Checks:

- New columns are nullable or have safe defaults.
- SQLite migrations avoid unsupported `ALTER COLUMN`.
- Raw SQL quotes reserved columns with `commonGroupCol`/`commonKeyCol`.
- MySQL/PostgreSQL boolean and string syntax are compatible.
- No destructive migration drops or truncates user data.

Commands:

```bash
rg -n 'AutoMigrate|ALTER TABLE|DROP TABLE|TRUNCATE|GROUP_CONCAT|STRING_AGG|JSONB|commonGroupCol|commonKeyCol|commonTrueVal|commonFalseVal' model
```

Blocked until DB fixtures exist:

- MySQL 5.7.8+ AutoMigrate run.
- PostgreSQL 9.6+ AutoMigrate run.

### Docker And Local Deployment

Audit files:

- `docker-compose.yml`
- `docker-compose.dev.yml`
- `docker-compose.local.yml`
- `scripts/regression.sh`
- docs quickstart/policy files

Checks:

- No real secrets are embedded.
- Local compose does not expose unsafe default admin token beyond documented local use.
- Regression script supports mock/local endpoints and does not require paid providers.
- Health check and smoke tests are deterministic.

Commands:

```bash
rg -n 'sk-|OPENAI_API_KEY|ANTHROPIC|GEMINI|SECRET|TOKEN|PASSWORD|localhost|BASE_URL|ADMIN_TOKEN' docker-compose*.yml scripts docs
docker compose -f docker-compose.local.yml config
```

Do not run compose against real provider configs. If the script requires real tokens, mark blocked.

## Security Test Order

1. Static secret scan and route authorization scan.
2. Unit tests for keys, provider accounts, internal user, experimental routing, billing, logs.
3. DB migration/schema checks with SQLite.
4. Mock integration for auth, channel selection, billing, and logs.
5. Admin route negative tests as normal user.
6. Local docker compose config validation.
7. Local smoke script only after mock upstream and local test tokens exist.
8. Optional MySQL/PostgreSQL migration tests with local disposable DBs.

## Blocked Security Items

- Missing requested security/routing/billing/logging spec docs.
- No confirmed local MySQL/PostgreSQL services for cross-DB migration security checks.
- No confirmed mock upstream server for full curl smoke tests.
- Frontend route/component authorization requires local Bun dependencies and frontend test fixtures.
- Any real provider call is blocked by audit rules.

## Expected Security Audit Output

The execution phase should produce:

- Finding list ordered by severity.
- Feature/assertion evidence matrix with pass/partial/fail/blocked status.
- Commands run and exit codes.
- Missing tests list.
- Runtime logs or excerpts for local/mocked flows.
- Explicit statement that no real provider keys or paid upstreams were used.

