# Limited Beta Checklist

Date: 2026-05-25

Current phase: `Phase 5 limited beta planning`
Deployment readiness: `limited_beta_ready`
Production readiness: `not_ready`

Use this checklist for small-scope limited beta execution. Evidence must be sanitized. Do not paste secrets, tokens, prompts, responses, or credentials into this file.

## Pre-Run Gate

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-001 | Confirm release owner and beta operator. | Pending |  |
| LBC-002 | Confirm environment is staging or beta-only, not production. | Pending |  |
| LBC-003 | Confirm fake provider is available. | Pending |  |
| LBC-004 | Confirm no high-privilege real provider key is configured. | Pending |  |
| LBC-005 | Confirm any real `official_cloud` test key is low-limit and manually approved before use. | Pending |  |
| LBC-006 | Confirm `experimental_proxy` is disabled and internal-only by default. | Pending |  |
| LBC-007 | Confirm `production_readiness` remains `not_ready`. | Pending |  |

## User, Organization, Project, And API Key

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-010 | Create or verify beta admin user. | Pending |  |
| LBC-011 | Create or verify beta normal user. | Pending |  |
| LBC-012 | Create or verify beta internal user. | Pending |  |
| LBC-013 | Create staging-only organization. | Pending |  |
| LBC-014 | Create staging-only project. | Pending |  |
| LBC-015 | Create API key bound to user, organization, and project as required. | Pending |  |
| LBC-016 | Verify API key plaintext is shown only at creation and is not returned afterward. | Pending |  |
| LBC-017 | Verify API key hash and prefix storage only. | Pending |  |
| LBC-018 | Verify disabled API key cannot call relay APIs. | Pending |  |
| LBC-019 | Verify organization/project token binding is enforced. | Pending |  |

## OpenAI-Compatible API

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-020 | Call `GET /v1/models` with normal beta API key. | Pending |  |
| LBC-021 | Call `POST /v1/chat/completions` with fake provider. | Pending |  |
| LBC-022 | If manually approved, call `POST /v1/chat/completions` with low-limit `official_cloud` test key. | Pending |  |
| LBC-023 | Verify response status and sanitized usage are recorded. | Pending |  |
| LBC-024 | Verify disabled model cannot be called. | Pending |  |
| LBC-025 | Verify `allowed_provider_types` rejects disallowed provider types. | Pending |  |

## Usage, Billing, And Balance

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-030 | Verify successful request creates `usage_log`. | Pending |  |
| LBC-031 | Verify failed request attribution is logged as designed. | Pending |  |
| LBC-032 | Verify token counts are present and reasonable. | Pending |  |
| LBC-033 | Verify cost calculation matches configured pricing. | Pending |  |
| LBC-034 | Verify balance is deducted for successful requests. | Pending |  |
| LBC-035 | Verify zero or insufficient balance is rejected before upstream call. | Pending |  |
| LBC-036 | Verify admin top-up works. | Pending |  |
| LBC-037 | Verify top-up is reflected in balance and operation history. | Pending |  |
| LBC-038 | Reconcile approved real-provider upstream fee against system deduction. | Pending |  |

## Provider, Channel, And Experimental Isolation

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-040 | Verify fake provider remains callable. | Pending |  |
| LBC-041 | Verify disabled provider cannot route. | Pending |  |
| LBC-042 | Verify disabled channel cannot route. | Pending |  |
| LBC-043 | Verify normal user cannot see `experimental_proxy` models. | Pending |  |
| LBC-044 | Verify normal user cannot call `experimental_proxy`. | Pending |  |
| LBC-045 | Verify internal user can test `experimental_proxy` only when explicitly enabled. | Pending |  |
| LBC-046 | Verify disabled `experimental_proxy` cannot route for any user. | Pending |  |
| LBC-047 | Verify `official_cloud` does not fallback to `experimental_proxy`. | Pending |  |

## Logs And Privacy

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-050 | Verify `usage_log` does not store full prompt by default. | Pending |  |
| LBC-051 | Verify `usage_log` does not store full response by default. | Pending |  |
| LBC-052 | Verify error logs do not contain full prompt or response. | Pending |  |
| LBC-053 | Verify `params.Other` is redacted. | Pending |  |
| LBC-054 | Verify logs do not contain API keys, bearer tokens, or provider credentials. | Pending |  |
| LBC-055 | Verify screenshots and issue evidence do not contain secrets or prompt/response content. | Pending |  |

## Admin Dashboard

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-060 | Verify admin provider page. | Pending |  |
| LBC-061 | Verify admin channel page. | Pending |  |
| LBC-062 | Verify admin API key page. | Pending |  |
| LBC-063 | Verify admin usage log page. | Pending |  |
| LBC-064 | Verify admin balance page. | Pending |  |
| LBC-065 | Verify admin top-up flow. | Pending |  |
| LBC-066 | Verify normal user cannot access admin pages or admin APIs. | Pending |  |

## Daily Observation

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-070 | Review daily `usage_log` attribution and cost. | Pending |  |
| LBC-071 | Review daily error logs and redaction. | Pending |  |
| LBC-072 | Review daily balance anomalies. | Pending |  |
| LBC-073 | Review provider failure rate. | Pending |  |
| LBC-074 | Review fallback behavior. | Pending |  |
| LBC-075 | Review token/cost consistency. | Pending |  |
| LBC-076 | Review upstream real fee versus system deduction for approved real-provider calls. | Pending |  |
| LBC-077 | Review for prompt/response leakage. | Pending |  |

## Exit Record

| ID | Check | Status | Evidence |
|---|---|---|---|
| LBC-090 | Confirm beta cycle completed or stopped with reason. | Pending |  |
| LBC-091 | Confirm critical/high issue count. | Pending |  |
| LBC-092 | Confirm medium issue disposition. | Pending |  |
| LBC-093 | Confirm rollback plan was verified. | Pending |  |
| LBC-094 | Confirm release-owner sign-off or stop decision. | Pending |  |
