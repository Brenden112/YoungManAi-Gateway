# Pre-Deployment Review Checklist

Date: 2026-05-20

Deployment readiness: `staging_ready_pending_runtime_signoff`
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

- [x] `go test ./...` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] `go vet ./...` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] `bash scripts/check-config-secrets.sh`
- [x] `LOCAL_FIXTURE=1 bash scripts/regression.sh` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] `docker compose config`
- [x] Docker runtime smoke with `docker-compose.fixture.yml` - passed by Pre-release verification #16 on commit `73ad2ff`; isolated staging runtime verification still required before production.
- [x] Frontend `bun run lint` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Frontend `bun run test` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Frontend `bun run build` - passed by Pre-release verification #16 on commit `73ad2ff`
- [ ] Cross-DB migration test for SQLite - blocked locally because `go` is not in PATH; passed in CI #16 migration job context
- [x] Cross-DB migration test for MySQL - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Cross-DB migration test for PostgreSQL - passed by Pre-release verification #16 on commit `73ad2ff`

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
- [x] CI `pre-release-verification` workflow artifacts reviewed for run #16 on commit `73ad2ff`.
- [ ] Staging verification runbook completed: `docs/STAGING_VERIFICATION_RUNBOOK.md`
- [ ] No real upstream API key, token, credential, prompt, or response appears in test fixtures or docs.
- [ ] Copy `.env.staging.example` to an untracked `.env.staging` and fill controlled staging secrets from the approved secret source.
- [ ] Confirm `.env`, `.env.*`, `secrets/`, `*.key`, `*.pem`, `*.p12`, and `*.pfx` are not committed.
- [ ] Run staging runtime verification in an isolated environment.
