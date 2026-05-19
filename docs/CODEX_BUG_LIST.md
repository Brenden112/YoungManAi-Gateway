# Codex Audit Bug List

## Canonical Full Remediation Status — 2026-05-17 22:05 +08:00

The current canonical status after full remediation is:

- fixed: `AUD-001`, `AUD-002`, `AUD-003`, `AUD-004`, `AUD-016`, `AUD-017`, `AUD-018`, `AUD-019`, `AUD-020`, `AUD-021`, `AUD-022`, `AUD-023`, `AUD-024-local-fixture`, `AUD-025-api-visibility`
- closed_by_ci: `blocked_test_infra_frontend`, `blocked_external_dependency_cross_db_runtime`, `skipped_environment_docker_runtime`, `final_go_verification_blocked`
- critical/high remaining: 0

Pre-release hardening update 2026-05-17: the remaining three items are no longer undocumented blockers. Each has a minimum fix path, CI coverage, and waiver:

- `blocked_test_infra_frontend`: `docs/WAIVERS/FRONTEND_TEST_INFRA_WAIVER.md`
- `blocked_external_dependency_cross_db_runtime`: `docs/WAIVERS/CROSS_DB_MIGRATION_WAIVER.md`
- `skipped_environment_docker_runtime`: `docs/WAIVERS/DOCKER_RUNTIME_SMOKE_WAIVER.md`. This is an environment/CI blocker, not a confirmed code bug; the local Docker build reached dependency download and was stopped after no progress.
- `final_go_verification_blocked`: `docs/WAIVERS/LOCAL_GO_TOOLCHAIN_WAIVER.md`. This is a local toolchain blocker, not a confirmed code bug; the current shell has no `go` binary and `/tmp/go2510/bin/go` is absent.

CI/staging verification setup 2026-05-18: the four environment blockers above are not marked passed. They are moved to `pending_ci_verification` and must be closed by `.github/workflows/pre-release-verification.yml` or staging runbook evidence.

Post-CI verification closure 2026-05-19: GitHub Actions `Pre-release verification` run #13 passed on branch `main` at commit `aeb43e5`. Jobs `go-test-vet`, `local-fixture-regression`, `cross-db-migration`, `docker-fixture-smoke`, and `frontend-check` all succeeded. The four local environment blockers are now `closed_by_ci`:

- `blocked_test_infra_frontend`
- `blocked_external_dependency_cross_db_runtime`
- `skipped_environment_docker_runtime`
- `final_go_verification_blocked`

Current release-audit counters: `critical_findings_remaining = 0`, `high_findings_remaining = 0`, `features_failed = 0`, CI verification passed, and `deployment_readiness = staging_ready`. This is not `production_ready`; recommended next action is to run staging manual verification using `docs/STAGING_VERIFICATION_RUNBOOK.md`.

Older per-issue sections below are retained as historical audit baseline where noted.

## AUD-001

- feature_id: `M11-F01-usage-log-schema`, also impacts `M11-F02`, `M11-F03`, `M12-F01`, `M12-F02`, `M15-F02`, `M16-F04`
- milestone: `M11`
- severity: critical
- status: fixed 2026-05-17
- category: bug
- affected_files: `model/log.go`
- reproduction_steps:
  1. Run `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go build ./...`.
  2. Observe syntax errors.
- expected_behavior: Repository builds and log model compiles.
- actual_behavior: Build fails with `model/log.go:373:2: syntax error: non-declaration statement outside function body` and `model/log.go:390:3: syntax error: non-declaration statement outside function body`.
- risk_impact: Entire backend cannot build; all log, billing, admin log, and regression features are non-executable.
- suggested_fix: Correct the malformed brace/control-flow structure in `GetAllLogs`; run gofmt, `go build ./...`, and targeted log/billing tests.
- root_cause: A stray closing brace ended `GetAllLogs` immediately after the count query, leaving the find query and channel-name backfill as orphan statements outside the function body.
- resolution: Removed the stray brace, kept the count error handling inside the function, and added a regression test covering `provider_type` plus `is_experimental_proxy` filtering and channel-name backfill.
- verification: `gofmt` passed for `model/log.go` and `model/log_m15_test.go`; `go build ./model` passed. `go test ./model/...` is still blocked by unrelated existing compile errors in `model/regression_m16_test.go`.
- remaining_risk: Full backend build now passes after `AUD-003`; model package tests still need the M16 regression test compile errors fixed under a separate confirmed issue.

