import http from 'node:http'

const port = Number(process.env.FAKE_PROVIDER_PORT || 4010)
let requestCount = 0

function send(res, status, body) {
  res.writeHead(status, {
    'content-type': 'application/json',
    'x-fixture-provider': 'fake-openai',
  })
  res.end(JSON.stringify(body))
}

const server = http.createServer(async (req, res) => {
  if (req.url === '/health') {
    send(res, 200, { ok: true })
    return
  }

  if (req.url === '/__fixture/requests') {
    send(res, 200, { request_count: requestCount })
    return
  }

  if (req.url === '/v1/models') {
    requestCount += 1
    send(res, 200, {
      object: 'list',
      data: [{ id: 'gpt-4o-mini', object: 'model', owned_by: 'fixture' }],
    })
    return
  }

  if (req.url === '/v1/chat/completions' && req.method === 'POST') {
    requestCount += 1
    req.resume()
    send(res, 200, {
      id: 'chatcmpl-fixture',
      object: 'chat.completion',
      created: Math.floor(Date.now() / 1000),
      model: 'gpt-4o-mini',
      choices: [
        {
          index: 0,
          message: { role: 'assistant', content: 'fixture response' },
          finish_reason: 'stop',
        },
      ],
      usage: { prompt_tokens: 1, completion_tokens: 1, total_tokens: 2 },
    })
    return
  }

  send(res, 404, { error: { message: 'fixture route not found' } })
})

server.listen(port, '0.0.0.0', () => {
  console.log(`fake OpenAI-compatible provider listening on ${port}`)
})
