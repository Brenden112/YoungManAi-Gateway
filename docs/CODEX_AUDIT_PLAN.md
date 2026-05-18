# Codex Audit Plan

## Audit Scope

This is a read-only independent audit plan for the B2B Multi-Provider AI API Gateway MVP after Droid Mission M0-M16 completion.

The audit will verify:

- M0-M16 feature implementation claims.
- Every validation assertion in `validation-contract.md`, plus extra assertions recorded in `.factory/mission-state.json`.
- Code evidence, test evidence, and runtime evidence for each feature.
- Security, privacy, billing, routing, migration, and admin-surface risks.
- Whether `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json` accurately reflect the current code.

Out of scope for this phase:

- Fixing bugs.
- Adding features.
- Editing business code.
- Deleting files.
- Changing DB schema.
- Calling real paid upstream providers.
- Using real upstream API keys.
- Committing changes.

## Inputs Read

Present and read:

- `AGENTS.md`
- `mission.md`
- `validation-contract.md`
- `features.json`
- `docs/DEVELOPMENT_LOG.md`
- `.factory/mission-state.json`
- `.factory/library/architecture.md`
- `.factory/library/environment.md`
- `.factory/library/user-testing.md`

Requested but missing. These are audit findings and evidence gaps:

- `docs/REPO_STRUCTURE.md`
- `docs/OPENAI_REQUEST_FLOW.md`
- `docs/PROVIDER_SPEC.md`
- `docs/ROUTING_SPEC.md`
- `docs/BILLING_SPEC.md`
- `docs/LOGGING_POLICY.md`
- `docs/SECURITY_POLICY.md`

Related replacement or partial docs found:

- `docs/provider-policy.md`
- `docs/openai-sdk-quickstart.md`
- `docs/B2B_PROVIDER_GATEWAY_MVP_PLAN.md`
- `docs/DROID_MISSION_EXECUTION_PLAN.md`

## Initial Consistency Findings To Verify

- `features.json` still marks most features after `F0.1` as `pending`, while `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json` mark M0-M16 complete.
- Feature IDs drift between files: `features.json` uses `F1.1`, while mission-state uses names like `M1-F01-provider-type-enum`.
- `validation-contract.md` ends at `A16.3.2`, but mission-state records extra `A16.4.1` through `A16.4.6`.
- `docs/DEVELOPMENT_LOG.md` records extra assertion groups not present in `validation-contract.md`, such as `A5.3.*`, `A14.3.*`, and `A16.4.*`.
- `validation-contract.md` includes `A15.3.2`, but mission-state currently records only `A15.3.1` for M15 balance management.
- `validation-contract.md` asks for `docker-compose.dev.yml`; M16 handoff emphasizes `docker-compose.local.yml`. Both must be audited.
- Requested spec docs are missing, so architecture/routing/billing/security assertions currently depend on mission docs and implementation evidence.

## Evidence Model

Each assertion will be graded with three evidence columns:

- Code evidence: implementation exists in expected file and follows project conventions.
- Test evidence: unit/integration/smoke test exists and targets the assertion.
- Runtime evidence: command or local mock scenario was executed successfully.

Final audit statuses:

- Pass: code + test + runtime evidence are all sufficient.
- Partial: code exists but tests/runtime are weak or indirect.
- Fail: implementation contradicts requirement.
- Blocked: verification requires missing fixture, missing doc, external provider, unavailable dependency, or explicit approval.
- Not applicable: requirement is obsolete or superseded, with documented reason.

## Milestone Test Goals

| Milestone | Test goal |
|---|---|
| M0 | Verify repository discovery artifacts, factory library, mission files, feature list, and architecture trace accuracy. |
| M1 | Verify Channel provider classification fields, defaults, security defaults, and SQLite/MySQL/PostgreSQL migration compatibility. |
| M2 | Verify ProviderAccount encrypted key storage, Channel nullable ProviderAccountId, and legacy relay path compatibility. |
| M3 | Verify OpenAI-compatible `/v1/models` and `/v1/chat/completions` behavior did not regress after schema additions. |
| M4 | Verify provider_type default mapping, backfill behavior, and Channel CRUD acceptance/return of provider_type. |
| M5 | Verify model mapping table, disabled-model rejection, pricing fields, and provider_type routing interaction. |
| M6 | Verify KiroGateway channel constants, adapter skeleton, factory registration, and no main-flow hardcoding. |
| M7 | Verify internal user detection and experimental_proxy visibility/call blocking for normal users. |
| M8 | Verify experimental_proxy routing isolation, retry/fallback safety, and disable/cache invalidation behavior. |
| M9 | Verify Organization/OrganizationMember/Project stub schemas and non-interference with existing flows. |
| M10 | Verify Token org/project binding, key hash/prefix, disable helper, and token-gated experimental access. |
| M11 | Verify Log metadata fields, privacy defaults, no full prompt/response storage, and sanitized error logging. |
| M12 | Verify token usage aggregation, cost attribution, success-only deduction, and failure refund behavior. |
| M13 | Verify insufficient balance rejection before upstream call and admin manual recharge audit logging. |
| M14 | Verify provider/channel admin APIs, provider_type filtering, experimental visibility rules, and disable endpoint. |
| M15 | Verify admin token/log/balance endpoints, masking, filters, and top-up visibility. |
| M16 | Verify local compose, docs, regression script, full test/vet, and smoke coverage using local/mocked services. |

## Commands To Run

All commands must use local fixtures or mocks only.

Baseline static/runtime checks:

```bash
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go build ./...
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./... -count=1
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go vet ./...
```

Focused package tests:

```bash
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./constant/... ./model/... -count=1
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./service/... -count=1
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./controller/... ./middleware/... -count=1
env GOCACHE=/tmp/go-build-cache GOPATH=/tmp/go-path /tmp/go2510/bin/go test ./relay/... -count=1
```

Static audit scans:

```bash
rg -n 'encoding/json|json\\.Marshal|json\\.Unmarshal|json\\.NewDecoder|json\\.NewEncoder' --glob '*.go'
rg -n 'ProviderType|experimental_proxy|AllowExperimental|ContextKeyTokenAllowExperimental|ContextKeyIsExperimentalProxy' constant model service middleware controller relay
rg -n 'KeyHash|KeyPrefix|EncryptedKey|GetMaskedKey|MaskTokenKey|sk-' model controller middleware service
rg -n 'StoreFullTextEnabled|Content|prompt|completion|RecordConsumeLog|RecordErrorLog|MaskSensitive' common model controller service
rg -n 'AutoMigrate|ALTER TABLE|GROUP_CONCAT|STRING_AGG|JSONB|@>|\\?\\|' model
rg -n 'DisableExperimentalProxy|provider_type|is_experimental_proxy|GetAdminAllTokens|GetAllLogs' controller model router
```

Local smoke and compose checks. These are blocked until local test credentials/fixtures are confirmed:

```bash
docker compose -f docker-compose.local.yml up -d --build
curl -fsS http://localhost:3000/api/status
ADMIN_TOKEN=<local_admin_token> BASE_URL=http://localhost:3000 bash scripts/regression.sh
docker compose -f docker-compose.local.yml down
```

Do not run `scripts/regression.sh` against production or real paid providers.

Frontend checks. Run only if dependencies are already available locally; otherwise mark blocked rather than installing from network:

```bash
cd web/default && bun run build
cd web/default && bun run i18n:sync
```

## Supplemental Test Files Needed

Do not create these in this planning phase. Recommended additions for the execution phase:

- `model/migration_crossdb_test.go`: SQLite plus optional MySQL/PostgreSQL migration contract tests for Channel, Token, Log, Organization, Project.
- `middleware/experimental_proxy_integration_test.go`: normal/internal user experimental_proxy access matrix with mock channels.
- `controller/channel_provider_filter_integration_test.go`: provider_type filter, non-admin visibility, disable-experimental endpoint.
- `controller/admin_token_log_balance_integration_test.go`: admin token masking, log filters, balance view, top-up audit log.
- `service/billing_session_mvp_test.go`: pre-consume, settle, refund, zero-balance rejection, no upstream call on 402.
- `controller/openai_mock_relay_test.go`: local mock upstream for `/v1/models` and `/v1/chat/completions`.
- `web/default/src/**/__tests__/provider-channel-visibility.test.tsx`: frontend evidence for `A14.2.2`.
- `scripts/mock_upstream.py` or Go httptest fixture: deterministic non-paid upstream for smoke tests.

## Potential High-Risk Areas

- API key security: `Token.Key` may still retain plaintext for backward compatibility; audit whether plaintext is persisted or returned.
- Provider key encryption: AES key handling depends on `CRYPTO_SECRET`; audit fallback behavior and logs.
- experimental_proxy bypass: access requires user internal status and token `AllowExperimental`; audit all retry and direct channel paths.
- Fallback isolation: retry loops must never cross from official/aggregator/authorized_proxy into experimental_proxy.
- Billing correctness: pre-consume, post-settle, refund, and error paths must be idempotent and race-resistant.
- Usage log privacy: `Content`, `Other`, error strings, and admin APIs must not leak prompts/responses by default.
- Admin endpoint authorization: M14/M15 endpoints must be behind admin/root middleware and must filter tenant-scoped data.
- Cross-DB compatibility: raw SQL must quote `group`/`key` and avoid DB-specific syntax without fallback.
- Mission evidence drift: completed status may be overstated where runtime/curl evidence was not actually collected.
- Docker/smoke drift: validation contract references `docker-compose.dev.yml`; M16 handoff references `docker-compose.local.yml`.

## Blockers

- Seven requested specification files are missing.
- Full MySQL/PostgreSQL migration verification requires local DB services or approved containers.
- Smoke tests require local admin/user tokens and mock upstream channels; do not use real provider keys.
- Frontend audit requires local Bun dependencies; network install should be treated as blocked unless explicitly approved.
- `features.json` status values conflict with mission-state, so feature completion cannot be trusted from one source alone.
- Some mission-state assertion IDs are not present in `validation-contract.md`; these require manual reconciliation before final pass/fail grading.

## Recommended Test Execution Order

1. Documentation and state reconciliation: compare `features.json`, `validation-contract.md`, `docs/DEVELOPMENT_LOG.md`, and `.factory/mission-state.json`.
2. Static code evidence scan: models, migrations, context keys, route registration, middleware guards, billing flow.
3. Fast unit tests: `go test ./constant/... ./model/... ./service/...`.
4. Full backend regression: `go build ./...`, `go test ./... -count=1`, `go vet ./...`.
5. Security-focused tests: key storage, privacy logging, experimental_proxy access, admin auth.
6. Local integration tests with httptest/mock upstream.
7. Docker local compose and smoke script using local fixtures only.
8. Optional frontend build/tests if local dependencies exist.
9. Final audit report with per-feature pass/partial/fail/blocked status and remediation list.