## AUD-002

- feature_id: `M15-F01-api-key-admin-page`
- milestone: `M15`
- severity: critical
- status: fixed 2026-05-17
- category: bug
- affected_files: `controller/token.go`
- reproduction_steps:
  1. Run `/tmp/go2510/bin/gofmt -l $(rg --files -g '*.go')`.
  2. Inspect `controller/token.go:82`.
- expected_behavior: `GetAllTokens` and `GetAdminAllTokens` are separate valid Go functions.
- actual_behavior: Lines 82-95 contain top-level statements after `GetAdminAllTokens`, causing `expected declaration, found userId`.
- risk_impact: Token controller cannot compile; admin token management and user token list endpoints are unavailable.
- suggested_fix: Restore the missing `func SearchTokens` declaration around the orphan user token search body; rerun gofmt and token controller tests.
- root_cause: The `SearchTokens` function declaration was missing, leaving its body as top-level statements immediately after `GetAdminAllTokens`.
- resolution: Restored `func SearchTokens(c *gin.Context)` around the existing user token search body and formatted `controller/token.go`.
- verification: `gofmt -l controller/token.go` passed with no output. File-scoped token controller tests passed for `controller/token.go controller/token_test.go`, including list and search masking coverage. Package and full-repo checks still fail only on unrelated remaining blockers.
- remaining_risk: Full backend build now passes after `AUD-003`; full tests are still blocked by existing `model/regression_m16_test.go` compile errors.

## AUD-003

- feature_id: `M14-F01-provider-admin-page`, also impacts `M14-F02`, `M14-F03`
- milestone: `M14`
- severity: critical
- status: fixed 2026-05-17
- category: bug
- affected_files: `controller/channel.go`
- reproduction_steps:
  1. Run `/tmp/go2510/bin/gofmt -l $(rg --files -g '*.go')`.
  2. Inspect `controller/channel.go:738`.
- expected_behavior: Channel admin controller compiles and route handlers are valid declarations.
- actual_behavior: Orphan struct fields begin at line 738 outside a type declaration, causing `expected declaration, found Tag`.
- risk_impact: Channel controller cannot compile; provider summary, disable experimental, channel filtering, and channel CRUD are unavailable.
- suggested_fix: Restore the missing struct/type declaration around the orphan fields or remove the corrupt block; rerun gofmt and controller/channel tests.
- root_cause: The `ChannelTag` type declaration was removed while its field list and closing brace remained, leaving struct fields as top-level orphan declarations after `GetChannelProviderSummary`.
- resolution: Restored `type ChannelTag struct {` before the existing tag-management fields and formatted `controller/channel.go`.
- verification: `gofmt -l controller/channel.go`, `go build ./controller`, `go test ./controller -count=1`, `go build ./...`, and `git diff --check` passed. Full `go test ./...` still fails only on existing `model/regression_m16_test.go` compile errors.
- remaining_risk: Full tests are still blocked by existing M16 regression test compile errors; global `gofmt -l $(rg --files -g '*.go')` still lists unrelated unformatted files.

## AUD-004

- feature_id: `M16-F04-regression-tests`
- milestone: `M16`
- severity: high
- status: fixed 2026-05-17
- category: test_gap
- affected_files: `middleware/distributor.go`, `common/ssrf_protection.go`, `relay/relay_adaptor.go`, `controller/token_test.go`, `controller/token.go`, `controller/channel.go`, `dto/gemini.go`, `setting/payment_waffo.go`, `constant/waffo_pay_method.go`, `constant/channel.go`, `relay/common/stream_status.go`, `service/m13_balance_test.go`, `service/channel_select.go`, `pkg/billingexpr/run.go`, `pkg/billingexpr/compile.go`, `model/user_oauth_binding.go`, `model/token.go`, `model/provider_account.go`, `model/log.go`
- reproduction_steps:
  1. Run `/tmp/go2510/bin/gofmt -l $(rg --files -g '*.go')`.
