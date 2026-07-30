<template>
  <div class="view-wrapper">
    <div class="page-container">
      <div class="page-header">
        <div>
          <p class="eyebrow">Audit Trail</p>
          <h1>Activity Log</h1>
          <p class="subtitle">Riwayat perbaikan & maintenance.</p>
        </div>

        <div class="header-action-group">
          <button class="primary-btn" @click="showReportModal = true" title="Prinjau & Export Laporan Audit Trail">
            <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
            <span>Laporan & Export</span>
          </button>
        </div>
      </div>

      <!-- Summary boxes -->
      <div class="summary-row">
        <div class="sbox green">
          <span class="sbox-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#16a34a" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
          </span>
          <div>
            <p class="sbox-label">WO Finish</p>
            <p class="sbox-value">{{ finishedWOs.length }} Tiket</p>
          </div>
        </div>
        <div class="sbox blue">
          <span class="sbox-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
          </span>
          <div>
            <p class="sbox-label">Maintenance Selesai</p>
            <p class="sbox-value">{{ maintenanceHistory.length }} Riwayat</p>
          </div>
        </div>
        <div class="sbox purple">
          <span class="sbox-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#9333ea" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"/><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"/></svg>
          </span>
          <div>
            <p class="sbox-label">Total Biaya Maintenance</p>
            <p class="sbox-value">Rp {{ formatNumber(totalMaintenanceCost) }}</p>
          </div>
        </div>
        <div class="sbox orange">
          <span class="sbox-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#ea580c" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 3h5v5"/><path d="M4 20L21 3"/><path d="M21 16v5h-5"/><path d="M15 15l6 6"/><path d="M4 4l5 5"/></svg>
          </span>
          <div>
            <p class="sbox-label">Total Mutasi Aset</p>
            <p class="sbox-value">{{ mutations.length }} Mutasi</p>
          </div>
        </div>
      </div>

      <!-- Search -->
      <div class="search-row">
        <input v-model="searchFilter" placeholder="Cari riwayat activity log..." class="search-input" />
      </div>

      <!-- Section 1: Work Order Timelines (Audit Trail) -->
      <div class="card-panel">
        <h2 class="section-title">Riwayat Seluruh Timeline Work Order (Audit Trail)</h2>
        <div class="table-responsive">
          <table>
            <thead>
              <tr>
                <th>Timeline ID</th>
                <th>WO ID</th>
                <th>Status</th>
                <th>Di-update Oleh</th>
                <th>Tindakan / Catatan Progres</th>
                <th>Biaya (Rp)</th>
                <th>Tanggal & Waktu Update</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="tl in filteredTimelines" :key="tl.id">
                <td><span class="wo-id">#TL-{{ tl.id }}</span></td>
                <td><span class="wo-id">#WO-{{ tl.work_order_id }}</span></td>
                <td><StatusBadge :status="tl.status" /></td>
                <td>
                  <span class="requester-chip">@{{ tl.updated_by || 'Sistem' }}</span>
                  <span class="user-role-sub-inline" v-if="tl.user_role"> ({{ tl.user_role }})</span>
                </td>
                <td class="desc-cell" :title="tl.action_taken">{{ tl.action_taken || '—' }}</td>
                <td>Rp {{ formatNumber(tl.cost || 0) }}</td>
                <td class="time-col">{{ formatDate(tl.created_at) }}</td>
              </tr>
              <tr v-if="filteredTimelines.length === 0">
                <td colspan="7" class="empty-state">Tidak ada data riwayat timeline Work Order.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Section 2: Maintenance History -->
      <div class="card-panel" style="margin-top: 24px;">
        <h2 class="section-title">Maintenance Selesai</h2>
        <div class="table-responsive">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Aset ID</th>
                <th>Tindakan Perawatan</th>
                <th>Biaya (Rp)</th>
                <th>Tanggal Pengerjaan</th>
                <th v-if="canManage">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mh in filteredMH" :key="mh.id">
                <td><span class="wo-id">#MH-{{ mh.id }}</span></td>
                <td>Aset #{{ mh.asset_id }}</td>
                <td class="desc-cell" :title="mh.action_taken">{{ mh.action_taken }}</td>
                <td>Rp {{ formatNumber(mh.cost || 0) }}</td>
                <td class="time-col">{{ formatDate(mh.created_at) }}</td>
                <td v-if="canManage" class="actions-cell">
                  <button class="icon-btn edit-btn" @click="openEditMhModal(mh)" title="Edit">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                    <span>Edit</span>
                  </button>
                  <button class="icon-btn delete-btn" @click="deleteMh(mh)" title="Hapus">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                    <span>Hapus</span>
                  </button>
                </td>
              </tr>
              <tr v-if="filteredMH.length === 0">
                <td :colspan="canManage ? 6 : 5" class="empty-state">Tidak ada data.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Section 3: Asset Mutation Timelines -->
      <div class="card-panel" style="margin-top: 24px;">
        <h2 class="section-title">Riwayat Seluruh Mutasi Aset (Audit Trail)</h2>
        <div class="table-responsive">
          <table>
            <thead>
              <tr>
                <th>Timeline ID</th>
                <th>Kode Aset</th>
                <th>Lokasi Asal</th>
                <th>Lokasi Tujuan (Baru)</th>
                <th>PIC / Penanggung Jawab</th>
                <th>Alasan Mutasi</th>
                <th>Tanggal & Waktu Mutasi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="mut in filteredAssetMutationTimelines" :key="mut.id">
                <td><span class="wo-id">#AMUT-{{ mut.id }}</span></td>
                <td><span class="wo-id">{{ mut.asset_code }}</span></td>
                <td>{{ mut.previous_location || '—' }}</td>
                <td><span class="location-new">{{ mut.new_location || '—' }}</span></td>
                <td>{{ mut.pic || 'Engineering' }}</td>
                <td class="desc-cell" :title="mut.reason">{{ mut.reason || '—' }}</td>
                <td class="time-col">{{ formatDate(mut.moved_at) }}</td>
              </tr>
              <tr v-if="filteredAssetMutationTimelines.length === 0">
                <td colspan="7" class="empty-state">Belum ada riwayat mutasi aset.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Custom UI Toast Notification -->
      <transition name="fade">
        <div v-if="showToast" :class="['custom-ui-toast', toastType]">
          <span class="toast-icon">{{ toastType === 'success' ? '✅' : '⚠️' }}</span>
          <span class="toast-text">{{ toastMsg }}</span>
          <button class="toast-close" @click="showToast = false">✕</button>
        </div>
      </transition>
    </div>

    <!-- Modals placed OUTSIDE .page-container so print media query hiding .page-container leaves report modal visible -->
    <ModalDialog :show="showWoModal" title="Edit Work Order" @close="showWoModal = false">
      <form @submit.prevent="submitEditWo" class="modal-form" v-if="selectedWo">
        <p class="modal-info"><strong>WO #{{ selectedWo.id }}:</strong> {{ selectedWo.location }}</p>
        <label>
          <span>Deskripsi</span>
          <textarea v-model="editWoDesc" rows="3" required></textarea>
        </label>
        <label>
          <span>Tindakan Perbaikan</span>
          <textarea v-model="editWoAction" rows="3" required></textarea>
        </label>
        <label>
          <span>Biaya (Rp)</span>
          <input v-model.number="editWoCost" type="number" min="0" required />
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Perubahan</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showMhModal" title="Edit Maintenance" @close="showMhModal = false">
      <form @submit.prevent="submitEditMh" class="modal-form" v-if="selectedMh">
        <p class="modal-info"><strong>Maintenance #MH-{{ selectedMh.id }} (Aset #{{ selectedMh.asset_id }})</strong></p>
        <label>
          <span>Tindakan Perawatan</span>
          <textarea v-model="editMhAction" rows="3" required></textarea>
        </label>
        <label>
          <span>Biaya (Rp)</span>
          <input v-model.number="editMhCost" type="number" min="0" required />
        </label>

        <button type="submit" class="submit-modal-btn">Simpan Perubahan</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showReportModal" title="Laporan Activity Log & Audit Trail" maxWidth="960px" @close="showReportModal = false">
      <!-- Report Filter Selection Bar (iPhone / iOS Segmented Control) -->
      <div class="report-filter-bar no-print">
        <span class="filter-bar-title">Tipe Dokumen Laporan</span>
        <div class="filter-btn-group">
          <button 
            :class="['filter-tab-btn', { active: reportTypeFilter === 'all' }]" 
            @click="reportTypeFilter = 'all'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="7" height="7" x="3" y="3" rx="1"/><rect width="7" height="7" x="14" y="3" rx="1"/><rect width="7" height="7" x="14" y="14" rx="1"/><rect width="7" height="7" x="3" y="14" rx="1"/></svg>
            <span>Semua</span>
          </button>
          <button 
            :class="['filter-tab-btn', { active: reportTypeFilter === 'wo' }]" 
            @click="reportTypeFilter = 'wo'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
            <span>Work Order</span>
          </button>
          <button 
            :class="['filter-tab-btn', { active: reportTypeFilter === 'maintenance' }]" 
            @click="reportTypeFilter = 'maintenance'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M22 12h-4l-3 9L9 3l-3 9H2"/></svg>
            <span>Maintenance</span>
          </button>
          <button 
            :class="['filter-tab-btn', { active: reportTypeFilter === 'mutation' }]" 
            @click="reportTypeFilter = 'mutation'"
          >
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M17 1l4 4-4 4"/><path d="M3 11V9a4 4 0 0 1 4-4h14"/><path d="M7 23l-4-4 4-4"/><path d="M21 13v2a4 4 0 0 1-4 4H3"/></svg>
            <span>Mutasi Aset</span>
          </button>
        </div>
      </div>

      <div class="monthly-report-printable" id="printableReportDocument">
        <div class="report-header">
          <h2>
            {{ 
              reportTypeFilter === 'wo' ? 'LAPORAN REKAPITULASI WORK ORDER' : 
              reportTypeFilter === 'maintenance' ? 'LAPORAN REKAPITULASI MAINTENANCE' : 
              reportTypeFilter === 'mutation' ? 'LAPORAN RIWAYAT MUTASI ASET' : 
              'LAPORAN AUDIT TRAIL & AKTIVITAS OPERASIONAL' 
            }}
          </h2>
          <p class="report-sub">Sistem AsetKu — Periode: {{ reportMonthYear }}</p>
          <hr class="report-divider" />
        </div>

        <div class="report-summary-boxes">
          <div class="rbox success" v-if="reportTypeFilter === 'all' || reportTypeFilter === 'wo'">
            <span>WO Finish</span>
            <strong>{{ finishedWOs.length }} Tiket</strong>
          </div>
          <div class="rbox blue" v-if="reportTypeFilter === 'all' || reportTypeFilter === 'maintenance'">
            <span>Maintenance Selesai</span>
            <strong>{{ maintenanceHistory.length }} Riwayat</strong>
          </div>
          <div class="rbox orange" v-if="reportTypeFilter === 'all' || reportTypeFilter === 'mutation'">
            <span>Total Mutasi Aset</span>
            <strong>{{ assetMutationTimelines.length || mutations.length }} Mutasi</strong>
          </div>
          <div class="rbox danger" v-if="reportTypeFilter === 'all' || reportTypeFilter === 'maintenance' || reportTypeFilter === 'wo'">
            <span>Total Biaya (WO / Maintenance)</span>
            <strong>Rp {{ formatNumber(reportTypeFilter === 'wo' ? totalWoCost : totalMaintenanceCost) }}</strong>
          </div>
        </div>

        <!-- Section 1: Work Order Selesai -->
        <template v-if="reportTypeFilter === 'all' || reportTypeFilter === 'wo'">
          <h3 class="report-section-heading" v-if="reportTypeFilter === 'all'">1. Rekapitulasi Work Order Selesai</h3>
          <div class="report-table-wrapper">
            <table class="report-table">
              <thead>
                <tr>
                  <th class="col-id">ID</th>
                  <th class="col-loc">Lokasi</th>
                  <th class="col-prio">Prioritas</th>
                  <th class="col-desc">Deskripsi</th>
                  <th class="col-desc">Tindakan Perbaikan</th>
                  <th class="col-stat">Status</th>
                  <th class="col-cost">Biaya (Rp)</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="wo in finishedWOs" :key="wo.id">
                  <td class="col-id">#WO-{{ wo.id }}</td>
                  <td class="col-loc">{{ wo.location || '—' }}</td>
                  <td class="col-prio">{{ wo.priority }}</td>
                  <td class="col-desc">{{ wo.description }}</td>
                  <td class="col-desc">{{ wo.action_taken || '—' }}</td>
                  <td class="col-stat">{{ wo.status }}</td>
                  <td class="col-cost">Rp {{ formatNumber(wo.cost || 0) }}</td>
                </tr>
                <tr v-if="finishedWOs.length === 0">
                  <td colspan="7" class="empty-state">Tidak ada data.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>

        <!-- Section 2: Maintenance Selesai -->
        <template v-if="reportTypeFilter === 'all' || reportTypeFilter === 'maintenance'">
          <h3 class="report-section-heading" :style="{ marginTop: reportTypeFilter === 'all' ? '24px' : '0' }">
            {{ reportTypeFilter === 'all' ? '2. Rekapitulasi Maintenance Selesai' : 'Rekapitulasi Maintenance Selesai' }}
          </h3>
          <div class="report-table-wrapper">
            <table class="report-table">
              <thead>
                <tr>
                  <th class="col-id">ID</th>
                  <th class="col-prio">Aset ID</th>
                  <th class="col-desc">Tindakan Perawatan</th>
                  <th class="col-cost">Biaya (Rp)</th>
                  <th class="col-loc">Tanggal</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="mh in maintenanceHistory" :key="mh.id">
                  <td class="col-id">#MH-{{ mh.id }}</td>
                  <td class="col-prio">Aset #{{ mh.asset_id }}</td>
                  <td class="col-desc">{{ mh.action_taken }}</td>
                  <td class="col-cost">Rp {{ formatNumber(mh.cost || 0) }}</td>
                  <td class="col-loc">{{ formatDate(mh.created_at) }}</td>
                </tr>
                <tr v-if="maintenanceHistory.length === 0">
                  <td colspan="5" class="empty-state">Tidak ada data.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>

        <!-- Section 3: Mutasi Aset -->
        <template v-if="reportTypeFilter === 'all' || reportTypeFilter === 'mutation'">
          <h3 class="report-section-heading" :style="{ marginTop: reportTypeFilter === 'all' ? '24px' : '0' }">
            {{ reportTypeFilter === 'all' ? '3. Rekapitulasi Mutasi Aset (Audit Trail Mutasi)' : 'Rekapitulasi Mutasi Aset (Audit Trail Mutasi)' }}
          </h3>
          <div class="report-table-wrapper">
            <table class="report-table">
              <thead>
                <tr>
                  <th class="col-id">ID Mutasi</th>
                  <th class="col-prio">Kode Aset</th>
                  <th class="col-loc">Lokasi Asal</th>
                  <th class="col-loc">Lokasi Baru</th>
                  <th class="col-stat">PIC</th>
                  <th class="col-desc">Alasan Mutasi</th>
                  <th class="col-loc">Waktu Mutasi</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="mut in (assetMutationTimelines.length ? assetMutationTimelines : mutations)" :key="mut.id">
                  <td class="col-id">#AMUT-{{ mut.id }}</td>
                  <td class="col-prio">{{ mut.asset_code }}</td>
                  <td class="col-loc">{{ mut.previous_location || '—' }}</td>
                  <td class="col-loc">{{ mut.new_location || '—' }}</td>
                  <td class="col-stat">{{ mut.pic || 'Engineering' }}</td>
                  <td class="col-desc">{{ mut.reason || '—' }}</td>
                  <td class="col-loc">{{ formatDate(mut.moved_at || mut.created_at) }}</td>
                </tr>
                <tr v-if="(assetMutationTimelines.length === 0 && mutations.length === 0)">
                  <td colspan="7" class="empty-state">Belum ada riwayat mutasi aset.</td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
      </div>

      <div class="report-actions no-print">
        <button class="excel-btn" @click="exportToExcel">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="8" y1="13" x2="16" y2="13"/><line x1="8" y1="17" x2="16" y2="17"/></svg>
          <span>Export ke Excel (.xlsx)</span>
        </button>
        <button class="print-btn" @click="printReport">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
          <span>Cetak Dokumen Laporan (PDF / Print)</span>
        </button>
      </div>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'
