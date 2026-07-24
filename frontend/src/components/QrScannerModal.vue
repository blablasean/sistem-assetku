<template>
  <ModalDialog :show="show" title="📱 Pemindai QR Code Aset" maxWidth="540px" @close="handleClose">
    <div class="qr-scanner-content">
      <!-- Camera Live Stream & Viewfinder Container -->
      <div class="camera-simulator">
        <div id="qr-reader" class="html5-qr-wrapper" v-show="isCameraActive"></div>

        <div v-if="!isCameraActive" class="scanner-frame-placeholder" :class="{ 'scan-success': scanSuccessAnim }">
          <div class="corner top-left"></div>
          <div class="corner top-right"></div>
          <div class="corner bottom-left"></div>
          <div class="corner bottom-right"></div>
          <p class="placeholder-text">📷 Tekan "Buka Kamera" atau "Upload Foto QR"</p>
        </div>

        <p class="scan-instruction">
          {{ scanStatusText }}
        </p>

        <div class="camera-controls">
          <button v-if="!isCameraActive" class="cam-btn" @click="startCamera">
            🎥 Buka Kamera WebCam
          </button>
          <button v-else class="cam-btn stop-cam" @click="stopCamera">
            ⏹️ Matikan Kamera
          </button>
          <label class="cam-btn file-btn">
            📁 Upload Foto QR
            <input type="file" accept="image/*" @change="handleFileUpload" style="display: none;" />
          </label>
        </div>
      </div>

      <!-- Manual Input Section -->
      <div class="manual-input-section">
        <label>Atau ketik Kode Aset / Hasil Pemindaian:</label>
        <div class="input-group">
          <input v-model="scannedCode" placeholder="Contoh: AST-RM301-AC, AST-LBY-SOFA" @keyup.enter="handleSearch" />
          <button @click="handleSearch" :disabled="loading">
            {{ loading ? 'Mencari...' : 'Cari Aset' }}
          </button>
        </div>
      </div>

      <!-- Error Notification -->
      <div v-if="errorMsg" class="error-banner">
        ⚠️ {{ errorMsg }}
      </div>

      <!-- Scanned Result Card -->
      <div v-if="assetDetail" class="scanned-result-card">
        <div class="res-header">
          <div>
            <h4>{{ assetDetail.asset_name }}</h4>
            <span class="asset-code-badge">{{ assetDetail.asset_code }}</span>
          </div>
          <StatusBadge :status="assetDetail.status || 'Active'" />
        </div>

        <div class="res-body">
          <p><strong>Lokasi:</strong> 📍 {{ assetDetail.location || 'Kamar / Area Hotel' }}</p>
          <p><strong>Kategori:</strong> 🏷️ {{ assetDetail.category || 'General' }}</p>
          <p><strong>PIC:</strong> 👤 {{ assetDetail.pic || 'Engineering Team' }}</p>
          <p><strong>Status Reservasi:</strong> {{ assetDetail.is_reserved ? '🟠 Ter-reservasi' : '🟢 Siap Digunakan' }}</p>
        </div>

        <div class="res-actions">
          <button class="report-wo-btn" @click="reportWorkOrder(assetDetail)">
            🚨 Laporkan Kerusakan Work Order
          </button>
        </div>
      </div>
    </div>
  </ModalDialog>
</template>

