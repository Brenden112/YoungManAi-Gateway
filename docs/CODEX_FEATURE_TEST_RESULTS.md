# Codex Feature Test Results

Generated: 2026-05-17 14:27 +08:00

Scope: M0-M16 feature-test-matrix retest based on `features.json`, `validation-contract.md`, `docs/DEVELOPMENT_LOG.md`, and `.factory/mission-state.json`. The original matrix made no business-code changes and used no real upstream provider key.

Security remediation update 2026-05-17: `AUD-016` fixed `F10.2`, `AUD-017` fixed disabled experimental routing, `AUD-018` fixed ProviderAccount active credentials, `AUD-019` fixed log privacy, and `AUD-021` fixed provider-type policy enforcement.

## Baseline Commands

| Check | Command | Result |
|---|---|---|
| Full Go tests | `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./... -count=1` | PASS |
| Targeted model tests | `go test ./model -run 'Test(Channel|ProviderAccount|Organization|Token|Log|Regression|GetGroup|DisableExperimental|StoreFullText)' -count=1` | PASS |
| Targeted controller tests | `go test ./controller -run 'Test(GetAllChannels|Channel|Provider|DisableExperimental|GetAllTokens|SearchTokens|Model)' -count=1` | PASS |
| Targeted service tests | `go test ./service -run 'Test(Internal|M13|Billing|PreConsume|Tiered|TextQuota)' -count=1` | PASS |
| Kiro skeleton tests | `go test ./relay/channel/kiro ./relay -count=1` | PASS |
| Go vet | `go vet ./...` | FAIL: mutex-copy warnings and multiple unreachable-code warnings |
| Docker dev config | `docker compose -f docker-compose.dev.yml config` | PASS |
| Docker local config | `docker compose -f docker-compose.local.yml config` | PASS |
| Frontend lint/test | `bun --version`; `bun run lint` | SKIPPED: `bun` is not installed; neither frontend package defines a `test` script |

## Summary

| Metric | Count |
|---|---:|
| Features tested from `features.json` | 37 |
| Pass | 11 |
| Fail | 11 |
| Blocked | 15 |
| Validation assertions checked | 70 |

Important distinction: `go test ./...` now passes, but many feature assertions are curl, cross-database, Docker runtime, frontend, or security-behavior assertions that are not proven by the Go unit suite.

## Per-Feature Matrix

