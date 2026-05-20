# Validation Contract

Every feature in features.json must satisfy the assertions listed here before it can be marked `done`.
Assertions are grouped by feature ID. Each assertion has a type: `unit`, `integration`, `curl`, or `db`.

## Current Status Note

This contract is the original M0-M16 validation baseline. The implementation has since gone through remediation, fixture hardening, and current-head CI verification. As of 2026-05-20, Pre-release verification #16 passed on commit `73ad2ff`, and the canonical current state is tracked in:

- `.factory/mission-state.json`
- `docs/CI_VERIFICATION_EVIDENCE.md`
- `docs/CODEX_FEATURE_TEST_RESULTS.md`
- `docs/DEVELOPMENT_LOG.md`

The assertions below remain useful as acceptance criteria and historical traceability, but current release readiness is determined by mission-state plus CI/staging evidence. Production readiness remains `not_ready` until isolated staging runtime verification and manual release sign-off are complete.

---

## M1 — Provider Type & Field Migration

### F1.1 — Channel model new fields
- **A1.1.1** `[unit]` `go test ./model/... -run TestChannelProviderType` passes: Channel struct has fields `ProviderType`, `RiskLevel`, `AvailableScope`, `Visibility`, `ManualEnableRequired`.
- **A1.1.2** `[db]` After AutoMigrate on SQLite, MySQL, and PostgreSQL, `PRAGMA table_info(channels)` / `DESCRIBE channels` / `\d channels` shows all five new columns.
- **A1.1.3** `[unit]` Creating a Channel with `ProviderType=""` defaults to `"official_cloud"` via GORM default tag.

### F1.2 — DB migration compatibility
- **A1.2.1** `[integration]` `go test ./model/... -run TestMigrate` passes against all three DB drivers without error.
- **A1.2.2** `[db]` Existing Channel rows after migration have `provider_type = 'official_cloud'` (backfill).

---

## M2 — Provider Account & Channel Association

### F2.1 — ProviderAccount model
- **A2.1.1** `[unit]` `go test ./model/... -run TestProviderAccount` passes: CRUD operations work.
- **A2.1.2** `[db]` `provider_accounts` table exists with columns: `id`, `name`, `provider_type`, `encrypted_key`, `created_at`.
- **A2.1.3** `[unit]` `ProviderAccount.EncryptedKey` is stored as AES-encrypted ciphertext; plaintext key is never written to DB.

### F2.2 — Channel ProviderAccountId FK
- **A2.2.1** `[db]` `channels` table has nullable `provider_account_id` column.
- **A2.2.2** `[unit]` Existing channels with `provider_account_id = NULL` continue to function (backward compat).

---

## M3 — OpenAI-compatible API Regression

### F3.1 — GET /v1/models
- **A3.1.1** `[curl]` `GET /v1/models` with valid Bearer token returns HTTP 200 and JSON body with `"object": "list"` and non-empty `data` array.
- **A3.1.2** `[curl]` `GET /v1/models` with invalid token returns HTTP 401.

### F3.2 — POST /v1/chat/completions
- **A3.2.1** `[curl]` `POST /v1/chat/completions` with valid token and `{"model":"<enabled_model>","messages":[...]}` returns HTTP 200 with `choices` array and `usage.prompt_tokens > 0`.
- **A3.2.2** `[unit]` `go test ./...` — all pre-existing tests still pass (no regression).

---

## M4 — Official Provider Minimal Integration

### F4.1 — Backfill provider_type
- **A4.1.1** `[db]` After migration, all channels with `type IN (1,3,14,24)` (OpenAI, Azure, Anthropic, Gemini) have `provider_type = 'official_cloud'`.
- **A4.1.2** `[unit]` `GetDefaultProviderType(channelType int) string` returns correct defaults for known types.

### F4.2 — Channel CRUD accepts provider_type
- **A4.2.1** `[curl]` `POST /api/channel` with `{"provider_type":"aggregator",...}` creates channel with correct provider_type.
- **A4.2.2** `[curl]` `GET /api/channel/:id` response includes `provider_type` field.

---

## M5 — Model Mapping & Price Basics

