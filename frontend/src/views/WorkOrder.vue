<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Work Order System</p>
        <h1>Work Order & Perbaikan</h1>
        <p class="subtitle">Pelaporan & perbaikan kerusakan.</p>
      </div>

      <div class="header-action-group">
        <button class="primary-btn btn-secondary-ios" @click="showReportModal = true" title="Prinjau & Export Laporan Bulanan">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
          <span>Laporan & Export</span>
        </button>
        <button class="primary-btn" @click="showCreateModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          <span>Buat Work Order Baru</span>
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
                <th>ID Work Order</th>
                <th>Kategori Aset</th>
                <th>Lokasi / Area</th>
                <th>Kode / ID Aset</th>
                <th>Pelapor (Username)</th>
                <th>Departemen Asal</th>
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
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-indigo">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2H2v10l9.29 9.29c.94.94 2.48.94 3.42 0l6.58-6.58c.94-.94.94-2.48 0-3.42L12 2Z"/><circle cx="7" cy="7" r=".5" fill="currentColor"/></svg>
                    <span>{{ wo.category || 'HVAC / AC' }}</span>
                  </div>
                </td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-blue">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>
                    <span>{{ wo.location || 'Ruangan / Area Operasional' }}</span>
                  </div>
                </td>
                <td class="nowrap-cell">Aset #{{ wo.asset_id }}</td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-slate">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                    <span>@{{ wo.requested_by || 'user_system' }}</span>
                  </div>
                </td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-orange">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><path d="M9 22v-4h6v4"/><path d="M8 6h.01"/><path d="M16 6h.01"/><path d="M12 6h.01"/><path d="M12 10h.01"/><path d="M12 14h.01"/><path d="M16 10h.01"/><path d="M16 14h.01"/><path d="M8 10h.01"/><path d="M8 14h.01"/></svg>
                    <span>{{ formatDepartmentLabel(wo.department) }}</span>
                  </div>
                </td>
                <td class="nowrap-cell"><StatusBadge :status="wo.priority || 'Medium'" /></td>
                <td class="desc-cell" :title="wo.description"><span class="desc-text">{{ wo.description }}</span></td>
                <td class="nowrap-cell">
                  <div v-if="wo.engineer_id > 0" class="ios-table-pill pill-emerald">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
                    <span>{{ getEngineerName(wo.engineer_id) }}</span>
                  </div>
                  <span v-else class="unassigned-chip">Belum Ditugaskan</span>
                </td>
                <td class="nowrap-cell"><StatusBadge :status="wo.status" /></td>
                <td class="actions-cell">
                  <button class="icon-btn log-btn" @click.stop="openLogsModal(wo)" title="Lihat Laporan Timeline Progres">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg> Timeline
                  </button>
                  <button v-if="canAssign && wo.status === 'Open'" class="icon-btn assign-btn" @click="openAssignModal(wo)" title="Assign Worker">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><polyline points="16 11 18 13 22 9"/></svg> Assign
                  </button>
                  <button v-if="canUpdateProgress && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn progress-btn" @click="openUpdateModal(wo)" title="Update Progres">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg> Update
                  </button>
                  <button v-if="canManageOrder && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn close-btn" @click="closeOrder(wo)" title="Selesaikan Work Order">
                    ✅ Selesai
                  </button>

                  <!-- HOD, Admin, and Management can cancel open work orders -->
                  <button v-if="canCancel && wo.status === 'Open'" class="icon-btn cancel-btn" @click="cancelOrder(wo)" title="Batal (Cancel Work Order)">
                    🚫 Batal
                  </button>

                <!-- HOD, Admin, and Management can delete work orders -->
                <button v-if="canDelete" class="icon-btn delete-btn" @click="deleteOrder(wo)" title="Hapus Work Order Permanen">
                  🗑️ Hapus
                </button>
              </td>
            </tr>
            <tr v-if="filteredWorkOrders.length === 0">
              <td colspan="11" class="empty-state">Tidak ada Work Order pada kategori ini.</td>
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
    <ModalDialog :show="showCreateModal" title="Ajukan Laporan Kerusakan Aset" @close="showCreateModal = false">
      <form @submit.prevent="submitWorkOrder" class="modal-form">
        <label>
          <span>Pilih Kode Aset Terdaftar</span>
          <select v-model.number="formWo.asset_id" @change="onAssetSelected" required>
            <option value="" disabled>-- Pilih Kode Aset Terdaftar --</option>
            <option v-for="asset in registeredAssets" :key="asset.id" :value="asset.id">
              {{ asset.asset_code }} — {{ asset.asset_name }} ({{ asset.location }})
            </option>
          </select>
        </label>
        <label>
          <span>Kategori Kerusakan / Aset</span>
          <input v-model="formWo.category" placeholder="Contoh: HVAC / AC, Elektronik, Plumbing, Dapur" required />
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

    <ModalDialog :show="showAssignModal" title="Penugasan Teknisi (Assign Worker)" @close="showAssignModal = false">
      <form @submit.prevent="submitAssign" class="modal-form" v-if="selectedWo">
        <p class="modal-info"><strong>WO #{{ selectedWo.id }}:</strong> {{ selectedWo.description }} ({{ selectedWo.location }})</p>
        <label>
          <span>Pilih Teknisi (Staff Engineer)</span>
          <select v-model.number="assignEngineerId" required>
            <option value="" disabled>-- Pilih Teknisi (Staff Engineer) --</option>
            <option v-for="eng in engineersList" :key="eng.id" :value="eng.id">
              {{ eng.name }} (@{{ eng.username }})
            </option>
          </select>
        </label>

        <button type="submit" class="submit-modal-btn">Tugaskan Teknisi</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showUpdateModal" title="Update Status Pengerjaan Teknisi" @close="showUpdateModal = false">
      <form @submit.prevent="submitUpdateStatus" class="modal-form" v-if="selectedWo">
        <label>
          <span>Status Baru</span>
          <select v-model="updateStatus">
            <option value="In Progress">In Progress (Sedang Dikerjakan)</option>
            <option value="Under Review">Under Review (Selesai Dikerjakan & Menunggu Review)</option>
            <option value="Finish">Finish / Completed (Selesai & Disetujui)</option>
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
    <ModalDialog :show="showReportModal" title="Laporan Bulanan Work Order & Maintenance" maxWidth="920px" @close="showReportModal = false">
      <div class="monthly-report-printable" id="printableReportDocument">
        <div class="report-header">
          <h2>LAPORAN BULANAN MANAJEMEN ASET & WORK ORDER</h2>
          <p class="report-sub">Sistem AsetKu — Periode: {{ reportMonthYear }}</p>
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
        <button class="excel-btn" @click="exportToExcel">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="8" y1="13" x2="16" y2="13"/><line x1="8" y1="17" x2="16" y2="17"/></svg>
          <span>Export ke Excel (.xlsx)</span>
        </button>
        <button class="print-btn" @click="printMonthlyReport">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
          <span>Cetak Dokumen Laporan (PDF / Print)</span>
        </button>
      </div>
    </ModalDialog>

    <!-- Modal Laporan Progres & Timeline Work Order -->
    <ModalDialog :show="showLogsModal" title="Laporan Timeline Progres Work Order" @close="showLogsModal = false">
      <div v-if="selectedWoForLogs" class="logs-modal-body">
        <div class="wo-info-banner">
          <div>
            <span class="wo-badge">#WO-{{ selectedWoForLogs.id }}</span>
            <h3 class="wo-banner-title">📍 {{ selectedWoForLogs.location }}</h3>
            <p class="wo-banner-sub">Pelapor: <strong>@{{ selectedWoForLogs.requested_by || 'user' }}</strong> (🏢 {{ formatDepartmentLabel(selectedWoForLogs.department) }})</p>
          </div>
          <StatusBadge :status="selectedWoForLogs.status" />
        </div>

        <div class="timeline-container">
          <h4 class="timeline-title">⏱️ Riwayat Perubahan & Progres Pengerjaan</h4>
          
          <div v-if="isLogsLoading" class="logs-loading">Memuat timeline progres...</div>
          
          <div v-else class="timeline-list">
            <div v-for="(log, idx) in woProgressLogs" :key="log.id || idx" class="timeline-item">
              <div class="timeline-node">
                <span class="node-icon">{{ getStatusIcon(log.status) }}</span>
              </div>
              <div class="timeline-content">
                <div class="timeline-header">
                  <StatusBadge :status="log.status" />
                  <span class="timeline-time">🕒 {{ formatDate(log.created_at) }}</span>
                </div>
                <p class="timeline-actor">
                  👤 Oleh: <strong>@{{ log.updated_by || 'Sistem' }}</strong> 
                  <span class="user-role-chip" v-if="log.user_role">({{ formatDepartmentLabel(log.user_role) }})</span>
                </p>
                <div class="timeline-notes" v-if="log.action_taken">
                  <p><strong>📝 Catatan / Tindakan:</strong> {{ log.action_taken }}</p>
                </div>
                <p v-if="log.cost > 0" class="timeline-cost">💰 Biaya Terkait: <strong>Rp {{ formatNumber(log.cost) }}</strong></p>
              </div>
            </div>

            <div v-if="woProgressLogs.length === 0" class="empty-logs">
              Belum ada riwayat progres tercatat untuk Work Order ini.
            </div>
          </div>
        </div>

        <!-- Form Tambah Catatan Timeline Baru -->
        <div v-if="canUpdateProgress && selectedWoForLogs.status !== 'Finish' && selectedWoForLogs.status !== 'Cancelled'" class="add-timeline-box">
          <h4 class="add-tl-title">➕ Tambah Catatan / Update Progres Timeline</h4>
          <form @submit.prevent="submitAddTimelineNote" class="add-tl-form">
            <div class="add-tl-row">
              <label class="tl-field">
                <span>Status Progres:</span>
                <select v-model="newTlStatus" class="modal-input">
                  <option value="In Progress">👷 In Progress (Sedang Dikerjakan)</option>
                  <option value="Under Review">🔍 Under Review (Menunggu Review HOD)</option>
                  <option value="Finish">✅ Finish / Completed (Selesai & Disetujui)</option>
                </select>
              </label>

              <label class="tl-field">
                <span>Biaya Perbaikan (Rp):</span>
                <input v-model.number="newTlCost" type="number" min="0" placeholder="Contoh: 150000" class="modal-input" />
              </label>
            </div>

            <label class="tl-field">
              <span>Catatan Progres / Tindakan Pengerjaan:</span>
              <textarea v-model="newTlActionTaken" rows="2" placeholder="Tuliskan catatan progres terbaru..." class="modal-input modal-textarea" required></textarea>
            </label>

            <button type="submit" class="submit-modal-btn add-tl-btn" :disabled="isSubmittingTl">
              {{ isSubmittingTl ? 'Menyimpan...' : '➕ Simpan Catatan Timeline' }}
            </button>
          </form>
        </div>
      </div>
    </ModalDialog>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const route = useRoute()

