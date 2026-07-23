<template>
  <div class="view-wrapper">
    <div class="page-container">
      <div class="page-header">
        <div>
          <p class="eyebrow">Work Order</p>
          <h1>🔧 Work Order Hotel</h1>
          <p class="subtitle">Pelaporan & perbaikan kerusakan.</p>
        </div>

        <div class="header-action-group">
          <button class="print-report-btn" @click="showReportModal = true" title="Prinjau & Export Laporan Bulanan">
            📄 Laporan & Export
          </button>
          <button class="primary-btn" @click="showCreateModal = true">
            🚨 Buat WO
          </button>
        </div>
      </div>

      <!-- Status Tabs -->
      <div class="status-tabs">
        <button v-for="tab in tabs" :key="tab.id" :class="['tab-btn', { active: activeTab === tab.id }]" @click="activeTab = tab.id">
          {{ tab.label }}
          <span class="tab-count">{{ getTabCount(tab.id) }}</span>
        </button>
      </div>

      <!-- Work Orders Table -->
      <div class="card-panel">
        <div class="table-responsive wo-table-wrapper">
          <table>
            <thead>
              <tr>
                <th>ID WO</th>
                <th>Lokasi / Kamar</th>
                <th>Aset ID</th>
                <th>Prioritas</th>
                <th>Deskripsi Kerusakan</th>
                <th>Teknisi (Engineer)</th>
                <th>Status</th>
                <th>Aksi Tiket</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in filteredWorkOrders" :key="wo.id">
                <td class="nowrap-cell"><span class="wo-id">#WO-{{ wo.id }}</span></td>
                <td class="nowrap-cell">📍 {{ wo.location || 'Kamar / Area Hotel' }}</td>
                <td class="nowrap-cell">Aset #{{ wo.asset_id }}</td>
                <td class="nowrap-cell"><StatusBadge :status="wo.priority || 'Medium'" /></td>
                <td class="desc-cell" :title="wo.description"><span class="desc-text">{{ wo.description }}</span></td>
                <td class="nowrap-cell">
                  <span v-if="wo.engineer_id > 0" class="engineer-chip">👷 Teknisi #{{ wo.engineer_id }}</span>
                  <span v-else class="unassigned-chip">Belum Ditugaskan</span>
                </td>
                <td class="nowrap-cell"><StatusBadge :status="wo.status" /></td>
                <td class="actions-cell">
                  <button v-if="canAssign && wo.status === 'Open'" class="icon-btn assign-btn" @click="openAssignModal(wo)" title="Assign Worker">
                    👷 Assign
                  </button>
                  <button v-if="canUpdateProgress && (wo.status === 'In Progress' || wo.status === 'Open')" class="icon-btn progress-btn" @click="openUpdateModal(wo)" title="Update Progres">
                    📝 Progres
                  </button>
                  <button v-if="canManageOrder && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn close-btn" @click="closeOrder(wo)" title="Selesaikan Work Order">
                    ✅ Selesai
                  </button>

                  <!-- External Staff, HOD, and Management can cancel open work orders -->
                <button v-if="wo.status === 'Open'" class="icon-btn cancel-btn" @click="cancelOrder(wo)" title="Batal (Cancel Work Order)">
                  🚫 Batal
                </button>

                <!-- HOD, Admin, and Management can delete work orders -->
                <button v-if="canDelete" class="icon-btn delete-btn" @click="deleteOrder(wo)" title="Hapus Work Order Permanen">
                  🗑️ Hapus
                </button>
              </td>
            </tr>
            <tr v-if="filteredWorkOrders.length === 0">
              <td colspan="8" class="empty-state">Tidak ada Work Order pada kategori ini.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Custom UI Toast Notification (No Browser Alert) -->
    <transition name="fade">
      <div v-if="showToast" :class="['custom-ui-toast', toastType]">
        <span class="toast-icon">{{ toastType === 'success' ? '✅' : '⚠️' }}</span>
        <span class="toast-text">{{ toastMsg }}</span>
        <button class="toast-close" @click="showToast = false">✕</button>
      </div>
    </transition>
    </div>

    <!-- Modals placed OUTSIDE .page-container so hiding .page-container during print never affects modal report -->
    <ModalDialog :show="showCreateModal" title="🚨 Ajukan Laporan Kerusakan Aset" @close="showCreateModal = false">
      <form @submit.prevent="submitWorkOrder" class="modal-form">
        <label>
          <span>ID Aset</span>
          <input v-model.number="formWo.asset_id" type="number" placeholder="Masukkan ID Aset" required />
        </label>
        <label>
          <span>Lokasi / Kamar Tempat Kerusakan</span>
          <input v-model="formWo.location" placeholder="Contoh: Kamar 301, Kitchen Dapur, Lift Gedung A" required />
        </label>
        <label>
          <span>Tingkat Prioritas (Priority)</span>
          <select v-model="formWo.priority">
            <option value="Low">Low (Kerusakan Ringan / Tidak Mendesak)</option>
            <option value="Medium">Medium (Kerusakan Standar)</option>
            <option value="High">High (Perlu Penanganan Cepat)</option>
            <option value="Emergency">Emergency (Darurat / Operasional Terganggu)</option>
          </select>
        </label>
        <label>
          <span>Deskripsi Rincian Kerusakan</span>
          <textarea v-model="formWo.description" rows="3" placeholder="Jelaskan detail gejala kerusakan aset..." required></textarea>
        </label>

        <button type="submit" class="submit-modal-btn">Kirim Tiket Perbaikan</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showAssignModal" title="👷 Penugasan Teknisi (Assign Worker)" @close="showAssignModal = false">
      <form @submit.prevent="submitAssign" class="modal-form" v-if="selectedWo">
        <p class="modal-info"><strong>WO #{{ selectedWo.id }}:</strong> {{ selectedWo.description }} (📍 {{ selectedWo.location }})</p>
        <label>
          <span>Pilih Teknisi (Staff Engineer)</span>
          <select v-model.number="assignEngineerId" required>
            <option value="101">Teknisi Deni (HVAC & AC)</option>
            <option value="102">Teknisi Budi (Plumbing & Sanitasi)</option>
            <option value="103">Teknisi Agus (Elektronik & Kelistrikan)</option>
            <option value="104">Teknisi Eko (General Maintenance)</option>
          </select>
        </label>

        <button type="submit" class="submit-modal-btn">Tugaskan Teknisi</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showUpdateModal" title="📝 Update Status Pengerjaan Teknisi" @close="showUpdateModal = false">
      <form @submit.prevent="submitUpdateStatus" class="modal-form" v-if="selectedWo">
        <label>
          <span>Status Baru</span>
          <select v-model="updateStatus">
            <option value="In Progress">In Progress (Sedang Dikerjakan)</option>
            <option value="Under Review">Under Review (Selesai Dikerjakan & Menunggu Review HOD)</option>
            <option value="Completed">Completed (Selesai Perbaikan)</option>
          </select>
        </label>

        <label>
          <span>Tindakan Perbaikan yang Dilakukan</span>
          <textarea v-model="updateActionTaken" rows="3" placeholder="Misal: Ganti kapasitor kompresor AC dan isi ulang freon R32."></textarea>
        </label>

        <label>
          <span>Estimasi Biaya Perbaikan (Rp)</span>
          <input v-model.number="updateCost" type="number" placeholder="Contoh: 150000" />
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Progres</button>
      </form>
    </ModalDialog>

    <!-- Printable Report Modal (Outside .page-container) -->
    <ModalDialog :show="showReportModal" title="🖨️ Laporan Bulanan Work Order & Maintenance" maxWidth="920px" @close="showReportModal = false">
      <div class="monthly-report-printable" id="printableReportDocument">
        <div class="report-header">
          <h2>🏨 LAPORAN BULANAN MANAJEMEN ASET & WORK ORDER</h2>
          <p class="report-sub">Sistem AsetKu Hotel — Periode: {{ reportMonthYear }}</p>
          <hr class="report-divider" />
        </div>

        <div class="report-summary-boxes">
          <div class="rbox">
            <span>Total Tiket</span>
            <strong>{{ workOrders.length }} Tiket</strong>
          </div>
          <div class="rbox success">
            <span>Finish / Selesai</span>
            <strong>{{ countStatus('Finish') + countStatus('Completed') + countStatus('Closed') }} Tiket</strong>
          </div>
          <div class="rbox warning">
            <span>Dalam Proses</span>
            <strong>{{ countStatus('In Progress') + countStatus('Open') }} Tiket</strong>
          </div>
          <div class="rbox danger">
            <span>Total Biaya Perbaikan</span>
            <strong>Rp {{ formatNumber(totalReportCost) }}</strong>
          </div>
        </div>

        <div class="report-table-wrapper">
          <table class="report-table">
            <thead>
              <tr>
                <th class="col-id">ID</th>
                <th class="col-loc">Lokasi / Kamar</th>
                <th class="col-prio">Prioritas</th>
                <th class="col-desc">Rincian Kerusakan</th>
                <th class="col-stat">Status</th>
                <th class="col-cost">Biaya (Rp)</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in workOrders" :key="wo.id">
                <td class="col-id">#WO-{{ wo.id }}</td>
                <td class="col-loc">{{ wo.location }}</td>
                <td class="col-prio">{{ wo.priority }}</td>
                <td class="col-desc">{{ wo.description }}</td>
                <td class="col-stat">{{ wo.status }}</td>
                <td class="col-cost">Rp {{ formatNumber(wo.cost || 0) }}</td>
              </tr>
              <tr v-if="workOrders.length === 0">
                <td colspan="6" class="empty-state">Belum ada rekapitulasi data Work Order bulan ini.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="report-actions no-print">
        <button class="submit-modal-btn excel-btn" @click="exportToExcel">
          📊 Export ke Excel (.xlsx)
        </button>
        <button class="submit-modal-btn print-btn" @click="printMonthlyReport">
          🖨️ Cetak Dokumen Laporan (PDF / Print)
        </button>
      </div>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const route = useRoute()

