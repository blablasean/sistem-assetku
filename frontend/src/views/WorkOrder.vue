<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Work Order System</p>
        <h1>Work Order & Perbaikan</h1>
        <p class="subtitle">Pelaporan & perbaikan kerusakan.</p>
      </div>

      <div class="header-action-group">
        <button class="primary-btn btn-secondary-ios" @click="showReportModal = true" title="Prinjau & Export Laporan Bulanan">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
          <span>Laporan & Export</span>
        </button>
        <button class="primary-btn" @click="openCreateModal">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          <span>Buat Work Order Baru</span>
        </button>
      </div>
    </div>

      <!-- 5-Box Split Status Grid (Top Row: 3 Boxes, Bottom Row: 2 Boxes) -->
      <div class="status-grid-container">
        <!-- Top Row (3 Boxes) -->
        <div class="status-grid-top">
          <button 
            :class="['status-grid-box', { active: activeTab === 'all' }]"
            @click="activeTab = 'all'"
          >
            <div class="sgb-header">
              <span class="sgb-dot dot-all"></span>
              <span class="sgb-count">{{ getTabCount('all') }}</span>
            </div>
            <span class="sgb-title">Semua Tiket</span>
          </button>

          <button 
            :class="['status-grid-box', { active: activeTab === 'Open' }]"
            @click="activeTab = 'Open'"
          >
            <div class="sgb-header">
              <span class="sgb-dot dot-open"></span>
              <span class="sgb-count">{{ getTabCount('Open') }}</span>
            </div>
            <span class="sgb-title">Open (Baru)</span>
          </button>

          <button 
            :class="['status-grid-box', { active: activeTab === 'In Progress' }]"
            @click="activeTab = 'In Progress'"
          >
            <div class="sgb-header">
              <span class="sgb-dot dot-in-progress"></span>
              <span class="sgb-count">{{ getTabCount('In Progress') }}</span>
            </div>
            <span class="sgb-title">In Progress</span>
          </button>
        </div>

        <!-- Bottom Row (2 Boxes) -->
        <div class="status-grid-bottom">
          <button 
            :class="['status-grid-box', { active: activeTab === 'Completed' }]"
            @click="activeTab = 'Completed'"
          >
            <div class="sgb-header">
              <span class="sgb-dot dot-completed"></span>
              <span class="sgb-count">{{ getTabCount('Completed') }}</span>
            </div>
            <span class="sgb-title">Selesai (Finish)</span>
          </button>

          <button 
            :class="['status-grid-box', { active: activeTab === 'Cancelled' }]"
            @click="activeTab = 'Cancelled'"
          >
            <div class="sgb-header">
              <span class="sgb-dot dot-cancelled"></span>
              <span class="sgb-count">{{ getTabCount('Cancelled') }}</span>
            </div>
            <span class="sgb-title">Dibatalkan</span>
          </button>
        </div>
      </div>

      <!-- Work Orders Table (Desktop Only) -->
      <div class="card-panel desktop-table-only">
        <div class="table-responsive wo-table-wrapper">
          <table>
            <thead>
              <tr>
                <th>ID Work Order</th>
                <th>Kategori Aset</th>
                <th>Lokasi / Area</th>
                <th>Kode / ID Aset</th>
                <th>Pelapor (Username)</th>
                <th>Departemen Asal</th>
                <th>Prioritas</th>
                <th>Deskripsi Kerusakan</th>
                <th>Teknisi (Engineer)</th>
                <th>Status</th>
                <th>Aksi Tiket</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in filteredWorkOrders" :key="wo.id" :class="{ 'highlight-assigned-wo': isAssignedToCurrentUser(wo) }">
                <td class="nowrap-cell"><span class="wo-id">#WO-{{ wo.id }}</span></td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-indigo">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2H2v10l9.29 9.29c.94.94 2.48.94 3.42 0l6.58-6.58c.94-.94.94-2.48 0-3.42L12 2Z"/><circle cx="7" cy="7" r=".5" fill="currentColor"/></svg>
                    <span>{{ wo.category || 'HVAC / AC' }}</span>
                  </div>
                </td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-blue">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>
                    <span>{{ wo.location || 'Ruangan / Area Operasional' }}</span>
                  </div>
                </td>
                <td class="nowrap-cell">Aset #{{ wo.asset_id }}</td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-slate">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                    <span>@{{ wo.requested_by || 'user_system' }}</span>
                  </div>
                </td>
                <td class="nowrap-cell">
                  <div class="ios-table-pill pill-orange">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="16" height="20" x="4" y="2" rx="2" ry="2"/><path d="M9 22v-4h6v4"/><path d="M8 6h.01"/><path d="M16 6h.01"/><path d="M12 6h.01"/><path d="M12 10h.01"/><path d="M12 14h.01"/><path d="M16 10h.01"/><path d="M16 14h.01"/><path d="M8 10h.01"/><path d="M8 14h.01"/></svg>
                    <span>{{ formatDepartmentLabel(wo.department) }}</span>
                  </div>
                </td>
                <td class="nowrap-cell"><StatusBadge :status="wo.priority || 'Medium'" /></td>
                <td class="desc-cell" :title="wo.description"><span class="desc-text">{{ wo.description }}</span></td>
                <td class="nowrap-cell">
                  <div v-if="wo.engineer_id > 0" :class="['ios-table-pill', isAssignedToCurrentUser(wo) ? 'pill-amber-assigned' : 'pill-emerald']">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14.7 6.3a1 1 0 0 0 0 1.4l1.6 1.6a1 1 0 0 0 1.4 0l3.77-3.77a6 6 0 0 1-7.94 7.94l-6.91 6.91a2.12 2.12 0 0 1-3-3l6.91-6.91a6 6 0 0 1 7.94-7.94l-3.76 3.76z"/></svg>
                    <span>{{ getEngineerName(wo.engineer_id) }}</span>
                  </div>
                  <span v-else class="unassigned-chip">Belum Ditugaskan</span>
                </td>
                <td class="nowrap-cell"><StatusBadge :status="wo.status" /></td>
                <td class="actions-cell">
                  <button class="icon-btn log-btn" @click.stop="openLogsModal(wo)" title="Lihat Laporan Timeline Progres">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                    <span>Timeline</span>
                  </button>
                  <button v-if="canAssign && wo.status === 'Open'" class="icon-btn assign-btn" @click="openAssignModal(wo)" title="Assign Worker">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><polyline points="16 11 18 13 22 9"/></svg>
                    <span>Assign</span>
                  </button>
                  <button v-if="canUpdateProgress && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn progress-btn" @click="openUpdateModal(wo)" title="Update Progres Pengerjaan">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                    <span>Update</span>
                  </button>
                  <button v-if="canManageOrder && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn close-btn" @click="closeOrder(wo)" title="Selesaikan Work Order">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                    <span>Selesai</span>
                  </button>
                  <button v-if="canCancelWo(wo)" class="icon-btn cancel-btn" @click="openCancelModal(wo)" title="Batal (Cancel Work Order)">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
                    <span>Batal</span>
                  </button>
                  <button v-if="canDelete" class="icon-btn delete-btn" @click="deleteOrder(wo)" title="Hapus Work Order Permanen">
                    <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                    <span>Hapus</span>
                  </button>
                </td>
              </tr>
              <tr v-if="filteredWorkOrders.length === 0">
                <td colspan="11" class="empty-state">Tidak ada Work Order pada kategori ini.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Mobile Work Orders Card List View (Visible on Mobile / Android) -->
      <div class="mobile-wo-list mobile-only">
        <div v-if="filteredWorkOrders.length === 0" class="mobile-empty-card">
          Tidak ada Work Order pada kategori ini.
        </div>
        <div v-else v-for="wo in filteredWorkOrders" :key="'mwo-' + wo.id" :class="['mwo-card', { 'highlight-assigned-wo': isAssignedToCurrentUser(wo) }]">
          <div class="mwo-header">
            <div class="mwo-id-wrap">
              <span class="wo-id">#WO-{{ wo.id }}</span>
              <StatusBadge :status="wo.priority || 'Medium'" />
            </div>
            <StatusBadge :status="wo.status" />
          </div>

          <div class="mwo-details">
            <div class="mwo-pill-row">
              <div class="ios-table-pill pill-blue">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>
                <span>{{ wo.location || 'Ruangan' }}</span>
              </div>
              <div class="ios-table-pill pill-indigo">
                <span>{{ wo.category || 'HVAC / AC' }}</span>
              </div>
            </div>

            <p class="mwo-desc">{{ wo.description }}</p>
            <div v-if="wo.status === 'Cancelled' && wo.alasan_pembatalan" class="mwo-cancel-reason">
              🚫 Alasan Batal: <em>"{{ wo.alasan_pembatalan }}"</em>
            </div>

            <div class="mwo-meta">
              <span>Pelapor: <strong>@{{ wo.requested_by || 'user' }}</strong> ({{ formatDepartmentLabel(wo.department) }})</span>
              <span>Teknisi: <strong>{{ wo.engineer_id > 0 ? getEngineerName(wo.engineer_id) : 'Belum Ditugaskan' }}</strong></span>
            </div>
          </div>

          <div class="mwo-actions-bar">
            <button class="icon-btn log-btn" @click.stop="openLogsModal(wo)">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              <span>Timeline</span>
            </button>
            <button v-if="canAssign && wo.status === 'Open'" class="icon-btn assign-btn" @click="openAssignModal(wo)">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><polyline points="16 11 18 13 22 9"/></svg>
              <span>Assign</span>
            </button>
            <button v-if="canUpdateProgress && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn progress-btn" @click="openUpdateModal(wo)">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
              <span>Update</span>
            </button>
            <button v-if="canManageOrder && wo.status !== 'Finish' && wo.status !== 'Cancelled'" class="icon-btn close-btn" @click="closeOrder(wo)">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
              <span>Selesai</span>
            </button>
            <button v-if="canCancelWo(wo)" class="icon-btn cancel-btn" @click="openCancelModal(wo)">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="15" y1="9" x2="9" y2="15"/><line x1="9" y1="9" x2="15" y2="15"/></svg>
              <span>Batal</span>
            </button>
            <button v-if="canDelete" class="icon-btn delete-btn" @click="deleteOrder(wo)">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
              <span>Hapus</span>
            </button>
          </div>
        </div>
      </div>

    <!-- Custom UI Toast Notification (No Browser Alert) -->
    <transition name="fade">
      <div v-if="showToast" :class="['custom-ui-toast', toastType]">
        <span class="toast-icon">{{ toastType === 'success' ? '✅' : '⚠️' }}</span>
        <span class="toast-text">{{ toastMsg }}</span>
        <button class="toast-close" @click="showToast = false">✕</button>
      </div>
    </transition>
    </div>

    <!-- Modals placed OUTSIDE .page-container so hiding .page-container during print never affects modal report -->
    <ModalDialog :show="showCreateModal" title="Ajukan Laporan Kerusakan Aset" @close="showCreateModal = false">
      <form @submit.prevent="submitWorkOrder" class="modal-form">
        <label>
          <span>Pilih Kode Aset Terdaftar</span>
          <select v-model.number="formWo.asset_id" @change="onAssetSelected" required>
            <option v-for="asset in registeredAssets" :key="asset.id" :value="asset.id">
              {{ asset.asset_code }} — {{ asset.asset_name }} ({{ asset.location }})
            </option>
          </select>
        </label>
        <label>
          <span>Kategori Kerusakan / Aset</span>
          <input v-model="formWo.category" placeholder="Contoh: HVAC / AC, Elektronik, Plumbing, Dapur" required />
        </label>
        <label>
          <span>Lokasi / Kamar Tempat Kerusakan</span>
          <input v-model="formWo.location" placeholder="Contoh: Kamar 301, Kitchen Dapur, Lift Gedung A" required />
        </label>
        <label>
          <span>Tingkat Prioritas (Priority)</span>
          <select v-model="formWo.priority">
            <option value="Emergency">Emergency (Darurat)</option>
            <option value="High">High (Perlu Penanganan Cepat)</option>
            <option value="Medium">Medium (Kerusakan Standar)</option>
            <option value="Low">Low (Tidak Mendesak)</option>
          </select>
        </label>
        <label>
          <span>Deskripsi Rincian Kerusakan</span>
          <textarea v-model="formWo.description" rows="3" placeholder="Jelaskan detail gejala kerusakan aset..." required></textarea>
        </label>

        <button type="submit" class="submit-modal-btn">Kirim Tiket Perbaikan</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showAssignModal" title="Penugasan Teknisi (Assign Worker)" @close="showAssignModal = false">
      <form @submit.prevent="submitAssign" class="modal-form" v-if="selectedWo">
        <p class="modal-info"><strong>WO #{{ selectedWo.id }}:</strong> {{ selectedWo.description }} ({{ selectedWo.location }})</p>
        <label>
          <span>Pilih Teknisi (Staff Engineer)</span>
          <select v-model.number="assignEngineerId" required>
            <option value="" disabled>-- Pilih Teknisi (Staff Engineer) --</option>
            <option v-for="eng in engineersList" :key="eng.id" :value="eng.id">
              {{ eng.name }} (@{{ eng.username }})
            </option>
          </select>
        </label>

        <button type="submit" class="submit-modal-btn">Tugaskan Teknisi</button>
      </form>
    </ModalDialog>

    <ModalDialog :show="showUpdateModal" title="Update Status Pengerjaan Teknisi" maxWidth="560px" @close="showUpdateModal = false">
      <form @submit.prevent="submitUpdateStatus" class="modal-form" v-if="selectedWo">
        <label class="modal-label">
          <span class="label-title">Status Baru <span style="color: #dc2626;">*</span></span>
          <select v-model="updateStatus" class="modal-input ios-select-input">
            <option value="In Progress">👷 In Progress (Sedang Dikerjakan)</option>
            <option value="Under Review">🔍 Under Review (Menunggu Review)</option>
            <option value="Finish">✅ Finish (Selesai & Disetujui)</option>
          </select>
        </label>

        <label class="modal-label" style="margin-top: 12px;">
          <span class="label-title">Tindakan Perbaikan yang Dilakukan <span style="color: #dc2626;">*</span></span>
          <textarea v-model="updateActionTaken" rows="3" placeholder="Misal: Ganti kapasitor kompresor AC dan isi ulang freon R32." class="modal-input modal-textarea" required></textarea>
        </label>

        <label class="modal-label" style="margin-top: 12px;">
          <span class="label-title">Estimasi Biaya Perbaikan (Rp)</span>
          <input 
            type="text" 
            :value="formattedCostInput" 
            @input="onCostInput" 
            placeholder="Rp 0" 
            class="modal-input" 
          />
        </label>

        <button type="submit" class="submit-modal-btn" style="margin-top: 16px; background: #007aff; border-color: #007aff;">Simpan Progres</button>
      </form>
    </ModalDialog>

    <!-- Printable Report Modal (Outside .page-container) -->
    <ModalDialog :show="showReportModal" title="Laporan Bulanan Work Order & Maintenance" maxWidth="920px" @close="showReportModal = false">
      <div class="monthly-report-printable" id="printableReportDocument">
        <div class="report-header">
          <h2>LAPORAN BULANAN MANAJEMEN ASET & WORK ORDER</h2>
          <p class="report-sub">Sistem AsetKu — Periode: {{ reportMonthYear }}</p>
          <hr class="report-divider" />
        </div>

        <div class="report-summary-boxes">
          <div class="rbox">
            <span>Total Tiket</span>
            <strong>{{ workOrders.length }} Tiket</strong>
          </div>
          <div class="rbox success">
            <span>Finish / Selesai</span>
            <strong>{{ countStatus('Finish') + countStatus('Completed') + countStatus('Closed') }} Tiket</strong>
          </div>
          <div class="rbox warning">
            <span>Dalam Proses</span>
            <strong>{{ countStatus('In Progress') + countStatus('Open') }} Tiket</strong>
          </div>
          <div class="rbox danger">
            <span>Total Biaya Perbaikan</span>
            <strong>Rp {{ formatNumber(totalReportCost) }}</strong>
          </div>
        </div>

        <div class="report-table-wrapper">
          <table class="report-table">
            <thead>
              <tr>
                <th class="col-id">ID</th>
                <th class="col-loc">Lokasi / Kamar</th>
                <th class="col-prio">Prioritas</th>
                <th class="col-desc">Rincian Kerusakan</th>
                <th class="col-stat">Status</th>
                <th class="col-cost">Biaya (Rp)</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wo in workOrders" :key="wo.id">
                <td class="col-id">#WO-{{ wo.id }}</td>
                <td class="col-loc">{{ wo.location }}</td>
                <td class="col-prio">{{ wo.priority }}</td>
                <td class="col-desc">{{ wo.description }}</td>
                <td class="col-stat">{{ wo.status }}</td>
                <td class="col-cost">Rp {{ formatNumber(wo.cost || 0) }}</td>
              </tr>
              <tr v-if="workOrders.length === 0">
                <td colspan="6" class="empty-state">Belum ada rekapitulasi data Work Order bulan ini.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="report-actions no-print">
        <button class="excel-btn" @click="exportToExcel">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="8" y1="13" x2="16" y2="13"/><line x1="8" y1="17" x2="16" y2="17"/></svg>
          <span>Export ke Excel (.xlsx)</span>
        </button>
        <button class="print-btn" @click="printMonthlyReport">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
          <span>Cetak Dokumen Laporan (PDF / Print)</span>
        </button>
      </div>
    </ModalDialog>

    <!-- Modal Laporan Progres & Timeline Work Order (iOS Grouped Page Card UI) -->
    <ModalDialog :show="showLogsModal" title="Timeline Progres Work Order" maxWidth="640px" @close="showLogsModal = false">
      <div v-if="selectedWoForLogs" class="ios-timeline-modal-body">
        
        <!-- Header Info Card (iOS Style) -->
        <div class="ios-asset-info-card">
          <div class="ios-aic-header">
            <span class="ios-code-badge">#WO-{{ selectedWoForLogs.id }}</span>
            <StatusBadge :status="selectedWoForLogs.status" />
          </div>
          <h3 class="ios-asset-name">📍 {{ selectedWoForLogs.location }}</h3>
          <div class="ios-registration-location">
            <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
            <span>Pelapor: <strong>@{{ selectedWoForLogs.requested_by || 'user' }}</strong> ({{ formatDepartmentLabel(selectedWoForLogs.department) }})</span>
          </div>
        </div>

        <!-- Single Consolidated Page Box Card for Timeline Items -->
        <div class="ios-timeline-page-card">
          <div class="ios-tl-header">
            <div style="display: flex; align-items: center; gap: 8px;">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              <h4>Riwayat Perubahan & Progres Pengerjaan</h4>
            </div>
            <button v-if="canUpdateProgress && selectedWoForLogs.status !== 'Finish' && selectedWoForLogs.status !== 'Cancelled'" class="icon-btn progress-btn" @click="handleNavigateToUpdateProgres(selectedWoForLogs)" title="Update Status Pengerjaan Teknisi" style="display: inline-flex; align-items: center; gap: 4px;">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="opacity: 0.85;"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
              <span style="font-size: 0.78rem;">Update Progres</span>
            </button>
          </div>
          
          <div v-if="isLogsLoading" class="ios-tl-loading">
            Memuat timeline progres...
          </div>
          
          <div v-else class="ios-tl-items-wrapper">
            <div v-for="(log, idx) in woProgressLogs" :key="log.id || idx" class="ios-tl-card-item">
              <!-- Item Top: Step Chip, Status Badge, & Time -->
              <div class="ios-tl-item-top">
                <div style="display: flex; align-items: center; gap: 6px;">
                  <span class="ios-step-chip">Langkah #{{ idx + 1 }}</span>
                  <StatusBadge :status="log.status" />
                </div>
                <span class="ios-tl-time">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                  {{ formatDate(log.created_at) }}
                </span>
              </div>

              <!-- Item Stacked Body (Vertical Stack for Perfect Inside Fit) -->
              <div class="ios-tl-stacked-body">
                <!-- Actor / User Info -->
                <div class="ios-pic-info-pill">
                  <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                  <span>Oleh: <strong>@{{ log.updated_by || 'Sistem' }}</strong> <span v-if="log.user_role">({{ formatDepartmentLabel(log.user_role) }})</span></span>
                </div>

                <!-- Action Taken / Catatan -->
                <div class="ios-reason-box" v-if="log.action_taken">
                  <p><strong>Catatan / Tindakan:</strong> {{ formatActionTaken(log.action_taken) }}</p>
                </div>

                <!-- Biaya Perbaikan (if any) -->
                <div v-if="log.cost > 0" class="ios-cost-pill">
                  💰 Biaya Terkait: <strong>Rp {{ formatNumber(log.cost) }}</strong>
                </div>
              </div>
            </div>

            <div v-if="woProgressLogs.length === 0" class="ios-empty-tl">
              Belum ada riwayat progres tercatat untuk Work Order ini.
            </div>
          </div>
        </div>

      </div>
    </ModalDialog>

    <!-- Modal Pembatalan Work Order (iOS Style Danger Badge) -->
    <ModalDialog :show="showCancelModal" title="Membatalkan Work Order" @close="showCancelModal = false">
      <form @submit.prevent="submitCancelOrder" class="modal-form" v-if="selectedWo">
        <div class="ios-cancel-warning-banner">
          <div class="ios-danger-circle-icon">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#dc2626" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
          </div>
          <div class="ios-warning-text">
            <h4>Membatalkan Tiket #WO-{{ selectedWo.id }}</h4>
            <p>{{ selectedWo.description }}</p>
          </div>
        </div>

        <label>
          <span>Alasan Pembatalan Tiket <span style="color: #dc2626;">*</span></span>
          <textarea v-model="cancelReason" rows="3" placeholder="Tuliskan alasan pembatalan tiket ini (misal: salah input lokasi / sudah diperbaiki sendiri / duplikat laporan)" required class="sharp-input" style="width: 100%; min-height: 80px; padding: 10px; border-radius: 8px;"></textarea>
        </label>

        <div class="modal-form-actions" style="display: flex; justify-content: flex-end; gap: 10px; margin-top: 14px;">
          <button type="button" class="cancel-btn-ios" @click="showCancelModal = false">Kembali</button>
          <button type="submit" class="danger-btn-ios" :disabled="isSubmittingCancel">
            <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="4.93" y1="4.93" x2="19.07" y2="19.07"/></svg>
            <span>{{ isSubmittingCancel ? 'Memproses...' : 'Konfirmasi Batal Tiket' }}</span>
          </button>
        </div>
      </form>
    </ModalDialog>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const route = useRoute()