### F5.1 — Channel selection respects provider_type
- **A5.1.1** `[unit]` `service.SelectChannel` with group filter returns only channels matching the requested provider_type when filter is set.
- **A5.1.2** `[unit]` Channel with `provider_type = 'experimental_proxy'` is excluded from default selection pool.

### F5.2 — Pricing for new provider types
- **A5.2.1** `[unit]` `service.GetQuotaByChannel` does not panic for channels with any of the four provider_type values.

---

## M6 — KiroGatewayAdapter Skeleton

### F6.1 — ChannelTypeKiroGateway constant
- **A6.1.1** `[unit]` `constant.ChannelTypeKiroGateway` is defined, value is unique (not conflicting with existing types).
- **A6.1.2** `[unit]` `constant.ChannelTypeNames[constant.ChannelTypeKiroGateway]` returns `"KiroGateway"`.

### F6.2 — KiroGatewayAdapter implements Adaptor
- **A6.2.1** `[unit]` `var _ channel.Adaptor = (*KiroGatewayAdaptor)(nil)` compiles without error (interface compliance).
- **A6.2.2** `[unit]` `go build ./...` passes with no errors.

### F6.3 — Registered in GetAdaptor factory
- **A6.3.1** `[unit]` `relay.GetAdaptor(constant.ChannelTypeKiroGateway)` returns non-nil adaptor.

---

## M7 — experimental_proxy Access Control

### F7.1 — Internal user identification
- **A7.1.1** `[unit]` `isInternalUser(user)` returns true when `user.Group == "internal"`.
- **A7.1.2** `[unit]` `isInternalUser(user)` returns false for `user.Group == "default"`.

### F7.2 — Distribute blocks non-internal users
- **A7.2.1** `[integration]` Normal user (group=default) calling a model served only by experimental_proxy channel receives HTTP 403 with error code `experimental_proxy_forbidden`.
- **A7.2.2** `[integration]` Internal user (group=internal) calling same model with channel enabled receives HTTP 200.

### F7.3 — experimental_proxy default disabled
- **A7.3.1** `[unit]` `NewExperimentalProxyChannel()` sets `Status = ChannelStatusDisabled` and `ManualEnableRequired = true`.
- **A7.3.2** `[db]` Channels with `provider_type = 'experimental_proxy'` created via API have `status = 2` (disabled) by default.

---

## M8 — experimental_proxy Routing & Fallback Isolation

### F8.1 — No fallback to experimental_proxy
- **A8.1.1** `[unit]` `service.SelectChannel` with `excludeExperimental=true` never returns a channel with `provider_type = 'experimental_proxy'`.
- **A8.1.2** `[integration]` When all official_cloud channels for a model fail, the retry loop does NOT attempt experimental_proxy channels; returns 503.

### F8.2 — One-click disable
- **A8.2.1** `[curl]` `PATCH /api/channel/:id` with `{"status":2}` on an experimental_proxy channel disables it; subsequent calls return 503/no-channel.
- **A8.2.2** `[unit]` Channel cache is invalidated after status update.

---

## M9 — Organization / Project Base Tables

### F9.1 — Organization table
- **A9.1.1** `[db]` `organizations` table exists with columns: `id`, `name`, `status`, `created_at`.
- **A9.1.2** `[unit]` AutoMigrate passes on all three DB drivers.

### F9.2 — Project table
- **A9.2.1** `[db]` `projects` table exists with columns: `id`, `org_id`, `name`, `status`, `created_at`.
- **A9.2.2** `[unit]` Existing functionality unaffected (no FK constraints that break existing tables).

---

## M10 — API Key Org/Project Binding & Security

### F10.1 — Token new fields
- **A10.1.1** `[db]` `tokens` table has nullable `org_id` and `project_id` columns.
- **A10.1.2** `[unit]` Existing tokens with NULL org_id/project_id continue to authenticate normally.

### F10.2 — Key hash storage
- **A10.2.1** `[unit]` `model.GetTokenByKey(rawKey)` finds the token by hashed lookup; raw key is not stored in DB.
- **A10.2.2** `[db]` No row in `tokens` table has a `key` column value starting with `sk-` in plaintext (key is hashed).

