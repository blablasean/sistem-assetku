<template>
  <header class="header-navbar">
    <div class="navbar-container">
      <div class="navbar-left">
        <button class="mobile-hamburger" v-if="isLoggedIn" @click="mobileMenuOpen = !mobileMenuOpen" aria-label="Toggle navigation">
          {{ mobileMenuOpen ? '✕' : '☰' }}
        </button>

        <div class="navbar-brand">
          <div class="brand-logo">
            <img src="/assets/logo.png" alt="AsetKu Logo" class="brand-logo-img" @error="logoFailed = true" v-if="!logoFailed" />
            <span v-else class="brand-text-logo">A</span>
          </div>
          <div>
            <span class="brand-title">AsetKu</span>
            <span class="brand-sub">Asset & Work Order</span>
          </div>
        </div>
      </div>

      <!-- Desktop Navigation Links -->
      <nav class="navbar-links desktop-only" v-if="isLoggedIn">
        <router-link to="/dashboard" class="nav-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="nav-svg"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg> Dashboard
        </router-link>
        <router-link to="/assets" class="nav-item" v-if="canAccessAdvanced">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="nav-svg"><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg> Manajemen Aset
        </router-link>
        <router-link to="/workorders" class="nav-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="nav-svg"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg> Work Order
        </router-link>
        <router-link to="/maintenance" class="nav-item" v-if="canAccessAdvanced">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="nav-svg"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg> Maintenance
        </router-link>
        <router-link to="/activitylogs" class="nav-item" v-if="canAccessAdvanced">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="nav-svg"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg> Activity Log
        </router-link>
        <router-link to="/users" class="nav-item" v-if="userRole === 'admin'">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="nav-svg"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg> User Management
        </router-link>
      </nav>

      <!-- Mobile Drawer Overlay -->
      <div class="mobile-drawer-overlay" v-if="mobileMenuOpen && isLoggedIn" @click="mobileMenuOpen = false">
        <div class="mobile-drawer" @click.stop>
          <div class="drawer-top-bar">
            <span class="drawer-menu-title">Menu Navigasi</span>
            <button class="close-drawer" @click.stop="mobileMenuOpen = false">✕</button>
          </div>

          <nav class="drawer-links">
            <router-link to="/dashboard" class="drawer-item" @click="mobileMenuOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg> Dashboard
            </router-link>
            <router-link to="/assets" class="drawer-item" v-if="canAccessAdvanced" @click="mobileMenuOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg> Manajemen Aset
            </router-link>
            <router-link to="/workorders" class="drawer-item" @click="mobileMenuOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg> Work Order
            </router-link>
            <router-link to="/maintenance" class="drawer-item" v-if="canAccessAdvanced" @click="mobileMenuOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg> Maintenance
            </router-link>
            <router-link to="/activitylogs" class="drawer-item" v-if="canAccessAdvanced" @click="mobileMenuOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg> Activity Log
            </router-link>
            <router-link to="/users" class="drawer-item" v-if="userRole === 'admin'" @click="mobileMenuOpen = false">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 1 0 7.75"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg> User Management
            </router-link>

            <div class="drawer-divider"></div>

            <button class="drawer-item logout" @click="handleLogout">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/></svg> Keluar (Logout)
            </button>
          </nav>
        </div>
      </div>

      <!-- User Menu & Scan Button -->
      <div class="navbar-user" v-if="isLoggedIn">
        <button class="qr-quick-btn desktop-only" @click="$emit('open-qr-scanner')" title="Scan QR Code Aset">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-svg"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/><path d="M3 12h.01"/><path d="M12 3h.01"/><path d="M12 16v.01"/><path d="M16 12h1"/><path d="M21 12v.01"/><path d="M12 21v-1"/></svg>
        </button>

        <!-- Account Dropdown (Desktop) -->
        <div class="account-dropdown-wrapper desktop-only">
          <button class="account-trigger" @click="toggleDropdown">
            <div class="user-avatar">
              <img v-if="userAvatar" :src="userAvatar" class="avatar-img-nav" />
              <span v-else>{{ userInitial }}</span>
            </div>
            <div class="user-info">
              <span class="user-name" :title="userName">{{ shortUserName }}</span>
              <span class="user-role-sub">{{ roleLabel }}</span>
            </div>
            <span class="dropdown-arrow">▼</span>
          </button>

          <div class="dropdown-menu" v-if="dropdownOpen" @click="dropdownOpen = false">
            <div class="dropdown-header">
              <div class="dh-avatar">
                <img v-if="userAvatar" :src="userAvatar" class="avatar-img-nav" />
                <span v-else>{{ userInitial }}</span>
              </div>
              <div class="dh-info">
                <p class="dh-title" :title="userName">{{ shortUserName }}</p>
                <span class="dh-role-badge" :class="userRole">{{ roleLabel }}</span>
              </div>
            </div>
            <div class="dropdown-divider"></div>
            <button class="dropdown-item" @click="showProfile">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="dd-icon"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
              <span>Profil Saya & Foto</span>
            </button>
            <button class="dropdown-item logout-item" @click="handleLogout">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="dd-icon"><path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" x2="9" y1="12" y2="12"/></svg>
              <span>Keluar (Logout)</span>
            </button>
          </div>
        </div>

        <!-- Profile Avatar Trigger Button (Mobile Only, Top Right Header) -->
        <button class="mobile-profile-avatar-btn mobile-only" @click="showProfile" title="Profil Saya & Foto">
          <div class="user-avatar">
            <img v-if="userAvatar" :src="userAvatar" class="avatar-img-nav" />
            <span v-else>{{ userInitial }}</span>
          </div>
        </button>
      </div>
    </div>

    <!-- Custom User Profile Modal UI -->
    <ModalDialog :show="showProfileModal" title="Profil Saya & Foto" maxWidth="440px" @close="showProfileModal = false">
      <form @submit.prevent="updateUserProfile" class="profile-card-content">
        <div class="avatar-edit-wrapper">
          <div class="profile-avatar-large">
            <img v-if="editProfileAvatar" :src="editProfileAvatar" class="avatar-img-large" />
            <span v-else>{{ userInitial }}</span>
          </div>
          <label class="upload-avatar-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
            <span>Upload Foto Profil</span>
            <input type="file" accept="image/*" @change="onAvatarFileSelected" style="display:none" />
          </label>
        </div>

        <!-- Read-only Info Display -->
        <div class="profile-read-only-box">
          <div class="pro-info-row">
            <span>Nama Lengkap</span>
            <strong>{{ userName }}</strong>
          </div>
          <div class="pro-info-row">
            <span>Role / Jabatan</span>
            <span class="role-badge" :class="userRole">{{ roleLabel }}</span>
          </div>
          <p class="admin-only-notice">
            ℹ️ Perubahan nama & password hanya dapat dilakukan oleh Administrator dari menu User Management.
          </p>
        </div>

        <div v-if="profileMsg" class="profile-msg-banner" :class="profileMsgType">
          {{ profileMsg }}
        </div>

        <div class="profile-actions-row">
          <button type="button" class="ios-btn ios-btn-cancel" @click="showProfileModal = false">Batal</button>
          <button type="submit" class="ios-btn ios-btn-primary" :disabled="isUpdatingProfile">
            {{ isUpdatingProfile ? 'Menyimpan...' : 'Simpan Foto Profil' }}
          </button>
        </div>
      </form>
    </ModalDialog>

    <!-- Mobile Bottom Dock (iOS App Style Navigation) -->
    <nav class="mobile-bottom-dock mobile-only" v-if="isLoggedIn">
      <!-- Staff / External Role Dock Layout (3 Items: Dashboard | Scan QR (Center) | Work Order) -->
      <template v-if="isStaffOnly">
        <router-link to="/dashboard" class="dock-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg>
          <span>Dashboard</span>
        </router-link>

        <button class="dock-fab-scan" @click="$emit('open-qr-scanner')" title="Scan QR Code Aset">
          <div class="fab-icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/></svg>
          </div>
          <span>Scan QR</span>
        </button>

        <router-link to="/workorders" class="dock-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
          <span>Work Order</span>
        </router-link>
      </template>

      <!-- Advanced Role Dock Layout (5 Items: Dashboard | Work Order | Scan QR (Center) | Aset | Maintenance) -->
      <template v-else>
        <router-link to="/dashboard" class="dock-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="9" x="3" y="3" rx="1"/><rect width="7" height="5" x="14" y="3" rx="1"/><rect width="7" height="9" x="14" y="12" rx="1"/><rect width="7" height="5" x="3" y="16" rx="1"/></svg>
          <span>Dashboard</span>
        </router-link>

        <router-link to="/workorders" class="dock-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
          <span>Work Order</span>
        </router-link>

        <button class="dock-fab-scan" @click="$emit('open-qr-scanner')" title="Scan QR Code Aset">
          <div class="fab-icon-wrapper">
            <svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/></svg>
          </div>
          <span>Scan QR</span>
        </button>

        <router-link to="/assets" class="dock-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 22V12"/></svg>
          <span>Aset</span>
        </router-link>

        <router-link to="/maintenance" class="dock-item">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" x2="16" y1="2" y2="6"/><line x1="8" x2="8" y1="2" y2="6"/><line x1="3" x2="21" y1="10" y2="10"/></svg>
          <span>Maintenance</span>
        </router-link>
      </template>
    </nav>
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

