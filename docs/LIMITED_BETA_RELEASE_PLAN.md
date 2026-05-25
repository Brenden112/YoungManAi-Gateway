# Limited Beta Release Plan

Date: 2026-05-25

Current phase: `Phase 5 limited beta planning`
Deployment readiness: `limited_beta_ready`
Production readiness: `not_ready`
Next recommended action: `execute limited beta checklist`

## Purpose

This plan defines a small, real-world, low-quota beta validation after isolated staging and internal gray runtime checks passed. It is a planning artifact only. It does not authorize production deployment, does not mark production readiness, and does not add or modify business logic.

## Hard Boundaries

- Do not use real high-privilege provider keys.
- Use real providers only with low-limit test keys after manual release-owner confirmation.
- Do not write real keys, bearer tokens, credentials, prompts, responses, customer data, or production data into documentation, logs, fixtures, screenshots, tests, or evidence.
- Keep `experimental_proxy` disabled and `internal_only` by default.
- Keep `production_readiness = not_ready`.
- Do not allow `official_cloud` fallback to `experimental_proxy`.

## Beta Scope

| Area | Limit |
|---|---|
| Test users | Up to 5 named users total: 1 admin, up to 3 normal users, and up to 1 internal user. |
| Test organizations | Up to 2 staging-only organizations. |
| Test projects | Up to 3 staging-only projects across the beta organizations. |
| Test cycle | 7 calendar days, extendable once by release-owner approval if medium issues need observation. |
| Traffic limit | Maximum 300 total relay calls per day across all beta users. |
| Per-user balance | Maximum low test balance equivalent to USD 5 per user or stricter local quota policy. |
| Daily call limit | Maximum 100 relay calls per user per day. |
| Real provider use | Optional. Fake provider remains required. Real `official_cloud` use requires low-limit test key and manual confirmation before first call. |

## Provider Strategy

| Provider type | Beta policy |
|---|---|
| Fake provider | Required and retained for smoke, regression, isolation, billing, and negative-path checks. |
| `official_cloud` | May be connected to a low-limit test key only after manual approval. High-privilege or production keys are forbidden. |
| `aggregator` | Not mandatory for limited beta. Configure only if there is a specific low-risk test objective. |
| `experimental_proxy` | Default disabled, internal-only, and hidden from normal users. |

Additional routing requirements:

- Only internal users may test `experimental_proxy`, and only after explicit enablement for that internal user and provider/channel.
- Normal users must not see or call experimental models.
- Disabled providers, disabled channels, disabled models, and disallowed provider types must not route.
- `official_cloud` errors, capacity failures, or disabled states must not fallback to `experimental_proxy`.

## Safety Requirements

- API keys must not be stored in plaintext after one-time creation display.
- `ProviderAccount` credentials must be encrypted at rest.
- Prompt and response bodies must not be saved by default.
- `params.Other`, headers, bearer tokens, credentials, prompt fragments, response fragments, and upstream errors must be redacted before logging or evidence capture.
- Zero-balance users must be rejected before any upstream call.
- Disabled provider, channel, or model records must not be callable.
- `allowed_provider_types` must be enforced for all API keys and relay paths.
- Organization and project token binding must be enforced and reflected in usage attribution.

## Execution Model

1. Start with fake provider only.
2. Complete all critical checklist items for user, organization, project, API key, billing, logging, and admin controls.
3. If fake-provider results are clean, request manual approval for a low-limit `official_cloud` test key.
4. Run a small approved `official_cloud` test set with sanitized prompts and no sensitive data.
5. Observe logs, balances, provider failures, fallback behavior, and cost reconciliation daily.
6. Stop immediately if any stop condition in `docs/LIMITED_BETA_EXIT_CRITERIA.md` occurs.

## Evidence Rules

Record only sanitized evidence:

- Commit SHA.
- Environment identifier.
- Checklist item ID and pass/fail status.
- Command or UI action name.
- HTTP status and sanitized request category.
- Sanitized usage totals, token totals, and cost deltas.
- Issue IDs and severity.
- Release-owner decisions.

Do not record raw API keys, bearer tokens, provider credentials, real prompts, real responses, customer data, or production data.

## Monitoring And Observation

- Review `usage_log` daily for attribution, token count, cost, status, provider type, channel, and provider account consistency.
- Review error logs daily for redaction and unexpected failures.
- Review balances daily for negative, unchanged, or unexpectedly changed balances.
- Track provider failure rate by provider type and channel.
- Confirm fallback behavior never crosses provider-type restrictions.
- Compare token and cost calculations against provider-reported usage.
- Reconcile upstream real provider charges against system deductions for any approved real-provider calls.
- Check daily for prompt or response leakage in logs, docs, screenshots, issue comments, and test artifacts.

## Production Boundary

Limited beta may produce evidence for production preparation, but it must not move the project into production. Production preparation requires the conditions in `docs/LIMITED_BETA_EXIT_CRITERIA.md` and a separate human release-owner sign-off.
