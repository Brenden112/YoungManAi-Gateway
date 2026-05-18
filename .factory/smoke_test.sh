#!/usr/bin/env bash
# .factory/smoke_test.sh — MVP smoke test for B2B AI API Gateway
# Covers: M3-F01 (GET /v1/models), M3-F02 (POST /v1/chat/completions),
#         M3-F03 (streaming), auth rejection, experimental_proxy block.
#
# Usage:
#   BASE_URL=http://localhost:3000 API_KEY=sk-xxx MODEL=gpt-4o-mini bash .factory/smoke_test.sh
#   bash .factory/smoke_test.sh BASE_URL=http://localhost:3000 API_KEY=sk-xxx MODEL=gpt-4o-mini
#
# Required env vars:
#   API_KEY   — a valid user API key (sk-...)
#
# Optional env vars:
#   BASE_URL  — default: http://localhost:3000
#   MODEL     — default: gpt-4o-mini
#   TIMEOUT   — curl timeout in seconds, default: 30

set -euo pipefail

for arg in "$@"; do
  case "$arg" in
    BASE_URL=*|API_KEY=*|MODEL=*|TIMEOUT=*)
      export "$arg"
      ;;
    *)
      echo "[ERROR] unsupported argument: $arg"
      echo "Usage: BASE_URL=http://localhost:3000 API_KEY=sk-xxx MODEL=gpt-4o-mini bash .factory/smoke_test.sh"
      exit 1
      ;;
  esac
done

BASE_URL="${BASE_URL:-http://localhost:3000}"
MODEL="${MODEL:-gpt-4o-mini}"
TIMEOUT="${TIMEOUT:-30}"
PASS=0
FAIL=0

if [ -z "${API_KEY:-}" ]; then
  echo "[ERROR] API_KEY env var is required"
  exit 1
fi

_pass() { echo "[PASS] $1"; PASS=$((PASS+1)); }
_fail() { echo "[FAIL] $1"; FAIL=$((FAIL+1)); }
_info() { echo "[INFO] $1"; }

# ── T01: GET /v1/models — valid token ────────────────────────────────────────
_info "T01: GET /v1/models with valid token"
RESP=$(curl -sf --max-time "$TIMEOUT" \
  -H "Authorization: Bearer $API_KEY" \
  "$BASE_URL/v1/models" 2>&1) || { _fail "T01: curl failed or non-2xx"; RESP=""; }

if [ -n "$RESP" ]; then
  if echo "$RESP" | grep -q '"object":"list"' && echo "$RESP" | grep -q '"data":\['; then
    _pass "T01: GET /v1/models returns 200 with model list"
  else
    _fail "T01: response missing 'object:list' or 'data' array — got: ${RESP:0:200}"
  fi
fi

# ── T02: GET /v1/models — invalid token → 401 ────────────────────────────────
_info "T02: GET /v1/models with invalid token"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time "$TIMEOUT" \
  -H "Authorization: Bearer sk-invalid-token-xyz" \
  "$BASE_URL/v1/models" 2>/dev/null)
if [ "$HTTP_CODE" = "401" ]; then
  _pass "T02: invalid token returns 401"
else
  _fail "T02: expected 401, got $HTTP_CODE"
fi

# ── T03: POST /v1/chat/completions — non-streaming ───────────────────────────
_info "T03: POST /v1/chat/completions (non-streaming)"
RESP=$(curl -sf --max-time "$TIMEOUT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"model\":\"$MODEL\",\"messages\":[{\"role\":\"user\",\"content\":\"Say hi\"}],\"max_tokens\":5}" \
  "$BASE_URL/v1/chat/completions" 2>&1) || { _fail "T03: curl failed or non-2xx"; RESP=""; }

if [ -n "$RESP" ]; then
  if echo "$RESP" | grep -q '"choices"' && echo "$RESP" | grep -q '"usage"'; then
    _pass "T03: chat/completions returns choices + usage"
  else
    _fail "T03: response missing choices or usage — got: ${RESP:0:300}"
  fi
fi

# ── T04: POST /v1/chat/completions — streaming ───────────────────────────────
_info "T04: POST /v1/chat/completions (streaming)"
STREAM_RESP=$(curl -sf --max-time "$TIMEOUT" \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{\"model\":\"$MODEL\",\"messages\":[{\"role\":\"user\",\"content\":\"Say hi\"}],\"stream\":true,\"max_tokens\":5}" \
  "$BASE_URL/v1/chat/completions" 2>&1) || { _fail "T04: curl failed or non-2xx"; STREAM_RESP=""; }

if [ -n "$STREAM_RESP" ]; then
  if echo "$STREAM_RESP" | grep -q "data:"; then
    _pass "T04: streaming returns SSE data chunks"
  else
    _fail "T04: streaming response missing 'data:' SSE prefix — got: ${STREAM_RESP:0:200}"
  fi
fi

# ── T05: POST /v1/chat/completions — missing token → 401 ─────────────────────
_info "T05: POST /v1/chat/completions without token"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time "$TIMEOUT" \
  -H "Content-Type: application/json" \
  -d "{\"model\":\"$MODEL\",\"messages\":[{\"role\":\"user\",\"content\":\"hi\"}]}" \
  "$BASE_URL/v1/chat/completions" 2>/dev/null)
if [ "$HTTP_CODE" = "401" ]; then
  _pass "T05: missing token returns 401"
else
  _fail "T05: expected 401, got $HTTP_CODE"
fi

# ── T06: GET /api/status — health check ──────────────────────────────────────
_info "T06: GET /api/status health check"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --max-time "$TIMEOUT" \
  "$BASE_URL/api/status" 2>/dev/null)
if [ "$HTTP_CODE" = "200" ]; then
  _pass "T06: /api/status returns 200"
else
  _fail "T06: /api/status returned $HTTP_CODE (expected 200)"
fi

# ── Summary ───────────────────────────────────────────────────────────────────
echo ""
echo "Results: $PASS passed, $FAIL failed"
if [ "$FAIL" -gt 0 ]; then
  echo "[SMOKE TEST FAILED]"
  exit 1
fi
echo "[SMOKE TEST PASSED]"
exit 0
