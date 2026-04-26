import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

const BACKEND_URL = process.env.BACKEND_URL ?? 'http://localhost:8080';

export default defineConfig({
	plugins: [sveltekit()],
	server: {
		port: 5173,
		strictPort: false,
		proxy: {
			'/api': {
				target: BACKEND_URL,
				changeOrigin: true,
				secure: false
			}
		}
	}
});