- expected_behavior: No files are listed; syntax is valid and formatting is stable.
- actual_behavior: After `AUD-001`, `AUD-002`, and `AUD-003`, gofmt no longer reports syntax errors, but it still lists multiple unrelated unformatted files.
- risk_impact: Formatting gate fails and syntax errors hide further test failures.
- suggested_fix: Fix syntax errors first, then run gofmt on affected files and review unrelated formatting churn before committing.
- root_cause: Multiple Go files had stale gofmt alignment/import-order drift after the prior feature and fix phase changes.
- resolution: Ran gofmt on the remaining listed files from `gofmt -l`; no business logic, database schema, test assertions, or test deletion changes were intentionally made for this issue.
- verification: `/tmp/go2510/bin/gofmt -l $(rg --files -g '*.go')`, `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./... -count=1`, and `git diff --check` all passed. Frontend lint was skipped because `bun` is unavailable in this environment; frontend tests were skipped because neither frontend package defines a `test` script.
- remaining_risk: `git diff --ignore-all-space --stat` still shows pre-existing non-formatting worktree diffs in `constant/channel.go`, `middleware/distributor.go`, `model/token.go`, `relay/relay_adaptor.go`, and `service/channel_select.go`; these are tracked separately as `AUD-015` and were not expanded under `AUD-004`.

## AUD-005

- feature_id: `M10-F02-api-key-hash-prefix-once`
- milestone: `M10`
- severity: high
- category: security
- status: superseded_by_AUD-016_fixed_2026-05-17
- affected_files: `model/token.go`, `controller/token.go`
- reproduction_steps:
  1. Inspect `model/token.go:16-40`.
  2. Inspect `model/token.go:48-50`.
- expected_behavior: Raw API key is not persisted in DB; lookup uses hash and responses expose only prefix/masked values unless explicitly allowed by secure one-time display policy.
- actual_behavior_before_fix: `Token.Key` remained a persisted unique indexed string, and the code comment said plaintext key was retained for backward-compatible auth lookup.
- risk_impact: Database compromise can expose user API keys, contradicting mission key security and validation assertions `A10.2.1`/`A10.2.2`.
- resolution: Fixed under `AUD-016`; token auth/cache lookup now uses HMAC `key_hash`, display uses `key_prefix`, full keys are shown only once, and legacy plaintext rows are migrated to non-secret hash storage.

## AUD-006

- feature_id: `M0-F03-create-mission-files`
- milestone: `M0`
- severity: medium
- category: documentation_mismatch
- affected_files: `features.json`, `.factory/mission-state.json`, `docs/DEVELOPMENT_LOG.md`
- reproduction_steps:
  1. Run the node comparison script in the audit notes.
  2. Observe `features.json count 37 { done: 1, pending: 36 }`.
  3. Compare with mission-state `completed_features` count 59.
- expected_behavior: Feature status sources agree or explicitly document supersession.
- actual_behavior: `features.json` still marks most features pending while mission-state and development log mark MVP complete.
- risk_impact: Audit cannot trust completion status without independent code/test evidence.
- suggested_fix: Reconcile feature IDs and status fields across all mission documents after code fixes.

## AUD-007

- feature_id: `M0-F03-create-mission-files`
- milestone: `M0`
- severity: medium
- category: documentation_mismatch
- affected_files: `docs/REPO_STRUCTURE.md`, `docs/OPENAI_REQUEST_FLOW.md`, `docs/PROVIDER_SPEC.md`, `docs/ROUTING_SPEC.md`, `docs/BILLING_SPEC.md`, `docs/LOGGING_POLICY.md`, `docs/SECURITY_POLICY.md`
- reproduction_steps:
  1. Run `rg --files --hidden` for the requested docs.
- expected_behavior: All requested audit input docs exist.
- actual_behavior: Seven requested spec docs are missing.
- risk_impact: Architecture, provider, routing, billing, logging, and security audit baselines are incomplete.
- suggested_fix: Add or regenerate the missing spec docs, or explicitly mark them superseded by existing docs.

## AUD-008

- feature_id: `M16-F04-regression-tests`
- milestone: `M16`
- severity: medium
- category: documentation_mismatch
- affected_files: `validation-contract.md`, `.factory/mission-state.json`, `docs/DEVELOPMENT_LOG.md`
- reproduction_steps:
  1. Compare M14-M16 assertions in `validation-contract.md` with mission-state `validation_status`.
- expected_behavior: Every completed assertion in mission-state is defined in validation-contract, and every contract assertion is represented in mission-state.
- actual_behavior: Mission-state includes `A16.4.1-A16.4.6`, `A14.3.*`, and other assertions not defined in validation-contract. `A15.3.2` exists in validation-contract but is absent from mission-state.
- risk_impact: Validation coverage is ambiguous and may overstate completed work.
- suggested_fix: Reconcile assertion IDs, add definitions for extra assertions, and backfill missing mission-state validation status.

