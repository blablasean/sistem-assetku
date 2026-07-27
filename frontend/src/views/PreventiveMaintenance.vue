<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Perawatan Rutin</p>
        <h1>Maintenance (PM)</h1>
        <p class="subtitle">Jadwal perawatan berkala.</p>
      </div>
      <div class="header-action-group">
        <button class="primary-btn btn-secondary-ios" @click="showReportModal = true">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
          <span>Laporan & Export</span>
        </button>
        <button class="primary-btn" v-if="canManageSchedule" @click="openAddModal">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 6px;"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg> Tambah Jadwal
        </button>
      </div>
    </div>

    <!-- Dedicated Reminder Banner for Staff Engineer -->
    <div v-if="userRole === 'engineer'" class="engineer-reminder-banner">
      <div class="erb-icon">
        <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M6 8a6 6 0 0 1 12 0c0 7 3 9 3 9H3s3-2 3-9"/><path d="M10.3 21a1.94 1.94 0 0 0 3.4 0"/></svg>
      </div>
      <div class="erb-content">
        <h3>Pengingat Perawatan Rutin (Staff Engineer)</h3>
        <p>Halaman ini berisi daftar pengingat <strong>Preventive Maintenance</strong> yang telah dijadwalkan oleh HOD Engineer, Supervisor, dan Administrator. Silakan laksanakan inspeksi fisik aset dan klik tombol <strong>"Selesaikan Checklist"</strong> setelah pengerjaan selesai.</p>
      </div>
    </div>

    <!-- View Switcher & Calendar Toolbar -->
    <div class="pm-toolbar">
      <div class="view-switcher-btns">
        <button :class="['tab-btn', { active: viewMode === 'calendar' }]" @click="viewMode = 'calendar'">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 6px;"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg> Kalender Reminder
        </button>
        <button :class="['tab-btn', { active: viewMode === 'grid' }]" @click="viewMode = 'grid'">
          <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 6px;"><rect width="7" height="7" x="3" y="3" rx="1"/><rect width="7" height="7" x="14" y="3" rx="1"/><rect width="7" height="7" x="14" y="14" rx="1"/><rect width="7" height="7" x="3" y="14" rx="1"/></svg> Kartu Detail
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
              :key="item.id + '_' + cell.dateStr" 
              :class="['cal-event-chip', isItemCompletedOnDate(item, cell.dateStr) ? 'status-completed' : 'status-pending']"
              @click="openDetailModalForDate(item, cell.dateStr)"
              :title="`Aset #${item.asset_id} (${item.schedule_type}) - ${isItemCompletedOnDate(item, cell.dateStr) ? 'Sudah Di-checklist' : 'Belum Di-checklist'}`"
            >
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
            <p><strong>Jatuh Tempo Berikutnya:</strong> {{ formatDate(item.next_run) }}</p>
            <div class="checklist-box">
              <p><strong>Checklist Inspeksi:</strong></p>
              <pre>{{ item.checklist_data }}</pre>
            </div>
          </div>

          <div class="pm-actions">
            <template v-if="canCompleteChecklist">
              <button 
                v-if="isPmItemDueOnDate(item, getTodayDateStr())"
                :class="['checklist-btn full-width-btn', isItemCompletedOnDate(item, getTodayDateStr()) ? 'btn-completed-style' : 'btn-pending-style']" 
                @click="submitChecklist(item, getTodayDateStr())"
              >
                {{ isItemCompletedOnDate(item, getTodayDateStr()) ? 'Checklist Hari Ini Selesai' : 'Selesaikan Checklist Hari Ini (' + getTodayDateStr() + ')' }}
              </button>
              <p v-else class="checklist-notice-readonly notice-future">
                Hari ini bukan tanggal inspeksi untuk aset ini.
              </p>
            </template>
            <p v-else class="checklist-notice-readonly">Penyelesaian checklist hanya oleh HOD, Supervisor & Admin</p>

            <div v-if="canManageSchedule" class="pm-admin-btns">
              <button class="icon-btn edit-btn" @click="openEditModal(item)" title="Edit Jadwal PM">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                <span>Edit</span>
              </button>
              <button class="icon-btn delete-btn" @click="deletePMSchedule(item)" title="Hapus Jadwal PM">
                <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                <span>Hapus</span>
              </button>
            </div>
          </div>
        </div>

        <div v-if="pmList.length === 0" class="empty-state">
          Belum ada jadwal Preventive Maintenance di database.
        </div>
      </div>
    </div>

    <!-- Modal Tambah / Edit Schedule PM -->
    <ModalDialog :show="showAddModal" :title="isEditMode ? 'Edit Jadwal Preventive Maintenance' : 'Tambah Jadwal Preventive Maintenance'" @close="showAddModal = false">
      <form @submit.prevent="savePMSchedule" class="pm-form">
        <label>
          <span>Pilih Kode Aset Terdaftar</span>
          <select v-model.number="newAssetId" :disabled="isEditMode" required>
            <option value="" disabled>-- Pilih Kode Aset Terdaftar --</option>
            <option v-for="asset in registeredAssets" :key="asset.id" :value="asset.id">
              {{ asset.asset_code }} — {{ asset.asset_name }} ({{ asset.location }})
            </option>
          </select>
        </label>
        <label>
          <span>Frekuensi Rutin (Schedule Type)</span>
          <select v-model="newScheduleType">
            <option value="Daily">Daily (Harian)</option>
            <option value="Weekly">Weekly (Mingguan)</option>
            <option value="Monthly">Monthly (Bulanan)</option>
            <option value="Yearly">Yearly (Tahunan)</option>
          </select>
        </label>
        <label>
          <span>Tanggal Mulai / Jatuh Tempo Pertama</span>
          <input type="date" v-model="newNextRun" required />
        </label>
        <label>
          <span>Rincian Poin Checklist Inspeksi Fisik</span>
          <textarea v-model="newChecklistData" rows="4" placeholder="Contoh: 1. Cek tekanan freon AC&#10;2. Bersihkan filter evaporator&#10;3. Cek drainase air kondensasi" required></textarea>
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Jadwal PM</button>
      </form>
    </ModalDialog>

    <!-- Modal Detail & Eksekusi Checklist PM -->
    <ModalDialog :show="showDetailModal" title="Rincian Reminder Maintenance" @close="showDetailModal = false">
      <div v-if="selectedPmItem" class="pm-detail-modal-body">
        <div class="pm-detail-header-box">
          <div class="pm-header-row">
            <span class="schedule-type-badge">{{ selectedPmItem.schedule_type }}</span>
            <span :class="['status-chip-badge', isItemCompletedOnDate(selectedPmItem, selectedDetailDate) ? 'badge-completed' : 'badge-pending']">
              {{ isItemCompletedOnDate(selectedPmItem, selectedDetailDate) ? 'Sudah Di-checklist' : 'Belum Di-checklist' }}
            </span>
          </div>
          <h3>Aset #{{ selectedPmItem.asset_id }} — {{ getAssetName(selectedPmItem.asset_id) }}</h3>
        </div>

        <div class="pm-detail-info">
          <p><strong>Tanggal Inspeksi Target:</strong> {{ selectedDetailDate }}</p>
          <p><strong>Target Jatuh Tempo Akhir:</strong> {{ formatDate(selectedPmItem.next_run) }}</p>
          <div class="checklist-box">
            <p><strong>Catatan & Checklist Inspeksi:</strong></p>
            <pre>{{ selectedPmItem.checklist_data }}</pre>
          </div>
        </div>

        <div class="pm-detail-actions">
          <template v-if="canCompleteChecklist">
            <button 
              v-if="isDateReached(selectedDetailDate)"
              :class="['checklist-btn full-width-btn', isItemCompletedOnDate(selectedPmItem, selectedDetailDate) ? 'btn-completed-style' : 'btn-pending-style']" 
              @click="submitChecklistFromModal(selectedPmItem, selectedDetailDate)"
            >
              {{ isItemCompletedOnDate(selectedPmItem, selectedDetailDate) ? 'Checklist Sudah Selesai (Klik untuk update)' : 'Selesaikan Checklist Perawatan' }}
            </button>
            <p v-else class="checklist-notice-readonly notice-future">
              Checklist perawatan hanya dapat diisi ketika tanggal jadwal ({{ selectedDetailDate }}) telah tiba.
            </p>
          </template>
          <p v-else class="checklist-notice-readonly">Penyelesaian checklist maintenance hanya dapat dikonfirmasi oleh HOD, Supervisor, atau Admin.</p>

          <div v-if="canManageSchedule" class="pm-admin-modal-btns">
            <button class="icon-btn edit-btn" @click="openEditModalFromDetail(selectedPmItem)">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
              <span>Edit Jadwal</span>
            </button>
            <button class="icon-btn delete-btn" @click="deleteFromDetail(selectedPmItem)">
              <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
              <span>Hapus Jadwal</span>
            </button>
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

  <!-- Laporan Maintenance Modal -->
  <ModalDialog :show="showReportModal" title="Laporan Preventive Maintenance" maxWidth="960px" @close="showReportModal = false">
    <div class="monthly-report-printable" id="printablePMReport">
      <div class="report-header">
        <h2>LAPORAN PREVENTIVE MAINTENANCE</h2>
        <p class="report-sub">Sistem AsetKu &mdash; Periode: {{ reportMonthYear }}</p>
        <hr class="report-divider" />
      </div>

      <div class="report-summary-boxes">
        <div class="rbox">
          <span>Total Jadwal</span>
          <strong>{{ pmList.length }} Jadwal</strong>
        </div>
        <div class="rbox success">
          <span>Daily</span>
          <strong>{{ countPMType('Daily') }} Jadwal</strong>
        </div>
        <div class="rbox warning">
          <span>Weekly / Monthly</span>
          <strong>{{ countPMType('Weekly') + countPMType('Monthly') }} Jadwal</strong>
        </div>
        <div class="rbox danger">
          <span>Yearly</span>
          <strong>{{ countPMType('Yearly') }} Jadwal</strong>
        </div>
      </div>

      <div class="report-table-wrapper">
        <table class="report-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Aset</th>
              <th>Frekuensi</th>
              <th>Jatuh Tempo</th>
              <th>Checklist Inspeksi</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in pmList" :key="item.id">
              <td>#PM-{{ item.id }}</td>
              <td>{{ getAssetName(item.asset_id) }} (Aset #{{ item.asset_id }})</td>
              <td>{{ item.schedule_type }}</td>
              <td>{{ formatDate(item.next_run) }}</td>
              <td style="white-space: pre-wrap; font-size: 0.8rem;">{{ item.checklist_data }}</td>
              <td>{{ item.status || 'Active' }}</td>
            </tr>
            <tr v-if="pmList.length === 0">
              <td colspan="6" class="empty-state">Belum ada jadwal maintenance terdaftar.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div class="report-actions no-print">
      <button class="excel-btn" @click="exportPMToExcel">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="8" y1="13" x2="16" y2="13"/><line x1="8" y1="17" x2="16" y2="17"/></svg>
        <span>Export ke Excel (.xlsx)</span>
      </button>
      <button class="print-btn" @click="() => window.print()">
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
        <span>Cetak Dokumen Laporan (PDF / Print)</span>
      </button>
    </div>
  </ModalDialog>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'
import { useAuth } from '../composables/useAuth'
import { useNotification } from '../composables/useNotification'
import { formatDate, getTodayDateStr } from '../utils/formatters'

const { canManageAssets: canManageSchedule, canManageAssets: canCompleteChecklist, userRole, syncAuth } = useAuth()
const { showToast, toastMsg, toastType, notify } = useNotification()

const todayDateStr = getTodayDateStr()

function isItemCompletedOnDate(item, dateStr) {
  if (!item || !dateStr) return false
  const cleanDateStr = String(dateStr).includes('T') ? String(dateStr).split('T')[0] : String(dateStr).trim()
  if (!cleanDateStr) return false

  if (item.completed_dates && typeof item.completed_dates === 'string') {
    const datesArr = item.completed_dates.split(',').map(d => d.trim().split('T')[0])
    if (datesArr.includes(cleanDateStr)) return true
  }

  const nextRunClean = typeof item.next_run === 'string' ? item.next_run.split('T')[0] : ''
  if (item.status === 'Completed' && (nextRunClean === cleanDateStr)) return true
  return false
}

function isDateReached(dateStr) {
  if (!dateStr) return true
  const cleanDate = String(dateStr).includes('T') ? String(dateStr).split('T')[0] : String(dateStr).trim()
  return cleanDate <= getTodayDateStr()
}

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

const selectedDetailDate = ref('')

function openDetailModalForDate(item, dateStr) {
  selectedPmItem.value = item
  selectedDetailDate.value = dateStr || formatDate(item.next_run)
  showDetailModal.value = true
}

function openDetailModal(item) {
  openDetailModalForDate(item, formatDate(item.next_run))
}

async function submitChecklistFromModal(item, targetDateStr) {
  showDetailModal.value = false
  await submitChecklist(item, targetDateStr)
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

const showReportModal = ref(false)
const reportMonthYear = ref(new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }))

