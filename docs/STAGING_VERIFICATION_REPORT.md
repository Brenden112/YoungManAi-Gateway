# Staging Verification Report

## Summary

| Field | Value |
|---|---|
| Verification time | 2026-05-22 |
| Environment | GitHub Codespaces isolated staging runtime; local workspace `/mnt/d/Projects/new-api`, shell timezone Asia/Shanghai |
| Commit verified locally | `b9179a38` |
| CI evidence reviewed | Pre-release verification #16, commit `73ad2ff`, success |
| Status | `passed` |
| Phase 2 status | `passed` |
| Isolated staging runtime verification status | `passed` |
| Codespaces staging status | `passed` |
| Docker fixture runtime status | `passed_in_codespaces` |
| Deployment readiness | `internal_gray_ready` |
| Production readiness | `not_ready` |

This staging verification used only static inspection, local fixture scripts, compose validation, and existing fake-provider CI evidence. No real upstream provider key was used and no paid provider was called.

Phase 2 closure update: isolated staging runtime verification passed in GitHub Codespaces. See `docs/CODESPACES_STAGING_EVIDENCE.md` for the Codespaces evidence record. This report does not mark production ready. The next recommended action is to prepare the internal gray test plan.

## Commands Run

| Command | Result |
|---|---|
| `bash scripts/check-config-secrets.sh` | Exit 0. Default compose and env examples passed static secret-pattern checks. |
| `bash scripts/ci-verify.sh` | Exit 2. No security-check failures, but blocked because `go` is missing in PATH and `web/default/node_modules` is missing. `git diff --check` passed. |
| `LOCAL_FIXTURE=1 bash scripts/regression.sh` | Exit 2. Blocked: Go binary not found. |
| `bash scripts/ci-migration-check.sh` | Exit 2. Blocked: Go binary not found. |
| `docker compose config` | Exit 0. Compose syntax valid; default compose now renders only environment-driven values and local-only placeholders. |
| `docker compose -f docker-compose.fixture.yml config` | Exit 0. Fixture compose syntax valid and uses `fake-upstream`, Redis, SQLite tmpfs, and `STORE_FULL_TEXT_ENABLED=false`. |
| `docker compose -f docker-compose.fixture.yml up -d --build` | Exit 1. Blocked: Docker CLI failed with `Failed to initialize: protocol not available`. |
| `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes` | Exit 1. Blocked by same Docker CLI initialization error. |
| `docker ps` | Exit 1. Blocked: Docker CLI failed with `Failed to initialize: protocol not available`. |
| `command -v go` | Exit 1. Go is unavailable in the current shell. |
| `command -v bun` | Exit 1. Bun is unavailable in the current shell. |
| `command -v jq` | Exit 1. jq is unavailable in the current shell. |
| `command -v npm` | Exit 0. npm exists, but frontend dependencies are not installed. |

## Codespaces Verification Results

| Command / Check | Result |
|---|---|
| `bash scripts/check-config-secrets.sh` | Passed: `Config secret check passed`. |
| `bash scripts/ci-verify.sh` | Core verification passed: Go available, Go tests passed, Go vet passed, local fixture regression passed, git diff check passed. Summary reported `passed=7 failed=0 blocked=5`. |
| Go toolchain | `go1.26.1 linux/amd64`. |
| `go test ./model/... -count=1` | Passed. |
| `go test ./middleware/... -count=1` | Passed. |
| `go test ./... -count=1` | Passed. |
| `go vet ./...` | Passed. |
| `LOCAL_FIXTURE=1 bash scripts/regression.sh` | Passed. |
| `bash scripts/ci-migration-check.sh` | Passed. |
| `docker compose config` | Passed. |
| `docker compose -f docker-compose.fixture.yml config` | Passed. |
| `.factory/mission-state.json` JSON parse | Passed. |
| `docker compose -f docker-compose.fixture.yml up -d --build` | Passed; `fake-upstream`, `redis`, and `new-api` containers started. |
| Docker fixture runtime smoke | Passed with fake-provider fixture only. |
| `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes` | Cleanup passed. |

Frontend local script checks remained blocked in `ci-verify.sh` because root package scripts were missing and `web/default` dependencies were not installed in that local script context. This is recorded as a non-blocking note because GitHub Actions Pre-release verification #16 already passed `frontend-check`.

## Environment And Configuration Checks

| Check | Result |
|---|---|
| `docker-compose.yml` contains no reusable local credentials | Passed after staging-hardening. Sensitive values now come from environment variables with local-only placeholders for syntax validation. |
| `.env.example` contains no real secrets | Passed. Values are comments/placeholders such as `random_string`, `user:password`, and `your-password`. |
| `.env.staging.example` contains no real secrets | Passed. Values are staging placeholders and must be replaced in an untracked `.env.staging`. |
| Local `.env` or config leakage | Passed. `.gitignore` ignores `.env`, `.env.*`, `secrets/`, and key/certificate bundles while allowing `.env.example` and `.env.staging.example`. |
| Secret scan | Passed with fixture/test-only notes. Matches were placeholder fixture secrets, documentation examples, or synthetic test keys; no real-looking committed credential was identified. |
| Encryption secret documentation | Passed with staging caveat. `CRYPTO_SECRET` is used by fixture compose and migration checks; production must set a non-fixture value before real deployment. |
| ProviderAccount credential encryption | Passed by code inspection and CI evidence. `ProviderAccount.Key` is `gorm:"-" json:"-"`; `EncryptedKey` is stored and decrypted after find. |
| `experimental_proxy` default policy | Passed by code inspection and CI evidence. New experimental channels default high-risk/internal/manual-enable and manually disabled unless explicitly enabled. |
| `allow_experimental` default | Passed. Token field defaults false and request policy requires internal user plus token opt-in. |
| Fake/local fixture configuration | Passed. `scripts/fake-openai-provider.mjs`, `scripts/seed-local-fixture.sh`, and `docker-compose.fixture.yml` are present. |
| Fixture does not depend on real provider | Passed. Fixture compose uses `fake-upstream` and local Redis/SQLite only. |

