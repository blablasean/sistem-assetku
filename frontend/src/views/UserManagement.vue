<template>
  <div class="page-container">
    <!-- Non-Admin Access Denied Notice -->
    <div v-if="userRole !== 'admin'" class="access-denied-card">
      <div class="denied-icon">🚫</div>
      <h2>Akses Terbatas (Admin Only)</h2>
      <p>Halaman Manajemen Pengguna ini hanya dapat diakses oleh Administrator Hotel.</p>
      <router-link to="/dashboard" class="action-btn-sharp">Kembali ke Dashboard</router-link>
    </div>

    <!-- Admin User Management Panel -->
    <div v-else class="user-management-panel">
      <!-- Standard Page Header Bar -->
      <div class="page-header">
        <div>
          <p class="eyebrow">Kontrol Pengguna & Peranan</p>
          <h1>👥 Manajemen User & Hak Akses</h1>
          <p class="subtitle">Kelola pendaftaran akun, nama, password, dan peranan (role) staf hotel secara aman.</p>
        </div>

        <button class="primary-btn" @click="openAddModal">
          ➕ Tambah User Baru
        </button>
      </div>

      <!-- Toast Notification -->
      <div v-if="showToast" class="toast-banner" :class="toastType">
        {{ toastMsg }}
      </div>

      <!-- Metric Cards -->
      <div class="summary-cards-grid">
        <div class="sum-card">
          <div class="sc-icon">👥</div>
          <div>
            <p class="sc-label">Total Pengguna</p>
            <h3 class="sc-val">{{ userList.length }}</h3>
          </div>
        </div>

        <div class="sum-card">
          <div class="sc-icon">👑</div>
          <div>
            <p class="sc-label">Admin & HOD</p>
            <h3 class="sc-val">{{ adminCount }}</h3>
          </div>
        </div>

        <div class="sum-card">
          <div class="sc-icon">🛠️</div>
          <div>
            <p class="sc-label">Staff & Engineer</p>
            <h3 class="sc-val">{{ staffCount }}</h3>
          </div>
        </div>
      </div>

      <!-- Search & Filter Bar -->
      <div class="filter-card">
        <div class="search-input-box">
          <span class="s-icon">🔍</span>
          <input
            v-model="searchQuery"
            placeholder="Cari nama, username, atau role pengguna..."
            class="sharp-input"
          />
        </div>

        <select v-model="roleFilter" class="sharp-select">
          <option value="">Semua Role / Jabatan</option>
          <option value="admin">Administrator</option>
          <option value="hod">HOD Engineer</option>
          <option value="management">Supervisor Engineer</option>
          <option value="engineer">Staff Engineer</option>
          <option value="dept_akunting">Departement Akunting</option>
          <option value="dept_spa">Departement Spa</option>
          <option value="dept_sales">Department Sales</option>
          <option value="dept_hr">Department Human Resources</option>
          <option value="dept_fb_kitchen">Department Food Beverage Kitchen</option>
          <option value="dept_fb_service">Department Food Beverage Service</option>
          <option value="dept_housekeeping">Department House Keeping</option>
          <option value="dept_frontoffice">Department Front Office</option>
        </select>
      </div>

      <!-- Users Data Table -->
      <div class="table-container-sharp">
        <table class="data-table-sharp">
          <thead>
            <tr>
              <th>ID</th>
              <th>Nama Lengkap</th>
              <th>Username</th>
              <th>Role / Jabatan</th>
              <th style="text-align: center;">Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="isLoading">
              <td colspan="5" class="empty-cell">Memuat data pengguna...</td>
            </tr>
            <tr v-else-if="filteredUsers.length === 0">
              <td colspan="5" class="empty-cell">Tidak ada pengguna yang sesuai pencarian.</td>
            </tr>
            <tr v-else v-for="user in filteredUsers" :key="user.id">
              <td><strong>#{{ user.id }}</strong></td>
              <td>
                <div class="user-cell">
                  <div class="user-avatar-sm">
                    <img v-if="user.avatar" :src="user.avatar" class="u-avatar-img" />
                    <span v-else>{{ getInitial(user.name) }}</span>
                  </div>
                  <span class="u-name">{{ user.name }}</span>
                </div>
              </td>
              <td><code>{{ user.username }}</code></td>
              <td>
                <span class="role-badge-sharp" :class="user.role">
                  {{ getRoleLabel(user.role) }}
                </span>
              </td>
              <td>
                <div class="action-flex">
                  <button class="tbl-btn edit-btn" @click="openEditModal(user)">
                    ✏️ Edit
                  </button>
                  <button class="tbl-btn delete-btn" @click="deleteUser(user)">
                    🗑️ Hapus
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Form User (Tambah / Edit) -->
    <ModalDialog
      :show="showUserModal"
      :title="isEditMode ? '✏️ Edit Data Pengguna' : '➕ Registrasi Pengguna Baru'"
      maxWidth="500px"
      @close="showUserModal = false"
    >
      <form @submit.prevent="saveUser" class="form-sharp">
        <label>
          Nama Lengkap Staff:
          <input
            v-model="formUser.name"
            required
            placeholder="Contoh: Budi Santoso, S.T."
            class="sharp-input"
          />
        </label>

        <label>
          Username Login:
          <input
            v-model="formUser.username"
            required
            placeholder="Contoh: budi_eng"
            class="sharp-input"
          />
        </label>

        <label>
          Password {{ isEditMode ? '(Kosongkan jika tidak diubah)' : '' }}:
          <input
            v-model="formUser.password"
            :required="!isEditMode"
            type="password"
            placeholder="Password rahasia..."
            class="sharp-input"
          />
        </label>

        <label>
          Role / Hak Akses Sistem:
          <select v-model="formUser.role" required class="sharp-select">
            <option value="admin">👑 Administrator (Akses Penuh)</option>
            <option value="hod">⭐ Head of Department (HOD Engineer)</option>
            <option value="management">👔 Supervisor Engineer</option>
            <option value="engineer">🛠️ Staff Engineer (Teknisi Lapangan)</option>
            <option value="dept_akunting">📊 Departement Akunting</option>
            <option value="dept_spa">💆 Departement Spa</option>
            <option value="dept_sales">📈 Department Sales</option>
            <option value="dept_hr">👥 Department Human Resources</option>
            <option value="dept_fb_kitchen">🍳 Department Food Beverage Kitchen</option>
            <option value="dept_fb_service">🍽️ Department Food Beverage Service</option>
            <option value="dept_housekeeping">🧹 Department House Keeping</option>
            <option value="dept_frontoffice">🛎️ Department Front Office</option>
          </select>
        </label>

        <div class="modal-form-actions">
          <button type="button" class="action-btn-sharp cancel-btn" @click="showUserModal = false">
            Batal
          </button>
          <button type="submit" class="action-btn-sharp primary-btn">
            {{ isEditMode ? 'Simpan Perubahan' : 'Daftarkan User' }}
          </button>
        </div>
      </form>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')