const userRole = ref(sessionStorage.getItem('user_role') || 'external')

const canAssign = computed(() => ['hod', 'management', 'supervisor', 'admin'].includes(userRole.value.toLowerCase()))
const canManageOrder = computed(() => ['hod', 'management', 'supervisor', 'admin'].includes(userRole.value.toLowerCase()))
const canUpdateProgress = computed(() => ['engineer', 'hod', 'management', 'supervisor', 'admin'].includes(userRole.value.toLowerCase()))
const canDelete = computed(() => ['hod', 'admin', 'management', 'supervisor'].includes(userRole.value.toLowerCase()))
const canCancel = computed(() => ['hod', 'admin', 'management', 'supervisor'].includes(userRole.value.toLowerCase()))

function isWoRequester(wo) {
  if (!wo || !currentUsername.value) return false
  return wo.requested_by && wo.requested_by.toLowerCase() === currentUsername.value.toLowerCase()
}

function canCancelWo(wo) {
  if (!wo || wo.status !== 'Open') return false
  if (canCancel.value) return true
  return isWoRequester(wo)
}

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

const activeTab = ref('all')
const tabs = [
  { id: 'all', label: 'Semua Tiket' },
  { id: 'open', label: 'Open' },
  { id: 'in_progress', label: 'In Progress' },
  { id: 'under_review', label: 'Under Review' },
  { id: 'completed', label: 'Finish / Completed' }
]

