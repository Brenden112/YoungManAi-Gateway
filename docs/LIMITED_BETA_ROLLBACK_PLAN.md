# Limited Beta Rollback Plan

Date: 2026-05-25

Current phase: `Phase 5 limited beta planning`
Deployment readiness: `limited_beta_ready`
Production readiness: `not_ready`

This rollback plan applies only to limited beta. It does not authorize production deployment or production rollback procedures.

## Rollback Owners

| Role | Responsibility |
|---|---|
| Release owner | Approves beta start, stop, rollback, and production preparation decision. |
| Beta operator | Executes provider/channel/API key disablement and collects sanitized evidence. |
| Admin reviewer | Verifies balances, logs, usage attribution, and privacy after rollback. |

## Rollback Triggers

Rollback is required if any stop condition in `docs/LIMITED_BETA_EXIT_CRITERIA.md` occurs, including:

- Key, token, credential, prompt, or response leakage.
- Normal user access to `experimental_proxy`.
- Zero-balance upstream call.
- Incorrect billing or balance deduction.
- Admin permission bypass.
- Organization or project authorization bypass.
- Fallback bypass of provider-type restrictions.
- Unapproved real-provider use.

## Immediate Containment

1. Disable affected API keys.
2. Disable affected provider accounts.
3. Disable affected channels.
4. Disable affected models if the issue is model-specific.
5. Keep `experimental_proxy` disabled and internal-only.
6. Remove low-limit real provider key from runtime configuration if it was involved.
7. Stop beta traffic and notify release owner.
8. Preserve sanitized logs and issue metadata without secrets, prompts, or responses.

## Service Rollback

Use the last known passed beta/staging artifact or commit approved by the release owner.

1. Record current commit SHA, environment, and rollback reason.
2. Stop beta traffic.
3. Deploy the approved previous artifact or revert to the approved previous commit through the normal deployment path.
4. Clear or refresh route/provider/channel caches if required by the deployment path.
5. Verify `/api/status`.
6. Verify `GET /v1/models` with fake provider.
7. Verify `POST /v1/chat/completions` with fake provider.
8. Verify disabled providers/channels/models cannot route.
9. Verify normal user cannot access `experimental_proxy`.
10. Verify zero-balance rejection happens before upstream call.

## Data And Billing Correction

- Do not delete usage logs unless the release owner approves a privacy-driven purge.
- Mark affected usage records for review if billing or attribution is wrong.
- Apply manual balance correction only after admin review.
- Record sanitized correction reason and operation ID.
- Reconcile any approved real-provider upstream fee against corrected system balance.

## Secret Handling

- Rotate any exposed API key, provider credential, or bearer token immediately.
- Remove exposed credentials from screenshots, issue bodies, documents, logs, and artifacts where feasible.
- Treat any full prompt or full response exposure as a privacy incident and stop beta until reviewed.
- Never paste replacement credentials into documents or issue comments.

## Rollback Verification

Rollback is complete only when:

- Affected API keys, providers, channels, or models are disabled or reverted.
- Fake-provider smoke passes.
- Real-provider traffic is stopped unless release owner explicitly re-approves low-limit testing.
- No normal user can access `experimental_proxy`.
- Zero-balance requests do not call upstream.
- Logs remain redacted.
- Balance and usage corrections are recorded.
- Release owner records a continue, extend, or stop decision.
