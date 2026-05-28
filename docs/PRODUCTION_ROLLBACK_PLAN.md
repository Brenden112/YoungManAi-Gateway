# Production Rollback Plan

Date: 2026-05-28
Phase: `Phase 7 production preparation sign-off pack`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`

This rollback plan must be reviewed and owned before production deployment. It is a preparation artifact only.

## Rollback Triggers

Rollback or disable affected capability immediately if any of the following occur:

- Elevated 5xx rate or sustained latency regression.
- Provider failure rate exceeds the agreed threshold.
- Billing, token usage, or balance deduction anomaly.
- Zero-balance or insufficient-balance request reaches an upstream provider.
- Normal user can see or call `experimental_proxy`.
- `official_cloud` falls back to `experimental_proxy`.
- Prompt, response, API key, bearer token, provider credential, or database credential leakage.
- Database migration, connection, or data-integrity incident.
- Redis outage affects auth, rate limiting, routing, or billing behavior.

## Version Rollback

1. Freeze production changes and assign an incident owner.
2. Disable risky provider/channel paths if needed before version rollback.
3. Redeploy the previous known-good image or artifact.
4. Confirm `/api/status`, login, model listing, chat completion, billing log write, and admin access on the rolled-back version.
5. Keep `production_readiness` blocked until post-incident review is complete.

## Database Backup Restore

1. Stop write traffic or put the service into maintenance mode.
2. Snapshot the current database before restore for forensic retention.
3. Restore the latest verified pre-deployment backup.
4. Run migration/schema compatibility checks for SQLite, MySQL, or PostgreSQL as applicable.
5. Validate user balances, token records, provider/channel records, and usage-log continuity.

## Provider And Channel Controls

- Disable the affected ProviderAccount, provider, channel, or model immediately from admin controls.
- Confirm disabled channels cannot route through cache, preferred-channel, fallback, or retry paths.
- Switch traffic to a known-good provider only if provider type and allowed model policy still apply.
- Keep real provider keys out of incident notes, screenshots, and logs.

## Experimental Proxy Kill Switch

- Set `experimental_proxy` channels disabled or internal-only.
- Confirm `allow_experimental=false` for normal tokens.
- Confirm normal model list hides experimental-only models.
- Confirm direct normal-user experimental calls fail closed.

## API Key Disable

- Disable compromised or suspicious API keys immediately.
- Confirm disabled keys return unauthorized responses and cannot reach upstream.
- Rotate affected user credentials and notify affected owners according to incident policy.

## Billing And Balance Recovery

- Freeze affected user balance changes if overcharge, undercharge, or misattribution is suspected.
- Export sanitized usage IDs and quota deltas only.
- Reconcile system deductions against upstream charges without storing real prompts or responses.
- Apply manual corrections through audited admin top-up or deduction records.

## Log Leakage Response

- Stop the leaking log path or disable full-text storage.
- Rotate exposed API keys, provider keys, session secrets, or encryption secrets as applicable.
- Remove leaked material from accessible log sinks and preserve restricted forensic copies.
- Record the incident without copying secrets, prompts, or responses into the report.

## Upstream Provider Failure

- Disable the failing provider/channel.
- Route only to approved providers of the same allowed provider type.
- Do not allow `official_cloud` traffic to fall back to `experimental_proxy`.
- Monitor failure rate, token usage, and billing deltas during provider recovery.
