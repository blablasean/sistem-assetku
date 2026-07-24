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
      <router-link to="/assets" class="nav-item" v-if="canAccessAdvanced">
        <span class="nav-icon">📦</span> Manajemen Aset
      </router-link>
      <router-link to="/workorders" class="nav-item">
        <span class="nav-icon">🔧</span> Work Order
      </router-link>
      <router-link to="/maintenance" class="nav-item" v-if="canAccessAdvanced">
        <span class="nav-icon">📅</span> Maintenance
      </router-link>
      <router-link to="/activitylogs" class="nav-item" v-if="canAccessAdvanced">
        <span class="nav-icon">📋</span> Activity Log
      </router-link>
      <router-link to="/users" class="nav-item" v-if="userRole === 'admin'">
        <span class="nav-icon">👥</span> User Management
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
          <router-link to="/assets" class="drawer-item" v-if="canAccessAdvanced">
            <span>📦</span> Manajemen Aset
          </router-link>
          <router-link to="/workorders" class="drawer-item">
            <span>🔧</span> Work Order
          </router-link>
          <router-link to="/maintenance" class="drawer-item" v-if="canAccessAdvanced">
            <span>📅</span> Maintenance
          </router-link>
          <router-link to="/activitylogs" class="drawer-item" v-if="canAccessAdvanced">
            <span>📋</span> Activity Log
          </router-link>
          <router-link to="/users" class="drawer-item" v-if="userRole === 'admin'">
            <span>👥</span> User Management
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
          <div class="user-avatar">
            <img v-if="userAvatar" :src="userAvatar" class="avatar-img-nav" />
            <span v-else>{{ userInitial }}</span>
          </div>
          <div class="user-info">
            <span class="user-name">{{ userName }}</span>
            <span class="user-role-sub">{{ roleLabel }}</span>
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
            <span>👤</span> Profil Saya & Foto
          </button>
          <button class="dropdown-item logout-item" @click="handleLogout">
            <span>🚪</span> Keluar (Logout)
          </button>
        </div>
      </div>
    </div>
    <!-- Custom User Profile Modal UI -->
    <ModalDialog :show="showProfileModal" title="👤 Foto & Profil Saya" maxWidth="450px" @close="showProfileModal = false">
      <form @submit.prevent="updateUserProfile" class="profile-card-content">
        <div class="avatar-edit-wrapper">
          <div class="profile-avatar-large">
            <img v-if="editProfileAvatar" :src="editProfileAvatar" class="avatar-img-large" />
            <span v-else>{{ userInitial }}</span>
          </div>
          <label class="upload-avatar-btn">
            📷 Upload Foto Profil
            <input type="file" accept="image/*" @change="onAvatarFileSelected" style="display:none" />
          </label>
        </div>

        <!-- Read-only Info Display -->
        <div class="profile-read-only-box">
          <div class="pro-info-row">
            <span>Nama Lengkap:</span>
            <strong>{{ userName }}</strong>
          </div>
          <div class="pro-info-row">
            <span>Role / Jabatan:</span>
            <span class="role-badge" :class="userRole">{{ roleLabel }}</span>
          </div>
          <div class="pro-info-row">
            <span>Status Akun:</span>
            <strong class="status-active-badge">🟢 Aktif</strong>
          </div>
          <p class="admin-only-notice">
            ℹ️ Perubahan nama & password hanya dapat dilakukan oleh Administrator dari menu User Management.
          </p>
        </div>

        <div v-if="profileMsg" class="profile-msg-banner" :class="profileMsgType">
          {{ profileMsg }}
        </div>

        <div class="profile-actions-row">
          <button type="button" class="action-btn-sharp cancel-btn" @click="showProfileModal = false">Batal</button>
          <button type="submit" class="action-btn-sharp primary-btn" :disabled="isUpdatingProfile">
            {{ isUpdatingProfile ? 'Menyimpan...' : 'Simpan Foto Profil' }}
          </button>
        </div>
      </form>
    </ModalDialog>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import ModalDialog from './ModalDialog.vue'
import api from '../api'

const router = useRouter()

const dropdownOpen = ref(false)
const mobileMenuOpen = ref(false)
const showProfileModal = ref(false)

const userName = ref(sessionStorage.getItem('user_name') || localStorage.getItem('user_name') || 'User Hotel')
const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')
const userAvatar = ref(sessionStorage.getItem('user_avatar') || localStorage.getItem('user_avatar') || '')

