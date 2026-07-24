<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Perawatan Rutin</p>
        <h1>📅 Maintenance (PM)</h1>
        <p class="subtitle">Jadwal perawatan berkala.</p>
      </div>
      <button class="primary-btn" v-if="canManageSchedule" @click="openAddModal">
        ➕ Tambah Jadwal
      </button>
    </div>

    <!-- Dedicated Reminder Banner for Staff Engineer -->
    <div v-if="userRole === 'engineer'" class="engineer-reminder-banner">
      <div class="erb-icon">🔔</div>
      <div class="erb-content">
        <h3>Pengingat Perawatan Rutin (Staff Engineer)</h3>
        <p>Halaman ini berisi daftar pengingat <strong>Preventive Maintenance</strong> yang telah dijadwalkan oleh HOD Engineer, Supervisor, dan Administrator. Silakan laksanakan inspeksi fisik aset dan klik tombol <strong>"✅ Selesaikan Checklist"</strong> setelah pengerjaan selesai.</p>
      </div>
    </div>

    <!-- View Switcher & Calendar Toolbar -->
    <div class="pm-toolbar">
      <div class="view-switcher-btns">
        <button :class="['tab-btn', { active: viewMode === 'calendar' }]" @click="viewMode = 'calendar'">
          📅 Kalender Reminder
        </button>
        <button :class="['tab-btn', { active: viewMode === 'grid' }]" @click="viewMode = 'grid'">
          📇 Kartu Detail
        </button>
      </div>

      <div class="calendar-nav-controls" v-if="viewMode === 'calendar'">
        <button class="cal-nav-btn" @click="prevMonth">◀</button>
        <strong class="cal-month-label">{{ calendarMonthTitle }}</strong>
        <button class="cal-nav-btn" @click="nextMonth">▶</button>
        <button class="cal-today-btn" @click="resetToToday">Hari Ini</button>
      </div>
    </div>

    <!-- CALENDAR VIEW -->
    <div v-if="viewMode === 'calendar'" class="calendar-card-panel">
      <div class="calendar-header-weekdays">
        <div class="weekday"><span class="full-day">Minggu</span><span class="short-day">Min</span></div>
        <div class="weekday"><span class="full-day">Senin</span><span class="short-day">Sen</span></div>
        <div class="weekday"><span class="full-day">Selasa</span><span class="short-day">Sel</span></div>
        <div class="weekday"><span class="full-day">Rabu</span><span class="short-day">Rab</span></div>
        <div class="weekday"><span class="full-day">Kamis</span><span class="short-day">Kam</span></div>
        <div class="weekday"><span class="full-day">Jumat</span><span class="short-day">Jum</span></div>
        <div class="weekday"><span class="full-day">Sabtu</span><span class="short-day">Sab</span></div>
      </div>

      <div class="calendar-grid-cells">
        <div 
          v-for="(cell, index) in calendarDays" 
          :key="index" 
          :class="['cal-day-cell', { 'empty-cell': !cell.isCurrentMonth, 'today-cell': cell.isToday }]"
        >
          <div class="cal-day-number" v-if="cell.isCurrentMonth">
            <span>{{ cell.dayNumber }}</span>
            <span v-if="cell.isToday" class="today-tag">Hari Ini</span>
          </div>

          <div class="cal-events-list" v-if="cell.isCurrentMonth">
            <div 
              v-for="item in cell.items" 
              :key="item.id" 
              :class="['cal-event-chip', item.status === 'Completed' ? 'status-completed' : 'status-pending']"
              @click="openDetailModal(item)"
              :title="`Aset #${item.asset_id} (${item.schedule_type}) - ${item.checklist_data}`"
            >
              <span class="event-icon">{{ item.status === 'Completed' ? '✅' : '📌' }}</span>
              <span class="event-title">#{{ item.asset_id }} — {{ getAssetName(item.asset_id) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- GRID CARD VIEW -->
    <div v-else class="card-panel">
      <div class="pm-grid">
        <div class="pm-card" v-for="item in pmList" :key="item.id">
          <div class="pm-card-header">
            <div>
              <span class="schedule-type-badge">{{ item.schedule_type }}</span>
              <h3 class="pm-asset-title">Aset #{{ item.asset_id }} — {{ getAssetName(item.asset_id) }}</h3>
            </div>
            <StatusBadge :status="item.status || 'Active'" />
          </div>

          <div class="pm-details">
            <p><strong>Jatuh Tempo Berikutnya:</strong> ⏰ {{ formatDate(item.next_run) }}</p>
            <div class="checklist-box">
              <p><strong>Checklist Inspeksi:</strong></p>
              <pre>{{ item.checklist_data }}</pre>
            </div>
          </div>

          <div class="pm-actions">
            <button v-if="canCompleteChecklist" class="checklist-btn" @click="submitChecklist(item)">
              ✅ Selesaikan Checklist
            </button>
            <p v-else class="checklist-notice-readonly">ℹ️ Penyelesaian checklist hanya oleh HOD, Supervisor & Admin</p>
            <div v-if="canManageSchedule" class="pm-admin-btns">
              <button class="icon-btn edit-btn" @click="openEditModal(item)" title="Edit Jadwal PM">✏️ Edit</button>
              <button class="icon-btn delete-btn" @click="deletePMSchedule(item)" title="Hapus Jadwal PM">🗑️ Hapus</button>
            </div>
          </div>
        </div>

        <div v-if="pmList.length === 0" class="empty-state">
          Belum ada jadwal Preventive Maintenance di database.
        </div>
      </div>
    </div>

    <!-- Modal Tambah / Edit Schedule PM -->
    <ModalDialog :show="showAddModal" :title="isEditMode ? '✏️ Edit Jadwal Preventive Maintenance' : '➕ Tambah Jadwal Preventive Maintenance'" @close="showAddModal = false">
      <form @submit.prevent="savePMSchedule" class="pm-form">
        <label>
          <span>Pilih Kode Aset Terdaftar</span>
          <select v-model.number="newAssetId" :disabled="isEditMode" required>
            <option value="" disabled>-- Pilih Kode Aset Terdaftar --</option>
            <option v-for="asset in registeredAssets" :key="asset.id" :value="asset.id">
              {{ asset.asset_code }} — {{ asset.asset_name }} (📍 {{ asset.location }})
            </option>
          </select>
        </label>

        <label>
          <span>Tipe Jadwal</span>
          <select v-model="newScheduleType">
            <option value="Daily">Daily (Harian)</option>
            <option value="Weekly">Weekly (Mingguan)</option>
            <option value="Monthly">Monthly (Bulanan)</option>
            <option value="Quarterly">Quarterly (Per 3 Bulan)</option>
            <option value="Yearly">Yearly (Tahunan)</option>
          </select>
        </label>

        <label>
          <span>Jatuh Tempo Berikutnya</span>
          <input v-model="newNextRun" type="date" required />
        </label>

        <label>
          <span>Checklist Item (Inspeksi Per Baris)</span>
          <textarea v-model="newChecklistData" rows="4" placeholder="Contoh: 1. Cek tekanan freon AC&#10;2. Bersihkan filter evaporator&#10;3. Cek drainase air kondensasi" required></textarea>
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Jadwal PM</button>
      </form>
    </ModalDialog>

    <!-- Modal Detail & Eksekusi Checklist PM -->
    <ModalDialog :show="showDetailModal" title="📌 Rincian Reminder Maintenance" @close="showDetailModal = false">
      <div v-if="selectedPmItem" class="pm-detail-modal-body">
        <div class="pm-detail-header-box">
          <span class="schedule-type-badge">{{ selectedPmItem.schedule_type }}</span>
          <h3>Aset #{{ selectedPmItem.asset_id }} — {{ getAssetName(selectedPmItem.asset_id) }}</h3>
          <StatusBadge :status="selectedPmItem.status || 'Active'" />
        </div>

        <div class="pm-detail-info">
          <p><strong>⏰ Jatuh Tempo Perawatan:</strong> {{ formatDate(selectedPmItem.next_run) }}</p>
          <div class="checklist-box">
            <p><strong>📝 Catatan & Checklist Inspeksi:</strong></p>
            <pre>{{ selectedPmItem.checklist_data }}</pre>
          </div>
        </div>

        <div class="pm-detail-actions">
          <button v-if="canCompleteChecklist" class="checklist-btn full-width-btn" @click="submitChecklistFromModal(selectedPmItem)">
            ✅ Selesaikan Checklist Perawatan
          </button>
          <p v-else class="checklist-notice-readonly">ℹ️ Penyelesaian checklist maintenance hanya dapat dikonfirmasi oleh HOD, Supervisor, atau Admin.</p>
          <div v-if="canManageSchedule" class="pm-admin-modal-btns">
            <button class="icon-btn edit-btn" @click="openEditModalFromDetail(selectedPmItem)">✏️ Edit Jadwal</button>
            <button class="icon-btn delete-btn" @click="deleteFromDetail(selectedPmItem)">🗑️ Hapus Jadwal</button>
          </div>
        </div>
      </div>
    </ModalDialog>

    <!-- Custom UI Toast Notification -->
    <transition name="fade">
      <div v-if="showToast" :class="['custom-ui-toast', toastType]">
        <span class="toast-icon">{{ toastType === 'success' ? '✅' : '⚠️' }}</span>
        <span class="toast-text">{{ toastMsg }}</span>
        <button class="toast-close" @click="showToast = false">✕</button>
      </div>
    </transition>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')
