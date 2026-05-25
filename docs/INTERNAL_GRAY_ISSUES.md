# Internal Gray Issues

## 2026-05-25 Runtime Retry Closure

Runtime retry closure status: `passed`

| ID | Severity | Status | Finding | Evidence | Closure |
|---|---|---|---|---|---|
| IG-2026-05-25-001 | Medium | `fixed_by_fixture_dockerfile` | Docker fixture build could be killed in the production `builder-classic` frontend stage under Codespaces resource pressure. | `Dockerfile.fixture` removes `builder-classic` from fixture builds; fixture build passed. | Closed. Fixture-only path is local fixture / staging smoke only, not production. |
| IG-2026-05-25-002 | Low | `environment_note` | Exact `localhost:3000` replay was blocked by unrelated container `aiclient2api`. | Validation used `FIXTURE_PORT=3001` and `new-api_fixture-network`; runtime regression passed `9 passed, 0 failed`. | Non-product environment note. |
| IG-2026-05-22-runtime-blockers | Medium | `passed_in_codespaces` | Previously blocked fixture runtime checks needed Docker-capable execution. | `/api/status`, fixture seed, and runtime regression passed with fake upstream and placeholder fixture keys. | Closed for internal gray runtime retry. |

Closure readiness:

- Critical findings: `0`
- High findings: `0`
- Exit criteria met: `true`
- Deployment readiness: `limited_beta_ready`
- Production readiness: `not_ready`
- Next recommended action: `prepare limited beta release plan`

## 2026-05-25 Runtime Retry Update

Fixture build instability is closed for the fixture path. `Dockerfile.fixture` removes the production Dockerfile's separate classic frontend build from `docker-compose.fixture.yml`, and the fixture image built successfully before runtime validation.

New environment note: exact host-port validation on `localhost:3000` is blocked in this workspace because unrelated container `aiclient2api` already publishes port 3000. Equivalent fixture validation passed on `FIXTURE_PORT=3001` and inside `new-api_fixture-network`.

| ID | Severity | Status | Finding | Evidence | Recommended next step |
|---|---|---|---|---|---|
| IG-2026-05-25-001 | Medium | Closed | Codespaces-style fixture build pressure from the production `builder-classic` stage caused build instability. | `new-api-fixture:latest` built successfully through `Dockerfile.fixture`; regression passed 9/9 against fake upstream. | Keep `Dockerfile.fixture` scoped to local fixture / staging smoke; continue using production `Dockerfile` for release images. |
| IG-2026-05-25-002 | Low | Open | Exact local `localhost:3000` fixture replay is blocked by an unrelated running service. | Docker reports `aiclient2api` publishing `0.0.0.0:3000->3000/tcp`. | Free port 3000 before exact command replay, or set `FIXTURE_PORT` for side-by-side local validation. |

Date: 2026-05-22
Phase: `internal-gray-test-execution`
Environment: local workspace `/mnt/d/Projects/new-api`
Test commit: `c961125c`

## Summary

No critical or high product issue was found. The Phase 4 runtime retry closure is `passed`; prior fixture runtime blockers are now `passed_in_codespaces`.

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 0 |
| Low | 1 |

## Open Issues

| ID | Severity | Status | Finding | Evidence | Recommended next step |
|---|---|---|---|---|---|
| IG-2026-05-22-001 | Medium | `passed_in_codespaces` | Local Go toolchain was unavailable, blocking `ci-verify.sh` Go checks and `LOCAL_FIXTURE=1 bash scripts/regression.sh`. | Runtime retry closure passed; regression passed `9 passed, 0 failed` in Docker fixture path. | Closed for internal gray runtime retry. |
| IG-2026-05-22-002 | Medium | `passed_in_codespaces` | Local `jq` was unavailable, blocking fixture seed. | Fixture seed succeeded in Docker fixture path with placeholder fixture keys. | Closed for internal gray runtime retry. |
| IG-2026-05-22-003 | Medium | `passed_in_codespaces` | Docker daemon operations were unavailable in the prior local context. | Fixture build passed through `Dockerfile.fixture`; runtime checks passed in Docker fixture path. | Closed for internal gray runtime retry. |
| IG-2026-05-22-004 | Medium | `passed_in_codespaces` | Fresh curl smoke and runtime admin/log checks could not execute because fixture startup was blocked. | `/api/status` succeeded in fixture network; runtime regression passed `9 passed, 0 failed`. | Closed for internal gray runtime retry. |
| IG-2026-05-22-005 | Low | Open | Frontend local script checks remain a non-blocking local environment note. | `ci-verify.sh` recorded frontend dependencies missing while GitHub Actions `frontend-check` evidence is passed. | Keep CI frontend-check evidence current; install frontend dependencies only if local frontend validation is required. |

## Stop Criteria Review

None of the documented stop criteria were observed:

- No API key leakage observed.
- No provider credential leakage observed.
- No bearer token leakage observed.
- No real prompt or response was stored in this report.
- No normal-user `experimental_proxy` access was observed.
- No zero-balance upstream call was observed.
- No billing corruption was observed.
- No admin permission bypass was observed.

The absence of observed stop criteria is not a full pass for runtime behavior because the local runtime fixture was blocked.