const workOrders = ref([])
const isLoading = ref(false)

const showCreateModal = ref(false)
const formWo = ref({ asset_id: '', category: '', location: '', priority: 'Emergency', description: '' })

function openCreateModal() {
  const topAsset = registeredAssets.value[0]
  formWo.value = {
    asset_id: topAsset?.id || '',
    category: topAsset?.category || '',
    location: topAsset?.location || '',
    priority: 'Emergency',
    description: ''
  }
  if (topAsset) {
    onAssetSelected()
  }
  showCreateModal.value = true
}

const showAssignModal = ref(false)
const selectedWo = ref(null)
const assignEngineerId = ref('')
const engineersList = ref([])

function getEngineerName(id) {
  if (!id) return 'Belum Ditugaskan'
  const eng = engineersList.value.find(e => Number(e.id) === Number(id))
  return eng ? (eng.name || `@${eng.username}`) : 'Teknisi'
}

function formatRupiahDisplay(val) {
  const num = Number(val) || 0
  if (num === 0) return '0,00 Rp'
  const formatted = num.toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
  return `${formatted} Rp`
}

function formatActionTaken(text) {
  if (!text) return ''
  return text.replace(/Teknisi\s*#(\d+)/gi, (match, idStr) => {
    const engName = getEngineerName(Number(idStr))
    return (engName && engName !== 'Teknisi' && engName !== 'Belum Ditugaskan') ? `Teknisi ${engName}` : match
  })
}

