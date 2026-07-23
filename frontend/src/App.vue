<template>
  <div class="app-layout">
    <HeaderNavbar v-if="showNavbar" @open-qr-scanner="showQrScanner = true" />

    <main class="main-wrapper" :class="{ 'with-bottom-nav': showNavbar }">
      <router-view @open-qr-scanner="showQrScanner = true" />
    </main>

    <!-- Android Bottom Navigation Bar (Mobile / Android) -->
    <BottomNav v-if="showNavbar" @open-qr-scanner="showQrScanner = true" />

    <QrScannerModal :show="showQrScanner" @close="showQrScanner = false" />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import HeaderNavbar from './components/HeaderNavbar.vue'
import BottomNav from './components/BottomNav.vue'
import QrScannerModal from './components/QrScannerModal.vue'

const route = useRoute()
const showQrScanner = ref(false)

const showNavbar = computed(() => {
  return route.path !== '/login'
})
</script>

<style>
:root {
  font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  color-scheme: light;
}

body {
  margin: 0;
  padding: 0;
  background-color: #f8fafc;
  color: #0f172a;
  -webkit-font-smoothing: antialiased;
  -webkit-tap-highlight-color: transparent;
}

.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
  width: 100%;
  max-width: 100vw;
  overflow-x: hidden;
}

.main-wrapper {
  flex: 1;
  width: 100%;
  max-width: 100vw;
  overflow-x: hidden;
}

/* Adjust bottom margin on mobile devices for Android bottom navigation bar */
@media (max-width: 850px) {
  .main-wrapper.with-bottom-nav {
    padding-bottom: 74px;
  }

  .page-container, .dashboard-screen {
    padding: 16px 12px !important;
  }
}

* {
  box-sizing: border-box;
}
</style>