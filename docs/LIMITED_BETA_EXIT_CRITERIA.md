# Limited Beta Exit Criteria

Date: 2026-05-25

Current phase: `Phase 5 limited beta planning`
Deployment readiness: `limited_beta_ready`
Production readiness: `not_ready`

These criteria govern whether limited beta may proceed to production preparation. They do not authorize production deployment.

## Continue Beta Criteria

Limited beta may continue only while all of the following remain true:

- No API key, bearer token, provider credential, or sensitive tenant data is exposed.
- `experimental_proxy` remains disabled and internal-only by default.
- Normal users cannot see or call `experimental_proxy`.
- `official_cloud` never falls back to `experimental_proxy`.
- Zero-balance and insufficient-balance users are rejected before upstream calls.
- Disabled providers, disabled channels, disabled models, and disallowed provider types cannot route.
- Organization and project token bindings are enforced.
- Prompt and response bodies are not saved by default.
- Billing, token count, and balance deductions are explainable for every sampled call.
- Daily usage, error, balance, provider failure, fallback, privacy, and reconciliation checks are completed.

## Stop Beta Immediately

Stop limited beta immediately if any of the following occur:

- API key leakage.
- Provider credential leakage.
- Bearer token leakage.
- Normal user can call `experimental_proxy`.
- Normal user can see experimental-only models.
- Zero-balance or insufficient-balance request reaches upstream.
- Billing or balance deduction is incorrect in a way that can overcharge, undercharge, or misattribute usage.
- `usage_log`, error log, screenshot, issue, document, fixture, or test artifact stores full prompt or full response.
- Admin permission bypass.
- Normal user can access admin pages or admin APIs.
- Organization or project authorization bypass.
- Fallback bypasses `provider_type`, `allowed_provider_types`, disabled provider, disabled channel, or disabled model restrictions.
- `official_cloud` falls back to `experimental_proxy`.
- A real high-privilege provider key is configured or used.
- An unapproved real provider call occurs.
- Rollback cannot be executed by the assigned owner.

## Enter Production Preparation Criteria

All of the following must be true before moving from limited beta to production preparation:

- Beta cycle completed for the approved test window or was explicitly closed by release-owner decision.
- Critical issue count is `0`.
- High issue count is `0`.
- Medium issues are assessed, documented, and either fixed, accepted with owner/date, or deferred with release-owner approval.
- Billing and balance reconciliation passed for fake-provider and any approved real-provider calls.
- Upstream real provider charges match system deductions within the accepted low-quota beta tolerance.
- Log privacy review passed with no prompt, response, key, token, or credential leakage.
- `experimental_proxy` isolation review passed.
- `official_cloud` to `experimental_proxy` fallback prohibition passed.
- Rollback plan was reviewed and verified in the beta environment.
- Admin dashboard checks passed.
- Organization/project token binding checks passed.
- Manual release-owner sign-off is recorded.

## Required Exit Record

Before production preparation, record:

- Commit SHA and environment identifier.
- Beta dates and traffic totals.
- User, organization, and project counts.
- Provider types used and whether any real provider was approved.
- Checklist summary from `docs/LIMITED_BETA_CHECKLIST.md`.
- Issue list with severity, owner, status, and disposition.
- Billing and balance reconciliation result.
- Log privacy review result.
- Experimental isolation review result.
- Rollback verification result.
- Release-owner sign-off decision.

`production_readiness` must remain `not_ready` until a separate production readiness review changes it.