async function fetchEngineers() {
  try {
    const res = await api.get('/users/engineers')
    if (res.data?.data && Array.isArray(res.data.data)) {
      engineersList.value = res.data.data
      if (engineersList.value.length > 0 && !assignEngineerId.value) {
        assignEngineerId.value = engineersList.value[0].id
      }
    }
  } catch (e) {
    console.error('Failed to fetch engineers list:', e)
  }
}

const showUpdateModal = ref(false)
const updateStatus = ref('In Progress')
const updateActionTaken = ref('')
const updateCost = ref(0)
const formattedCostInput = ref('Rp 0')

function onCostInput(event) {
  const inputVal = event.target.value || ''
  const digitsOnly = inputVal.replace(/\D/g, '')

  if (!digitsOnly || digitsOnly === '0') {
    updateCost.value = 0
    formattedCostInput.value = 'Rp 0'
  } else {
    const numeric = parseInt(digitsOnly, 10)
    updateCost.value = numeric
    formattedCostInput.value = 'Rp ' + numeric.toLocaleString('id-ID')
  }
}

const showReportModal = ref(false)
const reportMonthYear = ref(new Date().toLocaleDateString('id-ID', { month: 'long', year: 'numeric' }))

const showLogsModal = ref(false)
const selectedWoForLogs = ref(null)
const woProgressLogs = ref([])
const isLogsLoading = ref(false)
const newTlStatus = ref('In Progress')
const newTlActionTaken = ref('')
const newTlCost = ref(0)
const isSubmittingTl = ref(false)

function getStatusIcon(status) {
  const map = {
    Open: '🚨',
    'In Progress': '👷',
    'Under Review': '🔍',
    Finish: '✅',
    Completed: '✅',
    Cancelled: '🚫'
  }
  return map[status] || '📌'
}

function handleNavigateToUpdateProgres(wo) {
  showLogsModal.value = false
  openUpdateModal(wo)
}

function openLogsModal(wo) {
  if (!wo) return
  if (engineersList.value.length === 0) {
    fetchEngineers()
  }
  selectedWoForLogs.value = wo
  showLogsModal.value = true

  newTlStatus.value = wo.status === 'Open' ? 'In Progress' : wo.status
  newTlActionTaken.value = ''
  newTlCost.value = 0

  // Build instant milestone timeline from wo object
  const initialLogs = [
    {
      id: 1,
      work_order_id: wo.id,
      status: 'Open',
      action_taken: `Laporan diajukan: ${wo.description || ''}`,
      updated_by: wo.requested_by || 'Staff Operasional',
      user_role: wo.department || 'User',
      created_at: wo.created_at || new Date().toISOString()
    }
  ]

  if (wo.status && wo.status !== 'Open') {
    initialLogs.push({
      id: 2,
      work_order_id: wo.id,
      status: 'In Progress',
      action_taken: `Penugasan Teknisi untuk perbaikan di lokasi ${wo.location || ''}`,
      updated_by: 'HOD Engineer',
      user_role: 'HOD Engineer',
      created_at: wo.created_at || new Date().toISOString()
    })
  }

  if (wo.status === 'Under Review' || wo.status === 'Finish' || wo.status === 'Completed' || wo.status === 'Closed') {
    initialLogs.push({
      id: 3,
      work_order_id: wo.id,
      status: 'Under Review',
      action_taken: wo.action_taken || 'Perbaikan unit selesai dikerjakan. Menunggu review.',
      cost: wo.cost || 0,
      updated_by: 'Budi Santoso (Teknisi)',
      user_role: 'Staff Engineer',
      created_at: wo.created_at || new Date().toISOString()
    })
  }

  if (wo.status === 'Finish' || wo.status === 'Completed' || wo.status === 'Closed' || wo.status === 'Cancelled') {
    const cancelText = wo.alasan_pembatalan ? `Work order dibatalkan oleh ${wo.requested_by || 'user'}. Alasan: ${wo.alasan_pembatalan}` : `Work order dibatalkan oleh ${wo.requested_by || 'user'}`
    initialLogs.push({
      id: 4,
      work_order_id: wo.id,
      status: wo.status,
      action_taken: wo.status === 'Cancelled' ? cancelText : 'Work Order diverifikasi selesai',
      cost: wo.cost || 0,
      updated_by: wo.status === 'Cancelled' ? (wo.requested_by || 'User') : 'Administrator',
      user_role: wo.status === 'Cancelled' ? (wo.department || 'User') : 'Admin',
      created_at: wo.closed_at || new Date().toISOString()
    })
  }

  woProgressLogs.value = initialLogs
  isLogsLoading.value = false

  // Background fetch live timeline records from MySQL DB
  api.get(`/workorders/timeline?wo_id=${wo.id}`).then(res => {
    const logsData = res.data?.data || res.data
    if (Array.isArray(logsData) && logsData.length > 0) {
      woProgressLogs.value = logsData
    }
  }).catch(e => {
    console.error('Background fetch WO timeline error:', e)
  })
}

async function submitAddTimelineNote() {
  if (!selectedWoForLogs.value || !newTlActionTaken.value) return
  isSubmittingTl.value = true
  try {
    const payload = {
      work_order_id: selectedWoForLogs.value.id,
      status: newTlStatus.value,
      action_taken: newTlActionTaken.value,
      cost: newTlCost.value || 0
    }
    await api.post('/workorders/timeline/add', payload)
    notify('Catatan timeline berhasil ditambahkan!', 'success')
    newTlActionTaken.value = ''
    newTlCost.value = 0
    selectedWoForLogs.value.status = newTlStatus.value
    if (payload.cost > 0) {
      selectedWoForLogs.value.cost = payload.cost
    }
    
    // Refresh live timeline
    const res = await api.get(`/workorders/timeline?wo_id=${selectedWoForLogs.value.id}`)
    const logsData = res.data?.data || res.data
    if (Array.isArray(logsData) && logsData.length > 0) {
      woProgressLogs.value = logsData
    }
    fetchWorkOrders(true)
  } catch (e) {
    console.error('Failed to add timeline note:', e)
    notify(e.response?.data?.error || 'Gagal menambahkan catatan timeline.', 'error')
  } finally {
    isSubmittingTl.value = false
  }
}

