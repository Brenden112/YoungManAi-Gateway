# Production Release Risk Register

Date: 2026-05-28
Phase: `Phase 7 production preparation sign-off pack`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`

## Summary

Critical findings: `0`
High findings: `0`
LBI-003 status: `manual_required`

This register tracks production-preparation risks. It does not authorize production deployment.

| ID | Severity | Status | Risk | Mitigation | Owner |
|---|---|---|---|---|---|
| PRR-001 | Manual gate | `open` | No low-limit real `official_cloud` beta has been run. | Keep LBI-003 `manual_required`; run Phase 6C only with human-approved low-limit key and quota cap. | Release owner |
| PRR-002 | High if missed | `manual_required` | Real production secrets may not be configured or may be weak. | Configure strong `CRYPTO_SECRET`, `SESSION_SECRET`, database password, Redis password, and provider credentials in secret store only. | Release owner |
| PRR-003 | High if missed | `manual_required` | Database backup/restore may be unverified. | Configure encrypted backups and perform restore test before release. | Operations owner |
| PRR-004 | High if missed | `manual_required` | Monitoring and alert routing may be incomplete. | Implement alerts in `docs/PRODUCTION_MONITORING_AND_ALERTING_PLAN.md` and assign on-call owner. | Operations owner |
| PRR-005 | High if missed | `manual_required` | Rollback may not be executable under incident pressure. | Review `docs/PRODUCTION_ROLLBACK_PLAN.md`, keep previous artifact available, and assign rollback owner. | Release owner |
| PRR-006 | Critical if regressed | `controlled_by_policy` | `experimental_proxy` becomes visible or callable by normal users. | Keep disabled/internal-only default, `allow_experimental=false`, and alert on normal-user experimental calls. | Security owner |
| PRR-007 | Critical if regressed | `controlled_by_policy` | `official_cloud` falls back to `experimental_proxy`. | Keep provider-type guards on model listing, retry, fallback, preferred channel, and selected-channel setup. | Provider owner |
| PRR-008 | Critical if regressed | `controlled_by_policy` | Zero-balance request reaches upstream provider. | Monitor zero-balance bypass, verify pre-consume quota rejection, rollback on any bypass. | Billing owner |
| PRR-009 | High if regressed | `controlled_by_policy` | Prompt/response, API key, bearer token, or provider credential leaks into logs. | Keep `STORE_FULL_TEXT_ENABLED=false`, sanitizer enabled, secret scan passing, and incident rotation plan ready. | Security owner |
| PRR-010 | Medium | `manual_required` | Fake-provider beta may not represent real provider billing and failure behavior. | Release owner must decide whether fake-provider beta is sufficient for production preparation or run Phase 6C. | Release owner |

## LBI-003 Gate Impact

LBI-003 does not block this production preparation pack. It blocks `production_ready` because no human-approved low-limit real `official_cloud` provider key was supplied and no real paid provider was called. Clearing it requires Phase 6C low-limit real provider beta.
