# Limited Beta Issues

Date: 2026-05-25

Current phase: `limited-beta-notes-resolution`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`

## Summary

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 0 |
| Low | 0 |

No critical or high issue was found. Fake-provider limited beta notes are closed. Real-provider beta remains under a manual provider gate.

## Resolved Issues

| ID | Severity | Status | Issue | Evidence | Minimal fix path |
|---|---|---|---|---|---|
| LBI-001 | Medium | `closed_by_codespaces_or_ci` | Local Go remains unavailable in this workspace, but current CI/Codespaces evidence passed. | GitHub Actions `Pre-release verification` run 24 passed for commit `57ad3623`. Local `ci-verify` still reports Go/frontend blockers only. | Keep CI/Codespaces required for Go verification; Windows/local Go absence does not block fake-provider beta. |
| LBI-002 | Medium | `closed` | Runtime organization/project binding needed fixture execution. | Fixture-only org/project rows were created, org/project-bound API key called `/v1/models` and `/v1/chat/completions`, usage log recorded token-context `org_id=60001` and `project_id=60002`, spoofed client IDs were ignored, disabled org/project token creation was rejected, and model/provider-type limits were enforced. | None for fake-provider beta. |
| LBI-004 | Low | `closed` | Streaming API compatibility needed fixture execution. | `stream=true` chat returned SSE/chunk output, wrote a completion log, preserved org/project context, did not store full prompt/response, and did not bypass provider-type/model policy. | None for fake-provider beta. |
| LBI-005 | Low | `closed` | OpenAI SDK runtime compatibility needed fixture execution. | OpenAI Node SDK in a temporary container passed `models.list`, non-streaming chat, and streaming chat against the fixture endpoint. | None for fake-provider beta. |

## Manual Gate

| ID | Severity | Status | Issue | Evidence | Minimal fix path |
|---|---|---|---|---|---|
| LBI-003 | Manual gate | `manual_required` | Real low-limit `official_cloud` provider call was not executed because no human-approved low-limit key was supplied. | Fixture used fake upstream and placeholder non-secret channel key only. | Release owner decides whether to provide a low-limit test key and approve real-provider limited beta. Do not run real-provider calls without explicit approval. |

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
| Organization/project authorization bypass | No |
| Fallback bypasses provider type restrictions | No |
| Real high-privilege provider key configured or used | No |
| Unapproved real provider call | No |

No stop condition was triggered. Organization/project runtime binding was verified in Phase 6B.