const editProfileName = ref(userName.value)
const editProfilePassword = ref('')
const editProfileAvatar = ref(userAvatar.value)
const isUpdatingProfile = ref(false)
const profileMsg = ref('')
const profileMsgType = ref('success')

const isLoggedIn = computed(() => !!(sessionStorage.getItem('token') || localStorage.getItem('token')))

// Staff-only roles: dept_* dan external — hanya bisa akses Dashboard & Work Order
const STAFF_ROLES = ['external', 'dept_akunting', 'dept_spa', 'dept_sales', 'dept_hr', 'dept_fb_kitchen', 'dept_fb_service', 'dept_housekeeping', 'dept_frontoffice']
const isStaffOnly = computed(() => STAFF_ROLES.includes(userRole.value))
// Roles yang dapat akses halaman advanced (Aset, Maintenance, Activity Log)
const canAccessAdvanced = computed(() => !isStaffOnly.value)

const userInitial = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : 'U'
})

const roleLabel = computed(() => {
  const map = {
    admin: 'Administrator',
    hod: 'HOD Engineer',
    management: 'Supervisor Engineer',
    engineer: 'Staff Engineer',
    dept_akunting: 'Departement Akunting',
    dept_spa: 'Departement Spa',
    dept_sales: 'Department Sales',
    dept_hr: 'Department Human Resources',
    dept_fb_kitchen: 'Department Food Beverage Kitchen',
    dept_fb_service: 'Department Food Beverage Service',
    dept_housekeeping: 'Department House Keeping',
    dept_frontoffice: 'Department Front Office',
    external: 'Staff Hotel'
  }
  return map[userRole.value] || 'User Hotel'
})

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function showProfile() {
  editProfileName.value = userName.value
  editProfilePassword.value = ''
  editProfileAvatar.value = userAvatar.value
  profileMsg.value = ''
  showProfileModal.value = true
}

function onAvatarFileSelected(e) {
  const file = e.target.files[0]
  if (!file) return
  if (file.size > 2 * 1024 * 1024) {
    alert('Ukuran foto melebihi 2MB. Silakan gunakan foto yang lebih kecil.')
    return
  }
  const reader = new FileReader()
  reader.onload = (evt) => {
    editProfileAvatar.value = evt.target.result
  }
  reader.readAsDataURL(file)
}

async function updateUserProfile() {
  isUpdatingProfile.value = true
  profileMsg.value = ''
  try {
    const res = await api.post('/auth/profile', {
      name: editProfileName.value,
      password: editProfilePassword.value,
      avatar: editProfileAvatar.value
    })
    userName.value = editProfileName.value
    userAvatar.value = editProfileAvatar.value
    sessionStorage.setItem('user_name', editProfileName.value)
    localStorage.setItem('user_name', editProfileName.value)
    sessionStorage.setItem('user_avatar', editProfileAvatar.value)
    localStorage.setItem('user_avatar', editProfileAvatar.value)

    profileMsg.value = 'Profil & Foto Anda berhasil diperbarui!'
    profileMsgType.value = 'success'
    setTimeout(() => {
      showProfileModal.value = false
    }, 1500)
  } catch (e) {
    profileMsg.value = 'Gagal menyimpan profil: ' + (e.response?.data?.message || e.message)
    profileMsgType.value = 'error'
  } finally {
    isUpdatingProfile.value = false
  }
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
  padding: 10px 20px;
  border-bottom: 2px solid #d97706;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
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
  width: 38px;
  height: 38px;
  border-radius: 2px !important;
  cursor: pointer;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-logo {
  font-size: 1.3rem;
  background: #1e293b;
  padding: 4px 8px;
  border-radius: 2px !important;
  border: 1px solid #334155;
}

.brand-title {
  display: block;
  font-size: 1.15rem;
  font-weight: 800;
  color: #f59e0b;
  line-height: 1.1;
  letter-spacing: 0.5px;
}

.brand-sub {
  font-size: 0.7rem;
  color: #cbd5e1;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.navbar-links {
  display: flex;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #cbd5e1;
  text-decoration: none;
  padding: 8px 12px;
  border-radius: 2px !important;
  font-size: 0.85rem;
  font-weight: 600;
  border: 1px solid transparent;
  transition: all 0.15s ease;
}

.nav-item:hover, .nav-item.router-link-active {
  background: #1e293b;
  color: #f59e0b;
  border-color: #334155;
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.qr-quick-btn {
  background: #2563eb;
  color: white;
  border: 1px solid #1d4ed8;
  padding: 8px 14px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.82rem;
  cursor: pointer;
  min-height: 38px;
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
  border-radius: 2px !important;
  cursor: pointer;
  height: 44px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: #d97706;
  color: #ffffff;
  border-radius: 2px !important;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.9rem;
  flex-shrink: 0;
}

.user-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  justify-content: center;
  text-align: left;
  line-height: 1.2;
}

.user-name {
  font-size: 0.84rem;
  font-weight: 700;
  color: #ffffff;
  margin: 0;
  padding: 0;
  text-align: left;
  white-space: nowrap;
}

.user-role-sub {
  font-size: 0.72rem;
  font-weight: 600;
  color: #f59e0b;
  margin: 2px 0 0 0;
  padding: 0;
  text-align: left;
  white-space: nowrap;
  display: block;
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
  top: calc(100% + 6px);
  width: 200px;
  background: #0f172a;
  border: 1px solid #334155;
  border-radius: 2px !important;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.4);
  padding: 4px 0;
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
  background: rgba(15, 23, 42, 0.75);
  backdrop-filter: blur(2px);
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
  border-right: 2px solid #d97706;
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
  gap: 6px;
}

