# MockAPI 🦞

Lightweight local API mock server with Web UI. No dependencies. One binary.

## Features

- 🚀 **Zero dependencies** — single Go binary, no Node.js needed
- 🎨 **Web UI** — visually manage mock routes
- ⚡ **Instant setup** — `mockapi` and go
- ⏱️ **Simulated delay** — add latency to test loading states
- 📋 **Copy curl** — one-click copy for testing
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

# Or
./mockapi --port 3000
```

Then open http://localhost:8088/ for the Web UI.

## How It Works

1. Open Web UI → Add routes (method, path, status, body, delay)
2. Your mock endpoints are available at `http://localhost:<port>/mock/<path>`
3. Routes auto-save to `mockapi.json`

### Example

Add a route:
- Method: `GET`
- Path: `/users/1`
- Status: `200`
- Body: `{"id": 1, "name": "Alice"}`

Then test it:
```bash
curl http://localhost:8088/mock/users/1
# {"id": 1, "name": "Alice"}
```

### API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/_api/routes` | GET | List all routes |
| `/_api/routes` | POST | Create a route |
| `/_api/routes?id=xxx` | DELETE | Delete a route |

### Example: Add via API

```bash
curl -X POST http://localhost:8088/_api/routes \
  -H "Content-Type: application/json" \
  -d '{
    "method": "GET",
    "path": "/users",
    "status": 200,
    "body": "[{\"id\":1},{\"id\":2}]",
    "delay_ms": 200
  }'
```

## License

MIT