function countPMType(type) {
  return pmList.value.filter(p => p.schedule_type === type).length
}

import { exportToExcel, triggerPrint } from '../utils/exportUtils'

function exportPMToExcel() {
  const fileName = `Laporan_PM_${reportMonthYear.value.replace(/\s+/g, '_')}.xls`
  const headers = ['ID', 'Aset', 'Frekuensi', 'Jatuh Tempo', 'Checklist', 'Status']
  const rows = pmList.value.map(item => [
    `#PM-${item.id}`,
    `${getAssetName(item.asset_id)} (Aset #${item.asset_id})`,
    item.schedule_type || '',
    formatDate(item.next_run),
    (item.checklist_data || '').replace(/\n/g, '; '),
    item.status || 'Active'
  ])
  exportToExcel(fileName, headers, rows)
}

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

async function submitChecklist(item, targetDateStr) {
  try {
    const dateToSubmit = targetDateStr || selectedDetailDate.value || formatDate(item.next_run)
    await api.post(`/maintenance/${item.id}/checklist`, {
      target_date: dateToSubmit,
      checklist: item.checklist_data
    })

    if (!item.completed_dates) {
      item.completed_dates = dateToSubmit
    } else if (!item.completed_dates.includes(dateToSubmit)) {
      item.completed_dates += ',' + dateToSubmit
    }

    notify(`Checklist perawatan Aset #${item.asset_id} tanggal ${dateToSubmit} berhasil diselesaikan (Kuning ➔ Hijau)!`, 'success')
    await fetchPMSchedules()
  } catch (e) {
    notify('Gagal menyelesaikan checklist: ' + (e.response?.data?.message || e.message), 'error')
  }
}