.drawer-item {
  display: flex;
  align-items: center;
  gap: 10px;
  color: #cbd5e1;
  text-decoration: none;
  padding: 10px 12px;
  border-radius: 2px !important;
  font-weight: 600;
  font-size: 0.9rem;
  background: transparent;
  border: 1px solid transparent;
  width: 100%;
  text-align: left;
}

.drawer-item.router-link-active {
  background: #1e293b;
  color: #f59e0b;
  border-color: #334155;
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
  width: 64px;
  height: 64px;
  border-radius: 2px !important;
  background: #0f172a;
  border: 2px solid #d97706;
  color: #f59e0b;
  font-size: 1.8rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
}

.profile-card-content h3 {
  margin: 0 0 4px;
  font-size: 1.2rem;
  color: #0f172a;
  font-weight: 700;
}

.profile-role-badge {
  font-size: 0.8rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 2px !important;
  margin-bottom: 18px;
}

.profile-role-badge.admin, .profile-role-badge.hod {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
}

.profile-role-badge.management {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.profile-role-badge.engineer {
  background: #dcfce7;
  color: #15803d;
  border: 1px solid #bbf7d0;
}

.profile-role-badge.external {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.profile-details-grid {
  width: 100%;
  display: grid;
  gap: 8px;
  margin-bottom: 18px;
  text-align: left;
}

.pdetail-item {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  padding: 10px 14px;
  border-radius: 2px !important;
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
  background: #0f172a;
  color: white;
  border: 1px solid #1e293b;
  padding: 10px;
  border-radius: 2px !important;
  font-weight: 700;
  cursor: pointer;
}

.close-profile-btn {
  width: 100%;
}

.avatar-img-nav {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 2px !important;
}

.avatar-edit-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 16px;
}

.avatar-img-large {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 2px !important;
}

.upload-avatar-btn {
  margin-top: 8px;
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #0f172a;
  padding: 6px 12px;
  border-radius: 2px !important;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  display: inline-block;
}

.profile-form-sharp {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 12px;
  text-align: left;
  margin-bottom: 16px;
}

.profile-form-sharp label span {
  display: block;
  font-size: 0.82rem;
  font-weight: 700;
  color: #0f172a;
  margin-bottom: 4px;
}

.profile-msg-banner {
  width: 100%;
  padding: 8px 12px;
  border-radius: 2px !important;
  font-size: 0.82rem;
  font-weight: 700;
  margin-bottom: 12px;
}

.profile-msg-banner.success {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.profile-msg-banner.error {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fca5a5;
}

.profile-actions-row {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.profile-read-only-box {
  width: 100%;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  padding: 14px;
  margin-bottom: 16px;
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pro-info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.85rem;
}

.pro-info-row span {
  color: #64748b;
  font-weight: 600;
}

.pro-info-row strong {
  color: #0f172a;
  font-weight: 700;
}

.admin-only-notice {
  margin: 6px 0 0;
  font-size: 0.75rem;
  color: #64748b;
  font-style: italic;
  border-top: 1px dashed #cbd5e1;
  padding-top: 8px;
  line-height: 1.3;
}
</style>
