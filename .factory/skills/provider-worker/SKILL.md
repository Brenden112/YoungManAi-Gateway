# provider-worker SKILL

## Role
Implements Provider Adapter Framework: channel type constants, ProviderAccount model,
Adaptor interface implementations, and GetAdaptor factory registration.

## Scope
- `constant/channel.go`, `model/channel.go`, `model/provider_account.go`
- `relay/channel/<new>/`, `relay/relay_adaptor.go`
- `controller/channel.go`, `dto/`

## Rules
1. New channel type constants must be added before `ChannelTypeDummy`
2. Always add compile-time interface check: `var _ channel.Adaptor = (*YourAdaptor)(nil)`
3. Upstream keys must be encrypted via `common/crypto.go` — never store plaintext
4. Use `common.Marshal`/`Unmarshal` (not `encoding/json`) per `CLAUDE.md` Rule 1
5. Optional scalar fields in request DTOs must use pointer types with `omitempty` (Rule 6)
6. Run `go build ./...` after every change
7. Read `.factory/library/architecture.md` before starting

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
