<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Preventive Maintenance Management</p>
        <h1>📅 Jadwal & Perawatan Rutin Aset Hotel</h1>
        <p class="subtitle">Penjadwalan inspeksi rutin harian/mingguan/bulanan untuk mencegah kerusakan fatal pada aset hotel.</p>
      </div>
      <button class="primary-btn" v-if="canManageSchedule" @click="showAddModal = true">
        ➕ Tambah Jadwal PM Baru
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
              ✅ Selesaikan Checklist & Update Maintenance
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal Tambah Schedule PM -->
    <ModalDialog :show="showAddModal" title="Tambah Jadwal Preventive Maintenance" @close="showAddModal = false">
      <form @submit.prevent="createPMSchedule" class="pm-form">
        <label>
          <span>ID Aset</span>
          <input v-model.number="newAssetId" type="number" placeholder="Masukkan ID Aset" required />
        </label>

        <label>
          <span>Tipe Jadwal</span>
          <select v-model="newScheduleType">
            <option>Daily (Harian)</option>
            <option>Weekly (Mingguan)</option>
            <option>Monthly (Bulanan)</option>
            <option>Yearly (Tahunan)</option>
          </select>
        </label>

        <label>
          <span>Checklist Item (Per Baris)</span>
          <textarea v-model="newChecklistData" rows="4" placeholder="- Cek Freon AC&#10;- Bersihkan Filter Air&#10;- Tes Tegangan Generator"></textarea>
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Jadwal PM</button>
      </form>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'

const userRole = ref(localStorage.getItem('user_role') || 'external')
const canManageSchedule = computed(() => userRole.value === 'hod' || userRole.value === 'management')

const showAddModal = ref(false)
const newAssetId = ref(1)
const newScheduleType = ref('Monthly (Bulanan)')
const newChecklistData = ref('- Cek Oli Generator\n- Check Voltase Aki\n- Tes Otomatis Transfer Switch (ATS)')

const pmList = ref([
  {
    id: 1,
    asset_id: 101,
    schedule_type: 'Monthly',
    next_run: '2026-08-01',
    checklist_data: '1. Cek tekanan freon AC\n2. Bersihkan filter evaporator\n3. Cek drainase air kondensasi',
    status: 'Active'
  },
  {
    id: 2,
    asset_id: 104,
    schedule_type: 'Weekly',
    next_run: '2026-07-28',
    checklist_data: '1. Tes running generator 15 menit\n2. Cek level solar tangki harian\n3. Ukur tegangan aki starter',
    status: 'Active'
  },
  {
    id: 3,
    asset_id: 108,
    schedule_type: 'Daily',
    next_run: '2026-07-23',
    checklist_data: '1. Cek suhu Chiller Dapur (-18°C)\n2. Pastikan pintu seal karet rapat\n3. Bersihkan kondensor outdoor',
    status: 'In Progress'
  }
])

function createPMSchedule() {
  pmList.value.unshift({
    id: Date.now(),
    asset_id: newAssetId.value,
    schedule_type: newScheduleType.value,
    next_run: '2026-08-15',
    checklist_data: newChecklistData.value,
    status: 'Active'
  })
  showAddModal.value = false
  alert('Jadwal Preventive Maintenance berhasil dibuat!')
}

function submitChecklist(item) {
  item.status = 'Completed'
  alert(`Checklist perawatan rutin untuk Aset #${item.asset_id} selesai! Tanggal Last Maintenance diperbarui.`)
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

.pm-details p {
  margin: 4px 0;
  font-size: 0.85rem;
  color: #475569;
}

.checklist-box {
  margin-top: 10px;
  background: #f8fafc;
  border-radius: 10px;
  padding: 10px;
  border: 1px solid #f1f5f9;
}

.checklist-box pre {
  margin: 4px 0 0;
  white-space: pre-wrap;
  font-family: inherit;
  font-size: 0.8rem;
  color: #334155;
}

.pm-actions {
  margin-top: 16px;
}

.checklist-btn {
  width: 100%;
  background: #059669;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
}

.checklist-btn:hover {
  background: #047857;
}

.pm-form {
  display: grid;
  gap: 16px;
}

.pm-form label {
  display: grid;
  gap: 6px;
  font-weight: 600;
  color: #1e293b;
}

.pm-form input, .pm-form select, .pm-form textarea {
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
</style>
