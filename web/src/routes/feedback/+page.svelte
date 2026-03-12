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
	<title>{m.nav_feedback()} — MockAPI</title>
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
			<a href="/pricing" class="text-sm text-zinc-400 hover:text-white transition-colors">{m.nav_pricing()}</a>
			<a href="/feedback" class="text-sm text-white font-medium">{m.nav_feedback()}</a>
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
		<h1 class="text-5xl md:text-6xl font-black mb-6">{m.feedback_title()}</h1>
		<p class="text-xl text-zinc-400">{m.feedback_subtitle()}</p>
	</div>
</section>

<!-- Feedback Form -->
<section class="py-16 px-6">
	<div class="max-w-2xl mx-auto">
		<Card class="p-8 bg-zinc-900/50 border-zinc-800 rounded-2xl">
			<form action="https://formspree.io/f/xnnnnnnn" method="POST" class="space-y-6">
				<div class="grid md:grid-cols-2 gap-6">
					<div class="space-y-2">
						<label for="name" class="text-sm font-medium">{m.feedback_name()}</label>
						<input id="name" name="name" placeholder={m.feedback_name_placeholder()} class="w-full h-10 px-3 rounded-lg bg-zinc-950 border border-zinc-800 text-white placeholder-zinc-500" />
					</div>
					<div class="space-y-2">
						<label for="email" class="text-sm font-medium">{m.feedback_email()}</label>
						<input id="email" name="email" type="email" placeholder={m.feedback_email_placeholder()} class="w-full h-10 px-3 rounded-lg bg-zinc-950 border border-zinc-800 text-white placeholder-zinc-500" />
					</div>
				</div>
				
				<div class="space-y-2">
					<label for="type" class="text-sm font-medium">{m.feedback_type()}</label>
					<select id="type" name="type" class="w-full h-10 px-3 rounded-lg bg-zinc-950 border border-zinc-800 text-white">
						<option value="bug">🐛 {m.feedback_type_bug()}</option>
						<option value="feature">💡 {m.feedback_type_feature()}</option>
						<option value="question">❓ {m.feedback_type_question()}</option>
						<option value="other">💬 {m.feedback_type_other()}</option>
					</select>
				</div>
				
				<div class="space-y-2">
					<label for="message" class="text-sm font-medium">{m.feedback_message()}</label>
					<textarea id="message" name="message" placeholder={m.feedback_message_placeholder()} rows="6" class="w-full px-3 py-2 rounded-lg bg-zinc-950 border border-zinc-800 text-white placeholder-zinc-500 resize-none"></textarea>
				</div>
				
				<Button type="submit" size="lg" class="w-full bg-white text-zinc-900 hover:bg-zinc-200">
					{m.feedback_submit()}
				</Button>
			</form>
		</Card>
	</div>
</section>

<!-- Alternative Contact -->
<section class="py-16 px-6 border-t border-zinc-800">
	<div class="max-w-4xl mx-auto">
		<h2 class="text-2xl font-bold mb-8 text-center">{m.feedback_other()}</h2>
		
		<div class="grid md:grid-cols-3 gap-6">
			<a href="https://github.com/fynntang/MockAPI/issues" class="group">
				<Card class="p-6 bg-zinc-900/50 border-zinc-800 hover:border-zinc-700 transition-all text-center">
					<div class="text-4xl mb-4">🐙</div>
					<h3 class="font-semibold mb-2 group-hover:text-white">{m.feedback_github_issues()}</h3>
					<p class="text-sm text-zinc-400">{m.feedback_github_issues_desc()}</p>
				</Card>
			</a>
			
			<a href="mailto:contact@mockapi.work" class="group">
				<Card class="p-6 bg-zinc-900/50 border-zinc-800 hover:border-zinc-700 transition-all text-center">
					<div class="text-4xl mb-4">📧</div>
					<h3 class="font-semibold mb-2 group-hover:text-white">{m.feedback_email_contact()}</h3>
					<p class="text-sm text-zinc-400">contact@mockapi.work</p>
				</Card>
			</a>
			
			<a href="https://github.com/fynntang/MockAPI/discussions" class="group">
				<Card class="p-6 bg-zinc-900/50 border-zinc-800 hover:border-zinc-700 transition-all text-center">
					<div class="text-4xl mb-4">💬</div>
					<h3 class="font-semibold mb-2 group-hover:text-white">{m.feedback_discussions()}</h3>
					<p class="text-sm text-zinc-400">{m.feedback_discussions_desc()}</p>
				</Card>
			</a>
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