<script setup>
import { ref, watch, computed, nextTick, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { Html5Qrcode } from 'html5-qrcode'
import ModalDialog from './ModalDialog.vue'
import StatusBadge from './StatusBadge.vue'
import api from '../api'

const props = defineProps({
  show: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'report-wo'])
const router = useRouter()

const isCameraActive = ref(false)
const scanSuccessAnim = ref(false)
let html5QrcodeInstance = null

const scannedCode = ref('')
const loading = ref(false)
const errorMsg = ref('')
const assetDetail = ref(null)

const scanStatusText = computed(() => {
  if (scanSuccessAnim.value) return '✅ QR Code Terdeteksi & Diproses!'
  if (isCameraActive.value) return '🎥 Kamera Aktif — Arahkan ke Stiker QR Code Aset'
  return '📷 Tekan Buka Kamera atau Upload Foto QR'
})

watch(() => props.show, (newVal) => {
  if (newVal) {
    startCamera()
  } else {
    stopCamera()
    resetState()
  }
})

async function startCamera() {
  await stopCamera()
  errorMsg.value = ''
  try {
    isCameraActive.value = true
    await nextTick()
    html5QrcodeInstance = new Html5Qrcode("qr-reader")

    let cameraConfig = { facingMode: "environment" }
    try {
      const cameras = await Html5Qrcode.getCameras()
      if (cameras && cameras.length > 0) {
        const backCam = cameras.find(c => /back|rear|environment/i.test(c.label))
        cameraConfig = backCam ? backCam.id : cameras[0].id
      }
    } catch (e) {
      console.warn('Get cameras list fallback:', e)
    }

    await html5QrcodeInstance.start(
      cameraConfig,
      {
        fps: 10,
        qrbox: (w, h) => {
          const min = Math.min(w || 250, h || 250)
          return { width: Math.max(160, Math.floor(min * 0.75)), height: Math.max(160, Math.floor(min * 0.75)) }
        }
      },
      (decodedText) => {
        if (decodedText && decodedText !== scannedCode.value) {
          onQrDetected(decodedText)
        }
      },
      () => {}
    )
  } catch (err) {
    console.warn('Camera start error:', err)
    // Fallback try with facingMode user / default
    try {
      if (html5QrcodeInstance) {
        await html5QrcodeInstance.start(
          { facingMode: "user" },
          { fps: 10, qrbox: { width: 200, height: 200 } },
          (decodedText) => {
            if (decodedText && decodedText !== scannedCode.value) {
              onQrDetected(decodedText)
            }
          },
          () => {}
        )
        isCameraActive.value = true
        return
      }
    } catch (fallbackErr) {
      console.warn('Camera fallback failed:', fallbackErr)
    }
    isCameraActive.value = false
    errorMsg.value = 'Kamera WebCam tidak dapat diakses. Silakan gunakan tombol Upload Foto QR.'
  }
}

async function stopCamera() {
  if (html5QrcodeInstance) {
    try {
      if (html5QrcodeInstance.isScanning) {
        await html5QrcodeInstance.stop()
      }
      html5QrcodeInstance.clear()
    } catch (e) {
      console.warn('Error stopping html5Qrcode:', e)
    }
    html5QrcodeInstance = null
  }
  isCameraActive.value = false
}

function triggerBeep() {
  try {
    const ctx = new (window.AudioContext || window.webkitAudioContext)()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain)
    gain.connect(ctx.destination)
    osc.type = 'sine'
    osc.frequency.value = 880
    gain.gain.setValueAtTime(0.1, ctx.currentTime)
    osc.start()
    osc.stop(ctx.currentTime + 0.15)
  } catch (e) {
    // Audio Context optional
  }
}

function onQrDetected(code) {
  triggerBeep()
  scanSuccessAnim.value = true
  scannedCode.value = code
  handleSearch()
  setTimeout(() => {
    scanSuccessAnim.value = false
  }, 2000)
}

function handleClose() {
  stopCamera()
  emit('close')
}

function resetState() {
  scannedCode.value = ''
  errorMsg.value = ''
  assetDetail.value = null
  scanSuccessAnim.value = false
}

async function handleFileUpload(event) {
  const file = event.target.files?.[0]
  if (!file) return

  loading.value = true
  errorMsg.value = ''
  assetDetail.value = null

  try {
    await stopCamera()
    isCameraActive.value = true
    await nextTick()

    const tempScanner = new Html5Qrcode("qr-reader")
    const decodedText = await tempScanner.scanFile(file, true)
    if (decodedText) {
      onQrDetected(decodedText)
      return
    }
  } catch (err) {
    console.warn('File decode error:', err)
  }

  // Check filename fallback for AST-code pattern
  const fnMatch = file.name.match(/(AST-[A-Z0-9-]+)/i)
  if (fnMatch && fnMatch[1]) {
    onQrDetected(fnMatch[1].toUpperCase())
    return
  }

  loading.value = false
  errorMsg.value = 'Foto terunggah, namun QR Code tidak terbaca. Pastikan foto stiker QR terlihat jelas.'
}

