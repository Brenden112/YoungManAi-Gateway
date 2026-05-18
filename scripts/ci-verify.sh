#!/usr/bin/env bash
set -euo pipefail

PASS=0
FAIL=0
BLOCKED=0

green() { printf 'PASS    %s\n' "$*"; }
red() { printf 'FAIL    %s\n' "$*"; }
yellow() { printf 'BLOCKED %s\n' "$*"; }

run_check() {
  local label="$1"
  shift
  if "$@"; then
    green "$label"
    PASS=$((PASS + 1))
  else
    red "$label"
    FAIL=$((FAIL + 1))
  fi
}

block_check() {
  local label="$1"
  local reason="$2"
  yellow "$label: $reason"
  BLOCKED=$((BLOCKED + 1))
}

has_script() {
  local pkg="$1"
  local script="$2"
  node -e "const p=require(process.argv[1]); process.exit(p.scripts && p.scripts[process.argv[2]] ? 0 : 1)" "$pkg" "$script" >/dev/null 2>&1
}

detect_package_manager() {
  local dir="$1"
  if command -v bun >/dev/null 2>&1; then
    printf 'bun'
  elif command -v pnpm >/dev/null 2>&1; then
    printf 'pnpm'
  elif command -v npm >/dev/null 2>&1; then
    printf 'npm'
  elif command -v yarn >/dev/null 2>&1; then
    printf 'yarn'
  else
    return 1
  fi
}

run_package_script() {
  local dir="$1"
  local manager="$2"
  local script="$3"
  case "$manager" in
    bun) (cd "$dir" && bun run "$script") ;;
    pnpm) (cd "$dir" && pnpm run "$script") ;;
    npm) (cd "$dir" && npm run "$script") ;;
    yarn) (cd "$dir" && yarn run "$script") ;;
    *) return 1 ;;
  esac
}

install_frontend_deps_if_needed() {
  local dir="$1"
  local manager="$2"
  if [ -d "$dir/node_modules" ]; then
    return 0
  fi
  if [ "${CI:-}" != "true" ]; then
    return 2
  fi
  case "$manager" in
    bun) (cd "$dir" && bun install --frozen-lockfile) ;;
    pnpm) (cd "$dir" && pnpm install --frozen-lockfile) ;;
    npm) (cd "$dir" && npm ci) ;;
    yarn) (cd "$dir" && yarn install --frozen-lockfile) ;;
    *) return 1 ;;
  esac
}

echo "=== CI verification ==="
echo "No real upstream provider or API key is used by this script."

if command -v go >/dev/null 2>&1; then
  go version
  run_check "go test ./model/... -count=1" go test ./model/... -count=1
  run_check "go test ./middleware/... -count=1" go test ./middleware/... -count=1
  run_check "go test ./... -count=1" go test ./... -count=1
  run_check "go vet ./..." go vet ./...
  run_check "LOCAL_FIXTURE=1 bash scripts/regression.sh" env LOCAL_FIXTURE=1 bash scripts/regression.sh
else
  block_check "final_go_verification_blocked" "go binary not found in PATH"
  block_check "go test ./model/... -count=1" "go binary not found in PATH"
  block_check "go test ./middleware/... -count=1" "go binary not found in PATH"
  block_check "go test ./... -count=1" "go binary not found in PATH"
  block_check "go vet ./..." "go binary not found in PATH"
  block_check "LOCAL_FIXTURE=1 bash scripts/regression.sh" "go binary not found in PATH"
fi

run_check "git diff --check" git diff --check

frontend_dirs=()
[ -f package.json ] && frontend_dirs+=(".")
[ -f web/default/package.json ] && frontend_dirs+=("web/default")

if [ "${#frontend_dirs[@]}" -eq 0 ]; then
  block_check "blocked_test_infra_frontend" "no package.json found"
else
  for dir in "${frontend_dirs[@]}"; do
    pkg="$dir/package.json"
    if ! command -v node >/dev/null 2>&1; then
      block_check "blocked_test_infra_frontend ($dir)" "node is not available to inspect package scripts"
      continue
    fi
    if ! manager="$(detect_package_manager "$dir")"; then
      block_check "blocked_test_infra_frontend ($dir)" "no bun, pnpm, npm, or yarn available"
      continue
    fi
    dep_status=0
    install_frontend_deps_if_needed "$dir" "$manager" || dep_status=$?
    if [ "$dep_status" -eq 2 ]; then
      block_check "blocked_test_infra_frontend ($dir)" "frontend dependencies are missing; run dependency install in CI/staging with $manager"
      continue
    elif [ "$dep_status" -ne 0 ]; then
      block_check "blocked_test_infra_frontend ($dir)" "frontend dependencies could not be installed with $manager"
      continue
    fi
    ran_any=0
    for script in lint test build; do
      if has_script "$pkg" "$script"; then
        ran_any=1
        run_check "$dir $manager run $script" run_package_script "$dir" "$manager" "$script"
      else
        block_check "$dir $script" "package script is missing"
      fi
    done
    if [ "$ran_any" -eq 0 ]; then
      block_check "blocked_test_infra_frontend ($dir)" "no lint/test/build scripts are defined"
    fi
  done
fi

echo "=== CI verification summary ==="
echo "passed=$PASS failed=$FAIL blocked=$BLOCKED"

if [ "$FAIL" -gt 0 ]; then
  exit 1
fi
if [ "$BLOCKED" -gt 0 ]; then
  exit 2
fi
