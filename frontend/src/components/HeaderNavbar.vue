<template>
  <header class="header-navbar">
    <div class="navbar-left">
      <button class="mobile-hamburger" v-if="isLoggedIn" @click="mobileMenuOpen = !mobileMenuOpen" aria-label="Toggle navigation">
        {{ mobileMenuOpen ? '✕' : '☰' }}
      </button>

      <div class="navbar-brand">
        <div class="brand-logo">🏢</div>
        <div>
          <span class="brand-title">AsetKu</span>
          <span class="brand-sub">Hotel Asset & WO</span>
        </div>
      </div>
    </div>

    <!-- Desktop Navigation Links -->
    <nav class="navbar-links desktop-only" v-if="isLoggedIn">
      <router-link to="/dashboard" class="nav-item">
        <span class="nav-icon">📊</span> Dashboard
      </router-link>
      <router-link to="/assets" class="nav-item" v-if="userRole !== 'external'">
        <span class="nav-icon">📦</span> Manajemen Aset
      </router-link>
      <router-link to="/workorders" class="nav-item">
        <span class="nav-icon">🔧</span> Work Order
      </router-link>
      <router-link to="/maintenance" class="nav-item" v-if="userRole !== 'external'">
        <span class="nav-icon">📅</span> Maintenance
      </router-link>
      <router-link to="/activitylogs" class="nav-item">
        <span class="nav-icon">📋</span> Activity Log
      </router-link>
    </nav>

    <!-- Mobile Drawer Overlay -->
    <div class="mobile-drawer-overlay" v-if="mobileMenuOpen && isLoggedIn" @click="mobileMenuOpen = false">
      <div class="mobile-drawer" @click.stop>
        <div class="drawer-header">
          <div class="user-avatar">{{ userInitial }}</div>
          <div>
            <p class="drawer-name">{{ userName }}</p>
            <span class="role-badge" :class="userRole">{{ roleLabel }}</span>
          </div>
          <button class="close-drawer" @click="mobileMenuOpen = false">✕</button>
        </div>

        <nav class="drawer-links" @click="mobileMenuOpen = false">
          <router-link to="/dashboard" class="drawer-item">
            <span>📊</span> Dashboard
          </router-link>
          <router-link to="/assets" class="drawer-item" v-if="userRole !== 'external'">
            <span>📦</span> Manajemen Aset
          </router-link>
          <router-link to="/workorders" class="drawer-item">
            <span>🔧</span> Work Order
          </router-link>
          <router-link to="/maintenance" class="drawer-item" v-if="userRole !== 'external'">
            <span>📅</span> Maintenance
          </router-link>
          <router-link to="/activitylogs" class="drawer-item">
            <span>📋</span> Activity Log
          </router-link>

          <div class="drawer-divider"></div>
          <button class="drawer-item logout" @click="handleLogout">
            <span>🚪</span> Keluar (Logout)
          </button>
        </nav>
      </div>
    </div>

    <!-- User Menu & Scan Button -->
    <div class="navbar-user" v-if="isLoggedIn">
      <button class="qr-quick-btn" @click="$emit('open-qr-scanner')" title="Scan QR Code Aset">
        📱 <span class="btn-text-desktop">Scan QR</span>
      </button>

      <!-- Account Dropdown (Desktop) -->
      <div class="account-dropdown-wrapper desktop-only">
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
    <!-- Custom User Profile Modal UI -->
    <ModalDialog :show="showProfileModal" title="👤 Profil Pengguna" maxWidth="450px" @close="showProfileModal = false">
      <div class="profile-card-content">
        <div class="profile-avatar-large">{{ userInitial }}</div>
        <h3>{{ userName }}</h3>
        <p class="profile-role-badge" :class="userRole">{{ roleLabel }}</p>
        
        <div class="profile-details-grid">
          <div class="pdetail-item">
            <span>Role / Jabatan</span>
            <strong>{{ roleLabel }}</strong>
          </div>
          <div class="pdetail-item">
            <span>Status Sesi</span>
            <strong class="status-active-badge">🟢 Aktif</strong>
          </div>
          <div class="pdetail-item">
            <span>Hapus Work Order</span>
            <strong>{{ (userRole === 'hod' || userRole === 'admin' || userRole === 'management') ? 'Diizinkan' : 'Tidak' }}</strong>
          </div>
          <div class="pdetail-item">
            <span>Kontrol Sistem</span>
            <strong>{{ userRole === 'admin' ? 'Akses Penuh' : 'Akses Standar' }}</strong>
          </div>
        </div>

        <button class="submit-modal-btn close-profile-btn" @click="showProfileModal = false">Tutup Profil</button>
      </div>
    </ModalDialog>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import ModalDialog from './ModalDialog.vue'

const router = useRouter()

const dropdownOpen = ref(false)
const mobileMenuOpen = ref(false)
const showProfileModal = ref(false)

const userName = ref(sessionStorage.getItem('user_name') || localStorage.getItem('user_name') || 'User Hotel')
const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')

const isLoggedIn = computed(() => !!(sessionStorage.getItem('token') || localStorage.getItem('token')))

const userInitial = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : 'U'
})

const roleLabel = computed(() => {
  const map = {
    admin: 'Administrator',
    hod: 'HOD Engineer',
    management: 'Supervisor',
    engineer: 'Staff Engineer',
    external: 'Staff Hotel'
  }
  return map[userRole.value] || 'User Hotel'
})

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function showProfile() {
  showProfileModal.value = true
}

