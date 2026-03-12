import adapters from '@sveltejs/adapter-cloudflare';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	kit: {
		adapter: adapters(),
		alias: {
			$lib: './src/lib'
		}
	},
	preprocess: [],
	extensions: ['.svelte']
};

export default config;