# Codex Full Remediation Plan

Generated: 2026-05-17

Scope: remaining audit/security/code-risk/test-gap closure after AUD-016, AUD-017, AUD-018, AUD-019, and AUD-021. No real upstream provider keys or paid providers may be used. Runtime/provider tests must use local mocks, stubs, or fixtures.

## Open Issues

### AUD-022-org-project-token-binding-enforcement

- severity: high
- root cause: token organization/project fields were added as schema/context placeholders, but request authorization, billing, provider policy, and log attribution did not consistently enforce them as the trusted scope.
- affected features: F9.1, F9.2, F10.1, F10.4, F11.1, F12.1, F12.2, F13.1, F14.2
- affected validation assertions: A9.1.1, A9.1.2, A9.2.1, A9.2.2, A10.1.1, A10.1.2, A11.1.2, A12.1.1, A12.1.2, A12.2.1, A13.1.2, A14.2.1
- planned fix: resolve token-bound organization_id/project_id during token authentication; ignore client-supplied tenant scope; reject disabled bound organization/project; preserve legacy NULL bindings as legacy user scope only; ensure logs, billing context, allowed_models, and allowed_provider_types use the authenticated token scope.
- tests to add: project A cannot access project B; organization A cannot access organization B; usage_log tenant fields come from token context; spoofed request tenant fields are ignored; disabled project/organization bound tokens fail; legacy token remains user-scoped; provider/model/billing checks consume authenticated token scope.
- commands to run: `go test ./model/... -count=1`, `go test ./middleware/... -count=1`, `go test ./... -count=1`, `go vet ./...`
- expected completion criteria: token-bound tenant scope is authoritative and tested; no legacy token gains cross-tenant rights.
- blocked conditions: if existing organization/project models cannot represent disabled state, record the missing state as `blocked_product_decision` and fail closed for non-active statuses.

### AUD-023-admin-topup-contract

- severity: medium
- root cause: validation/docs expect `POST /api/user/topup`, while implementation exposes top-up through generic `POST /api/user/manage` with `action=add_quota`.
- affected features: F13.2, F15.3
- affected validation assertions: A13.2.1, A13.2.2, A15.3.2
- planned fix: add explicit admin-only `POST /api/user/topup` route that validates positive quota and reuses existing balance adjustment/logging logic.
- tests to add: admin success; normal user rejected; balance increases; manage log written; invalid/negative quota rejected; route/documentation alignment.
- commands to run: `go test ./controller/... -count=1`, `go test ./... -count=1`
- expected completion criteria: `/api/user/topup` is stable, admin-only, positive-only, and documented in audit outputs.
- blocked conditions: none expected.

### AUD-020-go-vet-failures

- severity: medium
- root cause: `common/custom-event.go` copies a value containing `sync.Mutex`; several relay adaptors contain code after unconditional returns.
- affected features: F16.3 and full regression quality gate
- affected validation assertions: A16.3.2
- planned fix: use pointer receivers or pointer parameters for `CustomEvent` render/write paths and remove unreachable statements without changing relay behavior.
- tests to add: custom event render/content-type behavior; relay adaptor package tests remain passing; `go vet ./...` as verification.
- commands to run: `go test ./common/... -count=1`, `go test ./relay/... -count=1`, `go vet ./...`
- expected completion criteria: `go vet ./...` exits 0 or only documented generated/third-party warnings remain.
- blocked conditions: none expected.

### AUD-024-deterministic-runtime-smoke-fixtures

- severity: medium
- root cause: runtime curl, Docker, and smoke assertions require a running gateway plus seeded users/channels and upstream responders, but no deterministic local fixture exists.
- affected features: F3.1, F3.2, F8.1, F8.2, F11.2, F12.1, F12.2, F13.1, F16.1, F16.2
- affected validation assertions: A3.1.1, A3.1.2, A3.2.1, A8.1.2, A8.2.1, A11.2.1, A12.1.1, A12.1.2, A12.2.1, A13.1.1, A16.1.1, A16.2.1
- planned fix: add local fake upstream/smoke fixture coverage where feasible without real keys; update smoke documentation to run deterministically or mark environment-only Docker execution as skipped when Docker/runtime is unavailable.
- tests to add: fake official provider success with usage; fake experimental provider requires internal + allow_experimental; normal user and disabled experimental rejection; zero balance prevents upstream call; usage_log excludes prompt/response; smoke script local mode.
- commands to run: `go test ./... -count=1`, `docker compose config`, local smoke script if environment permits.
- expected completion criteria: deterministic unit/integration evidence exists; any Docker/runtime non-execution is explicitly recorded as `skipped_environment`.
- blocked conditions: real Docker daemon or long-running server unavailable in the current environment.

