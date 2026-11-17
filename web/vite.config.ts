import { sveltekit } from '@sveltejs/kit/vite';
import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import linguiForSvelte from 'lingui-for-svelte/unplugin/vite';
import linguiMacro from 'unplugin-lingui-macro/vite';

export default defineConfig({
	plugins: [
		linguiMacro(),
		linguiForSvelte(),
		tailwindcss(),
		sveltekit({
			preprocess: vitePreprocess(),
			compilerOptions: { runes: true },
			paths: { relative: false },
			serviceWorker: { register: false },
			version: { pollInterval: 60000 },
			adapter: adapter({
				pages: '../internal/webui/dist',
				assets: '../internal/webui/dist',
				fallback: 'index.html'
			})
		})
	],
	server: {
		// The help shows ../CHANGELOG.md.
		fs: { allow: ['..'] },
		proxy: {
			'/api': { target: 'http://localhost:8080', ws: true },
			'/admin': 'http://localhost:8080'
		}
	}
});
