# Codex Security Findings

Audit date: 2026-05-17

Scope: feature-test-matrix retest after `AUD-001` through `AUD-004`, followed by focused `AUD-016` through `AUD-019` and `AUD-021` remediation. No real upstream provider call or real API key was used.

Security remediation update 2026-05-17: `SEC-CRIT-001` / `AUD-016`, `SEC-CRIT-002` / `AUD-017`, `SEC-HIGH-001` / `AUD-018`, `SEC-HIGH-005` / `AUD-019`, and `SEC-HIGH-002` / `AUD-021` are fixed. Remaining counts below exclude those fixed findings.

## Summary

The backend now compiles and `go test ./... -count=1` passes, so the previous build-blocker security finding is resolved. Runtime security certification is still incomplete because curl/Docker/mock-upstream fixtures were not available for every validation assertion.

| Severity | Count |
|---|---:|
| critical | 0 |
| high | 4 |
| medium | 10 |
| low | 1 |

## Full Remediation Security Update — 2026-05-17 22:05 +08:00

Backend critical/high security findings are now closed.

- `SEC-HIGH-003` experimental access control: fixed at API verification level; normal users cannot see `experimental_proxy` models, internal users require explicit `allow_experimental`, and disabled experimental remains unroutable in backend tests.
- `SEC-HIGH-004` / `AUD-022`: fixed. Token-bound organization/project scope is resolved and enforced during auth/setup; disabled or mismatched bindings fail closed; legacy NULL bindings remain user-scoped.
- `SEC-MED-002` / `AUD-023`: fixed. Admin manual top-up contract is covered through `POST /api/user/topup` request bodies with `user_id` and positive `quota`.
- `SEC-MED-006` / `AUD-024`: partially fixed. `LOCAL_FIXTURE=1 bash scripts/regression.sh` gives deterministic local smoke without real providers; fresh Docker/curl runtime remains environment-gated.
- `SEC-MED-007` / `AUD-020`: fixed. `go vet ./...` passes.
- `SEC-MED-005` / `AUD-025`: partially fixed. API-level visibility tests pass; frontend component execution remains `blocked_test_infra`.

Remaining counts after full remediation:

| Severity | Count |
|---|---:|
| critical | 0 |
| high | 0 |
| medium | 3 |
| low | 0 |

Remaining medium items are `blocked_test_infra_frontend`, `blocked_external_dependency_cross_db_runtime`, and `skipped_environment_docker_runtime`. No remaining item is a confirmed backend credential, routing, billing, or log-privacy defect.

## Pre-Release Hardening Security Update — 2026-05-17 23:30 +08:00

Historical status before CI closure: the three remaining medium items were marked `accepted_blocker`, not `passed`:

- `blocked_test_infra_frontend`: local Bun/dependency execution unavailable; CI frontend lint/test/build job added; waiver at `docs/WAIVERS/FRONTEND_TEST_INFRA_WAIVER.md`.
- `blocked_external_dependency_cross_db_runtime`: SQLite migration smoke passed locally; MySQL/PostgreSQL service proof is CI-gated; waiver at `docs/WAIVERS/CROSS_DB_MIGRATION_WAIVER.md`.
- `skipped_environment_docker_runtime`: fake-provider fixture and seed script added; runtime smoke must pass in CI/manual Docker; waiver at `docs/WAIVERS/DOCKER_RUNTIME_SMOKE_WAIVER.md`.

Security release gate remains manual review. Do not set deployment readiness to `ready` until these accepted blockers have reviewed CI/manual evidence.

## Post-CI Verification Closure — 2026-05-19

GitHub Actions `Pre-release verification` run #13 passed on branch `main` at commit `aeb43e5` in approximately 2m37s. Jobs `go-test-vet`, `local-fixture-regression`, `cross-db-migration`, `docker-fixture-smoke`, and `frontend-check` all succeeded.

Security/audit counters after CI closure:

| Metric | Count / Status |
|---|---|
| critical_findings_remaining | 0 |
| high_findings_remaining | 0 |
| medium local-environment blockers | `closed_by_ci` |
| features_failed | 0 |
| CI verification | passed |
| deployment_readiness | `staging_ready` |