const userName = ref(sessionStorage.getItem('user_name') || 'Pengguna System')
const userRole = ref(sessionStorage.getItem('user_role') || 'external')
const userAvatar = ref(sessionStorage.getItem('user_avatar') || '')

const editProfileName = ref(userName.value)
const editProfilePassword = ref('')
const editProfileAvatar = ref(userAvatar.value)
const isUpdatingProfile = ref(false)
const profileMsg = ref('')
const profileMsgType = ref('success')
const logoFailed = ref(false)
const avatarFailed = ref(false)

function onAvatarImgError() {
  avatarFailed.value = true
}

const isLoggedIn = computed(() => !!sessionStorage.getItem('token'))

// Staff-only roles: dept_* dan external — hanya bisa akses Dashboard & Work Order
const STAFF_ROLES = ['external', 'dept_akunting', 'dept_spa', 'dept_sales', 'dept_hr', 'dept_fb_kitchen', 'dept_fb_service', 'dept_housekeeping', 'dept_frontoffice']
const isStaffOnly = computed(() => STAFF_ROLES.includes(userRole.value))
// Roles yang dapat akses halaman advanced (Aset, Maintenance, Activity Log)
const canAccessAdvanced = computed(() => !isStaffOnly.value)

const userInitial = computed(() => {
  return userName.value ? userName.value.charAt(0).toUpperCase() : 'U'
})

