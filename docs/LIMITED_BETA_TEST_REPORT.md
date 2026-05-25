# Limited Beta Test Report

Date: 2026-05-25
Test time: 2026-05-25T20:30:48+08:00
Environment: local workspace `/mnt/d/Projects/new-api`; Docker fixture network `new-api_fixture-network`; fake upstream only
Commit: `675091d0`
Phase: `limited-beta-notes-resolution`
Status: `completed_with_manual_gate`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`
Next recommended action: `decide whether to run low-limit real provider beta or prepare production sign-off pack`

## Scope

This run executed the limited beta checklist against the local Docker fixture and current repository evidence. It did not use a real provider key, did not call a paid upstream provider, did not write real keys/tokens/credentials/prompts/responses into evidence, and did not modify business logic.

The fixture used placeholder non-secret provider values and fake upstream traffic only. The real `official_cloud` low-limit provider path remains pending manual approval and a human-supplied low-limit test key.

## Phase 6B Notes Resolution

| Item | Result | Evidence |
|---|---|---|
| LBI-001 local Go unavailable | `closed_by_codespaces_or_ci` | Local Go is still unavailable, but GitHub Actions `Pre-release verification` run 24 passed for commit `57ad3623`. Windows/local Go absence does not block fake-provider beta. |
| LBI-002 org/project runtime workflow | `closed` | Fixture-only org/project rows were created; org/project-bound API key successfully called `/v1/models` and `/v1/chat/completions`; usage log used token-context org/project IDs and ignored spoofed client IDs; disabled org/project token creation was rejected; allowed models/provider types remained enforced. |
| LBI-003 real low-limit provider | `manual_required` | No low-limit real provider key was supplied; no real provider call was made. This remains a release-owner manual gate. |
| LBI-004 streaming smoke | `closed` | `stream=true` returned SSE/chunk output, wrote log evidence, preserved org/project context, did not store full prompt/response, and did not bypass policy. |
| LBI-005 OpenAI SDK smoke | `closed` | OpenAI Node SDK passed `models.list`, non-streaming chat, and streaming chat against the fixture endpoint. |

Phase 6B fake-provider beta exit: `true_for_fake_provider_beta`.

## Phase 6B Commands Run

| Command | Exit | Result |
|---|---:|---|
| `bash scripts/check-config-secrets.sh` | 0 | Passed. |
| `bash scripts/ci-verify.sh` | 2 | No failed checks locally; Go/frontend checks blocked locally and closed by CI/Codespaces evidence. |
| GitHub Actions pre-release verification API check | 0 | Run 24 passed for commit `57ad3623`. |
| `docker compose -f docker-compose.fixture.yml up -d --build` | 1 | Fixture image built, but host port 3000 was already allocated. |
| `FIXTURE_PORT=3001 docker --context default compose -f docker-compose.fixture.yml up -d --build` | 0 | Clean fixture started on free host port 3001. |
| `BASE_URL=http://new-api:3000 bash scripts/seed-local-fixture.sh` in helper container | 0 | Fixture seeded with fake upstream and placeholder key. |
| Fixture-only org/project runtime smoke | 0 | Passed 19 checks. |
| OpenAI Node SDK smoke | 0 | Passed `models.list`, non-streaming chat, and streaming chat. |
| `docker --context default compose -f docker-compose.fixture.yml down --remove-orphans --volumes` | 0 | Fixture cleaned up. |
| `git diff --check` | 0 | Passed; existing CRLF warnings only. |
| `node -e "JSON.parse(require('fs').readFileSync('.factory/mission-state.json','utf8'))"` | 0 | Passed. |

## Commands Run