Closed by CI: `blocked_test_infra_frontend`, `blocked_external_dependency_cross_db_runtime`, `skipped_environment_docker_runtime`, and `final_go_verification_blocked`.

This does not make the release `production_ready`. Recommended next action: run staging manual verification using `docs/STAGING_VERIFICATION_RUNBOOK.md`, then perform environment-variable review, real deployment topology review, and manual security sign-off.

## Fixed In Security Remediation

### SEC-CRIT-001 / AUD-016: User API Keys Persisted And Queried As Plaintext

- feature_id: `F10.2`
- status: fixed 2026-05-17
- affected_files: `model/token.go`, `controller/token.go`, `middleware/auth.go`
- evidence: `Token.Key` remains a persisted unique column; `GetTokenByKey` queries `key = ?`; cache deletion still uses raw keys; `GetTokenKey` and batch key APIs can re-display full keys.
- expected_behavior: store only `key_hash` and `key_prefix`; raw key is shown once and is not recoverable.
- actual_behavior: `KeyHash` and `KeyPrefix` are additive, but raw keys remain the active storage and lookup path.
- risk_impact: DB compromise exposes active user API keys.
- fix: new and migrated tokens compute an HMAC `key_hash` and display-only `key_prefix`; auth/cache lookup uses `key_hash`; full keys are returned only by create; disabled keys fail auth; legacy plaintext rows are migrated to hash/prefix and the deprecated `key` column is overwritten with a non-secret hash for compatibility.
- verification: `go test ./model/... -count=1` and `go test ./... -count=1` pass; new token tests cover non-plaintext persistence, one-time full-key display, correct/wrong auth, disabled auth, error-log non-leakage, and legacy plaintext no longer acting as an auth source.

### SEC-CRIT-002 / AUD-017: Disabled `experimental_proxy` Can Remain Routable Through DB Ability Selection

- feature_id: `F8.2`, `F14.2`
- status: fixed 2026-05-17
- affected_files: `model/channel.go`, `model/ability.go`, `middleware/distributor.go`
- fix: routeability filtering now requires enabled channels across DB ability selection, enabled-model queries, retry/fallback candidate pools, preferred-channel compatibility checks, and final selected-channel setup.
- verification: AUD-017 selection/fallback tests and full `go test ./... -count=1` pass.

### SEC-HIGH-001 / AUD-018: ProviderAccount Encryption Was Not The Active Relay Credential Path

- feature_id: `F2.1`
- status: fixed 2026-05-17
- affected_files: `model/provider_account.go`, `model/channel.go`, `middleware/distributor.go`, task/video relay credential paths
- fix: linked channels now resolve active credentials from decrypted ProviderAccount records at runtime. ProviderAccount credentials remain encrypted at rest, are not serialized through `json`, do not silently fallback to `Channel.Key`, and disabled/decrypt-failed accounts are rejected before upstream calls.
- verification: Added tests for encrypted-at-rest storage, active runtime credential precedence, no legacy fallback, decrypt failure rejection, disabled account rejection, credential non-leakage in errors, and legacy compatibility. `go test ./model/... -count=1` and `go test ./... -count=1` pass.

### SEC-HIGH-005 / AUD-019: Log Privacy Did Not Sanitize `Other`

- feature_id: `F11.2`
- status: fixed 2026-05-17
- affected_files: `model/log.go`, `model/log_sanitize.go`, `controller/relay.go`
- evidence: `RecordConsumeLog` blanked `Content` when full-text storage was disabled, but serialized `params.Other` and some error strings without centralized prompt/secret filtering.
- risk_impact: prompt text, response text, bearer tokens, API keys, provider credentials, authorization headers, tool payloads, or nested metadata could be persisted in logs.
- fix: added centralized sanitizers for structured `Other` payloads, JSON-string payloads, plain strings, and error messages before logger output and database persistence. Sensitive keys and bearer/API-key/credential/token patterns are redacted recursively while accounting fields remain available.
- verification: Added tests for API key, bearer authorization, prompt/messages/input, response/output, nested field, JSON-string, and error-message redaction; persisted consume/error log sanitization; normal accounting-field preservation; and default full-text/debug payload closure. `go test ./model/... -count=1`, `go test ./middleware/... -count=1`, and `go test ./... -count=1` pass.

