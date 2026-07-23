<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Aset Hotel</p>
        <h1>📦 Manajemen Aset</h1>
        <p class="subtitle">Daftar inventaris aset, lokasi, PIC, dan riwayat mutasi.</p>
      </div>

      <div class="header-actions" v-if="canCreateAsset">
        <button class="primary-btn" @click="openAddModal">
          ➕ Tambah Aset
        </button>
      </div>
    </div>

    <div class="card-panel">
      <!-- Toolbar: Search, Filter, Sort -->
      <div class="toolbar-grid">
        <input v-model="searchQuery" placeholder="🔍 Cari nama, kode, lokasi, atau PIC aset..." class="search-input" @input="filterAssets" />
        
        <select v-model="filterStatus" @change="filterAssets" class="filter-select">
          <option value="">Semua Status</option>
          <option value="Active">Active</option>
          <option value="Maintenance">Maintenance</option>
          <option value="Damaged">Damaged</option>
          <option value="Reserved">Reserved</option>
          <option value="Retired">Retired</option>
        </select>

        <select v-model="sortBy" @change="filterAssets" class="sort-select">
          <option value="id-desc">Terbaru</option>
          <option value="name-asc">Nama Aset (A-Z)</option>
          <option value="location-asc">Lokasi / Kamar</option>
        </select>
      </div>

      <!-- Assets Table -->
      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>Kode Aset</th>
              <th>Nama Aset</th>
              <th>Kategori</th>
              <th>Lokasi / Kamar</th>
              <th>PIC</th>
              <th>Status</th>
              <th>Dokumen</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="asset in displayedAssets" :key="asset.id">
              <td><span class="code-badge">{{ asset.asset_code }}</span></td>
              <td>
                <span class="asset-name" @click="viewDetail(asset)">{{ asset.asset_name }}</span>
                <span v-if="asset.is_reserved" class="reserved-tag">Reserved</span>
              </td>
              <td>{{ asset.category || 'General' }}</td>
              <td>📍 {{ asset.location || 'Main Store' }}</td>
              <td>👤 {{ asset.pic || 'Engineering' }}</td>
              <td><StatusBadge :status="asset.status" /></td>
              <td>
                <a v-if="asset.document_url" :href="asset.document_url" target="_blank" class="doc-link">📄 Manual</a>
                <span v-else class="no-doc">-</span>
              </td>
              <td class="actions-cell">
                <button class="icon-btn qr-btn" @click="openQrPrint(asset)" title="Generate & Cetak QR Code">
                  🖨️ QR
                </button>
                <button class="icon-btn mut-btn" @click="openMutationModal(asset)" title="Mutasi Lokasi Barang">
                  🔄 Mutasi
                </button>
                <button class="icon-btn edit-btn" v-if="canCreateAsset" @click="openEditModal(asset)" title="Edit Aset">
                  ✏️
                </button>
                <button class="icon-btn delete-btn" v-if="canDeleteAsset" @click="deleteAsset(asset)" title="Hapus Aset Permanen">
                  🗑️
                </button>
              </td>
            </tr>
            <tr v-if="displayedAssets.length === 0">
              <td colspan="8" class="empty-state">Tidak ada data aset yang ditemukan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Custom UI Toast Notification -->
    <transition name="fade">
      <div v-if="showToast" :class="['custom-ui-toast', toastType]">
        <span class="toast-icon">{{ toastType === 'success' ? '✅' : '⚠️' }}</span>
        <span class="toast-text">{{ toastMsg }}</span>
        <button class="toast-close" @click="showToast = false">✕</button>
      </div>
    </transition>

    <!-- Modal Registrasi / Edit Aset -->
    <ModalDialog :show="showAssetModal" :title="isEditMode ? '✏️ Edit Data Aset' : '➕ Registrasi Aset Baru'" @close="showAssetModal = false">
      <form @submit.prevent="saveAsset" class="modal-form">
        <label>
          <span>Kode Aset Unik</span>
          <input v-model="formAsset.asset_code" placeholder="Contoh: AST-RM301-AC" :disabled="isEditMode" required />
        </label>
        <label>
          <span>Nama Aset</span>
          <input v-model="formAsset.asset_name" placeholder="Contoh: AC Split Daikin 1.5 PK" required />
        </label>
        <label>
          <span>Kategori</span>
          <select v-model="formAsset.category">
            <option>HVAC / AC</option>
            <option>Elektronik & TV</option>
            <option>Mebel & Furniture</option>
            <option>Plumbing & Sanitasi</option>
            <option>Mesin & Generator</option>
            <option>Kitchen Equipment</option>
          </select>
        </label>
        <label>
          <span>Lokasi Penempatan / Kamar</span>
          <input v-model="formAsset.location" placeholder="Contoh: Kamar 301, Kitchen Dapur Utama, Chiller Room" required />
        </label>
        <label>
          <span>PIC Penanggung Jawab</span>
          <input v-model="formAsset.pic" placeholder="Contoh: Supervisor Engineering" />
        </label>
        <label>
          <span>Status Kondisi Aset</span>
          <select v-model="formAsset.status">
            <option value="Active">Active (Berfungsi Normal)</option>
            <option value="Maintenance">Maintenance (Dalam Perawatan)</option>
            <option value="Damaged">Damaged (Rusak)</option>
            <option value="Reserved">Reserved (Cadangan)</option>
            <option value="Retired">Retired (Dihentikan)</option>
          </select>
        </label>
        <label>
          <span>URL Dokumen / Manual Book</span>
          <input v-model="formAsset.document_url" placeholder="https://..." />
        </label>

        <button type="submit" class="submit-modal-btn">
          {{ isEditMode ? 'Simpan Perubahan' : 'Daftarkan Aset' }}
        </button>
      </form>
    </ModalDialog>

    <!-- Modal Mutasi Lokasi Barang -->
    <ModalDialog :show="showMutModal" title="🔄 Mutasi Lokasi Aset" @close="showMutModal = false">
      <div v-if="selectedAssetForMut" class="mut-modal-body">
        <div class="current-info">
          <p><strong>Aset:</strong> {{ selectedAssetForMut.asset_name }} ({{ selectedAssetForMut.asset_code }})</p>
          <p><strong>Lokasi Sekarang:</strong> 📍 {{ selectedAssetForMut.location }}</p>
        </div>

        <form @submit.prevent="submitMutation" class="modal-form">
          <label>
            <span>Lokasi Baru (Kamar / Area Hotel)</span>
            <input v-model="mutNewLocation" placeholder="Contoh: Kamar 205, Lobby Lounge" required />
          </label>
          <label>
            <span>PIC Penanggung Jawab Baru</span>
            <input v-model="mutPIC" placeholder="Nama PIC baru" required />
          </label>
          <label>
            <span>Alasan Perpindahan / Mutasi</span>
            <input v-model="mutReason" placeholder="Contoh: Penyesuaian kebutuhan event / pemindahan room hotel" />
          </label>

          <button type="submit" class="submit-modal-btn warning-btn">Proses Mutasi Lokasi</button>
        </form>
      </div>
    </ModalDialog>

    <!-- Modal Cetak QR Code -->
    <ModalDialog :show="showQrModal" title="🖨️ Cetak & Preview QR Code Aset" @close="showQrModal = false">
      <div v-if="selectedAssetForQr" class="qr-print-card">
        <div class="printable-badge" id="qrPrintArea">
          <div class="qr-code-box">
            <span class="qr-placeholder">🔳</span>
          </div>
          <div class="qr-meta">
            <h3>AsetKu Hotel</h3>
            <p class="q-code">{{ selectedAssetForQr.asset_code }}</p>
            <p class="q-name">{{ selectedAssetForQr.asset_name }}</p>
            <p class="q-loc">📍 {{ selectedAssetForQr.location }}</p>
          </div>
        </div>

        <p class="qr-instruction">Tempelkan stiker QR Code ini pada unit fisik aset untuk akses cepat scan.</p>
        <button class="submit-modal-btn" @click="printQrCard">🖨️ Cetak Stiker QR</button>
      </div>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')
