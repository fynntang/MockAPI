# MockAPI 🦞

**Lightweight Local API Mock Server**

Zero dependencies. One binary. Full-featured mock server with Web UI, JavaScript scripts, WebSocket, GraphQL, gRPC, and Swagger import.

[![Go Report Card](https://goreportcard.com/badge/github.com/fynntang/MockAPI)](https://goreportcard.com/report/github.com/fynntang/MockAPI)
[![License: MIT](https://img.shields.io/badge/License-MIT-purple.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
<a href="https://www.nxgntools.com/tools/mockapi?utm_source=mockapi" target="_blank" rel="noopener" style="display: inline-block; width: auto;">
    <img src="https://www.nxgntools.com/api/embed/mockapi?type=FEATURED_ON&hideUpvotes=true" alt="Launching Soon on NextGen Tools" style="height: 48px; width: auto;" />
</a>
<a href="https://www.producthunt.com/products/mockapi-4?embed=true&amp;utm_source=badge-featured&amp;utm_medium=badge&amp;utm_campaign=badge-mockapi-4" target="_blank" rel="noopener noreferrer"><img alt="MockAPI - Lightweight Local API Mock Server One binary Full-featured. | Product Hunt" width="250" height="54" src="https://api.producthunt.com/widgets/embed-image/v1/featured.svg?post_id=1095393&amp;theme=neutral&amp;t=1773409390030"></a>

🌐 **Live Demo**: [mockapi.work](https://mockapi.work)

---

## ✨ Features

- 🎨 **Web UI** - Dark theme interface to visually manage mock routes
- 📜 **JavaScript Scripts** - Dynamic responses with full JavaScript engine (goja)
- 📥 **Swagger Import** - Import OpenAPI 2.0/3.x and auto-generate routes
- 🔌 **WebSocket Mock** - Mock WebSocket endpoints with message handlers
- ⚡ **GraphQL Support** - Mock queries and mutations
- 🔁 **gRPC-Web** - gRPC-Web compatible HTTP endpoints
- 🔄 **Proxy Mode** - Forward unmatched requests to real backend
- ⚙️ **Conditional Routes** - Match by headers or body content
- 🔀 **Dynamic Paths** - Path parameters (`:id`) and wildcards (`*`)
- 🚀 **Zero Dependencies** - Single Go binary, no runtime dependencies
- 🔥 **Hot Reload** - Automatic config reload on file changes
- 🔒 **HTTPS/CORS** - Built-in HTTPS and configurable CORS
- 📦 **12+ Templates** - Pre-built templates for common patterns

---

## 🚀 Quick Start

### Install

```bash
# From Go (requires Go 1.21+)
go install github.com/fynntang/MockAPI@latest

# Or clone and build
git clone https://github.com/fynntang/MockAPI.git
cd MockAPI
go build -o mockapi ./cmd/mockapi
```

### Run

```bash
# Start server with default config
mockapi serve

# Start with custom port and hot reload
mockapi serve --port 8088 --hot-reload

# Initialize a new project
mockapi init my-project

# Validate configuration
mockapi validate
```

### Test

```bash
# Create a route via API
curl -X POST http://localhost:8088/_api/routes \
  -H "Content-Type: application/json" \
  -d '{
  "method": "GET",
  "path": "/users/:id",
  "status": 200,
  "body": "{\"id\": {{id}}, \"name\": \"Alice\"}"
}'

# Test the route
curl http://localhost:8088/mock/users/42
# Response: {"id": 42, "name": "Alice"}
```

---

## 📖 Documentation

### Configuration

MockAPI uses JSON configuration files:

```json
{
  "port": 8088,
  "routes": [
    {
      "method": "GET",
      "path": "/users/:id",
      "status": 200,
      "body": "{\"id\": {{id}}, \"name\": \"User\"}"
    }
  ]
}
```

### Dynamic Responses

Use JavaScript for complex logic:

```javascript
// Access: method, path, params, query, body, headers
{
  "id": parseInt(params.id),
  "timestamp": new Date().toISOString(),
  "userAgent": headers["user-agent"]
}
```

### Path Matching

| Pattern | Matches |
|---------|---------|
| `/users` | `/users` |
| `/users/:id` | `/users/1`, `/users/42` |
| `/files/*` | `/files/a`, `/files/a/b/c` |

---

## 🏗️ Project Structure

```
mockapi/
├── cmd/mockapi/          # CLI entry point
├── pkg/
│   ├── cli/             # CLI commands
│   ├── config/          # Configuration management
│   ├── graphql/         # GraphQL mock engine
│   ├── grpcmock/        # gRPC-Web compatibility
│   ├── script/          # JavaScript engine (goja)
│   ├── server/          # HTTP server
│   ├── swagger/         # OpenAPI import
│   └── ws/              # WebSocket support
├── website/             # Landing page & assets
├── routes/              # Example route configs
├── protos/              # Protocol buffer definitions
└── scripts/             # Example JavaScript scripts
```

---

## 📊 Performance

Optimized with route indexing for O(1) lookup:

| Metric | Before | After | Improvement |
|--------|--------|-------|-------------|
| Latency | 5169 ns/op | 214.8 ns/op | **24x faster** |
| Memory | 6400 B/op | 68 B/op | **94% less** |

```bash
# Run benchmarks
go test -bench=. ./pkg/...
```

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/config/...
```

**Coverage**: 26 unit tests across all core packages

---

## 🐳 Deployment

### Docker

```bash
# Build image
docker build -t mockapi .

# Run container
docker run -p 8088:8088 mockapi
```

### Cloudflare Pages

The landing page is deployed to Cloudflare Pages:

```bash
# Install Wrangler CLI
npm install -g wrangler

# Deploy
wrangler pages deploy website/
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Reporting Issues

- 🐛 **Bug Report**: [GitHub Issues](https://github.com/fynntang/MockAPI/issues/new?template=bug_report.md)
- 💡 **Feature Request**: [GitHub Issues](https://github.com/fynntang/MockAPI/issues/new?template=feature_request.md)
- ❓ **Question**: [GitHub Issues](https://github.com/fynntang/MockAPI/issues/new?template=question.md)

---

## 📄 License

MockAPI is released under the [MIT License](LICENSE).

```
Copyright (c) 2026 MockAPI

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

---

## 🙏 Acknowledgments

- [Go](https://go.dev/) - The Go Programming Language
- [goja](https://github.com/dop251/goja) - JavaScript interpreter in Go
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework

---

## 📬 Contact

- **Website**: [mockapi.work](https://mockapi.work)
- **GitHub**: [fynntang/MockAPI](https://github.com/fynntang/MockAPI)
- **Email**: contact@mockapi.work

---

**Built with ❤️ using Go**
