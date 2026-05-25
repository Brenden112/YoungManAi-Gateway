# Limited Beta Issues

Date: 2026-05-25

Current phase: `limited-beta-execution`
Deployment readiness: `limited_beta_passed_with_notes`
Production readiness: `not_ready`

## Summary

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 3 |
| Low | 2 |

No critical or high issue was found. Limited beta should not move to production preparation until the medium/low notes below are resolved or explicitly accepted by the release owner.

## Open Issues

| ID | Severity | Status | Issue | Evidence | Minimal fix path |
|---|---|---|---|---|---|
| LBI-001 | Medium | Open | Local Go toolchain is unavailable, so `LOCAL_FIXTURE=1 bash scripts/regression.sh` and full local `ci-verify` Go checks are blocked in this workspace. | `command -v go` exit 1; `LOCAL_FIXTURE=1 bash scripts/regression.sh` exit 2; `ci-verify` passed 2 and blocked 7 with no failures. | Run these checks in the CI/Codespaces environment or install Go 1.22+ locally. GitHub Actions pre-release verification run 23 passed for commit `675091d0`. |
| LBI-002 | Medium | Open | Runtime organization/project creation and binding was not executed through a dashboard/API workflow in this fixture run. | No organization/project management route was identified during fixture execution; prior model/CI evidence covers tenant binding. | Execute org/project creation through the intended admin workflow in beta environment, or add a documented operator runbook for the existing workflow before production preparation. |
| LBI-003 | Medium | Open | Real low-limit `official_cloud` provider call was not executed because no human-approved low-limit key was supplied. | Fixture used fake upstream and placeholder non-secret channel key only. | Release owner supplies and approves a low-limit test key, then run the approved real-provider subset with sanitized prompts and no secret capture. |
| LBI-004 | Low | Open | Streaming API compatibility was not executed. | Fake upstream does not implement streaming; no approved stream-capable real provider was supplied. | Run streaming smoke against a stream-capable fake provider or an approved low-limit provider before production preparation if streaming is in beta scope. |
| LBI-005 | Low | Open | OpenAI SDK runtime compatibility was not executed in this environment. | HTTP OpenAI-compatible `/v1/models` and `/v1/chat/completions` passed; no SDK runtime package/approved endpoint was available for this run. | Run a minimal SDK smoke against the fixture or approved beta endpoint and record sanitized status only. |

## Non-Issues / Harness Notes

| Note | Disposition |
|---|---|
| First containerized smoke attempt returned admin/log failures after switching normal-user cookie. | Harness issue only. Corrected smoke preserved admin cookie and passed the affected admin/log checks. |
| Corrected smoke initially checked disabled-channel routing with shared model `gpt-4o-mini`, which routed to another enabled official fixture channel. | Harness issue only. Targeted unique-model disabled-channel check passed with HTTP 503. |
| Host `curl http://localhost:3001/api/status` failed while Docker-network curl succeeded and container was healthy. | Local host networking quirk; fixture network API status passed. |

## Stop Criteria Review

| Stop condition | Observed |
|---|---|
| API key leakage | No |
| Provider credential leakage | No |
| Bearer token leakage | No |
| Normal user can call `experimental_proxy` | No |
| Normal user can see experimental-only models | No |
| Zero-balance request reaches upstream | No |
| Billing/balance critical error | No |
| Full prompt/response saved in usage or error evidence | No |
| Admin permission bypass | No |
| Organization/project authorization bypass | Not observed; runtime workflow not executed |
| Fallback bypasses provider type restrictions | No |
| Real high-privilege provider key configured or used | No |
| Unapproved real provider call | No |

No stop condition was triggered.
