# Internal Gray Test Plan

Date: 2026-05-22

Current phase: `Phase 3 internal gray test planning`
Deployment readiness: `internal_gray_ready`
Production readiness: `not_ready`
Next recommended action: `execute internal gray test checklist`

## Scope

This plan defines the small-scope internal gray validation to run before limited beta or production preparation. It is a planning artifact only. It does not authorize production deployment, does not mark production readiness, and does not add product behavior.

## Hard Boundaries

- Do not use real high-privilege upstream API keys.
- Do not call real paid providers unless a human later supplies a low-limit staging-only test key and explicitly approves the call.
- Do not write real keys, tokens, credentials, prompts, responses, customer data, or production data into documentation, logs, test fixtures, screenshots, or evidence.
- Do not modify core business logic during this phase.
- Do not add business features during this phase.
- Keep `production_readiness = not_ready`.

## Test Objectives

- Verify the MVP is stable enough for small internal use.
- Verify user, API key, provider, billing, log, balance, admin, and `experimental_proxy` isolation behavior.
- Verify whether the release can proceed to limited beta or production preparation after internal gray execution and sign-off.

## Test Roles

| Role | Purpose |
|---|---|
| Admin user | Configure provider accounts, channels, API keys, balances, and review logs. |
| Normal user | Validate ordinary API access and prove experimental paths remain hidden and blocked. |
| Internal user | Validate explicitly opted-in experimental access. |
| Test organization | Validate org-bound API keys and usage attribution. |
| Test project | Validate project-bound API keys and usage attribution. |

All users, organizations, projects, tokens, and credentials must be staging-only and disposable.

## Test Providers

| Provider | Use in gray test |
|---|---|
| Fake provider | Required first path for all smoke and regression tests. |
| `official_cloud` test provider | Placeholder only by default. Use no real key unless a low-limit staging key is manually approved later. |
| `experimental_proxy` test provider | Placeholder, default disabled, internal-only, and never visible/callable to normal users. |
| `aggregator` test provider | Configurable for coverage, but no real call is required in this phase. |

## Required Test Areas

### User And API Key

- Create a normal user.
- Create an internal user.
- Create API keys for user, organization, and project scopes.
- Verify the full API key plaintext is displayed only once at creation.
- Verify only API key hash and prefix are stored or returned afterward.
- Verify a disabled API key cannot call relay APIs.
- Verify API key binding to user, org, and project is enforced.
- Verify `allowed_models` is enforced.
- Verify `allowed_provider_types` is enforced.

### OpenAI-Compatible API

- Verify `GET /v1/models`.
- Verify `POST /v1/chat/completions`.
- Verify non-streaming chat completion.
- If the selected channel supports streaming, verify streaming behavior.
- Verify OpenAI SDK compatibility with the staging endpoint and test key.

### Provider And Channel

- Verify `official_cloud` channels can be called through fake or approved low-limit staging provider only.
- Verify `aggregator` can be configured without requiring a real paid call.
- Verify `ProviderAccount` credentials are encrypted at rest.
- Verify `ProviderAccount` credentials decrypt only for runtime use.
- Verify disabled providers cannot route.
- Verify disabled channels cannot route.
- Verify legacy channel behavior remains compatible.

### Experimental Proxy

- Verify `experimental_proxy` defaults to disabled.
- Verify `experimental_proxy` defaults to internal-only.
- Verify normal users cannot see experimental models.
- Verify normal users cannot call experimental models.
- Verify internal users must have `allow_experimental=true` before calling experimental models.
- Verify disabled `experimental_proxy` cannot be called by any user.
- Verify `official_cloud` failure cannot fall back to `experimental_proxy`.
- Verify fallback, retry, preferred-channel, specific-channel, and legacy paths cannot bypass restrictions.

### Billing And Balance

- Verify successful requests record token counts.
- Verify successful requests calculate cost.
- Verify successful requests deduct balance.
- Verify failed requests do not deduct balance incorrectly.
- Verify zero-balance requests do not call upstream.
- Verify admin manual top-up works.
- Verify top-up creates an operation record.

### Logs And Privacy

- Verify successful requests write `usage_log`.
- Verify failed requests write `usage_log` or error log as designed.
- Verify `usage_log` includes `user_id`, `org_id`, `project_id`, and `api_key_id`.
- Verify `usage_log` includes `provider_type`, `channel_id`, `provider_account_id`, `model`, token counts, cost, and status.
- Verify full prompt content is not stored by default.
- Verify full response content is not stored by default.
- Verify `params.Other` is redacted.
- Verify `error_message` is redacted.
- Verify logs do not leak API keys, provider credentials, bearer tokens, prompts, or responses.

### Admin Dashboard

- Verify admin can view providers.
- Verify admin can view channels.
- Verify admin can view API keys.
- Verify admin can disable API keys.
- Verify admin can view `usage_log`.
- Verify admin can view user balance.
- Verify admin can perform manual top-up.
- Verify normal users cannot access admin pages or admin APIs.

### Deployment And Regression

- Run `bash scripts/check-config-secrets.sh`.
- Run `bash scripts/ci-verify.sh`.
- Run `LOCAL_FIXTURE=1 bash scripts/regression.sh`.
- Run Docker compose fixture smoke.
- Review GitHub Actions pre-release verification.
- Attach staging verification evidence.

## Suggested Schedule

| Day | Focus |
|---|---|
| Day 1 | Fake provider internal validation. |
| Day 2 | Low-limit `official_cloud` test provider validation only if a staging key is manually approved. |
| Day 3 | Admin, logs, billing, balance, and error scenario validation. |
| Day 4-7 | Low-traffic stability observation with internal users only. |

## Evidence Rules

Record command names, exit codes, commit SHA, environment name, sanitized observations, issue links, and sign-off decisions. Do not record real credentials, raw bearer tokens, raw API keys, real prompts, real responses, production data, or customer data.

## Exit Decision

Use `docs/INTERNAL_GRAY_EXIT_CRITERIA.md` for the formal proceed/stop criteria. Passing this plan can support limited beta or production preparation only; it must not set `production_readiness` to ready.
