# Production Deployment Checklist

Date: 2026-05-28
Phase: `Phase 7 production preparation sign-off pack`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`

This checklist must be completed by the release owner before any production deployment. It is not a deployment authorization.

## Environment And Secrets

- [ ] `.env.production` exists only in the deployment secret store or production host and is not committed to git.
- [ ] `CRYPTO_SECRET` is strong, random, production-specific, and unavailable to frontend code.
- [ ] `SESSION_SECRET` is strong, random, production-specific, and unavailable to frontend code.
- [ ] Database password is strong, random, production-specific, rotated from examples, and stored outside git.
- [ ] Redis password, if used, is strong, random, production-specific, and stored outside git.
- [ ] Real provider keys are injected from a controlled secret source only.
- [ ] No real key, token, credential, prompt, or response appears in docs, logs, tests, issue comments, screenshots, or frontend bundles.

## Security And Privacy

- [ ] ProviderAccount credentials are encrypted at rest.
- [ ] API keys are stored as hash/prefix only.
- [ ] API keys are displayed in full only once at creation time.
- [ ] `STORE_FULL_TEXT_ENABLED=false`.
- [ ] `ERROR_LOG` output is sanitized and does not include provider credentials, bearer tokens, API keys, prompts, or responses.
- [ ] `usage_log` does not save full prompt or full response content by default.
- [ ] `params.Other` and error messages are redacted before persistence.
- [ ] Admin access has a named owner and audited access path.

## Provider And Routing Policy

- [ ] `experimental_proxy` is disabled or internal-only by default.
- [ ] `allow_experimental` defaults to `false`.
- [ ] Normal users cannot see experimental-only models.
- [ ] Normal users cannot call `experimental_proxy`.
- [ ] `official_cloud` cannot fall back to `experimental_proxy`.
- [ ] Disabled providers, disabled channels, disabled models, and disallowed provider types cannot route.
- [ ] `allowed_provider_types` is enforced for model listing, preferred channel, retry, fallback, and selected-channel setup.
- [ ] Organization/project token binding is enforced and spoofed request context is ignored.

## Billing And Balance

- [ ] Zero-balance and insufficient-balance requests do not call upstream providers.
- [ ] Successful usage records provider, channel, model, token counts, quota/cost, user, token, organization, and project attribution.
- [ ] Failed requests do not incorrectly deduct balance.
- [ ] Admin top-up creates an auditable operation record.
- [ ] Billing reconciliation owner and review cadence are assigned.

## Operational Readiness

- [ ] Database backup is configured, encrypted, monitored, and restore-tested.
- [ ] Rollback version and previous image/artifact are available.
- [ ] Provider/channel disable switches are accessible to the release owner.
- [ ] API key disable path is verified.
- [ ] Monitoring dashboards and alert routing are configured.
- [ ] On-call owner and escalation path are assigned.
- [ ] LBI-003 manual provider gate is either left as `manual_required` with production not ready, or cleared by Phase 6C low-limit real provider beta.

## Final Gate

- [ ] Release owner has reviewed `docs/PRODUCTION_PREPARATION_SIGNOFF.md`.
- [ ] Release owner has reviewed `docs/PRODUCTION_RELEASE_RISK_REGISTER.md`.
- [ ] Release owner explicitly confirms whether fake-provider beta evidence is sufficient to continue production preparation.
- [ ] Release owner explicitly confirms `production_readiness` may move only after a separate production readiness review.