import { triggerPrint } from '../utils/exportUtils'

const userRole = ref(sessionStorage.getItem('user_role') || 'external')
const canManage = computed(() => userRole.value === 'hod' || userRole.value === 'management' || userRole.value === 'admin')

const showToast = ref(false)
const toastMsg = ref('')
const toastType = ref('success')

function notify(msg, type = 'success') {
  toastMsg.value = msg
  toastType.value = type
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 4000)
}

const searchFilter = ref('')
const finishedWOs = ref([])
const timelines = ref([])
const assetMutationTimelines = ref([])
const maintenanceHistory = ref([])
const mutations = ref([])
const isLoading = ref(false)

const showReportModal = ref(false)
const reportTypeFilter = ref('all') // 'all' | 'wo' | 'maintenance' | 'mutation'
const reportMonthYear = ref(new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }))

const showWoModal = ref(false)
const selectedWo = ref(null)
const editWoDesc = ref('')
const editWoAction = ref('')
const editWoCost = ref(0)

const showMhModal = ref(false)
const selectedMh = ref(null)
const editMhAction = ref('')
const editMhCost = ref(0)

const totalMaintenanceCost = computed(() =>
  maintenanceHistory.value.reduce((sum, mh) => sum + (mh.cost || 0), 0)
)

const totalWoCost = computed(() =>
  finishedWOs.value.reduce((sum, wo) => sum + (wo.cost || 0), 0)
)

