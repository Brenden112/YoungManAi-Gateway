# Docker Runtime Smoke Waiver

Date: 2026-05-18

Status: `closed_by_ci`

blocker_id: `skipped_environment_docker_runtime`

## Current Verified Evidence

- Docker daemon was accessible during the fixture smoke attempt.
- `docker-compose.fixture.yml` was able to start the build path.
- `redis:7-alpine` and `node:24-alpine` image pulls/build dependencies were able to begin.
- No real upstream provider was called.
- No real API key was used.

## Not Completed

- Full Docker runtime startup.
- Seeded fixture initialization.
- Curl smoke runtime verification.

## Reason Not Completed

The local build reached `go mod download` / dependency download and then remained without useful progress for an extended period. To avoid blocking the audit flow on a local environment, network, or build-cache condition, this run stopped waiting and records the Docker runtime smoke as an accepted blocker.

## Risk Nature

- Environment / CI blocker.
- Not a currently confirmed business-code bug.
- Not evidence of provider routing, billing, credential, or log privacy failure.

## Required CI / Staging Commands

```bash
docker compose -f docker-compose.fixture.yml config
docker compose -f docker-compose.fixture.yml up -d --build
ADMIN_LINE="$(BASE_URL=http://localhost:3000 bash scripts/seed-local-fixture.sh | tail -n 1)"
eval "$ADMIN_LINE"
BASE_URL=http://localhost:3000 ADMIN_TOKEN="$ADMIN_TOKEN" bash scripts/regression.sh
docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes
```

## Manual Acceptance Steps

1. Run the required commands in an isolated CI runner or staging host.
2. Confirm the fixture uses only fake/local provider credentials.
3. Confirm no production `.env`, real provider key, real customer token, prompt, or response is mounted into the fixture.
4. Confirm `/v1/models` and `/v1/chat/completions` pass through the fake OpenAI-compatible provider.
5. Confirm normal users cannot use `experimental_proxy`.
6. Confirm disabled `experimental_proxy` channels cannot route.
7. Confirm zero-balance requests are rejected before upstream routing.
8. Confirm usage logs do not persist full prompt/response text by default.
9. Tear down only fixture resources with `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes`.

## Pass Criteria

- Fixture compose config validates.
- Fixture runtime starts successfully.
- Seed script creates the required local admin/user/token/channel records.
- `LOCAL_FIXTURE` or seeded runtime regression smoke exits `0`.
- Fake provider receives only synthetic fixture traffic.
- No real provider call is observed.
- No real API key is present in compose config, logs, seed output, or regression output.

## Failure Handling

- Keep `deployment_readiness` out of `production_ready` until staging manual verification is complete.
- Keep production release blocked on staging environment review even after this CI closure.
- If CI/staging fails before dependency download completes, classify as environment/build-cache failure and add retry/cache evidence.
- If CI/staging reaches runtime and a functional assertion fails, open a confirmed code bug with logs, failing endpoint, expected behavior, actual behavior, and fixture seed state.
- Do not run `docker system prune -a`.
- Do not delete unrelated images or volumes.
- Do not call real upstream providers.
- Do not use real API keys.

## Cleanup Refresh — 2026-05-18

- Fixture cleanup attempted with `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes`.
- Docker CLI failed before compose action with `Failed to initialize: protocol not available`.
- Cleanup must be verified manually in CI/staging or when Docker CLI recovers.
- No broad Docker cleanup was attempted.

## CI / Staging Verification Path

- `.github/workflows/pre-release-verification.yml` job `docker-fixture-smoke`
- `docs/STAGING_VERIFICATION_RUNBOOK.md`

## CI Closure — 2026-05-19

Pre-release verification #13 on branch `main` at commit `aeb43e5` passed in GitHub Actions. The `docker-fixture-smoke` job succeeded, closing local Docker runtime blocker `skipped_environment_docker_runtime` as `closed_by_ci`.

Current HEAD refresh 2026-05-20: Pre-release verification #16 on branch `main` at commit `73ad2ff` also passed in GitHub Actions. The `docker-fixture-smoke` job succeeded for the current reviewed HEAD.

The original local Docker/runtime blocker remains historical evidence of an environment limitation: the local build path stalled at dependency download and later Docker CLI access was unreliable. CI provided the required fixture build, seed, curl smoke, and cleanup evidence.

Production readiness is not granted by this waiver closure. Keep a production preflight requirement for staging manual verification using `docs/STAGING_VERIFICATION_RUNBOOK.md`, environment-variable review, real deployment topology review, and manual security sign-off.