const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')

const canAssign = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canManageOrder = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canUpdateProgress = computed(() => userRole.value === 'engineer' || userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canDelete = computed(() => userRole.value === 'hod' || userRole.value === 'admin' || userRole.value === 'management')

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

const activeTab = ref('all')
const tabs = [
  { id: 'all', label: 'Semua Tiket' },
  { id: 'open', label: 'Open' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'under_review', label: 'Under Review' },
  { id: 'completed', label: 'Finish / Completed' }
]

const workOrders = ref([])
const isLoading = ref(false)

const showCreateModal = ref(false)
const formWo = ref({ asset_id: 1, location: '', priority: 'Medium', description: '' })

const showAssignModal = ref(false)
const selectedWo = ref(null)
const assignEngineerId = ref(101)

const showUpdateModal = ref(false)
const updateStatus = ref('In Progress')
const updateActionTaken = ref('')
const updateCost = ref(0)

const showReportModal = ref(false)
const reportMonthYear = ref(new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }))

const filteredWorkOrders = computed(() => {
  if (activeTab.value === 'all') return workOrders.value
  if (activeTab.value === 'open') return workOrders.value.filter(w => w.status === 'Open')
  if (activeTab.value === 'in_progress') return workOrders.value.filter(w => w.status === 'In Progress')
  if (activeTab.value === 'under_review') return workOrders.value.filter(w => w.status === 'Under Review')
  if (activeTab.value === 'completed') return workOrders.value.filter(w => w.status === 'Finish' || w.status === 'Completed' || w.status === 'Closed')
  return workOrders.value
})

