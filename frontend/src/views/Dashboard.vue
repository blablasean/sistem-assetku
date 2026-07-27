<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">AsetKu Portal</p>
        <h1>Halo, {{ userName }}</h1>
        <p class="subtitle">Ringkasan operasional hotel.</p>
      </div>

      <div class="header-action-group">
        <button class="primary-btn" @click="$emit('open-qr-scanner')">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/><path d="M21 21v.01"/><path d="M12 7v3a2 2 0 0 1-2 2H7"/><path d="M3 12h.01"/><path d="M12 3h.01"/><path d="M12 16v.01"/><path d="M16 12h1"/><path d="M21 12v.01"/><path d="M12 21v-1"/></svg>
          <span>Scan QR Code</span>
        </button>
        <button class="primary-btn btn-danger-ios" @click="goToReportWO">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3Z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>
          <span>Laporkan Kerusakan</span>
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

.header-action-group {
  display: flex;
  align-items: center;
  gap: 10px;
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
  white-space: nowrap;
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.primary-btn.btn-danger-ios {
  background: #ff3b30 !important;
  border-color: #ff3b30 !important;
  box-shadow: 0 4px 12px rgba(255, 59, 48, 0.25);
}

.primary-btn.btn-danger-ios:hover {
  background: #d70015 !important;
  border-color: #d70015 !important;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.summary-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 6px !important;
  padding: 22px;
  color: #0f172a;
  min-height: 120px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  cursor: pointer;
  transition: all 0.15s ease;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.04);
}

.summary-card:hover {
  border-color: #cbd5e1;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

.card-primary { border-top: 4px solid #2563eb; }
.card-secondary { border-top: 4px solid #d97706; }
.card-accent { border-top: 4px solid #059669; }

.card-title {
  margin: 0;
  font-size: 0.85rem;
  color: #64748b;
  font-weight: 700;
}

.card-value {
  margin: 4px 0 0;
  font-size: 2.2rem;
  font-weight: 800;
  color: #0f172a;
  letter-spacing: -0.02em;
}

.card-note {
  margin: 0;
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
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