const filteredTimelines = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return timelines.value
  return timelines.value.filter(tl =>
    (tl.action_taken && tl.action_taken.toLowerCase().includes(q)) ||
    (tl.updated_by && tl.updated_by.toLowerCase().includes(q)) ||
    (tl.status && tl.status.toLowerCase().includes(q)) ||
    String(tl.work_order_id).includes(q)
  )
})

const filteredAssetMutationTimelines = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return assetMutationTimelines.value
  return assetMutationTimelines.value.filter(tl =>
    (tl.asset_code && tl.asset_code.toLowerCase().includes(q)) ||
    (tl.previous_location && tl.previous_location.toLowerCase().includes(q)) ||
    (tl.new_location && tl.new_location.toLowerCase().includes(q)) ||
    (tl.pic && tl.pic.toLowerCase().includes(q)) ||
    (tl.reason && tl.reason.toLowerCase().includes(q))
  )
})

const filteredWOs = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return finishedWOs.value
  return finishedWOs.value.filter(wo =>
    (wo.description && wo.description.toLowerCase().includes(q)) ||
    (wo.action_taken && wo.action_taken.toLowerCase().includes(q)) ||
    (wo.location && wo.location.toLowerCase().includes(q)) ||
    String(wo.asset_id).includes(q)
  )
})