const userRole = ref(sessionStorage.getItem('user_role') || 'external')

const canAssign = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canManageOrder = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canUpdateProgress = computed(() => userRole.value === 'engineer' || userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canDelete = computed(() => userRole.value === 'hod' || userRole.value === 'admin' || userRole.value === 'management')
const canCancel = computed(() => userRole.value === 'hod' || userRole.value === 'admin' || userRole.value === 'management')

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
const formWo = ref({ asset_id: '', category: '', location: '', priority: 'Medium', description: '' })

const showAssignModal = ref(false)
const selectedWo = ref(null)
const assignEngineerId = ref('')
const engineersList = ref([])

function getEngineerName(id) {
  const eng = engineersList.value.find(e => e.id === id)
  return eng ? eng.name : `Teknisi #${id}`
}

async function fetchEngineers() {
  try {
    const res = await api.get('/users/engineers')
    if (res.data?.data && Array.isArray(res.data.data)) {
      engineersList.value = res.data.data
      if (engineersList.value.length > 0 && !assignEngineerId.value) {
        assignEngineerId.value = engineersList.value[0].id
      }
    }
  } catch (e) {
    console.error('Failed to fetch engineers list:', e)
  }
}

const showUpdateModal = ref(false)
const updateStatus = ref('In Progress')
const updateActionTaken = ref('')
const updateCost = ref(0)

const showReportModal = ref(false)
const reportMonthYear = ref(new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }))