const totalReportCost = computed(() => {
  return workOrders.value.reduce((sum, w) => sum + (w.cost || 0), 0)
})

function getTabCount(tabId) {
  if (tabId === 'all') return workOrders.value.length
  if (tabId === 'open') return workOrders.value.filter(w => w.status === 'Open').length
  if (tabId === 'in_progress') return workOrders.value.filter(w => w.status === 'In Progress').length
  if (tabId === 'under_review') return workOrders.value.filter(w => w.status === 'Under Review').length
  if (tabId === 'completed') return workOrders.value.filter(w => w.status === 'Finish' || w.status === 'Completed' || w.status === 'Closed').length
  return 0
}

function countStatus(st) {
  return workOrders.value.filter(w => w.status === st).length
}

function formatNumber(num) {
  return (num || 0).toLocaleString('id-ID')
}

async function fetchWorkOrders() {
  isLoading.value = true
  try {
    const res = await api.get('/workorders')
    if (res.data?.data && Array.isArray(res.data.data)) {
      workOrders.value = res.data.data
    } else {
      workOrders.value = []
    }
  } catch (e) {
    console.error('Failed to fetch work orders from DB:', e)
    workOrders.value = []
  } finally {
    isLoading.value = false
  }
}

async function submitWorkOrder() {
  try {
    await api.post('/workorders', {
      asset_id: formWo.value.asset_id,
      location: formWo.value.location || 'Area Hotel',
      priority: formWo.value.priority,
      description: formWo.value.description
    })
    showCreateModal.value = false
    notify('Tiket Work Order / Laporan Kerusakan berhasil diajukan ke database!', 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal mengajukan Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openAssignModal(wo) {
  selectedWo.value = wo
  assignEngineerId.value = 101
  showAssignModal.value = true
}

async function submitAssign() {
  if (!selectedWo.value) return
  try {
    await api.post('/workorders/assign', {
      wo_id: selectedWo.value.id,
      engineer_id: assignEngineerId.value
    })
    showAssignModal.value = false
    notify(`Teknisi #${assignEngineerId.value} berhasil ditugaskan ke WO #${selectedWo.value.id}!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal menugaskan teknisi: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openUpdateModal(wo) {
  selectedWo.value = wo
  updateStatus.value = wo.status === 'Open' ? 'In Progress' : wo.status
  updateActionTaken.value = ''
  updateCost.value = 0
  showUpdateModal.value = true
}

async function submitUpdateStatus() {
  if (!selectedWo.value) return
  try {
    await api.post('/workorders/status', {
      wo_id: selectedWo.value.id,
      status: updateStatus.value,
      action_taken: updateActionTaken.value,
      cost: updateCost.value
    })
    showUpdateModal.value = false
    notify(`Status Work Order #${selectedWo.value.id} diperbarui menjadi "${updateStatus.value}"!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal memperbarui status: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function closeOrder(wo) {
  try {
    await api.post('/workorders/close', { wo_id: wo.id })
    notify(`Work Order #${wo.id} berhasil diselesaikan (Finish)!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal menutup Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function cancelOrder(wo) {
  try {
    await api.post('/workorders/cancel', { wo_id: wo.id })
    notify(`Work Order #${wo.id} dibatalkan (Cancelled)!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal membatalkan Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteOrder(wo) {
  try {
    await api.post('/workorders/delete', { wo_id: wo.id })
    workOrders.value = workOrders.value.filter(w => w.id !== wo.id)
    notify(`Work Order #${wo.id} berhasil dihapus permanen!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal menghapus Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function printMonthlyReport() {
  window.print()
}

function exportToExcel() {
  const monthName = reportMonthYear.value.replace(/\s+/g, '_')
  const fileName = `Laporan_WorkOrder_Hotel_${monthName}.xls`

  let htmlTable = `
    <html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
    <head>
      <meta charset="utf-8">
      <!--[if gte mso 9]>
      <xml>
        <x:ExcelWorkbook>
          <x:ExcelWorksheets>
            <x:ExcelWorksheet>
              <x:Name>Laporan Work Order</x:Name>
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
        .total-row { font-weight: bold; background-color: #f1f5f9; }
      </style>
    </head>
    <body>
      <h2 class="summary-header">LAPORAN BULANAN MANAJEMEN ASET & WORK ORDER HOTEL</h2>
      <p><b>Sistem AsetKu Hotel</b> — Periode: ${reportMonthYear.value}</p>
      <br/>
      <table border="1" cellspacing="0" cellpadding="6">
        <thead>
          <tr>
            <th>ID Work Order</th>
            <th>Lokasi / Kamar</th>
            <th>Prioritas</th>
            <th>Deskripsi Kerusakan</th>
            <th>Status Tiket</th>
            <th>Biaya Perbaikan (Rp)</th>
          </tr>
        </thead>
        <tbody>
  `

  workOrders.value.forEach(wo => {
    htmlTable += `
      <tr>
        <td>#WO-${wo.id}</td>
        <td>${wo.location || ''}</td>
        <td>${wo.priority || ''}</td>
        <td>${wo.description || ''}</td>
        <td>${wo.status || ''}</td>
        <td align="right">${wo.cost || 0}</td>
      </tr>
    `
  })

  htmlTable += `
        <tr class="total-row">
          <td colspan="5" align="right"><b>TOTAL BIAYA PERBAIKAN:</b></td>
          <td align="right"><b>${totalReportCost.value}</b></td>
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
  if (route.query.assetId) {
    formWo.value.asset_id = parseInt(route.query.assetId)
    formWo.value.location = route.query.location || ''
    showCreateModal.value = true
  }
  fetchWorkOrders()
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
  background: #dc2626;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(220, 38, 38, 0.3);
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

.status-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 4px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tab-btn {
  background: transparent;
  border: none;
  padding: 8px 16px;
  font-size: 0.9rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  border-bottom: 3px solid transparent;
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.tab-btn.active {
  color: #2563eb;
  border-bottom-color: #2563eb;
}

.tab-count {
  background: #e2e8f0;
  color: #334155;
  padding: 2px 6px;
  border-radius: 999px;
  font-size: 0.75rem;
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.wo-table-wrapper {
  overflow-x: auto;
  max-width: 100%;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 900px;
}

th {
  text-align: left;
  padding: 12px 14px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.85rem;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

td {
  padding: 14px 14px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.9rem;
  color: #334155;
  vertical-align: middle;
}

.nowrap-cell {
  white-space: nowrap;
}

.desc-cell {
  white-space: nowrap;
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.desc-text {
  font-weight: 500;
  white-space: nowrap;
}

.wo-id {
  font-family: monospace;
  font-weight: 800;
  color: #0f172a;
}

.engineer-chip {
  background: #dcfce7;
  color: #15803d;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
}

.unassigned-chip {
  color: #94a3b8;
  font-size: 0.8rem;
  font-style: italic;
  white-space: nowrap;
}

.actions-cell {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.icon-btn {
  border: 1px solid #cbd5e1;
  background: white;
  padding: 6px 10px;
  border-radius: 8px;
  font-size: 0.8rem;
  cursor: pointer;
  white-space: nowrap;
}

.assign-btn { color: #2563eb; border-color: #93c5fd; }
.progress-btn { color: #d97706; border-color: #fcd34d; }
.close-btn { color: #16a34a; border-color: #86efac; }
.cancel-btn { color: #dc2626; border-color: #fca5a5; }

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

.modal-form input, .modal-form select, .modal-form textarea {
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
}

.modal-info {
  background: #f8fafc;
  padding: 10px;
  border-radius: 8px;
  font-size: 0.9rem;
}

.empty-state {
  text-align: center;
  color: #94a3b8;
  padding: 24px;
}

/* Printable Monthly Report Styling */
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
  font-weight: 600;
}

.report-divider {
  border: 0;
  height: 2px;
  background: #0f172a;
  margin: 12px 0 20px;
}

.report-summary-boxes {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.rbox {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.rbox span { font-size: 0.75rem; color: #64748b; font-weight: 600; }
.rbox strong { font-size: 1.1rem; color: #0f172a; margin-top: 4px; }
.rbox.success strong { color: #16a34a; }
.rbox.warning strong { color: #d97706; }
.rbox.danger strong { color: #dc2626; }

.report-table-wrapper {
  width: 100%;
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 24px;
  table-layout: fixed;
}

.report-table th, .report-table td {
  border: 1px solid #cbd5e1;
  padding: 8px 10px;
  font-size: 0.85rem;
  word-break: break-word;
  white-space: normal;
}

.report-table th {
  background: #f1f5f9;
  font-weight: 700;
}

.report-actions {
  margin-top: 20px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.excel-btn {
  background: #16a34a;
}

.print-btn {
  background: #0284c7;
}
</style>
