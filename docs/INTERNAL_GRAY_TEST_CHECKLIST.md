# Internal Gray Test Checklist

Date: 2026-05-22

Current phase: `Phase 3 internal gray test planning`
Deployment readiness: `internal_gray_ready`
Production readiness: `not_ready`

Use this checklist during internal gray execution. Record pass/fail evidence without storing secrets, real prompts, real responses, or production data.

## Environment And Secret Safety

- [ ] Confirm tested commit SHA:
- [ ] Confirm staging environment identifier:
- [ ] Confirm `.env.staging` or equivalent secret source is untracked.
- [ ] Confirm no real high-privilege upstream API key is used.
- [ ] Confirm no real paid provider call is made unless a low-limit staging key has explicit human approval.
- [ ] Run `bash scripts/check-config-secrets.sh`.
- [ ] Confirm logs and evidence contain no API key, provider credential, bearer token, prompt, or response.

## Roles And Tenancy

- [ ] Create or verify admin user.
- [ ] Create or verify normal user.
- [ ] Create or verify internal user.
- [ ] Create or verify test organization.
- [ ] Create or verify test project.
- [ ] Bind test API key to user.
- [ ] Bind test API key to organization.
- [ ] Bind test API key to project.

## API Key Controls

- [ ] Create API key.
- [ ] Verify full plaintext API key is shown only once.
- [ ] Verify later API key views show only prefix or masked value.
- [ ] Verify stored key uses hash/prefix, not plaintext.
- [ ] Disable API key and verify relay call fails.
- [ ] Verify `allowed_models` blocks disallowed model.
- [ ] Verify `allowed_provider_types` blocks disallowed provider type.

## OpenAI-Compatible API

- [ ] Verify `GET /v1/models`.
- [ ] Verify `POST /v1/chat/completions`.
- [ ] Verify non-streaming chat completion.
- [ ] Verify streaming if the provider/channel supports stream.
- [ ] Verify OpenAI SDK compatibility against staging endpoint.

## Provider And Channel

- [ ] Verify fake provider smoke path.
- [ ] Verify `official_cloud` placeholder config.
- [ ] Verify `official_cloud` low-limit call only if explicitly approved.
- [ ] Verify `aggregator` can be configured without real paid call.
- [ ] Verify `ProviderAccount` credential is encrypted at rest.
- [ ] Verify `ProviderAccount` credential is available for runtime use only.
- [ ] Disable provider and verify it cannot route.
- [ ] Disable channel and verify it cannot route.
- [ ] Verify legacy channel compatibility.

## Experimental Proxy Isolation

- [ ] Verify `experimental_proxy` default disabled.
- [ ] Verify `experimental_proxy` default internal-only.
- [ ] Verify normal user cannot see experimental models.
- [ ] Verify normal user cannot call experimental models.
- [ ] Verify internal user without `allow_experimental=true` cannot call experimental models.
- [ ] Verify internal user with `allow_experimental=true` can call only enabled approved experimental test path.
- [ ] Verify disabled `experimental_proxy` cannot be called by any user.
- [ ] Verify `official_cloud` failure does not fall back to `experimental_proxy`.
- [ ] Verify fallback/retry/preferred/specific/legacy paths cannot bypass restrictions.

## Billing And Balance

- [ ] Verify successful request records token counts.
- [ ] Verify successful request calculates cost.
- [ ] Verify successful request deducts balance.
- [ ] Verify failed request does not deduct incorrectly.
- [ ] Verify zero-balance request does not call upstream.
- [ ] Verify admin manual top-up works.
- [ ] Verify top-up creates operation record.

## Logs And Privacy

- [ ] Verify successful request writes `usage_log`.
- [ ] Verify failed request writes expected log or error record.
- [ ] Verify log contains `user_id`, `org_id`, `project_id`, and `api_key_id`.
- [ ] Verify log contains `provider_type`, `channel_id`, `provider_account_id`, `model`, tokens, cost, and status.
- [ ] Verify full prompt is not stored by default.
- [ ] Verify full response is not stored by default.
- [ ] Verify `params.Other` is redacted.
- [ ] Verify `error_message` is redacted.
- [ ] Verify logs do not leak API keys, provider credentials, or bearer tokens.

## Admin Dashboard

- [ ] Admin can view providers.
- [ ] Admin can view channels.
- [ ] Admin can view API keys.
- [ ] Admin can disable API keys.
- [ ] Admin can view usage logs.
- [ ] Admin can view user balances.
- [ ] Admin can manually top up balance.
- [ ] Normal user cannot access admin pages.
- [ ] Normal user cannot access admin APIs.

## Regression And CI Evidence

- [ ] Run `bash scripts/ci-verify.sh`.
- [ ] Run `LOCAL_FIXTURE=1 bash scripts/regression.sh`.
- [ ] Run Docker compose fixture smoke.
- [ ] Review GitHub Actions pre-release verification.
- [ ] Attach staging verification evidence.
- [ ] Record frontend local script blocker only as non-blocking if GitHub Actions `frontend-check` remains passed.

## Exit Review

- [ ] Review `docs/INTERNAL_GRAY_EXIT_CRITERIA.md`.
- [ ] Review `docs/INTERNAL_GRAY_RISK_REGISTER.md`.
- [ ] Confirm critical/high issues are zero.
- [ ] Confirm rollback path is documented.
- [ ] Confirm `production_readiness` remains `not_ready`.
