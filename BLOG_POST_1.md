# 我用 Go 写了一个零依赖的 API Mock 服务器，性能提升 24 倍

> **摘要**: 从零开始开发一个轻量级本地 API Mock 服务器，使用纯 Go 实现，零运行时依赖，单二进制部署。通过路由索引优化，性能提升 24 倍，内存减少 94%。

---

## 背景

作为一名全栈开发者，我在日常开发中经常需要 Mock API 接口。市面上的 Mock 工具要么太重（需要 Node.js 运行时），要么功能太简单（不支持动态响应）。我决定自己写一个：

**需求清单**：
- ✅ 零运行时依赖（单二进制部署）
- ✅ 支持动态路由和路径参数
- ✅ 支持 JavaScript 脚本引擎
- ✅ 支持 WebSocket、GraphQL、gRPC
- ✅ 支持 Swagger/OpenAPI 导入
- ✅ 内置 Web UI（暗色主题）
- ✅ 热重载配置
- ✅ 高性能路由匹配

**技术选型**：Go 1.21

理由：
1. 编译为单二进制，无运行时依赖
2. 并发性能好
3. 标准库强大（net/http, encoding/json 等）
4. 部署简单（scp 上传即可运行）

---

## 架构设计

### 整体架构

```
┌─────────────────────────────────────────────────┐
│                   MockAPI                        │
├─────────────────────────────────────────────────┤
│  CLI  │  Web UI  │  Config  │  Route Engine    │
├─────────────────────────────────────────────────┤
│              Core Packages                       │
│  ┌─────────┬─────────┬─────────┬─────────────┐ │
│  │ REST    │ GraphQL │ WebSocket│ gRPC-Web    │ │
│  │ Mock    │ Mock    │ Mock    │ Mock        │ │
│  └─────────┴─────────┴─────────┴─────────────┘ │
│  ┌─────────┬─────────┬─────────┬─────────────┐ │
│  │ Script  │ Swagger │ Config  │ Route       │ │
│  │ Engine  │ Import  │ Manager │ Index       │ │
│  └─────────┴─────────┴─────────┴─────────────┘ │
└─────────────────────────────────────────────────┘
```

### 核心模块

1. **pkg/config** - 配置管理（JSON/YAML 加载、热重载）
2. **pkg/router** - 路由引擎（精确匹配、参数匹配、通配符）
3. **pkg/script** - JavaScript 脚本引擎（goja）
4. **pkg/graphql** - GraphQL Mock 支持
5. **pkg/grpcmock** - gRPC-Web 兼容端点
6. **pkg/swagger** - Swagger/OpenAPI 导入
7. **pkg/websocket** - WebSocket Mock 支持

---

## 性能优化：从 O(n) 到 O(1)

### 问题

最初的路由匹配是线性搜索：

```go
func (s *Server) matchRoute(method, path string) *Route {
    for _, route := range s.routes {
        if route.Method == method && route.Match(path) {
            return route
        }
    }
    return nil
}
```

**性能测试**：
```
BenchmarkRouteMatch-8    231218    5169 ns/op    6400 B/op
```

### 优化方案：RouteIndex

我设计了一个两层索引结构：

```go
type RouteIndex struct {
    // 精确匹配：map[method:path] -> Route
    exact map[string]*Route
    
    // 参数路由：前缀索引 map[method:prefix] -> []Route
    param map[string][]*Route
    
    // 通配符路由：按方法分组
    wildcard map[string][]*Route
}
```

**匹配逻辑**：
1. 先查 `exact` map（O(1)）
2. 再查 `param` 前缀索引（O(k)，k 为候选路由数）
3. 最后查 `wildcard`（O(m)，m 为通配符路由数）

### 优化结果

