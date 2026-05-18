# routing-security-worker SKILL

## Role
Implements access control, channel selection logic, and fallback isolation.

## Scope
- `middleware/distributor.go`, `service/channel_select.go`
- `service/group.go`, `controller/relay.go`, `middleware/auth.go`

## Security Invariants — NEVER violate
- `experimental_proxy` channels MUST NEVER be visible to non-internal users
- `official_cloud` failure MUST NEVER auto-fallback to `experimental_proxy`
- Internal users can only call `experimental_proxy` when channel is explicitly enabled
- Error messages must not leak channel implementation details

## Rules
1. Add `IsInternalUser` helper before modifying Distribute middleware
2. `excludeExperimental` flag must be passed through the entire retry loop
3. Run `go test ./middleware/... ./service/...` after every change
4. Read `.factory/library/architecture.md` and `validation-contract.md` before starting

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
