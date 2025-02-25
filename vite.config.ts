import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  root: "ui",
  plugins: [react()],
  publicDir: "pub",
  build: {
    outDir: "../internal/ui/assets",
    chunkSizeWarningLimit: 10000,
    assetsDir: ".",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        edit: './ui/edit.html',
        links: './ui/links.html',
      }
    }
  },
  server: {
    proxy: {
      '/api': 'https://go.finch-mahi.ts.net',
    },
  },
})
