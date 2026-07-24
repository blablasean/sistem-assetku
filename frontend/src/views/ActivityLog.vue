<template>
  <div class="view-wrapper">
    <div class="page-container">
      <div class="page-header">
        <div>
          <p class="eyebrow">Audit Trail</p>
          <h1>📋 Activity Log</h1>
          <p class="subtitle">Riwayat perbaikan & maintenance.</p>
        </div>

        <div class="header-action-group">
          <button class="print-report-btn" @click="showReportModal = true" title="Prinjau & Export Laporan Audit Trail">
            📄 Laporan & Export
          </button>
        </div>
      </div>

      <!-- Summary boxes -->
      <div class="summary-row">
        <div class="sbox green">
          <span class="sbox-icon">✅</span>
          <div>
            <p class="sbox-label">WO Finish</p>
            <p class="sbox-value">{{ finishedWOs.length }} Tiket</p>
          </div>
        </div>
        <div class="sbox blue">
          <span class="sbox-icon">🔧</span>
          <div>
            <p class="sbox-label">Maintenance Selesai</p>
            <p class="sbox-value">{{ maintenanceHistory.length }} Riwayat</p>
          </div>
        </div>
        <div class="sbox purple">
          <span class="sbox-icon">💰</span>
          <div>
            <p class="sbox-label">Total Biaya Maintenance</p>
            <p class="sbox-value">Rp {{ formatNumber(totalMaintenanceCost) }}</p>
          </div>
        </div>
        <div class="sbox orange">
          <span class="sbox-icon">🔀</span>
          <div>
            <p class="sbox-label">Total Mutasi Aset</p>
            <p class="sbox-value">{{ mutations.length }} Mutasi</p>
          </div>
        </div>
      </div>

      <!-- Search -->
      <div class="search-row">
        <input v-model="searchFilter" placeholder="🔍 Cari..." class="search-input" />
      </div>

      <!-- Section 1: Finished Work Orders -->
      <div class="card-panel">
        <h2 class="section-title">✅ Work Order Selesai</h2>
        <div class="table-responsive">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Lokasi</th>
                <th>Aset ID</th>
                <th>Prioritas</th>
                <th>Deskripsi</th>
                <th>Tindakan Perbaikan</th>
                <th>Biaya (Rp)</th>
                <th>Status</th>
                <th>Tanggal Selesai</th>
                <th v-if="canManage">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in filteredWOs" :key="wo.id">
                <td><span class="wo-id">#WO-{{ wo.id }}</span></td>
                <td>📍 {{ wo.location || '—' }}</td>
                <td>Aset #{{ wo.asset_id }}</td>
                <td><StatusBadge :status="wo.priority || 'Medium'" /></td>
                <td class="desc-cell" :title="wo.description">{{ wo.description }}</td>
                <td class="desc-cell" :title="wo.action_taken">{{ wo.action_taken || '—' }}</td>
                <td>Rp {{ formatNumber(wo.cost || 0) }}</td>
                <td><StatusBadge :status="wo.status" /></td>
                <td class="time-col">{{ formatDate(wo.closed_at) }}</td>
                <td v-if="canManage" class="actions-cell">
                  <button class="icon-btn edit-btn" @click="openEditWoModal(wo)" title="Edit">✏️ Edit</button>
                  <button class="icon-btn delete-btn" @click="deleteWo(wo)" title="Hapus">🗑️ Hapus</button>
                </td>
              </tr>
              <tr v-if="filteredWOs.length === 0">
                <td :colspan="canManage ? 10 : 9" class="empty-state">Tidak ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Section 2: Maintenance History -->
      <div class="card-panel" style="margin-top: 24px;">
        <h2 class="section-title">🔧 Maintenance Selesai</h2>
        <div class="table-responsive">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Aset ID</th>
                <th>Tindakan Perawatan</th>
                <th>Biaya (Rp)</th>
                <th>Tanggal Pengerjaan</th>
                <th v-if="canManage">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mh in filteredMH" :key="mh.id">
                <td><span class="wo-id">#MH-{{ mh.id }}</span></td>
                <td>Aset #{{ mh.asset_id }}</td>
                <td class="desc-cell" :title="mh.action_taken">{{ mh.action_taken }}</td>
                <td>Rp {{ formatNumber(mh.cost || 0) }}</td>
                <td class="time-col">{{ formatDate(mh.created_at) }}</td>
                <td v-if="canManage" class="actions-cell">
                  <button class="icon-btn edit-btn" @click="openEditMhModal(mh)" title="Edit">✏️ Edit</button>
                  <button class="icon-btn delete-btn" @click="deleteMh(mh)" title="Hapus">🗑️ Hapus</button>
                </td>
              </tr>
              <tr v-if="filteredMH.length === 0">
                <td :colspan="canManage ? 6 : 5" class="empty-state">Tidak ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Section 3: Mutation History -->
      <div class="card-panel" style="margin-top: 24px;">
        <h2 class="section-title">🔀 Riwayat Mutasi Aset</h2>
        <div class="table-responsive">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Aset ID</th>
                <th>Lokasi Sebelumnya</th>
                <th>Lokasi Baru</th>
                <th>PIC / Departemen</th>
                <th>Alasan Mutasi</th>
                <th>Tanggal Mutasi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mut in filteredMutations" :key="mut.id">
                <td><span class="wo-id">#MUT-{{ mut.id }}</span></td>
                <td>Aset #{{ mut.asset_id }}</td>
                <td>📍 {{ mut.previous_location || '—' }}</td>
                <td><span class="location-new">📍 {{ mut.new_location || '—' }}</span></td>
                <td>🏢 {{ mut.new_pic || '—' }}</td>
                <td class="desc-cell" :title="mut.reason">{{ mut.reason || '—' }}</td>
                <td class="time-col">{{ formatDate(mut.mutation_date) }}</td>
              </tr>
              <tr v-if="filteredMutations.length === 0">
                <td colspan="7" class="empty-state">Belum ada riwayat mutasi aset.</td>
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
    </div>

    <!-- Modals placed OUTSIDE .page-container so print media query hiding .page-container leaves report modal visible -->
    <ModalDialog :show="showWoModal" title="✏️ Edit WO" @close="showWoModal = false">
      <form @submit.prevent="submitEditWo" class="modal-form" v-if="selectedWo">
        <p class="modal-info"><strong>WO #{{ selectedWo.id }}:</strong> 📍 {{ selectedWo.location }}</p>
        <label>
          <span>Deskripsi</span>
          <textarea v-model="editWoDesc" rows="3" required></textarea>
        </label>
        <label>
          <span>Tindakan Perbaikan</span>
          <textarea v-model="editWoAction" rows="3" required></textarea>
        </label>
        <label>
          <span>Biaya (Rp)</span>
          <input v-model.number="editWoCost" type="number" min="0" required />
        </label>

        <button type="submit" class="submit-modal-btn">Simpan</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showMhModal" title="✏️ Edit Maintenance" @close="showMhModal = false">
      <form @submit.prevent="submitEditMh" class="modal-form" v-if="selectedMh">
        <p class="modal-info"><strong>Maintenance #MH-{{ selectedMh.id }} (Aset #{{ selectedMh.asset_id }})</strong></p>
        <label>
          <span>Tindakan Perawatan</span>
          <textarea v-model="editMhAction" rows="3" required></textarea>
        </label>
        <label>
          <span>Biaya (Rp)</span>
          <input v-model.number="editMhCost" type="number" min="0" required />
        </label>

        <button type="submit" class="submit-modal-btn">Simpan</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showReportModal" title="🖨️ Laporan Activity Log" maxWidth="960px" @close="showReportModal = false">
      <div class="monthly-report-printable" id="printableReportDocument">
        <div class="report-header">
          <h2>🏨 LAPORAN BULANAN AUDIT TRAIL & AKTIVITAS</h2>
          <p class="report-sub">Sistem AsetKu Hotel — Periode: {{ reportMonthYear }}</p>
          <hr class="report-divider" />
        </div>

        <div class="report-summary-boxes">
          <div class="rbox success">
            <span>WO Finish</span>
            <strong>{{ finishedWOs.length }} Tiket</strong>
          </div>
          <div class="rbox blue">
            <span>Maintenance Selesai</span>
            <strong>{{ maintenanceHistory.length }} Riwayat</strong>
          </div>
          <div class="rbox danger">
            <span>Total Biaya Maintenance</span>
            <strong>Rp {{ formatNumber(totalMaintenanceCost) }}</strong>
          </div>
        </div>

        <h3 class="report-section-heading">1. Rekapitulasi Work Order Selesai</h3>
        <div class="report-table-wrapper">
          <table class="report-table">
            <thead>
              <tr>
                <th class="col-id">ID</th>
                <th class="col-loc">Lokasi</th>
                <th class="col-prio">Prioritas</th>
                <th class="col-desc">Deskripsi</th>
                <th class="col-desc">Tindakan Perbaikan</th>
                <th class="col-stat">Status</th>
                <th class="col-cost">Biaya (Rp)</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in finishedWOs" :key="wo.id">
                <td class="col-id">#WO-{{ wo.id }}</td>
                <td class="col-loc">{{ wo.location || '—' }}</td>
                <td class="col-prio">{{ wo.priority }}</td>
                <td class="col-desc">{{ wo.description }}</td>
                <td class="col-desc">{{ wo.action_taken || '—' }}</td>
                <td class="col-stat">{{ wo.status }}</td>
                <td class="col-cost">Rp {{ formatNumber(wo.cost || 0) }}</td>
              </tr>
              <tr v-if="finishedWOs.length === 0">
                <td colspan="7" class="empty-state">Tidak ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>

        <h3 class="report-section-heading" style="margin-top: 24px;">2. Rekapitulasi Maintenance Selesai</h3>
        <div class="report-table-wrapper">
          <table class="report-table">
            <thead>
              <tr>
                <th class="col-id">ID</th>
                <th class="col-prio">Aset ID</th>
                <th class="col-desc">Tindakan Perawatan</th>
                <th class="col-cost">Biaya (Rp)</th>
                <th class="col-loc">Tanggal</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mh in maintenanceHistory" :key="mh.id">
                <td class="col-id">#MH-{{ mh.id }}</td>
                <td class="col-prio">Aset #{{ mh.asset_id }}</td>
                <td class="col-desc">{{ mh.action_taken }}</td>
                <td class="col-cost">Rp {{ formatNumber(mh.cost || 0) }}</td>
                <td class="col-loc">{{ formatDate(mh.created_at) }}</td>
              </tr>
              <tr v-if="maintenanceHistory.length === 0">
                <td colspan="5" class="empty-state">Tidak ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="report-actions no-print">
        <button class="submit-modal-btn excel-btn" @click="exportToExcel">
          📊 Export ke Excel (.xlsx)
        </button>
        <button class="submit-modal-btn print-btn" @click="printReport">
          🖨️ Cetak Dokumen Laporan (PDF / Print)
        </button>
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
const canManage = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')

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

