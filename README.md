# MockAPI 🦞

Lightweight local API mock server with Web UI. Zero dependencies. One binary.

## Features

- 🚀 **Zero dependencies** — single Go binary, no Node.js, no database
- 🎨 **Web UI** — visually manage mock routes with dark theme
- 🔀 **Dynamic routes** — path params (`/users/:id`) with `{{id}}` substitution
- 🌐 **Wildcard** — catch-all routes (`/api/*`)
- 📋 **Request logging** — real-time log with auto-refresh
- 📤 **Import/Export** — backup and share routes as JSON
- 📝 **Preset templates** — one-click common REST API routes
- ⏱️ **Simulated delay** — add latency to test loading states
- 🔀 **CORS** — enabled by default for frontend development
- 🔍 **Search & filter** — quickly find routes
- ✏️ **Edit routes** — update without delete + recreate
- 💾 **Auto-save** — routes persisted to `mockapi.json`

## Install

```bash
go build -o mockapi ./cmd/mockapi
```

## Usage

```bash
# Start on default port 8088
./mockapi

# Custom port
./mockapi 3000
./mockapi --port 3000

# Custom config file
./mockapi --config my-routes.json
```

Open http://localhost:8088/ for the Web UI.

## How It Works

1. Open Web UI → Add routes (method, path, status, body, delay)
2. Mock endpoints at `http://localhost:<port>/mock/<path>`
3. Routes auto-save to `mockapi.json`

### Dynamic Path Params

Create a route with `:param` in the path:
- Path: `/users/:id`
- Body: `{"id": {{id}}, "name": "Alice"}`

Request: `GET /mock/users/42` → `{"id": 42, "name": "Alice"}`

### Wildcard

- Path: `/api/*` matches any `/api/...` request
- Method: `ALL` matches any HTTP method

### Templates

Click 📋 Templates in the UI for one-click setup of:
- REST users CRUD
- Blog posts + comments
- Auth login
- Health check
- Slow response simulation
- Generic 404

## API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/_api/routes` | GET | List all routes |
| `/_api/routes` | POST | Create a route |
| `/_api/routes` | PUT | Update a route |
| `/_api/routes?id=xxx` | DELETE | Delete a route |
| `/_api/logs` | GET | Get request logs |
| `/_api/clear-logs` | POST | Clear logs |
| `/_api/templates` | GET | Get preset templates |
| `/_api/import` | POST | Import routes (JSON array) |
| `/_api/export` | GET | Download routes as JSON |

### Example: Add via curl

```bash
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
```

## License

MIT
