# Pre-Deployment Review Checklist

Date: 2026-05-19

Deployment readiness: `staging_config_hardened_with_blockers`
Production readiness: `not_ready`

## Security Checks

- [ ] API keys are not stored in plaintext.
- [ ] ProviderAccount credentials are encrypted at rest.
- [ ] `experimental_proxy` is not callable by normal users by default.
- [ ] Disabled `experimental_proxy` channels cannot route traffic.
- [ ] `allowed_provider_types` is enforced.
- [ ] Organization/project token binding is enforced.
- [ ] Prompt/response content is not written to logs by default.
- [ ] Zero balance rejects before upstream routing.

## Runtime Checks

- [ ] `go test ./...` - `pending_ci_verification`
- [ ] `go vet ./...` - `pending_ci_verification`
- [x] `bash scripts/check-config-secrets.sh`
- [ ] `LOCAL_FIXTURE=1 bash scripts/regression.sh` - `pending_ci_verification`
- [ ] `docker compose config`
- [ ] Docker runtime smoke with `docker-compose.fixture.yml` - `pending_ci_verification`, requires CI/staging verification before production.
- [ ] Frontend `bun run lint` - `pending_ci_verification`
- [ ] Frontend `bun run test` - `pending_ci_verification`
- [ ] Frontend `bun run build` - `pending_ci_verification`
- [ ] Cross-DB migration test for SQLite
- [ ] Cross-DB migration test for MySQL - `pending_ci_verification`
- [ ] Cross-DB migration test for PostgreSQL - `pending_ci_verification`

## Configuration Checks

- [ ] `.env.example` contains no real secrets.
- [x] `.env.staging.example` exists and contains placeholders only.
- [x] Default `docker-compose.yml` reads database/cache/session/encryption secrets from environment variables.
- [ ] Upstream provider keys are not committed to git.
- [ ] Encryption secret is configured before production use.
- [ ] `provider_account_id` is configured correctly for linked channels.
- [ ] `experimental_proxy` uses `enabled=false` by default unless explicitly approved.
- [ ] `allow_experimental=false` by default for tokens.

## Manual Release Gate

- [ ] Frontend test waiver reviewed: `docs/WAIVERS/FRONTEND_TEST_INFRA_WAIVER.md`
- [ ] Cross-DB migration waiver reviewed: `docs/WAIVERS/CROSS_DB_MIGRATION_WAIVER.md`
- [ ] Docker runtime smoke waiver reviewed: `docs/WAIVERS/DOCKER_RUNTIME_SMOKE_WAIVER.md`
- [ ] Local Go toolchain waiver reviewed: `docs/WAIVERS/LOCAL_GO_TOOLCHAIN_WAIVER.md`
- [ ] CI `pre-release-verification` workflow artifacts reviewed.
- [ ] Staging verification runbook completed: `docs/STAGING_VERIFICATION_RUNBOOK.md`
- [ ] No real upstream API key, token, credential, prompt, or response appears in test fixtures or docs.
- [ ] Copy `.env.staging.example` to an untracked `.env.staging` and fill controlled staging secrets from the approved secret source.
- [ ] Confirm `.env`, `.env.*`, `secrets/`, `*.key`, `*.pem`, `*.p12`, and `*.pfx` are not committed.
- [ ] Run staging runtime verification in an isolated environment.