const filteredMH = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return maintenanceHistory.value
  return maintenanceHistory.value.filter(mh =>
    (mh.action_taken && mh.action_taken.toLowerCase().includes(q)) ||
    String(mh.asset_id).includes(q)
  )
})

const filteredMutations = computed(() => {
  const q = searchFilter.value.toLowerCase()
  if (!q) return mutations.value
  return mutations.value.filter(mut =>
    (mut.new_location && mut.new_location.toLowerCase().includes(q)) ||
    (mut.previous_location && mut.previous_location.toLowerCase().includes(q)) ||
    (mut.new_pic && mut.new_pic.toLowerCase().includes(q)) ||
    (mut.reason && mut.reason.toLowerCase().includes(q)) ||
    String(mut.asset_id).includes(q)
  )
})

function formatNumber(num) {
  return (num || 0).toLocaleString('id-ID')
}

function formatDate(dateStr) {
  if (!dateStr) return '—'
  try {
    return new Date(dateStr).toLocaleString('id-ID', {
      day: '2-digit', month: 'short', year: 'numeric', hour: '2-digit', minute: '2-digit'
    })
  } catch {
    return dateStr
  }
}

async function fetchLogs() {
  isLoading.value = true
  try {
    const res = await api.get('/activitylogs')
    if (res.data?.data) {
      finishedWOs.value = res.data.data.work_orders || []
      timelines.value = res.data.data.timelines || []
      assetMutationTimelines.value = res.data.data.asset_mutation_timelines || []
      maintenanceHistory.value = res.data.data.maintenance_history || []
      mutations.value = res.data.data.mutations || []
    }
  } catch (e) {
    console.error('Failed to fetch activity logs:', e)
  } finally {
    isLoading.value = false
  }
}