async function handleSearch() {
  if (!scannedCode.value.trim()) return
  loading.value = true
  errorMsg.value = ''
  assetDetail.value = null

  try {
    const res = await api.get('/assets/code?code=' + encodeURIComponent(scannedCode.value.trim()))
    if (res.data?.data) {
      assetDetail.value = res.data.data
    } else {
      errorMsg.value = 'Aset dengan kode tersebut tidak ditemukan.'
    }
  } catch (e) {
    errorMsg.value = 'Kode Aset ' + scannedCode.value + ' tidak ditemukan di database hotel.'
  } finally {
    loading.value = false
  }
}

function reportWorkOrder(asset) {
  stopCamera()
  emit('close')
  router.push({
    path: '/workorders',
    query: { assetId: asset.id, assetCode: asset.asset_code, location: asset.location }
  })
}

onUnmounted(() => {
  stopCamera()
})
</script>

<style scoped>
.qr-scanner-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.camera-simulator {
  background: #0f172a;
  border-radius: 2px !important;
  border: 1px solid #334155;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  position: relative;
  overflow: hidden;
  min-height: 240px;
}

.html5-qr-wrapper {
  width: 100%;
  max-width: 380px;
  border-radius: 2px !important;
  border: 1px solid #334155;
  overflow: hidden;
  margin-bottom: 10px;
}

.scanner-frame-placeholder {
  width: 180px;
  height: 180px;
  border: 2px dashed rgba(245, 158, 11, 0.6);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  border-radius: 2px !important;
  padding: 12px;
}

.placeholder-text {
  font-size: 0.78rem;
  color: #cbd5e1;
  text-align: center;
  font-weight: 600;
}

.scanner-frame-placeholder.scan-success {
  border-color: #10b981;
  box-shadow: 0 0 20px rgba(16, 185, 129, 0.6);
}

.corner {
  position: absolute;
  width: 16px;
  height: 16px;
  border-color: #f59e0b;
  border-style: solid;
}

.top-left { top: -2px; left: -2px; border-width: 3px 0 0 3px; }
.top-right { top: -2px; right: -2px; border-width: 3px 3px 0 0; }
.bottom-left { bottom: -2px; left: -2px; border-width: 0 0 3px 3px; }
.bottom-right { bottom: -2px; right: -2px; border-width: 0 3px 3px 0; }

.scan-instruction {
  margin: 8px 0 10px;
  font-size: 0.8rem;
  color: #e2e8f0;
  text-align: center;
  font-weight: 600;
}

.camera-controls {
  margin-top: 4px;
  display: flex;
  gap: 8px;
}

.cam-btn {
  background: #2563eb;
  color: white;
  border: 1px solid #1d4ed8;
  padding: 8px 14px;
  border-radius: 2px !important;
  font-size: 0.8rem;
  font-weight: 700;
  cursor: pointer;
}

.stop-cam {
  background: #dc2626;
  border-color: #b91c1c;
}

.file-btn {
  background: #0284c7;
  border-color: #0369a1;
  display: inline-flex;
  align-items: center;
}

.manual-input-section label {
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
  display: block;
  margin-bottom: 6px;
}

.input-group {
  display: flex;
  gap: 8px;
}

.input-group input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  font-size: 0.88rem;
}

.input-group button {
  background: #0f172a;
  color: white;
  border: 1px solid #1e293b;
  padding: 0 16px;
  border-radius: 2px !important;
  font-weight: 700;
  cursor: pointer;
}

.error-banner {
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fca5a5;
  padding: 10px 14px;
  border-radius: 2px !important;
  font-size: 0.85rem;
  font-weight: 600;
}

.scanned-result-card {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
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
  font-weight: 700;
}

.asset-code-badge {
  font-size: 0.75rem;
  background: #f1f5f9;
  color: #334155;
  border: 1px solid #cbd5e1;
  padding: 2px 8px;
  border-radius: 2px !important;
  font-family: monospace;
  font-weight: 700;
}

.res-body p {
  margin: 4px 0;
  font-size: 0.85rem;
  color: #334155;
}

.res-actions {
  margin-top: 14px;
}

.report-wo-btn {
  width: 100%;
  background: #dc2626;
  color: white;
  border: 1px solid #b91c1c;
  padding: 10px;
  border-radius: 2px !important;
  font-weight: 700;
  cursor: pointer;
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
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.25);
}
</style>