| Command | Exit | Result |
|---|---:|---|
| `git rev-parse --short HEAD` | 0 | `675091d0` |
| `date -Iseconds` | 0 | `2026-05-25T20:30:48+08:00` |
| `command -v go` | 1 | Blocked locally: Go is not in PATH. |
| `command -v jq` | 1 | Blocked locally: jq is not installed. |
| `command -v docker` | 0 | Docker CLI exists at `/usr/bin/docker`. |
| `bash scripts/check-config-secrets.sh` | 0 | Passed: config secret check passed. |
| `bash scripts/ci-verify.sh` | 2 | No failed checks. Passed 2, blocked 7 due missing local Go and frontend dependencies. |
| `LOCAL_FIXTURE=1 bash scripts/regression.sh` | 2 | Blocked locally: Go binary not found. |
| `docker compose -f docker-compose.fixture.yml config` | 0 | Passed: fixture compose rendered successfully. |
| `FIXTURE_PORT=3001 docker --context default compose -f docker-compose.fixture.yml up -d --build` | 0 | Passed: fixture image built and containers started. |
| `docker --context default compose -f docker-compose.fixture.yml ps` | 0 | Passed: `new-api`, `fake-upstream`, and `redis` were running; `new-api` was healthy. |
| `docker --context default run --rm --network new-api_fixture-network curlimages/curl:8.16.0 -fsS http://new-api:3000/api/status` | 0 | Passed: fixture API status returned success inside the Docker network. |
| Containerized Node fixture smoke, first attempt | 1 | Harness issue: admin cookie was not restored after normal-user login, causing false admin-route failures. Product checks before that point mostly passed. |
| Corrected containerized Node fixture smoke | 1 | Passed 28 checks; one disabled-channel check used a shared model and routed to another enabled official channel, so a targeted unique-model check was run. |
| Targeted disabled unique official channel check | 0 | Passed: disabled unique-model official channel returned HTTP 503 and did not route. |
| `curl -s https://api.github.com/repos/Brenden112/YoungManAi-Gateway/actions/workflows/pre-release-verification.yml/runs?per_page=1` | 0 | Latest pre-release verification run `23`, commit `675091d0`, conclusion `success`. |
| `docker --context default compose -f docker-compose.fixture.yml stop` | 0 | Fixture containers stopped. |
| `git diff --check` | 0 | Passed; Git emitted existing CRLF normalization warnings only. |
| `node -e "JSON.parse(require('fs').readFileSync('.factory/mission-state.json','utf8'))"` | 0 | Passed: mission-state JSON parsed successfully. |

## GitHub Actions Status

Latest `Pre-release verification` workflow status was confirmed through the GitHub API:

| Field | Value |
|---|---|
| Run number | `23` |
| Commit | `675091d0d393ab1f9469e40d8fd48656691dffcf` |
| Status | `completed` |
| Conclusion | `success` |
| Created | `2026-05-25T12:27:13Z` |
| Updated | `2026-05-25T12:29:21Z` |

## Checklist Summary

