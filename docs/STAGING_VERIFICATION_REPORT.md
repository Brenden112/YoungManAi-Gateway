# Staging Verification Report

## Summary

| Field | Value |
|---|---|
| Verification time | 2026-05-19T18:18:46+08:00 |
| Environment | Local workspace `/mnt/d/Projects/new-api`, shell timezone Asia/Shanghai |
| Commit verified locally | `5a5f136f` |
| CI evidence reviewed | Pre-release verification #13, commit `aeb43e5`, success |
| Status | `completed_with_blockers` |
| Deployment readiness | `staging_config_hardened_with_blockers` |
| Production readiness | `not_ready` |

This staging verification used only static inspection, local fixture scripts, compose validation, and existing fake-provider CI evidence. No real upstream provider key was used and no paid provider was called.

## Commands Run

| Command | Result |
|---|---|
| `bash scripts/check-config-secrets.sh` | Exit 0. Default compose and env examples passed static secret-pattern checks. |
| `bash scripts/ci-verify.sh` | Exit 2. No security-check failures, but blocked because `go` is missing in PATH and `web/default/node_modules` is missing. `git diff --check` passed. |
| `LOCAL_FIXTURE=1 bash scripts/regression.sh` | Exit 2. Blocked: Go binary not found. |
| `docker compose config` | Exit 0. Compose syntax valid; default compose now renders only environment-driven values and local-only placeholders. |
| `docker compose -f docker-compose.fixture.yml config` | Exit 0. Fixture compose syntax valid and uses `fake-upstream`, Redis, SQLite tmpfs, and `STORE_FULL_TEXT_ENABLED=false`. |
| `docker compose -f docker-compose.fixture.yml up -d --build` | Exit 1. Blocked: Docker CLI failed with `Failed to initialize: protocol not available`. |
| `docker compose -f docker-compose.fixture.yml down --remove-orphans --volumes` | Exit 1. Blocked by same Docker CLI initialization error. |
| `command -v go` | Exit 1. Go is unavailable in the current shell. |
| `command -v bun` | Exit 1. Bun is unavailable in the current shell. |
| `command -v npm` | Exit 0. npm exists, but frontend dependencies are not installed. |

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

No key leakage was found. Staging verification can continue only with blockers noted below because the current local runtime cannot execute Go or Docker fixture smoke.

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

Runtime API calls against the fixture are blocked locally because Docker cannot initialize. CI #13 already passed `local-fixture-regression` and `docker-fixture-smoke`, covering `GET /v1/models`, `POST /v1/chat/completions`, official provider success through the fake upstream, normal/internal experimental access control, disabled experimental rejection, zero balance rejection, and no full prompt/response storage.

Local result: `api_staging_check_blocked` due to `docker_staging_runtime_blocked`.

## Admin Verification Results

Backend route and UI source inspection confirms admin surfaces exist for channel/provider management, API keys, usage logs, user quota, and manual top-up. CI #13 passed frontend lint/test/build. Current local frontend runtime is blocked because Bun and `web/default/node_modules` are unavailable, and Docker fixture runtime cannot start.

Local result: `frontend_staging_check_blocked`.

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
| `local_go_toolchain_blocked` | `go` is not in PATH, so local Go tests and `LOCAL_FIXTURE=1` regression cannot run. | Run in CI/staging host with Go 1.22+ or restore the local Go toolchain. |
| `frontend_staging_check_blocked` | Bun is unavailable and `web/default/node_modules` is missing. | Run `bun install --frozen-lockfile`, `bun run lint`, `bun run test`, and `bun run build` in staging/CI. |
| `docker_staging_runtime_blocked` | Docker CLI fails with `Failed to initialize: protocol not available`. | Run fixture smoke on a Docker-capable staging host. |
| `api_staging_check_blocked` | API smoke depends on the Docker fixture runtime. | Re-run seeded fake-provider curl smoke once Docker is available. |

## Recommendation

Do not mark production ready. The code and CI evidence support staging configuration hardening, but this local staging run completed with environment blockers. Recommended readiness is `staging_config_hardened_with_blockers`; use `.env.staging.example` only as a template, never commit real secrets, and run staging runtime verification in an isolated environment before the internal gray test.
