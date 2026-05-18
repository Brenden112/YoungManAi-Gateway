# Frontend Test Infrastructure Waiver

Date: 2026-05-17

Status: `pending_ci_verification`

## Blocked Reason

The default frontend has `web/default/package.json` and `web/default/bun.lock`, and now defines `lint`, `test`, and `build` scripts. A minimal `bun test` coverage path was added for `experimental_proxy` channel visibility:

- normal users do not see `experimental_proxy`
- admins can see `experimental_proxy`
- disabled `experimental_proxy` rows remain disabled

The current local environment cannot execute the frontend suite because `bun` is not installed and `web/default/node_modules` is absent. No dependency installation was performed in this environment.

## Minimum Fix Path

1. Install Bun in the runner.
2. From `web/default/`, run `bun install --frozen-lockfile`.
3. Execute `bun run lint`.
4. Execute `bun run test`.
5. Execute `bun run build`.

## CI Commands

```bash
cd web/default
bun install --frozen-lockfile
bun run lint
bun run test
bun run build
```

## CI Coverage

`.github/workflows/pre-release-hardening.yml` contains a `frontend` job that runs the commands above with `oven-sh/setup-bun`.

## Acceptance

This waiver is acceptable only for environments without Bun/dependencies. Release sign-off must review the CI frontend job result before changing deployment readiness to `ready`.

## CI / Staging Verification Path

- `scripts/ci-verify.sh`
- `.github/workflows/pre-release-verification.yml` job `frontend-check`
- `docs/STAGING_VERIFICATION_RUNBOOK.md`

This waiver remains open until frontend lint, test, and build pass in CI/staging. Missing scripts or dependency failures must remain blocked and must not be marked passed.