const userList = ref([])
const isLoading = ref(false)

const searchQuery = ref('')
const roleFilter = ref('')

const showUserModal = ref(false)
const isEditMode = ref(false)
const formUser = ref({ user_id: 0, username: '', password: '', name: '', role: 'engineer' })

const showToast = ref(false)
const toastMsg = ref('')
const toastType = ref('success')

function notify(msg, type = 'success') {
  toastMsg.value = msg
  toastType.value = type
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

const adminCount = computed(() => {
  return userList.value.filter(u => u.role === 'admin' || u.role === 'hod').length
})

const staffCount = computed(() => {
  return userList.value.filter(u => u.role === 'engineer' || u.role === 'external' || u.role === 'management').length
})

const filteredUsers = computed(() => {
  let list = [...userList.value]

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(u =>
      (u.name && u.name.toLowerCase().includes(q)) ||
      (u.username && u.username.toLowerCase().includes(q)) ||
      (u.role && u.role.toLowerCase().includes(q))
    )
  }

  if (roleFilter.value) {
    list = list.filter(u => u.role === roleFilter.value)
  }

  return list
})

async function fetchUsers() {
  userRole.value = sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external'
  if (userRole.value !== 'admin') return

  isLoading.value = true
  try {
    const res = await api.get('/users')
    const rawData = res.data?.data || res.data || []
    if (Array.isArray(rawData)) {
      userList.value = rawData.map(u => ({
        id: u.id || u.ID,
        name: u.name || u.Name || 'User',
        username: u.username || u.Username || '',
        role: u.role || u.Role || 'external',
        avatar: u.avatar || u.Avatar || ''
      }))
    } else {
      userList.value = []
    }
  } catch (e) {
    console.error('Fetch users error:', e)
    notify('Gagal memuat data pengguna: ' + (e.response?.data?.message || e.message), 'error')
  } finally {
    isLoading.value = false
  }
}

function openAddModal() {
  isEditMode.value = false
  formUser.value = { user_id: 0, username: '', password: '', name: '', role: 'engineer' }
  showUserModal.value = true
}

function openEditModal(user) {
  isEditMode.value = true
  formUser.value = {
    user_id: user.id,
    username: user.username,
    password: '',
    name: user.name,
    role: user.role
  }
  showUserModal.value = true
}