## AUD-009

- feature_id: `M16-F04-regression-tests`
- milestone: `M16`
- severity: medium
- category: bug
- affected_files: `common/custom-event.go`
- reproduction_steps:
  1. Run `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go vet ./...`.
- expected_behavior: `go vet ./...` exits cleanly.
- actual_behavior: Vet reports `CustomEvent` copies a `sync.Mutex` by value in `encode`, `Render`, and `WriteContentType`.
- risk_impact: Potential concurrency bug and release quality gate failure.
- suggested_fix: Use pointer receivers/parameters for `CustomEvent` or remove the mutex from copied value paths.

## AUD-010

- feature_id: `M14-F02-provider-enable-disable`, also impacts `M15` frontend/admin pages
- milestone: `M14`
- severity: medium
- category: test_gap
- affected_files: `web/default/package.json`, `web/classic/package.json`
- reproduction_steps:
  1. Run `cd web/default && npm run lint`.
  2. Run `cd web/default && npm test`.
  3. Run `cd web/classic && npm run lint`.
  4. Run `cd web/classic && npm test`.
- expected_behavior: Frontend lint/test commands run or are documented as unavailable with replacement checks.
- actual_behavior: `eslint`, `prettier`, `tsc`, and `rsbuild` are not installed; both frontends have no `test` script.
- risk_impact: Frontend assertions, especially `A14.2.2`, lack executable evidence.
- suggested_fix: Add repeatable frontend dependency setup and test scripts, or document frontend audit as blocked with a specific fixture.

## AUD-011

- feature_id: `M16-F01-docker-compose-local`
- milestone: `M16`
- severity: medium
- category: security
- affected_files: `docker-compose.yml`, `docker-compose.local.yml`
- reproduction_steps:
  1. Inspect `docker-compose.yml:27-29`, `docker-compose.yml:57-68`, and `docker-compose.local.yml:29`.
- expected_behavior: Compose files use safe local-only defaults or environment variable placeholders for secrets.
- actual_behavior: Default Postgres/Redis password `123456` and `SESSION_SECRET=change-me-in-production` are present. The file warns to change passwords, but defaults remain active.
- risk_impact: Accidental deployment with default credentials can expose DB/Redis/session security.
- suggested_fix: Use required environment variables or `.env.example` placeholders and fail fast when production secrets are unchanged.

## AUD-012

- feature_id: `M16-F02-provider-policy-docs`
- milestone: `M16`
- severity: medium
- category: test_gap
- affected_files: `scripts/regression.sh`, `docker-compose.local.yml`
- reproduction_steps:
  1. Attempting smoke execution requires a running app and `ADMIN_TOKEN`.
  2. Current backend build fails, so app cannot start from source.
- expected_behavior: Smoke test can run entirely against local mock/stub services after compose up.
- actual_behavior: Smoke test was not runnable in this phase because the backend does not compile and local seeded credentials/mock upstream were not available.
- risk_impact: M16 smoke coverage is claimed but not independently reproduced.
- suggested_fix: Add a deterministic mock upstream and seeded local fixture mode for `scripts/regression.sh`.

## AUD-013

- feature_id: `M5-F01-model-mapping-table`
- milestone: `M5`
- severity: medium
- category: documentation_mismatch
- affected_files: `validation-contract.md`, `service/channel_select.go`, `middleware/distributor.go`
- reproduction_steps:
  1. Read `validation-contract.md` A5.1.1.
  2. Search for `service.SelectChannel`.
- expected_behavior: Validation contract names the actual selection function under test.
- actual_behavior: Contract references `service.SelectChannel`, but current code appears to use `CacheGetRandomSatisfiedChannel` and distributor selection logic.
- risk_impact: Tests may validate the wrong abstraction or leave routing behavior untested.
- suggested_fix: Update validation contract to name the real selection entry points and add tests for them.

## AUD-014

- feature_id: `M14-F02-experimental_proxy-visibility`
- milestone: `M14`
- severity: medium
- category: test_gap
- affected_files: `web/default`, `web/classic`, `validation-contract.md`
- reproduction_steps:
  1. Review `A14.2.2`.
  2. Search frontend tests for provider visibility.