const shortUserName = computed(() => {
  if (!userName.value) return 'User'
  const name = userName.value.trim()
  const parts = name.split(/\s+/)
  if (parts.length > 2) {
    return `${parts[0]} ${parts[1]}`
  }
  return name
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
    external: 'Staff Operasional'
  }
  return map[userRole.value] || 'Pengguna System'
})

function syncUserFromStorage() {
  userName.value = sessionStorage.getItem('user_name') || 'Pengguna System'
  userRole.value = sessionStorage.getItem('user_role') || 'external'
  userAvatar.value = sessionStorage.getItem('user_avatar') || ''
}

function toggleDropdown() {
  syncUserFromStorage()
  dropdownOpen.value = !dropdownOpen.value
}

function showProfileFromMobile() {
  mobileMenuOpen.value = false
  showProfile()
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
    sessionStorage.setItem('user_avatar', editProfileAvatar.value)

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

async function handleLogout() {
  try {
    await api.post('/auth/logout').catch(() => null)
  } finally {
    sessionStorage.clear()
    localStorage.clear()
    router.push('/login')
  }
}

function handleClickOutside(e) {
  if (!e.target.closest('.account-dropdown-wrapper')) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  syncUserFromStorage()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.header-navbar {
  background: #ffffff;
  color: #0f172a;
  border-bottom: 1px solid #e2e8f0;
  position: sticky;
  top: 0;
  z-index: 100;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.05);
  width: 100%;
}

.navbar-container {
  max-width: 1440px;
  margin: 0 auto;
  padding: 0 24px;
  height: 60px;
  min-height: 60px;
  max-height: 60px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  overflow: visible;
}

.navbar-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.mobile-hamburger {
  display: none;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  color: #0f172a;
  font-size: 1.2rem;
  width: 38px;
  height: 38px;
  border-radius: 4px !important;
  cursor: pointer;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
}

.brand-logo {
  font-size: 1.2rem;
  background: #f8fafc;
  padding: 4px;
  border-radius: 4px !important;
  border: 1px solid #e2e8f0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
}

.brand-logo-img {
  width: 26px;
  height: 26px;
  object-fit: contain;
}

.brand-text-logo {
  font-weight: 900;
  font-size: 1.1rem;
  color: #0f172a;
}

.brand-title {
  display: block;
  font-size: 1.15rem;
  font-weight: 800;
  color: #0f172a;
  line-height: 1.1;
  letter-spacing: -0.02em;
}

.brand-sub {
  font-size: 0.72rem;
  color: #64748b;
  font-weight: 600;
  letter-spacing: 0.5px;
}

.navbar-links {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
  flex-wrap: nowrap;
}

.nav-item {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  color: #475569;
  text-decoration: none;
  padding: 7px 10px;
  border-radius: 4px !important;
  font-size: 0.82rem;
  font-weight: 600;
  border: 1px solid transparent;
  transition: all 0.15s ease;
  white-space: nowrap;
  flex-shrink: 0;
}

.nav-item:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.nav-item.router-link-active {
  background: #eff6ff;
  color: #2563eb;
  border-color: #bfdbfe;
}

.navbar-user {
  display: flex;
  align-items: center;
  gap: 12px;
}

.qr-quick-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  background: #007aff;
  color: #ffffff;
  border: 1px solid #007aff;
  padding: 0 16px;
  border-radius: 10px !important;
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
  height: 42px;
  line-height: 1;
  transition: all 0.15s ease;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
}