const showLogsModal = ref(false)
const selectedWoForLogs = ref(null)
const woProgressLogs = ref([])
const isLogsLoading = ref(false)
const newTlStatus = ref('In Progress')
const newTlActionTaken = ref('')
const newTlCost = ref(0)
const isSubmittingTl = ref(false)

function getStatusIcon(status) {
  const map = {
    Open: '🚨',
    'In Progress': '👷',
    'Under Review': '🔍',
    Finish: '✅',
    Completed: '✅',
    Cancelled: '🚫'
  }
  return map[status] || '📌'
}

function openLogsModal(wo) {
  if (!wo) return
  selectedWoForLogs.value = wo
  showLogsModal.value = true

  newTlStatus.value = wo.status === 'Open' ? 'In Progress' : wo.status
  newTlActionTaken.value = ''
  newTlCost.value = 0

  // Build instant milestone timeline from wo object
  const initialLogs = [
    {
      id: 1,
      work_order_id: wo.id,
      status: 'Open',
      action_taken: `Laporan diajukan: ${wo.description || ''}`,
      updated_by: wo.requested_by || 'Staff Operasional',
      user_role: wo.department || 'User',
      created_at: wo.created_at || new Date().toISOString()
    }
  ]

  if (wo.status && wo.status !== 'Open') {
    initialLogs.push({
      id: 2,
      work_order_id: wo.id,
      status: 'In Progress',
      action_taken: `Penugasan Teknisi untuk perbaikan di lokasi ${wo.location || ''}`,
      updated_by: 'HOD Engineer',
      user_role: 'HOD Engineer',
      created_at: wo.created_at || new Date().toISOString()
    })
  }

  if (wo.status === 'Under Review' || wo.status === 'Finish' || wo.status === 'Completed' || wo.status === 'Closed') {
    initialLogs.push({
      id: 3,
      work_order_id: wo.id,
      status: 'Under Review',
      action_taken: wo.action_taken || 'Perbaikan unit selesai dikerjakan. Menunggu review.',
      cost: wo.cost || 0,
      updated_by: 'Budi Santoso (Teknisi)',
      user_role: 'Staff Engineer',
      created_at: wo.created_at || new Date().toISOString()
    })
  }

  if (wo.status === 'Finish' || wo.status === 'Completed' || wo.status === 'Closed' || wo.status === 'Cancelled') {
    initialLogs.push({
      id: 4,
      work_order_id: wo.id,
      status: wo.status,
      action_taken: wo.status === 'Cancelled' ? 'Work Order dibatalkan' : 'Work Order diverifikasi selesai',
      cost: wo.cost || 0,
      updated_by: 'Administrator',
      user_role: 'Admin',
      created_at: wo.closed_at || new Date().toISOString()
    })
  }

  woProgressLogs.value = initialLogs
  isLogsLoading.value = false

  // Background fetch live timeline records from MySQL DB
  api.get(`/workorders/timeline?wo_id=${wo.id}`).then(res => {
    const logsData = res.data?.data || res.data
    if (Array.isArray(logsData) && logsData.length > 0) {
      woProgressLogs.value = logsData
    }
  }).catch(e => {
    console.error('Background fetch WO timeline error:', e)
  })
}

