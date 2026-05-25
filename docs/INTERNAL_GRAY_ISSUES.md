# Internal Gray Issues

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

No critical or high product issue was found in this Phase 4 local execution. The run is `completed_with_notes` because local environment blockers prevented a complete fresh internal gray runtime execution.

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 4 |
| Low | 1 |

## Open Issues

| ID | Severity | Status | Finding | Evidence | Recommended next step |
|---|---|---|---|---|---|
| IG-2026-05-22-001 | Medium | Open | Local Go toolchain is unavailable, blocking `ci-verify.sh` Go checks and `LOCAL_FIXTURE=1 bash scripts/regression.sh`. | `command -v go` exited 1; regression exited 2 with `Go binary not found`. | Rerun in Codespaces or staging host with Go available. |
| IG-2026-05-22-002 | Medium | Open | Local `jq` is unavailable, blocking fixture seed. | `command -v jq` exited 1; `scripts/seed-local-fixture.sh` exited 1 with `missing required command: jq`. | Install `jq` in the internal gray executor or run in Codespaces. |
| IG-2026-05-22-003 | Medium | Open | Docker daemon operations are unavailable in this local environment. | `docker ps`, fixture `up`, and fixture `down` returned `Failed to initialize: protocol not available`. | Rerun Docker fixture runtime on Docker-capable staging host or Codespaces. |
| IG-2026-05-22-004 | Medium | Open | Fresh curl smoke and runtime admin/log checks could not execute because fixture startup was blocked. | `curl http://localhost:3000/api/status` exited 7 because no server was running. | Rerun after fixture startup succeeds. |
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