const filteredWorkOrders = computed(() => {
  if (activeTab.value === 'all') return workOrders.value
  if (activeTab.value === 'Open') return workOrders.value.filter(w => w.status === 'Open')
  if (activeTab.value === 'In Progress') return workOrders.value.filter(w => w.status === 'In Progress' || w.status === 'Under Review')
  if (activeTab.value === 'Completed') return workOrders.value.filter(w => w.status === 'Finish' || w.status === 'Completed' || w.status === 'Closed')
  if (activeTab.value === 'Cancelled') return workOrders.value.filter(w => w.status === 'Cancelled' || w.status === 'Batal')
  return workOrders.value
})

const totalReportCost = computed(() => {
  return workOrders.value.reduce((sum, w) => sum + (w.cost || 0), 0)
})

function getTabCount(tabId) {
  if (tabId === 'all') return workOrders.value.length
  if (tabId === 'Open') return workOrders.value.filter(w => w.status === 'Open').length
  if (tabId === 'In Progress') return workOrders.value.filter(w => w.status === 'In Progress' || w.status === 'Under Review').length
  if (tabId === 'Completed') return workOrders.value.filter(w => w.status === 'Finish' || w.status === 'Completed' || w.status === 'Closed').length
  if (tabId === 'Cancelled') return workOrders.value.filter(w => w.status === 'Cancelled' || w.status === 'Batal').length
  return 0
}

function countStatus(st) {
  return workOrders.value.filter(w => w.status === st).length
}

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

const currentUserId = ref(Number(sessionStorage.getItem('user_id')) || 0)
const currentUsername = ref(sessionStorage.getItem('username') || sessionStorage.getItem('user_name') || 'admin')
const currentUserRole = ref(sessionStorage.getItem('user_role') || 'admin')

function isAssignedToCurrentUser(wo) {
  if (!wo || !wo.engineer_id) return false
  
  if (currentUserId.value > 0 && Number(wo.engineer_id) === currentUserId.value) {
    return true
  }

  if (wo.engineer_id > 0 && engineersList.value.length > 0) {
    const eng = engineersList.value.find(e => Number(e.id) === Number(wo.engineer_id))
    if (eng) {
      if (eng.username && currentUsername.value && eng.username.toLowerCase() === currentUsername.value.toLowerCase()) {
        return true
      }
      if (eng.name && currentUsername.value && eng.name.toLowerCase() === currentUsername.value.toLowerCase()) {
        return true
      }
    }
  }

  return false
}

function formatDepartmentLabel(roleOrDept) {
  if (!roleOrDept) return 'Staff Operasional'
  const map = {
    dept_akunting: 'Departement Akunting',
    dept_spa: 'Departement Spa',
    dept_sales: 'Department Sales',
    dept_hr: 'Department Human Resources',
    dept_fb_kitchen: 'Department Food Beverage Kitchen',
    dept_fb_service: 'Department Food Beverage Service',
    dept_housekeeping: 'Department House Keeping',
    dept_frontoffice: 'Department Front Office',
    admin: 'Administrator',
    hod: 'HOD Engineer',
    management: 'Supervisor Engineer',
    engineer: 'Staff Engineer',
    external: 'Staff Operasional'
  }
  return map[roleOrDept] || roleOrDept
}

async function fetchWorkOrders(isSilent = false) {
  if (!isSilent) isLoading.value = true
  try {
    const res = await api.get('/workorders')
    if (res.data?.data && Array.isArray(res.data.data)) {
      workOrders.value = res.data.data
    } else {
      workOrders.value = []
    }
  } catch (e) {
    if (!isSilent) console.error('Failed to fetch work orders from DB:', e)
  } finally {
    if (!isSilent) isLoading.value = false
  }
}