const canManageSchedule = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')
const canCompleteChecklist = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')

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

const todayDateStr = new Date().toISOString().split('T')[0]

const viewMode = ref('calendar')
const calendarDate = ref(new Date())

const showDetailModal = ref(false)
const selectedPmItem = ref(null)

function getAssetName(assetId) {
  const found = registeredAssets.value.find(a => a.id === assetId)
  return found ? found.asset_name : `Aset #${assetId}`
}

const calendarMonthTitle = computed(() => {
  return calendarDate.value.toLocaleDateString('id-ID', { month: 'long', year: 'numeric' })
})

function prevMonth() {
  const d = new Date(calendarDate.value)
  d.setMonth(d.getMonth() - 1)
  calendarDate.value = d
}

function nextMonth() {
  const d = new Date(calendarDate.value)
  d.setMonth(d.getMonth() + 1)
  calendarDate.value = d
}

function resetToToday() {
  calendarDate.value = new Date()
}

function isPmItemDueOnDate(item, dateStr) {
  if (!item || !item.next_run || !dateStr) return false
  
  const cellDate = new Date(dateStr + 'T00:00:00')
  if (isNaN(cellDate.getTime())) return false

  const targetDate = new Date(item.next_run)
  if (isNaN(targetDate.getTime())) return false

  const targetDayStart = new Date(targetDate.getFullYear(), targetDate.getMonth(), targetDate.getDate())
  const cellDayStart = new Date(cellDate.getFullYear(), cellDate.getMonth(), cellDate.getDate())

  if (cellDayStart > targetDayStart) return false

  let startDayStart = new Date(targetDayStart)
  startDayStart.setFullYear(startDayStart.getFullYear() - 1)
  if (item.created_at) {
    const parsedStart = new Date(item.created_at)
    if (!isNaN(parsedStart.getTime())) {
      startDayStart = new Date(parsedStart.getFullYear(), parsedStart.getMonth(), parsedStart.getDate())
    }
  }

  if (cellDayStart < startDayStart) return false

  const st = item.schedule_type || 'Monthly'

  if (st === 'Daily') {
    return true
  }

  if (st === 'Weekly') {
    return cellDayStart.getDay() === targetDayStart.getDay()
  }

  if (st === 'Monthly') {
    const targetDom = targetDayStart.getDate()
    const lastDayInMonth = new Date(cellDayStart.getFullYear(), cellDayStart.getMonth() + 1, 0).getDate()
    const expectedDom = Math.min(targetDom, lastDayInMonth)
    return cellDayStart.getDate() === expectedDom
  }

  if (st === 'Quarterly') {
    const targetDom = targetDayStart.getDate()
    const lastDayInMonth = new Date(cellDayStart.getFullYear(), cellDayStart.getMonth() + 1, 0).getDate()
    const expectedDom = Math.min(targetDom, lastDayInMonth)
    const monthDiff = (cellDayStart.getMonth() - targetDayStart.getMonth() + 1200) % 3
    return cellDayStart.getDate() === expectedDom && monthDiff === 0
  }

  if (st === 'Yearly') {
    return cellDayStart.getDate() === targetDayStart.getDate() && cellDayStart.getMonth() === targetDayStart.getMonth()
  }

  return cellDayStart.getTime() === targetDayStart.getTime()
}

