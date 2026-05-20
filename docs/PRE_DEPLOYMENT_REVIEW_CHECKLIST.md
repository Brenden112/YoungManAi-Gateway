# Pre-Deployment Review Checklist

Date: 2026-05-20

Deployment readiness: `internal_gray_blocked`
Production readiness: `not_ready`

Current next step: rerun the Docker fixture regression with the hardened fixture scripts and record the result in `docs/INTERNAL_GRAY_TEST_REPORT.md`. This checklist does not mark production ready.

## Security Checks

- [x] API keys are not stored in plaintext. State: `passed_by_ci_and_code_inspection`.
- [x] ProviderAccount credentials are encrypted at rest. State: `passed_by_ci_and_code_inspection`.
- [x] `experimental_proxy` is not callable by normal users by default. State: `passed_by_ci_and_codespaces_fixture`.
- [x] Disabled `experimental_proxy` channels cannot route traffic. State: `passed_by_ci_and_codespaces_fixture`.
- [x] `allowed_provider_types` is enforced. State: `passed_by_ci_and_code_inspection`.
- [x] Organization/project token binding is enforced. State: `passed_by_ci_and_code_inspection`.
- [x] Prompt/response content is not written to logs by default. State: `passed_by_ci_and_codespaces_fixture`.
- [x] Zero balance rejects before upstream routing. State: `passed_by_ci_and_codespaces_fixture`.

## Runtime Checks

- [x] `go test ./...` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] `go vet ./...` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] `bash scripts/check-config-secrets.sh`
- [x] `LOCAL_FIXTURE=1 bash scripts/regression.sh` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] `docker compose config`
- [ ] Docker runtime smoke with `docker-compose.fixture.yml` - latest internal gray Codespaces run failed 3 regression checks with HTTP 503; see `docs/INTERNAL_GRAY_TEST_REPORT.md`.
- [x] Frontend `bun run lint` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Frontend `bun run test` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Frontend `bun run build` - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Cross-DB migration test for SQLite - passed in Codespaces via `bash scripts/ci-migration-check.sh`
- [x] Cross-DB migration test for MySQL - passed by Pre-release verification #16 on commit `73ad2ff`
- [x] Cross-DB migration test for PostgreSQL - passed by Pre-release verification #16 on commit `73ad2ff`

## Configuration Checks

- [x] `.env.example` contains no real secrets. State: `passed_by_config_secret_check`.
- [x] `.env.staging.example` exists and contains placeholders only.
- [x] Default `docker-compose.yml` reads database/cache/session/encryption secrets from environment variables.
- [x] Upstream provider keys are not committed to git. State: `passed_by_config_secret_check`.
- [ ] Encryption secret is configured before production use. State: `manual_required_before_production`.
- [ ] `provider_account_id` is configured correctly for linked channels. State: `manual_required_during_internal_gray`.
- [x] `experimental_proxy` uses `enabled=false` by default unless explicitly approved. State: `passed_by_ci_and_code_inspection`.
- [x] `allow_experimental=false` by default for tokens. State: `passed_by_ci_and_code_inspection`.

## Manual Release Gate

- [x] Frontend test waiver reviewed: `docs/WAIVERS/FRONTEND_TEST_INFRA_WAIVER.md` - remaining local script issue is non-blocking because CI frontend-check passed.
- [x] Cross-DB migration waiver reviewed: `docs/WAIVERS/CROSS_DB_MIGRATION_WAIVER.md` - Codespaces migration check passed.
- [x] Docker runtime smoke waiver reviewed: `docs/WAIVERS/DOCKER_RUNTIME_SMOKE_WAIVER.md` - Codespaces Docker fixture runtime passed.
- [x] Local Go toolchain waiver reviewed: `docs/WAIVERS/LOCAL_GO_TOOLCHAIN_WAIVER.md` - Codespaces Go verification passed.
- [x] CI `pre-release-verification` workflow artifacts reviewed for run #16 on commit `73ad2ff`.
- [x] Staging verification runbook completed in Codespaces; evidence recorded in `docs/CODESPACES_STAGING_EVIDENCE.md`.
- [x] Internal gray test plan prepared in `docs/INTERNAL_GRAY_TEST_PLAN.md`.
- [ ] Internal gray Docker fixture regression passed. State: `failed_pending_rerun_after_fixture_script_hardening`.
- [x] No real upstream API key, token, credential, prompt, or response appears in test fixtures or docs. State: `passed_by_config_secret_check`.
- [ ] Copy `.env.staging.example` to an untracked `.env.staging` and fill controlled staging secrets from the approved secret source. State: `manual_required_for_internal_gray`.
- [ ] Confirm `.env`, `.env.*`, `secrets/`, `*.key`, `*.pem`, `*.p12`, and `*.pfx` are not committed. State: `manual_required_for_internal_gray`.
- [x] Run staging runtime verification in an isolated environment - passed in GitHub Codespaces.
