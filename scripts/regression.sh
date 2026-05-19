#!/usr/bin/env bash
# M16-F04 Regression Test Suite
# Tests all six security and behaviour invariants against a running gateway.
#
# Usage:
#   BASE_URL=http://localhost:3000 ADMIN_TOKEN=<root-token> bash scripts/regression.sh
#   LOCAL_FIXTURE=1 bash scripts/regression.sh
#
# Prerequisites:
#   - Gateway running (docker compose -f docker-compose.local.yml up -d --build)
#   - curl, jq installed
#   - ADMIN_TOKEN: a token belonging to the root/admin user

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:3000}"
ADMIN_TOKEN="${ADMIN_TOKEN:-}"
ADMIN_USER_ID="${ADMIN_USER_ID:-}"
LOCAL_FIXTURE="${LOCAL_FIXTURE:-0}"
PASS=0
FAIL=0
GO_BIN="${GO_BIN:-}"
if [ -z "$GO_BIN" ]; then
  if [ -x /tmp/go2510/bin/go ]; then
    GO_BIN=/tmp/go2510/bin/go
  else
    GO_BIN=go
  fi
fi

# ── helpers ──────────────────────────────────────────────────────────────────

red()   { printf '\033[31m%s\033[0m\n' "$*"; }
green() { printf '\033[32m%s\033[0m\n' "$*"; }
bold()  { printf '\033[1m%s\033[0m\n' "$*"; }

if [ "$LOCAL_FIXTURE" = "1" ]; then
  bold "=== M16 Local Fixture Regression Suite ==="
  echo "No real upstream provider or API key is used in this mode."
  if ! command -v "$GO_BIN" >/dev/null 2>&1; then
    red "BLOCKED final_go_verification_blocked: Go binary not found ($GO_BIN)."
    exit 2
  fi
  env GOCACHE="${GOCACHE:-/tmp/go-build-cache}" GOPATH="${GOPATH:-/tmp/go-path}" "$GO_BIN" test ./model -run 'TestAUD021|TestAUD022|TestRecordConsumeLogSanitizesOther' -count=1
  env GOCACHE="${GOCACHE:-/tmp/go-build-cache}" GOPATH="${GOPATH:-/tmp/go-path}" "$GO_BIN" test ./middleware -run 'TestAUD021|TestAUD022' -count=1
  env GOCACHE="${GOCACHE:-/tmp/go-build-cache}" GOPATH="${GOPATH:-/tmp/go-path}" "$GO_BIN" test ./controller -run 'TestAUD025|TestAdminManualTopUp|TestNormalUserManualTopUp|TestAdminManualTopUpRejects' -count=1
  env GOCACHE="${GOCACHE:-/tmp/go-build-cache}" GOPATH="${GOPATH:-/tmp/go-path}" "$GO_BIN" test ./service -run 'TestPreConsumeQuotaRejectsOnZeroBalance|TestPreConsumeQuotaNoErrorLogOnInsufficientBalance' -count=1
  env GOCACHE="${GOCACHE:-/tmp/go-build-cache}" GOPATH="${GOPATH:-/tmp/go-path}" "$GO_BIN" vet ./...
  green "Local fixture regression checks passed."
  exit 0
fi

assert_http() {
  local label="$1" expected="$2" actual="$3"
  if [ "$actual" = "$expected" ]; then
    green "  PASS  $label (HTTP $actual)"
    PASS=$((PASS+1))
  else
    red   "  FAIL  $label — expected HTTP $expected, got HTTP $actual"
    FAIL=$((FAIL+1))
  fi
}

assert_contains() {
  local label="$1" needle="$2" haystack="$3"
  if echo "$haystack" | grep -q "$needle"; then
    green "  PASS  $label"
    PASS=$((PASS+1))
  else
    red   "  FAIL  $label — expected '$needle' in response"
    FAIL=$((FAIL+1))
  fi
}

assert_not_contains() {
  local label="$1" needle="$2" haystack="$3"
  if ! echo "$haystack" | grep -q "$needle"; then
    green "  PASS  $label"
    PASS=$((PASS+1))
  else
    red   "  FAIL  $label — '$needle' should NOT appear in response"
    FAIL=$((FAIL+1))
  fi
}