No key leakage was found. Codespaces runtime verification closed the Go, Docker fixture, and API staging blockers.

## Security Verification Results

| Invariant | Result |
|---|---|
| API Key not stored in plaintext | Passed by code inspection and CI evidence: `tokens.key_hash` is used; legacy `key` is overwritten with non-secret hash. |
| API Key not looked up in plaintext | Passed: `GetTokenByKey` computes `HashTokenKey` and queries `key_hash`. |
| Full API Key displayed only once | Passed by code inspection: list/search responses use `key_prefix`; `GetFullKey()` returns empty. |
| Disabled API Key cannot call | Passed by code inspection and CI evidence: `ValidateUserToken` rejects non-enabled status. |
| ProviderAccount credential encrypted at rest | Passed by code inspection and tests referenced in CI evidence. |
| ProviderAccount credential decrypts only at runtime | Passed: `ResolveActiveCredential` uses linked account after GORM `AfterFind`; no fallback to channel key when linked. |
| `params.Other` and error logs redacted | Passed: `SanitizeLogOther`, `SanitizeLogString`, and `SanitizeErrorMessage` are used before persistence. |
| Full prompt not stored by default | Passed: fixture sets `STORE_FULL_TEXT_ENABLED=false`; consume log content is blanked unless enabled. |
| Full response not stored by default | Passed by same log path and sanitizer. |
| Normal user cannot see `experimental_proxy` | Passed by CI evidence and model-list filtering tests. Local runtime check blocked. |
| Normal user cannot call `experimental_proxy` | Passed by CI evidence and distributor guard. Local runtime check blocked. |
| Disabled `experimental_proxy` cannot be called | Passed by CI evidence and channel status guards. Local runtime check blocked. |
| official_cloud does not fallback to `experimental_proxy` | Passed by provider policy enforcement and CI evidence. |
| `allowed_provider_types` enforced across routing paths | Passed by code inspection and CI evidence: DB/cache/fallback/preferred/specific/final guards use provider policy. |
| org/project token binding enforced | Passed by code inspection and CI evidence: `ResolveTokenTenantScope` validates org/project status and membership. |
| Zero balance does not call upstream | Passed by code inspection and CI evidence: wallet quota checks return 402 with skip-retry before upstream. Local runtime check blocked. |

## API Verification Results

Runtime API calls against the fixture passed in Codespaces using the fake-provider fixture. This closes `api_staging_check_blocked` as `closed_in_codespaces`.

## Admin Verification Results

Backend route and UI source inspection confirms admin surfaces exist for channel/provider management, API keys, usage logs, user quota, and manual top-up. CI #16 passed frontend lint/test/build. Frontend local script checks in Codespaces remain a non-blocking note because the GitHub Actions `frontend-check` job already passed.

## Logs And Privacy

Passed by code inspection and CI evidence. Consume logs record user, token, provider, channel, model, token counts, quota/cost fields, log status via log type, org/project context, and sanitized `Other`. Default full-text logging is disabled in the fixture and `RecordConsumeLog` blanks content unless `StoreFullTextEnabled` is true.

## Provider And Experimental Proxy

Passed by code inspection and CI evidence. `experimental_proxy` requires internal user context and token `allow_experimental=true`; disabled channels are rejected; `allowed_provider_types` applies to model listing, route selection, fallback, preferred/specific channel paths, and selected-channel setup.

## Balance And Billing

Passed by code inspection and CI evidence. Wallet quota checks reject zero or insufficient balance with HTTP 402 and skip retry/no error-log flags. Successful usage logging records quota and token counts; settlement updates user, token, and channel quota paths.

## Issues Found

- No new critical or high security issue found.
- Fixed/mitigated: default `docker-compose.yml` no longer contains local example credentials such as `123456`; sensitive values are read from env vars with explicit local-only placeholders.
- Production remains blocked until staging/prod secrets are supplied from a controlled secret manager or an untracked `.env.staging` derived from `.env.staging.example`.

## Blocked Items

| Blocker | Reason | Minimal fix path |
|---|---|---|
| `local_go_toolchain_blocked` | Closed in Codespaces. | `closed_in_codespaces` |
| `docker_staging_runtime_blocked` | Closed in Codespaces. | `closed_in_codespaces` |
| `api_staging_check_blocked` | Closed in Codespaces. | `closed_in_codespaces` |
| `frontend_local_script_check` | Local script checks still blocked by missing root package scripts / web dependencies, but GitHub Actions frontend-check passed. | `non_blocking_note_ci_frontend_check_passed` |

## Recommendation

Do not mark production ready. Codespaces evidence supports `internal_gray_ready`; use `.env.staging.example` only as a template, never commit real secrets, and prepare the internal gray test plan. Production readiness remains `not_ready`.