const canCreateAsset = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canDeleteAsset = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')

const showToast = ref(false)
const toastMsg = ref('')
const toastType = ref('success')

function notify(msg, type = 'success') {
  toastMsg.value = msg
  toastType.value = type
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 4000)
}

const assets = ref([])
const isLoading = ref(false)

const searchQuery = ref('')
const filterStatus = ref('')
const sortBy = ref('id-desc')

const showAssetModal = ref(false)
const isEditMode = ref(false)
const formAsset = ref({ id: 0, asset_code: '', asset_name: '', category: 'HVAC / AC', location: '', pic: '', status: 'Active', document_url: '' })

const showMutModal = ref(false)
const selectedAssetForMut = ref(null)
const mutNewLocation = ref('')
const mutPIC = ref('')
const mutReason = ref('')

const showQrModal = ref(false)
const selectedAssetForQr = ref(null)

const displayedAssets = computed(() => {
  let list = [...assets.value]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(a => 
      (a.asset_name && a.asset_name.toLowerCase().includes(q)) || 
      (a.asset_code && a.asset_code.toLowerCase().includes(q)) || 
      (a.location && a.location.toLowerCase().includes(q)) ||
      (a.pic && a.pic.toLowerCase().includes(q))
    )
  }
  if (filterStatus.value) {
    list = list.filter(a => a.status === filterStatus.value)
  }
  if (sortBy.value === 'name-asc') {
    list.sort((a, b) => (a.asset_name || '').localeCompare(b.asset_name || ''))
  } else if (sortBy.value === 'location-asc') {
    list.sort((a, b) => (a.location || '').localeCompare(b.location || ''))
  } else {
    list.sort((a, b) => b.id - a.id)
  }
  return list
})

async function fetchAssets() {
  isLoading.value = true
  try {
    const res = await api.get('/assets')
    if (res.data?.data && Array.isArray(res.data.data)) {
      assets.value = res.data.data
    } else {
      assets.value = []
    }
  } catch (e) {
    console.error('Failed to fetch assets from backend:', e)
    assets.value = []
  } finally {
    isLoading.value = false
  }
}

