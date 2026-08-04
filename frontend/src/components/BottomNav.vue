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

const userRole = ref(sessionStorage.getItem('user_role') || 'external')
const isLoggedIn = computed(() => !!sessionStorage.getItem('token'))
const isMobileView = computed(() => true)
</script>

<style scoped>
.android-bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 62px;
  background: rgba(15, 23, 42, 0.92);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: space-around;
  z-index: 250;
  box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.3);
  padding-bottom: env(safe-area-inset-bottom, 0px);
}

.bnav-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #94a3b8;
  text-decoration: none;
  font-size: 0.72rem;
  font-weight: 700;
  gap: 2px;
  flex: 1;
  height: 100%;
  transition: all 0.2s ease;
}

.bnav-icon {
  font-size: 1.25rem;
  transition: transform 0.2s ease;
}

.bnav-item.router-link-active {
  color: #38bdf8;
}

.bnav-item.router-link-active .bnav-icon {
  transform: translateY(-2px);
}

.fab-wrapper {
  position: relative;
  top: -12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fab-btn {
  width: 52px;
  height: 52px;
  border-radius: 50% !important;
  background: linear-gradient(135deg, #007aff, #0284c7);
  color: white;
  border: 3px solid #0f172a;
  font-size: 1.4rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 8px 20px rgba(0, 122, 255, 0.45);
  transition: transform 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.fab-btn:active {
  transform: scale(0.92);
}

@media (min-width: 850px) {
  .android-bottom-nav {
    display: none;
  }
}
</style>