---

## M11 — usage_log & Privacy Policy

### F11.1 — Log.ProviderType field
- **A11.1.1** `[db]` `logs` table has `provider_type` varchar column.
- **A11.1.2** `[unit]` `RecordConsumeLog` populates `provider_type` from channel info.

### F11.2 — Privacy enforcement
- **A11.2.1** `[unit]` When `STORE_FULL_TEXT_ENABLED=false`, `log.Other` does not contain `"prompt"` or `"completion"` keys with full text content.
- **A11.2.2** `[unit]` `log.PromptTokens` and `log.CompletionTokens` are always recorded regardless of privacy setting.

---

## M12 — Token Statistics & Base Deduction

### F12.1 — Token counting
- **A12.1.1** `[unit]` After a successful relay, `user.UsedQuota` increases by `log.Quota`.
- **A12.1.2** `[unit]` After a successful relay, `token.UsedQuota` increases by `log.Quota`.

### F12.2 — Quota deduction
- **A12.2.1** `[integration]` After `POST /v1/chat/completions` succeeds, DB query shows `log.quota > 0` for the request.
- **A12.2.2** `[unit]` `BillingSession.Settle(actualQuota)` correctly adjusts delta between pre-consumed and actual quota.

---

## M13 — Insufficient Balance & Admin Manual Top-up

### F13.1 — Balance check
- **A13.1.1** `[integration]` Token with `remain_quota=0` and `unlimited_quota=false` receives HTTP 402 or 429 on `POST /v1/chat/completions`.
- **A13.1.2** `[unit]` `PreConsumeQuota` returns error when estimated quota exceeds available balance.

### F13.2 — Admin manual top-up
- **A13.2.1** `[curl]` `POST /api/user/topup` by admin with `{"user_id":X,"quota":100000}` returns 200; `user.Quota` increases by 100000.
- **A13.2.2** `[db]` A `logs` row with `type=3` (LogTypeManage) is created for the top-up operation.

---

## M14 — Admin Backend Provider/Channel Minimal Page

### F14.1 — Channel list filter
- **A14.1.1** `[curl]` `GET /api/channel?provider_type=official_cloud` returns only official_cloud channels.
- **A14.1.2** `[curl]` `GET /api/channel?provider_type=experimental_proxy` requires admin role; returns 403 for non-admin.

### F14.2 — experimental_proxy visibility
- **A14.2.1** `[curl]` `GET /api/channel` for non-admin user does not include channels with `provider_type=experimental_proxy`.
- **A14.2.2** `[unit]` Frontend channel list component does not render experimental_proxy rows for non-admin users.

---

## M15 — Admin Backend API Key / usage_log / Balance Pages

### F15.1 — API Key management
- **A15.1.1** `[curl]` Admin `GET /api/token` returns token list with masked keys (no full sk- value).
- **A15.1.2** `[curl]` Admin `DELETE /api/token/:id` soft-deletes the token.

### F15.2 — usage_log query
- **A15.2.1** `[curl]` `GET /api/log` returns log entries without `prompt_text` / `completion_text` fields.
- **A15.2.2** `[curl]` `GET /api/log?token_id=X` filters correctly.

### F15.3 — Balance management
- **A15.3.1** `[curl]` Admin `GET /api/user/:id` returns current `quota` and `used_quota`.
- **A15.3.2** `[curl]` Admin top-up via UI triggers `POST /api/user/topup` and reflects updated balance.

---

## M16 — Docker, Docs & Regression Tests

### F16.1 — Docker compose
- **A16.1.1** `[integration]` `docker compose -f docker-compose.dev.yml up -d` starts without error; `GET /api/status` returns 200 within 30s.

### F16.2 — Smoke test script
- **A16.2.1** `[integration]` `.factory/smoke_test.sh` exits 0 covering: models list, chat completion, balance check, experimental_proxy block.

### F16.3 — Full regression
- **A16.3.1** `[unit]` `go test ./... -count=1` exits 0 with no FAIL lines.
- **A16.3.2** `[unit]` `go vet ./...` exits 0 with no warnings.