async function submitAddTimelineNote() {
  if (!selectedWoForLogs.value || !newTlActionTaken.value) return
  isSubmittingTl.value = true
  try {
    const payload = {
      work_order_id: selectedWoForLogs.value.id,
      status: newTlStatus.value,
      action_taken: newTlActionTaken.value,
      cost: newTlCost.value || 0
    }
    await api.post('/workorders/timeline/add', payload)
    notify('Catatan timeline berhasil ditambahkan!', 'success')
    newTlActionTaken.value = ''
    newTlCost.value = 0
    selectedWoForLogs.value.status = newTlStatus.value
    if (payload.cost > 0) {
      selectedWoForLogs.value.cost = payload.cost
    }
    
    // Refresh live timeline
    const res = await api.get(`/workorders/timeline?wo_id=${selectedWoForLogs.value.id}`)
    const logsData = res.data?.data || res.data
    if (Array.isArray(logsData) && logsData.length > 0) {
      woProgressLogs.value = logsData
    }
    fetchWorkOrders(true)
  } catch (e) {
    console.error('Failed to add timeline note:', e)
    notify(e.response?.data?.error || 'Gagal menambahkan catatan timeline.', 'error')
  } finally {
    isSubmittingTl.value = false
  }
}

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

const currentUsername = ref(sessionStorage.getItem('username') || sessionStorage.getItem('user_name') || 'admin')
const currentUserRole = ref(sessionStorage.getItem('user_role') || 'admin')

