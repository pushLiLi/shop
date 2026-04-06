import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

import path from 'path'

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
        ws: true
      },
      '/media': {
        target: 'http://localhost:9000',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/media/, '')
      }
    }
  },
  build: {
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (id.includes('node_modules/vue/') || id.includes('node_modules/vue-router/') || id.includes('node_modules/pinia/')) {
            return 'vendor-vue'
          }
          if (id.includes('node_modules/chart.js/') || id.includes('node_modules/vue-chartjs/')) {
            return 'vendor-chart'
          }
          if (id.includes('node_modules/marked/')) {
            return 'vendor-marked'
          }
        }
      }
    }
  }
})
