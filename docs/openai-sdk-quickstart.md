# OpenAI SDK Quickstart

The gateway exposes a fully OpenAI-compatible API. Any SDK or tool that supports a custom `base_url` works without modification.

## Base URL

```
http://localhost:3000/v1
```

Replace with your deployed host in production.

## Python (openai >= 1.0)

```python
from openai import OpenAI

client = OpenAI(
    api_key="sk-your-token-here",
    base_url="http://localhost:3000/v1",
)

# Chat completion
response = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "Hello!"}],
)
print(response.choices[0].message.content)

# Streaming
stream = client.chat.completions.create(
    model="gpt-4o",
    messages=[{"role": "user", "content": "Count to 5."}],
    stream=True,
)
for chunk in stream:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="", flush=True)
```

## Node.js (openai >= 4.0)

```javascript
import OpenAI from "openai";

const client = new OpenAI({
  apiKey: "sk-your-token-here",
  baseURL: "http://localhost:3000/v1",
});

// Chat completion
const response = await client.chat.completions.create({
  model: "gpt-4o",
  messages: [{ role: "user", content: "Hello!" }],
});
console.log(response.choices[0].message.content);

// Streaming
const stream = await client.chat.completions.create({
  model: "gpt-4o",
  messages: [{ role: "user", content: "Count to 5." }],
  stream: true,
});
for await (const chunk of stream) {
  process.stdout.write(chunk.choices[0]?.delta?.content ?? "");
}
```

## List available models

```bash
curl http://localhost:3000/v1/models \
  -H "Authorization: Bearer sk-your-token-here"
```

Internal users see all models including `experimental_proxy` channels. Normal users see only `official_cloud`, `aggregator`, and `authorized_proxy` models.

## API token flags

| Flag | Effect |
|---|---|
| `allow_experimental=false` (default) | Token cannot reach `experimental_proxy` channels |
| `allow_experimental=true` | Token can reach `experimental_proxy` channels (internal users only) |

## Error codes

| HTTP | Meaning |
|---|---|
| 401 | Invalid or expired API token |
| 402 | Insufficient quota — top up via admin panel |
| 403 | Access denied — `experimental_proxy` channel requires internal user + `allow_experimental` token |
| 429 | Rate limit exceeded |
| 503 | No available upstream channel for the requested model |
