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

	const isZh = $derived(getLocale() === 'zh-hans');
</script>

<svelte:head>
	<title>{isZh ? '我用 Go 写了一个 Mock 服务器，性能提升了 24 倍' : 'I Built a Mock Server in Go with 24x Performance Boost'} — MockAPI Blog</title>
	<meta name="description" content={isZh ? '零依赖、单二进制、功能完整、高性能的 API Mock 服务器' : 'Zero dependencies, single binary, full-featured, high-performance API mock server'} />
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
			{isZh ? '返回博客' : 'Back to Blog'}
		</a>
		
		{#if isZh}
			<!-- Chinese Version -->
			<h1 class="text-4xl font-bold mb-4">我用 Go 写了一个 Mock 服务器，性能提升了 24 倍</h1>
			
			<div class="flex items-center gap-4 text-sm text-zinc-500 mb-8">
				<span>2026-03-12</span>
				<span>•</span>
				<span>8 分钟阅读</span>
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
				
				<h3>2. 标准库足够强大</h3>
				<p>Go 的 <code>net/http</code> 性能优异，<code>encoding/json</code> 够用，<code>embed</code> 可以把 Web UI 嵌入二进制。</p>
				
				<h3>3. 跨平台编译</h3>
				<p>一次编译，到处运行。</p>
				
				<h3>4. 部署简单</h3>
				<p>没有复杂的依赖关系，不需要 Docker 也能轻松部署。</p>
				
				<h2>性能优化：从 O(n) 到 O(1)</h2>
				
				<p>项目初期，路由匹配是线性搜索。当路由数量增加时，性能急剧下降。</p>
				
				<h3>解决方案：RouteIndex</h3>
				
				<p>我设计了一个两层索引结构：</p>
				<ul>
					<li><strong>精确匹配</strong>：直接查 map，O(1)</li>
					<li><strong>参数路由</strong>：按前缀分组，大幅减少候选集</li>
					<li><strong>通配符路由</strong>：按方法分组</li>
				</ul>
				
				<h3>优化效果</h3>
				
				<p>性能提升 <strong>24 倍</strong>，内存减少 <strong>94%</strong>。这就是算法优化的力量。</p>
				
				<h2>核心功能</h2>
				
				<h3>动态路由</h3>
				<p>支持 <code>:param</code> 参数和 <code>*</code> 通配符。</p>
				
				<h3>JavaScript 脚本引擎</h3>
				<p>当静态响应不够用时，可以用 JavaScript 动态生成响应。</p>
				
				<h3>Swagger/OpenAPI 导入</h3>
				<p>有现成的 API 文档？一键导入，自动生成 Mock 路由。</p>
				
				<h2>内置 Web UI</h2>
				
				<p>打开浏览器就能管理：<code>http://localhost:8088/_ui</code></p>
				
				<p>功能：路由列表和管理、可视化添加路由、快速模板、请求日志查看、实时调试。</p>
				
				<h2>快速开始</h2>
				
				<pre><code class="language-bash"># 安装
go install github.com/fynntang/MockAPI@latest

# 启动
mockapi serve</code></pre>
				
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
				
				<p><em>Happy Mocking! 🦞</em></p>
			</div>
		{:else}
			<!-- English Version -->
			<h1 class="text-4xl font-bold mb-4">I Built a Mock Server in Go with 24x Performance Boost</h1>
			
			<div class="flex items-center gap-4 text-sm text-zinc-500 mb-8">
				<span>2026-03-12</span>
				<span>•</span>
				<span>8 min read</span>
			</div>
			
			<div class="prose prose-zinc dark:prose-invert max-w-none">
				<p class="lead">As a full-stack developer, I work with APIs every day. Frontend waiting for backend, third-party API rate limits, microservice coordination... all these scenarios need mock data.</p>
				
				<p>I've tried many mock tools, but something always felt off:</p>
				
				<ul>
					<li><strong>MockServer</strong>: Powerful but heavy, requires JVM</li>
					<li><strong>JSON Server</strong>: Lightweight but only REST, no dynamic responses</li>
					<li><strong>Postman Mock</strong>: Requires internet, paid for team collaboration</li>
				</ul>
				
				<p>So I decided to build my own. The goal was simple:</p>
				
				<blockquote>
					<strong>Zero dependencies, single binary, full-featured, high-performance</strong>
				</blockquote>
				
				<p>Three months later, MockAPI was born.</p>
				
				<h2>Results at a Glance</h2>
				
				<table>
					<thead>
						<tr><th>Metric</th><th>Value</th></tr>
					</thead>
					<tbody>
						<tr><td>Runtime Dependencies</td><td>0</td></tr>
						<tr><td>Installation</td><td><code>go install</code> one command</td></tr>
						<tr><td>Startup Time</td><td>&lt; 100ms</td></tr>
						<tr><td>Route Matching</td><td><strong>24x faster</strong></td></tr>
						<tr><td>Memory Usage</td><td><strong>94% less</strong></td></tr>
					</tbody>
				</table>
				
				<p>Features:</p>
				<ul>
					<li>✅ REST Mock (dynamic routes, path params, wildcards, conditional responses)</li>
					<li>✅ GraphQL Mock</li>
					<li>✅ WebSocket Mock</li>
					<li>✅ gRPC-Web Mock</li>
					<li>✅ JavaScript scripting engine</li>
					<li>✅ Swagger/OpenAPI import</li>
					<li>✅ Built-in Web UI</li>
					<li>✅ Hot reload</li>
				</ul>
				
				<h2>Why Go?</h2>
				
				<p>I chose Go not because I only know Go (though that's partly true), but because it fits perfectly:</p>
				
				<h3>1. Zero Runtime Dependencies</h3>
				<p>Compile to a single binary. Users don't need Node.js, JVM, or any other runtime. Download and run.</p>
				
				<h3>2. Powerful Standard Library</h3>
				<p>Go's <code>net/http</code> is fast, <code>encoding/json</code> is sufficient, <code>embed</code> lets you embed the Web UI in the binary.</p>
				
				<h3>3. Cross-Platform Compilation</h3>
				<p>Compile once, run everywhere.</p>
				
				<h3>4. Simple Deployment</h3>
				<p>No complex dependencies, deploy without Docker.</p>
				
				<h2>Performance: From O(n) to O(1)</h2>
				
				<p>Initially, route matching was linear search. Performance degraded as routes increased.</p>
				
				<h3>Solution: RouteIndex</h3>
				
				<p>I designed a two-layer index structure:</p>
				<ul>
					<li><strong>Exact match</strong>: Direct map lookup, O(1)</li>
					<li><strong>Param routes</strong>: Grouped by prefix, reduces candidates</li>
					<li><strong>Wildcard routes</strong>: Grouped by method</li>
				</ul>
				
				<h3>Results</h3>
				
				<p><strong>24x faster</strong>, <strong>94% less memory</strong>. That's the power of choosing the right data structure.</p>
				
				<h2>Core Features</h2>
				
				<h3>Dynamic Routes</h3>
				<p>Support <code>:param</code> and <code>*</code> wildcards.</p>
				
				<h3>JavaScript Engine</h3>
				<p>Generate dynamic responses with JavaScript when static isn't enough.</p>
				
				<h3>Swagger/OpenAPI Import</h3>
				<p>Import existing API docs, auto-generate mock routes.</p>
				
				<h2>Built-in Web UI</h2>
				
				<p>Manage in browser: <code>http://localhost:8088/_ui</code></p>
				
				<p>Features: Route management, visual editor, quick templates, request logs, real-time debugging.</p>
				
				<h2>Quick Start</h2>
				
				<pre><code class="language-bash"># Install
go install github.com/fynntang/MockAPI@latest

# Run
mockapi serve</code></pre>
				
				<h2>Project Links</h2>
				
				<ul>
					<li><strong>GitHub</strong>: https://github.com/fynntang/MockAPI</li>
					<li><strong>Website</strong>: https://mockapi.work</li>
				</ul>
				
				<p>If this project helps you, please Star ⭐</p>
				
				<h2>Final Thoughts</h2>
				
				<p>Building MockAPI taught me:</p>
				
				<ol>
					<li><strong>Simplicity wins</strong> - Zero dependencies beats feature bloat</li>
					<li><strong>Performance early</strong> - Good data structures beat late optimization</li>
					<li><strong>Developer experience matters</strong> - CLI and Web UI improve usability</li>
					<li><strong>Test-driven</strong> - 26 tests let me sleep well</li>
				</ol>
				
				<p><em>Happy Mocking! 🦞</em></p>
			</div>
		{/if}
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