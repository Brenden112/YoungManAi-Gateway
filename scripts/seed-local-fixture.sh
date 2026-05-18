#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:3000}"
COOKIE_JAR="${COOKIE_JAR:-/tmp/new-api-fixture-cookies.txt}"
ADMIN_USERNAME="${ADMIN_USERNAME:-root}"
ADMIN_PASSWORD="${ADMIN_PASSWORD:-123456}"
FAKE_UPSTREAM_BASE_URL="${FAKE_UPSTREAM_BASE_URL:-http://fake-upstream:4010}"

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || {
    echo "missing required command: $1" >&2
    exit 1
  }
}

api_cookie() {
  local method="$1" path="$2"
  shift 2
  curl -fsS -b "$COOKIE_JAR" -c "$COOKIE_JAR" -X "$method" "$BASE_URL$path" \
    -H "Content-Type: application/json" "$@"
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

api_cookie POST /api/user/login \
  -d "{\"username\":\"$ADMIN_USERNAME\",\"password\":\"$ADMIN_PASSWORD\"}" >/dev/null

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

echo "Fixture seeded. Run:"
echo "ADMIN_TOKEN=$ADMIN_TOKEN BASE_URL=$BASE_URL bash scripts/regression.sh"
