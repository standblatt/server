import { sveltekit } from '@sveltejs/kit/vite';
import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vite';
import linguiForSvelte from 'lingui-for-svelte/unplugin/vite';
import linguiMacro from 'unplugin-lingui-macro/vite';

export default defineConfig({
	plugins: [linguiMacro(), linguiForSvelte(), tailwindcss(), sveltekit()],
	server: {
		// The help shows ../CHANGELOG.md.
		fs: { allow: ['..'] },
		proxy: {
			'/api': { target: 'http://localhost:8080', ws: true },
			'/admin': 'http://localhost:8080'
		}
	}
});