.qr-quick-btn:hover {
  background: #0062cc;
  border-color: #0062cc;
  transform: translateY(-1px);
}

.qr-quick-btn .btn-svg {
  flex-shrink: 0;
  display: block;
}

.account-dropdown-wrapper {
  position: relative;
}

.account-trigger {
  display: flex;
  align-items: center;
  gap: 10px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #0f172a;
  padding: 5px 12px;
  border-radius: 10px !important;
  cursor: pointer;
  height: 42px;
  transition: all 0.15s ease;
}

.account-trigger:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: #007aff;
  color: #ffffff;
  border-radius: 8px !important;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 0.9rem;
  flex-shrink: 0;
  overflow: hidden;
  position: relative;
}

.avatar-img-nav {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
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
  font-size: 0.86rem;
  font-weight: 800;
  color: #0f172a !important;
  margin: 0;
  padding: 0;
  text-align: left;
  white-space: nowrap;
  max-width: 130px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role-sub {
  font-size: 0.72rem;
  font-weight: 600;
  color: #2563eb !important;
  margin: 2px 0 0 0;
  padding: 0;
  text-align: left;
  white-space: nowrap;
  max-width: 130px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: block;
}

.dropdown-arrow {
  font-size: 0.6rem;
  color: #64748b;
  margin-left: 2px;
}

.dropdown-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  width: 230px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px !important;
  box-shadow: 0 12px 32px rgba(15, 23, 42, 0.12);
  padding: 8px;
  z-index: 300;
  animation: dropdownFade 0.15s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes dropdownFade {
  from { opacity: 0; transform: translateY(-6px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

.dropdown-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 10px 10px;
}

.dh-avatar {
  width: 36px;
  height: 36px;
  background: #007aff;
  color: #ffffff;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 800;
  font-size: 1rem;
  flex-shrink: 0;
  overflow: hidden;
}

.dh-info {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 2px;
  overflow: hidden;
}

.dh-title {
  margin: 0;
  font-weight: 800;
  color: #0f172a;
  font-size: 0.88rem;
  white-space: nowrap;
  text-overflow: ellipsis;
  overflow: hidden;
  max-width: 140px;
}

.dh-role-badge {
  font-size: 0.68rem;
  font-weight: 700;
  padding: 2px 8px;
  border-radius: 6px;
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #bfdbfe;
}

.dropdown-divider {
  height: 1px;
  background: #f1f5f9;
  margin: 4px 0 6px;
}

.dropdown-item {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  background: transparent;
  border: none;
  color: #334155;
  font-size: 0.86rem;
  font-weight: 600;
  text-align: left;
  cursor: pointer;
  border-radius: 10px !important;
  transition: all 0.15s ease;
  box-sizing: border-box;
}

.dropdown-item:hover {
  background: #f8fafc;
  color: #0f172a;
}

.dropdown-item.logout-item {
  color: #dc2626;
}

.dropdown-item.logout-item:hover {
  background: #fef2f2;
  color: #b91c1c;
}

.dd-icon {
  flex-shrink: 0;
}

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
  background: #ffffff;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 20px;
  border-right: 1px solid #e2e8f0;
  box-shadow: 10px 0 25px -5px rgba(0, 0, 0, 0.1);
}

.drawer-top-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e2e8f0;
}

.drawer-menu-title {
  font-weight: 800;
  font-size: 1rem;
  color: #0f172a;
  letter-spacing: -0.01em;
}

.close-drawer {
  margin-left: auto;
  background: transparent;
  border: none;
  color: #64748b;
  font-size: 1.2rem;
  cursor: pointer;
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
  color: #334155;
  text-decoration: none;
  padding: 10px 14px;
  min-height: 44px;
  box-sizing: border-box;
  border-radius: 10px !important;
  font-size: 0.9rem;
  font-weight: 600;
  border: 1px solid transparent;
  width: 100%;
  cursor: pointer;
  background: transparent;
  text-align: left;
}

.drawer-item:hover, .drawer-item.router-link-active {
  background: #f1f5f9;
  color: #0f172a;
  border-color: #e2e8f0;
}

.profile-drawer-btn {
  color: #007aff !important;
  font-weight: 700 !important;
}

.profile-drawer-btn:hover {
  background: #eff6ff !important;
}

.drawer-divider {
  height: 1px;
  background: #e2e8f0;
  margin: 8px 0;
}

.drawer-item.logout {
  color: #dc2626;
}

.mobile-profile-avatar-btn {
  display: none;
  background: transparent;
  border: none;
  padding: 0;
  cursor: pointer;
  border-radius: 50%;
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

  .mobile-profile-avatar-btn {
    display: flex !important;
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
  padding: 4px 0 0;
}

.avatar-edit-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 20px;
}

.profile-avatar-large {
  width: 80px;
  height: 80px;
  border-radius: 20px !important;
  background: #007aff;
  color: #ffffff;
  font-size: 2.2rem;
  font-weight: 800;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  overflow: hidden;
  box-shadow: 0 8px 20px rgba(0, 122, 255, 0.25);
  border: 3px solid #ffffff;
}

.avatar-img-large {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.upload-avatar-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  background: #007aff;
  color: #ffffff;
  padding: 8px 16px;
  border-radius: 10px !important;
  font-size: 0.82rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.2);
}

