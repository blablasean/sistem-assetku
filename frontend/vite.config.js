import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Helper to bypass dev server API proxy for HTML page navigations & static asset files
const bypassHtml = (req) => {
  if (req.url && req.url.match(/\.(png|jpg|jpeg|gif|svg|ico|css|js|webp)(\?.*)?$/i)) {
    return req.url
  }
  if (req.headers.accept && req.headers.accept.includes('text/html')) {
    return '/index.html'
  }
}

export default defineConfig({
  plugins: [
    vue()
  ],
  server: {
    port: 3000,
    host: true,
    proxy: {
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        bypass: bypassHtml
      },
      '/assets': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        bypass: bypassHtml
      },
      '/workorders': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        bypass: bypassHtml
      },
      '/maintenance': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        bypass: bypassHtml
      },
      '/users': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        bypass: bypassHtml
      },
      '/activitylogs': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        bypass: bypassHtml
      },
      '/uploads': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})