api() {
  local method="$1" path="$2" token="$3"
  shift 3
  local headers=(-H "Authorization: Bearer $token" -H "Content-Type: application/json")
  if [[ "$path" == /api/* && "$token" = "$ADMIN_TOKEN" && -n "$ADMIN_USER_ID" ]]; then
    headers+=(-H "New-Api-User: $ADMIN_USER_ID")
  fi
  curl -s -X "$method" "$BASE_URL$path" \
    "${headers[@]}" \
    "$@"
}

api_cookie() {
  local method="$1" path="$2" cookie_jar="$3"
  shift 3
  curl -s -b "$cookie_jar" -c "$cookie_jar" -X "$method" "$BASE_URL$path" \
    -H "Content-Type: application/json" \
    "$@"
}

api_cookie_user() {
  local method="$1" path="$2" cookie_jar="$3" user_id="$4"
  shift 4
  curl -s -b "$cookie_jar" -c "$cookie_jar" -X "$method" "$BASE_URL$path" \
    -H "Content-Type: application/json" \
    -H "New-Api-User: $user_id" \
    "$@"
}

http_code() {
  local method="$1" path="$2" token="$3"
  shift 3
  curl -s -o /dev/null -w "%{http_code}" -X "$method" "$BASE_URL$path" \
    -H "Authorization: Bearer $token" \
    -H "Content-Type: application/json" \
    "$@"
}

# ── preflight ─────────────────────────────────────────────────────────────────

bold "=== M16 Regression Suite ==="
echo "BASE_URL: $BASE_URL"

if [ -z "$ADMIN_TOKEN" ]; then
  red "ERROR: ADMIN_TOKEN is required. Export it before running this script."
  exit 1
fi

# Wait for gateway to be ready
echo "Waiting for gateway..."
for i in $(seq 1 20); do
  if curl -sf "$BASE_URL/api/status" > /dev/null 2>&1; then
    break
  fi
  sleep 2
done
curl -sf "$BASE_URL/api/status" > /dev/null || { red "Gateway not reachable at $BASE_URL"; exit 1; }
echo "Gateway ready."

# ── setup: create test users and tokens ──────────────────────────────────────

bold "\n--- Setup: creating test users ---"

# Normal user
NORMAL_USER=$(api POST /api/user/register "" -d '{"username":"regtest_normal","password":"Regtest123!","name":"Regression Normal"}' 2>/dev/null || true)
NORMAL_USER_ID=$(api GET /api/user/search?keyword=regtest_normal "$ADMIN_TOKEN" | jq -r '.data.items[0].id // empty')

# Internal user (group=internal)
INTERNAL_USER=$(api POST /api/user/register "" -d '{"username":"regtest_internal","password":"Regtest123!","name":"Regression Internal"}' 2>/dev/null || true)
INTERNAL_USER_ID=$(api GET /api/user/search?keyword=regtest_internal "$ADMIN_TOKEN" | jq -r '.data.items[0].id // empty')
if [ -n "$INTERNAL_USER_ID" ]; then
  api PUT /api/user "$ADMIN_TOKEN" -d "{\"id\":$INTERNAL_USER_ID,\"username\":\"regtest_internal\",\"display_name\":\"regtest_internal\",\"group\":\"internal\"}" > /dev/null
fi

# Zero-quota user
ZERO_USER=$(api POST /api/user/register "" -d '{"username":"regtest_zero","password":"Regtest123!","name":"Regression Zero"}' 2>/dev/null || true)
ZERO_USER_ID=$(api GET /api/user/search?keyword=regtest_zero "$ADMIN_TOKEN" | jq -r '.data.items[0].id // empty')

if [ -n "$NORMAL_USER_ID" ]; then
  api POST /api/user/manage "$ADMIN_TOKEN" -d "{\"id\":$NORMAL_USER_ID,\"action\":\"add_quota\",\"mode\":\"add\",\"value\":100000}" > /dev/null
fi
if [ -n "$INTERNAL_USER_ID" ]; then
  api POST /api/user/manage "$ADMIN_TOKEN" -d "{\"id\":$INTERNAL_USER_ID,\"action\":\"add_quota\",\"mode\":\"add\",\"value\":100000}" > /dev/null
fi

# Tokens
NORMAL_COOKIE="$(mktemp)"
INTERNAL_COOKIE="$(mktemp)"
ZERO_COOKIE="$(mktemp)"
trap 'rm -f "$NORMAL_COOKIE" "$INTERNAL_COOKIE" "$ZERO_COOKIE"' EXIT

NORMAL_LOGIN_ID=$(api_cookie POST /api/user/login "$NORMAL_COOKIE" -d '{"username":"regtest_normal","password":"Regtest123!"}' | jq -r '.data.id // empty')
INTERNAL_LOGIN_ID=$(api_cookie POST /api/user/login "$INTERNAL_COOKIE" -d '{"username":"regtest_internal","password":"Regtest123!"}' | jq -r '.data.id // empty')
ZERO_LOGIN_ID=$(api_cookie POST /api/user/login "$ZERO_COOKIE" -d '{"username":"regtest_zero","password":"Regtest123!"}' | jq -r '.data.id // empty')

NORMAL_TOKEN_KEY=$(api_cookie_user POST /api/token "$NORMAL_COOKIE" "$NORMAL_LOGIN_ID" -d '{"name":"regtest-normal","remain_quota":100000,"allowed_provider_types":"official_cloud"}' | jq -r '.data.key // empty')
INTERNAL_TOKEN_KEY=$(api_cookie_user POST /api/token "$INTERNAL_COOKIE" "$INTERNAL_LOGIN_ID" -d '{"name":"regtest-internal","remain_quota":100000,"allow_experimental":true,"allowed_provider_types":"experimental_proxy"}' | jq -r '.data.key // empty')
ZERO_TOKEN_KEY=$(api_cookie_user POST /api/token "$ZERO_COOKIE" "$ZERO_LOGIN_ID" -d '{"name":"regtest-zero","remain_quota":0,"unlimited_quota":true,"allowed_provider_types":"official_cloud"}' | jq -r '.data.key // empty')

echo "Normal user ID:   ${NORMAL_USER_ID:-<not found>}"
echo "Internal user ID: ${INTERNAL_USER_ID:-<not found>}"
echo "Zero user ID:     ${ZERO_USER_ID:-<not found>}"

# ── T1: official_cloud channel accessible to normal user ─────────────────────

bold "\n--- T1: official_cloud channel accessible to normal user ---"
# /v1/models should return at least one model for a normal user
MODELS_RESP=$(http_code GET /v1/models "$NORMAL_TOKEN_KEY")
assert_http "T1 normal user can list models" "200" "$MODELS_RESP"

CHAT_CODE=$(http_code POST /v1/chat/completions "$NORMAL_TOKEN_KEY" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"fixture smoke"}]}')
assert_http "T1 normal user can chat through official provider" "200" "$CHAT_CODE"

# ── T2: experimental_proxy models hidden from normal user ────────────────────

bold "\n--- T2: experimental_proxy models hidden from normal user ---"
MODELS_BODY=$(api GET /v1/models "$NORMAL_TOKEN_KEY")
# experimental_proxy channels should not appear in normal user model list
# (they are filtered by GetGroupEnabledModelsExcludingExperimental)
assert_not_contains "T2 experimental models absent for normal user" "experimental" "$MODELS_BODY"

# ── T3: normal user rejected when routed to experimental_proxy channel ────────

bold "\n--- T3: normal user rejected on experimental_proxy channel ---"
# Attempt to call a model that only exists on an experimental_proxy channel.
# The gateway should return 403 at the post-selection guard.
# We use a sentinel model name that maps only to experimental channels.
EXP_CODE=$(http_code POST /v1/chat/completions "$NORMAL_TOKEN_KEY" \
  -d '{"model":"kiro-test-experimental","messages":[{"role":"user","content":"hi"}]}')
assert_http "T3 normal user gets 403 on experimental model" "403" "$EXP_CODE"

# ── T4: disabled experimental_proxy channel rejected ─────────────────────────

bold "\n--- T4: disabled experimental_proxy channel rejected ---"
# Disable all experimental_proxy channels via admin endpoint
DISABLE_RESP=$(api POST /api/channel/disable-experimental "$ADMIN_TOKEN")
assert_contains "T4 disable-experimental returns success" '"success":true' "$DISABLE_RESP"

# Internal user with allow_experimental token should now get 503 (no available channel)
DISABLED_CODE=$(http_code POST /v1/chat/completions "$INTERNAL_TOKEN_KEY" \
  -d '{"model":"kiro-test-experimental","messages":[{"role":"user","content":"hi"}]}')
assert_http "T4 disabled experimental channel returns 503" "503" "$DISABLED_CODE"

# ── T5: insufficient balance rejected with 402 ───────────────────────────────

bold "\n--- T5: insufficient balance rejected ---"
QUOTA_CODE=$(http_code POST /v1/chat/completions "$ZERO_TOKEN_KEY" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}')
assert_http "T5 zero-quota user gets 402" "402" "$QUOTA_CODE"

# ── T6: default no prompt/response storage ───────────────────────────────────

bold "\n--- T6: default no prompt/response storage ---"
# Check that STORE_FULL_TEXT_ENABLED is false by inspecting a recent log entry.
# After a request, the log content field should be empty/null.
# We make a real request first (use normal user against a real model if available).
# Then check the last log entry for this token.
LOG_RESP=$(api GET "/api/log?token_name=regtest-normal&page=1&page_size=1" "$ADMIN_TOKEN")
LOG_CONTENT=$(echo "$LOG_RESP" | jq -r '.data.items[0].content // ""')
if [ -z "$LOG_CONTENT" ] || [ "$LOG_CONTENT" = "null" ]; then
  green "  PASS  T6 prompt/response not stored in log (content is empty)"
  PASS=$((PASS+1))
else
  red   "  FAIL  T6 log content should be empty when STORE_FULL_TEXT_ENABLED=false, got: $LOG_CONTENT"
  FAIL=$((FAIL+1))
fi

# ── summary ──────────────────────────────────────────────────────────────────

bold "\n=== Results ==="
echo "Passed: $PASS"
echo "Failed: $FAIL"
echo "Total:  $((PASS+FAIL))"

if [ "$FAIL" -gt 0 ]; then
  red "\nSome tests FAILED."
  exit 1
else
  green "\nAll tests PASSED."
  exit 0
fi
