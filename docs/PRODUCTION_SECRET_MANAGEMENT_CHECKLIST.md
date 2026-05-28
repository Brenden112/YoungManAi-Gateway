# Production Secret Management Checklist

Date: 2026-05-28
Phase: `Phase 7 production preparation sign-off pack`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`

No real provider key, credential, token, prompt, or response was used to prepare this checklist.

## Secret Sources

- [ ] Production secrets are stored in an approved secret manager or deployment platform secret store.
- [ ] `.env.production` is not committed to git and is not copied into docs, test artifacts, screenshots, or issue comments.
- [ ] Real provider keys are never embedded in frontend code, build-time public variables, fixtures, tests, or documentation.
- [ ] Secret access is limited to named production operators.
- [ ] Secret rotation owner and rotation cadence are assigned.

## Required Production Secrets

- [ ] `CRYPTO_SECRET`: strong random value, production-specific, backed up only through approved secret process.
- [ ] `SESSION_SECRET`: strong random value, production-specific, rotated from all examples.
- [ ] Database password: strong random value, production-specific, least-privilege account.
- [ ] Redis password: strong random value if Redis auth is enabled.
- [ ] Provider credentials: stored encrypted through ProviderAccount or equivalent secret path.
- [ ] OAuth/OIDC/WebAuthn secrets: production-specific and stored outside git.

## Application Storage Rules

- [ ] ProviderAccount credentials are encrypted at rest.
- [ ] Channel legacy keys are not used when a linked encrypted ProviderAccount is configured.
- [ ] API keys are persisted as hash/prefix only.
- [ ] Full API key value is shown once and is not recoverable from list/search endpoints.
- [ ] Bearer tokens and session values are redacted from logs.
- [ ] Prompt/response full text remains disabled by default with `STORE_FULL_TEXT_ENABLED=false`.

## Release Checks

- [ ] `bash scripts/check-config-secrets.sh` passes before release.
- [ ] Git status and review confirm no `.env`, `.env.*`, `secrets/`, private keys, certificates, provider keys, or generated secret dumps are staged.
- [ ] Production deployment variables are reviewed in the secret store, not copied into review docs.
- [ ] Rollback plan includes key revocation and rotation if exposure is suspected.

## LBI-003 Constraint

LBI-003 stays `manual_required` until a human release owner supplies and approves a low-limit real `official_cloud` provider key for Phase 6C. Until then, no real provider key should be present in local docs, logs, tests, or this repository.