.upload-avatar-btn:hover {
  background: #0062cc;
  transform: translateY(-1px);
}

.profile-read-only-box {
  width: 100%;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 14px !important;
  padding: 16px;
  margin-bottom: 18px;
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 12px;
  box-sizing: border-box;
}

.pro-info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 0.88rem;
}

.pro-info-row span {
  color: #64748b;
  font-weight: 600;
}

.pro-info-row strong {
  color: #0f172a;
  font-weight: 700;
}

.role-badge {
  font-size: 0.75rem;
  font-weight: 700;
  padding: 3px 10px;
  border-radius: 8px !important;
  display: inline-block;
}

.role-badge.admin, .role-badge.hod {
  background: #fef3c7;
  color: #b45309;
  border: 1px solid #fde68a;
}

.role-badge.management {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #bfdbfe;
}

.role-badge.engineer {
  background: #dcfce7;
  color: #15803d;
  border: 1px solid #bbf7d0;
}

.role-badge.external {
  background: #f1f5f9;
  color: #475569;
  border: 1px solid #e2e8f0;
}

.status-active-badge {
  color: #16a34a !important;
}

.admin-only-notice {
  margin: 4px 0 0;
  font-size: 0.78rem;
  color: #64748b;
  border-top: 1px dashed #cbd5e1;
  padding-top: 10px;
  line-height: 1.4;
}

.profile-msg-banner {
  width: 100%;
  padding: 10px 14px;
  border-radius: 10px !important;
  font-size: 0.85rem;
  font-weight: 700;
  margin-bottom: 16px;
  box-sizing: border-box;
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
  gap: 10px;
}

.ios-btn {
  height: 42px;
  padding: 0 20px;
  border-radius: 10px !important;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  border: none;
  transition: all 0.15s ease;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}

.ios-btn-cancel {
  background: #f1f5f9;
  color: #0f172a;
  border: 1px solid #cbd5e1;
}

.ios-btn-cancel:hover {
  background: #e2e8f0;
}

.ios-btn-primary {
  background: #007aff;
  color: #ffffff;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
}

.ios-btn-primary:hover {
  background: #0062cc;
  transform: translateY(-1px);
}

.ios-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none !important;
}

/* ── Notification Bell & Popover (iOS Style) ─────────────── */
.notif-dropdown-wrapper {
  position: relative;
}

.notif-bell-btn {
  width: 42px;
  height: 42px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  color: #0f172a;
  border-radius: 10px !important;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  transition: all 0.15s ease;
}

.notif-bell-btn:hover {
  background: #f1f5f9;
  border-color: #cbd5e1;
}

.notif-badge-count {
  position: absolute;
  top: -4px;
  right: -4px;
  background: #ff3b30;
  color: #ffffff;
  font-size: 0.68rem;
  font-weight: 800;
  height: 18px;
  min-width: 18px;
  padding: 0 4px;
  border-radius: 999px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2px solid #ffffff;
  box-shadow: 0 2px 6px rgba(255, 59, 48, 0.35);
  line-height: 1;
}

