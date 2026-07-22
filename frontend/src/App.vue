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
  font-family: Inter, system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
  color-scheme: light;
}

body {
  margin: 0;
  padding: 0;
  background-color: #f8fafc;
  color: #0f172a;
  -webkit-font-smoothing: antialiased;
}

.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-wrapper {
  flex: 1;
}

* {
  box-sizing: border-box;
}
</style>