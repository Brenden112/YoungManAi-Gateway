# admin-ui-worker SKILL

## Role
Implements admin backend API endpoints and minimal frontend pages for
Provider/Channel management, API Key management, usage_log, and balance.

## Scope
- Backend: `controller/channel.go`, `controller/token.go`, `controller/log.go`,
  `controller/user.go`, `router/api-router.go`
- Frontend: `web/default/src/pages/` (React 19, Tailwind, Base UI, Bun)

## Rules
1. All admin endpoints must require admin role (`middleware.AdminAuth`)
2. `experimental_proxy` channels must never appear in non-admin responses
3. Token list must return masked keys only (`model.MaskTokenKey`)
4. `usage_log` must not expose prompt/completion text
5. Frontend: use `bun` as package manager; follow existing component patterns
6. i18n: add translation keys to `web/default/src/i18n/locales/` for new UI text
7. Run `go test ./controller/...` and `bun run build` after changes
8. Check `validation-contract.md` assertions A14.x–A15.x before marking done

---

## Mandatory Feature Completion Rule

After completing any feature, you MUST update both files before returning the handoff:

- `docs/DEVELOPMENT_LOG.md`
- `.factory/mission-state.json`

**Do NOT claim a feature is completed unless both files are updated.**

### What to record

1. feature_id
2. milestone
3. worker name
4. feature status: `completed` / `blocked` / `failed`
5. files_inspected
6. files_modified
7. implementation_summary
8. validation_assertions_fulfilled
9. commands_run (command, exit_code, observation)
10. tests_added
11. manual_checks
12. risks_or_todos
13. breaking_changes
14. blocker and minimal_fix_path (if blocked)
15. next_recommended_feature
16. updated_at timestamp

### On Completion

- Add to `completed_features` in `.factory/mission-state.json`
- Remove from `pending_features`
- Update `validation_status` for fulfilled assertions
- Append change record to `docs/DEVELOPMENT_LOG.md`
- Update `current_milestone`, `current_feature`, `next_recommended_action`

### On Blocked

- Add to `blocked_features`; keep out of `completed_features`
- Record blocker, impact, related files, `minimal_fix_path`
- Do NOT continue to next feature automatically

### On Failed

- Mark `status: failed`; record failure reason and failed commands with exit codes
- Do NOT continue to next feature automatically

### Handoff Template — completed

```json
{
  "feature_id": "",
  "milestone": "",
  "worker": "",
  "status": "completed",
  "files_inspected": [],
  "files_modified": [],
  "implementation_summary": "",
  "validation_assertions_fulfilled": [],
  "commands_run": [{"command": "", "exit_code": 0, "observation": ""}],
  "tests_added": [],
  "manual_checks": [],
  "risks_or_todos": [],
  "breaking_changes": [],
  "development_log_updated": true,
  "mission_state_updated": true,
  "next_recommended_feature": ""
}
```

### Handoff Template — blocked

```json
{
  "feature_id": "",
  "milestone": "",
  "worker": "",
  "status": "blocked",
  "blocker": "",
  "why_it_blocks_the_mission": "",
  "minimal_fix_path": "",
  "files_related": [],
  "development_log_updated": true,
  "mission_state_updated": true,
  "recommended_next_action": ""
}
```
