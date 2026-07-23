<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Audit Trail</p>
        <h1>📋 Activity Log</h1>
        <p class="subtitle">Riwayat Work Order Selesai dan Maintenance yang telah dikerjakan.</p>
      </div>
    </div>

    <!-- Summary boxes -->
    <div class="summary-row">
      <div class="sbox green">
        <span class="sbox-icon">✅</span>
        <div>
          <p class="sbox-label">Work Order Selesai</p>
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
    </div>

    <!-- Search -->
    <div class="search-row">
      <input v-model="searchFilter" placeholder="🔍 Cari berdasarkan deskripsi, lokasi, atau aset..." class="search-input" />
    </div>

    <!-- Section 1: Finished Work Orders -->
    <div class="card-panel">
      <h2 class="section-title">✅ Work Order Selesai (Finish)</h2>
      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Lokasi / Kamar</th>
              <th>Aset ID</th>
              <th>Prioritas</th>
              <th>Deskripsi Kerusakan</th>
              <th>Tindakan Perbaikan</th>
              <th>Biaya (Rp)</th>
              <th>Status</th>
              <th>Tanggal Selesai</th>
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
            </tr>
            <tr v-if="filteredWOs.length === 0">
              <td colspan="9" class="empty-state">Belum ada Work Order yang diselesaikan.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Section 2: Maintenance History -->
    <div class="card-panel" style="margin-top: 24px;">
      <h2 class="section-title">🔧 Riwayat Maintenance Selesai</h2>
      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Aset ID</th>
              <th>Tindakan Perawatan yang Dilakukan</th>
              <th>Biaya (Rp)</th>
              <th>Tanggal Pengerjaan</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="mh in filteredMH" :key="mh.id">
              <td><span class="wo-id">#MH-{{ mh.id }}</span></td>
              <td>Aset #{{ mh.asset_id }}</td>
              <td class="desc-cell" :title="mh.action_taken">{{ mh.action_taken }}</td>
              <td>Rp {{ formatNumber(mh.cost || 0) }}</td>
              <td class="time-col">{{ formatDate(mh.created_at) }}</td>
            </tr>
            <tr v-if="filteredMH.length === 0">
              <td colspan="5" class="empty-state">Belum ada riwayat maintenance.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import api from '../api'

const searchFilter = ref('')
const finishedWOs = ref([])
const maintenanceHistory = ref([])
const isLoading = ref(false)

const totalMaintenanceCost = computed(() =>
  maintenanceHistory.value.reduce((sum, mh) => sum + (mh.cost || 0), 0)
)

const filteredWOs = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return finishedWOs.value
  return finishedWOs.value.filter(wo =>
    (wo.description && wo.description.toLowerCase().includes(q)) ||
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
    }
  } catch (e) {
    console.error('Failed to fetch activity logs:', e)
    // Fallback demo data
    finishedWOs.value = [
      { id: 1, asset_id: 3, location: 'Kamar 201', priority: 'High', description: 'AC tidak dingin, freon habis', action_taken: 'Isi ulang freon R32, ganti kapasitor', cost: 250000, status: 'Finish', closed_at: '2026-07-22T10:30:00Z' },
      { id: 2, asset_id: 7, location: 'Lobby Utama', priority: 'Medium', description: 'Lampu koridor mati total', action_taken: 'Ganti ballast lampu LED 18W x3', cost: 85000, status: 'Finish', closed_at: '2026-07-21T14:20:00Z' }
    ]
    maintenanceHistory.value = [
      { id: 1, asset_id: 4, action_taken: 'Tes running generator 15 menit, cek level solar, cek aki starter', cost: 50000, created_at: '2026-07-20T08:00:00Z' },
      { id: 2, asset_id: 1, action_taken: 'Bersihkan filter evaporator AC, cek tekanan freon', cost: 35000, created_at: '2026-07-19T09:30:00Z' }
    ]
  } finally {
    isLoading.value = false
  }
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
  max-width: 280px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.time-col {
  color: #64748b;
  font-size: 0.82rem;
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: #94a3b8;
  font-size: 0.9rem;
}
</style>