### AUD-025-frontend-experimental-visibility

- severity: medium
- root cause: frontend packages lack executable test scripts and current environment lacks Bun; non-admin experimental visibility has only backend/API evidence.
- affected features: F14.2, F15 admin UI assertions
- affected validation assertions: A14.2.1, A14.2.2, A15.3.2
- planned fix: add deterministic API-level visibility tests for normal/admin/internal scopes; run available frontend scripts if dependencies exist; otherwise document `blocked_test_infra` with minimum test-infra path.
- tests to add: normal model/provider list excludes experimental_proxy; admin view can include experimental_proxy; disabled experimental visible as disabled to admin but unroutable; internal explicit allow is usable.
- commands to run: inspect `web/default/package.json`; run available `bun run lint`, `bun run build`, or document blocked if Bun/dependencies are absent.
- expected completion criteria: backend visibility tests pass and frontend execution status is accurately documented.
- blocked conditions: Bun/dependencies/test scripts absent.

### AUD-006/AUD-008/AUD-015/AUD-026-doc-state-reconciliation

- severity: medium/low
- root cause: feature source files, validation contract, development log, and mission-state contain stale statuses and mismatched assertion/API names after remediation.
- affected features: global M0-M16 feature-test-matrix
- affected validation assertions: all feature matrix assertions with changed status, especially A6.3.1, A13.2.1, A15.3.2, A16.3.2
- planned fix: update audit report, bug list, fix recommendations, feature results, security findings, development log, and mission-state to reflect the final verified state.
- tests to add: JSON parse check for mission-state; `git diff --check`.
- commands to run: `node -e "JSON.parse(require('fs').readFileSync('.factory/mission-state.json','utf8'))"`, `git diff --check`
- expected completion criteria: docs reflect real code/test state, including any blockers.
- blocked conditions: product decisions for missing spec docs or broad validation contract rewrites can be recorded as `blocked_product_decision`.

## Final Regression Commands

- `/tmp/go2510/bin/gofmt` on modified Go files
- `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./model/... -count=1`
- `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./middleware/... -count=1`
- `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./... -count=1`
- `env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go vet ./...`
- `git diff --check`
- `node -e "JSON.parse(require('fs').readFileSync('.factory/mission-state.json','utf8'))"`
- `docker compose config`
- available frontend lint/test/build scripts, or explicit `blocked_test_infra`

## Completion Criteria

- Critical/high findings are zero or explicitly blocked with minimum fix path.
- Every fixed issue has a matching test or deterministic verification command.
- Feature-test-matrix is rerun and audit documents are updated.
- `.factory/mission-state.json` records `all_remediation_status`, `fixed_issues`, `blocked_issues`, `remaining_issues`, `final_test_status`, `final_security_status`, `feature_matrix_status`, `updated_at`, and `next_recommended_action`.
- Final handoff includes `development_log_updated: true` and `mission_state_updated: true`.

## Execution Result — 2026-05-17 22:05 +08:00

Completed:

- `AUD-022-org-project-token-binding-enforcement`
- `AUD-023-admin-topup-contract`
- `AUD-020-go-vet-failures`
- `AUD-024-deterministic-runtime-smoke-fixtures` through local fixture mode
- `AUD-025-frontend-experimental-visibility` at deterministic API level

Blocked:

- `blocked_test_infra_frontend`: Bun is unavailable and frontend package test scripts are absent.
- `blocked_external_dependency_cross_db_runtime`: live MySQL/PostgreSQL migration proof requires external services.
- `skipped_environment_docker_runtime`: a fresh seeded Docker/curl runtime fixture was not executed against an isolated environment.

Final backend quality gates passed: `go test ./model/...`, `go test ./middleware/...`, `go test ./...`, `go vet ./...`, `git diff --check`, mission-state JSON parse, Docker compose config, and `LOCAL_FIXTURE=1 bash scripts/regression.sh`.