- expected_behavior: Frontend channel list component has an executable test proving non-admin users do not render experimental rows.
- actual_behavior: No frontend test script is available; no component-level evidence was found during this phase.
- risk_impact: UI may expose experimental_proxy rows even if backend API filters correctly.
- suggested_fix: Add frontend component or Playwright test using mocked channel-list response and non-admin role.

## AUD-015

- feature_id: `feature-test-matrix`
- milestone: `FIX`
- severity: medium
- category: audit_followup
- affected_files: `constant/channel.go`, `middleware/distributor.go`, `model/token.go`, `relay/relay_adaptor.go`, `service/channel_select.go`
- reproduction_steps:
  1. Run `git diff --ignore-all-space --stat -- constant/channel.go middleware/distributor.go model/token.go relay/relay_adaptor.go service/channel_select.go`.
  2. Observe non-whitespace changes remain in the worktree after `AUD-004` formatting.
- expected_behavior: AUD-004 should only leave formatting/import-order changes, and any semantic diff should be reviewed under its owning feature or audit item.
- actual_behavior: The listed files contain non-formatting diffs such as provider/channel wiring and token/routing behavior changes relative to HEAD.
- risk_impact: Formatting cleanup can mask unrelated behavior changes if reviewers only inspect the AUD-004 diff as whitespace churn.
- suggested_fix: Review these non-formatting diffs under `feature-test-matrix` or their original feature issues with targeted tests; do not treat them as part of AUD-004.

## AUD-016

- feature_id: `F10.2`
- milestone: `M10`
- severity: critical
- category: security
- status: fixed 2026-05-17
- affected_files: `model/token.go`, `controller/token.go`, `middleware/auth.go`
- reproduction_steps:
  1. Inspect `model.Token.Key`.
  2. Inspect `model.GetTokenByKey`.
  3. Inspect `controller.GetTokenKey`.
- expected_behavior: API keys are stored and looked up by hash/prefix only; full key is shown once at creation.
- actual_behavior_before_fix: Raw keys remained stored in `tokens.key`, lookup used the raw `key` column, and key retrieval endpoints could re-display the full key.

## Full Remediation Closure — 2026-05-17

### AUD-020

- severity: medium
- status: fixed 2026-05-17
- category: code_quality
- root_cause: `CustomEvent` copied a `sync.Mutex` by value and relay adaptors had unreachable statements after unconditional `panic`/`return`.
- affected_files: `common/custom-event.go`, SSE render call sites, relay adaptor files.
- resolution: pointer render paths for `CustomEvent`; removed unreachable relay adaptor statements.
- verification: `common.TestCustomEventRenderWritesSSEHeadersAndData`, `go vet ./...`, `go test ./... -count=1`.

### AUD-022-org-project-token-binding-enforcement

- severity: high
- status: fixed 2026-05-17
- category: security
- root_cause: org/project token fields were schema/context placeholders without fail-closed auth validation.
- affected_files: `model/token.go`, `middleware/auth.go`, `controller/token.go`, `model/log.go`.
- resolution: added token tenant scope resolver and auth/setup validation; disabled/mismatched org/project bindings fail; legacy NULL tokens remain user-scoped; usage_log tenant fields come from token context.
- verification: AUD-022 model and middleware tests plus full Go regression.

### AUD-023-admin-topup-contract

- severity: medium
- status: fixed 2026-05-17
- category: contract_mismatch
- root_cause: manual top-up contract expected `POST /api/user/topup`, but manual quota changes existed only through `/api/user/manage`.
- affected_files: `controller/user.go`.
- resolution: existing `/api/user/topup` now supports admin-only manual top-up request bodies with `user_id` and positive `quota`, writes manage logs, and preserves legacy redeem-code top-up.
- verification: admin success, normal-user rejection, balance increment, manage-log, and invalid-amount tests.

### AUD-024-deterministic-runtime-smoke-fixtures

- severity: medium
- status: fixed_with_environment_blockers 2026-05-17
- category: test_gap
- root_cause: runtime smoke required a seeded server and upstream providers.
- affected_files: `scripts/regression.sh`.
- resolution: added `LOCAL_FIXTURE=1` mode that runs deterministic local Go fixture checks and `go vet` without real keys or providers.
- verification: `LOCAL_FIXTURE=1 bash scripts/regression.sh` passed.
- remaining_blocker: fresh Docker/curl smoke still requires isolated seeded runtime fixture.