async function saveUser() {
  try {
    if (isEditMode.value) {
      await api.post('/users/edit', formUser.value)
      notify('Data pengguna berhasil diperbarui!', 'success')
    } else {
      await api.post('/users/create', formUser.value)
      notify('Pengguna baru berhasil didaftarkan!', 'success')
    }
    showUserModal.value = false
    await fetchUsers()
  } catch (e) {
    notify('Gagal menyimpan pengguna: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteUser(user) {
  if (!confirm(`Apakah Anda yakin ingin menghapus akun "${user.name}" (${user.username})?`)) return

  try {
    await api.post('/users/delete', { user_id: user.id })
    notify(`Pengguna "${user.name}" berhasil dihapus!`, 'success')
    await fetchUsers()
  } catch (e) {
    notify('Gagal menghapus pengguna: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function getInitial(name) {
  return name ? name.charAt(0).toUpperCase() : 'U'
}

function getRoleLabel(role) {
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
  return map[role] || role
}

onMounted(() => {
  fetchUsers()
})
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 16px;
}

.access-denied-card {
  background: #ffffff;
  border: 1px solid #fca5a5;
  border-radius: 2px !important;
  padding: 40px;
  text-align: center;
  margin-top: 40px;
}

.denied-icon {
  font-size: 3rem;
  margin-bottom: 12px;
}

.access-denied-card h2 {
  color: #991b1b;
  margin: 0 0 8px;
}

.access-denied-card p {
  color: #475569;
  margin-bottom: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  gap: 16px;
}

.eyebrow {
  text-transform: uppercase;
  font-size: 0.75rem;
  letter-spacing: 0.1em;
  color: #d97706;
  font-weight: 800;
  margin-bottom: 4px;
}

.page-header h1 {
  margin: 0;
  font-size: 1.75rem;
  color: #0f172a;
  font-weight: 800;
}

.subtitle {
  margin-top: 6px;
  color: #64748b;
  font-size: 0.95rem;
}

.primary-btn {
  background: #0f172a;
  color: white;
  border: 1px solid #1e293b;
  padding: 10px 18px;
  border-radius: 2px !important;
  font-weight: 700;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.primary-btn:hover {
  background: #1e293b;
  border-color: #d97706;
}

.action-btn-sharp {
  background: #0f172a;
  color: white;
  border: 1px solid #1e293b;
  padding: 10px 16px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.88rem;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.primary-btn {
  background: #2563eb;
  border-color: #1d4ed8;
}

.cancel-btn {
  background: #64748b;
  border-color: #475569;
}

.toast-banner {
  padding: 12px 16px;
  border-radius: 2px !important;
  margin-bottom: 16px;
  font-weight: 700;
  font-size: 0.88rem;
}

.toast-banner.success {
  background: #dcfce7;
  color: #166534;
  border: 1px solid #bbf7d0;
}

.toast-banner.error {
  background: #fee2e2;
  color: #991b1b;
  border: 1px solid #fca5a5;
}

.summary-cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 20px;
}

.sum-card {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.sc-icon {
  font-size: 2rem;
  background: #f1f5f9;
  width: 50px;
  height: 50px;
  border-radius: 2px !important;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #e2e8f0;
}

.sc-label {
  margin: 0;
  font-size: 0.78rem;
  color: #64748b;
  text-transform: uppercase;
  font-weight: 700;
}

.sc-val {
  margin: 2px 0 0;
  font-size: 1.4rem;
  color: #0f172a;
  font-weight: 800;
}

.filter-card {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  padding: 14px 16px;
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.search-input-box {
  flex: 1;
  min-width: 250px;
  display: flex;
  align-items: center;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  padding: 0 12px;
}

.s-icon {
  color: #64748b;
  margin-right: 8px;
}

.sharp-input, .sharp-select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  font-size: 0.88rem;
  background: #ffffff;
}

.search-input-box .sharp-input {
  border: none;
  background: transparent;
  padding: 10px 0;
}

.table-container-sharp {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  overflow-x: auto;
}

.data-table-sharp {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table-sharp th {
  background: #0f172a;
  color: #ffffff;
  padding: 12px 16px;
  font-size: 0.8rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table-sharp td {
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
  font-size: 0.88rem;
  color: #334155;
}

.empty-cell {
  text-align: center;
  padding: 30px !important;
  color: #64748b;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 10px;
}

.user-avatar-sm {
  width: 32px;
  height: 32px;
  background: #d97706;
  color: white;
  font-weight: 800;
  border-radius: 2px !important;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.85rem;
}

.u-name {
  font-weight: 700;
  color: #0f172a;
}

.role-badge-sharp {
  display: inline-block;
  padding: 3px 8px;
  border-radius: 2px !important;
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  border: 1px solid transparent;
}

.role-badge-sharp.admin { background: #fef3c7; color: #b45309; border-color: #fde68a; }
.role-badge-sharp.hod { background: #dbeafe; color: #1e40af; border-color: #bfdbfe; }
.role-badge-sharp.management { background: #ffedd5; color: #c2410c; border-color: #fed7aa; }
.role-badge-sharp.engineer { background: #dcfce7; color: #15803d; border-color: #bbf7d0; }
.role-badge-sharp.external { background: #f1f5f9; color: #475569; border-color: #cbd5e1; }

.action-flex {
  display: flex;
  gap: 6px;
  justify-content: center;
}

.tbl-btn {
  border: 1px solid #cbd5e1;
  padding: 4px 10px;
  border-radius: 2px !important;
  font-size: 0.75rem;
  font-weight: 700;
  cursor: pointer;
  background: #f8fafc;
}

.edit-btn:hover { background: #e0f2fe; color: #0369a1; border-color: #bae6fd; }
.delete-btn:hover { background: #fee2e2; color: #991b1b; border-color: #fca5a5; }

.form-sharp {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.form-sharp label {
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
}

.modal-form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 10px;
}

.u-avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 2px !important;
}
</style>
