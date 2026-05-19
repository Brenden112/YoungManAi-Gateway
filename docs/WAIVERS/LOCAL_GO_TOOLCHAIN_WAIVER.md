# Local Go Toolchain Waiver

Date: 2026-05-18

Status: `closed_by_ci`

blocker_id: `final_go_verification_blocked`

## Current Finding

The current shell does not have a usable Go binary:

- `which go` returned no path.
- `go version` failed with `go: command not found`.
- `/tmp/go2510/bin/go` does not exist.
- `find /tmp -path "*/bin/go" -type f` found no Go binary.

## Impact

The final verification refresh could not execute:

- `go test ./model/... -count=1`
- `go test ./middleware/... -count=1`
- `go test ./... -count=1`
- `go vet ./...`
- `LOCAL_FIXTURE=1 bash scripts/regression.sh`

This is a local environment/toolchain blocker, not evidence of a business-code failure.

## Prior Evidence

Earlier remediation records show `go test ./...` and `go vet ./...` passed before the local Go toolchain disappeared from the current shell/session. Those prior results are not upgraded to final-refresh evidence; CI/staging must rerun the checks.

## Required CI / Staging Verification

```bash
go version
go test ./model/... -count=1
go test ./middleware/... -count=1
go test ./... -count=1
go vet ./...
LOCAL_FIXTURE=1 bash scripts/regression.sh
```

## Pass Criteria

- Go 1.22+ is available.
- All listed Go tests pass.
- `go vet ./...` exits `0`.
- Local fixture regression exits `0` without real upstream provider calls or real API keys.

## Failure Handling

- If Go remains unavailable in CI/staging, fix the runner image/toolchain before production review.
- If Go is available and tests fail, classify the failure as code, test dependency, or environment based on the failing output.
- Keep production release blocked on staging manual verification even after CI passes final Go verification.

## CI / Staging Verification Path

- `scripts/ci-verify.sh`
- `.github/workflows/pre-release-verification.yml` job `go-test-vet`
- `.github/workflows/pre-release-verification.yml` job `local-fixture-regression`

## CI Closure — 2026-05-19

Pre-release verification #13 on branch `main` at commit `aeb43e5` passed in GitHub Actions. The `go-test-vet` and `local-fixture-regression` jobs succeeded, closing local Go toolchain blocker `final_go_verification_blocked` as `closed_by_ci`.

The original local shell was still a valid environment blocker: Go was unavailable locally, so the local refresh could not run the final checks. CI provided the required Go 1.22+ runner evidence.

Production readiness is not granted by this waiver closure. Keep a production preflight requirement for staging manual verification using `docs/STAGING_VERIFICATION_RUNBOOK.md`, environment-variable review, real deployment topology review, and manual security sign-off.