function openAddModal() {
  isEditMode.value = false
  formAsset.value = { id: 0, asset_code: 'AST-RM' + Math.floor(100 + Math.random() * 900) + '-UNIT', asset_name: '', category: 'HVAC / AC', location: '', pic: '', status: 'Active', document_url: '' }
  showAssetModal.value = true
}

function openEditModal(asset) {
  isEditMode.value = true
  formAsset.value = { ...asset }
  showAssetModal.value = true
}

async function saveAsset() {
  try {
    if (isEditMode.value) {
      await api.put(`/assets?id=${formAsset.value.id}`, formAsset.value)
      notify('Perubahan data aset berhasil disimpan!', 'success')
    } else {
      await api.post('/assets', formAsset.value)
      notify('Aset baru berhasil didaftarkan!', 'success')
    }
    showAssetModal.value = false
    await fetchAssets()
  } catch (e) {
    notify('Gagal menyimpan aset: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteAsset(asset) {
  try {
    await api.post('/assets/delete', { asset_id: asset.id })
    assets.value = assets.value.filter(a => a.id !== asset.id)
    notify(`Aset "${asset.asset_name}" berhasil dihapus permanen!`, 'success')
    await fetchAssets()
  } catch (e) {
    notify('Gagal menghapus aset: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openMutationModal(asset) {
  selectedAssetForMut.value = asset
  mutNewLocation.value = ''
  mutPIC.value = asset.pic || ''
  mutReason.value = ''
  showMutModal.value = true
}

async function submitMutation() {
  if (!selectedAssetForMut.value) return
  try {
    await api.post('/mutations', {
      asset_id: selectedAssetForMut.value.id,
      new_location: mutNewLocation.value,
      new_pic: mutPIC.value,
      reason: mutReason.value
    })
    showMutModal.value = false
    notify(`Mutasi aset "${selectedAssetForMut.value.asset_name}" ke ${mutNewLocation.value} berhasil dicatat!`, 'success')
    await fetchAssets()
  } catch (e) {
    notify('Gagal mencatat mutasi: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openQrPrint(asset) {
  selectedAssetForQr.value = asset
  showQrModal.value = true
}

function printQrCard() {
  window.print()
}

function viewDetail(asset) {
  alert(`Detail Aset Hotel:\n\nKode: ${asset.asset_code}\nNama: ${asset.asset_name}\nLokasi: ${asset.location}\nPIC: ${asset.pic}\nStatus: ${asset.status}`)
}

onMounted(() => {
  fetchAssets()
})
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.eyebrow {
  margin: 0 0 4px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.8rem;
  color: #2563eb;
  font-weight: 700;
}

h1 {
  margin: 0 0 8px;
  font-size: 1.8rem;
  color: #0f172a;
}

.subtitle {
  margin: 0;
  color: #64748b;
  font-size: 0.95rem;
}

.primary-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.toolbar-grid {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 260px;
  padding: 10px 16px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
}

.filter-select, .sort-select {
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  background: white;
}

.table-responsive {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  text-align: left;
  padding: 12px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.85rem;
  border-bottom: 1px solid #e2e8f0;
}

td {
  padding: 14px 12px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.9rem;
  color: #334155;
}

.code-badge {
  background: #e2e8f0;
  color: #0f172a;
  padding: 3px 8px;
  border-radius: 6px;
  font-family: monospace;
  font-weight: 700;
}

.asset-name {
  font-weight: 600;
  color: #2563eb;
  cursor: pointer;
}

.reserved-tag {
  margin-left: 6px;
  background: #ffedd5;
  color: #c2410c;
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 700;
}

.doc-link {
  color: #0284c7;
  text-decoration: none;
  font-weight: 600;
}

.no-doc { color: #94a3b8; }

.actions-cell {
  display: flex;
  gap: 6px;
}

.icon-btn {
  border: 1px solid #cbd5e1;
  background: white;
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background 0.15s;
}

.icon-btn:hover { background: #f1f5f9; }

.modal-form {
  display: grid;
  gap: 16px;
}

.modal-form label {
  display: grid;
  gap: 6px;
  font-weight: 600;
  color: #1e293b;
}

.modal-form input, .modal-form select {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
}

.submit-modal-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  margin-top: 8px;
}

.warning-btn { background: #d97706; }

.current-info {
  background: #f8fafc;
  padding: 12px;
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 0.9rem;
}

.qr-print-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.printable-badge {
  border: 2px solid #0f172a;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  background: #ffffff;
  width: 100%;
  max-width: 380px;
}

.qr-code-box {
  width: 90px;
  height: 90px;
  background: #0f172a;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 3rem;
}

.qr-meta h3 { margin: 0; font-size: 1.1rem; color: #0f172a; }
.q-code { margin: 4px 0; font-family: monospace; font-weight: 800; color: #2563eb; }
.q-name { margin: 2px 0; font-weight: 600; font-size: 0.9rem; }
.q-loc { margin: 4px 0 0; color: #64748b; font-size: 0.85rem; }

.qr-instruction {
  font-size: 0.85rem;
  color: #64748b;
  text-align: center;
}
</style>
