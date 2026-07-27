<template>
  <div class="app-layout">
    <HeaderNavbar v-if="showNavbar" @open-qr-scanner="showQrScanner = true" />

    <main class="main-wrapper">
      <router-view @open-qr-scanner="showQrScanner = true" />
    </main>

    <QrScannerModal :show="showQrScanner" @close="showQrScanner = false" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import HeaderNavbar from './components/HeaderNavbar.vue'
import QrScannerModal from './components/QrScannerModal.vue'

const route = useRoute()
const showQrScanner = ref(false)

const showNavbar = computed(() => {
  return route.path !== '/login'
})
</script>

<style>
:root {
  font-family: -apple-system, BlinkMacSystemFont, "SF Pro Display", "SF Pro Text", "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
  color-scheme: light;
  --bg-page: #f8fafc;
  --bg-card: #ffffff;
  --text-main: #0f172a;
  --text-muted: #64748b;
  --border-color: #e2e8f0;
  --accent-primary: #2563eb;
  --accent-success: #16a34a;
  --accent-warning: #d97706;
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--bg-page);
  color: var(--text-main);
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  -webkit-tap-highlight-color: transparent;
  letter-spacing: -0.01em;
}

/* Minimal Bold Sharp Design Tokens (Apple / Microsoft Style) */
button,
input,
select,
textarea,
.card,
.modal-container,
.modal-content,
.btn,
.badge,
.status-badge,
.search-bar,
.summary-card,
.stat-card,
.data-table {
  border-radius: 4px !important;
}

.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  width: 100%;
  max-width: 100vw;
  overflow-x: hidden;
  background-color: var(--bg-page);
}

.main-wrapper {
  flex: 1;
  width: 100%;
  max-width: 100vw;
  overflow-x: hidden;
}

/* Crisp Clean Buttons */
button, .btn {
  font-family: inherit;
  font-weight: 600;
  letter-spacing: -0.01em;
  transition: all 0.15s ease-in-out;
}

button:active, .btn:active {
  transform: scale(0.98);
}

@media (max-width: 640px) {
  .page-container, .dashboard-screen {
    padding: 12px 10px !important;
  }

  /* Full width mobile modal adjustments */
  .modal-container, .modal-content {
    width: 96% !important;
    max-width: 96% !important;
    margin: 10px auto !important;
    padding: 16px 12px !important;
    border-radius: 12px !important;
  }

  /* Touch scrolling utilities */
  .touch-scroll-x {
    display: flex;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
    gap: 8px;
    padding-bottom: 4px;
  }
  .touch-scroll-x::-webkit-scrollbar {
    display: none;
  }

  /* Mobile action button grid */
  .header-action-group {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    width: 100%;
  }
  .header-action-group button, .header-action-group .primary-btn {
    flex: 1 1 calc(50% - 4px);
    justify-content: center;
    font-size: 0.82rem !important;
    padding: 10px 8px !important;
  }
}

* {
  box-sizing: border-box;
}
</style>