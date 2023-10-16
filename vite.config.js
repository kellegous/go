import { defineConfig } from 'vite';
import { resolve } from 'path';

export default defineConfig({
	root: 'ui',
	publicDir: 'pub',
	build: {
		outDir: resolve(__dirname, 'pkg/web/ui'),
		assertsDir: '.',
		rollupOptions: {
			input: {
				main: resolve(__dirname, 'ui/index.html'),
			}
		},
	},
});