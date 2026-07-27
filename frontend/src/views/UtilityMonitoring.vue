<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Utility Monitoring & Cost Tracking</p>
        <h1>⚡ Pemantauan Utilitas & Energi (Listrik & Air PDAM)</h1>
        <p class="subtitle">Monitoring konsumsi energi dan pemakaian air harian/bulanan untuk efisiensi biaya operasional.</p>
      </div>

      <button class="primary-btn" @click="showLogModal = true">
        📝 Catat Meter Hari Ini
      </button>
    </div>

    <!-- Summary Metrics -->
    <div class="metrics-grid">
      <div class="metric-card electricity">
        <div class="metric-icon">⚡</div>
        <div>
          <p class="metric-label">Konsumsi Listrik Bulan Ini</p>
          <p class="metric-val">14,250 kWh</p>
          <p class="metric-cost">Estimasi Biaya: Rp 21,375,000</p>
        </div>
      </div>

      <div class="metric-card water">
        <div class="metric-icon">💧</div>
        <div>
          <p class="metric-label">Konsumsi Air PDAM Bulan Ini</p>
          <p class="metric-val">1,820 m³</p>
          <p class="metric-cost">Estimasi Biaya: Rp 9,100,000</p>
        </div>
      </div>
    </div>

    <!-- Meter Reading Logs -->
    <div class="card-panel">
      <h3>📊 Log Pencatatan Meter Utilitas Harian</h3>
      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>Tanggal</th>
              <th>Meter Listrik (kWh)</th>
              <th>Meter Air (m³)</th>
              <th>Pencatat / Petugas</th>
              <th>Catatan Operasional</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="log in utilityLogs" :key="log.id">
              <td>📅 {{ log.date }}</td>
              <td>⚡ {{ log.electricity_kwh }} kWh</td>
              <td>💧 {{ log.water_m3 }} m³</td>
              <td>👤 {{ log.recorder }}</td>
              <td>{{ log.notes }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Modal Catat Meter -->
    <ModalDialog :show="showLogModal" title="📝 Catat Angka Meter Utilitas Harian" @close="showLogModal = false">
      <form @submit.prevent="submitMeterLog" class="modal-form">
        <label>
          <span>Tanggal Pencatatan</span>
          <input v-model="logDate" type="date" required />
        </label>
        <label>
          <span>Angka Stand Meter Listrik (kWh)</span>
          <input v-model.number="logPower" type="number" placeholder="Contoh: 14250" required />
        </label>
        <label>
          <span>Angka Stand Meter Air PDAM (m³)</span>
          <input v-model.number="logWater" type="number" placeholder="Contoh: 1820" required />
        </label>
        <label>
          <span>Catatan / Observasi Suhu Chiller</span>
          <input v-model="logNotes" placeholder="Misal: Pemakaian normal, chiller 1 beroperasi penuh" />
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Log Meter</button>
      </form>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ModalDialog from '../components/ModalDialog.vue'

const showLogModal = ref(false)
const logDate = ref('2026-07-22')
const logPower = ref(14250)
const logWater = ref(1820)
const logNotes = ref('Operasional normal, beban fasilitas 85%')

const utilityLogs = ref([
  { id: 1, date: '2026-07-22', electricity_kwh: 14250, water_m3: 1820, recorder: 'Teknisi Eko', notes: 'Beban fasilitas 85%, sistem AC beroperasi penuh' },
  { id: 2, date: '2026-07-21', electricity_kwh: 13800, water_m3: 1750, recorder: 'Teknisi Deni', notes: 'Pemakaian normal harian' },
  { id: 3, date: '2026-07-20', electricity_kwh: 14100, water_m3: 1790, recorder: 'Teknisi Budi', notes: 'Ada event banquet di ballroom' }
])

function submitMeterLog() {
  utilityLogs.value.unshift({
    id: Date.now(),
    date: logDate.value,
    electricity_kwh: logPower.value,
    water_m3: logWater.value,
    recorder: 'User Logged In',
    notes: logNotes.value
  })
  showLogModal.value = false
  alert('Pencatatan stand meter utilitas harian berhasil disimpan!')
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

.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 20px;
  margin-bottom: 28px;
}

.metric-card {
  border-radius: 20px;
  padding: 24px;
  color: white;
  display: flex;
  align-items: center;
  gap: 20px;
}

.metric-card.electricity { background: linear-gradient(135deg, #0284c7, #2563eb); }
.metric-card.water { background: linear-gradient(135deg, #0d9488, #059669); }

.metric-icon {
  font-size: 2.5rem;
  background: rgba(255, 255, 255, 0.2);
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.metric-label {
  margin: 0;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.85);
}

.metric-val {
  margin: 4px 0;
  font-size: 2rem;
  font-weight: 800;
}

.metric-cost {
  margin: 0;
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.8);
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.card-panel h3 {
  margin: 0 0 16px;
  font-size: 1.15rem;
  color: #0f172a;
}

.table-responsive { overflow-x: auto; }

table {
  width: 100%;
  border-collapse: collapse;
}

th {
  text-align: left;
  padding: 12px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.85rem;
  border-bottom: 1px solid #e2e8f0;
}

td {
  padding: 14px 12px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.9rem;
  color: #334155;
}

.modal-form {
  display: grid;
  gap: 16px;
}

.modal-form label {
  display: grid;
  gap: 6px;
  font-weight: 600;
  color: #1e293b;
}

.modal-form input {
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
