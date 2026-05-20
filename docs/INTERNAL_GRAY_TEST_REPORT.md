# Internal Gray Test Report

Date: 2026-05-20

Deployment readiness: `internal_gray_blocked`
Production readiness: `not_ready`

## Summary

An internal gray Docker fixture runtime attempt was executed in GitHub Codespaces against the fake-provider fixture. The fixture stack built and started successfully, and the seed script emitted an admin token. The regression suite then failed 3 of 8 checks with HTTP 503 responses.

No real upstream provider or API key was used.

## Evidence

| Step | Result |
|---|---|
| `docker compose -f docker-compose.fixture.yml up -d --build` | Passed. Image built and `redis`, `fake-upstream`, and `new-api` containers started. |
| `ADMIN_LINE="$(BASE_URL=http://localhost:3000 bash scripts/seed-local-fixture.sh \| tail -n 1)"` | Passed. Seed emitted `ADMIN_TOKEN`, `ADMIN_USER_ID=1`, and `BASE_URL=http://localhost:3000`. |
| `BASE_URL=http://localhost:3000 ADMIN_TOKEN="$ADMIN_TOKEN" ADMIN_USER_ID="$ADMIN_USER_ID" bash scripts/regression.sh` | Failed. 5 passed, 3 failed, exit 1. |
| `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes` | Requested cleanup command was part of the run sequence. |

## Failed Checks

| Check | Expected | Actual |
|---|---|---|
| T1 normal user can chat through official provider | HTTP 200 | HTTP 503 |
| T3 normal user gets 403 on experimental model | HTTP 403 | HTTP 503 |
| T5 zero-quota user gets 402 | HTTP 402 | HTTP 503 |

## Passing Checks

| Check | Result |
|---|---|
| T1 normal user can list models | HTTP 200 |
| T2 experimental models absent for normal user | Passed |
| T4 disable-experimental returns success | Passed |
| T4 disabled experimental channel returns 503 | Passed |
| T6 prompt/response not stored in log | Passed |

## Immediate Remediation

The fixture scripts were hardened after this failure:

- `scripts/seed-local-fixture.sh` now calls `/api/channel/fix` after creating fixture channels. This rebuilds ability rows and refreshes runtime channel cache state before regression traffic starts.
- `scripts/regression.sh` now checks setup API responses, verifies created API tokens are non-empty, asserts the model list contains `gpt-4o-mini`, and prints response bodies for T1, T3, and T5 failures.

## Current Status

Internal gray execution is blocked until the Docker fixture regression is rerun with the hardened scripts and returns all checks passing. Production readiness remains `not_ready`.

## Next Recommended Action

Rerun the same Codespaces fixture command sequence. If any 503 remains, use the newly printed response body plus container logs to determine whether the failure is channel selection, upstream reachability, quota pre-consume, or provider-policy handling.