function formatDepartmentLabel(roleOrDept) {
  if (!roleOrDept) return 'Staff Operasional'
  const map = {
    dept_akunting: 'Departement Akunting',
    dept_spa: 'Departement Spa',
    dept_sales: 'Department Sales',
    dept_hr: 'Department Human Resources',
    dept_fb_kitchen: 'Department Food Beverage Kitchen',
    dept_fb_service: 'Department Food Beverage Service',
    dept_housekeeping: 'Department House Keeping',
    dept_frontoffice: 'Department Front Office',
    admin: 'Administrator',
    hod: 'HOD Engineer',
    management: 'Supervisor Engineer',
    engineer: 'Staff Engineer',
    external: 'Staff Operasional'
  }
  return map[roleOrDept] || roleOrDept
}

async function fetchWorkOrders(isSilent = false) {
  if (!isSilent) isLoading.value = true
  try {
    const res = await api.get('/workorders')
    if (res.data?.data && Array.isArray(res.data.data)) {
      workOrders.value = res.data.data
    } else {
      workOrders.value = []
    }
  } catch (e) {
    if (!isSilent) console.error('Failed to fetch work orders from DB:', e)
  } finally {
    if (!isSilent) isLoading.value = false
  }
}

async function submitWorkOrder() {
  try {
    await api.post('/workorders', {
      asset_id: Number(formWo.value.asset_id) || 1,
      category: formWo.value.category || 'HVAC / AC',
      location: formWo.value.location || 'Ruangan / Area Operasional',
      priority: formWo.value.priority,
      description: formWo.value.description,
      requested_by: currentUsername.value,
      department: currentUserRole.value
    })
    showCreateModal.value = false
    formWo.value = { asset_id: '', category: '', location: '', priority: 'Medium', description: '' }
    notify('Tiket Work Order / Laporan Kerusakan berhasil diajukan ke database!', 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal mengajukan Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openAssignModal(wo) {
  selectedWo.value = wo
  if (engineersList.value.length > 0) {
    assignEngineerId.value = engineersList.value[0].id
  } else {
    assignEngineerId.value = ''
  }
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

import { exportToExcel as exportToExcelHelper, triggerPrint } from '../utils/exportUtils'

function printMonthlyReport() {
  triggerPrint()
}

function exportToExcel() {
  const monthName = reportMonthYear.value.replace(/\s+/g, '_')
  const fileName = `Laporan_WorkOrder_${monthName}.xls`
  const headers = ['Kode WO', 'Lokasi / Kamar', 'Prioritas', 'Deskripsi Kerusakan', 'Status Tiket', 'Biaya Perbaikan (Rp)']
  const rows = workOrders.value.map(wo => [
    wo.wo_code || `#WO-${wo.id}`,
    wo.location || '-',
    wo.priority || '-',
    wo.description || '-',
    wo.status || '-',
    wo.cost || 0
  ])
  exportToExcelHelper(fileName, headers, rows)
}

const registeredAssets = ref([])

async function fetchRegisteredAssets() {
  try {
    const res = await api.get('/assets')
    if (res.data?.data && Array.isArray(res.data.data)) {
      registeredAssets.value = res.data.data
    }
  } catch (e) {
    console.error('Failed to fetch registered assets:', e)
  }
}

function onAssetSelected() {
  const selected = registeredAssets.value.find(a => a.id === formWo.value.asset_id)
  if (selected) {
    if (!formWo.value.category) formWo.value.category = selected.category || ''
    if (!formWo.value.location) formWo.value.location = selected.location || ''
  }
}

let pollTimer = null

onMounted(() => {
  if (route.query.assetId) {
    formWo.value.asset_id = parseInt(route.query.assetId)
    formWo.value.location = route.query.location || ''
    showCreateModal.value = true
  }
  fetchWorkOrders()
  fetchRegisteredAssets()
  fetchEngineers()

  // Real-time status sync across all roles & screens
  pollTimer = setInterval(() => {
    fetchWorkOrders(true)
  }, 3000)
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
})
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 16px;
}

.header-action-group {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 37.6px;
}

.primary-btn {
  background: #007aff !important;
  color: #ffffff !important;
  border: 1px solid #007aff !important;
  height: 37.6px !important;
  padding: 0 18px !important;
  border-radius: 10px !important;
  font-size: 0.88rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px !important;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25) !important;
  transition: all 0.15s ease !important;
  line-height: 1 !important;
  white-space: nowrap !important;
  box-sizing: border-box !important;
  margin: 0 !important;
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.primary-btn svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  display: block;
}

.primary-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.primary-btn.btn-secondary-ios {
  background: #ffffff !important;
  color: #0f172a !important;
  border: 1px solid #cbd5e1 !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04) !important;
}

.primary-btn.btn-secondary-ios:hover {
  background: #f1f5f9 !important;
  border-color: #94a3b8 !important;
}

.print-report-btn {
  background: #ffffff;
  color: #0f172a;
  border: 1px solid #cbd5e1;
  padding: 10px 18px;
  border-radius: 4px !important;
  font-weight: 700;
  font-size: 0.88rem;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.print-report-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
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

.ios-table-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 6px !important;
  font-size: 0.82rem;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
}

.pill-indigo { background: #e0e7ff; color: #3730a3; border: 1px solid #c7d2fe; }
.pill-blue { background: #eff6ff; color: #1e40af; border: 1px solid #bfdbfe; }
.pill-slate { background: #f1f5f9; color: #334155; border: 1px solid #cbd5e1; }
.pill-orange { background: #fff7ed; color: #c2410c; border: 1px solid #ffedd5; }
.pill-emerald { background: #ecfdf5; color: #065f46; border: 1px solid #a7f3d0; }

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
  max-width: 220px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.desc-text {
  font-weight: 500;
  display: block;
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
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

.icon-btn svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  display: block;
}

.icon-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.log-btn {
  background: #eff6ff;
  color: #2563eb;
  border-color: #bfdbfe;
}

.log-btn:hover {
  background: #dbeafe;
  border-color: #93c5fd;
}

.assign-btn {
  background: #f0fdf4;
  color: #16a34a;
  border-color: #bbf7d0;
}

.assign-btn:hover {
  background: #dcfce7;
}

.progress-btn {
  background: #fffbeb;
  color: #d97706;
  border-color: #fef3c7;
}

.progress-btn:hover {
  background: #fef3c7;
}

.close-btn {
  background: #ecfdf5;
  color: #059669;
  border-color: #a7f3d0;
}

.close-btn:hover {
  background: #d1fae5;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
  padding: 0 8px;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.modal-form {
  display: grid;
  gap: 16px;
}

.modal-form label {
  display: grid;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: -0.01em;
}

.modal-form input, .modal-form select, .modal-form textarea {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 4px !important;
  font-size: 0.92rem;
  color: #0f172a;
  background: #ffffff;
  outline: none;
  transition: all 0.15s ease;
}

.modal-form input:focus, .modal-form select:focus, .modal-form textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
}

.submit-modal-btn {
  background: #0f172a;
  color: #ffffff;
  border: 1px solid #0f172a;
  padding: 13px 18px;
  border-radius: 4px !important;
  font-size: 0.92rem;
  font-weight: 700;
  cursor: pointer;
  margin-top: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.submit-modal-btn:hover {
  background: #1e293b;
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
  margin-top: 24px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  width: 100%;
}

.excel-btn {
  background: #34c759 !important;
  color: #ffffff !important;
  border: 1px solid #34c759 !important;
  padding: 12px 20px !important;
  border-radius: 10px !important;
  font-size: 0.9rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(52, 199, 89, 0.25);
  transition: all 0.15s ease;
  line-height: 1;
  margin: 0;
}

.excel-btn:hover {
  background: #28a745 !important;
  border-color: #28a745 !important;
  transform: translateY(-1px);
}

.print-btn {
  background: #007aff !important;
  color: #ffffff !important;
  border: 1px solid #007aff !important;
  padding: 12px 20px !important;
  border-radius: 10px !important;
  font-size: 0.9rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
  transition: all 0.15s ease;
  line-height: 1;
  margin: 0;
}

.print-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.requester-chip {
  background: #f1f5f9;
  color: #0f172a;
  padding: 3px 8px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.8rem;
  border: 1px solid #cbd5e1;
}

.dept-chip {
  background: #dbeafe;
  color: #1e40af;
  padding: 3px 8px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.8rem;
  border: 1px solid #bfdbfe;
}

.log-btn {
  background: #f8fafc;
  color: #0284c7;
  border-color: #7dd3fc;
}

.logs-modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.wo-info-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  padding: 14px 16px;
}

.wo-badge {
  font-size: 0.78rem;
  background: #0f172a;
  color: white;
  padding: 2px 8px;
  border-radius: 2px !important;
  font-weight: 800;
}

.wo-banner-title {
  margin: 6px 0 2px;
  font-size: 1.1rem;
  color: #0f172a;
}

.wo-banner-sub {
  margin: 0;
  font-size: 0.85rem;
  color: #64748b;
}

.timeline-container {
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  padding: 16px;
  background: #ffffff;
}

.timeline-title {
  margin: 0 0 16px;
  font-size: 0.95rem;
  color: #0f172a;
  font-weight: 800;
}

.timeline-list {
  display: flex;
  flex-direction: column;
  position: relative;
  padding-left: 20px;
}

.timeline-list::before {
  content: '';
  position: absolute;
  left: 9px;
  top: 8px;
  bottom: 8px;
  width: 2px;
  background: #e2e8f0;
}

.timeline-item {
  position: relative;
  margin-bottom: 20px;
  display: flex;
  gap: 14px;
}

.timeline-item:last-child {
  margin-bottom: 0;
}

.timeline-node {
  position: absolute;
  left: -20px;
  top: 0;
  width: 20px;
  height: 20px;
  background: #ffffff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  z-index: 1;
}

.timeline-content {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  padding: 12px 14px;
  width: 100%;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.timeline-time {
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 600;
}

.timeline-actor {
  margin: 0 0 6px;
  font-size: 0.85rem;
  color: #334155;
}

.user-role-chip {
  color: #0284c7;
  font-weight: 700;
}

.timeline-notes {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 8px 10px;
  border-radius: 2px !important;
  font-size: 0.85rem;
  color: #0f172a;
}

.timeline-notes p {
  margin: 0;
}

.timeline-cost {
  margin: 8px 0 0;
  font-size: 0.85rem;
  color: #16a34a;
}

.logs-loading, .empty-logs {
  text-align: center;
  padding: 24px;
  color: #64748b;
  font-size: 0.9rem;
}

.add-timeline-box {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  padding: 16px;
}

.add-tl-title {
  margin: 0 0 12px;
  font-size: 0.95rem;
  color: #0f172a;
  font-weight: 800;
}

.add-tl-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.add-tl-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.tl-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #334155;
}

.add-tl-btn {
  align-self: flex-end;
  padding: 8px 16px;
  font-weight: 700;
}

/* === Mobile Responsive CSS (Android & iOS) === */
@media (max-width: 640px) {
  .page-container { padding: 16px 14px !important; }
  .page-header { flex-direction: column; align-items: stretch; gap: 12px; }
  .header-action-group { width: 100%; display: flex; flex-wrap: wrap; gap: 8px; }
  .header-action-group .primary-btn { flex: 1; min-width: 130px; justify-content: center; height: 40px !important; font-size: 0.82rem !important; }
  .status-tabs { overflow-x: auto; flex-wrap: nowrap; padding-bottom: 4px; -webkit-overflow-scrolling: touch; }
  .tab-btn { padding: 8px 12px; font-size: 0.8rem; flex-shrink: 0; white-space: nowrap; }
  .card-panel { padding: 16px !important; border-radius: 14px !important; }
  .report-summary-boxes { grid-template-columns: repeat(2, 1fr); gap: 8px; }
  .report-actions { flex-direction: column; gap: 8px; }
  .excel-btn, .print-btn { width: 100%; justify-content: center; }
}
</style>
