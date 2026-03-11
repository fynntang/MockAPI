# MockAPI 🦞

> Lightweight local API mock server with Web UI. Zero dependencies. One binary.

## ✨ Features

- 🚀 **Zero dependencies** — single Go binary, no Node.js, no database
- 🎨 **Web UI** — dark theme, visually manage routes
- 🔀 **Dynamic routes** — path params (`/users/:id`) with `{{id}}` substitution
- 🌐 **Wildcard** — catch-all routes (`/api/*`)
- 🔀 **ALL method** — one route matches any HTTP method
- ⚡ **Conditional responses** — match by request headers or body content
- 📜 **JavaScript scripts** — dynamic responses with full scripting power
- 🔄 **Proxy mode** — forward unmatched requests to a real backend
- 📥 **Swagger/OpenAPI import** — auto-generate mocks from API specs
- 🔌 **WebSocket mock** — mock WS endpoints with auto-reply or scripts
- 📋 **Request logging** — real-time log with auto-refresh
- 📤 **Import/Export** — backup and share routes as JSON
- 📝 **Preset templates** — one-click common REST API routes
- ⏱️ **Simulated delay** — add latency to test loading states
- 🌐 **CORS** — enabled by default
- 🔒 **HTTPS** — optional TLS support
- 🐳 **Docker** — ready-to-use Dockerfile
- 💾 **Auto-save** — routes persisted to `mockapi.json`

## Install

```bash
# Build from source
go build -o mockapi ./cmd/mockapi

# Or use Docker
docker build -t mockapi .
docker run -p 8088:8088 mockapi
```

## Usage

```bash
mockapi                       # Default port 8088
mockapi 3000                  # Custom port
mockapi --port 3000           # Custom port (flag style)
mockapi --config routes.json  # Custom config file
mockapi --proxy http://backend:3000  # Proxy unmatched to backend
mockapi --https --cert cert.pem --key key.pem  # Enable HTTPS
```

Open **http://localhost:8088/** for the Web UI.

## Features in Detail

### Dynamic Path Params

Create a route with `:param` in the path:
- Path: `/users/:id`
- Body: `{"id": {{id}}, "name": "Alice"}`

```
GET /mock/users/42 → {"id": 42, "name": "Alice"}
```

### JavaScript Scripts

Write JS code for dynamic responses. Available variables:
- `method` — HTTP method
- `path` — request path
- `headers` — request headers map
- `body` — request body string
- `params` — path parameters map
- `query` — query parameters map

```javascript
// Random response
var random = Math.floor(Math.random() * 1000);
respond({
  status: 200,
  body: JSON.stringify({ id: random, timestamp: Date.now() })
});
```

```javascript
// Use path params
respond({
  status: 200,
  body: JSON.stringify({ userId: params.id, query: query.page })
});
```

### Conditional Responses

Match based on request headers or body:
- **Match Headers**: `{"Authorization": "Bearer test"}`
- **Match Body**: `error` (substring match)

Multiple routes on the same path with different conditions = different responses!

### Swagger/OpenAPI Import

Import OpenAPI 2.0 or 3.x specs (JSON or YAML):
1. Click **📥 Swagger** in the UI
2. Paste spec or upload file
3. Routes auto-generated from paths

### WebSocket Mock

Create WebSocket handlers:
- **Auto Reply** — static JSON response to any message
- **On Connect** — message sent when client connects
- **On Message Script** — JS to process incoming messages

Connect at `ws://localhost:8088/ws/<path>`

### Proxy Mode

Set a proxy URL in Settings. Unmatched requests are forwarded to the real backend:
```
Frontend → MockAPI (mocked routes intercepted) → Real Backend
```

Perfect for frontend development against a partially-ready API.

## API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/_api/routes` | GET/POST/PUT/DELETE | CRUD for mock routes |
| `/_api/logs` | GET | Get request logs |
| `/_api/clear-logs` | POST | Clear all logs |
| `/_api/templates` | GET | Get preset templates |
| `/_api/import` | POST | Import routes (JSON array) |
| `/_api/import-swagger` | POST | Import from OpenAPI spec |
| `/_api/export` | GET | Download routes as JSON |
| `/_api/config` | GET/PUT | Server configuration |
| `/_api/ws` | GET/POST/DELETE | WebSocket handlers |

### Example: Add Script Route

```bash
curl -X POST http://localhost:8088/_api/routes \
  -H "Content-Type: application/json" \
  -d '{
    "method": "GET",
    "path": "/random",
    "script": "respond({status: 200, body: JSON.stringify({rand: Math.random()})})"
  }'
```

### Example: Import Swagger

```bash
curl -X POST http://localhost:8088/_api/import-swagger \
  -H "Content-Type: application/json" \
  -d @openapi.json
```

## Cross Compilation

```bash
make build-all
```

Builds for: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64

## Tech Stack

- **Go 1.21** — no runtime dependencies
- **goja** — JavaScript engine for dynamic responses
- **gorilla/websocket** — WebSocket support
- **yaml.v3** — OpenAPI/YAML parsing

## License

MIT