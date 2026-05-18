# User Testing Reference

## Test Users

| Username | Role | Group | Notes |
|----------|------|-------|-------|
| root | root (100) | default | created on first boot, password 123456 |
| admin_test | admin (10) | default | manual top-up, channel management |
| internal_user | common (1) | internal | can call experimental_proxy after explicit enable |
| normal_user | common (1) | default | cannot call experimental_proxy |

## API Key Patterns

All user API keys follow the `sk-` prefix pattern.
Keys are stored as SHA-256 hash in DB; never log or return the raw key after creation.

## Test Scenarios

### 1. Basic Chat Completion (Happy Path)
```bash
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-<user_token>" \
  -H "Content-Type: application/json" \
  -d '{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hello"}]}'
# Expected: 200 with choices array, usage object
```

### 2. Model List
```bash
curl http://localhost:3000/v1/models \
  -H "Authorization: Bearer sk-<user_token>"
# Expected: 200 with data array; experimental_proxy models NOT visible to normal_user
```

### 3. Insufficient Balance Rejection
```bash
# Create token with RemainQuota=0, UnlimitedQuota=false
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-<zero_quota_token>" ...
# Expected: 402 or 429 with insufficient_quota error
```

### 4. experimental_proxy — Normal User Blocked
```bash
# normal_user tries to call a model served only by experimental_proxy channel
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-<normal_user_token>" \
  -d '{"model":"kiro-test-model",...}'
# Expected: 403 forbidden
```

### 5. experimental_proxy — Internal User Allowed
```bash
# internal_user with explicit enable
curl -X POST http://localhost:3000/v1/chat/completions \
  -H "Authorization: Bearer sk-<internal_user_token>" \
  -d '{"model":"kiro-test-model",...}'
# Expected: 200 (channel enabled + user is internal)
```

### 6. experimental_proxy One-Click Disable
```bash
# Admin disables the experimental_proxy channel
PATCH /api/channel/:id  {"status": 2}
# Then internal_user call → Expected: 503 or 404 no available channel
```

### 7. Admin Manual Top-up
```bash
POST /api/user/topup  {"user_id": X, "quota": 100000}
# Expected: user.Quota increases by 100000, log entry type=manage recorded
```

### 8. usage_log Privacy
```bash
GET /api/log?token_id=X
# Expected: no prompt_text, no completion_text fields in response
# prompt_tokens and completion_tokens counts are present
```

## Acceptance Checklist per Milestone

- M3: curl /v1/models returns 200; curl /v1/chat/completions returns 200 with usage
- M7: normal_user blocked from experimental_proxy; internal_user allowed
- M8: disabling experimental_proxy channel blocks all calls; no auto-fallback
- M12: after request, user.UsedQuota increases; token.UsedQuota increases; log.Quota > 0
- M13: zero-quota token gets 402/429; admin top-up restores access
- M16: all above pass in CI smoke test

---

## Mandatory Feature Completion Rule

After completing **any** feature, the Worker MUST update both files before returning the final handoff:

- `docs/DEVELOPMENT_LOG.md`
- `.factory/mission-state.json`

**Do NOT claim a feature is completed unless both files are updated.**

The final handoff must explicitly include:
```json
{
  "development_log_updated": true,
  "mission_state_updated": true
}
```

Full handoff templates and current mission state are in `.factory/mission-state.json`.
This rule applies to all milestones M0–M16 and all workers.