### SEC-HIGH-002 / AUD-021: `allowed_provider_types` Was Not Enforced In Routing

- feature_id: `F5.1`, `F10.4`, `F8.1`
- status: fixed 2026-05-17
- affected_files: `model/token.go`, `middleware/auth.go`, `middleware/distributor.go`, `service/channel_select.go`, `model/ability.go`, `model/channel_cache.go`, `model/channel_satisfy.go`
- evidence: tokens exposed only `allow_experimental`; arbitrary official/aggregator/authorized_proxy allowlists were not enforced across selection, retry/fallback, DB ability, memory, preferred, or final setup paths.
- risk_impact: a token intended for one provider family could be routed to another provider family, including `experimental_proxy` through fallback.
- fix: added token `allowed_provider_types` and centralized `ProviderTypePolicy`, then applied it to selection candidates, retry/fallback, DB ability selection, memory candidates, preferred/specific channel paths, model listing, and final setup. `experimental_proxy` remains gated by both internal/admin identity and token `allow_experimental=true`.
- verification: Added AUD-021 model and middleware tests for allowlists, experimental default denial, internal explicit allow/deny, fallback/retry exclusion, DB ability exclusion, preferred/final rejection, empty provider_type inference, disabled channels, and no credential leakage. `go test ./model/... -count=1`, `go test ./middleware/... -count=1`, and `go test ./... -count=1` pass.

## Critical Open

No critical security findings remain open after the AUD-017 remediation.

### SEC-CRIT-002: Disabled `experimental_proxy` Can Remain Routable Through DB Ability Selection

- status: fixed 2026-05-17; retained below as historical baseline from the feature-test-matrix retest.
- feature_id: `F8.2`, `F14.2`
- affected_files: `model/channel.go`, `model/ability.go`, `middleware/distributor.go`
- evidence: `DisableExperimentalProxyChannels` updates only `channels.status`; DB selector chooses enabled `abilities`, then loads the channel without checking `channel.Status`; only some middleware branches check status.
- expected_behavior: disabled experimental channels are impossible to route for all users and all cache modes.
- actual_behavior: DB ability rows can still identify a disabled channel when memory cache is off or stale.
- risk_impact: the one-click experimental kill switch may not actually stop traffic.
- recommended_fix: update related ability rows in the disable transaction and add a final enabled-status check after every channel load.

## High

### SEC-HIGH-001: ProviderAccount Encryption Is Not The Active Relay Credential Path

- status: fixed 2026-05-17; retained below as historical baseline from the feature-test-matrix retest.
- feature_id: `F2.1`
- affected_files: `model/provider_account.go`, `model/channel.go`, `middleware/distributor.go`, `controller/channel.go`
- evidence: ProviderAccount encrypts `EncryptedKey`, but relay setup still reads `Channel.Key`; channel APIs can still expose channel key material.
- risk_impact: upstream provider credentials remain plaintext in the active path.

### SEC-HIGH-002: `allowed_provider_types` / Provider-Type Restriction Is Not Implemented

- status: fixed 2026-05-17; retained below as historical baseline from the feature-test-matrix retest.
- feature_id: `F5.1`, `F10.4`
- affected_files: `model/token.go`, `middleware/auth.go`, `service/channel_select.go`
- evidence: there is no token field or selector enforcement for allowed provider types; only `AllowExperimental` exists.
- risk_impact: a key intended for official providers can be routed to aggregators or authorized proxies.

### SEC-HIGH-003: Experimental Access Control Is Incomplete And Not Runtime-Proven

- feature_id: `F7.2`, `F7.3`, `F8.1`
- affected_files: `middleware/distributor.go`, `controller/relay.go`, `service/channel_select.go`
- evidence: normal selection uses internal-user plus token flag, but post-selection guard checks only internal user; required error code `experimental_proxy_forbidden` is not emitted; integration matrix was not run.
- risk_impact: edge paths can diverge from the documented double gate.

### SEC-HIGH-004: API Key Org/Project Binding Is Schema/Context Only

