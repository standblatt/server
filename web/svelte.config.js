import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
	preprocess: vitePreprocess(),
	compilerOptions: { runes: true },
	kit: {
		paths: { relative: false },
		serviceWorker: { register: false },
		version: { pollInterval: 60000 },
		alias: { $locales: 'src/locales' },
		adapter: adapter({
			pages: '../internal/webui/dist',
			assets: '../internal/webui/dist',
			fallback: 'index.html'
		})
	}
};
