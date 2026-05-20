# Remote Codex Development Plan

Date: 2026-05-20

## Purpose

This document turns the current project review into an executable development plan for remote GitHub Codex workers.

Current baseline:

- Current branch: `main`
- Current HEAD reviewed by the user: `73ad2ff`
- Pre-release verification run: `#16`
- CI status for current HEAD: `Success`
- Production readiness: `not_ready`
- Required next release gate: isolated staging runtime verification and manual sign-off

Remote workers must work through pull requests. They must not push directly to `main`.

## Mandatory Rules For Remote Workers

Every remote Codex task must follow these rules:

- Create a dedicated branch and pull request.
- Do not use real upstream provider keys.
- Do not call paid providers.
- Do not commit `.env`, `.env.*`, key files, certificates, real prompts, real responses, or real customer data.
- Preserve all protected project and author identity references.
- Use project JSON wrappers from `common/json.go` for business-code marshal and unmarshal work.
- Keep changes scoped to the assigned phase.
- Update both `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json` before marking any feature complete.
- Include validation commands and results in the pull request body.

Required pull request handoff fields:

```json
{
  "development_log_updated": true,
  "mission_state_updated": true
}
```

## Phase 0: Current HEAD CI Evidence Closure

Goal: Record Pre-release verification `#16` on commit `73ad2ff` as the current trusted CI evidence.

Suggested branch:

```text
chore/ci-16-verification-closure
```

Tasks:

1. Update `docs/CI_VERIFICATION_EVIDENCE.md` with Pre-release verification `#16`, commit `73ad2ff`, branch `main`, and successful jobs.
2. Update audit and status documents that still point only to older CI run `#13` / commit `aeb43e5`.
3. Update `.factory/mission-state.json`:
   - `ci_verification_status = passed`
   - `ci_run_number = 16`
   - `ci_commit = 73ad2ff`
   - `production_readiness = not_ready`
   - `next_recommended_action = run isolated staging runtime verification and release sign-off`
4. Append a change record to `docs/DEVELOPMENT_LOG.md`.
5. Do not change business logic.

Validation:

```bash
node -e "JSON.parse(require('fs').readFileSync('.factory/mission-state.json','utf8'))"
git diff --check
bash scripts/check-config-secrets.sh
```

## Phase 1: Mission Status Reconciliation

Goal: Remove status drift between `features.json`, `validation-contract.md`, audit reports, development log, and mission-state.

Suggested branch:

```text
chore/reconcile-mission-status
```

Tasks:

1. Review `features.json`, `validation-contract.md`, `docs/CODEX_FEATURE_TEST_RESULTS.md`, `docs/CODEX_AUDIT_REPORT.md`, `docs/DEVELOPMENT_LOG.md`, and `.factory/mission-state.json`.
2. Mark completed features in `features.json`, or explicitly mark the file as superseded by mission-state and CI evidence.
3. Add a clear current-state note to `validation-contract.md`.
4. Mark historical failed audit sections as historical baseline where later remediation and CI evidence closed them.
5. Keep `production_readiness` as `not_ready`.
6. Update `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json`.

Validation:

```bash
node -e "JSON.parse(require('fs').readFileSync('features.json','utf8'))"
node -e "JSON.parse(require('fs').readFileSync('.factory/mission-state.json','utf8'))"
git diff --check
```

## Phase 2: Isolated Staging Runtime Verification

Goal: Turn CI success into staging runtime evidence.

Suggested branch:

```text
chore/staging-runtime-verification-report
```

Tasks:

1. Run the verification flow from `docs/STAGING_VERIFICATION_RUNBOOK.md` in an isolated staging or CI environment.
2. Use fake provider traffic only.
3. Do not mount production `.env` files.
4. Do not use real upstream provider keys.
5. Record command outputs and observations in `docs/STAGING_VERIFICATION_REPORT.md`.
6. Update `docs/PRE_DEPLOYMENT_REVIEW_CHECKLIST.md`, `docs/DEVELOPMENT_LOG.md`, and `.factory/mission-state.json`.

Required commands:

```bash
bash scripts/check-config-secrets.sh
bash scripts/ci-verify.sh
docker compose config
LOCAL_FIXTURE=1 bash scripts/regression.sh
bash scripts/ci-migration-check.sh
docker compose -f docker-compose.fixture.yml config
docker compose -f docker-compose.fixture.yml up -d --build
ADMIN_LINE="$(BASE_URL=http://localhost:3000 bash scripts/seed-local-fixture.sh | tail -n 1)"
eval "$ADMIN_LINE"
BASE_URL=http://localhost:3000 ADMIN_TOKEN="$ADMIN_TOKEN" ADMIN_USER_ID="$ADMIN_USER_ID" bash scripts/regression.sh
docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes
```

Required behavior evidence:

- `/v1/models` returns fixture-visible models.
- `/v1/chat/completions` succeeds through the fake OpenAI-compatible provider.
- Normal users cannot see or call `experimental_proxy`.
- Internal users with `allow_experimental=true` can call enabled `experimental_proxy`.
- Disabled `experimental_proxy` channels cannot route.
- Zero-balance requests reject before upstream routing.
- `usage_log` does not contain full prompt or response content by default.