onMounted(() => {
  syncAuth()
  fetchPMSchedules()
  fetchRegisteredAssets()
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
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.pm-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 20px;
}

.pm-card {
  background: #ffffff;
  border-radius: 14px !important;
  border: 1px solid #e2e8f0;
  padding: 20px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  transition: all 0.2s ease;
}

.pm-card:hover {
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
  border-color: #cbd5e1;
}

.pm-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 14px;
}

.schedule-type-badge {
  font-size: 0.72rem;
  background: #e0e7ff;
  color: #3730a3;
  border: 1px solid #c7d2fe;
  padding: 3px 8px;
  border-radius: 6px !important;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  display: inline-block;
  margin-bottom: 4px;
}

.pm-asset-title {
  margin: 4px 0 0;
  font-size: 1.15rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.01em;
}

.pm-details {
  margin-bottom: 16px;
  font-size: 0.9rem;
}

.checklist-box {
  background: #f8fafc;
  padding: 12px 14px;
  border-radius: 8px !important;
  margin-top: 12px;
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
  background: #007aff;
  color: #ffffff;
  border: 1px solid #007aff;
  padding: 11px 16px;
  border-radius: 8px !important;
  font-size: 0.88rem;
  font-weight: 700;
  cursor: pointer;
  width: 100%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.15s ease;
}