```go
func (idx *RouteIndex) Match(method, path string) *Route {
    // 1. 精确匹配 O(1)
    key := method + ":" + path
    if route, ok := idx.exact[key]; ok {
        return route
    }
    
    // 2. 参数路由前缀匹配
    segments := strings.Split(path, "/")
    for i := len(segments); i > 0; i-- {
        prefix := strings.Join(segments[:i], "/")
        key := method + ":" + prefix
        if candidates, ok := idx.param[key]; ok {
            for _, route := range candidates {
                if route.Match(path) {
                    return route
                }
            }
        }
    }
    
    // 3. 通配符匹配
    if candidates, ok := idx.wildcard[method]; ok {
        for _, route := range candidates {
            if route.Match(path) {
                return route
            }
        }
    }
    
    return nil
}
```

**性能对比**：

| 指标 | 优化前 | 优化后 | 提升 |
|------|--------|--------|------|
| 单次操作耗时 | 5169 ns | 214.8 ns | **24x** |
| 内存分配 | 6400 B/op | 68 B/op | **94%↓** |
| 时间复杂度 | O(n) | O(1) | - |

```
BenchmarkRouteMatch-8    5589621    214.8 ns/op    68 B/op
```

---

## 核心功能实现

### 1. 动态路由匹配

支持 `:param` 和 `*` 通配符：

```go
func (r *Route) Match(path string) bool {
    routeSegments := strings.Split(r.Path, "/")
    pathSegments := strings.Split(path, "/")
    
    if len(routeSegments) != len(pathSegments) && !r.HasWildcard {
        return false
    }
    
    for i, routeSeg := range routeSegments {
        if routeSeg == "*" {
            return true // 通配符匹配剩余所有
        }
        if strings.HasPrefix(routeSeg, ":") {
            continue // 参数匹配任意值
        }
        if routeSeg != pathSegments[i] {
            return false
        }
    }
    
    return true
}
```

### 2. JavaScript 脚本引擎