## Phase 3: JSON Wrapper Governance

Goal: Enforce the project rule that business-code JSON marshal and unmarshal calls go through `common/json.go`.

### Phase 3A: Core Business Code

Suggested branch:

```text
refactor/json-wrapper-core
```

Scope:

- `controller/`
- `service/`
- `model/`
- `middleware/`

Tasks:

1. Replace direct `json.Marshal`, `json.Unmarshal`, and `json.NewDecoder(...).Decode` calls with `common.Marshal`, `common.Unmarshal`, and `common.DecodeJson` where they are business-code operations.
2. Preserve `encoding/json` imports where only JSON types or custom JSON interfaces are required.
3. Do not change provider adapter behavior in this phase.
4. Add or adjust targeted tests where behavior is not already covered.

Validation:

```bash
go test ./controller ./service ./model ./middleware -count=1
go test ./... -count=1
go vet ./...
git diff --check
```

### Phase 3B: Relay And Provider Paths

Suggested branch:

```text
refactor/json-wrapper-relay
```

Scope:

- `relay/`

Tasks:

1. Migrate relay business JSON parsing and encoding to `common` wrappers.
2. Preserve DTO `json.RawMessage` and custom marshal/unmarshal semantics.
3. Pay special attention to stream parsing, tool calls, provider response parsing, and error payloads.
4. Do not change upstream request/response protocol semantics.

Validation:

```bash
go test ./relay/... -count=1
go test ./... -count=1
LOCAL_FIXTURE=1 bash scripts/regression.sh
```

## Phase 4: CI Rule Enforcement

Goal: Prevent project-rule drift after JSON wrapper cleanup.

Suggested branch:

```text
ci/enforce-project-rules
```

Tasks:

1. Add `scripts/check-json-wrapper.sh`.
2. Detect direct business-code calls to `json.Marshal`, `json.Unmarshal`, and `json.NewDecoder(...).Decode` in controlled directories.
3. Add explicit allowlists for:
   - `common/json.go`
   - test files
   - DTO custom JSON methods where direct `encoding/json` is required
4. Wire the script into `scripts/ci-verify.sh`.
5. Wire the script into `.github/workflows/pre-release-verification.yml`.
6. Update `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json`.

Validation:

```bash
bash scripts/check-json-wrapper.sh
bash scripts/ci-verify.sh
```

## Phase 5: Release Sign-Off Pack

Goal: Prepare gray release materials without claiming production readiness too early.

Suggested branch:

```text
docs/release-signoff-pack
```

Tasks:

1. Update `docs/PRE_DEPLOYMENT_REVIEW_CHECKLIST.md` so every item has a concrete state:
   - `passed`
   - `blocked`
   - `manual_required`
2. Add or update `docs/RELEASE_SIGNOFF.md` with:
   - CI run `#16` evidence
   - staging runtime evidence
   - secret-source review
   - deployment topology review
   - experimental_proxy release policy
   - rollback plan
3. Keep `production_readiness = not_ready` unless a human release owner has explicitly signed off.
4. Do not include real secrets.
5. Update `docs/DEVELOPMENT_LOG.md` and `.factory/mission-state.json`.

Validation:

```bash
git diff --check
bash scripts/check-config-secrets.sh
```

## Phase 6: Post-MVP Feature Development

Only start this phase after Phases 0-5 have been completed or explicitly signed off.

Recommended order:

1. Enterprise tenant permission system
   - Expand Organization / Project from stubs into full RBAC.
   - Add project-level API keys, usage, quota, and member permissions.
   - Ensure cross-organization and cross-project access fails closed.

2. Provider management enhancement
   - Add richer Provider / ProviderAccount admin workflows.
   - Add health checks, cooldown, failure counters, account-level quota, and credential rotation.
   - Ensure credentials never appear in API responses, logs, or frontend state.

3. Dynamic billing expression system
   - Read `pkg/billingexpr/expr.md` before making any billing expression changes.
   - Add model, channel, provider, and user-group pricing expression support.
   - Keep pre-consume, settlement, usage logs, and display consistent.

4. Claude / Gemini / Responses API expansion
   - Implement each protocol as an independent milestone.
   - Keep OpenAI-compatible behavior non-regressing.
   - Confirm `StreamOptions` support per provider before adding any channel to `streamSupportedChannels`.

5. Production operations
   - Add operational audit logs.
   - Add alerting and rollback documentation.
   - Add backup and restore proof.
   - Add deployment health dashboards and runbooks.

## Recommended Immediate Remote Prompt

Start with Phase 0 and Phase 1. They do not change business logic, but they make the repository state trustworthy for later remote implementation work.

```text
Create branch chore/ci-16-verification-closure. Record Pre-release verification #16 success for current HEAD 73ad2ff, update CI evidence and mission-state, keep production_readiness not_ready, update DEVELOPMENT_LOG, run JSON parse, git diff --check, and config secret scan, then open a PR.
```
