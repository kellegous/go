import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  root: "ui",
  plugins: [react()],
  publicDir: "pub",
  base: "/s/",
  build: {
    outDir: "../internal/ui/assets",
    chunkSizeWarningLimit: 10000,
    assetsDir: ".",
    emptyOutDir: true,
    rollupOptions: {
      input: {
        edit: "./ui/edit/index.html",
        links: "./ui/index.html",
      },
    },
  },
  server: {
    proxy: {
      "/api": "https://go.finch-mahi.ts.net",
    },
  },
});
