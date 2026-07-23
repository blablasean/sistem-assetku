<template>
  <ModalDialog :show="show" title="📱 Pemindai QR Code Aset" maxWidth="540px" @close="handleClose">
    <div class="qr-scanner-content">
      <!-- Camera Live Stream & Viewfinder -->
      <div class="camera-simulator">
        <video ref="videoRef" class="live-video" autoplay playsinline muted v-show="isCameraActive"></video>
        <canvas ref="canvasRef" style="display: none;"></canvas>

        <div class="scanner-frame" :class="{ 'scan-success': scanSuccessAnim }">
          <div class="laser-line"></div>
          <div class="corner top-left"></div>
          <div class="corner top-right"></div>
          <div class="corner bottom-left"></div>
          <div class="corner bottom-right"></div>
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
import { ref, watch, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import jsQR from 'jsqr'
import ModalDialog from './ModalDialog.vue'
import StatusBadge from './StatusBadge.vue'
import api from '../api'

const props = defineProps({
  show: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'report-wo'])
const router = useRouter()

const videoRef = ref(null)
const canvasRef = ref(null)
const isCameraActive = ref(false)
const scanSuccessAnim = ref(false)
let mediaStream = null
let scanAnimFrame = null

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
  try {
    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
      mediaStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: 'environment', width: { ideal: 640 }, height: { ideal: 480 } }
      })
      if (videoRef.value) {
        videoRef.value.srcObject = mediaStream
      }
      isCameraActive.value = true
      startScanLoop()
    }
  } catch (err) {
    console.warn('WebCam camera access not allowed or unavailable:', err)
    isCameraActive.value = false
  }
}

function decodeQrFromCanvas(canvas) {
  if (!canvas || !canvas.width || !canvas.height) return null
  try {
    const ctx = canvas.getContext('2d', { willReadFrequently: true })
    const imgData = ctx.getImageData(0, 0, canvas.width, canvas.height)
    if (!imgData || !imgData.data || !imgData.width || !imgData.height) return null

    // Attempt 1: Standard QR scan
    let code = jsQR(imgData.data, imgData.width, imgData.height, {
      inversionAttempts: 'dontInvert'
    })
    if (code && code.data) return code.data

    // Attempt 2: Inverted QR scan
    code = jsQR(imgData.data, imgData.width, imgData.height, {
      inversionAttempts: 'onlyInvert'
    })
    if (code && code.data) return code.data
  } catch (e) {
    console.warn('jsQR decoding error:', e)
  }
  return null
}

let scanTimer = null

function startScanLoop() {
  if (!isCameraActive.value) return
  if (scanTimer) clearInterval(scanTimer)

  scanTimer = setInterval(() => {
    if (!isCameraActive.value || loading.value) return

    if (videoRef.value && videoRef.value.readyState === videoRef.value.HAVE_ENOUGH_DATA) {
      if (canvasRef.value) {
        const canvas = canvasRef.value
        const w = videoRef.value.videoWidth || 320
        const h = videoRef.value.videoHeight || 240
        if (w > 0 && h > 0) {
          canvas.width = w
          canvas.height = h
          const ctx = canvas.getContext('2d', { willReadFrequently: true })
          ctx.drawImage(videoRef.value, 0, 0, w, h)

          const decoded = decodeQrFromCanvas(canvas)
          if (decoded && decoded !== scannedCode.value) {
            onQrDetected(decoded)
          }
        }
      }
    }
  }, 400) // Throttled to 400ms -> smooth & fast!
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

function stopCamera() {
  if (scanTimer) {
    clearInterval(scanTimer)
    scanTimer = null
  }
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => track.stop())
    mediaStream = null
  }
  isCameraActive.value = false
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

  // 1. Check filename for asset code pattern (e.g. AST-LBY-SOFA, AST-RM301-AC)
  const fnMatch = file.name.match(/(AST-[A-Z0-9-]+)/i)
  if (fnMatch && fnMatch[1]) {
    onQrDetected(fnMatch[1].toUpperCase())
    return
  }

  const reader = new FileReader()
  reader.onload = (e) => {
    const img = new Image()
    img.crossOrigin = 'Anonymous'
    img.onload = () => {
      if (canvasRef.value) {
        const canvas = canvasRef.value
        const w = img.naturalWidth || img.width || 400
        const h = img.naturalHeight || img.height || 400
        if (w > 0 && h > 0) {
          canvas.width = w
          canvas.height = h
          const ctx = canvas.getContext('2d', { willReadFrequently: true })
          ctx.drawImage(img, 0, 0, w, h)

          const decoded = decodeQrFromCanvas(canvas)
          if (decoded) {
            onQrDetected(decoded)
            return
          }
        }
      }

      loading.value = false
      errorMsg.value = 'Foto terunggah, namun QR Code tidak terdeteksi otomatis. Silakan ketik Kode Aset (Contoh: AST-LBY-SOFA) di bawah.'
    }
    img.src = e.target.result
  }
  reader.readAsDataURL(file)
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
  border-radius: 16px;
  padding: 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: white;
  position: relative;
  overflow: hidden;
  min-height: 220px;
}

.live-video {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  opacity: 0.65;
}

.scanner-frame {
  width: 150px;
  height: 150px;
  border: 2px dashed rgba(56, 189, 248, 0.6);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
  transition: all 0.3s ease;
}

.scanner-frame.scan-success {
  border-color: #10b981;
  box-shadow: 0 0 20px rgba(16, 185, 129, 0.6);
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
  margin: 14px 0 6px;
  font-size: 0.8rem;
  color: #e2e8f0;
  text-align: center;
  z-index: 2;
  font-weight: 600;
  text-shadow: 0 1px 4px rgba(0,0,0,0.8);
}

.camera-controls {
  z-index: 2;
  margin-top: 4px;
  display: flex;
  gap: 8px;
}

.cam-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 6px 14px;
  border-radius: 8px;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
}

.stop-cam {
  background: #dc2626;
}

.file-btn {
  background: #0284c7;
  display: inline-flex;
  align-items: center;
}

.manual-input-section label {
  font-size: 0.85rem;
  font-weight: 700;
  color: #334155;
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
  border-radius: 10px;
  font-size: 0.9rem;
}

.input-group button {
  background: #2563eb;
  color: white;
  border: none;
  padding: 0 16px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
}

.error-banner {
  background: #fef2f2;
  color: #991b1b;
  padding: 12px;
  border-radius: 10px;
  font-size: 0.85rem;
  font-weight: 600;
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
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.25);
}
</style>
