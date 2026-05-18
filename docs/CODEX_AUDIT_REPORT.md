# Codex Audit Report

## Scope

This audit executed the approved audit plan against M0-M16. It used local static inspection and local commands only. No real upstream provider was called. No paid provider key was used.

Fix phase note: `AUD-001`, `AUD-002`, `AUD-003`, `AUD-004`, critical `AUD-016`, critical `AUD-017`, high `AUD-018`, high `AUD-019`, and high `AUD-021` were fixed on 2026-05-17. The M16 regression test compile errors were also fixed on 2026-05-17. The original audit findings remain the baseline for the remaining confirmed issues.

## Overall Result

The full backend build and Go test suite now pass after the controller syntax blockers, M16 regression test compile errors, global gofmt drift, `AUD-016` token key storage remediation, `AUD-017` disabled experimental routing remediation, `AUD-018` ProviderAccount credential wiring, `AUD-019` log privacy remediation, and `AUD-021` provider-type routing policy enforcement were fixed. The feature-test-matrix retest is now executable, and it shows that several feature claims remain failed or blocked by missing runtime/cross-DB/frontend evidence. The highest-risk open failures are org/project token binding enforcement, missing deterministic runtime fixtures, and failing `go vet`.

`AUD-001` fixed the malformed `GetAllLogs` block in `model/log.go`; `go build ./model` now passes. `AUD-002` fixed the orphan token search handler body in `controller/token.go`; isolated token controller tests now pass. `AUD-003` restored the missing `ChannelTag` type declaration in `controller/channel.go`; `go build ./...` and `go test ./controller` now pass. The M16 regression test sync updated `model/regression_m16_test.go` to current interfaces, and `go test ./model/...` plus `go test ./...` now pass. `AUD-004` ran gofmt across the remaining listed files and `gofmt -l $(rg --files -g '*.go')` now returns no files. `AUD-016` moved user API key auth/cache lookup to HMAC `key_hash`, preserved only display `key_prefix`, removed repeat full-key retrieval, and migrated legacy plaintext rows to non-secret hash storage. `AUD-017` added routeability guards so disabled channels cannot be selected through DB ability, cache candidate, fallback, or final setup paths. `AUD-018` now resolves active upstream credentials from encrypted ProviderAccount records when `provider_account_id` is present, rejects disabled/decryption-failed provider accounts before upstream calls, and preserves legacy `Channel.Key` only when no ProviderAccount is linked. `AUD-019` added centralized log sanitization for `params.Other`, error messages, and log content so prompt/response bodies, bearer tokens, API keys, credentials, secrets, authorization headers, nested sensitive fields, and JSON-string payloads are redacted before logger or usage-log persistence. `AUD-021` added token `allowed_provider_types` parsing and unified provider policy enforcement across memory candidates, DB ability selection, retry/fallback, preferred/affinity channels, specific channel selection, model listing, and final selected-channel setup.

The regression test compile errors were caused by stale test contracts: `GetGroupEnabledModelsExcludingExperimental` now returns `[]string`, the insufficient-quota error code lives in `types`, and the StoreFullText default test name collided with `model/log_m11_test.go`.

## Feature Results

Detailed per-feature results are in `docs/CODEX_FEATURE_TEST_RESULTS.md`.

| Metric | Count |
|---|---:|
| Feature rows tested from `features.json` | 37 |
| Pass | 11 |
| Fail | 11 |
| Blocked | 15 |
| Validation assertions checked | 70 |

## Open Issue Counts After Feature Matrix Retest

| Severity | Count |
|---|---:|
| critical | 1 |
| high | 4 |
| medium | 10 |
| low | 1 |

Security-only findings are summarized in `docs/CODEX_SECURITY_FINDINGS.md`.

## Commands Executed

```bash
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./...
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go build ./...
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test -race ./...
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go vet ./...
/tmp/go2510/bin/gofmt -l $(rg --files -g '*.go')
docker compose -f docker-compose.local.yml config
docker compose -f docker-compose.dev.yml config
docker compose -f docker-compose.yml config
cd web/default && npm run lint
cd web/default && npm test
cd web/default && npm run typecheck
cd web/default && npm run build -- --help
cd web/classic && npm run lint
cd web/classic && npm test
command -v govulncheck
command -v staticcheck
command -v golangci-lint
```

