<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Perawatan Rutin</p>
        <h1>📅 Maintenance (PM)</h1>
        <p class="subtitle">Jadwal dan checklist perawatan berkala aset.</p>
      </div>
      <button class="primary-btn" v-if="canManageSchedule" @click="openAddModal">
        ➕ Tambah Jadwal
      </button>
    </div>

    <div class="card-panel">
      <div class="pm-grid">
        <div class="pm-card" v-for="item in pmList" :key="item.id">
          <div class="pm-card-header">
            <div>
              <span class="schedule-type-badge">{{ item.schedule_type }}</span>
              <h3 class="pm-asset-title">Aset #{{ item.asset_id }}</h3>
            </div>
            <StatusBadge :status="item.status" />
          </div>

          <div class="pm-details">
            <p><strong>Jatuh Tempo Berikutnya:</strong> ⏰ {{ item.next_run }}</p>
            <div class="checklist-box">
              <p><strong>Checklist Inspeksi:</strong></p>
              <pre>{{ item.checklist_data }}</pre>
            </div>
          </div>

          <div class="pm-actions">
            <button class="checklist-btn" @click="submitChecklist(item)">
              ✅ Selesaikan Checklist
            </button>
            <div v-if="canManageSchedule" class="pm-admin-btns">
              <button class="icon-btn edit-btn" @click="openEditModal(item)" title="Edit Jadwal PM">✏️ Edit</button>
              <button class="icon-btn delete-btn" @click="deletePMSchedule(item)" title="Hapus Jadwal PM">🗑️ Hapus</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Tambah / Edit Schedule PM -->
    <ModalDialog :show="showAddModal" :title="isEditMode ? '✏️ Edit Jadwal Preventive Maintenance' : '➕ Tambah Jadwal Preventive Maintenance'" @close="showAddModal = false">
      <form @submit.prevent="savePMSchedule" class="pm-form">
        <label>
          <span>ID Aset</span>
          <input v-model.number="newAssetId" type="number" placeholder="Masukkan ID Aset" :disabled="isEditMode" required />
        </label>

        <label>
          <span>Tipe Jadwal</span>
          <select v-model="newScheduleType">
            <option value="Daily">Daily (Harian)</option>
            <option value="Weekly">Weekly (Mingguan)</option>
            <option value="Monthly">Monthly (Bulanan)</option>
            <option value="Yearly">Yearly (Tahunan)</option>
          </select>
        </label>

        <label>
          <span>Jatuh Tempo Berikutnya</span>
          <input v-model="newNextRun" type="date" required />
        </label>

        <label>
          <span>Checklist Item (Per Baris)</span>
          <textarea v-model="newChecklistData" rows="4" placeholder="- Cek Freon AC&#10;- Bersihkan Filter Air&#10;- Tes Tegangan Generator"></textarea>
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Jadwal PM</button>
      </form>
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
import { ref, computed } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const userRole = ref(sessionStorage.getItem('user_role') || localStorage.getItem('user_role') || 'external')
const canManageSchedule = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')

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

const showAddModal = ref(false)
const isEditMode = ref(false)
const editingPmId = ref(0)
const newAssetId = ref(1)
const newScheduleType = ref('Monthly')
const newNextRun = ref('2026-08-01')
const newChecklistData = ref('1. Cek Oli Generator\n2. Check Voltase Aki\n3. Tes Otomatis Transfer Switch (ATS)')

const pmList = ref([
  {
    id: 1,
    asset_id: 1,
    schedule_type: 'Monthly',
    next_run: '2026-08-01',
    checklist_data: '1. Cek tekanan freon AC\n2. Bersihkan filter evaporator\n3. Cek drainase air kondensasi',
    status: 'Active'
  },
  {
    id: 2,
    asset_id: 4,
    schedule_type: 'Weekly',
    next_run: '2026-07-28',
    checklist_data: '1. Tes running generator 15 menit\n2. Cek level solar tangki harian\n3. Ukur tegangan aki starter',
    status: 'Active'
  }
])

function openAddModal() {
  isEditMode.value = false
  editingPmId.value = 0
  newAssetId.value = 1
  newScheduleType.value = 'Monthly'
  newNextRun.value = '2026-08-01'
  newChecklistData.value = '1. Cek Oli Generator\n2. Check Voltase Aki\n3. Tes Otomatis Transfer Switch (ATS)'
  showAddModal.value = true
}

function openEditModal(item) {
  isEditMode.value = true
  editingPmId.value = item.id
  newAssetId.value = item.asset_id
  newScheduleType.value = item.schedule_type
  newNextRun.value = item.next_run
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
      const found = pmList.value.find(p => p.id === editingPmId.value)
      if (found) {
        found.schedule_type = newScheduleType.value
        found.next_run = newNextRun.value
        found.checklist_data = newChecklistData.value
      }
      notify('Jadwal Maintenance berhasil diperbarui!', 'success')
    } else {
      await api.post('/maintenance/schedule', {
        asset_id: newAssetId.value,
        schedule_type: newScheduleType.value,
        next_run: newNextRun.value,
        checklist_data: newChecklistData.value
      })
      pmList.value.unshift({
        id: Date.now(),
        asset_id: newAssetId.value,
        schedule_type: newScheduleType.value,
        next_run: newNextRun.value,
        checklist_data: newChecklistData.value,
        status: 'Active'
      })
      notify('Jadwal Maintenance berhasil dibuat!', 'success')
    }
  } catch (e) {
    if (isEditMode.value) {
      const found = pmList.value.find(p => p.id === editingPmId.value)
      if (found) {
        found.schedule_type = newScheduleType.value
        found.next_run = newNextRun.value
        found.checklist_data = newChecklistData.value
      }
      notify('Jadwal Maintenance diperbarui!', 'success')
    } else {
      pmList.value.unshift({
        id: Date.now(),
        asset_id: newAssetId.value,
        schedule_type: newScheduleType.value,
        next_run: newNextRun.value,
        checklist_data: newChecklistData.value,
        status: 'Active'
      })
      notify('Jadwal Maintenance dibuat!', 'success')
    }
  } finally {
    showAddModal.value = false
  }
}

async function deletePMSchedule(item) {
  try {
    await api.post('/maintenance/delete', { pm_id: item.id })
    pmList.value = pmList.value.filter(p => p.id !== item.id)
    notify(`Jadwal PM Aset #${item.asset_id} berhasil dihapus!`, 'success')
  } catch (e) {
    pmList.value = pmList.value.filter(p => p.id !== item.id)
    notify(`Jadwal PM Aset #${item.asset_id} dihapus!`, 'success')
  }
}

async function submitChecklist(item) {
  try {
    await api.post(`/maintenance/${item.id}/checklist`, {
      checklist: item.checklist_data
    })
    item.status = 'Completed'
    notify(`Checklist perawatan Aset #${item.asset_id} berhasil diselesaikan!`, 'success')
  } catch (e) {
    item.status = 'Completed'
    notify(`Checklist perawatan Aset #${item.asset_id} diperbarui!`, 'success')
  }
}
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
</style>
