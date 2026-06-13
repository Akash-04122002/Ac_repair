import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vuetify from 'vite-plugin-vuetify'

// Build output goes into ../backend/dist so the Go binary can embed it.
export default defineConfig({
  plugins: [vue(), vuetify({ autoImport: true })],
  build: {
    outDir: '../backend/dist',
    emptyOutDir: true,
  },
  server: {
    // During `npm run dev`, proxy API calls to the local Go server.
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