| Feature | Status | Evidence | Blocking / Failure Reason |
|---|---|---|---|
| F0.1 M0 repo discovery | pass | `.factory/library/*.md` exists; development log contains relay call chain | No validation-contract assertions for M0 |
| F1.1 Channel model fields | blocked | `model.Channel` has new fields; model tests pass | A1.1.2 requires SQLite, MySQL, and PostgreSQL schema evidence; only local Go/SQLite-style evidence was available |
| F1.2 DB migration compatibility | blocked | `model/main.go` contains backfill logic | A1.2.1/A1.2.2 require all three DB drivers and migrated existing rows |
| F2.1 ProviderAccount model | pass | encryption hooks, DB-backed provider-account tests, and active credential resolver tests pass | ProviderAccount credentials are encrypted at rest and decrypted only for runtime relay context |
| F2.2 Channel ProviderAccountId | pass | `ProviderAccountId *int` exists; nil compatibility and middleware setup tests pass | Linked channels use ProviderAccount credentials; nil ProviderAccountId preserves legacy channel key compatibility |
| F3.1 GET `/v1/models` | blocked | controller tests cover model-list filtering paths | A3.1.1/A3.1.2 require a running server and seeded valid/invalid tokens |
| F3.2 POST `/v1/chat/completions` | blocked | relay path compiles; no real upstream used | A3.2.1 requires mock/stub upstream fixture; none is available |
| F4.1 default provider type | blocked | `GetDefaultProviderType` tests pass | A4.1.1 requires migrated DB rows for OpenAI/Azure/Anthropic/Gemini |
| F4.2 Channel CRUD provider_type | blocked | controller builds; validation code exists | A4.2.1/A4.2.2 are curl assertions requiring a running admin API fixture |
| F5.1 provider-type-aware selection | pass | AUD-021 tests cover token `allowed_provider_types`, DB ability selection, memory candidates, retry/fallback, preferred/final setup, model listing policy, and experimental default denial | Runtime curl fixture still remains separate from the unit/integration coverage |
| F5.2 pricing for provider types | pass | service quota/pricing tests pass; provider_type is not used as a panic-prone input | No provider_type-specific panic observed in local tests |
| F6.1 KiroGateway constant | pass | `constant.ChannelTypeKiroGateway == 58`; names/default-provider tests pass | Assertion covered by unit tests |
| F6.2 KiroGateway adaptor skeleton | pass | `relay/channel/kiro` tests pass; full Go tests pass | Skeleton rejects relay methods by design |
| F6.3 GetAdaptor registration | fail | `GetAdaptor(APITypeKiroGateway)` is wired | A6.3.1 says `GetAdaptor(ChannelTypeKiroGateway)`, but `GetAdaptor` accepts API type, not channel type |
| F7.1 internal user helper | pass | `service/internal_user_test.go` passes | Contract names `isInternalUser(user)`, implementation is `service.IsInternalUser(*gin.Context)`; behavior is still covered |
| F7.2 block non-internal experimental | fail | post-selection guard blocks non-internal users | A7.2.1 requires error code `experimental_proxy_forbidden`; code returns plain message. No runtime integration evidence |
| F7.3 experimental default disabled | fail | `BeforeCreate` defaults status disabled when status is zero | Code preserves explicit enabled status for experimental channels; API-created default-disabled curl evidence is absent |
| F8.1 no fallback to experimental | blocked | in-memory selector tests exclude experimental when not allowed | A8.1.2 retry-loop integration with failing official channels was not run |
| F8.2 one-click disable | fail | bulk disable helper exists | DB selection path can load a disabled channel via enabled abilities because final DB load does not check channel status |
| F9.1 Organization table | blocked | struct tests pass | A9.1.1/A9.1.2 require DB schema and all-three-driver AutoMigrate evidence |
| F9.2 Project table | blocked | struct tests pass | A9.2.1/A9.2.2 require DB schema/migration evidence beyond local struct tests |
| F10.1 token org/project fields | fail | fields exist; auth copies context when present | token create/update paths do not persist org/project bindings; no DB-column/auth compatibility runtime evidence |
| F10.2 API key hash/prefix only | pass | `Token.BeforeCreate`, `Insert`, auth lookup, cache lookup, migration, and controller tests now use HMAC `key_hash` plus `key_prefix`; full key is returned only on create | Legacy `tokens.key` column remains as deprecated non-secret compatibility storage containing the hash, not plaintext |
| F11.1 usage_log provider_type | pass | `RecordConsumeLog` writes provider type from context; log tests pass | SQLite/local unit evidence is sufficient for this schema-level assertion |
| F11.2 no prompt/response logging | pass | `Content` is blank when `StoreFullTextEnabled=false`; AUD-019 sanitizer tests cover `params.Other`, nested maps, JSON strings, persisted consume logs, persisted error logs, and error-message redaction | Full-text/debug payload storage is default off through `StoreFullTextEnabled=false`; sanitized accounting metadata remains available |
| F12.1 token usage increment | blocked | quota settlement code compiles and tests pass | assertion requires a successful relay with log quota and DB usage increments |
| F12.2 quota deduction | blocked | `BillingSession.Settle` tests pass | A12.2.1 requires successful chat completion and DB log evidence |
| F13.1 insufficient balance | fail | unit path returns insufficient-quota errors before relay | HTTP status is 403 in code, not 402/429 as validation contract requires |
| F13.2 admin manual top-up | fail | `POST /api/user/manage` action `add_quota` exists and logs `LogTypeManage` | validation requires `POST /api/user/topup`, which is not the implemented admin endpoint |
| F14.1 channel provider filter | blocked | backend filter code and admin route exist | curl/admin/non-admin fixture was not run |
| F14.2 non-admin experimental visibility | fail | relay model list filters experimental models | `/api/channel` is admin-only, and frontend component test evidence is absent |
| F15.1 admin API key management | fail | admin token list exists and masking tests pass | no admin `DELETE /api/token/:id`; user `GetTokenKey` can re-display the full key |
| F15.2 usage_log query | blocked | log filter unit tests pass | curl assertions require running API fixture |
| F15.3 admin balance management | fail | admin `GET /api/user/:id` exists | UI/top-up assertion expects `POST /api/user/topup`; implementation uses manage action |
| F16.1 Docker compose runtime | blocked | compose config validates | `docker compose up` plus `/api/status` runtime check was not run |
| F16.2 smoke script | blocked | `.factory/smoke_test.sh` and `scripts/regression.sh` exist | requires running gateway, admin token, and deterministic fixture; no real keys used |
| F16.3 full regression | fail | `go test ./... -count=1` passes | `go vet ./...` fails |

## Validation Assertions

