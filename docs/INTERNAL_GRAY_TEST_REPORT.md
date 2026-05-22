# Internal Gray Test Report

Date: 2026-05-22
Test time: 2026-05-22T22:49:24+08:00
Environment: local workspace `/mnt/d/Projects/new-api`, shell timezone Asia/Shanghai
Test commit: `c961125c`
Phase: `internal-gray-test-execution`
Status: `completed_with_notes`
Deployment readiness: `internal_gray_passed_with_notes`
Production readiness: `not_ready`

## Scope

This run executed the Phase 4 internal gray checklist as far as the local environment allowed. It did not use real upstream provider keys, did not call paid providers, did not write real keys/tokens/credentials/prompts/responses into evidence, and did not modify business logic.

Phase 2 GitHub Codespaces evidence remains the runtime evidence for Go, fixture, Docker, migration, and fake-provider paths. This local Phase 4 run found no critical or high product issue, but several local environment blockers prevented a complete fresh runtime execution here.

## Commands Executed

| Command | Exit | Result |
|---|---:|---|
| `git rev-parse --short HEAD` | 0 | `c961125c` |
| `date -Iseconds` | 0 | `2026-05-22T22:49:24+08:00` |
| `bash scripts/check-config-secrets.sh` | 0 | Passed: config secret check passed. |
| `command -v go` | 1 | Blocked: Go binary is not in PATH. |
| `command -v docker` | 0 | Docker CLI exists at `/usr/bin/docker`. |
| `command -v jq` | 1 | Blocked: `jq` is not installed. |
| `bash scripts/ci-verify.sh` | 2 | Completed with no failed checks, but blocked Go and frontend dependency checks. Summary: `passed=2 failed=0 blocked=7`. |
| `LOCAL_FIXTURE=1 bash scripts/regression.sh` | 2 | Blocked: Go binary not found. |
| `docker compose -f docker-compose.fixture.yml config` | 0 | Passed: fixture compose rendered successfully. |
| `docker ps` | 1 | Blocked: Docker daemon operation failed with `Failed to initialize: protocol not available`. |
| `docker compose -f docker-compose.fixture.yml up -d --build` | 1 | Blocked: Docker daemon operation failed with `Failed to initialize: protocol not available`. |
| `curl -fsS http://localhost:3000/api/status` | 7 | Blocked: fixture server was not running. |
| `BASE_URL=http://localhost:3000 bash scripts/seed-local-fixture.sh` | 1 | Blocked: missing required command `jq`. |
| `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes` | 1 | Blocked: Docker daemon operation failed with `Failed to initialize: protocol not available`. |

## Passed Items

- Config secret scan passed.
- Fixture compose config rendered successfully.
- No real upstream provider key was used.
- No real paid provider was called.
- No real API key, token, provider credential, prompt, or response was written into this report.
- Existing Phase 2 Codespaces evidence shows `scripts/ci-verify.sh`, Go tests, Go vet, local fixture regression, Docker fixture runtime, fixture cleanup, and migration check passed with fake-provider traffic only.
- Existing GitHub Actions pre-release verification #16 remains the frontend-check evidence for this review.

## Failed Items

No product behavior check failed in this local run. The run was incomplete because required local tools and Docker daemon access were unavailable.

## Blocked Items

| Blocked item | Reason | Severity | Minimal next step |
|---|---|---|---|
| Go test/vet and `LOCAL_FIXTURE=1` regression | `go` binary not found in PATH. | Medium | Rerun in Codespaces or a staging host with Go available. |
| Fixture seed | `jq` is missing. | Medium | Install `jq` in the staging executor or run in Codespaces image that includes it. |
| Docker fixture runtime | Docker daemon operations fail with `Failed to initialize: protocol not available`. | Medium | Rerun on a Docker-capable staging host or GitHub Codespaces. |
| Curl smoke | Fixture server was not running because Docker startup was blocked. | Medium | Rerun after Docker fixture starts successfully. |
| Fresh Phase 4 runtime checklist completion | Dependent on Go, `jq`, and Docker fixture runtime. | Medium | Rerun full checklist in Docker-capable internal gray environment. |
| Frontend local script checks | Existing non-blocking note: local frontend dependencies are unavailable; GitHub Actions `frontend-check` passed. | Low | Keep CI frontend-check as current evidence unless it fails. |