- feature_id: `F10.1`
- affected_files: `model/token.go`, `controller/token.go`, `middleware/auth.go`
- evidence: `OrgId`/`ProjectId` fields exist and auth copies them to context, but create/update paths do not persist them and no membership boundary is enforced.
- risk_impact: tenant/project isolation can be assumed by logs without being enforced.

### SEC-HIGH-005: Log Privacy Does Not Sanitize `Other`

- status: fixed 2026-05-17; retained below as historical baseline from the feature-test-matrix retest.
- feature_id: `F11.2`
- affected_files: `model/log.go`, relay logging callers
- evidence: `RecordConsumeLog` blanks `Content` when full-text storage is disabled, but serializes `params.Other` without filtering prompt/completion keys.
- risk_impact: prompt or response text can still be written through structured side payloads.

### SEC-HIGH-006: Billing Deductions Are Not Atomic Under Concurrency

- feature_id: `F12.1`, `F12.2`, `F13.1`
- affected_files: `service/billing_session.go`, `service/quota.go`, `model/user.go`, `model/token.go`
- evidence: quota checks and decrements are separate operations; DB updates do not condition on sufficient remaining quota.
- risk_impact: concurrent requests can overdraw user or token quota.

### SEC-HIGH-007: Upstream Task/Async Paths Can Persist Sensitive Credential Data

- feature_id: `F2.1`, `F11.2`
- affected_files: `model/task.go`, `controller/video_proxy.go`, `controller/video_proxy_gemini.go`
- evidence: task private data stores provider keys for some video/Gemini flows.
- risk_impact: task table compromise can expose upstream credentials.

## Medium

### SEC-MED-001: Cross-DB Migration Assertions Are Not Proven

- feature_id: `F1.1`, `F1.2`, `F9.1`, `F9.2`, `F10.1`
- evidence: local Go tests pass, but MySQL/PostgreSQL migration evidence was not produced.

### SEC-MED-002: Admin Top-Up Contract Does Not Match Implementation

- feature_id: `F13.2`, `F15.3`
- evidence: validation expects `POST /api/user/topup`; implementation uses `/api/user/manage` with `action=add_quota`.

### SEC-MED-003: Insufficient Quota HTTP Status Contract Does Not Match Implementation

- feature_id: `F13.1`
- evidence: validation expects HTTP 402 or 429; implementation returns HTTP 403 for insufficient user/token quota paths.

### SEC-MED-004: Admin Token Management Is Incomplete

- feature_id: `F15.1`
- evidence: admin list exists, but no admin delete route matching the validation assertion exists.

### SEC-MED-005: Frontend Security Visibility Tests Are Unavailable

- feature_id: `F14.2`
- evidence: `bun` is unavailable and no frontend `test` scripts exist; non-admin experimental row visibility is not component-tested.

### SEC-MED-006: Smoke And Runtime Tests Lack Deterministic Fixtures

- feature_id: `F3.1`, `F3.2`, `F16.2`
- evidence: smoke scripts require a running gateway, admin token, `jq`, seeded users/channels, and an upstream responder.

### SEC-MED-007: `go vet` Regression Gate Fails

- feature_id: `F16.3`
- evidence: `go vet ./...` reports mutex-copy warnings and unreachable-code warnings.

### SEC-MED-008: Docker Compose Uses Active Default Secrets

- feature_id: `F16.1`
- evidence: compose config includes Postgres password `123456` and local `SESSION_SECRET=change-me-in-production`.

### SEC-MED-009: Validation Contract References Nonexistent Or Mismatched APIs

- feature_id: `global-validation-contract`
- evidence: `service.SelectChannel`, `service.GetQuotaByChannel`, `NewExperimentalProxyChannel`, `/api/user/topup`, and `relay.GetAdaptor(ChannelTypeKiroGateway)` do not match current code.

### SEC-MED-010: Provider Account Audit Identity Is Missing From Usage Logs

- feature_id: `F8.4`, `F11.1`
- evidence: logs record channel id/provider type/experimental flag, but not `provider_account_id`.

## Low

### SEC-LOW-001: Mission State And Feature Source Drift

- feature_id: `F0.1`
- evidence: `features.json` marks 36 of 37 features pending while development log and mission-state mark the expanded mission complete.