const calendarDays = computed(() => {
  const year = calendarDate.value.getFullYear()
  const month = calendarDate.value.getMonth()

  const firstDayIndex = new Date(year, month, 1).getDay()
  const totalDays = new Date(year, month + 1, 0).getDate()

  const days = []

  for (let i = 0; i < firstDayIndex; i++) {
    days.push({ dayNumber: '', dateStr: '', isCurrentMonth: false, items: [] })
  }

  for (let day = 1; day <= totalDays; day++) {
    const monthStr = String(month + 1).padStart(2, '0')
    const dayStr = String(day).padStart(2, '0')
    const dateStr = `${year}-${monthStr}-${dayStr}`

    const items = pmList.value.filter(item => {
      return isPmItemDueOnDate(item, dateStr)
    })

    const isToday = dateStr === todayDateStr

    days.push({
      dayNumber: day,
      dateStr: dateStr,
      isCurrentMonth: true,
      isToday: isToday,
      items: items
    })
  }

  return days
})

function openDetailModal(item) {
  selectedPmItem.value = item
  showDetailModal.value = true
}

async function submitChecklistFromModal(item) {
  showDetailModal.value = false
  await submitChecklist(item)
}

function openEditModalFromDetail(item) {
  showDetailModal.value = false
  openEditModal(item)
}

