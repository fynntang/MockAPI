<script lang="ts">
	import logo from '$lib/assets/logo.svg';
	import { Button } from '$lib/components/ui/button';
	import * as m from '$lib/paraglide/messages';
	import { getLocale, setLocale } from '$lib/paraglide/runtime';
	import { theme } from '$lib/stores/theme';

	function switchLocale(newLocale: string) {
		setLocale(newLocale as 'en' | 'zh-hans', { reload: true });
	}

	function toggleTheme() {
		theme.toggle();
	}
</script>

<svelte:head>
	<title>我用 Go 写了一个 Mock 服务器，性能提升了 24 倍 — MockAPI Blog</title>
	<meta name="description" content="零依赖、单二进制、功能完整、高性能的 API Mock 服务器" />
</svelte:head>

<!-- Navigation -->
<nav class="fixed top-0 w-full z-50 backdrop-blur-xl bg-white/80 dark:bg-zinc-950/80 border-b border-zinc-200 dark:border-zinc-800">
	<div class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
		<a href="/" class="flex items-center gap-2">
			<img src={logo} alt="MockAPI" class="h-6 w-6">
			<span class="text-xl font-bold">MockAPI</span>
		</a>
		<div class="hidden md:flex items-center gap-8">
			<a href="/#features" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_features()}</a>
			<a href="/#screenshots" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_screenshots()}</a>
			<a href="/pricing" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_pricing()}</a>
			<a href="/blog" class="text-sm text-zinc-900 dark:text-white font-medium">Blog</a>
			<a href="/feedback" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_feedback()}</a>
			<a href="https://github.com/fynntang/MockAPI" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">GitHub</a>
		</div>
		<div class="flex items-center gap-3">
			<button onclick={toggleTheme} class="p-2 rounded-lg bg-zinc-100 dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">
				{#if $theme === 'dark'}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
				{:else}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path></svg>
				{/if}
			</button>
			<div class="flex gap-1 bg-zinc-100 dark:bg-zinc-900 rounded-lg p-1">
				<button onclick={() => switchLocale('en')} class="px-2 py-1 text-xs rounded {getLocale() === 'en' ? 'bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">EN</button>
				<button onclick={() => switchLocale('zh-hans')} class="px-2 py-1 text-xs rounded {getLocale() === 'zh-hans' ? 'bg-white dark:bg-zinc-800 text-zinc-900 dark:text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">中文</button>
			</div>
			<a href="https://github.com/fynntang/MockAPI/releases">
				<Button class="bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 hover:bg-zinc-800 dark:hover:bg-zinc-200">{m.nav_download()}</Button>
			</a>
		</div>
	</div>
</nav>

<!-- Article -->
<article class="pt-32 pb-24 px-6">
	<div class="max-w-3xl mx-auto">
		<a href="/blog" class="inline-flex items-center gap-2 text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors mb-8">
			<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path></svg>
			Back to Blog
		</a>
		
		<h1 class="text-4xl font-bold mb-4">我用 Go 写了一个 Mock 服务器，性能提升了 24 倍</h1>
		
		<div class="flex items-center gap-4 text-sm text-zinc-500 mb-8">
			<span>2026-03-12</span>
			<span>•</span>
			<span>8 min read</span>
		</div>
		
		<div class="prose prose-zinc dark:prose-invert max-w-none">
			<p class="lead">作为一名全栈开发者，我每天都在和 API 打交道。前端开发等后端接口、第三方 API 限流、微服务联调......这些场景都需要 Mock 数据。</p>
			
			<p>试过不少 Mock 工具，但总觉得差点意思：</p>
			
			<ul>
				<li><strong>MockServer</strong>：功能强大但太重，需要 JVM</li>
				<li><strong>JSON Server</strong>：轻量但只支持 REST，不支持动态响应</li>
				<li><strong>Postman Mock</strong>：需要联网，付费才能团队协作</li>
			</ul>
			
			<p>于是我决定自己造轮子。目标很简单：</p>
			
			<blockquote>
				<strong>零依赖、单二进制、功能完整、高性能</strong>
			</blockquote>
			
			<p>三个月后，MockAPI 诞生了。</p>
			
			<h2>成果一览</h2>
			
			<table>
				<thead>
					<tr><th>指标</th><th>数值</th></tr>
				</thead>
				<tbody>
					<tr><td>运行时依赖</td><td>0</td></tr>
					<tr><td>安装方式</td><td><code>go install</code> 一条命令</td></tr>
					<tr><td>启动时间</td><td>&lt; 100ms</td></tr>
					<tr><td>路由匹配性能</td><td>提升 <strong>24 倍</strong></td></tr>
					<tr><td>内存占用</td><td>减少 <strong>94%</strong></td></tr>
				</tbody>
			</table>
			
			<p>支持的功能：</p>
			<ul>
				<li>✅ REST Mock（动态路由、路径参数、条件响应）</li>
				<li>✅ GraphQL Mock</li>
				<li>✅ WebSocket Mock</li>
				<li>✅ gRPC-Web Mock</li>
				<li>✅ JavaScript 脚本引擎</li>
				<li>✅ Swagger/OpenAPI 导入</li>
				<li>✅ 内置 Web UI</li>
				<li>✅ 热重载</li>
			</ul>
			
			<h2>技术选型：为什么是 Go？</h2>
			
			<p>选择 Go 不是因为我只会 Go（虽然这也有关系），而是它完美契合这个场景：</p>
			
			<h3>1. 零运行时依赖</h3>
			<p>编译成单个二进制文件，用户不需要安装 Node.js、JVM 或任何其他运行时。下载即用，上传即跑。</p>
			
			<pre><code class="language-bash"># 安装
go install github.com/fynntang/MockAPI@latest

# 启动
mockapi serve

# 完成！</code></pre>
			
			<h3>2. 标准库足够强大</h3>
			<p>Go 的 <code>net/http</code> 性能优异，<code>encoding/json</code> 够用，<code>embed</code> 可以把 Web UI 嵌入二进制。不需要引入重型框架。</p>
			
			<h3>3. 跨平台编译</h3>
			<p>一次编译，到处运行：</p>
			
			<pre><code class="language-bash">GOOS=darwin GOARCH=amd64 go build  # macOS
GOOS=linux GOARCH=amd64 go build   # Linux
GOOS=windows GOARCH=amd64 go build # Windows</code></pre>
			
			<h3>4. 部署简单</h3>
			<p>没有复杂的依赖关系，不需要 Docker 也能轻松部署。当然，如果你喜欢容器化：</p>
			
			<pre><code class="language-bash">docker run -p 8088:8088 mockapi</code></pre>
			
			<h2>性能优化：从 O(n) 到 O(1)</h2>
			
			<p>项目初期，路由匹配是这样的：</p>
			
			<pre><code class="language-go">// 线性搜索 O(n)
func (s *Server) matchRoute(method, path string) *Route {
    for _, route := range s.routes {
        if route.Method == method && route.Match(path) {
            return route
        }
    }
    return nil
}</code></pre>
			
			<p>当路由数量少时没问题，但随着功能增加，性能急剧下降：</p>
			
			<pre><code>100 个路由：  ~500ns
1000 个路由： ~5ms
10000 个路由：~50ms  // 不可接受！</code></pre>
			
			<h3>解决方案：RouteIndex</h3>
			
			<p>我设计了一个两层索引结构：</p>
			
			<pre><code>┌─────────────────────────────────────────┐
│              RouteIndex                  │
├─────────────────────────────────────────┤
│ exact: map["GET:/users"] -> Route       │  ← O(1)
├─────────────────────────────────────────┤
│ param: map["GET:/users"] -> []Route     │  ← O(k)
│        map["GET:/posts"] -> []Route     │
├─────────────────────────────────────────┤
│ wildcard: map["GET"] -> []Route         │  ← O(m)
└─────────────────────────────────────────┘</code></pre>
			
			<p>匹配逻辑：</p>
			<ol>
				<li><strong>精确匹配</strong>：直接查 map，O(1)</li>
				<li><strong>参数路由</strong>：按前缀分组，大幅减少候选集</li>
				<li><strong>通配符路由</strong>：按方法分组</li>
			</ol>
			
			<h3>优化效果</h3>
			
			<pre><code>优化前：BenchmarkRouteMatch-8    231218    5169 ns/op    6400 B/op
优化后：BenchmarkRouteMatch-8    5589621   214 ns/op     68 B/op

性能提升：24 倍
内存减少：94%</code></pre>
			
			<p>这就是算法优化的力量。没有黑魔法，只是把数据结构选对了。</p>
			
			<h2>核心功能实现</h2>
			
			<h3>1. 动态路由</h3>
			<p>支持 <code>:param</code> 参数和 <code>*</code> 通配符：</p>
			
			<pre><code class="language-yaml"># routes.json
- path: /users/:id
  method: GET
  response:
    body: |
      {
        "id": {{params.id}},
        "name": "User {{params.id}}"
      }

- path: /api/*
  method: ANY
  proxy: https://real-api.com</code></pre>
			
			<h3>2. JavaScript 脚本引擎</h3>
			<p>当静态响应不够用时，可以用 JavaScript 动态生成：</p>
			
			<pre><code class="language-javascript">// 条件响应
if (headers["x-api-key"] === "secret") {
  return { authorized: true, user: "admin" };
}

// 模拟延迟
sleep(100);
return { data: "delayed response" };</code></pre>
			
			<h3>3. Swagger/OpenAPI 导入</h3>
			<p>有现成的 API 文档？一键导入：</p>
			
			<pre><code class="language-bash">mockapi import swagger.yaml</code></pre>
			
			<p>自动解析所有端点，生成 Mock 路由。</p>
			
			<h2>内置 Web UI</h2>
			
			<p>不想写配置文件？打开浏览器就能管理：</p>
			
			<pre><code>http://localhost:8088/_ui</code></pre>
			
			<p>功能：</p>
			<ul>
				<li>📋 路由列表和管理</li>
				<li>➕ 可视化添加路由</li>
				<li>📝 快速模板（12+ 预设）</li>
				<li>📊 请求日志查看</li>
				<li>🔍 实时调试</li>
			</ul>
			
			<h2>对比其他方案</h2>
			
			<table>
				<thead>
					<tr><th>特性</th><th>MockAPI</th><th>MockServer</th><th>JSON Server</th><th>Postman Mock</th></tr>
				</thead>
				<tbody>
					<tr><td>零依赖</td><td>✅</td><td>❌ (JVM)</td><td>❌ (Node.js)</td><td>❌ (联网)</td></tr>
					<tr><td>动态响应</td><td>✅</td><td>✅</td><td>❌</td><td>✅</td></tr>
					<tr><td>GraphQL</td><td>✅</td><td>❌</td><td>❌</td><td>✅</td></tr>
					<tr><td>WebSocket</td><td>✅</td><td>✅</td><td>❌</td><td>❌</td></tr>
					<tr><td>gRPC</td><td>✅</td><td>❌</td><td>❌</td><td>❌</td></tr>
					<tr><td>Web UI</td><td>✅</td><td>❌</td><td>❌</td><td>✅</td></tr>
					<tr><td>开源</td><td>✅</td><td>✅</td><td>✅</td><td>❌</td></tr>
				</tbody>
			</table>
			
			<h2>快速开始</h2>
			
			<pre><code class="language-bash"># 安装
go install github.com/fynntang/MockAPI@latest

# 启动
mockapi serve

# 打开 Web UI
open http://localhost:8088/_ui</code></pre>
			
			<p>就这么简单！</p>
			
			<h2>项目地址</h2>
			
			<ul>
				<li><strong>GitHub</strong>: https://github.com/fynntang/MockAPI</li>
				<li><strong>官网</strong>: https://mockapi.work</li>
			</ul>
			
			<p>如果这个项目对你有帮助，欢迎 Star ⭐</p>
			
			<h2>写在最后</h2>
			
			<p>开发 MockAPI 的过程让我学到了很多：</p>
			
			<ol>
				<li><strong>简单往往更好</strong> - 零依赖比功能堆砌更重要</li>
				<li><strong>性能要趁早考虑</strong> - 好的数据结构胜过后期优化</li>
				<li><strong>开发者体验很重要</strong> - CLI 和 Web UI 大幅提升易用性</li>
				<li><strong>测试驱动开发</strong> - 26 个测试用例让我睡得安稳</li>
			</ol>
			
			<p>如果你也需要一个轻量级 Mock 服务器，不妨试试 MockAPI。有问题或建议，欢迎在 GitHub 提 Issue！</p>
			
			<p><em>Happy Mocking! 🦞</em></p>
		</div>
	</div>
</article>

<!-- Footer -->
<footer class="py-16 px-6 border-t border-zinc-200 dark:border-zinc-800">
	<div class="max-w-6xl mx-auto">
		<div class="flex flex-col md:flex-row justify-between items-center gap-8">
			<div class="flex items-center gap-2">
				<img src={logo} alt="MockAPI" class="h-6 w-6">
				<span class="text-xl font-bold">MockAPI</span>
			</div>
			<div class="flex items-center gap-6 text-sm text-zinc-600 dark:text-zinc-400">
				<a href="https://github.com/fynntang/MockAPI" class="hover:text-zinc-900 dark:hover:text-white transition-colors">{m.footer_github()}</a>
				<a href="https://github.com/fynntang/MockAPI/issues" class="hover:text-zinc-900 dark:hover:text-white transition-colors">{m.footer_issues()}</a>
				<a href="mailto:contact@mockapi.work" class="hover:text-zinc-900 dark:hover:text-white transition-colors">{m.footer_contact()}</a>
			</div>
		</div>
		<div class="mt-8 pt-8 border-t border-zinc-200 dark:border-zinc-800 text-center text-sm text-zinc-500">
			{m.footer_copyright()}
		</div>
	</div>
</footer>