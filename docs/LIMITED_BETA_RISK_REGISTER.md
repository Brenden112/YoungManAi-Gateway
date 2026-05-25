# Limited Beta Risk Register

Date: 2026-05-25

Current phase: `Phase 5 limited beta planning`
Deployment readiness: `limited_beta_ready`
Production readiness: `not_ready`

| Risk ID | Risk | Severity | Trigger | Mitigation | Stop Condition |
|---|---|---|---|---|---|
| LBR-001 | API key plaintext exposure | Critical | Full API key appears after one-time creation display, in logs, docs, screenshots, or API responses | Verify hash/prefix storage and masked responses before beta traffic | Any full key appears outside one-time creation response |
| LBR-002 | Provider credential leakage | Critical | Provider key appears in DB dumps, logs, UI, errors, docs, screenshots, or evidence | Use fake provider first; use only manually approved low-limit test key; verify encryption/redaction | Any provider credential is exposed |
| LBR-003 | Prompt or response leakage | High | `usage_log`, error log, document, issue, screenshot, or artifact contains full prompt/response | Keep prompt/response logging disabled; inspect daily logs and evidence | Full prompt or response is stored or published |
| LBR-004 | Normal user reaches `experimental_proxy` | Critical | Normal user sees or calls experimental model | Keep default disabled/internal-only; test visibility and routing with normal user | Normal user can see or call experimental path |
| LBR-005 | `official_cloud` falls back to `experimental_proxy` | Critical | Official provider failure routes to experimental provider | Test failure and disabled-channel paths; inspect selected provider type | Any official-to-experimental fallback occurs |
| LBR-006 | Zero-balance request calls upstream | High | Upstream receives request from zero/insufficient-balance user | Test balance rejection before real provider use | Any zero/insufficient-balance request reaches upstream |
| LBR-007 | Incorrect billing or balance deduction | High | Quota delta does not match token/cost calculation or upstream fee | Daily usage and balance reconciliation; low per-user caps | Incorrect debit/credit or wrong attribution is observed |
| LBR-008 | Disabled provider/channel/model still routes | High | Traffic reaches disabled object | Disable each object and retry with fake provider | Disabled route is selected |
| LBR-009 | `allowed_provider_types` bypass | High | API key routes to provider type outside its allowed list | Test disallowed provider-type API key before beta traffic | Provider type restriction is bypassed |
| LBR-010 | Organization/project authorization bypass | High | API key usage is accepted for wrong organization or project | Test scoped API keys and inspect `usage_log` attribution | Cross-org or cross-project call succeeds incorrectly |
| LBR-011 | Admin permission bypass | Critical | Normal user accesses admin page/API or performs admin action | Test normal-user dashboard and API access | Normal user can perform admin capability |
| LBR-012 | Real provider spend exceeds beta cap | High | Upstream fee or call volume exceeds approved low quota | Use low-limit test key, daily cap, and manual approval | Spend exceeds approved cap or high-privilege key is used |
| LBR-013 | Error log redaction gap | High | `params.Other`, upstream error, header, or credential appears in logs | Daily error log review and sanitized evidence rules | Token, credential, prompt, or response appears in error logs |
| LBR-014 | Provider failure rate hides routing issue | Medium | Elevated failures with unclear provider/channel attribution | Track provider failure rate by provider type and channel | Failure pattern indicates critical routing or billing risk |
| LBR-015 | Rollback cannot be executed quickly | High | Operator cannot disable traffic, revert artifact, or identify owner | Review rollback plan before beta; assign owner | Rollback is blocked during incident |

## Review Cadence

- Review before beta start.
- Review after first fake-provider beta pass.
- Review before any manually approved real `official_cloud` call.
- Review daily during beta.
- Review before production preparation decision.

Do not proceed to production preparation with unresolved critical or high risks.