function handleLogout() {
  sessionStorage.clear()
  localStorage.clear()
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
  padding: 10px 16px;
  border-bottom: 1px solid #1e293b;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mobile-hamburger {
  display: none;
  background: #1e293b;
  border: 1px solid #334155;
  color: white;
  font-size: 1.3rem;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  cursor: pointer;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-logo {
  font-size: 1.5rem;
  background: #1e293b;
  padding: 4px 8px;
  border-radius: 10px;
}

.brand-title {
  display: block;
  font-size: 1.15rem;
  font-weight: 800;
  color: #38bdf8;
  line-height: 1.1;
}

.brand-sub {
  font-size: 0.7rem;
  color: #94a3b8;
}

.navbar-links {
  display: flex;
  gap: 6px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #cbd5e1;
  text-decoration: none;
  padding: 8px 12px;
  border-radius: 10px;
  font-size: 0.85rem;
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
  gap: 12px;
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
  min-height: 40px;
}

.account-dropdown-wrapper {
  position: relative;
}

.account-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #1e293b;
  border: 1px solid #334155;
  color: white;
  padding: 5px 10px;
  border-radius: 10px;
  cursor: pointer;
}

.user-avatar {
  width: 30px;
  height: 30px;
  background: #38bdf8;
  color: #0f172a;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.9rem;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}

.user-name {
  font-size: 0.8rem;
  font-weight: 600;
}

.role-badge {
  font-size: 0.65rem;
  padding: 1px 5px;
  border-radius: 4px;
  font-weight: 600;
  text-transform: uppercase;
}

.role-badge.management { background: #ea580c; color: white; }
.role-badge.hod { background: #2563eb; color: white; }
.role-badge.engineer { background: #059669; color: white; }
.role-badge.external { background: #64748b; color: white; }

.dropdown-arrow {
  font-size: 0.6rem;
  color: #94a3b8;
}

.dropdown-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  width: 200px;
  background: #1e293b;
  border: 1px solid #334155;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.3);
  padding: 6px 0;
  z-index: 200;
}

.dropdown-header {
  padding: 8px 14px;
}

.dh-title {
  margin: 0;
  font-weight: 700;
  color: white;
  font-size: 0.85rem;
}

.dh-sub {
  margin: 2px 0 0;
  color: #94a3b8;
  font-size: 0.7rem;
}

.dropdown-divider {
  height: 1px;
  background: #334155;
  margin: 4px 0;
}

.dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  background: transparent;
  border: none;
  color: #cbd5e1;
  font-size: 0.85rem;
  text-align: left;
  cursor: pointer;
}

.logout-item { color: #f87171; }

/* Mobile Drawer Overlay */
.mobile-drawer-overlay {
  position: fixed;
  inset: 0;
  background: rgba(15, 23, 42, 0.7);
  backdrop-filter: blur(4px);
  z-index: 300;
  display: flex;
}

.mobile-drawer {
  width: 280px;
  background: #0f172a;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
  border-right: 1px solid #1e293b;
}

.drawer-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid #1e293b;
}

.drawer-name {
  margin: 0 0 4px;
  font-weight: 700;
  color: white;
  font-size: 0.95rem;
}

.close-drawer {
  margin-left: auto;
  background: transparent;
  border: none;
  color: #94a3b8;
  font-size: 1.2rem;
}

.drawer-links {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.drawer-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #cbd5e1;
  text-decoration: none;
  padding: 12px 14px;
  border-radius: 10px;
  font-weight: 600;
  font-size: 0.95rem;
  background: transparent;
  border: none;
  width: 100%;
  text-align: left;
}

.drawer-item.router-link-active {
  background: #1e293b;
  color: #38bdf8;
}

.drawer-divider {
  height: 1px;
  background: #1e293b;
  margin: 12px 0;
}

.drawer-item.logout {
  color: #f87171;
}

/* Responsive Queries for Mobile (Android & iOS) */
@media (max-width: 850px) {
  .desktop-only {
    display: none !important;
  }

  .mobile-hamburger {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

/* User Profile Modal Styling */
.profile-card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  padding: 10px 0;
}

.profile-avatar-large {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  color: white;
  font-size: 2rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  box-shadow: 0 8px 20px rgba(37, 99, 235, 0.3);
}

.profile-card-content h3 {
  margin: 0 0 4px;
  font-size: 1.3rem;
  color: #0f172a;
}

.profile-role-badge {
  font-size: 0.82rem;
  font-weight: 700;
  padding: 4px 12px;
  border-radius: 999px;
  margin-bottom: 20px;
}

.profile-role-badge.admin, .profile-role-badge.hod {
  background: #fef3c7;
  color: #b45309;
}

.profile-role-badge.management {
  background: #dbeafe;
  color: #1e40af;
}

.profile-role-badge.engineer {
  background: #dcfce7;
  color: #15803d;
}

.profile-role-badge.external {
  background: #f1f5f9;
  color: #475569;
}

.profile-details-grid {
  width: 100%;
  display: grid;
  gap: 10px;
  margin-bottom: 20px;
  text-align: left;
}

.pdetail-item {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  padding: 10px 14px;
  border-radius: 10px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pdetail-item span {
  font-size: 0.8rem;
  color: #64748b;
}

.pdetail-item strong {
  font-size: 0.85rem;
  color: #0f172a;
}

.status-active-badge {
  color: #16a34a !important;
}

.submit-modal-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
}

.close-profile-btn {
  width: 100%;
}
</style>