.checklist-btn.btn-completed-style {
  background: #ecfdf5;
  color: #059669;
  border-color: #a7f3d0;
}

.checklist-btn.btn-pending-style {
  background: #007aff;
  color: #ffffff;
  border-color: #007aff;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
}

.pm-admin-btns {
  display: flex;
  gap: 8px;
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

.edit-btn {
  background: #f8fafc;
  color: #475569;
  border-color: #cbd5e1;
  padding: 0 10px;
}

.edit-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
  padding: 0 10px;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.pm-form {
  display: grid;
  gap: 14px;
}

.pm-form label {
  display: grid;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: -0.01em;
}

.pm-form input, .pm-form select, .pm-form textarea {
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

.pm-form input:focus, .pm-form select:focus, .pm-form textarea:focus {
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
  width: 100%;
  transition: all 0.15s ease;
}

.submit-modal-btn:hover {
  background: #1e293b;
}

.empty-state {
  grid-column: 1 / -1;
  text-align: center;
  padding: 40px;
  color: #94a3b8;
  font-size: 0.95rem;
}

.engineer-reminder-banner {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-radius: 10px !important;
  padding: 16px 20px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 16px;
}

.erb-icon {
  width: 36px;
  height: 36px;
  background: #dbeafe;
  color: #2563eb;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.erb-content h3 {
  margin: 0 0 4px 0;
  font-size: 0.98rem;
  color: #1e40af;
  font-weight: 800;
}

.erb-content p {
  margin: 0;
  font-size: 0.88rem;
  color: #1e3a8a;
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
  border-radius: 4px !important;
  font-weight: 700;
  font-size: 0.85rem;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.tab-btn.active {
  background: #0f172a;
  color: #ffffff;
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

.pm-header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 8px;
  width: 100%;
}

.status-chip-badge {
  font-size: 0.78rem;
  font-weight: 800;
  padding: 4px 10px;
  border-radius: 4px;
}

.badge-pending {
  background: #fef3c7;
  color: #92400e;
  border: 1px solid #fcd34d;
}

.badge-completed {
  background: #dcfce7;
  color: #15803d;
  border: 1px solid #86efac;
}

.btn-pending-style {
  background: #d97706 !important;
  color: #ffffff !important;
  border: 1px solid #b45309 !important;
}

.btn-completed-style {
  background: #16a34a !important;
  color: #ffffff !important;
  border: 1px solid #15803d !important;
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

/* === Report Modal CSS === */
.header-action-group { display: flex; align-items: center; gap: 10px; }
.monthly-report-printable { padding: 4px 0 16px; }
.report-header { text-align: center; margin-bottom: 20px; }
.report-header h2 { font-size: 1.1rem; font-weight: 800; color: #0f172a; margin: 0 0 4px; }
.report-sub { color: #64748b; font-size: 0.85rem; margin: 0; }
.report-divider { border: none; border-top: 2px solid #e2e8f0; margin: 16px 0 0; }
.report-summary-boxes { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 20px; }
.rbox { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 12px 14px; }
.rbox span { display: block; font-size: 0.75rem; color: #64748b; font-weight: 600; margin-bottom: 4px; }
.rbox strong { font-size: 1.1rem; font-weight: 800; color: #0f172a; }
.rbox.success strong { color: #16a34a; }
.rbox.warning strong { color: #d97706; }
.rbox.danger strong { color: #dc2626; }
.report-table-wrapper { border: 1px solid #e2e8f0; border-radius: 10px; overflow: hidden; }
.report-table { width: 100%; border-collapse: collapse; font-size: 0.82rem; }
.report-table th { background: #0f172a; color: #fff; padding: 9px 12px; text-align: left; font-weight: 700; white-space: nowrap; }
.report-table td { padding: 8px 12px; border-bottom: 1px solid #f1f5f9; color: #334155; vertical-align: top; }
.report-table tbody tr:last-child td { border-bottom: none; }
.report-table tbody tr:hover td { background: #f8fafc; }
.empty-state { text-align: center; color: #94a3b8; font-style: italic; padding: 20px !important; }
.report-actions { display: flex; gap: 12px; padding-top: 16px; justify-content: flex-end; border-top: 1px solid #e2e8f0; margin-top: 16px; }
.excel-btn, .print-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 8px;
  height: 42px; padding: 0 18px; border-radius: 10px; font-size: 0.88rem;
  font-weight: 700; cursor: pointer; border: none; transition: all 0.15s ease;
  line-height: 1; white-space: nowrap;
}
.excel-btn { background: #16a34a; color: #fff; }
.excel-btn:hover { background: #15803d; }
.print-btn { background: #007aff; color: #fff; }
.print-btn:hover { background: #0062cc; }
.btn-icon { flex-shrink: 0; }
.btn-secondary-ios {
  background: #ffffff !important; color: #0f172a !important;
  border: 1px solid #cbd5e1 !important; box-shadow: 0 2px 8px rgba(0,0,0,0.04) !important;
}
.btn-secondary-ios:hover { background: #f1f5f9 !important; border-color: #94a3b8 !important; }
@media print { .no-print { display: none !important; } }
</style>
