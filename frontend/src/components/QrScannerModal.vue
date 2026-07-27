<template>
  <ModalDialog :show="show" title="Pemindai QR Code Aset" maxWidth="540px" @close="handleClose">
    <div class="qr-scanner-content">
      <!-- Camera Live Stream & Viewfinder Container -->
      <div class="camera-simulator">
        <div id="qr-reader" class="html5-qr-wrapper" v-show="isCameraActive"></div>

        <div v-if="!isCameraActive" class="scanner-frame-placeholder" :class="{ 'scan-success': scanSuccessAnim }">
          <div class="corner top-left"></div>
          <div class="corner top-right"></div>
          <div class="corner bottom-left"></div>
          <div class="corner bottom-right"></div>
          <div class="placeholder-center">
            <svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="#64748b" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"><path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z"/><circle cx="12" cy="13" r="3"/></svg>
            <p class="placeholder-text">Arahkan Kamera ke Stiker QR Code</p>
          </div>
        </div>

        <p class="scan-instruction">
          {{ scanStatusText }}
        </p>

        <div class="camera-controls">
          <button v-if="!isCameraActive" class="cam-btn start-cam" @click="startCamera">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="cam-icon"><path d="m22 8-6 4 6 4V8z"/><rect width="14" height="12" x="2" y="6" rx="2" ry="2"/></svg>
            <span>Buka Kamera</span>
          </button>
          <button v-else class="cam-btn stop-cam" @click="stopCamera">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="cam-icon"><rect width="18" height="18" x="3" y="3" rx="2"/></svg>
            <span>Matikan Kamera</span>
          </button>
          <label class="cam-btn file-btn">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="cam-icon"><path d="M4 20h16a2 2 0 0 0 2-2V8a2 2 0 0 0-2-2h-7.93a2 2 0 0 1-1.66-.9l-.82-1.2A2 2 0 0 0 7.93 3H4a2 2 0 0 0-2 2v13c0 1.1.9 2 2 2z"/></svg>
            <span>Upload Foto QR</span>
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
  if (scanSuccessAnim.value) return 'QR Code Terdeteksi & Diproses!'
  if (isCameraActive.value) return 'Kamera Aktif — Arahkan ke Stiker QR Code Aset'
  return 'Pilih metode pemindaian di bawah'
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

.qr-scanner-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.camera-simulator {
  background: #f8fafc;
  border-radius: 6px !important;
  border: 1px solid #e2e8f0;
  padding: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #0f172a;
  position: relative;
  overflow: hidden;
  min-height: 240px;
}

.html5-qr-wrapper {
  width: 100%;
  max-width: 380px;
  border-radius: 4px !important;
  border: 1px solid #cbd5e1;
  overflow: hidden;
  margin-bottom: 10px;
}

.scanner-frame-placeholder {
  width: 180px;
  height: 180px;
  border: 2px dashed #94a3b8;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12px;
  border-radius: 6px !important;
  padding: 12px;
  background: #ffffff;
}

.placeholder-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.placeholder-text {
  font-size: 0.78rem;
  color: #64748b;
  text-align: center;
  font-weight: 700;
  margin: 0;
}

.scanner-frame-placeholder.scan-success {
  border-color: #16a34a;
  box-shadow: 0 0 20px rgba(22, 163, 74, 0.3);
}

.corner {
  position: absolute;
  width: 16px;
  height: 16px;
  border-color: #0f172a;
  border-style: solid;
}

.top-left { top: -2px; left: -2px; border-width: 3px 0 0 3px; }
.top-right { top: -2px; right: -2px; border-width: 3px 3px 0 0; }
.bottom-left { bottom: -2px; left: -2px; border-width: 0 0 3px 3px; }
.bottom-right { bottom: -2px; right: -2px; border-width: 0 3px 3px 0; }

.scan-instruction {
  margin: 8px 0 12px;
  font-size: 0.84rem;
  color: #0f172a;
  text-align: center;
  font-weight: 700;
}

.camera-controls {
  margin-top: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  width: 100%;
  max-width: 420px;
}

.cam-btn {
  flex: 1 1 0%;
  width: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  height: 42px;
  padding: 0 12px;
  border-radius: 4px !important;
  font-size: 0.85rem;
  font-weight: 700;
  cursor: pointer;
  box-sizing: border-box;
  text-align: center;
  vertical-align: middle;
  line-height: 1;
  user-select: none;
  margin: 0;
  outline: none;
  transition: all 0.15s ease;
}

.cam-btn .cam-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  display: block;
}

.cam-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.cam-btn.start-cam {
  background: #0f172a;
  color: #ffffff;
  border: 1px solid #0f172a;
}

.cam-btn.start-cam:hover {
  background: #1e293b;
}

.cam-btn.stop-cam {
  background: #dc2626;
  color: #ffffff;
  border: 1px solid #dc2626;
}

.cam-btn.stop-cam:hover {
  background: #b91c1c;
}

.cam-btn.file-btn {
  background: #2563eb;
  color: #ffffff;
  border: 1px solid #2563eb;
}

.cam-btn.file-btn:hover {
  background: #1d4ed8;
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
  border-radius: 4px !important;
  font-size: 0.88rem;
  color: #0f172a;
  outline: none;
}

.input-group input:focus {
  border-color: #2563eb;
}

.input-group button {
  background: #0f172a;
  color: white;
  border: 1px solid #0f172a;
  padding: 0 16px;
  border-radius: 4px !important;
  font-weight: 700;
  cursor: pointer;
}

.error-banner {
  background: #fef2f2;
  color: #991b1b;
  border: 1px solid #fca5a5;
  padding: 10px 14px;
  border-radius: 4px !important;
  font-size: 0.85rem;
  font-weight: 600;
}

.scanned-result-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 4px !important;
  padding: 16px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
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
  font-weight: 800;
}

.asset-code-badge {
  font-size: 0.75rem;
  background: #f1f5f9;
  color: #0f172a;
  border: 1px solid #cbd5e1;
  padding: 2px 8px;
  border-radius: 4px !important;
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
  border: 1px solid #dc2626;
  padding: 10px;
  border-radius: 4px !important;
  font-weight: 700;
  cursor: pointer;
}
</style>