| Assertion | Status | Evidence |
|---|---|---|
| A1.1.1 | pass | Channel fields exist; model tests pass |
| A1.1.2 | blocked | all-three-DB schema check not run |
| A1.1.3 | pass | default provider type tests pass |
| A1.2.1 | blocked | all-three-driver migration test not run |
| A1.2.2 | blocked | migrated existing-row backfill not verified |
| A2.1.1 | fail | ProviderAccount CRUD function coverage is incomplete |
| A2.1.2 | blocked | DB table schema not independently inspected |
| A2.1.3 | pass | AES encrypt/decrypt hook tests pass |
| A2.2.1 | blocked | DB column schema not independently inspected |
| A2.2.2 | pass | nil ProviderAccountId compatibility tests pass |
| A3.1.1 | blocked | requires running server/valid token |
| A3.1.2 | blocked | requires running server/invalid token |
| A3.2.1 | blocked | requires mock upstream chat fixture |
| A3.2.2 | pass | `go test ./... -count=1` passes |
| A4.1.1 | blocked | migrated DB row check not run |
| A4.1.2 | pass | constant tests pass |
| A4.2.1 | blocked | curl admin API fixture not run |
| A4.2.2 | blocked | curl admin API fixture not run |
| A5.1.1 | pass | `ProviderTypePolicy` is enforced by selection, fallback/retry, DB ability, memory candidate, preferred/specific channel, model list, and final setup paths |
| A5.1.2 | pass | selector tests exclude experimental by default |
| A5.2.1 | pass | service tests pass without provider_type panic |
| A6.1.1 | pass | constant tests pass |
| A6.1.2 | pass | constant tests pass |
| A6.2.1 | pass | Kiro adaptor tests compile and pass |
| A6.2.2 | pass | full Go tests pass |
| A6.3.1 | fail | factory is keyed by API type, not channel type |
| A7.1.1 | pass | service internal-user tests pass |
| A7.1.2 | pass | service internal-user tests pass |
| A7.2.1 | fail | no required `experimental_proxy_forbidden` code; no runtime fixture |
| A7.2.2 | blocked | requires internal user plus enabled experimental runtime fixture |
| A7.3.1 | fail | no `NewExperimentalProxyChannel`; explicit enabled status is preserved |
| A7.3.2 | blocked | API-created channel DB status not tested |
| A8.1.1 | pass | selector unit tests pass |
| A8.1.2 | blocked | retry-loop runtime fixture not run |
| A8.2.1 | blocked | curl disable flow not run |
| A8.2.2 | fail | helper updates channel status only; ability status/cache invalidation evidence incomplete |
| A9.1.1 | blocked | DB table schema not inspected |
| A9.1.2 | blocked | all-three-DB AutoMigrate not run |
| A9.2.1 | blocked | DB table schema not inspected |
| A9.2.2 | blocked | migration compatibility not proven |
| A10.1.1 | blocked | DB column schema not independently inspected |
| A10.1.2 | blocked | auth compatibility runtime not independently proven |
| A10.2.1 | pass | `GetTokenByKey` computes HMAC and queries `key_hash`; wrong keys and legacy plaintext-only values fail auth |
| A10.2.2 | pass | create/insert/migration tests verify plaintext is not persisted and only `key_hash`/`key_prefix` carry credential material |
| A11.1.1 | pass | Log schema field exists and tests pass |
| A11.1.2 | pass | `RecordConsumeLog` populates provider type from context |
| A11.2.1 | pass | AUD-019 sanitizer redacts prompt/messages/input, response/output, bearer tokens, API keys, credentials, secrets, headers, nested fields, and JSON-string payloads before usage-log persistence |
| A11.2.2 | pass | prompt/completion token fields are recorded in consume log |
| A12.1.1 | blocked | requires successful relay DB evidence |
| A12.1.2 | blocked | requires successful relay DB evidence |
| A12.2.1 | blocked | requires successful chat completion fixture |
| A12.2.2 | pass | `BillingSession.Settle` tests pass |
| A13.1.1 | fail | code returns HTTP 403 for insufficient quota |
| A13.1.2 | pass | pre-consume unit path returns error |
| A13.2.1 | fail | expected `/api/user/topup` route is absent |
| A13.2.2 | pass | manage add_quota writes `LogTypeManage` |
| A14.1.1 | blocked | curl admin API fixture not run |
| A14.1.2 | blocked | curl non-admin API fixture not run |
| A14.2.1 | blocked | curl non-admin API fixture not run |
| A14.2.2 | fail | frontend component test/script unavailable |
| A15.1.1 | pass | admin token list and masking tests pass |
| A15.1.2 | fail | no admin delete route matching assertion |
| A15.2.1 | blocked | curl log API fixture not run |
| A15.2.2 | blocked | curl log API fixture not run |
| A15.3.1 | blocked | curl user API fixture not run |
| A15.3.2 | fail | UI/topup route assertion is not implemented |
| A16.1.1 | blocked | compose runtime not started |
| A16.2.1 | blocked | smoke script not run; requires fixture/admin token |
| A16.3.1 | pass | full Go tests pass |
| A16.3.2 | fail | `go vet ./...` fails |

## Documentation Consistency Findings

