<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Card } from '$lib/components/ui/card';
	import * as m from '$lib/paraglide/messages';
	import { getLocale, setLocale } from '$lib/paraglide/runtime';
	import logo from '$lib/assets/logo.svg';

	function switchLocale(newLocale: string) {
		setLocale(newLocale as 'en' | 'zh-hans', { reload: true });
	}
</script>

<svelte:head>
	<title>{m.nav_pricing()} — MockAPI</title>
</svelte:head>

<!-- Navigation -->
<nav class="fixed top-0 w-full z-50 backdrop-blur-xl bg-zinc-950/80 border-b border-zinc-800">
	<div class="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
		<a href="/" class="flex items-center gap-2">
			<img src={logo} alt="MockAPI" class="h-6 w-6">
			<span class="text-xl font-bold">MockAPI</span>
		</a>
		<div class="hidden md:flex items-center gap-8">
			<a href="/#features" class="text-sm text-zinc-400 hover:text-white transition-colors">{m.nav_features()}</a>
			<a href="/#screenshots" class="text-sm text-zinc-400 hover:text-white transition-colors">{m.nav_screenshots()}</a>
			<a href="/pricing" class="text-sm text-white font-medium">{m.nav_pricing()}</a>
			<a href="/feedback" class="text-sm text-zinc-400 hover:text-white transition-colors">{m.nav_feedback()}</a>
			<a href="https://github.com/fynntang/MockAPI" class="text-sm text-zinc-400 hover:text-white transition-colors">GitHub</a>
		</div>
		<div class="flex items-center gap-3">
			<div class="flex gap-1 bg-zinc-900 rounded-lg p-1">
				<button onclick={() => switchLocale('en')} class="px-2 py-1 text-xs rounded {getLocale() === 'en' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}">EN</button>
				<button onclick={() => switchLocale('zh-hans')} class="px-2 py-1 text-xs rounded {getLocale() === 'zh-hans' ? 'bg-zinc-800 text-white' : 'text-zinc-400 hover:text-white'}">中文</button>
			</div>
			<a href="https://github.com/fynntang/MockAPI/releases">
				<Button class="bg-white text-zinc-900 hover:bg-zinc-200">{m.nav_download()}</Button>
			</a>
		</div>
	</div>
</nav>

<!-- Hero -->
<section class="pt-32 pb-16 px-6">
	<div class="max-w-4xl mx-auto text-center">
		<h1 class="text-5xl md:text-6xl font-black mb-6">{m.pricing_title()}</h1>
		<p class="text-xl text-zinc-400">{m.pricing_subtitle()}</p>
	</div>
</section>

<!-- Pricing Cards -->
<section class="py-16 px-6">
	<div class="max-w-5xl mx-auto">
		<div class="grid md:grid-cols-2 gap-8">
			<!-- Free -->
			<Card class="p-8 bg-zinc-900/50 border-zinc-800 rounded-3xl">
				<div class="mb-8">
					<h2 class="text-2xl font-bold mb-2">{m.pricing_free()}</h2>
					<div class="flex items-baseline gap-2">
						<span class="text-6xl font-black">{m.pricing_free_price()}</span>
						<span class="text-zinc-400">{m.pricing_free_period()}</span>
					</div>
				</div>
				
				<p class="text-zinc-400 mb-8">{m.pricing_free_desc()}</p>
				
				<ul class="space-y-4 mb-8">
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_core()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_swagger()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_js()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_webui()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_mit()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-emerald-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_community()}</span>
					</li>
				</ul>
				
				<a href="https://github.com/fynntang/MockAPI/releases" class="block">
					<Button size="lg" class="w-full bg-white text-zinc-900 hover:bg-zinc-200">{m.pricing_download_free()}</Button>
				</a>
			</Card>
			
			<!-- Pro -->
			<Card class="relative p-8 bg-zinc-900/50 border-zinc-700 rounded-3xl">
				<div class="absolute -top-4 right-8 bg-white text-zinc-900 px-4 py-1 rounded-full text-sm font-bold">{m.pricing_coming_soon()}</div>
				
				<div class="mb-8">
					<h2 class="text-2xl font-bold mb-2">{m.pricing_pro()}</h2>
					<div class="flex items-baseline gap-2">
						<span class="text-6xl font-black">{m.pricing_pro_price()}</span>
						<span class="text-zinc-400">{m.pricing_pro_period()}</span>
					</div>
				</div>
				
				<p class="text-zinc-400 mb-8">{m.pricing_pro_desc()}</p>
				
				<ul class="space-y-4 mb-8">
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-white mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_all()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-white mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_team()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-white mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_cloud()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-white mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_priority()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-white mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_analytics()}</span>
					</li>
					<li class="flex items-start gap-3">
						<svg class="w-5 h-5 text-white mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
						<span>{m.pricing_feature_domains()}</span>
					</li>
				</ul>
				
				<Button size="lg" disabled class="w-full bg-zinc-800 text-zinc-400 cursor-not-allowed">{m.pricing_coming_soon()}</Button>
			</Card>
		</div>
	</div>
</section>

<!-- FAQ -->
<section class="py-24 px-6 border-t border-zinc-800">
	<div class="max-w-3xl mx-auto">
		<h2 class="text-3xl font-bold mb-12 text-center">{m.pricing_faq_title()}</h2>
		
		<div class="space-y-8">
			<div>
				<h3 class="text-lg font-semibold mb-2">{m.pricing_faq_q1()}</h3>
				<p class="text-zinc-400">{m.pricing_faq_a1()}</p>
			</div>
			<div>
				<h3 class="text-lg font-semibold mb-2">{m.pricing_faq_q2()}</h3>
				<p class="text-zinc-400">{m.pricing_faq_a2()}</p>
			</div>
			<div>
				<h3 class="text-lg font-semibold mb-2">{m.pricing_faq_q3()}</h3>
				<p class="text-zinc-400">{m.pricing_faq_a3()}</p>
			</div>
			<div>
				<h3 class="text-lg font-semibold mb-2">{m.pricing_faq_q4()}</h3>
				<p class="text-zinc-400">{m.pricing_faq_a4()}</p>
			</div>
		</div>
	</div>
</section>

<!-- Footer -->
<footer class="py-16 px-6 border-t border-zinc-800">
	<div class="max-w-6xl mx-auto">
		<div class="flex flex-col md:flex-row justify-between items-center gap-8">
			<div class="flex items-center gap-2">
				<img src={logo} alt="MockAPI" class="h-6 w-6">
				<span class="text-xl font-bold">MockAPI</span>
			</div>
			<div class="flex items-center gap-6 text-sm text-zinc-400">
				<a href="https://github.com/fynntang/MockAPI" class="hover:text-white transition-colors">{m.footer_github()}</a>
				<a href="https://github.com/fynntang/MockAPI/issues" class="hover:text-white transition-colors">{m.footer_issues()}</a>
				<a href="mailto:contact@mockapi.work" class="hover:text-white transition-colors">{m.footer_contact()}</a>
			</div>
		</div>
		<div class="mt-8 pt-8 border-t border-zinc-800 text-center text-sm text-zinc-500">
			{m.footer_copyright()}
		</div>
	</div>
</footer>