<template>
  <nav class="android-bottom-nav" v-if="isLoggedIn && isMobileView">
    <router-link to="/dashboard" class="bnav-item">
      <span class="bnav-icon">📊</span>
      <span class="bnav-label">Dashboard</span>
    </router-link>

    <router-link to="/assets" class="bnav-item" v-if="userRole !== 'external'">
      <span class="bnav-icon">📦</span>
      <span class="bnav-label">Aset</span>
    </router-link>

    <!-- Center Floating Action Button (FAB) for QR Scanner -->
    <div class="fab-wrapper" @click="$emit('open-qr-scanner')">
      <button class="fab-btn" title="Scan QR Code Aset">
        📱
      </button>
    </div>

    <router-link to="/workorders" class="bnav-item">
      <span class="bnav-icon">🔧</span>
      <span class="bnav-label">Work Order</span>
    </router-link>

    <router-link to="/maintenance" class="bnav-item" v-if="userRole !== 'external'">
      <span class="bnav-icon">📅</span>
      <span class="bnav-label">Maintenance</span>
    </router-link>

    <router-link to="/activitylogs" class="bnav-item" v-if="userRole === 'external'">
      <span class="bnav-icon">📋</span>
      <span class="bnav-label">Log</span>
    </router-link>
  </nav>
</template>

<script setup>
import { ref, computed } from 'vue'

defineEmits(['open-qr-scanner'])

const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')
const isLoggedIn = computed(() => !!(sessionStorage.getItem('token') || localStorage.getItem('token')))
const isMobileView = computed(() => true)
</script>

<style scoped>
.android-bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 60px;
  background: #0f172a;
  border-top: 2px solid #d97706;
  display: flex;
  align-items: center;
  justify-content: space-around;
  z-index: 250;
  box-shadow: 0 -4px 15px rgba(0, 0, 0, 0.4);
  padding-bottom: env(safe-area-inset-bottom, 0px);
}

.bnav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #cbd5e1;
  text-decoration: none;
  font-size: 0.7rem;
  font-weight: 600;
  gap: 2px;
  flex: 1;
  height: 100%;
  transition: color 0.15s ease;
}

.bnav-icon {
  font-size: 1.2rem;
}

.bnav-item.router-link-active {
  color: #f59e0b;
}

.fab-wrapper {
  position: relative;
  top: -10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fab-btn {
  width: 48px;
  height: 48px;
  border-radius: 2px !important;
  background: #d97706;
  color: white;
  border: 2px solid #0f172a;
  font-size: 1.4rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 12px rgba(217, 119, 6, 0.4);
  transition: transform 0.15s ease;
}

.fab-btn:active {
  transform: scale(0.95);
}

@media (min-width: 850px) {
  .android-bottom-nav {
    display: none;
  }
}
</style>
