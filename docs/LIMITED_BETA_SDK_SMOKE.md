# Limited Beta SDK Smoke

Date: 2026-05-25
Environment: Docker fixture network `new-api_fixture-network`
Provider: fake upstream only
SDK: OpenAI Node SDK, installed in a temporary container
Real provider: not used
Production readiness: `not_ready`

## Result

| Check | Status |
|---|---|
| SDK install in temporary container | Passed |
| `models.list` | Passed |
| `chat.completions.create` non-streaming | Passed |
| `chat.completions.create` streaming | Passed |
| Real OpenAI key used | No |
| Real upstream provider called | No |
| Full prompt/response written to docs | No |

## Sanitized Runtime Command

The executed smoke created a fixture-only user and API key in memory, then used:

```js
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: fixtureApiKey,
  baseURL: "http://new-api:3000/v1",
});

await client.models.list();
await client.chat.completions.create({
  model: "gpt-4o-mini",
  messages: [{ role: "user", content: "fixture sdk" }],
});

const stream = await client.chat.completions.create({
  model: "gpt-4o-mini",
  stream: true,
  messages: [{ role: "user", content: "fixture sdk stream" }],
});

for await (const chunk of stream) {
  break;
}
```

The real fixture API key was not printed or written to this document.
