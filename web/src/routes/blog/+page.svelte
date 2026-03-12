<script lang="ts">
	import * as m from '$lib/paraglide/messages';
	import { getLocale, setLocale } from '$lib/paraglide/runtime';
	import logo from '$lib/assets/logo.svg';
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import { theme } from '$lib/stores/theme';

	// 博客文章列表
	const posts = [
		{
			slug: 'performance-optimization-24x',
			title: '我用 Go 写了一个 Mock 服务器，性能提升了 24 倍',
			description: '零依赖、单二进制、功能完整、高性能的 API Mock 服务器',
			date: '2026-03-12',
			tags: ['Go', 'MockAPI', 'Performance', 'Open Source'],
			readTime: '8 min'
		}
	];

	function switchLocale(newLocale: string) {
		setLocale(newLocale as 'en' | 'zh-hans', { reload: true });
	}

	function toggleTheme() {
		theme.toggle();
	}
</script>

<svelte:head>
	<title>Blog — MockAPI</title>
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

<!-- Hero -->
<section class="pt-32 pb-16 px-6">
	<div class="max-w-4xl mx-auto text-center">
		<h1 class="text-5xl md:text-6xl font-black mb-6">Blog</h1>
		<p class="text-xl text-zinc-600 dark:text-zinc-400">开发日志、技术分享、产品更新</p>
	</div>
</section>

<!-- Blog Posts -->
<section class="py-16 px-6">
	<div class="max-w-4xl mx-auto">
		<div class="space-y-6">
			{#each posts as post}
				<a href="/blog/{post.slug}">
					<Card class="group p-6 bg-zinc-50 dark:bg-zinc-900/50 border-zinc-200 dark:border-zinc-800 hover:border-zinc-400 dark:hover:border-zinc-700 transition-all">
						<div class="flex flex-col md:flex-row md:items-start gap-4">
							<div class="flex-1">
								<h2 class="text-xl font-semibold mb-2 group-hover:text-emerald-600 dark:group-hover:text-emerald-400 transition-colors">{post.title}</h2>
								<p class="text-zinc-600 dark:text-zinc-400 mb-4">{post.description}</p>
								<div class="flex flex-wrap items-center gap-3 text-sm">
									<span class="text-zinc-500">{post.date}</span>
									<span class="text-zinc-300 dark:text-zinc-700">•</span>
									<span class="text-zinc-500">{post.readTime} read</span>
									{#each post.tags as tag}
										<span class="px-2 py-0.5 rounded bg-zinc-200 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400">{tag}</span>
									{/each}
								</div>
							</div>
							<svg class="w-6 h-6 text-zinc-400 group-hover:text-emerald-500 transition-colors flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
							</svg>
						</div>
					</Card>
				</a>
			{/each}
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