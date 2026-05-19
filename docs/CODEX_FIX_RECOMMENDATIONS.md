# Codex Fix Recommendations

## Priority 0: Restore Buildability

1. Fixed 2026-05-17: `model/log.go` malformed `GetAllLogs` block (`AUD-001`).
   - Result: `go build ./model` passes after the brace/control-flow repair.
   - Added test: `TestGetAllLogsFiltersProviderTypeAndExperimentalProxy`.
   - Remaining verification blocker: `go test ./model/...` fails in `model/regression_m16_test.go`, which must be handled separately.

2. Fixed 2026-05-17: `controller/token.go` orphan statements after `GetAdminAllTokens` (`AUD-002`).
   - Result: restored the missing `SearchTokens` handler declaration around the existing user token search body.
   - Verify: `gofmt -w controller/token.go`; file-scoped token controller tests pass. Full controller package tests remain blocked by `AUD-003`.

3. Fixed 2026-05-17: `controller/channel.go` orphan struct fields around line 738 (`AUD-003`).
   - Result: restored the missing `ChannelTag` struct declaration before the existing tag-management fields.
   - Verify: `gofmt -w controller/channel.go`, `go test ./controller -count=1`, and `go build ./...` pass.

4. Fixed 2026-05-17: remaining M16 regression test compile errors and `AUD-004` formatting drift.
   - Result: `/tmp/go2510/bin/gofmt -l $(rg --files -g '*.go')`, `go test ./... -count=1`, and `git diff --check` pass.
   - Frontend note: `bun run lint` could not run because `bun` is unavailable in the current environment; no frontend `test` script exists in `web/default` or `web/classic`.
   - Follow-up: `AUD-015` records non-formatting worktree diffs found during review and keeps them out of the AUD-004 scope.

## Priority 1: Security Correctness

5. Fixed 2026-05-17: rework token key storage so new tokens do not persist plaintext `sk-` values.
   - Result: token auth/cache lookup uses HMAC `key_hash`; display uses `key_prefix`; full keys are shown only once; legacy plaintext rows are migrated to non-secret hash storage.
   - Tests: DB row does not contain full raw key, auth lookup uses hash, disabled tokens fail auth, legacy plaintext is not an auth source, error logs do not include full keys, and user responses do not re-display full keys.

6. Fixed 2026-05-17: make disabled `experimental_proxy` channels unroutable in every selection mode.
   - Result: final routeability checks cover DB ability, memory candidate, fallback/retry, preferred channel, and setup paths.

7. Fixed 2026-05-17: wire encrypted ProviderAccount credentials into the active relay path.
   - Result: `provider_account_id` channels resolve runtime credentials from decrypted ProviderAccount data and do not silently fallback to plaintext `Channel.Key`.
   - Tests: encrypted-at-rest storage, provider credential precedence, no legacy fallback, decrypt failure rejection, disabled account rejection, credential non-leakage in errors, and legacy channel compatibility.

8. Fixed 2026-05-17: enforce `allowed_provider_types` and provider-type fallback restrictions (`AUD-021`).
   - Result: API token provider-type policy is applied to DB ability selection, memory candidates, retry/fallback, preferred/affinity, specific channel, model list, and final setup paths.
   - Tests: official-only, aggregator-only, normal empty policy, internal allow/deny experimental, fallback/retry, DB ability, preferred/final setup, empty provider_type, disabled channel, and credential non-leakage.

9. Fixed 2026-05-17: redact all log metadata under default privacy mode (`AUD-019`).
   - Result: `params.Other`, JSON strings, plain strings, and error messages are centrally sanitized.

10. Harden compose secrets.
   - Replace production defaults with required env vars.
   - Keep `docker-compose.local.yml` explicitly local-only.

## Priority 2: Validation And Test Coverage

11. Reconcile mission documents.
   - Update `features.json` statuses or mark it superseded.
   - Align assertion IDs across `validation-contract.md`, mission-state, and development log.
   - Add missing requested spec docs or document replacements.

12. Add deterministic local mock upstream and seeded test fixture mode.
   - Needed for `/v1/models`, `/v1/chat/completions`, zero quota, top-up, privacy log, and experimental_proxy smoke tests.
   - Ensure no real provider keys are needed.

13. Add frontend test/lint reproducibility.
   - Provide local dependency setup instructions.
   - Add `test` scripts or document why not applicable.
   - Add frontend tests for non-admin experimental_proxy visibility.

## Priority 3: Quality Gates

14. Fix `go vet` warnings.
   - Fix mutex-copy warnings in `common/custom-event.go`.
   - Fix or restructure unreachable code warnings in relay channel adaptors.

15. Add cross-DB migration tests.
   - SQLite should run in unit tests.
   - MySQL/PostgreSQL can run via disposable local containers in CI.
   - Cover Channel, ProviderAccount, Organization, Project, Token, Log.

16. Add admin route negative tests.
   - Normal user must receive 401/403 for admin token/log/channel/provider routes.
   - Admin route responses must not leak secrets.

17. Fixed 2026-05-17: add log privacy deep checks (`AUD-019`).
   - Nested `Other`, JSON-string `Other`, persisted consume logs, persisted error logs, and sanitized error messages are covered.
   - Full-text/debug payload storage remains default off through `StoreFullTextEnabled=false`.

