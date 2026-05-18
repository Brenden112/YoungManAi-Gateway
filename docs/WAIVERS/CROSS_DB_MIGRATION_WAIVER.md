# Cross-DB Migration Waiver

Date: 2026-05-17

Status: `pending_ci_verification`

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

This waiver is acceptable only for local environments without MySQL/PostgreSQL services. Deployment readiness must remain `needs_manual_review` until CI artifacts show all three migration jobs passing.

## CI / Staging Verification Path

- `scripts/ci-migration-check.sh`
- `.github/workflows/pre-release-verification.yml` job `cross-db-migration`
- `docs/STAGING_VERIFICATION_RUNBOOK.md`

This waiver remains open until SQLite, MySQL, and PostgreSQL migration checks pass in CI/staging.
