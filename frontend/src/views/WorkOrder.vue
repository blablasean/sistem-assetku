<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Ticket & Repair Request Management</p>
        <h1>🔧 Manajemen Work Order Hotel</h1>
        <p class="subtitle">Pelaporan perbaikan kerusakan aset berdasarkan lokasi & prioritas, penugasan teknisi, dan pemantauan real-time.</p>
      </div>

      <button class="primary-btn" @click="showCreateModal = true">
        🚨 Ajukan Laporan Kerusakan
      </button>
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
      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>ID WO</th>
              <th>Lokasi / Kamar</th>
              <th>Aset ID / Detail</th>
              <th>Prioritas</th>
              <th>Deskripsi Kerusakan</th>
              <th>Teknisi (Engineer)</th>
              <th>Status</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="wo in filteredWorkOrders" :key="wo.id">
              <td><span class="wo-id">#WO-{{ wo.id }}</span></td>
              <td>📍 {{ wo.location || 'Kamar / Area Hotel' }}</td>
              <td>Aset #{{ wo.asset_id }}</td>
              <td><StatusBadge :status="wo.priority || 'Medium'" /></td>
              <td><span class="desc-text">{{ wo.description }}</span></td>
              <td>
                <span v-if="wo.engineer_id > 0" class="engineer-chip">👷 Teknisi #{{ wo.engineer_id }}</span>
                <span v-else class="unassigned-chip">Belum Ditugaskan</span>
              </td>
              <td><StatusBadge :status="wo.status" /></td>
              <td class="actions-cell">
                <button v-if="canAssign && wo.status === 'Open'" class="icon-btn assign-btn" @click="openAssignModal(wo)" title="Assign Worker">
                  👷 Assign
                </button>
                <button v-if="canUpdateProgress && (wo.status === 'In Progress' || wo.status === 'Open')" class="icon-btn progress-btn" @click="openUpdateModal(wo)" title="Update Progres">
                  📝 Progres
                </button>
                <button v-if="canManageOrder && wo.status !== 'Closed' && wo.status !== 'Cancelled'" class="icon-btn close-btn" @click="closeOrder(wo)" title="Close Work Order">
                  ✅ Close
                </button>
                <button v-if="canManageOrder && wo.status === 'Open'" class="icon-btn cancel-btn" @click="cancelOrder(wo)" title="Cancel Work Order">
                  🚫 Cancel
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

    <!-- Modal Buat Work Order -->
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

    <!-- Modal Assign Engineer -->
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

    <!-- Modal Update Status / Catatan Pengerjaan -->
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const route = useRoute()

const userRole = ref(localStorage.getItem('user_role') || 'external')

const canAssign = computed(() => userRole.value === 'hod' || userRole.value === 'management')
const canManageOrder = computed(() => userRole.value === 'hod' || userRole.value === 'management')
const canUpdateProgress = computed(() => userRole.value === 'engineer' || userRole.value === 'hod' || userRole.value === 'management')

const activeTab = ref('all')
const tabs = [
  { id: 'all', label: 'Semua Tiket' },
  { id: 'open', label: 'Open' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'under_review', label: 'Under Review' },
  { id: 'completed', label: 'Closed / Completed' }
]

const workOrders = ref([
  { id: 1, asset_id: 1, location: 'Kamar 301', priority: 'Emergency', description: 'AC Kamar 301 bocor air dan tidak dingin', status: 'In Progress', requester_id: 3, engineer_id: 101 },
  { id: 2, asset_id: 3, location: 'Kitchen Dapur', priority: 'High', description: 'Chiller Dapur Utama suhu naik menjadi -5°C', status: 'Open', requester_id: 4, engineer_id: 0 },
  { id: 3, asset_id: 2, location: 'Kamar 102', priority: 'Medium', description: 'Smart TV HDMI port tidak terdeteksi', status: 'Closed', requester_id: 2, engineer_id: 103 },
  { id: 4, asset_id: 5, location: 'Lobby Lounge', priority: 'Low', description: 'Kaki sofa kendor perlu dikencangkan', status: 'Open', requester_id: 5, engineer_id: 0 }
])

