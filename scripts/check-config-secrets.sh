#!/usr/bin/env bash
set -euo pipefail

FILES=(
  "docker-compose.yml"
  ".env.example"
  ".env.staging.example"
)

ALLOW_RE='(change-me|example|fake|local-fixture|fixture|dummy|placeholder|your-|localhost|127\.0\.0\.1|\$\{|<[^>]+>)'
HIGH_RISK_RE='(sk-[A-Za-z0-9_-]{12,}|Bearer[[:space:]]+[A-Za-z0-9._~+/-]{12,}|api_key[[:space:]]*=|provider_key[[:space:]]*[:=]|access_token[[:space:]]*[:=]|refresh_token[[:space:]]*[:=])'
SECRET_NAME_RE='(password|passwd|pwd|secret|token|credential|provider[_-]?key|api[_-]?key|private[_-]?key|encryption[_-]?key)'
WEAK_VALUE_RE='(^|[^0-9])123456([^0-9]|$)|rootpass|postgres:postgres|root:rootpass'

fail=0

check_file_exists() {
  local file="$1"
  if [ ! -f "$file" ]; then
    printf 'FAIL missing required config file: %s\n' "$file" >&2
    fail=1
  fi
}

report_line() {
  local file="$1"
  local line_no="$2"
  local reason="$3"
  local line="$4"
  printf 'FAIL %s:%s %s\n' "$file" "$line_no" "$reason" >&2
  printf '     %s\n' "$line" >&2
  fail=1
}

for file in "${FILES[@]}"; do
  check_file_exists "$file"
done

for file in "${FILES[@]}"; do
  [ -f "$file" ] || continue
  line_no=0
  while IFS= read -r line || [ -n "$line" ]; do
    line_no=$((line_no + 1))
    trimmed="${line#"${line%%[![:space:]]*}"}"
    [ -z "$trimmed" ] && continue
    [[ "$trimmed" == \#* ]] && continue
    [[ "$trimmed" =~ (_ENDPOINT|_URL)= ]] && continue

    if [[ "$line" =~ $HIGH_RISK_RE ]] && ! [[ "$line" =~ $ALLOW_RE ]]; then
      report_line "$file" "$line_no" "contains high-risk secret pattern" "$line"
      continue
    fi

    if [[ "$line" =~ $WEAK_VALUE_RE ]] && ! [[ "$line" =~ $ALLOW_RE ]]; then
      report_line "$file" "$line_no" "contains weak/example credential that is not marked placeholder" "$line"
      continue
    fi

    if [[ "$line" =~ $SECRET_NAME_RE ]] && ! [[ "$line" =~ $ALLOW_RE ]]; then
      report_line "$file" "$line_no" "contains secret-like assignment without an allowed placeholder or env interpolation" "$line"
      continue
    fi
  done < "$file"
done

if [ "$fail" -ne 0 ]; then
  printf 'Config secret check failed. Replace real-looking values with env interpolation or placeholders.\n' >&2
  exit 1
fi

printf 'Config secret check passed.\n'
