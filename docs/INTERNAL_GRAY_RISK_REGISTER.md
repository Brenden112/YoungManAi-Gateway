# Internal Gray Risk Register

Date: 2026-05-22

Current phase: `Phase 3 internal gray test planning`
Deployment readiness: `internal_gray_ready`
Production readiness: `not_ready`

| Risk ID | Risk | Severity | Trigger | Mitigation | Stop Condition |
|---|---|---|---|---|---|
| IGR-001 | API key plaintext exposure after creation | Critical | API response, admin UI, log, or evidence shows full key after one-time display | Verify hash/prefix storage and masked responses before gray traffic | Any full key appears outside one-time creation response |
| IGR-002 | Provider credential leakage | Critical | Provider key appears in DB dumps, logs, UI, errors, docs, or evidence | Use staging-only secrets, verify `ProviderAccount` encryption and redaction | Any provider credential is exposed |
| IGR-003 | Normal user reaches `experimental_proxy` | Critical | Normal user sees or calls experimental model | Verify visibility, route guards, fallback, retry, preferred, specific, and legacy paths | Normal user can see or call experimental path |
| IGR-004 | Zero-balance request calls upstream | High | Upstream receives request from zero-balance user | Test zero-balance before any real provider call; use fake provider first | Any zero-balance call reaches upstream |
| IGR-005 | Incorrect billing or balance deduction | High | Success/failure request changes quota incorrectly | Compare request result, token count, cost, and quota delta | Incorrect debit or credit is observed |
| IGR-006 | Prompt or response stored by default | High | `usage_log`, error log, or evidence contains full prompt/response | Keep full-text logging disabled; inspect logs after fixture calls | Full prompt/response is stored or published |
| IGR-007 | Admin permission bypass | Critical | Normal user accesses admin page/API or performs admin action | Verify admin middleware and UI/API access with normal user | Normal user can access admin capability |
| IGR-008 | Tenant ownership attribution error | High | Logs or API keys attach to wrong user/org/project | Test scoped API keys and inspect sanitized usage records | Wrong user/org/project attribution appears |
| IGR-009 | Disabled provider/channel still routes | High | Traffic reaches disabled provider or channel | Disable each object and retry fake-provider call | Disabled route is selected |
| IGR-010 | Real paid provider used unintentionally | High | Real upstream provider receives traffic without explicit approval | Default to fake provider; require manual low-limit staging key approval | Unapproved paid provider call occurs |
| IGR-011 | Migration or DB compatibility issue | High | Staging migration fails or corrupts ownership/billing/log data | Review migration evidence and isolate staging DB from production | Migration corrupts data or cannot be safely recovered |
| IGR-012 | Rollback is unclear or not executable | High | Operator cannot identify artifact, DB path, Redis cleanup, or owner | Document rollback owner, approval, artifact, and data path before execution | Rollback cannot be executed during gray test |
| IGR-013 | Frontend local script blocker hides UI regression | Medium | Local frontend script remains blocked while CI frontend-check is stale or failing | Require GitHub Actions `frontend-check` passed evidence; record local blocker as non-blocking only then | CI frontend-check fails or is missing |
| IGR-014 | Logs leak bearer token or sensitive `params.Other` data | High | Sensitive headers or payload fields appear in logs | Inspect redaction for success and error paths | Token or credential appears in logs |
| IGR-015 | OpenAI SDK compatibility regression | Medium | SDK fails against staging endpoint despite direct curl success | Include SDK smoke with staging key and fake/approved provider | SDK call fails and no workaround is accepted |

## Risk Review Cadence

- Review this register before Day 1 fake-provider execution.
- Update risk status after Day 3 admin, log, billing, and balance verification.
- Re-review before any Day 4-7 low-traffic observation.
- Do not proceed to limited beta or production preparation with open critical or high risks.