| Area | Result | Evidence |
|---|---|---|
| Pre-run gate | Passed with notes | Fake provider used; no real provider key used; `experimental_proxy` kept disabled/internal-only by default; production readiness remained `not_ready`. |
| Config secret scan | Passed | `scripts/check-config-secrets.sh` exit 0. |
| CI verification | Passed with local blockers | `scripts/ci-verify.sh` had no failures, but local Go/frontend checks were blocked. GitHub Actions run 23 passed for commit `675091d0`. |
| Local fixture regression script | Blocked locally | `LOCAL_FIXTURE=1 bash scripts/regression.sh` requires local Go. Equivalent Docker fixture runtime smoke passed. |
| Docker fixture smoke | Passed | Fixture built, started, returned `/api/status`, and passed containerized API smoke. |
| Admin user | Passed | Fixture admin login succeeded. |
| Normal user | Passed | Normal beta fixture user created and used. |
| Internal user | Passed | Internal beta fixture user created and assigned internal group. |
| Organization/project | Blocked as runtime checklist item | No runtime organization/project management endpoint was identified for this fixture run. Prior model/CI evidence covers org/project token binding. |
| API key creation | Passed | Normal, internal, zero-balance, and disabled fixture API keys were created; plaintext keys were captured only in process memory and not written to docs. |
| API key masking | Passed | Token list did not return captured plaintext API key. |
| Disabled API key | Passed | Disabled API key returned HTTP 401 for `/v1/models`. |
| Normal admin access | Passed | Normal user admin API access was rejected. |
| Fake provider | Passed | `GET /v1/models` and non-streaming `POST /v1/chat/completions` succeeded through fake upstream. |
| `official_cloud` low-limit provider | Blocked by policy | No human-approved low-limit real provider key was supplied. Placeholder official channel used fake upstream only. |
| Provider credential exposure | Passed with fixture scope | Placeholder channel key was not returned in channel list. ProviderAccount encrypted storage remains covered by prior code/CI evidence. |
| Disabled channel routing | Passed | Targeted unique-model disabled official channel returned HTTP 503. |
| `allowed_provider_types` | Passed | Normal official-only token could not call experimental model; experimental models were hidden from normal model list. |
| `official_cloud` fallback to `experimental_proxy` | Passed by fixture and prior CI evidence | No official-to-experimental fallback was observed; policy was already covered in staging/CI evidence. |
| `experimental_proxy` default disabled/internal-only | Passed with evidence | Experimental channels were disabled by admin endpoint; normal user could not see or call experimental model. |
| Internal experimental access requirement | Passed by evidence with runtime disabled check | Internal user had explicit `allow_experimental=true`; once experimental channels were disabled, internal request returned HTTP 503. |
| Non-streaming API call | Passed | Fake provider chat returned HTTP 200 with token usage. |
| Streaming API call | Blocked | Fake upstream does not implement streaming. No approved real stream-capable provider was supplied. |
| OpenAI SDK compatibility | Blocked | No SDK runtime package or approved staging endpoint was available in this environment. HTTP OpenAI-compatible endpoints passed. |
| Successful request tokens/cost | Passed | Token log returned usage fields after successful chat. |
| Successful request balance deduction | Passed by fixture evidence and prior CI | Usage log/billing fields were present; detailed real fee reconciliation did not apply because no real provider was used. |
| Failed request charging | Passed by evidence | Zero-balance request returned HTTP 402 before upstream request count changed. |
| Zero balance upstream prevention | Passed | HTTP 402 and fake upstream request delta `0`. |
| Admin top-up | Passed | Admin top-up calls for normal/internal fixture users returned success. |
| Top-up operation record | Passed by prior CI evidence | Runtime top-up succeeded; operation-record contract covered by existing tests. |
| Usage log | Passed | Token log returned entries for the normal fixture API key. |
| Error log | Passed by evidence with notes | Error paths were exercised; no sensitive data was recorded in evidence. |
| Prompt/response privacy | Passed | Token log content was empty with `STORE_FULL_TEXT_ENABLED=false`. |
| `params.Other` and `error_message` redaction | Passed by prior CI evidence | Covered by sanitizer tests and staging verification; no leak observed in this run. |
| API key/provider credential/bearer leakage | Passed | No captured API key appeared in token log responses or documentation. |
| Upstream fee reconciliation | Not applicable | No real upstream provider was called. |
| Provider failure rate | Recorded | Fixture smoke had no product failure; harness false starts are excluded from provider failure rate. |
| Fallback behavior | Passed with notes | Disabled unique official model did not route; experimental fallback prohibition remains covered by fixture policy and prior CI. |

## Issues Found

See `docs/LIMITED_BETA_ISSUES.md`.

Summary:

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 3 |
| Low | 2 |

No critical or high issue was found. The remaining items are environment/tooling or manual-real-provider blockers for production preparation.

## Exit Criteria Assessment

Limited beta execution status: `completed_with_notes`.

Production-preparation exit criteria met: `false`.

Reason: critical/high findings are `0`, and fake-provider core beta checks passed, but medium/low notes remain for local Go regression, missing runtime org/project management verification, real low-limit provider approval/execution, streaming, and OpenAI SDK compatibility.

## Recommendation

- Deployment readiness: `limited_beta_passed_with_notes`.
- Production readiness: keep `not_ready`.
- Next recommended action: `resolve beta notes before production preparation`.