async function deleteFromDetail(item) {
  showDetailModal.value = false
  await deletePMSchedule(item)
}

const registeredAssets = ref([])
const showAddModal = ref(false)
const isEditMode = ref(false)
const editingPmId = ref(0)
const newAssetId = ref('')
const newScheduleType = ref('Monthly')
const newNextRun = ref(todayDateStr)
const newChecklistData = ref('')

const pmList = ref([])
const isLoading = ref(false)

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

function formatDate(dateStr) {
  if (!dateStr) return '—'
  if (typeof dateStr === 'string' && dateStr.includes('T')) {
    return dateStr.split('T')[0]
  }
  return dateStr
}

async function fetchPMSchedules() {
  isLoading.value = true
  try {
    const res = await api.get('/maintenance/schedules')
    if (res.data?.data && Array.isArray(res.data.data)) {
      pmList.value = res.data.data
    } else {
      pmList.value = []
    }
  } catch (e) {
    console.error('Failed to fetch PM schedules:', e)
  } finally {
    isLoading.value = false
  }
}

function openAddModal() {
  isEditMode.value = false
  editingPmId.value = 0
  newAssetId.value = registeredAssets.value[0]?.id || ''
  newScheduleType.value = 'Monthly'
  newNextRun.value = new Date().toISOString().split('T')[0]
  newChecklistData.value = ''
  showAddModal.value = true
}

function openEditModal(item) {
  isEditMode.value = true
  editingPmId.value = item.id
  newAssetId.value = item.asset_id
  newScheduleType.value = item.schedule_type
  newNextRun.value = formatDate(item.next_run)
  newChecklistData.value = item.checklist_data
  showAddModal.value = true
}

async function savePMSchedule() {
  try {
    if (isEditMode.value) {
      await api.post('/maintenance/edit', {
        pm_id: editingPmId.value,
        schedule_type: newScheduleType.value,
        next_run: newNextRun.value,
        checklist_data: newChecklistData.value
      })
      notify('Jadwal Maintenance berhasil diperbarui!', 'success')
    } else {
      await api.post('/maintenance/schedule', {
        asset_id: newAssetId.value,
        schedule_type: newScheduleType.value,
        next_run: newNextRun.value,
        checklist_data: newChecklistData.value
      })
      notify('Jadwal Maintenance berhasil dibuat!', 'success')
    }
    await fetchPMSchedules()
  } catch (e) {
    notify('Gagal menyimpan Jadwal PM: ' + (e.response?.data?.message || e.message), 'error')
  } finally {
    showAddModal.value = false
  }
}

