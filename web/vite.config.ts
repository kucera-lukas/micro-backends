import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import checker from "vite-plugin-checker";
import { VitePluginRadar } from "vite-plugin-radar";

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    checker({
      typescript: true,
      eslint: {
        lintCommand: 'eslint "./src/**/*.{ts,tsx}" --fix',
      },
    }),
    VitePluginRadar({
      analytics: { id: process.env.VITE_GA_MEASUREMENT_ID as string },
    }),
  ],
});