async function submitWorkOrder() {
  try {
    await api.post('/workorders', {
      asset_id: Number(formWo.value.asset_id) || 1,
      category: formWo.value.category || 'HVAC / AC',
      location: formWo.value.location || 'Ruangan / Area Operasional',
      priority: formWo.value.priority,
      description: formWo.value.description,
      requested_by: currentUsername.value,
      department: currentUserRole.value
    })
    showCreateModal.value = false
    formWo.value = { asset_id: '', category: '', location: '', priority: 'Medium', description: '' }
    notify('Tiket Work Order / Laporan Kerusakan berhasil diajukan ke database!', 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal mengajukan Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openAssignModal(wo) {
  selectedWo.value = wo
  if (engineersList.value.length > 0) {
    assignEngineerId.value = engineersList.value[0].id
  } else {
    assignEngineerId.value = ''
  }
  showAssignModal.value = true
}

async function submitAssign() {
  if (!selectedWo.value) return
  try {
    await api.post('/workorders/assign', {
      wo_id: selectedWo.value.id,
      engineer_id: assignEngineerId.value
    })
    showAssignModal.value = false
    const engName = getEngineerName(assignEngineerId.value)
    notify(`Teknisi "${engName}" berhasil ditugaskan ke WO #${selectedWo.value.id}!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal menugaskan teknisi: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openUpdateModal(wo) {
  selectedWo.value = wo
  updateStatus.value = wo.status === 'Open' ? 'In Progress' : wo.status
  updateActionTaken.value = ''
  const initialCost = wo.cost || 0
  updateCost.value = initialCost
  if (!initialCost || initialCost === 0) {
    formattedCostInput.value = 'Rp 0'
  } else {
    formattedCostInput.value = 'Rp ' + initialCost.toLocaleString('id-ID')
  }
  showUpdateModal.value = true
}

async function submitUpdateStatus() {
  if (!selectedWo.value) return
  try {
    await api.post('/workorders/status', {
      wo_id: selectedWo.value.id,
      status: updateStatus.value,
      action_taken: updateActionTaken.value,
      cost: updateCost.value
    })
    showUpdateModal.value = false
    notify(`Status Work Order #${selectedWo.value.id} diperbarui menjadi "${updateStatus.value}"!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal memperbarui status: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function closeOrder(wo) {
  try {
    await api.post('/workorders/close', { wo_id: wo.id })
    notify(`Work Order #${wo.id} berhasil diselesaikan (Finish)!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal menutup Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

const showCancelModal = ref(false)
const cancelReason = ref('')
const isSubmittingCancel = ref(false)

function openCancelModal(wo) {
  selectedWo.value = wo
  cancelReason.value = ''
  showCancelModal.value = true
}

async function submitCancelOrder() {
  if (!selectedWo.value || !cancelReason.value.trim()) {
    notify('Harap sertakan alasan pembatalan tiket.', 'error')
    return
  }
  isSubmittingCancel.value = true
  try {
    await api.post('/workorders/cancel', {
      wo_id: selectedWo.value.id,
      reason: cancelReason.value.trim()
    })
    showCancelModal.value = false
    notify(`Work Order #WO-${selectedWo.value.id} berhasil dibatalkan!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal membatalkan Work Order: ' + (e.response?.data?.message || e.message), 'error')
  } finally {
    isSubmittingCancel.value = false
  }
}

async function deleteOrder(wo) {
  try {
    await api.post('/workorders/delete', { wo_id: wo.id })
    workOrders.value = workOrders.value.filter(w => w.id !== wo.id)
    notify(`Work Order #${wo.id} berhasil dihapus permanen!`, 'success')
    await fetchWorkOrders()
  } catch (e) {
    notify('Gagal menghapus Work Order: ' + (e.response?.data?.message || e.message), 'error')
  }
}

import { exportToExcel as exportToExcelHelper, triggerPrint } from '../utils/exportUtils'

function printMonthlyReport() {
  triggerPrint()
}

function exportToExcel() {
  const monthName = reportMonthYear.value.replace(/\s+/g, '_')
  const fileName = `Laporan_WorkOrder_${monthName}.xls`
  const headers = ['Kode WO', 'Lokasi / Kamar', 'Prioritas', 'Deskripsi Kerusakan', 'Status Tiket', 'Biaya Perbaikan (Rp)']
  const rows = workOrders.value.map(wo => [
    wo.wo_code || `#WO-${wo.id}`,
    wo.location || '-',
    wo.priority || '-',
    wo.description || '-',
    wo.status || '-',
    wo.cost || 0
  ])
  exportToExcelHelper(fileName, headers, rows)
}

const registeredAssets = ref([])

async function fetchRegisteredAssets() {
  try {
    const res = await api.get('/assets')
    if (res.data?.data && Array.isArray(res.data.data)) {
      registeredAssets.value = res.data.data
    }
  } catch (e) {
    console.error('Failed to fetch registered assets:', e)
  }
}

function onAssetSelected() {
  const selected = registeredAssets.value.find(a => a.id === formWo.value.asset_id)
  if (selected) {
    if (!formWo.value.category) formWo.value.category = selected.category || ''
    if (!formWo.value.location) formWo.value.location = selected.location || ''
  }
}

let pollTimer = null

onMounted(() => {
  if (route.query.assetId) {
    formWo.value.asset_id = parseInt(route.query.assetId)
    formWo.value.location = route.query.location || ''
    showCreateModal.value = true
  }
  fetchWorkOrders()
  fetchRegisteredAssets()
  fetchEngineers()

  // Real-time status sync across all roles & screens
  pollTimer = setInterval(() => {
    fetchWorkOrders(true)
  }, 3000)
})

onUnmounted(() => {
  if (pollTimer) {
    clearInterval(pollTimer)
  }
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
  flex-wrap: wrap;
  gap: 16px;
}

.header-action-group {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 37.6px;
}

.primary-btn {
  background: #007aff !important;
  color: #ffffff !important;
  border: 1px solid #007aff !important;
  height: 37.6px !important;
  padding: 0 18px !important;
  border-radius: 10px !important;
  font-size: 0.88rem !important;
  font-weight: 700 !important;
  cursor: pointer;
  display: inline-flex !important;
  align-items: center !important;
  justify-content: center !important;
  gap: 6px !important;
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.25) !important;
  transition: all 0.15s ease !important;
  line-height: 1 !important;
  white-space: nowrap !important;
  box-sizing: border-box !important;
  margin: 0 !important;
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.primary-btn svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
  display: block;
}

.primary-btn span {
  display: inline-block;
  line-height: 1;
  white-space: nowrap;
}

.primary-btn.btn-secondary-ios, .btn-secondary-ios {
  background: #ffffff !important;
  color: #0f172a !important;
  border: 1px solid #cbd5e1 !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04) !important;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1) !important;
}

.primary-btn.btn-secondary-ios:hover, .btn-secondary-ios:hover {
  background: #f8fafc !important;
  border-color: #94a3b8 !important;
  color: #0284c7 !important;
  transform: translateY(-1.5px) !important;
  box-shadow: 0 4px 14px rgba(2, 132, 199, 0.15) !important;
}

.primary-btn.btn-secondary-ios:active, .btn-secondary-ios:active {
  transform: scale(0.97) !important;
}

.print-report-btn {
  background: #ffffff;
  color: #0f172a;
  border: 1px solid #cbd5e1;
  padding: 10px 18px;
  border-radius: 4px !important;
  font-weight: 700;
  font-size: 0.88rem;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
}

.print-report-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8;
}

.status-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  border-bottom: 2px solid #e2e8f0;
  padding-bottom: 4px;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tab-btn {
  background: transparent;
  border: none;
  padding: 8px 16px;
  font-size: 0.9rem;
  font-weight: 600;
  color: #64748b;
  cursor: pointer;
  border-bottom: 3px solid transparent;
  display: flex;
  align-items: center;
  gap: 6px;
  white-space: nowrap;
}

.tab-btn.active {
  color: #2563eb;
  border-bottom-color: #2563eb;
}

.tab-count {
  background: #e2e8f0;
  color: #334155;
  padding: 2px 6px;
  border-radius: 999px;
  font-size: 0.75rem;
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.ios-table-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  border-radius: 6px !important;
  font-size: 0.82rem;
  font-weight: 700;
  line-height: 1;
  white-space: nowrap;
}

.pill-indigo { background: #e0e7ff; color: #3730a3; border: 1px solid #c7d2fe; }
.pill-blue { background: #eff6ff; color: #1e40af; border: 1px solid #bfdbfe; }
.pill-slate { background: #f1f5f9; color: #334155; border: 1px solid #cbd5e1; }
.pill-orange { background: #fff7ed; color: #c2410c; border: 1px solid #ffedd5; }
.pill-emerald { background: #ecfdf5; color: #065f46; border: 1px solid #a7f3d0; }

.wo-table-wrapper {
  overflow-x: auto;
  max-width: 100%;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 900px;
}

th {
  text-align: left;
  padding: 12px 14px;
  background: #f8fafc;
  color: #475569;
  font-size: 0.85rem;
  border-bottom: 1px solid #e2e8f0;
  white-space: nowrap;
}

td {
  padding: 14px 14px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.9rem;
  color: #334155;
  vertical-align: middle;
}

.nowrap-cell {
  white-space: nowrap;
}

.desc-cell {
  max-width: 220px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.desc-text {
  font-weight: 500;
  display: block;
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.wo-id {
  font-family: monospace;
  font-weight: 800;
  color: #0f172a;
}

.engineer-chip {
  background: #dcfce7;
  color: #15803d;
  padding: 3px 8px;
  border-radius: 6px;
  font-size: 0.8rem;
  font-weight: 600;
  white-space: nowrap;
}

.unassigned-chip {
  color: #94a3b8;
  font-size: 0.8rem;
  font-style: italic;
  white-space: nowrap;
}

.actions-cell {
  display: flex;
  flex-direction: row;
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

.log-btn {
  background: #f8fafc;
  color: #475569;
  border: 1px solid #cbd5e1 !important;
}

.log-btn:hover {
  background: #f1f5f9;
  border-color: #94a3b8 !important;
  color: #1e293b;
}

.assign-btn {
  background: #f0fdf4;
  color: #15803d;
  border: 1px solid #86efac !important;
}

.assign-btn:hover {
  background: #dcfce7;
  border-color: #22c55e !important;
}

.progress-btn {
  background: #eff6ff;
  color: #2563eb;
  border: 1px solid #60a5fa !important;
  box-shadow: 0 1px 2px rgba(37, 99, 235, 0.08);
}

.progress-btn:hover {
  background: #dbeafe;
  border-color: #2563eb !important;
  color: #1d4ed8;
}

.close-btn {
  background: #f0fdf4;
  color: #15803d;
  border: 1px solid #4ade80 !important;
  box-shadow: 0 1px 2px rgba(21, 128, 61, 0.08);
}

.close-btn:hover {
  background: #dcfce7;
  border-color: #16a34a !important;
  color: #166534;
}

.cancel-btn {
  background: #fff7ed;
  color: #ea580c;
  border: 1px solid #fdba74 !important;
}

.cancel-btn:hover {
  background: #ffedd5;
  border-color: #f97316 !important;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fca5a5 !important;
  box-shadow: 0 1px 2px rgba(220, 38, 38, 0.08);
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #ef4444 !important;
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
  max-width: 100%;
  box-sizing: border-box;
  padding: 12px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 4px !important;
  font-size: 0.92rem;
  color: #0f172a;
  background: #ffffff;
  outline: none;
  transition: all 0.15s ease;
}

.modal-form select, .ios-select-input {
  max-width: 100%;
  box-sizing: border-box;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
}

.modal-form select option, .ios-select-input option {
  max-width: 100%;
  font-size: 0.88rem;
  padding: 6px 10px;
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
  transition: all 0.15s ease;
}

.submit-modal-btn:hover {
  background: #1e293b;
}

.modal-info {
  background: #f8fafc;
  padding: 10px;
  border-radius: 8px;
  font-size: 0.9rem;
}

.empty-state {
  text-align: center;
  color: #94a3b8;
  padding: 24px;
}

/* Printable Monthly Report Styling */
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
  font-weight: 600;
}

.report-divider {
  border: 0;
  height: 2px;
  background: #0f172a;
  margin: 12px 0 20px;
}

.report-summary-boxes {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.rbox {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  padding: 12px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.rbox span { font-size: 0.75rem; color: #64748b; font-weight: 600; }
.rbox strong { font-size: 1.1rem; color: #0f172a; margin-top: 4px; }
.rbox.success strong { color: #16a34a; }
.rbox.warning strong { color: #d97706; }
.rbox.danger strong { color: #dc2626; }

.report-table-wrapper {
  width: 100%;
  overflow-x: auto;
}

.report-table {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 24px;
  table-layout: fixed;
}

.report-table th, .report-table td {
  border: 1px solid #cbd5e1;
  padding: 8px 10px;
  font-size: 0.85rem;
  word-break: break-word;
  white-space: normal;
}

.report-table th {
  background: #f1f5f9;
  font-weight: 700;
}

.report-actions {
  margin-top: 24px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
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
}

.print-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.requester-chip {
  background: #f1f5f9;
  color: #0f172a;
  padding: 3px 8px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.8rem;
  border: 1px solid #cbd5e1;
}

.dept-chip {
  background: #dbeafe;
  color: #1e40af;
  padding: 3px 8px;
  border-radius: 2px !important;
  font-weight: 700;
  font-size: 0.8rem;
  border: 1px solid #bfdbfe;
}

.log-btn {
  background: #f8fafc;
  color: #0284c7;
  border-color: #7dd3fc;
}

.logs-modal-body {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.wo-info-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  padding: 14px 16px;
}

.wo-badge {
  font-size: 0.78rem;
  background: #0f172a;
  color: white;
  padding: 2px 8px;
  border-radius: 2px !important;
  font-weight: 800;
}

.wo-banner-title {
  margin: 6px 0 2px;
  font-size: 1.1rem;
  color: #0f172a;
}

.wo-banner-sub {
  margin: 0;
  font-size: 0.85rem;
  color: #64748b;
}

.timeline-container {
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  padding: 16px;
  background: #ffffff;
}

.timeline-title {
  margin: 0 0 16px;
  font-size: 0.95rem;
  color: #0f172a;
  font-weight: 800;
}

.timeline-list {
  display: flex;
  flex-direction: column;
  position: relative;
  padding-left: 20px;
}

.timeline-list::before {
  content: '';
  position: absolute;
  left: 9px;
  top: 8px;
  bottom: 8px;
  width: 2px;
  background: #e2e8f0;
}

.timeline-item {
  position: relative;
  margin-bottom: 20px;
  display: flex;
  gap: 14px;
}

.timeline-item:last-child {
  margin-bottom: 0;
}

.timeline-node {
  position: absolute;
  left: -20px;
  top: 0;
  width: 20px;
  height: 20px;
  background: #ffffff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
  z-index: 1;
}

.timeline-content {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 2px !important;
  padding: 12px 14px;
  width: 100%;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.timeline-time {
  font-size: 0.78rem;
  color: #64748b;
  font-weight: 600;
}

.timeline-actor {
  margin: 0 0 6px;
  font-size: 0.85rem;
  color: #334155;
}

.user-role-chip {
  color: #0284c7;
  font-weight: 700;
}

.timeline-notes {
  background: #ffffff;
  border: 1px solid #cbd5e1;
  padding: 8px 10px;
  border-radius: 2px !important;
  font-size: 0.85rem;
  color: #0f172a;
}

.timeline-notes p {
  margin: 0;
}

.timeline-cost {
  margin: 8px 0 0;
  font-size: 0.85rem;
  color: #16a34a;
}

.logs-loading, .empty-logs {
  text-align: center;
  padding: 24px;
  color: #64748b;
  font-size: 0.9rem;
}

.add-timeline-box {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  border-radius: 2px !important;
  padding: 16px;
}

.add-tl-title {
  margin: 0 0 12px;
  font-size: 0.95rem;
  color: #0f172a;
  font-weight: 800;
}

.add-tl-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.add-tl-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.tl-field {
  display: flex;
  flex-direction: column;
  gap: 4px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #334155;
}

.add-tl-btn {
  align-self: flex-end;
  padding: 8px 16px;
  font-weight: 700;
}

/* === 5-Box Split Status Grid === */
.status-grid-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 20px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px !important;
  padding: 10px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.03);
}

.status-grid-top {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 8px;
}

.status-grid-bottom {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.status-grid-box {
  background: #f8fafc;
  border: 1.5px solid #e2e8f0;
  border-radius: 10px !important;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 6px;
  cursor: pointer;
  transition: all 0.18s ease;
  text-align: left;
}

.status-grid-box:hover {
  background: #ffffff;
  border-color: #cbd5e1;
  transform: translateY(-1px);
}

.status-grid-box.active {
  background: #007aff !important;
  border-color: #007aff !important;
  box-shadow: 0 4px 14px rgba(0, 122, 255, 0.25);
}

.sgb-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sgb-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  display: inline-block;
}

.dot-all { background: #2563eb; }
.dot-open { background: #3b82f6; }
.dot-in-progress { background: #f59e0b; }
.dot-completed { background: #10b981; }
.dot-cancelled { background: #ef4444; }

.status-grid-box.active .sgb-dot {
  background: #ffffff !important;
  box-shadow: 0 0 6px rgba(255, 255, 255, 0.8);
}

.sgb-count {
  font-size: 1.1rem;
  font-weight: 800;
  color: #0f172a;
}

.status-grid-box.active .sgb-count {
  color: #ffffff !important;
}

.sgb-title {
  font-size: 0.78rem;
  font-weight: 700;
  color: #64748b;
  line-height: 1.2;
}

.status-grid-box.active .sgb-title {
  color: #ffffff !important;
}

/* === Desktop vs Mobile Display Toggle === */
.mobile-only {
  display: none !important;
}

/* === Mobile Responsive CSS (Android & iOS) === */
@media (max-width: 640px) {
  .status-grid-container {
    padding: 8px;
    gap: 6px;
    margin-bottom: 14px;
  }

  .status-grid-top, .status-grid-bottom {
    gap: 6px;
  }

  .status-grid-box {
    padding: 8px 10px;
  }

  .sgb-count {
    font-size: 1rem;
  }

  .sgb-title {
    font-size: 0.74rem;
  }

  .desktop-table-only {
    display: none !important;
  }

  .mobile-only {
    display: flex !important;
    flex-direction: column;
    gap: 12px;
  }

  .page-container { padding: 14px 10px !important; }
  .page-header { flex-direction: column; align-items: stretch; gap: 12px; }
  .header-action-group { width: 100%; display: flex; flex-wrap: wrap; gap: 8px; }
  .header-action-group .primary-btn { flex: 1; min-width: 130px; justify-content: center; height: 40px !important; font-size: 0.82rem !important; }
  .card-panel { padding: 16px !important; border-radius: 14px !important; }
  .report-summary-boxes { grid-template-columns: repeat(2, 1fr); gap: 8px; }
  .report-actions { flex-direction: column; gap: 8px; }
  .excel-btn, .print-btn { width: 100%; justify-content: center; }

  .mobile-wo-list {
    width: 100%;
  }

  .mwo-card {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 14px !important;
    padding: 14px 16px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .mwo-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f1f5f9;
  }

  .mwo-id-wrap {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .mwo-details {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .mwo-pill-row {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
  }

  .mwo-desc {
    margin: 2px 0;
    font-size: 0.84rem;
    color: #0f172a;
    font-weight: 600;
    line-height: 1.35;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .mwo-meta {
    display: flex;
    flex-direction: column;
    gap: 4px;
    font-size: 0.78rem;
    color: #64748b;
    background: #f8fafc;
    padding: 8px 10px;
    border-radius: 8px;
  }

  .mwo-cancel-reason {
    font-size: 0.78rem;
    color: #991b1b;
    background: #fef2f2;
    border: 1px solid #fecaca;
    padding: 6px 10px;
    border-radius: 8px;
    font-weight: 600;
    margin: 2px 0 4px;
  }

  .mwo-actions-bar {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
    padding-top: 8px;
    border-top: 1px solid #f1f5f9;
  }

  .mwo-actions-bar .icon-btn {
    padding: 6px 10px !important;
    font-size: 0.75rem !important;
    border-radius: 6px !important;
  }

  .mobile-empty-card {
    background: #ffffff;
    border: 1px dashed #cbd5e1;
    border-radius: 12px;
    padding: 24px;
    text-align: center;
    color: #64748b;
    font-size: 0.88rem;
  }
}

.cancel-btn {
  background: #fef2f2 !important;
  color: #dc2626 !important;
  border: 1px solid #fecaca !important;
  display: inline-flex !important;
  align-items: center !important;
  gap: 5px !important;
  border-radius: 8px !important;
  font-weight: 700 !important;
  transition: all 0.15s ease;
}

.cancel-btn:hover {
  background: #fee2e2 !important;
  border-color: #fca5a5 !important;
  color: #b91c1c !important;
}

.ios-cancel-warning-banner {
  background: #fff5f5;
  border: 1px solid #fed7d7;
  border-radius: 12px;
  padding: 12px 14px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 12px;
}

.ios-danger-circle-icon {
  width: 42px;
  height: 42px;
  background: #fee2e2;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid #fca5a5;
}

.ios-warning-text h4 {
  margin: 0 0 3px;
  font-size: 0.92rem;
  color: #991b1b;
  font-weight: 800;
}

.ios-warning-text p {
  margin: 0;
  font-size: 0.8rem;
  color: #7f1d1d;
  line-height: 1.3;
}

.danger-btn-ios {
  background: #dc2626 !important;
  color: #ffffff !important;
  border: 1px solid #dc2626 !important;
  padding: 8px 16px !important;
  border-radius: 10px !important;
  font-weight: 700 !important;
  font-size: 0.85rem !important;
  display: inline-flex !important;
  align-items: center !important;
  gap: 6px !important;
  box-shadow: 0 4px 12px rgba(220, 38, 38, 0.25) !important;
  cursor: pointer;
}

.danger-btn-ios:hover {
  background: #b91c1c !important;
  border-color: #b91c1c !important;
}

.cancel-btn-ios {
  background: #f1f5f9 !important;
  color: #475569 !important;
  border: 1px solid #cbd5e1 !important;
  padding: 8px 16px !important;
  border-radius: 10px !important;
  font-weight: 700 !important;
  font-size: 0.85rem !important;
  cursor: pointer;
}

/* === iOS Style Grouped Card Work Order Timeline === */
.ios-timeline-modal-body {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.ios-asset-info-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.ios-aic-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8px;
}

.ios-code-badge {
  font-family: monospace;
  font-size: 0.82rem;
  font-weight: 800;
  color: #1e40af;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 3px 8px;
  border-radius: 6px;
}

.ios-asset-name {
  margin: 0 0 6px;
  font-size: 1.05rem;
  font-weight: 800;
  color: #0f172a;
}

.ios-registration-location {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.82rem;
  color: #64748b;
}

.ios-timeline-page-card {
  background: #f8fafc;
  border: 1px solid #e2e8f0;
  border-radius: 16px;
  padding: 16px;
}

.ios-tl-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-bottom: 14px;
  padding-bottom: 10px;
  border-bottom: 1px dashed #cbd5e1;
}

.ios-tl-header h4 {
  margin: 0;
  font-size: 0.92rem;
  font-weight: 800;
  color: #0f172a;
}

.ios-tl-items-wrapper {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.ios-tl-card-item {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.03);
  display: flex;
  flex-direction: column;
  gap: 10px;
  transition: all 0.15s ease;
  overflow: hidden;
}

.ios-tl-card-item:hover {
  border-color: #cbd5e1;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

.ios-tl-item-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.ios-step-chip {
  font-size: 0.76rem;
  font-weight: 800;
  color: #2563eb;
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  padding: 3px 8px;
  border-radius: 6px;
}

.ios-tl-time {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.76rem;
  color: #64748b;
  background: #f8fafc;
  padding: 3px 8px;
  border-radius: 6px;
  border: 1px solid #e2e8f0;
}

.ios-tl-stacked-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.ios-pic-info-pill {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.78rem;
  color: #334155;
  background: #fff7ed;
  border: 1px solid #fed7aa;
  padding: 6px 10px;
  border-radius: 8px;
  word-break: break-word;
}

.ios-reason-box {
  background: #eff6ff;
  border: 1px solid #bfdbfe;
  border-left: 3px solid #2563eb;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 0.78rem;
  color: #1e40af;
  word-break: break-word;
}

.ios-reason-box p {
  margin: 0;
  line-height: 1.4;
}

.ios-cost-pill {
  font-size: 0.78rem;
  color: #166534;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  padding: 6px 10px;
  border-radius: 8px;
  font-weight: 600;
}

.ios-empty-tl {
  text-align: center;
  color: #64748b;
  font-size: 0.85rem;
  padding: 16px;
  background: #ffffff;
  border-radius: 10px;
  border: 1px dashed #cbd5e1;
}

@media (max-width: 640px) {
  .modal-form select, .ios-select-input {
    font-size: 0.84rem !important;
    padding: 10px 12px !important;
    max-width: 100% !important;
    box-sizing: border-box !important;
  }

  .modal-form select option, .ios-select-input option {
    font-size: 0.82rem !important;
  }

  .monthly-report-printable {
    padding: 4px 0 12px;
    width: 100%;
    box-sizing: border-box;
  }

  .report-header {
    text-align: center;
    margin-bottom: 12px;
  }

  .report-header h2 {
    font-size: 1.05rem !important;
  }

  .report-header p {
    font-size: 0.78rem !important;
  }

  .report-summary-boxes {
    grid-template-columns: repeat(2, 1fr) !important;
    gap: 8px !important;
  }

  .rbox {
    padding: 8px 10px !important;
  }

  .rbox span {
    font-size: 0.72rem !important;
  }

  .rbox strong {
    font-size: 1rem !important;
  }

  .report-table-wrapper {
    overflow-x: auto !important;
    -webkit-overflow-scrolling: touch;
    width: 100%;
    margin-bottom: 12px;
    border-radius: 8px;
    border: 1px solid #e2e8f0;
  }

  .report-table {
    min-width: 600px;
    font-size: 0.76rem !important;
  }

  .report-table th, .report-table td {
    padding: 6px 8px !important;
    white-space: nowrap;
  }

  .report-actions {
    flex-direction: column !important;
    gap: 8px !important;
  }

  .report-actions button {
    width: 100% !important;
    justify-content: center !important;
  }
}

/* Full Faded Yellow-Blue Box Highlight for Assigned Work Orders */
tr.highlight-assigned-wo, .mwo-card.highlight-assigned-wo {
  background: linear-gradient(135deg, #fffbeb 0%, #eff6ff 100%) !important;
  border: 2px solid #93c5fd !important;
  box-shadow: 0 4px 16px rgba(147, 197, 253, 0.3) !important;
}

tr.highlight-assigned-wo td {
  background: transparent !important;
}

@media print {
  body * {
    visibility: hidden;
  }
  .monthly-report-printable, .monthly-report-printable * {
    visibility: visible;
  }
  .monthly-report-printable {
    position: absolute;
    left: 0;
    top: 0;
    width: 100% !important;
    padding: 0 !important;
    margin: 0 !important;
  }
  .no-print {
    display: none !important;
  }
  .report-table-wrapper {
    overflow: visible !important;
    border: none !important;
  }
  .report-table {
    min-width: 100% !important;
    width: 100% !important;
  }
  .report-table th, .report-table td {
    white-space: normal !important;
  }
}
</style>
