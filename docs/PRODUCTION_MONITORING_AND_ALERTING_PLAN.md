# Production Monitoring And Alerting Plan

Date: 2026-05-28
Phase: `Phase 7 production preparation sign-off pack`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`

This plan defines required production monitoring before release. It does not configure monitoring and does not make the release production ready.

## Required Alerts

| Signal | Alert condition | Required response |
|---|---|---|
| 5xx error rate | Sustained increase over release threshold | Triage route, rollback if product-impacting. |
| Provider failure rate | Provider/channel failures exceed baseline | Disable affected provider/channel or switch to approved same-type provider. |
| Token usage anomaly | Sudden spike, impossible totals, or model mismatch | Freeze affected key/provider and reconcile usage. |
| Balance deduction anomaly | Overcharge, undercharge, negative balance drift, or missing deduction | Stop affected traffic path and reconcile ledger. |
| Zero-balance bypass | Any zero/insufficient-balance request reaches upstream | Stop traffic path immediately and treat as release blocker. |
| Experimental normal-user call | Any normal user calls or sees `experimental_proxy` | Disable experimental channels and investigate policy bypass. |
| Prompt/response leakage | Full prompt or full response appears in usage/error logs | Stop leaking path, rotate affected secrets if needed, incident response. |
| Provider credential leakage | Provider key or credential pattern appears in logs/docs/artifacts | Disable provider, rotate credential, restrict logs. |
| Database connection errors | Elevated connection failures, migration errors, or slow queries | Failover/restore according to database runbook. |
| Redis errors | Redis unavailable or error rate affects auth/rate-limit/cache | Fail closed where required and restore Redis. |
| Queue/task errors | Background jobs or async tasks fail or backlog grows | Pause release, drain/retry jobs, inspect failed task class. |

## Dashboards

- API request volume, status codes, latency, and 5xx rate.
- Provider/channel request counts, failure rates, retries, and fallback decisions.
- Token usage by model, provider type, channel, user, organization, and project.
- Billing quota/cost deltas and balance changes.
- Zero-balance rejection count and upstream request delta.
- Experimental model visibility/call attempts by user type and token flag.
- Error-log sanitizer and redaction counters if available.
- Database connections, query latency, migration status, and disk capacity.
- Redis availability, latency, memory, and eviction/error counters.
- Worker/queue depth, task failures, and retry counts.

## Alert Ownership

- Release owner: unassigned manual gate.
- On-call owner: unassigned manual gate.
- Billing reconciliation owner: unassigned manual gate.
- Security incident owner: unassigned manual gate.
- Provider operations owner: unassigned manual gate.

## Privacy Constraints

Alerts and dashboards must not include full API keys, bearer tokens, provider credentials, prompts, or responses. Use IDs, prefixes, hashes, counts, status codes, and sanitized error categories only.