const searchFilter = ref('')
const finishedWOs = ref([])
const maintenanceHistory = ref([])
const mutations = ref([])
const isLoading = ref(false)

const showReportModal = ref(false)
const reportMonthYear = ref(new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }))

const showWoModal = ref(false)
const selectedWo = ref(null)
const editWoDesc = ref('')
const editWoAction = ref('')
const editWoCost = ref(0)

const showMhModal = ref(false)
const selectedMh = ref(null)
const editMhAction = ref('')
const editMhCost = ref(0)

const totalMaintenanceCost = computed(() =>
  maintenanceHistory.value.reduce((sum, mh) => sum + (mh.cost || 0), 0)
)

const totalWoCost = computed(() =>
  finishedWOs.value.reduce((sum, wo) => sum + (wo.cost || 0), 0)
)

const filteredWOs = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return finishedWOs.value
  return finishedWOs.value.filter(wo =>
    (wo.description && wo.description.toLowerCase().includes(q)) ||
    (wo.action_taken && wo.action_taken.toLowerCase().includes(q)) ||
    (wo.location && wo.location.toLowerCase().includes(q)) ||
    String(wo.asset_id).includes(q)
  )
})

const filteredMH = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return maintenanceHistory.value
  return maintenanceHistory.value.filter(mh =>
    (mh.action_taken && mh.action_taken.toLowerCase().includes(q)) ||
    String(mh.asset_id).includes(q)
  )
})

