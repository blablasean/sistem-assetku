<template>
  <div class="dashboard-screen page-container">
    <div class="dashboard-topbar">
      <div>
        <p class="eyebrow">AsetKu Portal</p>
        <h1>Halo, {{ userName }}</h1>
        <p class="subtitle">Ringkasan operasional aset & tiket perbaikan terkini.</p>
      </div>

      <div class="quick-action-bar">
        <button class="action-btn primary" @click="$emit('open-qr-scanner')">
          📱 Scan QR Code Aset
        </button>
        <button class="action-btn danger" @click="goToReportWO">
          🚨 Laporkan Kerusakan
        </button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="summary-cards">
      <article class="summary-card card-primary" @click="$router.push('/assets')">
        <div>
          <p class="card-title">Total Aset Hotel</p>
          <p class="card-value">{{ totalAssetsCount }}</p>
        </div>
        <p class="card-note">Data Real Terdaftar di Database</p>
      </article>

      <article class="summary-card card-secondary" @click="$router.push('/workorders')">
        <div>
          <p class="card-title">Work Order Aktif</p>
          <p class="card-value">{{ activeWorkOrdersCount }}</p>
        </div>
        <p class="card-note">{{ emergencyCount }} Tiket Darurat (Emergency)</p>
      </article>

      <article class="summary-card card-accent" @click="$router.push('/maintenance')">
        <div>
          <p class="card-title">Jadwal Perawatan Rutin</p>
          <p class="card-value">{{ workOrders.length }}</p>
        </div>
        <p class="card-note">Tiket Kerusakan Terdaftar</p>
      </article>
    </div>

    <!-- Main Grid -->
    <div class="dashboard-grid">
      <section class="activity-panel card-panel">
        <div class="panel-header">
          <div>
            <p class="panel-label">Work Order Terkini</p>
            <h2>Tiket Perbaikan Kerusakan</h2>
          </div>
          <button class="secondary-button" @click="$router.push('/workorders')">Lihat Semua</button>
        </div>

        <div class="table-responsive dashboard-wo-scroll">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Lokasi</th>
                <th>Prioritas</th>
                <th>Kerusakan</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in workOrders" :key="item.id">
                <td class="nowrap-cell">#WO-{{ item.id }}</td>
                <td class="nowrap-cell">📍 {{ item.location }}</td>
                <td class="nowrap-cell"><StatusBadge :status="item.priority" /></td>
                <td class="desc-cell">{{ item.description }}</td>
                <td class="nowrap-cell"><StatusBadge :status="item.status" /></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="insights-panel card-panel">
        <div class="panel-header">
          <div>
            <p class="panel-label">Ringkasan Kondisi</p>
            <h2>Status Operasional</h2>
          </div>
        </div>

        <div class="insights-scroll-container">
          <div class="insights-list">
            <div class="insight-card" v-for="stat in quickStats" :key="stat.label">
              <p class="insight-title">{{ stat.label }}</p>
              <p class="insight-value">{{ stat.value }}</p>
              <p class="insight-note">{{ stat.note }}</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import api from '../api'

const router = useRouter()
const userName = ref(sessionStorage.getItem('user_name') || localStorage.getItem('user_name') || 'User Hotel')

defineEmits(['open-qr-scanner'])

function goToReportWO() {
  router.push('/workorders')
}

const workOrders = ref([])
const totalAssetsCount = ref(0)
const activeWorkOrdersCount = ref(0)
const emergencyCount = ref(0)

const quickStats = ref([
  { label: 'Total Aset Terdaftar', value: '0 Unit', note: 'Jumlah total aset di database' },
  { label: 'Aset Berfungsi Normal', value: '0 Unit', note: 'Aset berstatus Active' },
  { label: 'Dalam Perawatan', value: '0 Unit', note: 'Aset berstatus Maintenance' },
  { label: 'Aset Rusak', value: '0 Unit', note: 'Aset berstatus Damaged' }
])

