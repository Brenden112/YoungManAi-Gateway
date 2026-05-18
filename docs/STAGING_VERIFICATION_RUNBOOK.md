# Staging Verification Runbook

Date: 2026-05-18

Status: `pending_ci_verification`

This runbook verifies release readiness in a stable CI or staging environment. Do not use real upstream provider keys, real customer tokens, real credentials, real prompts, or real responses. All provider traffic must use the fake provider or local fixture.

## Environment Preparation

- Go 1.22+ or the version declared by `go.mod`.
- Docker daemon with Compose v2.
- Node plus Bun, pnpm, npm, or yarn. The default frontend is Bun-first.
- SQLite available locally.
- MySQL and PostgreSQL staging services or disposable CI services.
- Fake OpenAI-compatible provider from `scripts/fake-openai-provider.mjs`.
- Fixture compose file: `docker-compose.fixture.yml`.
- No production `.env` mounted into the fixture.

## Required Commands

```bash
bash scripts/ci-verify.sh
docker compose config
LOCAL_FIXTURE=1 bash scripts/regression.sh
bash scripts/ci-migration-check.sh
docker compose -f docker-compose.fixture.yml config
docker compose -f docker-compose.fixture.yml up -d --build
ADMIN_LINE="$(BASE_URL=http://localhost:3000 bash scripts/seed-local-fixture.sh | tail -n 1)"
eval "$ADMIN_LINE"
BASE_URL=http://localhost:3000 ADMIN_TOKEN="$ADMIN_TOKEN" bash scripts/regression.sh
docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes
```

For MySQL and PostgreSQL migration checks, set `SQL_DSN` to the staging or CI service DSN before running `bash scripts/ci-migration-check.sh`.

## Fake Provider Verification

Verify these behaviors against the seeded fixture:

- `/v1/models` returns fixture-visible models.
- `/v1/chat/completions` succeeds through the fake OpenAI-compatible provider.
- Official provider success uses only fixture credentials.
- Normal user calls to `experimental_proxy` are rejected.
- Internal user with `allow_experimental=true` can call `experimental_proxy` when the channel is enabled.
- Disabled `experimental_proxy` channels are rejected.
- Zero-balance requests reject before any upstream call.
- `usage_log` does not contain full prompt/response content by default.

## Manual Checks

- `.env` and CI secrets contain no real upstream API keys for this fixture.
- ProviderAccount credentials are encrypted at rest.
- API keys are not stored in plaintext.
- `experimental_proxy` is disabled/internal-only by default.
- `allowed_provider_types` restrictions are enforced.
- Organization/project token binding is enforced.
- CI logs do not print real keys, tokens, credentials, prompts, or responses.

## Pass Criteria

- Pre-release verification CI is green.
- `critical_findings_remaining = 0`.
- `high_findings_remaining = 0`.
- All accepted blockers are closed by CI/staging evidence or explicitly signed off by a human release owner.
- `deployment_readiness` may move from `needs_manual_review` to `staging_ready` only after the above evidence is reviewed.

## Failure Handling

- If a tool or service is missing, mark the result as a CI/staging blocker; do not mark it passed.
- If fake-provider runtime assertions fail, open a confirmed code bug with endpoint, seed state, expected behavior, actual behavior, and logs.
- If Docker cleanup fails, only use fixture-scoped cleanup. Do not run `docker system prune -a`, `docker volume prune`, or delete unrelated images/volumes.
