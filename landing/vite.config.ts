import tailwindcss from '@tailwindcss/vite';
import { defineConfig } from 'vitest/config';
import { loadEnv } from 'vite';
import { playwright } from '@vitest/browser-playwright';
import { sveltekit } from '@sveltejs/kit/vite';

export default defineConfig(({ mode }) => {
	// Load env file based on `mode` in the current working directory.
	// This ensures .env file is loaded before the config is used
	const env = loadEnv(mode, process.cwd(), 'PUBLIC_');
	
	// Log to verify env is loaded (only in dev mode)
	if (mode === 'development') {
		console.log('Loaded PUBLIC_API_KEY:', env.PUBLIC_API_KEY ? `${env.PUBLIC_API_KEY.substring(0, 8)}...` : 'NOT FOUND');
		console.log('Loaded PUBLIC_API_BASE_URL:', env.PUBLIC_API_BASE_URL || 'NOT FOUND');
	}
	
	return {
		plugins: [tailwindcss(), sveltekit()],
		test: {
			expect: { requireAssertions: true },
			projects: [
				{
					extends: './vite.config.ts',
					test: {
						name: 'client',
						browser: {
							enabled: true,
							provider: playwright(),
							instances: [{ browser: 'chromium', headless: true }]
						},
						include: ['src/**/*.svelte.{test,spec}.{js,ts}'],
						exclude: ['src/lib/server/**']
					}
				},
				{
					extends: './vite.config.ts',
					test: {
						name: 'server',
						environment: 'node',
						include: ['src/**/*.{test,spec}.{js,ts}'],
						exclude: ['src/**/*.svelte.{test,spec}.{js,ts}']
					}
				}
			]
		}
	};
});
