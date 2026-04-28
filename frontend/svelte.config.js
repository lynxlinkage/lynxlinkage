import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		adapter: adapter({
			pages: 'build',
			assets: 'build',
			// SPA shell for routes that opt out of prerender (e.g. /admin, /login).
			// Kept distinct from index.html so the prerendered home page isn't
			// overwritten.
			fallback: '200.html',
			precompress: false,
			strict: true
		}),
		prerender: {
			handleHttpError: ({ path, referrer, message }) => {
				if (path.startsWith('/api/')) return;
				console.warn(`Prerender error at ${path} (referred by ${referrer}): ${message}`);
			}
		},
		alias: {
			$lib: 'src/lib'
		}
	}
};

export default config;