const showCreateModal = ref(false)
const formWo = ref({ asset_id: 1, location: '', priority: 'Medium', description: '' })

const showAssignModal = ref(false)
const selectedWo = ref(null)
const assignEngineerId = ref(101)

const showUpdateModal = ref(false)
const updateStatus = ref('In Progress')
const updateActionTaken = ref('')
const updateCost = ref(0)

const filteredWorkOrders = computed(() => {
  if (activeTab.value === 'all') return workOrders.value
  if (activeTab.value === 'open') return workOrders.value.filter(w => w.status === 'Open')
  if (activeTab.value === 'in_progress') return workOrders.value.filter(w => w.status === 'In Progress')
  if (activeTab.value === 'under_review') return workOrders.value.filter(w => w.status === 'Under Review')
  if (activeTab.value === 'completed') return workOrders.value.filter(w => w.status === 'Closed' || w.status === 'Completed')
  return workOrders.value
})

function getTabCount(tabId) {
  if (tabId === 'all') return workOrders.value.length
  if (tabId === 'open') return workOrders.value.filter(w => w.status === 'Open').length
  if (tabId === 'in_progress') return workOrders.value.filter(w => w.status === 'In Progress').length
  if (tabId === 'under_review') return workOrders.value.filter(w => w.status === 'Under Review').length
  if (tabId === 'completed') return workOrders.value.filter(w => w.status === 'Closed' || w.status === 'Completed').length
  return 0
}

async function fetchWorkOrders() {
  try {
    const res = await api.get('/workorders')
    if (res.data?.data && Array.isArray(res.data.data) && res.data.data.length > 0) {
      workOrders.value = res.data.data
    }
  } catch (e) {
    // Keep initial mock list
  }
}

function submitWorkOrder() {
  workOrders.value.unshift({
    id: Date.now(),
    asset_id: formWo.value.asset_id,
    location: formWo.value.location || 'Area Hotel',
    priority: formWo.value.priority,
    description: formWo.value.description,
    status: 'Open',
    requester_id: 1,
    engineer_id: 0
  })
  showCreateModal.value = false
  alert('Tiket Work Order / Laporan Kerusakan berhasil diajukan!')
}

function openAssignModal(wo) {
  selectedWo.value = wo
  assignEngineerId.value = 101
  showAssignModal.value = true
}

function submitAssign() {
  if (!selectedWo.value) return
  selectedWo.value.engineer_id = assignEngineerId.value
  selectedWo.value.status = 'In Progress'
  showAssignModal.value = false
  alert(`Teknisi #${assignEngineerId.value} berhasil ditugaskan ke WO #${selectedWo.value.id}!`)
}

function openUpdateModal(wo) {
  selectedWo.value = wo
  updateStatus.value = wo.status === 'Open' ? 'In Progress' : wo.status
  updateActionTaken.value = ''
  updateCost.value = 0
  showUpdateModal.value = true
}

function submitUpdateStatus() {
  if (!selectedWo.value) return
  selectedWo.value.status = updateStatus.value
  showUpdateModal.value = false
  alert(`Status Work Order #${selectedWo.value.id} diperbarui menjadi "${updateStatus.value}"!`)
}

function closeOrder(wo) {
  wo.status = 'Closed'
  alert(`Work Order #${wo.id} telah ditutup (Closed) oleh Supervisor/HOD!`)
}

function cancelOrder(wo) {
  wo.status = 'Cancelled'
  alert(`Work Order #${wo.id} dibatalkan (Cancelled)!`)
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
  background: #dc2626;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(220, 38, 38, 0.3);
}

.status-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 4px;
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

.wo-id {
  font-family: monospace;
  font-weight: 800;
  color: #0f172a;
}

.desc-text {
  font-weight: 500;
}

.engineer-chip {
  background: #dcfce7;
  color: #15803d;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
}

.unassigned-chip {
  color: #94a3b8;
  font-size: 0.8rem;
  font-style: italic;
}

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
</style>
