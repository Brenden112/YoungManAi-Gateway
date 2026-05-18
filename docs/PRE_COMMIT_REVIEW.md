# Pre-Commit Review

Date: 2026-05-18

Status: `ready_for_ci_handoff`

Deployment readiness: `needs_manual_review`

## What This Work Completed

- Completed the full remediation and pre-release hardening documentation trail.
- Added CI/staging verification assets:
  - `scripts/ci-verify.sh`
  - `scripts/ci-migration-check.sh`
  - `.github/workflows/pre-release-verification.yml`
  - `docs/STAGING_VERIFICATION_RUNBOOK.md`
- Moved the remaining environment blockers to `pending_ci_verification` instead of marking them passed.
- Kept all provider smoke paths scoped to fake/local fixtures. No real upstream provider key, real token, credential, prompt, or response is required.

## Fixed Critical / High Security Issues

The audit status records zero remaining critical and high findings. Fixed areas include:

- API key storage hardening: token auth uses non-secret hash/prefix behavior, with full key display limited to one-time creation flow.
- ProviderAccount credential encryption and relay wiring.
- Disabled `experimental_proxy` route blocking.
- `allowed_provider_types` enforcement through model listing, selection, fallback, retry, preferred channels, and final setup paths.
- Organization/project token binding enforcement.
- Log privacy sanitization so prompt/response and sensitive metadata are not persisted by default.
- Go vet issues tracked during remediation were previously resolved in the remediation record.

## Blockers Still Requiring CI Verification

These are not marked passed:

- `blocked_test_infra_frontend` -> `pending_ci_verification`
- `blocked_external_dependency_cross_db_runtime` -> `pending_ci_verification`
- `skipped_environment_docker_runtime` -> `pending_ci_verification`
- `final_go_verification_blocked` -> `pending_ci_verification`

## Local Checks Passed

- `git diff --check`
- `.factory/mission-state.json` JSON parse
- `docs/CODEX_AUDIT_REPORT.md` contains `pending_ci_verification`
- `docs/PRE_DEPLOYMENT_REVIEW_CHECKLIST.md` exists
- `docs/STAGING_VERIFICATION_RUNBOOK.md` exists
- `.github/workflows/pre-release-verification.yml` exists
- `scripts/ci-verify.sh` exists, is executable, and documents behavior through clear pass/fail/blocked output
- `scripts/ci-migration-check.sh` exists, is executable, and documents SQL_DSN-based CI usage
- `docker compose config`
- `bash -n scripts/ci-verify.sh`
- `bash -n scripts/ci-migration-check.sh`
- `bash -n scripts/regression.sh`

## Local Checks Not Executed

- Go tests and `go vet` were not rerun locally because the current shell has no Go toolchain.
- `LOCAL_FIXTURE=1 bash scripts/regression.sh` remains blocked locally because it requires Go.
- Docker runtime smoke was not run locally because Docker runtime/cleanup previously failed with `Failed to initialize: protocol not available`.
- Frontend lint/test/build were not run locally because dependencies are absent and no dependency installation was performed.

## Required CI Workflow After Push

Run:

- `.github/workflows/pre-release-verification.yml`

Required jobs:

- `go-test-vet`
- `local-fixture-regression`
- `cross-db-migration`
- `docker-fixture-smoke`
- `frontend-check`

## If CI Fails, Check These Jobs First

1. `go-test-vet`: closes `final_go_verification_blocked` for backend test/vet evidence.
2. `local-fixture-regression`: verifies local fake-provider regression behavior without real upstream calls.
3. `cross-db-migration`: closes MySQL/PostgreSQL migration runtime evidence.
4. `docker-fixture-smoke`: closes seeded Docker runtime smoke and fixture-scoped cleanup evidence.
5. `frontend-check`: closes Bun/frontend lint, test, and build evidence.

## Current Git Summary

- Branch `main` is behind `origin/main` by 12 commits.
- Worktree contains modified tracked files and many untracked remediation, fixture, test, docs, script, and workflow files.
- `git diff --stat` for tracked files reports 67 files changed, 1854 insertions, 397 deletions.

## Release Gate

Do not move `deployment_readiness` beyond `needs_manual_review` until the pre-release verification workflow runs in CI/staging and the pending blockers are closed or explicitly signed off by a release owner.
