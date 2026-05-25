# Limited Beta Signoff

Date: 2026-05-25

Phase: `limited-beta-notes-resolution`
Status: `completed_with_manual_gate`
Environment: local Docker fixture and current CI evidence
Commit: `57ad3623`
Deployment readiness: `production_preparation_ready_with_manual_provider_gate`
Production readiness: `not_ready`
Next recommended action: `decide whether to run low-limit real provider beta or prepare production sign-off pack`

## Decision

Fake-provider limited beta notes are closed. Production preparation may start only after the release owner decides whether to run a low-limit real-provider beta or keep the real-provider path behind a manual gate.

## Evidence Summary

- Config secret scan passed.
- Docker fixture build and runtime smoke passed with fake upstream only.
- GitHub Actions `Pre-release verification` run 23 passed for commit `675091d0`.
- Admin, normal, internal, and zero-balance fixture users were exercised.
- API key creation, masking after creation, disabled-key rejection, and normal-user admin rejection passed.
- Fake-provider `GET /v1/models` and non-streaming `POST /v1/chat/completions` passed.
- Normal user could not see or call experimental model.
- Disabled experimental channel blocked internal experimental request.
- Zero-balance request returned HTTP 402 without increasing fake upstream request count.
- Usage log was written and did not store full prompt/response content.
- No real provider key was used.
- No real upstream provider was called.
- No API key, bearer token, provider credential, prompt, or response was written into documentation.
- Organization/project runtime binding passed with fixture-only tenant records.
- Streaming smoke passed against the fake-provider fixture.
- OpenAI Node SDK smoke passed for model listing, non-streaming chat, and streaming chat.
- Local Go absence is closed by CI/Codespaces evidence; GitHub Actions `Pre-release verification` run 24 passed for commit `57ad3623`.

## Findings

| Severity | Count |
|---|---:|
| Critical | 0 |
| High | 0 |
| Medium | 0 |
| Low | 0 |

Open notes are tracked in `docs/LIMITED_BETA_ISSUES.md`.

## Signoff Status

| Signoff item | Status |
|---|---|
| Limited beta fake-provider core path | Passed |
| Critical/high stop criteria | Passed: none observed |
| Real low-limit `official_cloud` provider | Manual gate; requires release-owner approval and low-limit key |
| Streaming compatibility | Passed with fake-provider fixture |
| OpenAI SDK compatibility | Passed with fake-provider fixture |
| Organization/project runtime workflow | Passed with fixture-only tenant records |
| Fake-provider beta exit | Passed |
| Production preparation | Allowed only after release-owner decision on manual provider gate |

## Manual Provider Gate

| Item | Status |
|---|---|
| Low-limit `official_cloud` key provided | No |
| Real-provider limited beta approved | No |
| Daily real-provider quota cap | Not set |
| Real-provider rollback conditions | Use `docs/LIMITED_BETA_ROLLBACK_PLAN.md`; stop immediately on credential leakage, billing mismatch, privacy leak, zero-balance upstream call, or provider-type fallback bypass. |

## Release Owner Decision Needed

Before production preparation, a release owner must decide whether to:

- run a low-limit real-provider beta with a manually supplied test key and quota cap, or
- keep real-provider execution out of scope and proceed with fake-provider beta evidence only.

`production_readiness` must remain `not_ready`.