- `features.json` still has 36 of 37 features as `pending`, while `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json` mark the expanded 53-feature mission set completed.
- `validation-contract.md` references APIs that do not match code: `service.SelectChannel`, `service.GetQuotaByChannel`, `NewExperimentalProxyChannel`, `/api/user/topup`, admin `DELETE /api/token/:id`, and `relay.GetAdaptor(ChannelTypeKiroGateway)`.
- `scripts/regression.sh` depends on a running gateway, `ADMIN_TOKEN`, `jq`, and seeded channels/users; it does not provide a deterministic local mock upstream fixture.

## Top 10 Issues From Matrix

1. Admin/manual top-up validation contract uses `/api/user/topup`, but implementation uses `/api/user/manage` with action `add_quota`.
2. `go vet ./...` fails, so F16 full regression is not satisfied.
3. Cross-DB migration assertions for SQLite/MySQL/PostgreSQL remain unverified.
4. Curl/runtime assertions for `/v1/models`, chat completion, admin channel/token/log APIs, and Docker health are blocked by missing seeded fixture/server run.
5. Frontend lint/test and non-admin experimental visibility component tests are unavailable in the current environment.
6. Runtime proof is still needed for several blocked assertions after the AUD-016/AUD-017/AUD-018/AUD-019/AUD-021 unit-level fixes.

## Full Remediation Retest — 2026-05-17 22:05 +08:00

| Check | Command | Result |
|---|---|---|
| Targeted model tests | `go test ./model/... -count=1` | PASS |
| Targeted middleware tests | `go test ./middleware/... -count=1` | PASS |
| Full Go tests | `go test ./... -count=1` | PASS |
| Go vet | `go vet ./...` | PASS |
| Local fixture smoke | `LOCAL_FIXTURE=1 bash scripts/regression.sh` | PASS |
| Docker compose config | `docker compose config` | PASS |
| Mission-state JSON | `node -e "JSON.parse(...)"` | PASS |
| Frontend | `bun --version` | BLOCKED: `bun` not installed; no frontend `test` scripts |

Updated feature matrix summary:

| Metric | Count |
|---|---:|
| Features retested | 37 |
| Pass | 29 |
| Fail | 0 |
| Blocked | 8 |
| Validation assertions checked | 70 |

Newly closed assertions include `A10.1.2`, `A11.1.2`, `A13.2.1`, `A13.2.2`, `A14.2.1`, local-fixture evidence for `A16.2.1`, and `A16.3.2`.

Remaining blocked assertions are environment or test-infrastructure scoped: all-three-DB migration proof, fresh seeded Docker/curl runtime checks, and frontend component-level `experimental_proxy` visibility.

## Pre-Release Hardening Retest — 2026-05-17 23:30 +08:00

| Check | Command | Result |
|---|---|---|
| SQLite migration smoke | `bash scripts/check-migrations.sh` | PASS |
| Fixture compose config | `docker compose -f docker-compose.fixture.yml config` | PASS |
| Frontend script presence | inspect `web/default/package.json` | PASS: `lint`, `test`, and `build` scripts exist |
| Frontend local execution | `bun --version` | ACCEPTED_BLOCKER: Bun absent locally |
| MySQL/PostgreSQL migration runtime | CI services | ACCEPTED_BLOCKER locally; CI workflow added |
| Docker fixture runtime | `docker compose -f docker-compose.fixture.yml up -d --build` | ACCEPTED_BLOCKER locally/manual; CI workflow added |

Accepted blocker waivers:

- `docs/WAIVERS/FRONTEND_TEST_INFRA_WAIVER.md`
- `docs/WAIVERS/CROSS_DB_MIGRATION_WAIVER.md`
- `docs/WAIVERS/DOCKER_RUNTIME_SMOKE_WAIVER.md`

Historical pre-CI status: deployment readiness remained `needs_manual_review` until CI evidence was available.

## Post-CI Verification Closure — 2026-05-19

GitHub Actions `Pre-release verification` run #13 passed on branch `main` at commit `aeb43e5` in approximately 2m37s.

| Check | CI job | Result |
|---|---|---|
| Go test/vet | `go-test-vet` | PASS |
| Local fixture regression | `local-fixture-regression` | PASS |
| Cross-DB migration | `cross-db-migration` | PASS |
| Docker fixture smoke | `docker-fixture-smoke` | PASS |
| Frontend lint/test/build | `frontend-check` | PASS |

Updated feature/audit summary:

| Metric | Count / Status |
|---|---|
| Features failed | 0 |
| Critical findings remaining | 0 |
| High findings remaining | 0 |
| Local-environment medium blockers | `closed_by_ci` |
| CI verification | passed |
| Deployment readiness | `staging_ready` |

Closed by CI: `blocked_test_infra_frontend`, `blocked_external_dependency_cross_db_runtime`, `skipped_environment_docker_runtime`, and `final_go_verification_blocked`.

This is not production readiness. Recommended next action: run staging manual verification using `docs/STAGING_VERIFICATION_RUNBOOK.md`.
