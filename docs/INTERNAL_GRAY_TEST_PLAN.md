# Internal Gray Test Plan

Date: 2026-05-20

Deployment readiness: `internal_gray_ready`
Production readiness: `not_ready`

## Purpose

This plan defines the next controlled internal gray test after Phase 2 isolated staging runtime verification passed in GitHub Codespaces. The gray test is intended to validate deployment wiring, account setup, provider policy, billing guards, logs, and rollback behavior in an isolated staging environment before any production release decision.

This document does not mark the project as production ready.

## Boundaries

- Use an isolated staging host, staging database, and staging Redis only.
- Do not mount or copy production `.env` files.
- Do not use production customer data, real prompts, or real responses.
- Do not commit `.env`, `.env.*`, key files, certificates, or controlled staging secrets.
- Use fake-provider traffic for the required smoke path.
- Optional real-provider smoke is allowed only if a human release owner supplies an approved staging-only key through an untracked secret source and explicitly approves the provider call.
- Keep `production_readiness` as `not_ready` until human release sign-off is recorded separately.

## Preconditions

| Item | Required state |
|---|---|
| Pre-release verification | GitHub Actions Pre-release verification #16 on commit `73ad2ff` passed. |
| Phase 2 isolated staging runtime verification | Passed in GitHub Codespaces. |
| Docker fixture runtime | `passed_in_codespaces`. |
| Go/runtime blockers | `closed_in_codespaces`. |
| API staging blocker | `closed_in_codespaces`. |
| Frontend local script blocker | `non_blocking_note_ci_frontend_check_passed`. |
| Deployment readiness | `internal_gray_ready`. |
| Production readiness | `not_ready`. |
| Secret source | Manual required: controlled staging secrets must be supplied outside git. |
| Staging env file | Manual required: copy `.env.staging.example` to untracked `.env.staging`. |

## Test Environment

The internal gray environment must use:

- One isolated application deployment built from the reviewed branch or commit.
- Dedicated staging database with no production data.
- Dedicated staging Redis instance.
- Staging-only `SESSION_SECRET` and `CRYPTO_SECRET`.
- Fake upstream provider endpoint for required smoke tests.
- Optional staging-only real provider key only after explicit release-owner approval.

## Test Accounts And Tokens

Create or verify these staging-only identities:

| Identity | Purpose |
|---|---|
| Admin user | Configure channels, provider accounts, user quota, and review logs. |
| Normal user | Validate official provider access and verify `experimental_proxy` remains hidden and blocked. |
| Internal test user | Validate explicitly approved `experimental_proxy` behavior with `allow_experimental=true`. |
| Zero-balance user | Validate insufficient balance rejects before upstream routing. |
| Normal user token | Validate default `allow_experimental=false` behavior. |
| Internal user token | Validate opt-in `allow_experimental=true` behavior. |

All accounts and tokens must be staging-only and disposable.

## Execution Steps

1. Secret and configuration review
   - Run `bash scripts/check-config-secrets.sh`.
   - Confirm `.env.staging` is untracked.
   - Confirm `.env`, `.env.*`, `secrets/`, `*.key`, `*.pem`, `*.p12`, and `*.pfx` are not committed.
   - Confirm `CRYPTO_SECRET`, `SESSION_SECRET`, database credentials, and Redis credentials are staging-only.

2. Deploy isolated staging
   - Deploy the reviewed commit or branch.
   - Start application, database, and Redis services.
   - Confirm health and application startup logs do not print secrets.

3. Migration and service verification
   - Run migration verification appropriate for the staging database.
   - Run `docker compose config` or the equivalent environment-render check for the staging topology.
   - Confirm the application can connect to database and Redis.

4. Required fake-provider smoke
   - Use the fake-provider fixture path first.
   - Verify `/v1/models` returns only models visible to the caller.
   - Verify `/v1/chat/completions` succeeds through the fake OpenAI-compatible provider.
   - Verify no real upstream provider or API key is used.

5. Provider policy checks
   - Normal users cannot see `experimental_proxy` models.
   - Normal users cannot call `experimental_proxy`.
   - Internal users can call enabled `experimental_proxy` only with `allow_experimental=true`.
   - Disabled `experimental_proxy` channels cannot route.
   - `allowed_provider_types` blocks disallowed route candidates.
   - Official provider traffic does not fall back to `experimental_proxy`.

6. Billing and quota checks
   - Zero-balance requests reject before upstream routing.
   - Successful fixture requests create usage logs with token/quota fields.
   - Manual top-up changes staging user quota only.
   - Failed or rejected requests do not deduct quota incorrectly.

7. Log and privacy checks
   - Usage logs do not contain full prompt or response content by default.
   - Error logs redact provider credentials and sensitive payload fields.
   - Admin list/search responses expose key prefixes only, not full API keys.

8. Frontend and admin checks
   - Confirm GitHub Actions `frontend-check` evidence is attached to the gray record.
   - Validate admin pages for provider accounts, channels, API keys, usage logs, and balance workflows.
   - Record the frontend local script blocker only as a non-blocking note unless it fails in the release CI workflow.

9. Rollback drill
   - Document the exact previous artifact or branch to restore.
   - Confirm database migration rollback or forward-fix decision path.
   - Confirm Redis/session cleanup steps.
   - Confirm who can execute rollback and who approves it.

10. Evidence capture and sign-off
    - Save command names, exit codes, environment name, commit SHA, and operator notes.
    - Do not store secrets, real prompts, real responses, or customer data in evidence.
    - Record pass/fail results in the staging report or release sign-off document.

## Pass Criteria

- Required fake-provider smoke passes.
- Security and provider policy checks pass.
- Zero-balance rejection happens before upstream routing.
- Logs do not store prompt/response content by default.
- No committed secret or real upstream key is found.
- Rollback owner and rollback path are documented.
- Production readiness remains `not_ready`.

## Fail Or Stop Criteria

Stop the gray test and do not proceed toward release if any of these occur:

- A real secret, real provider key, production prompt, production response, or customer data appears in committed files or test evidence.
- Normal users can see or call `experimental_proxy`.
- Disabled experimental channels can route.
- Zero-balance requests reach upstream routing.
- Provider credentials appear in API responses, frontend state, logs, or evidence.
- Staging deployment cannot be rolled back or safely destroyed.

## Evidence To Record

Record the following after execution:

- Commit SHA and branch under test.
- Staging environment identifier.
- Commands run and exit codes.
- Fake-provider smoke result.
- Provider policy check result.
- Billing/quota check result.
- Log privacy check result.
- Rollback drill result.
- Non-blocking notes, including the frontend local script note if still present.
- Human release-owner decision.

## Next Recommended Action

Execute this internal gray test plan with controlled staging secrets and fake-provider traffic first. Keep production readiness as `not_ready` until the release sign-off pack and human release-owner approval are complete.