const filteredMutations = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return mutations.value
  return mutations.value.filter(mut =>
    (mut.new_location && mut.new_location.toLowerCase().includes(q)) ||
    (mut.previous_location && mut.previous_location.toLowerCase().includes(q)) ||
    (mut.new_pic && mut.new_pic.toLowerCase().includes(q)) ||
    (mut.reason && mut.reason.toLowerCase().includes(q)) ||
    String(mut.asset_id).includes(q)
  )
})

function formatNumber(num) {
  return (num || 0).toLocaleString('id-ID')
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  try {
    return new Date(dateStr).toLocaleString('id-ID', {
      day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

async function fetchLogs() {
  isLoading.value = true
  try {
    const res = await api.get('/activitylogs')
    if (res.data?.data) {
      finishedWOs.value = res.data.data.work_orders || []
      maintenanceHistory.value = res.data.data.maintenance_history || []
      mutations.value = res.data.data.mutations || []
    }
  } catch (e) {
    console.error('Failed to fetch activity logs:', e)
  } finally {
    isLoading.value = false
  }
}

function openEditWoModal(wo) {
  selectedWo.value = wo
  editWoDesc.value = wo.description || ''
  editWoAction.value = wo.action_taken || ''
  editWoCost.value = wo.cost || 0
  showWoModal.value = true
}

async function submitEditWo() {
  if (!selectedWo.value) return
  try {
    await api.post('/workorders/edit', {
      wo_id: selectedWo.value.id,
      description: editWoDesc.value,
      action_taken: editWoAction.value,
      cost: editWoCost.value
    })
    showWoModal.value = false
    notify('WO diperbarui!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal mengubah WO: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteWo(wo) {
  try {
    await api.post('/workorders/delete', { wo_id: wo.id })
    finishedWOs.value = finishedWOs.value.filter(w => w.id !== wo.id)
    notify('WO dihapus!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal menghapus WO: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openEditMhModal(mh) {
  selectedMh.value = mh
  editMhAction.value = mh.action_taken || ''
  editMhCost.value = mh.cost || 0
  showMhModal.value = true
}

async function submitEditMh() {
  if (!selectedMh.value) return
  try {
    await api.post('/maintenance/history/edit', {
      history_id: selectedMh.value.id,
      action_taken: editMhAction.value,
      cost: editMhCost.value
    })
    showMhModal.value = false
    notify('Riwayat maintenance diperbarui!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal mengubah riwayat: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteMh(mh) {
  try {
    await api.post('/maintenance/history/delete', { history_id: mh.id })
    maintenanceHistory.value = maintenanceHistory.value.filter(m => m.id !== mh.id)
    notify('Riwayat maintenance dihapus!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal menghapus riwayat: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function printReport() {
  window.print()
}

function exportToExcel() {
  const monthName = reportMonthYear.value.replace(/\s+/g, '_')
  const fileName = `Laporan_ActivityLog_Hotel_${monthName}.xls`

  let htmlTable = `
    <html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
    <head>
      <meta charset="utf-8">
      <!--[if gte mso 9]>
      <xml>
        <x:ExcelWorkbook>
          <x:ExcelWorksheets>
            <x:ExcelWorksheet>
              <x:Name>Activity Log & Audit Trail</x:Name>
              <x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions>
            </x:ExcelWorksheet>
          </x:ExcelWorksheets>
        </x:ExcelWorkbook>
      </xml>
      <![endif]-->
      <style>
        th { background-color: #2563eb; color: #ffffff; font-weight: bold; border: 1px solid #000000; text-align: center; }
        td { border: 1px solid #cccccc; vertical-align: middle; }
        .summary-header { font-size: 14pt; font-weight: bold; color: #0f172a; }
        .section-header { font-size: 12pt; font-weight: bold; color: #1e293b; background-color: #e2e8f0; }
        .total-row { font-weight: bold; background-color: #f1f5f9; }
      </style>
    </head>
    <body>
      <h2 class="summary-header">LAPORAN BULANAN AUDIT TRAIL & AKTIVITAS HOTEL</h2>
      <p><b>Sistem AsetKu Hotel</b> — Periode: ${reportMonthYear.value}</p>
      <br/>
      <h3 class="section-header">1. REKAPITULASI WORK ORDER SELESAI (FINISH)</h3>
      <table border="1" cellspacing="0" cellpadding="6">
        <thead>
          <tr>
            <th>ID Work Order</th>
            <th>Lokasi / Kamar</th>
            <th>Prioritas</th>
            <th>Deskripsi Kerusakan</th>
            <th>Tindakan Perbaikan</th>
            <th>Status Tiket</th>
            <th>Biaya Perbaikan (Rp)</th>
          </tr>
        </thead>
        <tbody>
  `

  finishedWOs.value.forEach(wo => {
    htmlTable += `
      <tr>
        <td>#WO-${wo.id}</td>
        <td>${wo.location || ''}</td>
        <td>${wo.priority || ''}</td>
        <td>${wo.description || ''}</td>
        <td>${wo.action_taken || ''}</td>
        <td>${wo.status || ''}</td>
        <td align="right">${wo.cost || 0}</td>
      </tr>
    `
  })

  htmlTable += `
        <tr class="total-row">
          <td colspan="6" align="right"><b>TOTAL BIAYA WO SELESAI:</b></td>
          <td align="right"><b>${totalWoCost.value}</b></td>
        </tr>
      </tbody>
      </table>
      <br/><br/>
      <h3 class="section-header">2. REKAPITULASI RIWAYAT MAINTENANCE SELESAI</h3>
      <table border="1" cellspacing="0" cellpadding="6">
        <thead>
          <tr>
            <th>ID History</th>
            <th>ID Aset</th>
            <th>Tindakan Perawatan yang Dilakukan</th>
            <th>Biaya Maintenance (Rp)</th>
            <th>Tanggal Pengerjaan</th>
          </tr>
        </thead>
        <tbody>
  `

  maintenanceHistory.value.forEach(mh => {
    htmlTable += `
      <tr>
        <td>#MH-${mh.id}</td>
        <td>Aset #${mh.asset_id}</td>
        <td>${mh.action_taken || ''}</td>
        <td align="right">${mh.cost || 0}</td>
        <td>${formatDate(mh.created_at)}</td>
      </tr>
    `
  })

  htmlTable += `
        <tr class="total-row">
          <td colspan="3" align="right"><b>TOTAL BIAYA MAINTENANCE:</b></td>
          <td align="right"><b>${totalMaintenanceCost.value}</b></td>
          <td></td>
        </tr>
      </tbody>
      </table>
    </body>
    </html>
  `

  const blob = new Blob([htmlTable], { type: 'application/vnd.ms-excel;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', fileName)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.page-container {
  max-width: 1240px;
  margin: 0 auto;
  padding: 32px 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-action-group {
  display: flex;
  gap: 12px;
}

.print-report-btn {
  background: #0284c7;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(2, 132, 199, 0.3);
}

.export-excel-btn {
  background: #16a34a;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(22, 163, 74, 0.3);
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

.summary-row {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.sbox {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 20px;
  border-radius: 16px;
  background: #fff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.04);
}

.sbox.green { border-left: 4px solid #16a34a; }
.sbox.blue  { border-left: 4px solid #2563eb; }
.sbox.purple { border-left: 4px solid #7c3aed; }
.sbox.orange { border-left: 4px solid #d97706; }

.location-new {
  color: #16a34a;
  font-weight: 700;
}

.sbox-icon {
  font-size: 1.8rem;
  line-height: 1;
}

.sbox-label {
  margin: 0 0 4px;
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
}

.sbox-value {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

.search-row {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  max-width: 480px;
  padding: 11px 18px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  font-size: 0.9rem;
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.section-title {
  margin: 0 0 20px;
  font-size: 1.1rem;
  color: #0f172a;
  font-weight: 800;
}

.table-responsive {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 700px;
}

th {
  text-align: left;
  padding: 12px 14px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.82rem;
  font-weight: 700;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

td {
  padding: 13px 14px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 0.88rem;
  vertical-align: middle;
}

.wo-id {
  font-family: monospace;
  font-weight: 800;
  color: #0f172a;
  font-size: 0.85rem;
}

.desc-cell {
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.time-col {
  color: #64748b;
  font-size: 0.82rem;
  white-space: nowrap;
}

.actions-cell {
  display: flex;
  gap: 6px;
  white-space: nowrap;
}

.icon-btn {
  border: 1px solid #cbd5e1;
  background: white;
  padding: 5px 10px;
  border-radius: 8px;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
}

.edit-btn { color: #2563eb; border-color: #93c5fd; }
.delete-btn { color: #dc2626; border-color: #fca5a5; }

.modal-form {
  display: grid;
  gap: 14px;
}

.modal-form label {
  display: grid;
  gap: 6px;
  font-weight: 600;
  color: #1e293b;
  font-size: 0.9rem;
}

.modal-form input, .modal-form textarea {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  font-size: 0.9rem;
}

.submit-modal-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  margin-top: 6px;
}

.modal-info {
  background: #f8fafc;
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 0.9rem;
  color: #334155;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: #94a3b8;
  font-size: 0.9rem;
}

/* Printable Monthly Activity Log Styling */
.monthly-report-printable {
  background: #ffffff;
  padding: 20px;
  color: #0f172a;
}

.report-header {
  text-align: center;
  margin-bottom: 20px;
}

.report-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: #0f172a;
}

.report-sub {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 0.9rem;
}

.report-divider {
  border: 0;
  border-top: 2px solid #0f172a;
  margin: 12px 0 16px;
}

.report-summary-boxes {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.rbox {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.rbox span { font-size: 0.75rem; color: #475569; font-weight: 600; }
.rbox strong { font-size: 1.1rem; color: #0f172a; font-weight: 800; }
.rbox.success { border-color: #86efac; background: #f0fdf4; }
.rbox.blue { border-color: #93c5fd; background: #eff6ff; }
.rbox.danger { border-color: #fca5a5; background: #fef2f2; }

.report-section-heading {
  font-size: 1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 10px;
}

.report-table-wrapper {
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.report-table th {
  background: #f1f5f9;
  color: #0f172a;
  font-weight: 700;
  border: 1px solid #94a3b8;
  padding: 8px 10px;
  text-align: left;
}

.report-table td {
  border: 1px solid #cbd5e1;
  padding: 8px 10px;
  color: #334155;
  vertical-align: middle;
}

.report-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.excel-btn { background: #16a34a; }
.print-btn { background: #0284c7; }
</style>