## User And API Key

| Check | Result | Evidence |
|---|---|---|
| Create admin user | Blocked locally | Fixture seed could not run because `jq` is missing and Docker fixture was not started. |
| Create normal user | Passed by existing fixture evidence; blocked locally | `scripts/regression.sh` covers normal user creation when fixture runs; local execution blocked by Docker/Go. |
| Create internal user | Passed by existing fixture evidence; blocked locally | `scripts/regression.sh` covers internal user creation and group assignment when fixture runs. |
| Create API key | Passed by existing fixture evidence; blocked locally | `scripts/regression.sh` creates normal/internal/zero-balance keys. |
| Full plaintext shown once | Passed by code/CI evidence; blocked locally | Token responses use one-time `data.key`; list/search paths use prefix/masked data. |
| Hash/prefix storage | Passed by code/CI evidence | `tokens.key_hash` and `key_prefix` evidence from staging report. |
| Disabled API key cannot call | Passed by code/CI evidence; blocked locally | Covered by auth validation evidence; fresh local runtime was blocked. |
| User/org/project binding | Passed by code/CI evidence; blocked locally | Covered by `ResolveTokenTenantScope` evidence and prior tests. |
| `allowed_models` | Passed by code/CI evidence; blocked locally | Covered by prior token model/provider limit tests. |
| `allowed_provider_types` | Passed by code/CI evidence; blocked locally | Covered by distributor guard evidence. |

## OpenAI-Compatible API

| Check | Result | Evidence |
|---|---|---|
| `GET /v1/models` | Passed by Phase 2 Codespaces fixture; blocked locally | Local fixture server unavailable. |
| `POST /v1/chat/completions` | Passed by Phase 2 Codespaces fixture; blocked locally | Local fixture server unavailable. |
| Non-streaming call | Passed by Phase 2 Codespaces fixture; blocked locally | Regression fixture covers non-streaming chat. |
| Streaming if supported | Blocked | No fresh runtime provider execution in this local environment. |
| OpenAI SDK compatibility | Blocked | Requires running fixture or approved staging endpoint. |

## Provider And Channel

| Check | Result | Evidence |
|---|---|---|
| `official_cloud` test provider callable | Passed by Phase 2 Codespaces fixture; blocked locally | Fake upstream official channel path passed in Codespaces. |
| ProviderAccount credential encrypted | Passed by code/CI evidence | `ProviderAccount.Key` is not persisted; encrypted key path is covered by tests. |
| ProviderAccount runtime decrypt | Passed by code/CI evidence | Runtime credential resolver evidence from staging report. |
| Disabled provider cannot route | Passed by code/CI evidence; blocked locally | Fresh local runtime disabled-provider check was not possible. |
| Disabled channel cannot route | Passed by code/CI evidence; blocked locally | Fresh local runtime disabled-channel check was not possible. |
| Legacy channel compatibility | Passed by code/CI evidence | Legacy channel compatibility tests and no-relay-regression evidence. |

## Experimental Proxy

| Check | Result | Evidence |
|---|---|---|
| Default disabled | Passed by code/CI evidence | Channel defaults and disable-experimental fixture setup. |
| Default internal-only | Passed by code/CI evidence | Provider policy and token opt-in guard evidence. |
| Normal user not visible | Passed by Phase 2/CI evidence; blocked locally | Model list filtering covered by tests and fixture evidence. |
| Normal user not callable | Passed by Phase 2/CI evidence; blocked locally | Regression expects 403 for normal experimental call. |
| Internal user requires `allow_experimental=true` | Passed by code/CI evidence; blocked locally | Token opt-in guard evidence. |
| Disabled experimental not callable | Passed by Phase 2/CI evidence; blocked locally | Regression disables experimental channels and checks no route. |
| No official fallback to experimental | Passed by code/CI evidence | Provider policy guard applies to fallback paths. |
| Fallback/retry/preferred/legacy cannot bypass | Passed by code/CI evidence | Staging report records enforcement across route paths. |