.notif-popover-menu {
  position: absolute;
  right: 0;
  top: calc(100% + 8px);
  width: 340px;
  max-width: 90vw;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 16px !important;
  box-shadow: 0 16px 36px rgba(15, 23, 42, 0.14);
  z-index: 400;
  overflow: hidden;
  animation: popoverFade 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes popoverFade {
  from { opacity: 0; transform: translateY(-8px) scale(0.96); }
  to   { opacity: 1; transform: translateY(0) scale(1); }
}

.notif-popover-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  background: #f8fafc;
  border-bottom: 1px solid #f1f5f9;
}

.nph-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.9rem;
  font-weight: 800;
  color: #0f172a;
}

.nph-unread-tag {
  font-size: 0.68rem;
  font-weight: 700;
  background: #eff6ff;
  color: #2563eb;
  padding: 2px 8px;
  border-radius: 999px;
  border: 1px solid #bfdbfe;
}

.nph-read-all-btn {
  background: transparent;
  border: none;
  color: #2563eb;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 6px;
}

.nph-read-all-btn:hover {
  background: #eff6ff;
}

.notif-list-container {
  max-height: 380px;
  overflow-y: auto;
}

.notif-item-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 14px;
  border-bottom: 1px solid #f8fafc;
  cursor: pointer;
  transition: background 0.15s ease;
  position: relative;
}

.notif-item-card:hover {
  background: #f8fafc;
}

.notif-item-card.unread {
  background: #f0f9ff;
}

.notif-type-icon {
  width: 30px;
  height: 30px;
  border-radius: 9px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: #f1f5f9;
  color: #475569;
}

.notif-type-icon.workorder {
  background: #eff6ff;
  color: #2563eb;
}

.notif-type-icon.maintenance {
  background: #f0fdf4;
  color: #16a34a;
}

.notif-type-icon.mutation {
  background: #fff7ed;
  color: #ea580c;
}

.notif-content {
  flex: 1;
  min-width: 0;
}

.notif-title-row {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  gap: 8px;
  margin-bottom: 2px;
}

.notif-item-title {
  font-size: 0.84rem;
  font-weight: 800;
  color: #0f172a;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notif-time-ago {
  font-size: 0.7rem;
  color: #94a3b8;
  font-weight: 600;
  flex-shrink: 0;
}

.notif-item-body {
  margin: 0;
  font-size: 0.78rem;
  color: #475569;
  line-height: 1.35;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.unread-blue-dot {
  width: 7px;
  height: 7px;
  background: #007aff;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 6px;
}

.notif-empty-state {
  padding: 32px 16px;
  text-align: center;
  color: #94a3b8;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.notif-empty-state p {
  margin: 0;
  font-size: 0.85rem;
  font-weight: 600;
}

/* === iOS App Style Mobile Bottom Navigation Dock === */
.mobile-bottom-dock {
  display: none;
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 64px;
  background: rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-top: 1px solid rgba(226, 232, 240, 0.8);
  z-index: 250;
  padding: 0 8px calc(env(safe-area-inset-bottom, 0px) + 2px);
  box-shadow: 0 -4px 24px rgba(15, 23, 42, 0.08);
  justify-content: space-around;
  align-items: center;
}

@media (max-width: 850px) {
  .mobile-bottom-dock {
    display: flex !important;
  }
}

.dock-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 3px;
  color: #64748b;
  text-decoration: none;
  font-size: 0.68rem;
  font-weight: 700;
  flex: 1;
  background: transparent;
  border: none;
  padding: 6px 0;
  transition: all 0.15s ease;
  cursor: pointer;
}

.dock-item.router-link-active {
  color: #007aff;
}

.dock-item:active {
  transform: scale(0.92);
}

.dock-fab-scan {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
  background: transparent;
  border: none;
  color: #007aff;
  font-size: 0.68rem;
  font-weight: 800;
  cursor: pointer;
  margin-top: -18px;
  flex: 1;
}

.fab-icon-wrapper {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #007aff, #0056b3);
  color: #ffffff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 6px 18px rgba(0, 122, 255, 0.4);
  border: 3px solid #ffffff;
  transition: all 0.15s ease;
}

.dock-fab-scan:active .fab-icon-wrapper {
  transform: scale(0.92);
  box-shadow: 0 2px 8px rgba(0, 122, 255, 0.4);
}
</style>