## Command Outcomes

- `go test ./...`: PASS after the M16 regression test sync, `AUD-004` formatting pass, and `AUD-016` token storage remediation.
- `go build ./...`: PASS after `AUD-003`.
- `go test -race ./...`: Not rerun after `AUD-003`; original audit failed due compile blockers.
- `go vet ./...`: FAIL during feature-test-matrix. Remaining warnings include mutex-copy warnings in `common/custom-event.go` plus unreachable-code warnings in multiple relay channel adaptors.
- `gofmt -l`: PASS after `AUD-004`; no Go files are listed.
- Docker compose config checks: PASS for default, dev, and local compose files.
- Frontend lint: SKIPPED because `bun` is not installed in the current environment (`bun: command not found`).
- Frontend tests: SKIPPED because no `test` script exists in either frontend package.
- govulncheck/staticcheck/golangci-lint: SKIPPED, not installed.

## Top 10 Fix Priorities

1. Fixed 2026-05-17: `model/log.go` syntax errors in `GetAllLogs` (`AUD-001`).
2. Fixed 2026-05-17: `controller/token.go` orphan top-level statements (`AUD-002`).
3. Fixed 2026-05-17: `controller/channel.go` orphan struct fields (`AUD-003`).
4. Fixed 2026-05-17: `model/regression_m16_test.go` stale compile-time contracts; `go test ./...` now passes.
5. Fixed 2026-05-17: global gofmt drift under `AUD-004`; `gofmt -l $(rg --files -g '*.go')` now returns no files.
6. Fixed 2026-05-17: reworked token storage so new and migrated API keys are not persisted or looked up in plaintext (`AUD-016`).
7. Fixed 2026-05-17: disabled `experimental_proxy` channels are unroutable in DB/cache/fallback/final setup paths (`AUD-017`).
8. Fixed 2026-05-17: active upstream credentials are resolved through encrypted ProviderAccount when linked (`AUD-018`).
9. Fixed 2026-05-17: structured log side payloads and error messages are sanitized before logging (`AUD-019`).
10. Fixed 2026-05-17: `allowed_provider_types` and provider-type fallback restrictions are enforced in channel selection (`AUD-021`).

## Documentation And State Findings

- `features.json` says 36 of 37 features are still `pending`; mission-state says the expanded mission set is completed.
- Requested docs are missing: `docs/REPO_STRUCTURE.md`, `docs/OPENAI_REQUEST_FLOW.md`, `docs/PROVIDER_SPEC.md`, `docs/ROUTING_SPEC.md`, `docs/BILLING_SPEC.md`, `docs/LOGGING_POLICY.md`, `docs/SECURITY_POLICY.md`.
- Assertion IDs drift between validation contract, development log, and mission-state.

## Recommendation

Do not proceed with additional feature development. Continue focused remediation with `AUD-022-org-project-token-binding-enforcement`, then add deterministic runtime fixtures and rerun the feature-test-matrix.

## Full Remediation Update — 2026-05-17 22:05 +08:00

Status: `completed_with_blockers`.

Backend critical/high remediation is closed. `AUD-022` now enforces token-bound organization/project scope during authentication and trusted context setup; disabled organizations/projects and mismatched project/org bindings fail closed, legacy NULL tenant tokens remain user-scoped, usage_log tenant fields come from token context, and spoofed request context is cleared. `AUD-023` adds the admin manual top-up contract on `POST /api/user/topup` for request bodies containing `user_id` and positive `quota`, while preserving legacy redeem-code top-up behavior. `AUD-020` is fixed and `go vet ./...` passes. `AUD-024` adds deterministic local fixture mode to `scripts/regression.sh`. `AUD-025` adds API-level experimental visibility tests.

Final command outcomes:

