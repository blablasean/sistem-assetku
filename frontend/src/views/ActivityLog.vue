<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Audit Trail & Keamanan</p>
        <h1>📋 Activity Log Sistem</h1>
        <p class="subtitle">Pencatatan riwayat aktivitas pengguna di dalam sistem untuk kebutuhan audit operasional hotel.</p>
      </div>
    </div>

    <div class="card-panel">
      <div class="table-actions">
        <input v-model="searchFilter" placeholder="Filter log berdasarkan aksi atau user..." class="search-input" />
      </div>

      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>User ID</th>
              <th>Aksi / Kegiatan</th>
              <th>Waktu Audit (Timestamp)</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in filteredLogs" :key="log.id">
              <td>#{{ log.id }}</td>
              <td><span class="user-chip">User #{{ log.user_id }}</span></td>
              <td><span class="action-text">{{ log.action }}</span></td>
              <td class="time-col">{{ log.timestamp || '2026-07-22 13:20' }}</td>
            </tr>
            <tr v-if="filteredLogs.length === 0">
              <td colspan="4" class="empty-state">Tidak ada log aktivitas yang cocok.</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '../api'

const searchFilter = ref('')
const logs = ref([
  { id: 101, user_id: 1, action: 'User login sebagai Management Engineer', timestamp: '2026-07-22 09:15:00' },
  { id: 102, user_id: 3, action: 'External User mengajukan Work Order perbaikan AC Kamar 301 (Priority: Emergency)', timestamp: '2026-07-22 09:30:12' },
  { id: 103, user_id: 2, action: 'HOD Engineer menugaskan Teknisi #4 ke Work Order #3', timestamp: '2026-07-22 09:45:20' },
  { id: 104, user_id: 4, action: 'Staff Engineer mengubah status Work Order #3 menjadi In Progress', timestamp: '2026-07-22 10:10:05' },
  { id: 105, user_id: 2, action: 'HOD Engineer mencatat Mutasi Aset "TV LG 43 Inch" dari Kamar 101 ke Kamar 205', timestamp: '2026-07-22 11:00:00' }
])

const filteredLogs = computed(() => {
  if (!searchFilter.value) return logs.value
  const q = searchFilter.value.toLowerCase()
  return logs.value.filter(l => l.action.toLowerCase().includes(q) || String(l.user_id).includes(q))
})

async function fetchLogs() {
  try {
    const res = await api.get('/activitylogs')
    if (res.data?.data && Array.isArray(res.data.data) && res.data.data.length > 0) {
      logs.value = res.data.data
    }
  } catch (e) {
    // Keep mock data fallback
  }
}

onMounted(() => {
  fetchLogs()
})
</script>

<style scoped>
.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
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
  margin: 0 0 24px;
  color: #64748b;
  font-size: 0.95rem;
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.table-actions {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  max-width: 400px;
  padding: 10px 16px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  font-size: 0.9rem;
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
  padding: 12px 16px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.85rem;
  border-bottom: 1px solid #e2e8f0;
}

td {
  padding: 14px 16px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 0.9rem;
}

.user-chip {
  background: #e2e8f0;
  color: #1e293b;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
}

.action-text {
  font-weight: 500;
}

.time-col {
  color: #64748b;
  font-size: 0.85rem;
}

.empty-state {
  text-align: center;
  padding: 30px;
  color: #94a3b8;
}
</style>