18. Review non-formatting worktree diffs found during AUD-004 (`AUD-015`).
   - Scope: `constant/channel.go`, `middleware/distributor.go`, `model/token.go`, `relay/relay_adaptor.go`, `service/channel_select.go`.
   - Treat as a feature/test-matrix review, not as formatting cleanup.

## Suggested Fix Sequence

1. Syntax fixes and global gofmt drift are fixed through `AUD-004`.
2. Fixed 2026-05-17: `AUD-016` token plaintext storage.
3. Fixed 2026-05-17: `AUD-017` disabled experimental routing.
4. Fixed 2026-05-17: `AUD-018` ProviderAccount credential wiring.
5. Fixed 2026-05-17: `AUD-019` log privacy gap.
6. Fixed 2026-05-17: `AUD-021-allowed-provider-types-and-fallback-restriction`.
7. Fix `AUD-022-org-project-token-binding-enforcement`.
8. Add deterministic runtime fixtures for blocked curl/Docker/smoke assertions.
9. Reconcile mission docs and validation contract.
10. Fix `go vet` (`AUD-020`), then run optional `govulncheck`, `staticcheck`/`golangci-lint` if installed.

## Do Not Do During Fixes

- Do not call real paid providers.
- Do not use real upstream API keys.
- Do not silently change DB schema without migration tests.
- Do not mark any feature complete until code, test, and runtime evidence are all updated.

## Full Remediation Result — 2026-05-17

Completed:

- `AUD-022-org-project-token-binding-enforcement`: enforce token tenant scope, disabled org/project rejection, legacy user-scope compatibility, usage_log attribution from token context.
- `AUD-023-admin-topup-contract`: admin-only positive manual top-up path through `POST /api/user/topup` request bodies with operation logging.
- `AUD-020-go-vet-failures`: `go vet ./...` now passes.
- `AUD-024-deterministic-runtime-smoke-fixtures`: `LOCAL_FIXTURE=1 bash scripts/regression.sh` added and passed without real providers.
- `AUD-025-frontend-experimental-visibility`: deterministic API-level normal/internal experimental visibility test added and passed.

Still recommended before deployment:

- Review the accepted frontend waiver and require CI `bun run lint`, `bun run test`, and `bun run build` before final release.
- Review the accepted cross-DB waiver and require CI SQLite/MySQL/PostgreSQL migration jobs before final release.
- Review the accepted Docker runtime waiver and require seeded `docker-compose.fixture.yml` smoke in CI or staging before final release.
- Add Go module and runtime dependency caches for fixture builds.
- Set explicit timeout and retry behavior around fixture image build/dependency download.
- Prepare an independent staging environment for the complete seeded curl smoke before production.

## Pre-Release Hardening Recommendations — 2026-05-17

Added implementation/CI paths:

- `web/default` has `bun test` plus minimal `experimental_proxy` visibility tests.
- `scripts/check-migrations.sh` and `tests/migration` provide schema migration smoke for SQLite locally and MySQL/PostgreSQL via CI services.
- `docker-compose.fixture.yml`, `scripts/fake-openai-provider.mjs`, and `scripts/seed-local-fixture.sh` provide an isolated fake-provider runtime.
- `.github/workflows/pre-release-hardening.yml` wires the above into CI.

Remaining action: run and archive CI artifacts, then manually accept or reject the three waiver documents listed in `docs/PRE_DEPLOYMENT_REVIEW_CHECKLIST.md`.

## Final Verification Refresh Recommendations — 2026-05-18

- Restore a Go 1.22+ toolchain in the local shell or run final verification in CI/staging.
- Rerun `go test ./model/... -count=1`, `go test ./middleware/... -count=1`, `go test ./... -count=1`, `go vet ./...`, and `LOCAL_FIXTURE=1 bash scripts/regression.sh`.
- Verify fixture cleanup with `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes` when Docker CLI recovers.
- Keep frontend, cross-DB runtime, Docker runtime smoke, and final Go verification blockers as accepted environment blockers until CI/staging produces passing evidence.

## CI / Staging Verification Setup — 2026-05-18

Use the new verification assets to close environment blockers in a stable runner:

- Run `.github/workflows/pre-release-verification.yml`.
- Run `bash scripts/ci-verify.sh` in staging.
- Run `bash scripts/ci-migration-check.sh` for SQLite, MySQL, and PostgreSQL with service DSNs.
- Run the Docker fixture smoke only against the fake provider and fixture credentials.
- Review `docs/STAGING_VERIFICATION_RUNBOOK.md` before changing deployment readiness.

Do not mark any blocker passed until CI/staging evidence is available.

## Post-CI Verification Closure — 2026-05-19

GitHub Actions `Pre-release verification` run #13 passed on branch `main` at commit `aeb43e5` in approximately 2m37s. The passing jobs were `go-test-vet`, `local-fixture-regression`, `cross-db-migration`, `docker-fixture-smoke`, and `frontend-check`.

The prior local-environment blockers are now `closed_by_ci`: `blocked_test_infra_frontend`, `blocked_external_dependency_cross_db_runtime`, `skipped_environment_docker_runtime`, and `final_go_verification_blocked`.

Current recommendation state:

- `critical_findings_remaining = 0`
- `high_findings_remaining = 0`
- `features_failed = 0`
- CI verification passed.
- `deployment_readiness = staging_ready`
- `recommended_next_action = run staging manual verification using docs/STAGING_VERIFICATION_RUNBOOK.md`

Do not mark the release `production_ready` from CI alone. Production still requires staging verification, environment-variable review, real deployment topology review, and manual security sign-off.
