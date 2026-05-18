# Codex Feature Test Matrix

## Matrix Scope

This matrix maps M0-M16 features to code evidence, test evidence, and runtime validation methods. It is a plan, not a completed audit verdict.

Known source-of-truth drift:

- `features.json` lists most features as `pending`.
- `.factory/mission-state.json` lists M0-M16 complete.
- `docs/DEVELOPMENT_LOG.md` lists M0-M16 complete, but uses feature IDs that differ from `features.json`.
- `validation-contract.md` does not define every assertion later recorded in mission-state.

The audit must reconcile these before assigning final status.

## Feature-Level Test Matrix

| Milestone | Feature | Claimed implementation evidence | Planned verification |
|---|---|---|---|
| M0 | `F0.1` / repo discovery | `.factory/library/*.md`, `mission.md`, `features.json`, `validation-contract.md` | Check files exist, are non-empty, and match current call chain. Cross-check architecture against router/middleware/controller/relay/service/model code. |
| M1 | Provider type enum and Channel provider fields | `constant/channel.go`, `model/channel.go`, `model/main.go`, `constant/provider_type_test.go`, `model/channel_provider_type_test.go` | Unit tests for constants/defaults; DB migration column checks; inspect `BeforeCreate`; verify backfill is DB-compatible. |
| M1 | Risk/scope/visibility/manual-enable fields | `model/channel.go`, `constant/channel.go` | Verify defaults for official and experimental provider types; ensure explicit values are preserved; DB schema check. |
| M1 | experimental_proxy default policy | `model/channel.go`, tests in `model/channel_provider_type_test.go` | Verify experimental channels default disabled and manual-enable-required; API create path must use the hook. |
| M2 | ProviderAccount encrypted key table | `model/provider_account.go`, `common/crypto.go`, `model/provider_account_test.go` | Validate AES encryption/decryption, no plaintext key stored, CRUD behavior, migration table columns. |
| M2 | Channel ProviderAccountId link | `model/channel.go`, `model/main.go`, `model/channel_provider_type_test.go` | Verify nullable field, no FK breakage, legacy channels still route without provider account. |
| M2 | Legacy compatibility | relay/middleware/service grep, regression tests | Ensure `ProviderAccountId` is not required in relay path; old `Channel.Key` behavior still works. |
| M3 | `/v1/models` regression | `controller/model.go`, `model/channel_m3_regression_test.go`, smoke script | Unit plus local curl with mock DB/token; invalid token must be 401. |
| M3 | `/v1/chat/completions` regression | `controller/relay.go`, `relay/compatible_handler.go`, smoke script | Local mock upstream; verify 200, `choices`, `usage`, no external provider. |
| M4 | provider_type default mapping/backfill | `constant/channel.go`, `model/main.go`, `constant/provider_type_test.go` | Unit mapping checks; DB backfill scenario with old channel rows. |
| M4 | Channel CRUD provider_type | `controller/channel.go` | API tests for create/update/get; invalid provider_type rejected; response includes field. |
| M5 | ChannelModelMapping table | `model/channel_model_mapping.go`, `model/channel_model_mapping_test.go`, `model/main.go` | DB schema, enabled/disabled behavior, zero price handling, provider model override. |
| M5 | Disabled-model rejection and pricing | `middleware/distributor.go`, pricing helpers | Unit/integration tests for disabled mapping rejection and quota calculation with all provider types. |
| M6 | KiroGateway constants | `constant/channel.go`, `constant/api_type.go`, `common/api_type.go` | Verify type uniqueness, name, base URL, API mapping. |
| M6 | Kiro adapter skeleton | `relay/channel/kiro/adaptor.go`, tests | Compile-time `channel.Adaptor` assertion; every relay method returns `ErrNotImplemented`; no traffic forwarding. |
| M6 | GetAdaptor registration | `relay/relay_adaptor.go` | Unit test `relay.GetAdaptor` returns Kiro adapter; scan main flow for hardcoded Kiro references. |
| M7 | Internal user detection | `service/internal_user.go`, tests | Verify `GroupInternal`, nil context/user behavior, default users not internal. |
| M7 | Non-internal experimental block | `middleware/distributor.go`, `controller/model.go` | Integration test normal user cannot see/call experimental-only model; internal user path allowed when channel enabled. |
| M7 | experimental default disabled | `model/channel.go`, `model/channel_provider_type_test.go` | API create and model hook tests; DB status should be disabled unless explicit status provided. |
| M8 | No fallback to experimental_proxy | `model/channel_cache.go`, `service/channel_select.go`, `controller/relay.go`, `middleware/distributor.go` | Retry/fallback tests with official failure and experimental candidate present; ensure no experimental candidate when not allowed. |
| M8 | One-click disable/cache invalidation | `model/channel.go`, `controller/channel.go`, `router/api-router.go` | API test disable endpoint; cache invalidation and subsequent no-channel response. |
| M9 | Organization table | `model/organization.go`, `model/main.go`, tests | DB schema check, CRUD helpers, no business logic dependency. |
| M9 | OrganizationMember table | `model/organization.go`, tests | Unique membership intent, roles, CRUD helpers, cross-DB schema. |
| M9 | Project table | `model/organization.go`, tests | `org_id` schema, CRUD helpers, no FK side effects on existing flows. |
| M10 | Token org/project binding | `model/token.go`, `middleware/auth.go`, tests | Nullable fields, context propagation, existing token auth with nil values. |
| M10 | Key hash/prefix | `model/token.go`, tests | Verify hash lookup behavior, prefix only exposure, no full `sk-` in DB/API except allowed owner key retrieval path. |
| M10 | Token disable | `model/token.go`, `controller/token.go` | Disable helper and soft-delete/disabled behavior; auth rejects disabled token. |
| M10 | Token allow_experimental | `model/token.go`, `middleware/auth.go`, `service/channel_select.go` | Experimental access requires both internal user and token flag. |
| M11 | Log schema metadata | `model/log.go`, `model/main.go`, tests | DB columns for org/project/provider/experimental; RecordConsumeLog populates fields. |
| M11 | Success/error logs | `model/log.go`, `controller/relay.go` | Verify success consume logs and sanitized error logs; no quota deduction on error logs. |
| M11 | No prompt/response logging | `common/constants.go`, `model/log.go`, tests | Default false; `Content` cleared; `Other` does not include full text. |
| M12 | Token usage extraction | `constant/context_key.go`, `middleware/auth.go`, `middleware/distributor.go`, `model/log.go` | Verify context keys flow from token/channel into logs. |
| M12 | Cost calculation/stats | `model/log.go`, `model/log_m12_test.go` | `GetTokenUsageStats` filters by user/org/project/time and excludes error logs. |
| M12 | Balance deduction | existing billing flow | Verify success-only `PostConsumeQuota`; failure path `ReturnPreConsumedQuota`; user/token UsedQuota deltas. |
| M13 | Insufficient balance rejection | `service/pre_consume_quota.go`, `service/m13_balance_test.go` | Unit and integration tests: zero quota returns 402/429 before upstream call. |
| M13 | Admin manual recharge | existing user top-up path, `service/m13_balance_test.go` | Admin API increments quota and writes `LogTypeManage`; non-admin rejected. |
| M14 | Provider admin summary/page API | `controller/channel.go`, `model/channel.go`, `router/api-router.go` | Verify provider summary endpoint auth and counts by provider/status. |
| M14 | Provider enable/disable | `DisableExperimentalProxyChannels` | Verify only enabled experimental_proxy channels disabled; cache invalidated. |
| M14 | Channel provider filter | `controller/channel.go`, `model/channel_m14_test.go` | Verify `provider_type` query filter and non-admin experimental visibility. |
| M15 | Admin API key page/API | `controller/token.go`, `model/token.go`, routes | Verify admin route, masked keys, filters, delete/disable behavior. |
| M15 | Admin usage log page/API | `controller/log.go`, `model/log.go` | Verify provider and experimental filters; no prompt/completion text in responses. |
| M15 | Admin balance page/API | user controller/top-up path | Verify quota/used_quota visible to admin and top-up reflected. |
| M16 | Local Docker compose | `docker-compose.local.yml`, `docker-compose.dev.yml` | Local compose startup, `/api/status` health, no production secrets required. |
| M16 | Provider policy docs | `docs/provider-policy.md` | Review consistency with enforcement in code. |
| M16 | OpenAI SDK docs | `docs/openai-sdk-quickstart.md` | Review examples use local/mock-safe endpoints and documented errors. |
| M16 | Regression tests | `scripts/regression.sh`, `model/regression_m16_test.go` | Run against local mock stack; verify six core scenarios and `go test ./... -count=1`. |

