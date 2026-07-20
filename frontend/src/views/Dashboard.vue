<template>
  <div class="dashboard-screen">
    <div class="dashboard-topbar">
      <div>
        <p class="eyebrow">Dashboard AsetKu</p>
        <h1>Halo, Admin</h1>
        <p class="subtitle">Ringkasan operasional dan status aset terbaru.</p>
      </div>
      <button class="logout-button" @click="$emit('logout')">Keluar</button>
    </div>

    <div class="summary-cards">
      <article class="summary-card card-primary">
        <p class="card-title">Total Aset</p>
        <p class="card-value">128</p>
      </article>
      <article class="summary-card card-secondary">
        <p class="card-title">Work Order</p>
        <p class="card-value">34</p>
      </article>
      <article class="summary-card card-accent">
        <p class="card-title">Maintenance</p>
        <p class="card-value">12</p>
      </article>
    </div>

    <div class="dashboard-grid">
      <section class="activity-panel">
        <div class="panel-header">
          <div>
            <p class="panel-label">Aktivitas Terbaru</p>
            <h2>Work order terbaru</h2>
          </div>
          <button class="secondary-button">Lihat semua</button>
        </div>

        <table>
          <thead>
            <tr>
              <th>No.</th>
              <th>Nama Aset</th>
              <th>Status</th>
              <th>Deadline</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in workOrders" :key="item.id">
              <td>{{ index + 1 }}</td>
              <td>{{ item.asset }}</td>
              <td><span :class="['status-chip', item.variant]">{{ item.status }}</span></td>
              <td>{{ item.deadline }}</td>
            </tr>
          </tbody>
        </table>
      </section>

      <section class="insights-panel">
        <div class="panel-header">
          <div>
            <p class="panel-label">Ringkasan Cepat</p>
            <h2>Kondisi aset</h2>
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
const workOrders = [
  { id: 1, asset: 'Generator Unit 1', status: 'Dalam Proses', variant: 'warning', deadline: '25 Jul' },
  { id: 2, asset: 'AC Ruang Server', status: 'Selesai', variant: 'success', deadline: '22 Jul' },
  { id: 3, asset: 'Lift Gedung A', status: 'Menunggu', variant: 'default', deadline: '27 Jul' },
  { id: 4, asset: 'Pompa Air', status: 'Dalam Proses', variant: 'warning', deadline: '29 Jul' }
]

const quickStats = [
  { label: 'Aset Aktif', value: '106', note: '85% dari total aset aktif' },
  { label: 'Permintaan Maintenance', value: '8', note: 'Permintaan dalam 7 hari terakhir' },
  { label: 'Inspeksi Terlewat', value: '3', note: 'Butuh perhatian segera' }
]
</script>

<style scoped>
.dashboard-screen {
  max-width: 1200px;
  margin: 0 auto;
  padding: 32px;
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
  margin: 0 0 8px;
  text-transform: uppercase;
  letter-spacing: 0.14em;
  font-size: 0.8rem;
  color: #2563eb;
}

h1 {
  margin: 0;
  font-size: clamp(2rem, 2.5vw, 3rem);
}

.subtitle {
  margin: 8px 0 0;
  color: #6b7280;
  max-width: 620px;
}

.logout-button,
.secondary-button {
  border: none;
  border-radius: 14px;
  padding: 12px 18px;
  background: #1d4ed8;
  color: #ffffff;
  font-weight: 700;
  cursor: pointer;
}

.logout-button {
  min-width: 120px;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.summary-card {
  border-radius: 24px;
  padding: 26px;
  color: #ffffff;
  min-height: 132px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.card-primary {
  background: linear-gradient(135deg, #2563eb, #3b82f6);
}

.card-secondary {
  background: linear-gradient(135deg, #0f766e, #14b8a6);
}

.card-accent {
  background: linear-gradient(135deg, #be185d, #ec4899);
}

.card-title {
  margin: 0;
  font-size: 0.95rem;
  color: rgba(255, 255, 255, 0.85);
}

.card-value {
  margin: 0;
  font-size: 2.5rem;
  font-weight: 700;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 24px;
}

.activity-panel,
.insights-panel {
  background: #ffffff;
  border-radius: 28px;
  padding: 28px;
  box-shadow: 0 18px 45px rgba(15, 23, 42, 0.08);
}

.panel-header {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: center;
  margin-bottom: 22px;
}

.panel-label {
  margin: 0 0 8px;
  color: #2563eb;
  text-transform: uppercase;
  letter-spacing: 0.12em;
  font-size: 0.75rem;
}

.activity-panel h2,
.insights-panel h2 {
  margin: 0;
  font-size: 1.35rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 100%;
}

thead th {
  text-align: left;
  color: #475569;
  font-size: 0.9rem;
  padding: 14px 12px;
  border-bottom: 1px solid #e2e8f0;
}

tbody tr:hover {
  background: #f8fafc;
}

tbody td {
  padding: 16px 12px;
  color: #334155;
  border-bottom: 1px solid #e2e8f0;
}

.status-chip {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 999px;
  padding: 8px 12px;
  font-size: 0.85rem;
  font-weight: 600;
}

.status-chip.success {
  background: #d1fae5;
  color: #065f46;
}

.status-chip.warning {
  background: #fef3c7;
  color: #92400e;
}

.status-chip.default {
  background: #e2e8f0;
  color: #334155;
}

.insights-list {
  display: grid;
  gap: 18px;
}

.insight-card {
  background: #f8fafc;
  border-radius: 22px;
  padding: 20px;
}

.insight-title {
  margin: 0 0 10px;
  color: #475569;
  font-size: 0.95rem;
}

.insight-value {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
}

.insight-note {
  margin: 10px 0 0;
  color: #6b7280;
  font-size: 0.95rem;
}

@media (max-width: 960px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .summary-cards {
    grid-template-columns: 1fr;
  }
}
</style>