async function fetchDashboardData() {
  try {
    const [woRes, assetRes] = await Promise.all([
      api.get('/workorders').catch(() => null),
      api.get('/assets').catch(() => null)
    ])

    if (woRes?.data?.data && Array.isArray(woRes.data.data)) {
      workOrders.value = woRes.data.data
      activeWorkOrdersCount.value = woRes.data.data.filter(w => w.status !== 'Closed' && w.status !== 'Cancelled').length
      emergencyCount.value = woRes.data.data.filter(w => w.priority === 'Emergency' && w.status !== 'Closed').length
    }

    if (assetRes?.data?.data && Array.isArray(assetRes.data.data)) {
      totalAssetsCount.value = assetRes.data.data.length
      const activeCount = assetRes.data.data.filter(a => a.status === 'Active').length
      const maintCount = assetRes.data.data.filter(a => a.status === 'Maintenance').length
      const damagedCount = assetRes.data.data.filter(a => a.status === 'Damaged').length

      quickStats.value = [
        { label: 'Total Aset Terdaftar', value: `${totalAssetsCount.value} Unit`, note: 'Terdaftar di database hotel' },
        { label: 'Aset Berfungsi Normal', value: `${activeCount} Unit`, note: 'Status: Active' },
        { label: 'Dalam Perawatan (PM)', value: `${maintCount} Unit`, note: 'Status: Maintenance' },
        { label: 'Aset Rusak (Damaged)', value: `${damagedCount} Unit`, note: 'Status: Damaged' }
      ]
    }
  } catch (e) {
    console.error('Error fetching dashboard data:', e)
  }
}

onMounted(() => {
  fetchDashboardData()
})
</script>

<style scoped>
.dashboard-topbar {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 24px;
  align-items: center;
  margin-bottom: 28px;
}

.quick-action-bar {
  display: flex;
  gap: 12px;
}

.action-btn {
  border: none;
  border-radius: 12px;
  padding: 12px 18px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.15s ease;
}

.action-btn.primary {
  background: #2563eb;
  color: white;
}

.action-btn.danger {
  background: #dc2626;
  color: white;
}

.action-btn:hover {
  transform: translateY(-2px);
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.summary-card {
  border-radius: 20px;
  padding: 22px;
  color: #ffffff;
  min-height: 120px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: transform 0.2s;
}

.summary-card:hover {
  transform: translateY(-4px);
}

.card-primary { background: linear-gradient(135deg, #2563eb, #3b82f6); }
.card-secondary { background: linear-gradient(135deg, #d97706, #f59e0b); }
.card-accent { background: linear-gradient(135deg, #059669, #10b981); }

.card-title {
  margin: 0;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.85);
}

.card-value {
  margin: 4px 0 0;
  font-size: 2.2rem;
  font-weight: 800;
}

.card-note {
  margin: 0;
  font-size: 0.8rem;
  color: rgba(255, 255, 255, 0.75);
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.panel-label {
  margin: 0 0 4px;
  color: #2563eb;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.75rem;
  font-weight: 700;
}

.activity-panel h2,
.insights-panel h2 {
  margin: 0;
  font-size: 1.2rem;
  color: #0f172a;
}

.secondary-button {
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 6px 12px;
  background: #ffffff;
  color: #334155;
  font-weight: 600;
  font-size: 0.85rem;
  cursor: pointer;
}

/* Dashboard Scrollable Tables and Scrollable Insights Panel */
.dashboard-wo-scroll {
  max-height: 320px;
  overflow-y: auto;
  overflow-x: auto;
  border-radius: 12px;
}

.dashboard-wo-scroll table th {
  position: sticky;
  top: 0;
  background: #f8fafc;
  z-index: 5;
  white-space: nowrap;
}

.nowrap-cell {
  white-space: nowrap;
}

.desc-cell {
  white-space: nowrap;
  max-width: 250px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.insights-scroll-container {
  max-height: 320px;
  overflow-y: auto;
  padding-right: 4px;
}

.insights-list {
  display: grid;
  gap: 12px;
}

.insight-card {
  background: #f8fafc;
  border-radius: 14px;
  padding: 14px;
  border: 1px solid #f1f5f9;
}

.insight-title {
  margin: 0 0 4px;
  color: #64748b;
  font-size: 0.8rem;
  font-weight: 600;
}

.insight-value {
  margin: 0;
  font-size: 1.4rem;
  font-weight: 800;
  color: #0f172a;
}

.insight-note {
  margin: 2px 0 0;
  color: #94a3b8;
  font-size: 0.78rem;
}

@media (max-width: 960px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style>
