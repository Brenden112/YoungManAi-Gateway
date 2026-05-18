# Provider Policy

## Provider Types

The gateway classifies every channel into one of four provider types. The type is set automatically via `BeforeCreate` based on the channel's upstream API type, and can be overridden by an admin.

| Provider Type | Description | Default Risk | Visible To |
|---|---|---|---|
| `official_cloud` | Direct calls to official AI cloud APIs (OpenAI, Anthropic, Gemini, Azure, AWS Bedrock, etc.) | low | all users |
| `aggregator` | Third-party aggregator / reseller APIs (OpenRouter, SiliconFlow, etc.) | medium | all users |
| `authorized_proxy` | Self-hosted or on-premise models (Ollama, vLLM, custom OpenAI-compatible) | low | all users |
| `experimental_proxy` | Experimental or unverified proxy endpoints (e.g. KiroGateway) | high | internal users only |

## Access Control Rules

### Normal users
- Can call channels with provider type `official_cloud`, `aggregator`, or `authorized_proxy`.
- Requests routed to an `experimental_proxy` channel are rejected with **HTTP 403**.
- The `/v1/models` endpoint hides all models served exclusively by `experimental_proxy` channels.

### Internal users
- Identified by user group `internal` **or** admin role (`role >= 10`).
- Can call all four provider types.
- API tokens must have `allow_experimental = true` to reach `experimental_proxy` channels.

## experimental_proxy Safety Invariants

1. **Default disabled**: every new channel whose provider type resolves to `experimental_proxy` is created with `status = manually_disabled`. An admin must explicitly enable it.
2. **Double gate**: a request reaches an `experimental_proxy` channel only when **both** conditions are true:
   - The caller is an internal user (`group = internal` or `role >= 10`).
   - The API token has `allow_experimental = true`.
3. **Fallback isolation**: the channel selection algorithm never falls back to an `experimental_proxy` channel for a non-experimental request.
4. **Bulk disable**: admins can disable all enabled `experimental_proxy` channels in one call via `POST /api/channel/disable-experimental`.

## Risk Levels

| Level | Meaning |
|---|---|
| `low` | Production-grade, SLA-backed provider |
| `medium` | Generally reliable but no formal SLA |
| `high` | Experimental, unverified, or potentially unstable |

## Visibility / Available Scope

| Value | Meaning |
|---|---|
| `public` | Visible to all authenticated users |
| `internal` | Visible to internal users only |
| `private` | Visible to the owning org/project only |