## Validation Assertion Methods

| Assertion | Verification method |
|---|---|
| A1.1.1 | Inspect `model.Channel` fields and run `go test ./model/... -run TestChannelProviderType`. |
| A1.1.2 | Run AutoMigrate on SQLite/MySQL/PostgreSQL test DBs; inspect `channels` columns. MySQL/PostgreSQL blocked until local DBs available. |
| A1.1.3 | Create `Channel{ProviderType:""}` through GORM and direct hook test; expect `official_cloud`. |
| A1.2.1 | Run migration tests across all three DB drivers. Requires supplemental cross-DB test. |
| A1.2.2 | Seed old channel row, migrate, verify backfilled `provider_type='official_cloud'`. |
| A2.1.1 | Run ProviderAccount CRUD unit tests; add DB-backed CRUD if current tests are hook-only. |
| A2.1.2 | AutoMigrate and inspect `provider_accounts` columns. |
| A2.1.3 | Verify stored `EncryptedKey` is ciphertext and decrypted after find; grep/log scan for plaintext key. |
| A2.2.1 | Inspect `channels.provider_account_id` nullable column after AutoMigrate. |
| A2.2.2 | Run legacy channel tests and mock relay using `Channel.Key` with `ProviderAccountId=nil`. |
| A3.1.1 | Local curl/httptest `GET /v1/models` with valid test token; expect 200 list. |
| A3.1.2 | Local curl/httptest `GET /v1/models` with invalid token; expect 401. |
| A3.2.1 | Local mock upstream chat completion; expect 200, choices, prompt usage. |
| A3.2.2 | Run full `go test ./... -count=1`. |
| A4.1.1 | Seed channel types OpenAI/Azure/Anthropic/Gemini, migrate, verify official_cloud. |
| A4.1.2 | Run `constant` provider default tests. |
| A4.2.1 | Admin API create channel with `provider_type=aggregator`; verify persisted value. |
| A4.2.2 | API get channel by id; verify provider_type response field. |
| A5.1.1 | Audit actual selection API because `service.SelectChannel` may not exist; map assertion to `CacheGetRandomSatisfiedChannel`/distributor tests. |
| A5.1.2 | Verify default candidate pool excludes experimental_proxy. |
| A5.2.1 | Quota/pricing unit test over all four provider types. |
| A6.1.1 | Unit test constant uniqueness and sentinel ordering. |
| A6.1.2 | Unit test `ChannelTypeNames[ChannelTypeKiroGateway] == "KiroGateway"`. |
| A6.2.1 | Compile-time `channel.Adaptor` assertion for Kiro adapter. |
| A6.2.2 | Run `go build ./...`. |
| A6.3.1 | Unit test `relay.GetAdaptor` for Kiro type/API mapping. |
| A7.1.1 | Unit test internal group returns true. |
| A7.1.2 | Unit test default group returns false. |
| A7.2.1 | Integration/httptest normal user experimental-only model returns 403. |
| A7.2.2 | Integration/httptest internal user with enabled channel and token flag reaches mock upstream. |
| A7.3.1 | Unit test experimental channel defaults disabled/manual-enable. |
| A7.3.2 | API create experimental_proxy channel, inspect DB status 2. |
| A7.4.1 | Mission-state extra: reconcile definition, likely internal user explicit experimental call. Verify code/test evidence. |
| A7.4.2 | Mission-state extra: reconcile definition, likely normal user/default route isolation. Verify code/test evidence. |
| A8.1.1 | Unit test candidate selection never returns experimental when not allowed. |
| A8.1.2 | Integration retry test: official failures do not retry experimental channel. |
| A8.2.1 | Admin PATCH/disable endpoint then relay call returns no-channel. |
| A8.2.2 | Verify channel cache refresh/invalidation after status change. |
| A8.3.1 | Mission-state extra: explicit experimental candidate allowed only when user+token opt in. |
| A8.3.2 | Mission-state extra: retry param propagation from all call sites. |
| A8.4.1 | Mission-state extra: context/log tag for experimental_proxy set correctly. |
| A9.1.1 | AutoMigrate and inspect `organizations` columns. |
| A9.1.2 | Cross-DB AutoMigrate; blocked until MySQL/PostgreSQL fixtures exist. |
| A9.2.1 | AutoMigrate and inspect `projects` columns. |
| A9.2.2 | Run existing regression tests with org/project tables present. |
| A9.3.1 | Mission-state extra: OrganizationMember role/table evidence. |
| A9.3.2 | Mission-state extra: Organization/Project CRUD helper evidence. |
| A10.1.1 | AutoMigrate inspect `tokens.org_id` and `tokens.project_id`. |
| A10.1.2 | Existing token with nil org/project authenticates via httptest auth. |
| A10.2.1 | Verify raw key lookup by hash and no plaintext DB dependency. |
| A10.2.2 | DB scan for `tokens.key LIKE 'sk-%'`; expected none unless legacy compatibility explicitly stores plaintext. |
| A10.3.1 | Mission-state extra: `DisableToken` helper disables auth. |
| A10.4.1 | Mission-state extra: `AllowExperimental` defaults false and is persisted. |
| A10.4.2 | Mission-state extra: experimental access requires internal user plus token flag. |
| A11.1.1 | AutoMigrate inspect `logs.provider_type`. |
| A11.1.2 | `RecordConsumeLog` populates provider_type from context/channel. |
| A11.2.1 | With `STORE_FULL_TEXT_ENABLED=false`, inspect log content/other for no full prompt/completion. |
| A11.2.2 | Verify prompt/completion token counts remain recorded. |
| A11.3.1 | Mission-state extra: `StoreFullTextEnabled` default false. |
| A11.3.2 | Mission-state extra: error log sanitization/masking behavior. |
| A12.1.1 | Successful local relay increases `user.UsedQuota` by logged quota. |
| A12.1.2 | Successful local relay increases `token.UsedQuota` by logged quota. |
| A12.2.1 | Local chat completion writes `logs.quota > 0`. |
| A12.2.2 | Unit test BillingSession settle delta. |
| A12.3.1 | Mission-state extra: failure path returns pre-consumed quota. |
| A12.3.2 | Mission-state extra: error logs excluded from usage stats. |
| A13.1.1 | Integration zero quota token returns 402/429 before upstream call. |
| A13.1.2 | Unit test `PreConsumeQuota` insufficient balance error. |
| A13.2.1 | Admin top-up API increments user quota. |
| A13.2.2 | DB log row type manage created for top-up. |
| A14.1.1 | `GET /api/channel?provider_type=official_cloud` returns only official channels. |
| A14.1.2 | Non-admin experimental filter returns 403 or hides data according to final spec. |
| A14.2.1 | Non-admin `GET /api/channel` excludes experimental_proxy. |
| A14.2.2 | Frontend component test or API-backed UI test; currently needs supplemental frontend test evidence. |
| A14.3.1 | Mission-state extra: provider summary endpoint behavior. |
| A14.3.2 | Mission-state extra: disable experimental endpoint behavior. |
| A15.1.1 | Admin token list returns masked keys only. |
| A15.1.2 | Token delete/disable route prevents future auth. |
| A15.2.1 | Log API response omits prompt/completion text. |
| A15.2.2 | Log API `token_id` filter works. |
| A15.3.1 | Admin user detail returns quota and used_quota. |
| A15.3.2 | Contract assertion missing from mission-state: UI top-up reflects updated balance; needs explicit test. |
| A16.1.1 | Local compose startup and `/api/status` 200. Contract says dev compose; mission says local compose; test both or reconcile. |
| A16.2.1 | Run smoke/regression script against local mock stack. |
| A16.3.1 | Run `go test ./... -count=1`. |
| A16.3.2 | Run `go vet ./...`. |
| A16.4.1-A16.4.6 | Mission-state extras not in contract. Reconcile against `scripts/regression.sh` six scenarios before grading. |

## Required Runtime Evidence

The final audit should include captured output for:

- `go build ./...`
- `go test ./... -count=1`
- `go vet ./...`
- DB migration schema checks for SQLite, MySQL, PostgreSQL where available.
- Mock `/v1/models` and `/v1/chat/completions` flows.
- experimental_proxy normal/internal user matrix.
- zero-quota and admin top-up flows.
- admin channel/token/log/balance endpoint flows.
- docker compose local health check and regression script.

## Blocked Assertions Until Fixtures Exist

- All MySQL/PostgreSQL migration assertions.
- Curl assertions requiring a running app and seeded users/tokens/channels.
- Frontend visibility assertion `A14.2.2` unless frontend tests or Playwright fixtures are added.
- Any smoke scenario requiring real upstream provider keys; must be converted to mock upstream first.

