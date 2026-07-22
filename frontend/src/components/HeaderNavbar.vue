<template>
  <header class="header-navbar">
    <div class="navbar-brand">
      <div class="brand-logo">🏢</div>
      <div>
        <span class="brand-title">AsetKu</span>
        <span class="brand-sub">Hotel Asset & WO Management</span>
      </div>
    </div>

    <nav class="navbar-links" v-if="isLoggedIn">
      <router-link to="/dashboard" class="nav-item">
        <span class="nav-icon">📊</span> Dashboard
      </router-link>
      <router-link to="/assets" class="nav-item">
        <span class="nav-icon">📦</span> Manajemen Aset
      </router-link>
      <router-link to="/workorders" class="nav-item">
        <span class="nav-icon">🔧</span> Work Order
      </router-link>
      <router-link to="/maintenance" class="nav-item">
        <span class="nav-icon">📅</span> Maintenance
      </router-link>
      <router-link to="/utility" class="nav-item">
        <span class="nav-icon">⚡</span> Utility
      </router-link>
      <router-link to="/activitylogs" class="nav-item" v-if="userRole === 'management' || userRole === 'hod'">
        <span class="nav-icon">📋</span> Activity Log
      </router-link>
    </nav>

    <div class="navbar-user" v-if="isLoggedIn">
      <button class="qr-quick-btn" @click="$emit('open-qr-scanner')" title="Scan QR Code Aset">
        📱 Scan QR
      </button>

      <!-- Account Dropdown -->
      <div class="account-dropdown-wrapper">
        <button class="account-trigger" @click="toggleDropdown">
          <div class="user-avatar">{{ userInitial }}</div>
          <div class="user-info">
            <span class="user-name">{{ userName }}</span>
            <span class="role-badge" :class="userRole">{{ roleLabel }}</span>
          </div>
          <span class="dropdown-arrow">▼</span>
        </button>

        <div class="dropdown-menu" v-if="dropdownOpen" @click="dropdownOpen = false">
          <div class="dropdown-header">
            <p class="dh-title">{{ userName }}</p>
            <p class="dh-sub">Role: {{ roleLabel }}</p>
          </div>
          <div class="dropdown-divider"></div>
          <button class="dropdown-item" @click="showProfile">
            <span>👤</span> Profil Saya
          </button>

          <button class="dropdown-item logout-item" @click="handleLogout">
            <span>🚪</span> Keluar (Logout)
          </button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const dropdownOpen = ref(false)
const userName = ref(localStorage.getItem('user_name') || 'User Hotel')
const userRole = ref(localStorage.getItem('user_role') || 'external')

const isLoggedIn = computed(() => !!localStorage.getItem('token'))

const userInitial = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : 'U'
})

const roleLabel = computed(() => {
  const map = {
    management: 'Management Engineer',
    hod: 'HOD Engineer',
    engineer: 'Staff Engineer',
    external: 'External User (Staff)'
  }
  return map[userRole.value] || 'User Hotel'
})

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function showProfile() {
  alert(`Detail Pengguna:\n\nNama: ${userName.value}\nRole: ${roleLabel.value}\nSistem: Sistem AsetKu Hotel`)
}

function handleLogout() {
  localStorage.removeItem('token')
  localStorage.removeItem('user_role')
  localStorage.removeItem('user_name')
  router.push('/login')
}

function handleClickOutside(e) {
  if (!e.target.closest('.account-dropdown-wrapper')) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.header-navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: #0f172a;
  color: #ffffff;
  padding: 12px 24px;
  border-bottom: 1px solid #1e293b;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-logo {
  font-size: 1.8rem;
  background: #1e293b;
  padding: 6px 10px;
  border-radius: 12px;
}

.brand-title {
  display: block;
  font-size: 1.25rem;
  font-weight: 800;
  letter-spacing: -0.5px;
  color: #38bdf8;
}

.brand-sub {
  font-size: 0.75rem;
  color: #94a3b8;
}

.navbar-links {
  display: flex;
  gap: 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #cbd5e1;
  text-decoration: none;
  padding: 8px 14px;
  border-radius: 10px;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.nav-item:hover, .nav-item.router-link-active {
  background: #1e293b;
  color: #38bdf8;
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 16px;
}

.qr-quick-btn {
  background: linear-gradient(135deg, #0284c7, #0369a1);
  color: white;
  border: none;
  padding: 8px 14px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.qr-quick-btn:hover {
  transform: scale(1.03);
}

.account-dropdown-wrapper {
  position: relative;
}

.account-trigger {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #1e293b;
  border: 1px solid #334155;
  color: white;
  padding: 6px 12px;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s;
}

.account-trigger:hover {
  background: #334155;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: #38bdf8;
  color: #0f172a;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.95rem;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.user-name {
  font-size: 0.85rem;
  font-weight: 600;
}

.role-badge {
  font-size: 0.7rem;
  padding: 1px 6px;
  border-radius: 4px;
  font-weight: 600;
  text-transform: uppercase;
}

.role-badge.management { background: #ea580c; color: white; }
.role-badge.hod { background: #2563eb; color: white; }
.role-badge.engineer { background: #059669; color: white; }
.role-badge.external { background: #64748b; color: white; }

.dropdown-arrow {
  font-size: 0.65rem;
  color: #94a3b8;
}

.dropdown-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  width: 220px;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 14px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  padding: 8px 0;
  z-index: 200;
}

.dropdown-header {
  padding: 10px 16px;
}

.dh-title {
  margin: 0;
  font-weight: 700;
  color: white;
  font-size: 0.9rem;
}

.dh-sub {
  margin: 2px 0 0;
  color: #94a3b8;
  font-size: 0.75rem;
}

.dropdown-divider {
  height: 1px;
  background: #334155;
  margin: 6px 0;
}

.dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 16px;
  background: transparent;
  border: none;
  color: #cbd5e1;
  font-size: 0.85rem;
  text-align: left;
  cursor: pointer;
  transition: background 0.15s;
}

.dropdown-item:hover {
  background: #334155;
  color: white;
}

.logout-item {
  color: #f87171;
}

.logout-item:hover {
  background: #451a1a;
  color: #fca5a5;
}
</style>