async function deletePMSchedule(item) {
  try {
    await api.post('/maintenance/delete', { pm_id: item.id })
    notify(`Jadwal PM Aset #${item.asset_id} berhasil dihapus dari database!`, 'success')
    await fetchPMSchedules()
  } catch (e) {
    notify('Gagal menghapus Jadwal PM: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function submitChecklist(item) {
  try {
    await api.post(`/maintenance/${item.id}/checklist`, {
      checklist: item.checklist_data
    })
    item.status = 'Completed'
    notify(`Checklist perawatan Aset #${item.asset_id} berhasil diselesaikan!`, 'success')
    await fetchPMSchedules()
  } catch (e) {
    notify('Gagal menyelesaikan checklist: ' + (e.response?.data?.message || e.message), 'error')
  }
}

onMounted(() => {
  fetchPMSchedules()
  fetchRegisteredAssets()
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

.pm-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;
}

.pm-card {
  background: #ffffff;
  border-radius: 18px;
  border: 1px solid #e2e8f0;
  padding: 20px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.pm-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 14px;
}

.schedule-type-badge {
  font-size: 0.75rem;
  background: #e0f2fe;
  color: #0369a1;
  padding: 2px 8px;
  border-radius: 6px;
  font-weight: 700;
  text-transform: uppercase;
}

.pm-asset-title {
  margin: 4px 0 0;
  font-size: 1.2rem;
  color: #0f172a;
}

.pm-details {
  margin-bottom: 16px;
  font-size: 0.9rem;
}

.checklist-box {
  background: #f8fafc;
  padding: 10px 12px;
  border-radius: 10px;
  margin-top: 10px;
  border: 1px solid #e2e8f0;
}

.checklist-box pre {
  margin: 4px 0 0;
  font-family: inherit;
  font-size: 0.85rem;
  white-space: pre-wrap;
  color: #334155;
}

.pm-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.checklist-btn {
  background: #16a34a;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  width: 100%;
}

.pm-admin-btns {
  display: flex;
  gap: 8px;
}

.icon-btn {
  border: 1px solid #cbd5e1;
  background: white;
  padding: 6px 12px;
  border-radius: 8px;
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
  flex: 1;
}

.edit-btn { color: #2563eb; border-color: #93c5fd; }
.delete-btn { color: #dc2626; border-color: #fca5a5; }

.pm-form {
  display: grid;
  gap: 14px;
}

.pm-form label {
  display: grid;
  gap: 6px;
  font-weight: 600;
  font-size: 0.9rem;
}

.pm-form input, .pm-form select, .pm-form textarea {
  padding: 10px 12px;
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

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-size: 0.95rem;
}

.engineer-reminder-banner {
  background: #fef3c7;
  border: 1px solid #fcd34d;
  border-left: 5px solid #d97706;
  border-radius: 2px !important;
  padding: 16px 20px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.erb-icon {
  font-size: 2rem;
}

.erb-content h3 {
  margin: 0 0 4px 0;
  font-size: 1.05rem;
  color: #92400e;
  font-weight: 800;
}

.erb-content p {
  margin: 0;
  font-size: 0.9rem;
  color: #78350f;
  line-height: 1.4;
}

.pm-toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  flex-wrap: wrap;
  gap: 12px;
}

.view-switcher-btns {
  display: flex;
  gap: 8px;
}

.tab-btn {
  background: #f1f5f9;
  border: 1px solid #cbd5e1;
  color: #475569;
  padding: 8px 16px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.88rem;
  cursor: pointer;
}

.tab-btn.active {
  background: #0f172a;
  color: white;
  border-color: #0f172a;
}

.calendar-nav-controls {
  display: flex;
  align-items: center;
  gap: 10px;
}

.cal-nav-btn {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  color: #0f172a;
  width: 34px;
  height: 34px;
  border-radius: 2px !important;
  font-weight: 800;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cal-month-label {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
  min-width: 140px;
  text-align: center;
}

.cal-today-btn {
  background: #dbeafe;
  color: #1e40af;
  border: 1px solid #bfdbfe;
  padding: 6px 12px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.8rem;
  cursor: pointer;
}

.calendar-card-panel {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.calendar-header-weekdays {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  background: #0f172a;
  color: white;
  text-align: center;
  font-weight: 700;
  font-size: 0.85rem;
}

.weekday {
  padding: 10px 4px;
  border-right: 1px solid #1e293b;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.weekday:last-child {
  border-right: none;
}

.calendar-grid-cells {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  background: #cbd5e1;
  gap: 1px;
}

.cal-day-cell {
  background: #ffffff;
  min-height: 115px;
  padding: 6px;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  min-width: 0;
  box-sizing: border-box;
  overflow: hidden;
}

.cal-day-cell.empty-cell {
  background: #f8fafc;
}

.cal-day-cell.today-cell {
  background: #fefce8;
}

.cal-day-number {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 800;
  font-size: 0.85rem;
  color: #334155;
  margin-bottom: 4px;
  width: 100%;
}

.today-tag {
  background: #eab308;
  color: #0f172a;
  font-size: 0.65rem;
  padding: 1px 4px;
  border-radius: 4px;
  font-weight: 800;
  white-space: nowrap;
}

.cal-events-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
  max-height: 85px;
  width: 100%;
  box-sizing: border-box;
}

.cal-event-chip {
  padding: 4px 6px;
  border-radius: 4px;
  font-size: 0.72rem;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 4px;
  min-width: 0;
  width: 100%;
  box-sizing: border-box;
  overflow: hidden;
  transition: transform 0.1s ease;
}

.cal-event-chip:hover {
  transform: scale(1.02);
}

.cal-event-chip.status-pending {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fcd34d;
}

.cal-event-chip.status-completed {
  background: #dcfce7;
  color: #15803d;
  border: 1px solid #86efac;
}

.event-icon {
  flex-shrink: 0;
}

.event-title {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
  flex: 1;
}

.pm-detail-modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.pm-detail-header-box {
  display: flex;
  flex-direction: column;
  gap: 6px;
  background: #f8fafc;
  padding: 14px;
  border-radius: 2px !important;
  border: 1px solid #e2e8f0;
}

.pm-detail-header-box h3 {
  margin: 0;
  font-size: 1.1rem;
  color: #0f172a;
}

.pm-detail-info {
  display: flex;
  flex-direction: column;
  gap: 12px;
  font-size: 0.9rem;
}

.full-width-btn {
  width: 100%;
  padding: 12px;
  font-size: 0.95rem;
}

.pm-detail-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.pm-admin-modal-btns {
  display: flex;
  gap: 10px;
}

.checklist-notice-readonly {
  font-size: 0.8rem;
  color: #78350f;
  background: #fef3c7;
  border: 1px solid #fcd34d;
  padding: 8px 12px;
  border-radius: 2px !important;
  font-weight: 600;
  margin: 0;
  text-align: center;
}

.short-day {
  display: none;
}

@media (max-width: 640px) {
  .full-day {
    display: none;
  }
  .short-day {
    display: inline;
  }
  .cal-day-cell {
    min-height: 70px !important;
    padding: 3px !important;
  }
  .cal-day-number {
    font-size: 0.75rem !important;
  }
  .today-tag {
    font-size: 0.55rem !important;
    padding: 0 2px !important;
  }
  .cal-event-chip {
    font-size: 0.65rem !important;
    padding: 2px 4px !important;
  }
  .pm-toolbar {
    flex-direction: column;
    align-items: stretch;
  }
  .view-switcher-btns {
    width: 100%;
  }
  .tab-btn {
    flex: 1;
    text-align: center;
  }
  .cal-month-label {
    min-width: auto;
    font-size: 0.95rem;
  }
}
</style>
