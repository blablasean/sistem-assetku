<template>
  <div class="page-container">
    <!-- Non-Admin Access Denied Notice -->
    <div v-if="userRole !== 'admin'" class="access-denied-card">
      <div class="denied-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="#dc2626" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
      </div>
      <h2>Akses Terbatas (Admin Only)</h2>
      <p>Halaman Manajemen Pengguna ini hanya dapat diakses oleh Administrator Sistem.</p>
      <router-link to="/dashboard" class="action-btn-sharp">Kembali ke Dashboard</router-link>
    </div>

    <!-- Admin User Management Panel -->
    <div v-else class="user-management-panel">
      <!-- Standard Page Header Bar -->
      <div class="page-header">
        <div>
          <p class="eyebrow">Kontrol Pengguna & Peranan</p>
          <h1>Manajemen User & Hak Akses</h1>
          <p class="subtitle">Kelola pendaftaran akun, nama, password, dan peranan (role) pengguna secara aman.</p>
        </div>

        <button class="primary-btn" @click="openAddModal">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 6px;"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          <span>Tambah User Baru</span>
        </button>
      </div>

      <!-- Toast Notification -->
      <div v-if="showToast" class="toast-banner" :class="toastType">
        {{ toastMsg }}
      </div>

      <!-- Metric Cards -->
      <div class="summary-cards-grid">
        <div class="sum-card">
          <div class="sc-icon blue">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M22 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>
          </div>
          <div>
            <p class="sc-label">Total Pengguna</p>
            <h3 class="sc-val">{{ userList.length }}</h3>
          </div>
        </div>

        <div class="sum-card">
          <div class="sc-icon amber">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#d97706" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M11.562 3.266a.5.5 0 0 1 .876 0L15.3 8.87l6.2.902a.5.5 0 0 1 .277.853l-4.486 4.372 1.059 6.175a.5.5 0 0 1-.725.527L12 18.78l-5.625 2.957a.5.5 0 0 1-.725-.527l1.059-6.175L2.223 10.625a.5.5 0 0 1 .277-.853l6.2-.902z"/></svg>
          </div>
          <div>
            <p class="sc-label">Admin & HOD</p>
            <h3 class="sc-val">{{ adminCount }}</h3>
          </div>
        </div>

        <div class="sum-card">
          <div class="sc-icon emerald">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#059669" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
          </div>
          <div>
            <p class="sc-label">Staff & Engineer</p>
            <h3 class="sc-val">{{ staffCount }}</h3>
          </div>
        </div>
      </div>

      <!-- Search & Filter Bar -->
      <div class="filter-card">
        <div class="search-input-box">
          <span class="s-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="11" r="8"/><path d="m21 21-4.3-4.3"/></svg>
          </span>
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
                  <button class="tbl-btn edit-btn" @click="openEditModal(user)" title="Edit User">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                    <span>Edit</span>
                  </button>
                  <button class="tbl-btn delete-btn" @click="deleteUser(user)" title="Hapus User">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                    <span>Hapus</span>
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
      :title="isEditMode ? 'Edit Data Pengguna' : 'Registrasi Pengguna Baru'"
      maxWidth="500px"
      @close="showUserModal = false"
    >
      <form @submit.prevent="saveUser" class="form-sharp">
        <label>
          <span>Nama Lengkap Staff:</span>
          <input
            v-model="formUser.name"
            required
            placeholder="Contoh: Budi Santoso, S.T."
            class="sharp-input"
          />
        </label>

        <label>
          <span>Username Login:</span>
          <input
            v-model="formUser.username"
            required
            placeholder="Contoh: budi_eng"
            class="sharp-input"
          />
        </label>

        <label>
          <span>Password {{ isEditMode ? '(Kosongkan jika tidak diubah)' : '' }}:</span>
          <input
            v-model="formUser.password"
            :required="!isEditMode"
            type="password"
            placeholder="Password rahasia..."
            class="sharp-input"
          />
        </label>

        <label>
          <span>Role / Hak Akses Sistem:</span>
          <select v-model="formUser.role" required class="sharp-select">
            <option value="admin">Administrator (Akses Penuh)</option>
            <option value="hod">Head of Department (HOD Engineer)</option>
            <option value="management">Supervisor Engineer</option>
            <option value="engineer">Staff Engineer (Teknisi Lapangan)</option>
            <option value="dept_akunting">Departement Akunting</option>
            <option value="dept_spa">Departement Spa</option>
            <option value="dept_sales">Department Sales</option>
            <option value="dept_hr">Department Human Resources</option>
            <option value="dept_fb_kitchen">Department Food Beverage Kitchen</option>
            <option value="dept_fb_service">Department Food Beverage Service</option>
            <option value="dept_housekeeping">Department House Keeping</option>
            <option value="dept_frontoffice">Department Front Office</option>
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