使用 [goja](https://github.com/dop251/goja) 实现：

```go
func (s *ScriptEngine) Execute(script string, ctx *ScriptContext) (interface{}, error) {
    vm := goja.New()
    
    // 注入上下文
    vm.Set("method", ctx.Method)
    vm.Set("path", ctx.Path)
    vm.Set("params", ctx.Params)
    vm.Set("query", ctx.Query)
    vm.Set("body", ctx.Body)
    vm.Set("headers", ctx.Headers)
    
    // 执行脚本
    result, err := vm.RunString(script)
    if err != nil {
        return nil, err
    }
    
    return result.Export(), nil
}
```

**使用示例**：

```javascript
// 动态响应
{
  "id": parseInt(params.id),
  "name": "User " + params.id,
  "timestamp": new Date().toISOString()
}

// 条件响应
if (headers["x-api-key"] === "secret") {
  return { authorized: true }
}
return { authorized: false }
```

### 3. Swagger/OpenAPI 导入

解析 OpenAPI 规范，自动生成 Mock 路由：

```go
func (i *SwaggerImporter) Import(spec *openapi3.T) ([]*Route, error) {
    var routes []*Route
    
    for path, pathItem := range spec.Paths {
        for method, operation := range pathItem.Operations() {
            route := &Route{
                Method: strings.ToUpper(method),
                Path:   convertPath(path),
            }
            
            // 生成示例响应
            if operation.Responses != nil {
                response := operation.Responses.Value("200")
                if response != nil {
                    route.Body = generateMockBody(response.Value)
                }
            }
            
            routes = append(routes, route)
        }
    }
    
    return routes, nil
}
```

### 4. GraphQL Mock 支持

解析 GraphQL Schema，支持 Query/Mutation：

```go
func (g *GraphQLMock) HandleOperation(operation string, variables map[string]interface{}) interface{} {
    // 解析操作名
    parsed := graphql.Parse(operation)
    
    // 匹配预定义的 Mock 响应
    for _, mock := range g.mocks {
        if mock.OperationName == parsed.Operation.Name {
            return mock.GenerateResponse(variables)
        }
    }
    
    // 默认响应
    return g.GenerateDefault(parsed)
}
```

---

## CLI 设计

提供完整的命令行工具：

```bash
# 初始化新项目
mockapi init my-project

# 验证配置
mockapi validate

# 启动服务
mockapi serve --port 8088 --hot-reload

# 查看版本
mockapi version

# 帮助
mockapi help
```

**实现**（使用 cobra）：

```go
func main() {
    var rootCmd = &cobra.Command{
        Use:   "mockapi",
        Short: "Lightweight Local API Mock Server",
    }
    
    var serveCmd = &cobra.Command{
        Use:   "serve",
        Short: "Start mock server",
        Run: func(cmd *cobra.Command, args []string) {
            server := NewServer(config)
            server.Start()
        },
    }
    serveCmd.Flags().IntP("port", "p", 8088, "Server port")
    serveCmd.Flags().Bool("hot-reload", false, "Enable hot reload")
    
    rootCmd.AddCommand(serveCmd)
    rootCmd.Execute()
}
```

---

## 单元测试

26 个测试用例，覆盖所有核心模块：

```go
// pkg/config/config_test.go
func TestLoadConfig(t *testing.T) {
    config, err := LoadConfig("testdata/valid.json")
    assert.NoError(t, err)
    assert.Equal(t, 8088, config.Port)
}

// pkg/router/index_test.go
func TestRouteIndex_Match(t *testing.T) {
    idx := NewRouteIndex()
    idx.Add(&Route{Method: "GET", Path: "/users/:id"})
    
    route := idx.Match("GET", "/users/42")
    assert.NotNil(t, route)
    assert.Equal(t, "42", route.ExtractParam("id", "/users/42"))
}

// pkg/swagger/import_test.go
func TestSwaggerImporter(t *testing.T) {
    spec, _ := openapi3.NewLoader().LoadFromFile("testdata/petstore.yaml")
    routes, err := Import(spec)
    assert.NoError(t, err)
    assert.Greater(t, len(routes), 0)
}
```

---

## 部署方案

### 1. 本地运行
```bash
go install github.com/fynntang/MockAPI@latest
mockapi serve
```

### 2. Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o mockapi ./cmd/mockapi

FROM alpine:latest
COPY --from=builder /app/mockapi /mockapi
EXPOSE 8088
CMD ["/mockapi"]
```

### 3. GitHub Actions CI/CD
```yaml
name: Release
on:
  push:
    tags: ['v*']
jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go build -o mockapi ./cmd/mockapi
      - uses: softprops/action-gh-release@v1
        with:
          files: mockapi
```

---

## 项目地址

- **GitHub**: https://github.com/fynntang/MockAPI
- **落地页**: https://mockapi.work
- **文档**: https://github.com/fynntang/MockAPI#readme

---

## 下一步计划

1. **Product Hunt 发布** - 获取早期用户反馈
2. **Pro 版本规划** - 团队协作、云同步
3. **生态系统** - VS Code 插件、CLI 增强
4. **社区建设** - Discord 社群、贡献者计划

---

## 总结

开发 MockAPI 的过程中，我深刻体会到：

1. **性能优化要趁早** - 路由索引让性能提升 24 倍
2. **零依赖是优势** - 单二进制部署极大简化运维
3. **测试驱动开发** - 26 个测试用例保证质量
4. **用户体验重要** - CLI 和 Web UI 提升易用性

如果你也需要一个轻量级 Mock 服务器，欢迎试用 MockAPI！有任何问题或建议，欢迎在 GitHub 提 Issue。

---

**相关资源**：
- [MockAPI GitHub](https://github.com/fynntang/MockAPI)
- [Go 官方文档](https://golang.org/doc/)
- [goja JavaScript 引擎](https://github.com/dop251/goja)
- [OpenAPI 规范](https://swagger.io/specification/)