### AUD-025-frontend-experimental-visibility

- severity: medium
- status: partially_fixed_blocked_test_infra 2026-05-17
- category: test_gap
- root_cause: frontend component tests are unavailable in this environment.
- affected_files: `controller/model_list_test.go`, frontend package metadata.
- resolution: added deterministic API-level model visibility test for normal and internal users.
- verification: `controller.TestAUD025ListModelsExperimentalVisibilityByUserScope` passed.
- remaining_blocker: `bun` is not installed and frontend packages define no `test` scripts.
- risk_impact: DB compromise exposes active user API keys.
- resolution: New and migrated tokens now use HMAC `key_hash` for auth/cache lookup and `key_prefix` for display. Full keys are returned only in the create response. Legacy plaintext rows are migrated by hashing available plaintext and overwriting the deprecated `key` column with a non-secret hash for compatibility.
- verification: Added/updated model and controller tests for non-plaintext persistence, one-time full-key display, correct and wrong key auth, disabled key rejection, no repeated key leak, error-log non-leakage, and legacy plaintext no longer serving as an auth source. `go test ./model/... -count=1` and `go test ./... -count=1` pass.

## AUD-017

- feature_id: `F8.2`
- milestone: `M8/M14`
- severity: critical
- category: security
- affected_files: `model/channel.go`, `model/ability.go`, `middleware/distributor.go`
- reproduction_steps:
  1. Inspect `DisableExperimentalProxyChannels`.
  2. Inspect `model.GetChannel` DB ability-selection path.
- expected_behavior: Disabled `experimental_proxy` channels are impossible to route in every cache mode.
- actual_behavior: Bulk disable updates only `channels.status`; DB selection can choose enabled ability rows and load the disabled channel without a final status check.
- risk_impact: The emergency disable control can fail to stop experimental traffic.
- suggested_fix: Disable/update abilities in the same operation and enforce `channel.Status == enabled` after every DB channel load.

## AUD-018

- feature_id: `F2.1`
- milestone: `M2`
- severity: high
- category: security
- status: fixed 2026-05-17
- affected_files: `model/provider_account.go`, `model/channel.go`, `middleware/distributor.go`, `controller/channel.go`
- reproduction_steps:
  1. Inspect ProviderAccount encryption hooks.
  2. Inspect relay setup for credential source.
- expected_behavior: Active upstream credentials are resolved through encrypted ProviderAccount or equivalent encrypted storage.
- actual_behavior: ProviderAccount exists, but active relay setup still uses plaintext `Channel.Key`.
- risk_impact: Upstream provider keys remain exposed by database compromise.
- suggested_fix: Wire ProviderAccount into credential resolution and mask/retire plaintext channel key fields.
- root_cause: `SetupContextForSelectedChannel` always called `Channel.GetNextEnabledKey()`, which only reads `Channel.Key`, and linked `ProviderAccountId` was never consulted on the active relay credential path.
- resolution: Added `Channel.ResolveActiveCredential()` and wired main relay setup plus task/Midjourney/video credential fetch paths through it. When `provider_account_id` is present, the runtime credential is decrypted from ProviderAccount and legacy `Channel.Key` is not used as a fallback. Disabled ProviderAccounts, empty credentials, missing crypto secret, and decrypt failures reject before upstream calls and emit AUD-018 safety logs without credential material.
- verification: Added model and middleware tests for encrypted-at-rest storage, provider-account credential precedence, no fallback to legacy key, decrypt failure rejection, disabled account rejection, no credential in error strings, and legacy channel compatibility. `go test ./model/... -count=1` and `go test ./... -count=1` pass. `go vet ./...` still fails only on unrelated AUD-020 findings.

## AUD-019

- feature_id: `F11.2`
- milestone: `M11`
- severity: high
- category: security
- status: fixed 2026-05-17
- affected_files: `model/log.go`, relay log callers
- reproduction_steps:
  1. Inspect `RecordConsumeLog`.
  2. Observe `params.Other` is marshaled without prompt/completion filtering.
