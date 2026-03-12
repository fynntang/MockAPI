<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as m from '$lib/paraglide/messages';
	import { getLocale, setLocale, locales } from '$lib/paraglide/runtime';
	import logo from '$lib/assets/logo.svg';
	import { theme } from '$lib/stores/theme';

	let activeTab = $state('rest');
	let downloadTab = $state('go');

	function switchLocale(newLocale: string) {
		setLocale(newLocale as 'en' | 'zh-hans', { reload: true });
	}

	function toggleTheme() {
		theme.toggle();
	}
</script>

<svelte:head>
	<title>MockAPI — {m.hero_title_1()} {m.hero_title_2()}</title>
	<link rel="preconnect" href="https://fonts.googleapis.com" />
	<link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800;900&display=swap" rel="stylesheet" />
</svelte:head>

<!-- Navigation -->
<nav class="fixed top-0 w-full z-50 backdrop-blur-xl bg-white/80 dark:bg-zinc-950/80 border-b border-zinc-200 dark:border-zinc-800">
	<div class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
		<a href="/" class="flex items-center gap-2">
			<img src={logo} alt="MockAPI" class="h-6 w-6">
			<span class="text-xl font-bold">MockAPI</span>
		</a>
		<div class="hidden md:flex items-center gap-8">
			<a href="#features" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_features()}</a>
			<a href="#screenshots" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_screenshots()}</a>
			<a href="/pricing" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_pricing()}</a>
			<a href="/feedback" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">{m.nav_feedback()}</a>
			<a href="https://github.com/fynntang/MockAPI" class="text-sm text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">GitHub</a>
		</div>
		<div class="flex items-center gap-3">
			<!-- Theme Switcher -->
			<button onclick={toggleTheme} class="p-2 rounded-lg bg-zinc-100 dark:bg-zinc-900 text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white transition-colors">
				{#if $theme === 'dark'}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z"></path></svg>
				{:else}
					<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z"></path></svg>
				{/if}
			</button>
			<!-- Language Switcher -->
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

<!-- Hero Section -->
<section class="relative min-h-screen flex items-center justify-center overflow-hidden pt-16">
	<div class="absolute inset-0 bg-white dark:bg-zinc-950"></div>
	<div class="absolute inset-0 bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(200,200,200,0.3),transparent)] dark:bg-[radial-gradient(ellipse_80%_50%_at_50%_-20%,rgba(63,63,70,0.3),transparent)]"></div>
	<div class="absolute inset-0 bg-[linear-gradient(rgba(200,200,200,0.1)_1px,transparent_1px),linear-gradient(90deg,rgba(200,200,200,0.1)_1px,transparent_1px)] dark:bg-[linear-gradient(rgba(63,63,70,0.1)_1px,transparent_1px),linear-gradient(90deg,rgba(63,63,70,0.1)_1px,transparent_1px)] bg-[size:64px_64px] [mask-image:radial-gradient(ellipse_60%_50%_at_50%_0%,black_40%,transparent_100%)]"></div>
	
	<div class="relative z-10 text-center px-6 py-20">
		<div class="inline-flex items-center gap-2 px-4 py-2 rounded-full bg-zinc-100 dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 mb-8">
			<span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
			<span class="text-sm text-zinc-700 dark:text-zinc-300">{m.hero_badge()}</span>
		</div>
		
		<h1 class="text-5xl md:text-7xl font-black mb-6 leading-tight">
			<span class="bg-gradient-to-r from-zinc-900 via-zinc-700 to-zinc-500 dark:bg-gradient-to-r dark:from-white dark:via-zinc-200 dark:to-zinc-400 bg-clip-text text-transparent">{m.hero_title_1()}</span>
			<br>
			<span class="text-zinc-900 dark:text-white">{m.hero_title_2()}</span>
		</h1>
		
		<p class="text-lg md:text-xl text-zinc-600 dark:text-zinc-400 max-w-2xl mx-auto mb-10 leading-relaxed">
			{m.hero_subtitle()}
			<span class="text-zinc-900 dark:text-white">{m.hero_subtitle_highlight()}</span>
		</p>
		
		<div class="flex flex-col sm:flex-row gap-4 justify-center mb-8">
			<a href="https://github.com/fynntang/MockAPI/releases">
				<Button size="lg" class="bg-zinc-900 dark:bg-white text-white dark:text-zinc-900 hover:bg-zinc-800 dark:hover:bg-zinc-200 text-lg px-8 py-6">
					{m.hero_cta_download()}
					<svg class="w-5 h-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"></path>
					</svg>
				</Button>
			</a>
			<a href="https://github.com/fynntang/MockAPI">
				<Button size="lg" variant="outline" class="px-8 py-6 border-zinc-300 dark:border-zinc-700 hover:bg-zinc-100 dark:hover:bg-zinc-900 gap-2">
					<svg class="w-5 h-5" fill="currentColor" viewBox="0 0 24 24"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-6.27 0-1.38.48-2.535 1.26-3.435-.12-.315-.54-1.65.12-3.435 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.785.24 3.12.12 3.435.78.9 1.26 2.055 1.26 3.435 0 4.95-2.805 5.97-5.475 6.27.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0024 12c0-6.63-5.37-12-12-12z"/></svg>
					<span>{m.hero_cta_github()}</span>
				</Button>
			</a>
		</div>
		
		<div class="max-w-5xl mx-auto mt-8">
			<div class="rounded-2xl overflow-hidden border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/50 p-1">
				<img src="/screenshots/00-main-routes.png" alt="MockAPI Main Interface" class="w-full rounded-xl">
			</div>
		</div>
	</div>
</section>

<!-- Stats Section -->
<section class="py-20 px-6 border-y border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/30">
	<div class="max-w-6xl mx-auto">
		<div class="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
			<div>
				<div class="text-5xl font-black text-zinc-900 dark:text-white">24x</div>
				<div class="text-zinc-600 dark:text-zinc-400 mt-2">{m.stats_faster()}</div>
			</div>
			<div>
				<div class="text-5xl font-black text-zinc-900 dark:text-white">94%</div>
				<div class="text-zinc-600 dark:text-zinc-400 mt-2">{m.stats_memory()}</div>
			</div>
			<div>
				<div class="text-5xl font-black text-zinc-900 dark:text-white">0</div>
				<div class="text-zinc-600 dark:text-zinc-400 mt-2">{m.stats_deps()}</div>
			</div>
			<div>
				<div class="text-5xl font-black text-zinc-900 dark:text-white">12+</div>
				<div class="text-zinc-600 dark:text-zinc-400 mt-2">{m.stats_templates()}</div>
			</div>
		</div>
	</div>
</section>

<!-- Screenshots Gallery -->
<section id="screenshots" class="py-24 px-6">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-16">
			<h2 class="text-4xl md:text-5xl font-bold mb-4">{m.screenshots_title()}</h2>
			<p class="text-xl text-zinc-600 dark:text-zinc-400">{m.screenshots_subtitle()}</p>
		</div>
		
		<div class="flex flex-wrap justify-center gap-2 mb-8">
			<button onclick={() => activeTab = 'rest'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'rest' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">⚡ {m.screenshots_rest()}</button>
			<button onclick={() => activeTab = 'graphql'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'graphql' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">◈ {m.screenshots_graphql()}</button>
			<button onclick={() => activeTab = 'grpc'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'grpc' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">🔄 {m.screenshots_grpc()}</button>
			<button onclick={() => activeTab = 'websocket'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'websocket' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">🔌 {m.screenshots_websocket()}</button>
			<button onclick={() => activeTab = 'logs'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {activeTab === 'logs' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">📋 {m.screenshots_logs()}</button>
		</div>
		
		{#if activeTab === 'rest'}
			<div class="grid md:grid-cols-3 gap-6">
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/00-main-routes.png" alt="Routes Dashboard" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_routes_dashboard()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_routes_dashboard_desc()}</p>
					</div>
				</Card>
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/01-add-route-modal.png" alt="Add Route Modal" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_add_route()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_add_route_desc()}</p>
					</div>
				</Card>
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/02-quick-templates.png" alt="Quick Templates" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_quick_templates()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_quick_templates_desc()}</p>
					</div>
				</Card>
			</div>
		{:else if activeTab === 'graphql'}
			<div class="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto">
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/03-graphql-tab.png" alt="GraphQL Tab" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_graphql_tab()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_graphql_tab_desc()}</p>
					</div>
				</Card>
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/04-add-graphql.png" alt="Add GraphQL Mock" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_add_graphql()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_add_graphql_desc()}</p>
					</div>
				</Card>
			</div>
		{:else if activeTab === 'grpc'}
			<div class="grid md:grid-cols-3 gap-6">
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/05-grpc-tab.png" alt="gRPC Tab" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_grpc_tab()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_grpc_tab_desc()}</p>
					</div>
				</Card>
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/06-add-grpc.png" alt="Add gRPC Mock" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_add_grpc()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_add_grpc_desc()}</p>
					</div>
				</Card>
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/07-import-proto-grpc.png" alt="Import Proto" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_import_proto()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_import_proto_desc()}</p>
					</div>
				</Card>
			</div>
		{:else if activeTab === 'websocket'}
			<div class="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto">
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/08-websocket-tab.png" alt="WebSocket Tab" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_websocket_tab()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_websocket_tab_desc()}</p>
					</div>
				</Card>
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/09-add-websocket.png" alt="Add WebSocket Mock" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_add_websocket()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_add_websocket_desc()}</p>
					</div>
				</Card>
			</div>
		{:else if activeTab === 'logs'}
			<div class="max-w-2xl mx-auto">
				<Card class="group bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all overflow-hidden">
					<img src="/screenshots/10-request-logs.png" alt="Request Logs" class="w-full">
					<div class="p-4 border-t border-zinc-200 dark:border-zinc-800">
						<h4 class="font-semibold mb-1">{m.screenshots_request_logs()}</h4>
						<p class="text-sm text-zinc-600 dark:text-zinc-400">{m.screenshots_request_logs_desc()}</p>
					</div>
				</Card>
			</div>
		{/if}
	</div>
</section>

<!-- Features Grid -->
<section id="features" class="py-24 px-6 border-t border-zinc-200 dark:border-zinc-800">
	<div class="max-w-6xl mx-auto">
		<div class="text-center mb-16">
			<h2 class="text-4xl md:text-5xl font-bold mb-4">{m.features_title()}</h2>
			<p class="text-xl text-zinc-600 dark:text-zinc-400">{m.features_subtitle()}</p>
		</div>
		
		<div class="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">⚡</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_rest()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_rest_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">📜</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_js()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_js_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">📥</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_swagger()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_swagger_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">🔌</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_websocket()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_websocket_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">◈</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_graphql()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_graphql_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">🔄</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_grpc()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_grpc_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">🎨</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_webui()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_webui_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">🔀</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_proxy()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_proxy_desc()}</p>
			</Card>
			<Card class="p-8 bg-zinc-50 dark:bg-zinc-900/30 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
				<div class="text-4xl mb-4">🚀</div>
				<h3 class="text-xl font-semibold mb-2">{m.features_zerodeps()}</h3>
				<p class="text-zinc-600 dark:text-zinc-400">{m.features_zerodeps_desc()}</p>
			</Card>
		</div>
	</div>
</section>

<!-- Quick Start -->
<section class="py-24 px-6">
	<div class="max-w-4xl mx-auto">
		<div class="text-center mb-16">
			<h2 class="text-4xl md:text-5xl font-bold mb-4">{m.quickstart_title()}</h2>
		</div>
		
		<div class="rounded-2xl overflow-hidden border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/50 shadow-2xl">
			<div class="flex items-center gap-2 px-4 py-3 bg-zinc-100 dark:bg-zinc-800/50 border-b border-zinc-200 dark:border-zinc-800">
				<div class="w-3 h-3 rounded-full bg-red-500"></div>
				<div class="w-3 h-3 rounded-full bg-yellow-500"></div>
				<div class="w-3 h-3 rounded-full bg-green-500"></div>
				<span class="ml-4 text-xs text-zinc-500 font-mono">terminal</span>
			</div>
			<div class="p-6 font-mono text-sm space-y-4">
				<div><span class="text-zinc-500"># {m.quickstart_install()}</span></div>
				<div><span class="text-zinc-600 dark:text-zinc-400">$</span> <span class="text-zinc-900 dark:text-white">go install github.com/fynntang/MockAPI@latest</span></div>
				<div class="pt-4"><span class="text-zinc-500"># {m.quickstart_start()}</span></div>
				<div><span class="text-zinc-600 dark:text-zinc-400">$</span> <span class="text-zinc-900 dark:text-white">mockapi serve</span></div>
				<div class="pt-4 text-emerald-500">✓ {m.quickstart_running()}</div>
				<div class="text-emerald-500">✓ {m.quickstart_webui()}</div>
			</div>
		</div>
	</div>
</section>

<!-- Download -->
<section id="download" class="py-24 px-6 border-t border-zinc-200 dark:border-zinc-800">
	<div class="max-w-3xl mx-auto text-center">
		<h2 class="text-4xl md:text-5xl font-bold mb-4">{m.download_title()}</h2>
		<p class="text-xl text-zinc-600 dark:text-zinc-400 mb-12">{m.download_subtitle()}</p>
		
		<!-- Download Tabs -->
		<div class="flex justify-center gap-2 mb-8">
			<button onclick={() => downloadTab = 'go'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {downloadTab === 'go' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">🐹 {m.download_go()}</button>
			<button onclick={() => downloadTab = 'docker'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {downloadTab === 'docker' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">🐳 {m.download_docker()}</button>
			<button onclick={() => downloadTab = 'binary'} class="px-4 py-2 rounded-lg text-sm font-medium transition-all {downloadTab === 'binary' ? 'bg-zinc-800 text-white' : 'text-zinc-600 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white'}">📦 {m.download_binary()}</button>
		</div>
		
		<!-- Download Content -->
		<div class="rounded-2xl border border-zinc-200 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-900/50 p-8">
			{#if downloadTab === 'go'}
				<div class="text-left">
					<h3 class="text-lg font-semibold mb-3">🐹 {m.download_go()}</h3>
					<code class="block bg-white dark:bg-zinc-950 p-4 rounded-lg text-sm text-zinc-700 dark:text-zinc-300 font-mono border border-zinc-200 dark:border-zinc-800">go install github.com/fynntang/MockAPI@latest</code>
				</div>
			{:else if downloadTab === 'docker'}
				<div class="text-left">
					<h3 class="text-lg font-semibold mb-3">🐳 {m.download_docker()}</h3>
					<code class="block bg-white dark:bg-zinc-950 p-4 rounded-lg text-sm text-zinc-700 dark:text-zinc-300 font-mono border border-zinc-200 dark:border-zinc-800">docker run -p 8088:8088 mockapi</code>
				</div>
			{:else if downloadTab === 'binary'}
				<div class="text-left">
					<h3 class="text-lg font-semibold mb-3">📦 {m.download_binary()}</h3>
					<a href="https://github.com/fynntang/MockAPI/releases" class="inline-flex items-center gap-2 text-zinc-600 dark:text-zinc-300 hover:text-zinc-900 dark:hover:text-white transition-colors">
						{m.download_github()} →
					</a>
				</div>
			{/if}
		</div>
	</div>
</section>

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