- `go test ./model/... -count=1`: PASS
- `go test ./middleware/... -count=1`: PASS
- `go test ./... -count=1`: PASS
- `go vet ./...`: PASS
- `git diff --check`: PASS
- mission-state JSON parse: PASS
- `docker compose config`: PASS
- `LOCAL_FIXTURE=1 bash scripts/regression.sh`: PASS
- frontend: BLOCKED, `bun` is not installed and frontend packages have no `test` scripts.

Remaining blockers are medium severity and environment/test-infrastructure scoped: frontend component test execution, cross-DB live migration proof, and isolated seeded Docker/curl runtime smoke. No critical or high backend security findings remain.

## Pre-Release Hardening Update — 2026-05-17 23:30 +08:00

Status: `completed_with_accepted_blockers`.

No business feature was added. The hardening pass added executable fixture/CI paths and waiver documentation for the three remaining environment blockers:

- `blocked_test_infra_frontend`: accepted blocker. `web/default` now has a `test` script and minimal Bun tests for `experimental_proxy` visibility, but local execution is blocked because Bun and frontend dependencies are absent. Waiver: `docs/WAIVERS/FRONTEND_TEST_INFRA_WAIVER.md`.
- `blocked_external_dependency_cross_db_runtime`: accepted blocker. `scripts/check-migrations.sh` and `tests/migration` verify the required schema fields; SQLite ran locally. MySQL/PostgreSQL remain CI-service gated. Waiver: `docs/WAIVERS/CROSS_DB_MIGRATION_WAIVER.md`.
- `skipped_environment_docker_runtime`: accepted blocker. `docker-compose.fixture.yml` started the local build path with Docker daemon access and without real provider calls or real API keys, but the build stalled at `go mod download` / dependency download. This is not a business test failure. Waiver: `docs/WAIVERS/DOCKER_RUNTIME_SMOKE_WAIVER.md`.

CI path added: `.github/workflows/pre-release-hardening.yml`, covering backend gates, frontend Bun lint/test/build, SQLite/MySQL/PostgreSQL migration smoke, and Docker fixture smoke.

Deployment readiness remains `needs_manual_review`; do not mark ready until the waiver items are accepted with CI/manual evidence.

## Final Verification Refresh — 2026-05-18

Status: `completed_with_environment_blockers`.

No business logic was changed. The refresh verified that the current shell has no usable Go toolchain: `which go` returned no path, `go version` failed, `/tmp/go2510/bin/go` is absent, and no Go binary was found under `/tmp`. Final Go checks are therefore blocked locally as `final_go_verification_blocked`, with waiver `docs/WAIVERS/LOCAL_GO_TOOLCHAIN_WAIVER.md`.

Docker fixture cleanup was attempted with `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes`, but Docker failed before compose action with `Failed to initialize: protocol not available`. This remains an environment cleanup blocker, not a business test failure.

Deployment readiness remains `needs_manual_review`; final verification must run in CI/staging before production.

## CI / Staging Verification Setup — 2026-05-18

Status: `completed_with_blockers`.

No business logic or provider behavior was changed. Added the CI/staging verification path for the remaining environment blockers:

- `scripts/ci-verify.sh` runs Go tests, Go vet, local fixture regression, git diff check, and available frontend lint/test/build scripts. Missing tools or dependencies are reported as blocked, not passed.
- `scripts/ci-migration-check.sh` wraps the existing migration smoke for SQLite/MySQL/PostgreSQL service-backed CI runs.
- `.github/workflows/pre-release-verification.yml` adds jobs for Go test/vet, local fixture regression, cross-DB migration, Docker fixture smoke, and frontend checks.
- `docs/STAGING_VERIFICATION_RUNBOOK.md` documents staging prerequisites, fake-provider smoke checks, manual checks, pass criteria, and failure handling.

The remaining blockers are now `pending_ci_verification`: `blocked_test_infra_frontend`, `blocked_external_dependency_cross_db_runtime`, `skipped_environment_docker_runtime`, and `final_go_verification_blocked`.

Deployment readiness remains `needs_manual_review`. Next recommended action: Run pre-release verification workflow in CI/staging and close accepted blockers before production.