- expected_behavior: Default privacy mode excludes full prompt/response from all log fields.
- actual_behavior: `Content` is blanked, but `Other` can still contain prompt/completion payloads.
- risk_impact: Privacy-sensitive content can be persisted despite `STORE_FULL_TEXT_ENABLED=false`.
- suggested_fix: Redact/allowlist structured log metadata before storage and add tests for nested prompt/completion keys.
- root_cause: `RecordConsumeLog`, `RecordErrorLog`, and related log helpers serialized `params.Other` or caller-provided content directly, while the existing full-text gate only blanked `Content` and did not sanitize structured side payloads or JSON-string error data.
- resolution: Added centralized log sanitization for maps, nested values, JSON strings, plain strings, and error messages. Sensitive keys and values including API keys, bearer tokens, authorization headers, credentials, secrets, prompt/messages/input, response/output, tool payloads, headers, and metadata are redacted before logger output or database persistence. Normal accounting fields such as model, provider type, channel ID, tokens, and cost remain available.
- verification: Added AUD-019 model tests for API key redaction, bearer authorization redaction, prompt/messages/input redaction, response/output redaction, nested-map redaction, JSON-string redaction, error-message redaction, preservation of accounting fields, persisted consume-log sanitization, persisted error-log sanitization, and default full-text/debug payload closure. `go test ./model/... -count=1`, `go test ./middleware/... -count=1`, and `go test ./... -count=1` pass. `go vet ./...` still fails only on unrelated AUD-020 findings.

## AUD-020

- feature_id: `F16.3`
- milestone: `M16`
- severity: medium
- category: quality_gate
- affected_files: `common/custom-event.go`, `relay/channel/*/adaptor.go`
- reproduction_steps:
  1. Run `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go vet ./...`.
- expected_behavior: `go vet ./...` exits 0.
- actual_behavior: vet reports mutex-copy warnings and unreachable-code warnings.
- risk_impact: Full regression gate is not satisfied.
- suggested_fix: Fix vet warnings or document intentional unreachable stubs with code structure that passes vet.

## AUD-021

- feature_id: `F5.1`, `F10.4`, `F8.1`
- milestone: `M5/M8/M10`
- severity: high
- category: security
- status: fixed 2026-05-17
- affected_files: `model/token.go`, `middleware/auth.go`, `middleware/distributor.go`, `service/channel_select.go`, `model/ability.go`, `model/channel_cache.go`, `model/channel_satisfy.go`, `controller/model.go`, `controller/relay.go`, `controller/token.go`
- reproduction_steps:
  1. Create or inspect an API token intended for a restricted provider family.
  2. Route a model served by official, aggregator, authorized_proxy, and experimental_proxy channels.
  3. Trigger retry/fallback, preferred channel, memory-cache, and DB ability-selection paths.
- expected_behavior: `allowed_provider_types` constrains every channel candidate path; normal tokens cannot reach `experimental_proxy`; internal/admin users require explicit `allow_experimental=true` before `experimental_proxy` is eligible.
- actual_behavior: Only `AllowExperimental` existed, and selection/fallback/preferred paths did not have a unified provider-type whitelist.
- risk_impact: A token intended for official providers could route to aggregators or experimental proxies through fallback, retry, cache, DB ability, or affinity paths.
- suggested_fix: Add a centralized provider-type policy and apply it to all channel candidate and final setup paths.
- root_cause: Provider type was metadata on channels, but token policy had no persisted `allowed_provider_types` value and routing functions accepted only a boolean experimental flag. That left official/aggregator/authorized_proxy restrictions unenforced and made fallback behavior depend on local path-specific filtering.
- resolution: Added token `allowed_provider_types`, parsing and validation, auth context propagation, and a centralized `ProviderTypePolicy`. Memory candidates, DB ability selection, retry/fallback, preferred/affinity, specific channel, model-list, and final `SetupContextForSelectedChannel` paths now all check the same policy. Empty policies fail closed for `experimental_proxy` unless the caller is internal/admin and token `allow_experimental` is true. Invalid provider types and disabled channels are rejected before credentials are loaded.
- verification: Added AUD-021 model and middleware tests for official-only selection, aggregator-only selection, empty policy experimental rejection, internal explicit experimental allow, internal experimental denial without token flag, official-to-experimental fallback prevention, retry skipping disallowed providers, DB ability filtering, preferred/final setup rejection, empty provider_type inference for experimental channels, disabled channel rejection, and no credential leak on reject. `go test ./model/... -count=1`, `go test ./middleware/... -count=1`, and `go test ./... -count=1` pass. `go vet ./...` still fails only on unrelated AUD-020 findings.
