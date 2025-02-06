import path from 'path';
import react from '@vitejs/plugin-react';
import { defineConfig } from 'vite';

export default defineConfig({
	base: '/',
	plugins: [react()],
	resolve: {
		alias: {
			'@': path.resolve(__dirname, './src'),
		},
	},
	server: {
		host: true,
		strictPort: true,
		port: 3000,
		origin: 'http://0.0.0.0:3000',
	},
	preview: {
		// host: true,
		strictPort: true,
		port: 3000,
	},
});
