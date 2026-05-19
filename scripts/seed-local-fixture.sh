#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:3000}"
COOKIE_JAR="${COOKIE_JAR:-/tmp/new-api-fixture-cookies.txt}"
ADMIN_USERNAME="${ADMIN_USERNAME:-root}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-RegtestRoot123!}"
LEGACY_ADMIN_PASSWORD="${LEGACY_ADMIN_PASSWORD:-123456}"
FAKE_UPSTREAM_BASE_URL="${FAKE_UPSTREAM_BASE_URL:-http://fake-upstream:4010}"
ADMIN_USER_ID=""

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1" >&2
    exit 1
  }
}

api_cookie() {
  local method="$1" path="$2"
  shift 2
  local headers=(-H "Content-Type: application/json")
  if [ -n "${ADMIN_USER_ID:-}" ]; then
    headers+=(-H "New-Api-User: $ADMIN_USER_ID")
  fi
  curl -fsS -b "$COOKIE_JAR" -c "$COOKIE_JAR" -X "$method" "$BASE_URL$path" \
    "${headers[@]}" "$@"
}

api_no_cookie() {
  local method="$1" path="$2"
  shift 2
  curl -fsS -X "$method" "$BASE_URL$path" \
    -H "Content-Type: application/json" "$@"
}

login_admin() {
  local password="$1"
  local login_response
  local user_id

  if ! login_response="$(api_cookie POST /api/user/login \
    -d "{\"username\":\"$ADMIN_USERNAME\",\"password\":\"$password\"}")"; then
    return 1
  fi
  if [ "$(echo "$login_response" | jq -r '.success // false')" != "true" ]; then
    return 1
  fi
  user_id="$(echo "$login_response" | jq -r '.data.id // empty')"
  if [ -z "$user_id" ] || [ "$user_id" = "null" ]; then
    return 1
  fi
  ADMIN_USER_ID="$user_id"
}

require_cmd curl
require_cmd jq

rm -f "$COOKIE_JAR"

for _ in $(seq 1 60); do
  if curl -fsS "$BASE_URL/api/status" >/dev/null 2>&1; then
    break
  fi
  sleep 2
done
curl -fsS "$BASE_URL/api/status" >/dev/null

SETUP_STATUS="$(api_no_cookie GET /api/setup)"
SETUP_DONE="$(echo "$SETUP_STATUS" | jq -r '.data.status // false')"
ROOT_INIT="$(echo "$SETUP_STATUS" | jq -r '.data.root_init // false')"
if [ "$SETUP_DONE" != "true" ]; then
  if [ "$ROOT_INIT" = "true" ]; then
    api_no_cookie POST /api/setup \
      -d '{"SelfUseModeEnabled":true,"DemoSiteEnabled":false}' | jq -e '.success == true' >/dev/null
  else
    api_no_cookie POST /api/setup \
      -d "{
        \"username\":\"$ADMIN_USERNAME\",
        \"password\":\"$ADMIN_PASSWORD\",
        \"confirmPassword\":\"$ADMIN_PASSWORD\",
        \"SelfUseModeEnabled\":true,
        \"DemoSiteEnabled\":false
      }" | jq -e '.success == true' >/dev/null
  fi
fi

if ! login_admin "$ADMIN_PASSWORD"; then
  if ! login_admin "$LEGACY_ADMIN_PASSWORD"; then
    echo "failed to login fixture admin" >&2
    exit 1
  fi
fi

api_cookie POST /api/channel/disable-experimental >/dev/null || true

api_cookie POST /api/channel/ \
  -d "{
    \"mode\":\"single\",
    \"channel\":{
      \"name\":\"fixture-official-openai\",
      \"type\":1,
      \"key\":\"fixture-non-secret-key\",
      \"base_url\":\"$FAKE_UPSTREAM_BASE_URL\",
      \"models\":\"gpt-4o-mini\",
      \"group\":\"default\",
      \"status\":1,
      \"provider_type\":\"official_cloud\"
    }
  }" | jq -e '.success == true' >/dev/null

api_cookie POST /api/channel/ \
  -d "{
    \"mode\":\"single\",
    \"channel\":{
      \"name\":\"fixture-experimental-proxy\",
      \"type\":57,
      \"key\":\"fixture-non-secret-key\",
      \"base_url\":\"$FAKE_UPSTREAM_BASE_URL\",
      \"models\":\"kiro-test-experimental\",
      \"group\":\"default\",
      \"status\":1,
      \"provider_type\":\"experimental_proxy\"
    }
  }" | jq -e '.success == true' >/dev/null

ADMIN_TOKEN="$(api_cookie POST /api/token/ \
  -d '{"name":"fixture-admin-smoke","remain_quota":100000,"unlimited_quota":true,"allow_experimental":true,"allowed_provider_types":"official_cloud,experimental_proxy"}' \
  | jq -r '.data.key')"

if [ -z "$ADMIN_TOKEN" ] || [ "$ADMIN_TOKEN" = "null" ]; then
  echo "failed to create fixture admin token" >&2
  exit 1
fi

echo "Fixture seeded. Export:"
printf 'ADMIN_TOKEN=%q BASE_URL=%q\n' "$ADMIN_TOKEN" "$BASE_URL"