## Billing And Balance

| Check | Result | Evidence |
|---|---|---|
| Successful request records input/output/total tokens | Passed by code/CI evidence; blocked locally | Fresh usage log inspection was blocked. |
| Successful request calculates cost | Passed by code/CI evidence; blocked locally | Settlement path covered by prior tests and CI. |
| Successful request deducts balance | Passed by code/CI evidence; blocked locally | Settlement path covered by prior tests and CI. |
| Failed request does not incorrectly deduct | Passed by code/CI evidence; blocked locally | Prior regression covers insufficient-balance behavior. |
| Zero balance does not call upstream | Passed by Phase 2/CI evidence; blocked locally | Regression expects HTTP 402 before upstream. |
| Admin manual top-up | Passed by code/CI evidence; blocked locally | Admin top-up tests passed in CI evidence. |
| Top-up operation record | Passed by code/CI evidence; blocked locally | Admin top-up contract evidence. |

## Logs And Privacy

| Check | Result | Evidence |
|---|---|---|
| Successful request writes `usage_log` | Passed by code/CI evidence; blocked locally | Fresh local log query blocked. |
| Failed request writes log/error record | Passed by code/CI evidence; blocked locally | Existing log/error evidence. |
| Log includes user/org/project/API key IDs | Passed by code/CI evidence | Tenant attribution tests and staging report. |
| Log includes provider/channel/model/tokens/cost/status | Passed by code/CI evidence | Staging report and billing evidence. |
| Prompt not stored by default | Passed by Phase 2/CI evidence; blocked locally | Fixture sets `STORE_FULL_TEXT_ENABLED=false`; regression checks empty content. |
| Response not stored by default | Passed by code/CI evidence | Same default full-text logging path. |
| `params.Other` redacted | Passed by code/CI evidence | Sanitizer tests/evidence. |
| `error_message` redacted | Passed by code/CI evidence | Error sanitizer evidence. |
| No key/credential/bearer leakage | Passed by config scan and code/CI evidence | Config scan passed; no new leakage found in this run. |

## Admin Dashboard

| Check | Result | Evidence |
|---|---|---|
| Admin can view Provider | Passed by source/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Admin can view Channel | Passed by source/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Admin can view API Key | Passed by source/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Admin can disable API Key | Passed by source/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Admin can view `usage_log` | Passed by source/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Admin can view user balance | Passed by source/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Admin can top up | Passed by code/CI evidence; blocked locally | Runtime UI/API check blocked. |
| Normal user cannot access admin | Passed by auth evidence; blocked locally | Runtime UI/API check blocked. |

## Issue Counts

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 4 |
| Low | 1 |

## Issues Found

See `docs/INTERNAL_GRAY_ISSUES.md`.

## Exit Criteria Assessment

Exit criteria met: `false` for this local execution alone.

Reason: no critical/high issues were found, but fresh local execution of the Docker fixture, seed, curl smoke, Go-backed regression, streaming/API SDK checks, and several runtime admin/log checks was blocked by missing tools and Docker daemon failure. Existing Phase 2 Codespaces evidence supports the core fake-provider runtime path, but this Phase 4 run should be re-executed on a Docker-capable internal gray environment before a limited beta decision.

## Recommendations

- Limited beta recommendation: `not_recommended_from_this_local_execution_alone`.
- Production preparation recommendation: `not_recommended`.
- Production readiness: keep `not_ready`.
- Next recommended action: rerun blocked internal gray runtime checks in a Docker-capable environment with Go and `jq`, then update this report and sign-off.