function openEditWoModal(wo) {
  selectedWo.value = wo
  editWoDesc.value = wo.description || ''
  editWoAction.value = wo.action_taken || ''
  editWoCost.value = wo.cost || 0
  showWoModal.value = true
}

async function submitEditWo() {
  if (!selectedWo.value) return
  try {
    await api.post('/workorders/edit', {
      wo_id: selectedWo.value.id,
      description: editWoDesc.value,
      action_taken: editWoAction.value,
      cost: editWoCost.value
    })
    showWoModal.value = false
    notify('WO diperbarui!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal mengubah WO: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteWo(wo) {
  try {
    await api.post('/workorders/delete', { wo_id: wo.id })
    finishedWOs.value = finishedWOs.value.filter(w => w.id !== wo.id)
    notify('WO dihapus!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal menghapus WO: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openEditMhModal(mh) {
  selectedMh.value = mh
  editMhAction.value = mh.action_taken || ''
  editMhCost.value = mh.cost || 0
  showMhModal.value = true
}

async function submitEditMh() {
  if (!selectedMh.value) return
  try {
    await api.post('/maintenance/history/edit', {
      history_id: selectedMh.value.id,
      action_taken: editMhAction.value,
      cost: editMhCost.value
    })
    showMhModal.value = false
    notify('Riwayat maintenance diperbarui!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal mengubah riwayat: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteMh(mh) {
  try {
    await api.post('/maintenance/history/delete', { history_id: mh.id })
    maintenanceHistory.value = maintenanceHistory.value.filter(m => m.id !== mh.id)
    notify('Riwayat maintenance dihapus!', 'success')
    await fetchLogs()
  } catch (e) {
    notify('Gagal menghapus riwayat: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function printReport() {
  triggerPrint()
}

function exportToExcel() {
  const monthName = reportMonthYear.value.replace(/\s+/g, '_')
  const typeTag = reportTypeFilter.value.toUpperCase()
  const fileName = `Laporan_ActivityLog_${typeTag}_${monthName}.xls`

  let htmlTable = `
    <html xmlns:o="urn:schemas-microsoft-com:office:office" xmlns:x="urn:schemas-microsoft-com:office:excel" xmlns="http://www.w3.org/TR/REC-html40">
    <head>
      <meta charset="utf-8">
      <!--[if gte mso 9]>
      <xml>
        <x:ExcelWorkbook>
          <x:ExcelWorksheets>
            <x:ExcelWorksheet>
              <x:Name>Activity Log</x:Name>
              <x:WorksheetOptions><x:DisplayGridlines/></x:WorksheetOptions>
            </x:ExcelWorksheet>
          </x:ExcelWorksheets>
        </x:ExcelWorkbook>
      </xml>
      <![endif]-->
      <style>
        th { background-color: #2563eb; color: #ffffff; font-weight: bold; border: 1px solid #000000; text-align: center; }
        td { border: 1px solid #cccccc; vertical-align: middle; }
        .summary-header { font-size: 14pt; font-weight: bold; color: #0f172a; }
        .section-header { font-size: 12pt; font-weight: bold; color: #1e293b; background-color: #e2e8f0; }
        .total-row { font-weight: bold; background-color: #f1f5f9; }
      </style>
    </head>
    <body>
      <h2 class="summary-header">LAPORAN BULANAN AUDIT TRAIL & AKTIVITAS OPERASIONAL</h2>
      <p><b>Sistem AsetKu</b> — Periode: ${reportMonthYear.value}</p>
      <br/>
  `

  if (reportTypeFilter.value === 'all' || reportTypeFilter.value === 'wo') {
    htmlTable += `
      <h3 class="section-header">1. REKAPITULASI WORK ORDER SELESAI (FINISH)</h3>
      <table border="1" cellspacing="0" cellpadding="6">
        <thead>
          <tr>
            <th>ID Work Order</th>
            <th>Lokasi / Kamar</th>
            <th>Prioritas</th>
            <th>Deskripsi Kerusakan</th>
            <th>Tindakan Perbaikan</th>
            <th>Status Tiket</th>
            <th>Biaya Perbaikan (Rp)</th>
          </tr>
        </thead>
        <tbody>
    `
    finishedWOs.value.forEach(wo => {
      htmlTable += `
        <tr>
          <td>#WO-${wo.id}</td>
          <td>${wo.location || ''}</td>
          <td>${wo.priority || ''}</td>
          <td>${wo.description || ''}</td>
          <td>${wo.action_taken || ''}</td>
          <td>${wo.status || ''}</td>
          <td align="right">${wo.cost || 0}</td>
        </tr>
      `
    })
    htmlTable += `
          <tr class="total-row">
            <td colspan="6" align="right"><b>TOTAL BIAYA WO SELESAI:</b></td>
            <td align="right"><b>${totalWoCost.value}</b></td>
          </tr>
        </tbody>
      </table>
      <br/><br/>
    `
  }

  if (reportTypeFilter.value === 'all' || reportTypeFilter.value === 'maintenance') {
    htmlTable += `
      <h3 class="section-header">2. REKAPITULASI RIWAYAT MAINTENANCE SELESAI</h3>
      <table border="1" cellspacing="0" cellpadding="6">
        <thead>
          <tr>
            <th>ID History</th>
            <th>ID Aset</th>
            <th>Tindakan Perawatan yang Dilakukan</th>
            <th>Biaya Maintenance (Rp)</th>
            <th>Tanggal Pengerjaan</th>
          </tr>
        </thead>
        <tbody>
    `
    maintenanceHistory.value.forEach(mh => {
      htmlTable += `
        <tr>
          <td>#MH-${mh.id}</td>
          <td>Aset #${mh.asset_id}</td>
          <td>${mh.action_taken || ''}</td>
          <td align="right">${mh.cost || 0}</td>
          <td>${formatDate(mh.created_at)}</td>
        </tr>
      `
    })
    htmlTable += `
          <tr class="total-row">
            <td colspan="3" align="right"><b>TOTAL BIAYA MAINTENANCE:</b></td>
            <td align="right"><b>${totalMaintenanceCost.value}</b></td>
            <td></td>
          </tr>
        </tbody>
      </table>
      <br/><br/>
    `
  }

  if (reportTypeFilter.value === 'all' || reportTypeFilter.value === 'mutation') {
    htmlTable += `
      <h3 class="section-header">3. REKAPITULASI RIWAYAT MUTASI ASET (AUDIT TRAIL MUTASI)</h3>
      <table border="1" cellspacing="0" cellpadding="6">
        <thead>
          <tr>
            <th>ID Mutasi</th>
            <th>Kode Aset</th>
            <th>Lokasi Asal</th>
            <th>Lokasi Baru</th>
            <th>PIC Penanggung Jawab</th>
            <th>Alasan Mutasi</th>
            <th>Waktu Mutasi</th>
          </tr>
        </thead>
        <tbody>
    `
    const mutList = assetMutationTimelines.value.length ? assetMutationTimelines.value : mutations.value
    mutList.forEach(mut => {
      htmlTable += `
        <tr>
          <td>#AMUT-${mut.id}</td>
          <td>${mut.asset_code || ''}</td>
          <td>${mut.previous_location || '—'}</td>
          <td>${mut.new_location || '—'}</td>
          <td>${mut.pic || 'Engineering'}</td>
          <td>${mut.reason || '—'}</td>
          <td>${formatDate(mut.moved_at || mut.created_at)}</td>
        </tr>
      `
    })
    htmlTable += `
        </tbody>
      </table>
    `
  }

  htmlTable += `
    </body>
    </html>
  `

  const blob = new Blob([htmlTable], { type: 'application/vnd.ms-excel;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', fileName)
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

onMounted(() => {
  fetchLogs()
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

.eyebrow {
  margin: 0 0 4px;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  font-size: 0.8rem;
  color: #2563eb;
  font-weight: 700;
}

h1, .page-header h1 {
  margin: 0 0 8px;
  font-size: 1.8rem;
  color: #0f172a;
  font-weight: 800;
}

.subtitle {
  margin: 0;
  color: #64748b;
  font-size: 0.95rem;
}

.header-action-group {
  display: flex;
  gap: 12px;
}

.print-report-btn {
  background: #0284c7;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(2, 132, 199, 0.3);
}

.export-excel-btn {
  background: #16a34a;
  color: white;
  border: none;
  padding: 12px 20px;
  border-radius: 12px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 4px 14px rgba(22, 163, 74, 0.3);
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

.summary-row {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(220px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.sbox {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 20px;
  border-radius: 16px;
  background: #fff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 12px rgba(0,0,0,0.04);
}

.sbox.green { border-left: 4px solid #16a34a; }
.sbox.blue  { border-left: 4px solid #2563eb; }
.sbox.purple { border-left: 4px solid #7c3aed; }
.sbox.orange { border-left: 4px solid #d97706; }

.location-new {
  color: #16a34a;
  font-weight: 700;
}

.sbox-icon {
  font-size: 1.8rem;
  line-height: 1;
}

.sbox-label {
  margin: 0 0 4px;
  font-size: 0.8rem;
  color: #64748b;
  font-weight: 600;
}

.sbox-value {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

.search-row {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
  max-width: 480px;
  padding: 11px 18px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  font-size: 0.9rem;
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.section-title {
  margin: 0 0 20px;
  font-size: 1.1rem;
  color: #0f172a;
  font-weight: 800;
}

.table-responsive {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 700px;
}

th {
  text-align: left;
  padding: 12px 14px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.82rem;
  font-weight: 700;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

td {
  padding: 13px 14px;
  border-bottom: 1px solid #f1f5f9;
  color: #334155;
  font-size: 0.88rem;
  vertical-align: middle;
}

.wo-id {
  font-family: monospace;
  font-weight: 800;
  color: #0f172a;
  font-size: 0.85rem;
}

.desc-cell {
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.time-col {
  color: #64748b;
  font-size: 0.82rem;
  white-space: nowrap;
}

.actions-cell {
  display: flex;
  gap: 6px;
  white-space: nowrap;
}

.icon-btn {
  border: 1px solid #cbd5e1;
  background: white;
  padding: 5px 10px;
  border-radius: 8px;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
}

.edit-btn { color: #2563eb; border-color: #93c5fd; }
.delete-btn { color: #dc2626; border-color: #fca5a5; }

.modal-form {
  display: grid;
  gap: 14px;
}

.modal-form label {
  display: grid;
  gap: 6px;
  font-weight: 600;
  color: #1e293b;
  font-size: 0.9rem;
}

.modal-form input, .modal-form textarea {
  width: 100%;
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  font-size: 0.9rem;
}

.submit-modal-btn {
  background: #2563eb;
  color: white;
  border: none;
  padding: 12px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  margin-top: 6px;
}

.modal-info {
  background: #f8fafc;
  padding: 10px 14px;
  border-radius: 10px;
  font-size: 0.9rem;
  color: #334155;
}

.empty-state {
  text-align: center;
  padding: 32px;
  color: #94a3b8;
  font-size: 0.9rem;
}

/* Printable Monthly Activity Log Styling */
.monthly-report-printable {
  background: #ffffff;
  padding: 20px;
  color: #0f172a;
}

.report-header {
  text-align: center;
  margin-bottom: 20px;
}

.report-header h2 {
  margin: 0;
  font-size: 1.25rem;
  color: #0f172a;
}

.report-sub {
  margin: 4px 0 0;
  color: #64748b;
  font-size: 0.9rem;
}

.report-divider {
  border: 0;
  border-top: 2px solid #0f172a;
  margin: 12px 0 16px;
}

.report-summary-boxes {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.rbox {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 10px;
  padding: 10px 14px;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.rbox span { font-size: 0.75rem; color: #475569; font-weight: 600; }
.rbox strong { font-size: 1.1rem; color: #0f172a; font-weight: 800; }
.rbox.success { border-color: #86efac; background: #f0fdf4; }
.rbox.blue { border-color: #93c5fd; background: #eff6ff; }
.rbox.danger { border-color: #fca5a5; background: #fef2f2; }

.report-section-heading {
  font-size: 1rem;
  font-weight: 800;
  color: #0f172a;
  margin: 0 0 10px;
}

.report-table-wrapper {
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 0.85rem;
}

.report-table th {
  background: #f1f5f9;
  color: #0f172a;
  font-weight: 700;
  border: 1px solid #94a3b8;
  padding: 8px 10px;
  text-align: left;
}

.report-table td {
  border: 1px solid #cbd5e1;
  padding: 8px 10px;
  color: #334155;
  vertical-align: middle;
}

.primary-btn {
  background: #ffffff !important;
  color: #0f172a !important;
  border: 1px solid #cbd5e1 !important;
  padding: 10px 18px !important;
  border-radius: 10px !important;
  font-size: 0.88rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
  transition: all 0.15s ease;
  line-height: 1;
  white-space: nowrap;
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.primary-btn .btn-icon {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  display: block;
}

.print-report-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.actions-cell {
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 5px;
  height: 32px;
  padding: 0 10px;
  border-radius: 6px !important;
  font-size: 0.78rem;
  font-weight: 700;
  cursor: pointer;
  box-sizing: border-box;
  line-height: 1;
  text-decoration: none;
  transition: all 0.15s ease;
  border: 1px solid transparent;
  user-select: none;
}

.icon-btn svg {
  width: 14px;
  height: 14px;
  flex-shrink: 0;
  display: block;
}

.icon-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.edit-btn {
  background: #f8fafc;
  color: #475569;
  border-color: #cbd5e1;
  padding: 0 10px;
}

.edit-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
  padding: 0 10px;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.modal-form {
  display: grid;
  gap: 16px;
}

.modal-form label {
  display: grid;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 700;
  color: #0f172a;
  letter-spacing: -0.01em;
}

.modal-form input, .modal-form select, .modal-form textarea {
  width: 100%;
  padding: 12px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 4px !important;
  font-size: 0.92rem;
  color: #0f172a;
  background: #ffffff;
  outline: none;
  transition: all 0.15s ease;
}

.modal-form input:focus, .modal-form select:focus, .modal-form textarea:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
}

.submit-modal-btn {
  background: #0f172a;
  color: #ffffff;
  border: 1px solid #0f172a;
  padding: 13px 18px;
  border-radius: 4px !important;
  font-size: 0.92rem;
  font-weight: 700;
  cursor: pointer;
  margin-top: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  transition: all 0.15s ease;
}

.submit-modal-btn:hover {
  background: #1e293b;
}

.report-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
  width: 100%;
}

.excel-btn {
  background: #34c759 !important;
  color: #ffffff !important;
  border: 1px solid #34c759 !important;
  padding: 12px 20px !important;
  border-radius: 10px !important;
  font-size: 0.9rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(52, 199, 89, 0.25);
  transition: all 0.15s ease;
  line-height: 1;
  margin: 0;
  white-space: nowrap;
}

.excel-btn:hover {
  background: #28a745 !important;
  border-color: #28a745 !important;
  transform: translateY(-1px);
}

.print-btn {
  background: #007aff !important;
  color: #ffffff !important;
  border: 1px solid #007aff !important;
  padding: 12px 20px !important;
  border-radius: 10px !important;
  font-size: 0.9rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25);
  transition: all 0.15s ease;
  line-height: 1;
  margin: 0;
  white-space: nowrap;
}

.report-filter-bar {
  margin-bottom: 20px;
  padding: 14px 16px;
  background: #f8fafc;
  border-radius: 12px;
  border: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.filter-bar-title {
  font-size: 0.78rem;
  font-weight: 800;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: #64748b;
}

.filter-btn-group {
  display: flex;
  gap: 4px;
  background: #e2e8f0;
  padding: 4px;
  border-radius: 10px;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.05);
}

.filter-tab-btn {
  flex: 1;
  padding: 9px 12px;
  font-size: 0.82rem;
  font-weight: 600;
  background: transparent;
  border: none !important;
  color: #64748b;
  border-radius: 8px !important;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  white-space: nowrap;
}

.filter-tab-btn:hover {
  color: #0f172a;
  background: rgba(255, 255, 255, 0.5);
}

.filter-tab-btn.active {
  background: #ffffff !important;
  color: #0f172a !important;
  font-weight: 800 !important;
  box-shadow: 0 3px 10px rgba(15, 23, 42, 0.12), 0 1px 2px rgba(0, 0, 0, 0.04) !important;
}

/* === Mobile Responsive CSS (Android & iOS) === */
@media (max-width: 640px) {
  .page-container { padding: 16px 14px !important; }
  .page-header { flex-direction: column; align-items: stretch; gap: 12px; }
  .page-header .primary-btn { width: 100%; justify-content: center; height: 40px !important; font-size: 0.82rem !important; }
  .status-tabs { overflow-x: auto; flex-wrap: nowrap; padding-bottom: 4px; -webkit-overflow-scrolling: touch; }
  .tab-btn { padding: 8px 12px; font-size: 0.8rem; flex-shrink: 0; white-space: nowrap; }
  .card-panel { padding: 16px !important; border-radius: 14px !important; }
  .report-summary-boxes { grid-template-columns: repeat(1, 1fr); gap: 8px; }
  .report-actions { flex-direction: column; gap: 8px; }
  .excel-btn, .print-btn { width: 100%; justify-content: center; }
}
</style>
