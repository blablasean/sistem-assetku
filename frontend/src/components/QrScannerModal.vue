<template>
  <ModalDialog :show="show" title="📱 Pemindai QR Code Aset" maxWidth="500px" @close="$emit('close')">
    <div class="qr-scanner-content">
      <div class="camera-simulator">
        <div class="scanner-frame">
          <div class="laser-line"></div>
          <div class="corner top-left"></div>
          <div class="corner top-right"></div>
          <div class="corner bottom-left"></div>
          <div class="corner bottom-right"></div>
        </div>
        <p class="scan-instruction">Arahkan kamera ke Kode QR Aset di unit fisik...</p>
      </div>

      <div class="manual-input-section">
        <label>Atau masukkan Kode Aset / Hasil Scan QR:</label>
        <div class="input-group">
          <input v-model="scannedCode" placeholder="Contoh: AST-KCH-01, AST-RM301-AC" @keyup.enter="handleSearch" />
          <button @click="handleSearch" :disabled="loading">
            {{ loading ? 'Mencari...' : 'Cari Aset' }}
          </button>
        </div>
      </div>

      <div v-if="errorMsg" class="error-banner">
        ⚠️ {{ errorMsg }}
      </div>

      <div v-if="assetDetail" class="scanned-result-card">
        <div class="res-header">
          <div>
            <h4>{{ assetDetail.asset_name }}</h4>
            <span class="asset-code-badge">{{ assetDetail.asset_code }}</span>
          </div>
          <StatusBadge :status="assetDetail.status" />
        </div>

        <div class="res-body">
          <p><strong>Lokasi:</strong> 📍 {{ assetDetail.location || 'Kamar / Area Hotel' }}</p>
          <p><strong>Kategori:</strong> 🏷️ {{ assetDetail.category || 'General' }}</p>
          <p><strong>PIC:</strong> 👤 {{ assetDetail.pic || 'Engineering Team' }}</p>
          <p><strong>Last Maintenance:</strong> 📅 22 Juli 2026</p>
        </div>

        <div class="res-actions">
          <button class="report-wo-btn" @click="reportWorkOrder(assetDetail)">
            🚨 Laporkan Kerusakan Aset Ini
          </button>
        </div>
      </div>
    </div>
  </ModalDialog>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import ModalDialog from './ModalDialog.vue'
import StatusBadge from './StatusBadge.vue'
import api from '../api'

defineProps({
  show: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'report-wo'])
const router = useRouter()

const scannedCode = ref('')
const loading = ref(false)
const errorMsg = ref('')
const assetDetail = ref(null)

async function handleSearch() {
  if (!scannedCode.value.trim()) return
  loading.value = true
  errorMsg.value = ''
  assetDetail.value = null

  try {
    const res = await api.get('/assets/code?code=' + encodeURIComponent(scannedCode.value.trim()))
    assetDetail.value = res.data.data
  } catch (e) {
    errorMsg.value = 'Aset dengan kode QR tersebut tidak ditemukan di database server.'
  } finally {
    loading.value = false
  }
}

function reportWorkOrder(asset) {
  emit('close')
  router.push({
    path: '/workorders',
    query: { assetId: asset.id, assetCode: asset.asset_code, location: asset.location }
  })
}
</script>

<style scoped>
.qr-scanner-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.camera-simulator {
  background: #0f172a;
  border-radius: 16px;
  padding: 30px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  position: relative;
  overflow: hidden;
}

.scanner-frame {
  width: 160px;
  height: 160px;
  border: 2px dashed rgba(56, 189, 248, 0.4);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.laser-line {
  position: absolute;
  width: 100%;
  height: 3px;
  background: #38bdf8;
  box-shadow: 0 0 12px #38bdf8;
  animation: scan 2s infinite ease-in-out;
}

@keyframes scan {
  0% { top: 5%; }
  50% { top: 90%; }
  100% { top: 5%; }
}

.corner {
  position: absolute;
  width: 16px;
  height: 16px;
  border-color: #38bdf8;
  border-style: solid;
}

.top-left { top: -2px; left: -2px; border-width: 3px 0 0 3px; }
.top-right { top: -2px; right: -2px; border-width: 3px 3px 0 0; }
.bottom-left { bottom: -2px; left: -2px; border-width: 0 0 3px 3px; }
.bottom-right { bottom: -2px; right: -2px; border-width: 0 3px 3px 0; }

.scan-instruction {
  margin: 16px 0 0;
  font-size: 0.8rem;
  color: #94a3b8;
  text-align: center;
}

.manual-input-section label {
  font-size: 0.85rem;
  font-weight: 600;
  color: #334155;
}

.input-group {
  display: flex;
  gap: 8px;
  margin-top: 6px;
}

.input-group input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  font-size: 0.9rem;
}

.input-group button {
  background: #2563eb;
  color: white;
  border: none;
  padding: 0 16px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
}

.error-banner {
  background: #fef2f2;
  color: #991b1b;
  padding: 12px;
  border-radius: 10px;
  font-size: 0.85rem;
}

.scanned-result-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 16px;
}

.res-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.res-header h4 {
  margin: 0;
  font-size: 1.05rem;
  color: #0f172a;
}

.asset-code-badge {
  font-size: 0.75rem;
  background: #e2e8f0;
  color: #334155;
  padding: 2px 8px;
  border-radius: 6px;
  font-family: monospace;
  font-weight: 700;
}

.res-body p {
  margin: 4px 0;
  font-size: 0.85rem;
  color: #475569;
}

.res-actions {
  margin-top: 14px;
}

.report-wo-btn {
  width: 100%;
  background: #dc2626;
  color: white;
  border: none;
  padding: 10px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: background 0.15s;
}

.report-wo-btn:hover {
  background: #b91c1c;
}
</style>
