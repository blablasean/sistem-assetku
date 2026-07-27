import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/auth': 'http://localhost:8080',
      '/assets': 'http://localhost:8080',
      '/workorders': 'http://localhost:8080',
      '/maintenance': 'http://localhost:8080',
      '/users': 'http://localhost:8080',
      '/activitylogs': 'http://localhost:8080',
      '/uploads': 'http://localhost:8080'
    }
  }
})