# Internal Gray Exit Criteria

Date: 2026-05-22

Current phase: `Phase 3 internal gray test planning`
Deployment readiness: `internal_gray_ready`
Production readiness: `not_ready`

These criteria govern whether internal gray execution may proceed to limited beta or production preparation. They do not authorize production deployment.

## Proceed Criteria

All of the following must be true before moving to the next stage:

- Critical issue count is `0`.
- High issue count is `0`.
- API call success rate meets the agreed internal target for the test window.
- Balance deduction is correct for successful requests.
- Failed requests do not deduct balance incorrectly.
- Zero-balance requests do not call upstream.
- Normal users cannot see or call `experimental_proxy`.
- Internal users can call experimental paths only when explicitly allowed with `allow_experimental=true`.
- Disabled `experimental_proxy`, disabled providers, and disabled channels cannot route.
- `official_cloud` failures do not fall back to `experimental_proxy`.
- Logs do not store full prompts by default.
- Logs do not store full responses by default.
- Logs do not leak API keys, bearer tokens, or provider credentials.
- API keys are not exposed in plaintext after one-time creation display.
- API key user/org/project bindings are enforced.
- `allowed_models` and `allowed_provider_types` are enforced.
- Admin provider, channel, API key, usage log, balance, and top-up operations are usable.
- Normal users cannot access admin pages or admin APIs.
- Rollback path, rollback owner, and rollback approval path are documented.
- GitHub Actions pre-release verification remains passed or any new failure is resolved before exit.
- Staging verification evidence is attached.

## Stop Criteria

Stop internal gray execution immediately if any of the following occur:

- API key leakage.
- Provider credential leakage.
- Bearer token leakage.
- Real prompt or response stored in logs, docs, fixtures, screenshots, or evidence.
- Normal user can call `experimental_proxy`.
- Normal user can see experimental-only models.
- Disabled `experimental_proxy` can route.
- `official_cloud` falls back to `experimental_proxy`.
- Zero-balance request calls upstream.
- Billing error causes incorrect balance deduction or credit.
- Admin permission bypass.
- Normal user can access admin pages or APIs.
- Database migration corrupts data or breaks user/org/project ownership.
- Usage logs attribute traffic to the wrong user, organization, project, API key, provider, channel, or provider account.
- Rollback cannot be executed or has no accountable owner.
- Any critical or high security issue is discovered.

## Required Exit Record

Before the next phase, record:

- Commit SHA and environment identifier.
- Checklist result summary.
- Issue list with severity and disposition.
- API success-rate summary.
- Billing/balance verification result.
- Experimental isolation verification result.
- Log privacy verification result.
- Admin verification result.
- Rollback plan and owner.
- Human decision: continue to limited beta, continue to production preparation, extend gray test, or stop.

`production_readiness` must remain `not_ready` until a separate production readiness review and human release sign-off explicitly change it.
