# Cross-DB Migration Waiver

Date: 2026-05-17

Status: `closed_by_ci`

## Blocked Reason

SQLite migration schema smoke was added and executed locally through `bash scripts/check-migrations.sh`. The current local environment does not provide disposable MySQL or PostgreSQL services, so live migration execution against those two database engines remains CI-gated.

## Local Evidence

The migration schema smoke checks these required tables and fields:

- `channels.provider_type`
- `provider_accounts`
- `channel_model_mappings`
- `organizations`
- `organization_members`
- `projects`
- `tokens.org_id`, `tokens.project_id`, `tokens.allowed_provider_types`
- `logs` usage-log extension fields: `org_id`, `project_id`, `is_experimental_proxy`, `provider_type`, `request_id`, `upstream_request_id`

## Minimum Fix Path

1. Start disposable MySQL 8.0 and PostgreSQL 15 services.
2. Set `SQL_DSN` to the target service.
3. Run `bash scripts/check-migrations.sh`.
4. Repeat for SQLite, MySQL, and PostgreSQL.

## CI Commands

```bash
bash scripts/check-migrations.sh
SQL_DSN='root:rootpass@tcp(127.0.0.1:3306)/new_api_migration?charset=utf8mb4&parseTime=true' bash scripts/check-migrations.sh
SQL_DSN='postgresql://postgres:postgres@127.0.0.1:5432/new_api_migration?sslmode=disable' bash scripts/check-migrations.sh
```

## CI Coverage

`.github/workflows/pre-release-hardening.yml` contains separate `migration-sqlite`, `migration-mysql`, and `migration-postgres` jobs with disposable database services.

## Acceptance

This waiver was acceptable only for local environments without MySQL/PostgreSQL services. Deployment readiness must remain below `production_ready` until staging manual verification is complete.

## CI / Staging Verification Path

- `scripts/ci-migration-check.sh`
- `.github/workflows/pre-release-verification.yml` job `cross-db-migration`
- `docs/STAGING_VERIFICATION_RUNBOOK.md`

## CI Closure — 2026-05-19

Pre-release verification #13 on branch `main` at commit `aeb43e5` passed in GitHub Actions. The `cross-db-migration` job succeeded, closing local external dependency blocker `blocked_external_dependency_cross_db_runtime` as `closed_by_ci`.

Current HEAD refresh 2026-05-20: Pre-release verification #16 on branch `main` at commit `73ad2ff` also passed in GitHub Actions. The `cross-db-migration` job succeeded for the current reviewed HEAD.

The original local blocker remains historical evidence of a local service limitation: disposable MySQL and PostgreSQL services were unavailable in the shell where the local audit ran. CI provided the required SQLite, MySQL, and PostgreSQL migration evidence.

Production readiness is not granted by this waiver closure. Keep a production preflight requirement for staging manual verification using `docs/STAGING_VERIFICATION_RUNBOOK.md`, environment-variable review, real deployment topology review, and manual security sign-off.
