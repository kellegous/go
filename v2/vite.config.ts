import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  root: "ui",
  plugins: [react()],
  publicDir: "pub",
  base: "/s/",
  esbuild: {
    // TODO(kellegous): everything in npm is always broken garbage.
    // esbuild has a security advisory and I need to override esbuild
    // to be >= 0.28.1. But there is no vite version that has that
    // update yet. So I have to put an override in package.json ...
    // but, of course, there is some jacked up shit in esbuild that
    // is papered over in a later vite ... which I am not on, so now
    // I have to paper over it too. Long story short, remove all this
    // crap when the override is removed.
    supported: {
      destructuring: true,
    },
  },
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
});