const userRole = ref(sessionStorage.getItem('user_role') || 'external')
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
    external: 'Staff Operasional'
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
  padding: 24px 24px;
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
  flex-wrap: wrap;
}

.eyebrow {
  margin: 0 0 4px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.8rem;
  color: #2563eb;
  font-weight: 700;
}

h1, .page-header h1 {
  margin: 0 0 8px;
  font-size: 1.8rem;
  color: #0f172a;
  font-weight: 800;
}

.subtitle {
  margin: 0;
  color: #64748b;
  font-size: 0.95rem;
}

.primary-btn {
  background: #007aff !important;
  color: #ffffff !important;
  border: 1px solid #007aff !important;
  padding: 10px 18px !important;
  border-radius: 10px !important;
  font-size: 0.88rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
  transition: all 0.15s ease;
  line-height: 1;
  white-space: nowrap;
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.action-btn-sharp {
  background: #0f172a;
  color: white;
  border: 1px solid #1e293b;
  padding: 10px 16px;
  border-radius: 6px !important;
  font-weight: 700;
  font-size: 0.88rem;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  text-decoration: none;
}

.cancel-btn {
  background: #64748b;
  border-color: #475569;
}

.toast-banner {
  padding: 12px 16px;
  border-radius: 8px !important;
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
  border: 1px solid #e2e8f0;
  border-radius: 12px !important;
  padding: 16px;
  display: flex;
  align-items: center;
  gap: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.sc-icon {
  width: 44px;
  height: 44px;
  border-radius: 10px !important;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  flex-shrink: 0;
}

.sc-icon.blue { background: #eff6ff; border-color: #bfdbfe; }
.sc-icon.amber { background: #fff7ed; border-color: #ffedd5; }
.sc-icon.emerald { background: #ecfdf5; border-color: #a7f3d0; }

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
  border: 1px solid #e2e8f0;
  border-radius: 12px !important;
  padding: 14px 16px;
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.search-input-box {
  flex: 1;
  min-width: 250px;
  display: flex;
  align-items: center;
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 8px !important;
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
  border-radius: 8px !important;
  font-size: 0.9rem;
  background: #ffffff;
  color: #0f172a;
  outline: none;
  transition: all 0.15s ease;
}

.sharp-input:focus, .sharp-select:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
}

.search-input-box .sharp-input {
  border: none;
  background: transparent;
  padding: 10px 0;
}

.table-container-sharp {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 12px !important;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
}

.data-table-sharp {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.data-table-sharp th {
  background: #f8fafc;
  color: #0f172a;
  padding: 12px 16px;
  font-size: 0.8rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  border-bottom: 1px solid #e2e8f0;
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
  background: #0f172a;
  color: white;
  font-weight: 800;
  border-radius: 8px !important;
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
  padding: 4px 10px;
  border-radius: 6px !important;
  font-size: 0.75rem;
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
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  height: 32px;
  padding: 0 10px;
  border-radius: 6px !important;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  box-sizing: border-box;
  line-height: 1;
  text-decoration: none;
  transition: all 0.15s ease;
  border: 1px solid transparent;
  user-select: none;
}

.tbl-btn svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  display: block;
}

.tbl-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.edit-btn {
  background: #f8fafc;
  color: #475569;
  border-color: #cbd5e1;
}

.edit-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

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
  border-radius: 8px !important;
}

/* === Mobile Responsive CSS (Android & iOS) === */
@media (max-width: 640px) {
  .page-container { padding: 16px 14px !important; }
  .page-header { flex-direction: column; align-items: stretch; gap: 12px; }
  .page-header .primary-btn { width: 100%; justify-content: center; height: 40px !important; font-size: 0.82rem !important; }
  .card-panel { padding: 16px !important; border-radius: 14px !important; }
}
</style>
