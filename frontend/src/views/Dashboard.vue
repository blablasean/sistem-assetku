<template>
  <div class="dashboard-screen">
    <div class="dashboard-topbar">
      <div>
        <p class="eyebrow">AsetKu Hotel Portal</p>
        <h1>Halo, {{ userName }}</h1>
        <p class="subtitle">Ringkasan operasional manajemen aset hotel & pelacakan tiket perbaikan terkini.</p>
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
          <p class="card-value">128</p>
        </div>
        <p class="card-note">85% Aset Berfungsi Normal</p>
      </article>

      <article class="summary-card card-secondary" @click="$router.push('/workorders')">
        <div>
          <p class="card-title">Work Order Aktif</p>
          <p class="card-value">14</p>
        </div>
        <p class="card-note">4 Tiket Darurat (Emergency)</p>
      </article>

      <article class="summary-card card-accent" @click="$router.push('/maintenance')">
        <div>
          <p class="card-title">Jadwal Perawatan Rutin</p>
          <p class="card-value">8</p>
        </div>
        <p class="card-note">Penjadwalan PM Minggu Ini</p>
      </article>

      <article class="summary-card card-dark" @click="$router.push('/utility')">
        <div>
          <p class="card-title">Utility Monitoring</p>
          <p class="card-value">Normal</p>
        </div>
        <p class="card-note">Listrik & Air PDAM Terkendali</p>
      </article>
    </div>

    <!-- Main Grid -->
    <div class="dashboard-grid">
      <section class="activity-panel">
        <div class="panel-header">
          <div>
            <p class="panel-label">Work Order Terkini</p>
            <h2>Tiket Perbaikan Kerusakan</h2>
          </div>
          <button class="secondary-button" @click="$router.push('/workorders')">Lihat Semua</button>
        </div>

        <div class="table-responsive">
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
                <td>#WO-{{ item.id }}</td>
                <td>📍 {{ item.location }}</td>
                <td><StatusBadge :status="item.priority" /></td>
                <td>{{ item.description }}</td>
                <td><StatusBadge :status="item.status" /></td>
              </tr>
            </tbody>
          </table>
        </div>
      </section>

      <section class="insights-panel">
        <div class="panel-header">
          <div>
            <p class="panel-label">Ringkasan Kondisi</p>
            <h2>Status Operasional</h2>
          </div>
        </div>

        <div class="insights-list">
          <div class="insight-card" v-for="stat in quickStats" :key="stat.label">
            <p class="insight-title">{{ stat.label }}</p>
            <p class="insight-value">{{ stat.value }}</p>
            <p class="insight-note">{{ stat.note }}</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'

const router = useRouter()
const userName = ref(localStorage.getItem('user_name') || 'User Hotel')

defineEmits(['open-qr-scanner'])

function goToReportWO() {
  router.push('/workorders')
}

const workOrders = [
  { id: 1, location: 'Kamar 301', priority: 'Emergency', description: 'AC Kamar 301 bocor air dan tidak dingin', status: 'In Progress' },
  { id: 2, location: 'Kitchen Dapur', priority: 'High', description: 'Chiller Dapur Utama suhu naik ke -5°C', status: 'Open' },
  { id: 3, location: 'Kamar 102', priority: 'Medium', description: 'Smart TV HDMI port tidak terdeteksi', status: 'Closed' },
  { id: 4, location: 'Lobby Lounge', priority: 'Low', description: 'Kaki sofa kendor perlu dikencangkan', status: 'Open' }
]

const quickStats = [
  { label: 'Aset Aktif', value: '106 Unit', note: 'Siap digunakan di area hotel' },
  { label: 'Aset Maintenance', value: '12 Unit', note: 'Sedang dalam perawatan rutin' },
  { label: 'Kerusakan Terlapor', value: '4 Unit', note: 'Menunggu respon teknisi' }
]
</script>

<style scoped>
.dashboard-screen {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px 24px;
}

.dashboard-topbar {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 24px;
  align-items: center;
  margin-bottom: 32px;
}

.eyebrow {
  margin: 0 0 4px;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 0.8rem;
  color: #2563eb;
  font-weight: 700;
}

h1 {
  margin: 0;
  font-size: clamp(1.8rem, 2.5vw, 2.5rem);
  color: #0f172a;
}

.subtitle {
  margin: 6px 0 0;
  color: #64748b;
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
.card-dark { background: linear-gradient(135deg, #334155, #475569); }

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

.activity-panel,
.insights-panel {
  background: #ffffff;
  border-radius: 24px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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
  font-size: 1.25rem;
  color: #0f172a;
}

.secondary-button {
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 8px 14px;
  background: #ffffff;
  color: #334155;
  font-weight: 600;
  cursor: pointer;
}

table {
  width: 100%;
  border-collapse: collapse;
}

thead th {
  text-align: left;
  color: #475569;
  font-size: 0.85rem;
  padding: 10px 12px;
  border-bottom: 1px solid #e2e8f0;
}

tbody td {
  padding: 12px;
  color: #334155;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.85rem;
}

.insights-list {
  display: grid;
  gap: 14px;
}

.insight-card {
  background: #f8fafc;
  border-radius: 16px;
  padding: 16px;
}

.insight-title {
  margin: 0 0 6px;
  color: #475569;
  font-size: 0.85rem;
}

.insight-value {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 700;
  color: #0f172a;
}

.insight-note {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 0.8rem;
}

@media (max-width: 960px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style>
