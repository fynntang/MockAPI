# MockAPI 🦞

> Lightweight local API mock server with Web UI. Zero dependencies. One binary.

## ✨ Features

- 🚀 **Zero dependencies** — single Go binary, no Node.js, no database
- 🎨 **Web UI** — dark theme, visually manage routes
- 🔀 **Dynamic routes** — path params (`/users/:id`) with `{{id}}` substitution
- 🌐 **Wildcard** — catch-all routes (`/api/*`)
- 🔀 **ALL method** — one route matches any HTTP method
- ⚡ **Conditional responses** — match by request headers or body content
- 🔄 **Proxy mode** — forward unmatched requests to a real backend
- 📋 **Request logging** — real-time log with auto-refresh and proxy indicator
- 📤 **Import/Export** — backup and share routes as JSON
- 📝 **Preset templates** — one-click common REST API routes
- ⏱️ **Simulated delay** — add latency to test loading states
- 🌐 **CORS** — enabled by default
- 🔍 **Search & filter** — quickly find routes
- ✏️ **Edit routes** — update without delete + recreate
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

## How It Works

### Basic Routes

1. Open Web UI → Add routes (method, path, status, body, delay)
2. Mock endpoints at `http://localhost:<port>/mock/<path>`
3. Routes auto-save to `mockapi.json`

### Dynamic Path Params

Create a route with `:param` in the path:
- Path: `/users/:id`
- Body: `{"id": {{id}}, "name": "Alice"}`

```
GET /mock/users/42 → {"id": 42, "name": "Alice"}
```

### Wildcard & ALL

- `/api/*` matches any `/api/...` request
- Method `ALL` matches any HTTP method

### Conditional Responses

Match based on request headers or body:
- **Match Headers**: `{"Authorization": "Bearer test"}` — only responds when this header is present
- **Match Body**: `error` — only responds when request body contains "error"

Multiple routes on the same path with different conditions = different responses!

### Proxy Mode

Set a proxy URL in Settings. Unmatched requests are forwarded to the real backend:
```
Frontend → MockAPI (mocked routes intercepted) → Real Backend (everything else)
```

Perfect for frontend development against a partially-ready API.

### Templates

Click 📋 Templates for one-click setup of:
- REST users CRUD
- Blog posts + comments
- Auth login (success + conditional fail)
- Health check
- Slow response simulation
- Generic 404
- Paginated list

## API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/_api/routes` | GET | List all routes |
| `/_api/routes` | POST | Create a route |
| `/_api/routes` | PUT | Update a route |
| `/_api/routes?id=xxx` | DELETE | Delete a route |
| `/_api/logs` | GET | Get request logs |
| `/_api/clear-logs` | POST | Clear all logs |
| `/_api/templates` | GET | Get preset templates |
| `/_api/import` | POST | Import routes (JSON array) |
| `/_api/export` | GET | Download routes as JSON |
| `/_api/config` | GET | Get current config |
| `/_api/config` | PUT | Update config (proxy, CORS, etc.) |

### Example

```bash
# Add a route
curl -X POST http://localhost:8088/_api/routes \
  -H "Content-Type: application/json" \
  -d '{
    "method": "GET",
    "path": "/users/:id",
    "status": 200,
    "body": "{\"id\": {{id}}, \"name\": \"Alice\"}",
    "delay_ms": 200,
    "description": "Get user by ID"
  }'

# Conditional route (only match when header present)
curl -X POST http://localhost:8088/_api/routes \
  -H "Content-Type: application/json" \
  -d '{
    "method": "POST",
    "path": "/auth/login",
    "status": 200,
    "body": "{\"token\": \"secret\"}",
    "match_headers": {"X-Auth-Key": "test"}
  }'

# Set proxy
curl -X PUT http://localhost:8088/_api/config \
  -H "Content-Type: application/json" \
  -d '{"proxy_url": "http://localhost:3000"}'
```

## Cross Compilation

```bash
make build-all
```

Builds for: linux-amd64, linux-arm64, darwin-amd64, darwin-arm64, windows-amd64

## License

MIT
