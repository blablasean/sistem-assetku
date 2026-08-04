<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <p class="eyebrow">Inventaris Aset</p>
        <h1>Manajemen Aset</h1>
        <p class="subtitle">Manajemen & mutasi aset.</p>
      </div>

      <div class="header-actions">
        <button class="primary-btn btn-secondary-ios" @click="showReportModal = true" v-if="canCreateAsset">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
          <span>Laporan & Export</span>
        </button>
        <button class="primary-btn" @click="openAddModal" v-if="canCreateAsset">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
          <span>Tambah Aset Baru</span>
        </button>
      </div>
    </div>

    <div class="card-panel">
      <!-- Toolbar: Search, Filter, Sort -->
      <div class="toolbar-grid">
        <input v-model="searchQuery" placeholder="Cari nama, kode, lokasi, atau PIC aset..." class="search-input" @input="filterAssets" />
        
        <select v-model="filterStatus" @change="filterAssets" class="filter-select">
          <option value="">Semua Status</option>
          <option value="Active">Active</option>
          <option value="Maintenance">Maintenance</option>
          <option value="Damaged">Damaged</option>
          <option value="Reserved">Reserved</option>
          <option value="Retired">Retired</option>
        </select>

        <select v-model="sortBy" @change="filterAssets" class="sort-select">
          <option value="id-desc">Terbaru</option>
          <option value="name-asc">Nama Aset (A-Z)</option>
          <option value="location-asc">Lokasi / Kamar</option>
        </select>
      </div>

      <!-- Assets Table (Desktop Only) -->
      <div class="table-responsive desktop-table-only">
        <table>
          <thead>
            <tr>
              <th>Kode Aset</th>
              <th>Nama Aset</th>
              <th>Kategori</th>
              <th>Lokasi Registrasi</th>
              <th>Lokasi Mutasi Terbaru</th>
              <th>Terakhir Dipindahkan</th>
              <th>PIC</th>
              <th>Status</th>
              <th>Dokumen</th>
              <th>Aksi</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="asset in displayedAssets" :key="asset.id">
              <td><span class="code-badge">{{ asset.asset_code }}</span></td>
              <td>
                <span class="asset-name" @click="viewDetail(asset)">{{ asset.asset_name }}</span>
                <span v-if="asset.is_reserved" class="reserved-tag">Reserved</span>
              </td>
              <td>{{ asset.category || 'General' }}</td>
              <td>{{ asset.registration_location || asset.location }}</td>
              <td>{{ asset.location }}</td>
              <td>
                <span v-if="asset.location === (asset.registration_location || asset.location)" class="dash-text">-</span>
                <span v-else class="time-text">{{ formatDate(asset.last_moved_at) }}</span>
              </td>
              <td>{{ asset.pic || 'Engineering' }}</td>
              <td><StatusBadge :status="asset.status" /></td>
              <td>
                <a v-if="asset.document_url" :href="asset.document_url" target="_blank" class="doc-link">Manual</a>
                <span v-else class="no-doc">-</span>
              </td>
              <td class="actions-cell">
                <button class="icon-btn log-btn" @click.stop="openMutationTimelineModal(asset)" title="Lihat Timeline Mutasi Aset">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                  <span>Timeline</span>
                </button>
                <button class="icon-btn qr-btn" @click="openQrPrint(asset)" title="Generate & Cetak QR Code">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/></svg>
                  <span>QR</span>
                </button>
                <button class="icon-btn mut-btn" v-if="canMutate" @click="openMutationModal(asset)" title="Mutasi Lokasi Barang">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 3h5v5"/><path d="M4 20L21 3"/><path d="M21 16v5h-5"/><path d="M15 15l6 6"/><path d="M4 4l5 5"/></svg>
                  <span>Mutasi</span>
                </button>
                <button class="icon-btn edit-btn" v-if="canCreateAsset" @click="openEditModal(asset)" title="Edit Aset">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
                  <span>Edit</span>
                </button>
                <button class="icon-btn delete-btn" v-if="canDeleteAsset" @click="deleteAsset(asset)" title="Hapus Aset Permanen">
                  <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
                  <span>Hapus</span>
                </button>
              </td>
            </tr>
            <tr v-if="displayedAssets.length === 0">
              <td colspan="10" class="empty-state">Tidak ada data aset yang ditemukan.</td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Mobile Asset Cards List View (Visible on Mobile / Android Only) -->
      <div class="mobile-asset-list mobile-only">
        <div v-if="displayedAssets.length === 0" class="mobile-empty-card">
          Tidak ada data aset yang ditemukan.
        </div>
        <div v-else v-for="asset in displayedAssets" :key="'ma-' + asset.id" class="mobile-asset-card">
          <div class="mac-header">
            <span class="code-badge">{{ asset.asset_code }}</span>
            <StatusBadge :status="asset.status" />
          </div>

          <div class="mac-body" @click="viewDetail(asset)">
            <div class="mac-title-row">
              <h3 class="mac-title">{{ asset.asset_name }}</h3>
              <span v-if="asset.is_reserved" class="reserved-tag">Reserved</span>
            </div>

            <div class="mac-pills-row">
              <span class="mac-pill category-pill">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="m7.5 4.27 9 5.15"/><path d="M21 8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16Z"/><path d="m3.3 7 8.7 5 8.7-5"/><path d="M12 12v9.5"/></svg>
                <span>{{ asset.category || 'General' }}</span>
              </span>
              <span class="mac-pill location-pill">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>
                <span>{{ asset.location }}</span>
              </span>
              <span class="mac-pill pic-pill">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                <span>{{ asset.pic || 'Engineering' }}</span>
              </span>
            </div>

            <div v-if="asset.location !== (asset.registration_location || asset.location)" class="mac-mutation-info">
              <span>Registrasi: {{ asset.registration_location }}</span>
              <span class="mac-time">
                <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 4px;"><rect width="18" height="18" x="3" y="4" rx="2" ry="2"/><line x1="16" y1="2" x2="16" y2="6"/><line x1="8" y1="2" x2="8" y2="6"/><line x1="3" y1="10" x2="21" y2="10"/></svg>
                <span>{{ formatDate(asset.last_moved_at) }}</span>
              </span>
            </div>
          </div>

          <div class="mac-actions-bar">
            <button class="icon-btn log-btn" @click.stop="openMutationTimelineModal(asset)" title="Timeline Mutasi">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              <span>Timeline</span>
            </button>
            <button class="icon-btn qr-btn" @click="openQrPrint(asset)" title="Cetak QR">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect width="5" height="5" x="3" y="3" rx="1"/><rect width="5" height="5" x="16" y="3" rx="1"/><rect width="5" height="5" x="3" y="16" rx="1"/><path d="M21 16h-3a2 2 0 0 0-2 2v3"/></svg>
              <span>QR</span>
            </button>
            <button class="icon-btn mut-btn" v-if="canMutate" @click="openMutationModal(asset)" title="Mutasi">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M16 3h5v5"/><path d="M4 20L21 3"/><path d="M21 16v5h-5"/><path d="M15 15l6 6"/><path d="M4 4l5 5"/></svg>
              <span>Mutasi</span>
            </button>
            <button class="icon-btn edit-btn" v-if="canCreateAsset" @click="openEditModal(asset)" title="Edit">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 20h9"/><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"/></svg>
              <span>Edit</span>
            </button>
            <button class="icon-btn delete-btn" v-if="canDeleteAsset" @click="deleteAsset(asset)" title="Hapus">
              <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6h18"/><path d="M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6"/><path d="M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2"/><line x1="10" x2="10" y1="11" y2="17"/><line x1="14" x2="14" y1="11" y2="17"/></svg>
              <span>Hapus</span>
            </button>
          </div>
        </div>
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

    <!-- Modal Registrasi / Edit Aset -->
    <ModalDialog :show="showAssetModal" :title="isEditMode ? 'Edit Data Aset' : 'Registrasi Aset Baru'" @close="showAssetModal = false">
      <form @submit.prevent="saveAsset" class="modal-form">
        <label>
          <span>Kode Aset Unik</span>
          <input v-model="formAsset.asset_code" placeholder="Contoh: AST-RM301-AC" :disabled="isEditMode" required />
        </label>
        <label>
          <span>Nama Aset</span>
          <input v-model="formAsset.asset_name" placeholder="Contoh: AC Split Daikin 1.5 PK" required />
        </label>
        <label>
          <span>Kategori Aset</span>
          <input v-model="formAsset.category" placeholder="Contoh: HVAC / AC, Kendaraan Operasional, Audio System, dll." required />
        </label>
        <label>
          <span>Lokasi Penempatan / Kamar</span>
          <input v-model="formAsset.location" placeholder="Contoh: Kamar 301, Kitchen Dapur Utama, Chiller Room" required />
        </label>
        <label>
          <span>PIC Penanggung Jawab (Departemen)</span>
          <select v-model="formAsset.pic" required>
            <option value="Front Office">Front Office</option>
            <option value="House Keeping">House Keeping</option>
            <option value="Food Beverage Service">Food Beverage Service</option>
            <option value="Food Beverage Kitchen">Food Beverage Kitchen</option>
            <option value="Human Resource">Human Resource</option>
            <option value="Sales">Sales</option>
            <option value="Engineering">Engineering</option>
            <option value="Spa">Spa</option>
            <option value="Akunting">Akunting</option>
            <option value="IT">IT</option>
          </select>
        </label>
        <label>
          <span>Status Kondisi Aset</span>
          <select v-model="formAsset.status">
            <option value="Active">Active (Berfungsi Normal)</option>
            <option value="Maintenance">Maintenance (Dalam Perawatan)</option>
            <option value="Damaged">Damaged (Rusak)</option>
            <option value="Reserved">Reserved (Cadangan)</option>
            <option value="Retired">Retired (Dihentikan)</option>
          </select>
        </label>
        <label>
          <span>URL Dokumen / Manual Book</span>
          <input v-model="formAsset.document_url" placeholder="https://..." />
        </label>

        <button type="submit" class="submit-modal-btn">
          {{ isEditMode ? 'Simpan Perubahan' : 'Daftarkan Aset Baru' }}
        </button>
      </form>
    </ModalDialog>

    <!-- Modal Mutasi Lokasi Barang -->
    <ModalDialog :show="showMutModal" title="Mutasi Lokasi Aset" @close="showMutModal = false">
      <div v-if="selectedAssetForMut" class="mut-modal-body">
        <div class="current-info">
          <p><strong>Aset:</strong> {{ selectedAssetForMut.asset_name }} ({{ selectedAssetForMut.asset_code }})</p>
          <p><strong>Lokasi Sekarang:</strong> 📍 {{ selectedAssetForMut.location }}</p>
        </div>

        <form @submit.prevent="submitMutation" class="modal-form">
          <label>
            <span>Lokasi Baru (Ruangan / Area Operasional)</span>
            <input v-model="mutNewLocation" placeholder="Contoh: Ruang 205, Hall Utama" required />
          </label>
          <label>
            <span>PIC Penanggung Jawab Baru (Departemen)</span>
            <select v-model="mutPIC" required>
              <option value="Front Office">Front Office</option>
              <option value="House Keeping">House Keeping</option>
              <option value="Food Beverage Service">Food Beverage Service</option>
              <option value="Food Beverage Kitchen">Food Beverage Kitchen</option>
              <option value="Human Resource">Human Resource</option>
              <option value="Sales">Sales</option>
              <option value="Engineering">Engineering</option>
              <option value="Spa">Spa</option>
              <option value="Akunting">Akunting</option>
              <option value="IT">IT</option>
            </select>
          </label>
          <label>
            <span>Alasan Perpindahan / Mutasi</span>
            <input v-model="mutReason" placeholder="Contoh: Penyesuaian kebutuhan operasional / pemindahan unit" />
          </label>

          <button type="submit" class="submit-modal-btn warning-btn">Proses Mutasi Lokasi</button>
        </form>
      </div>
    </ModalDialog>

    <!-- Modal Cetak QR Code -->
    <ModalDialog :show="showQrModal" title="🖨️ Cetak & Preview QR Code Aset" @close="showQrModal = false">
      <div v-if="selectedAssetForQr" class="qr-print-card">
        <div class="printable-badge" id="qrPrintArea">
          <div class="qr-code-box">
            <img
              :src="generatedQrUrl || `https://api.qrserver.com/v1/create-qr-code/?size=180x180&data=${encodeURIComponent(selectedAssetForQr.asset_code)}`"
              :alt="selectedAssetForQr.asset_code"
              class="real-qr-image"
            />
          </div>
          <div class="qr-meta">
            <h3>AsetKu Hotel</h3>
            <p class="q-code">{{ selectedAssetForQr.asset_code }}</p>
            <p class="q-name">{{ selectedAssetForQr.asset_name }}</p>
            <p class="q-loc">📍 Lokasi Registrasi: {{ selectedAssetForQr.registration_location || selectedAssetForQr.location }}</p>
          </div>
        </div>

        <p class="qr-instruction">Tempelkan stiker QR Code ini pada unit fisik aset untuk akses cepat scan.</p>
        <div class="qr-btn-group">
          <button class="submit-modal-btn" @click="downloadQrStickerPng">📸 Unduh Stiker PNG</button>
        </div>
      </div>
    </ModalDialog>

    <!-- Modal Timeline Mutasi Aset (iOS Grouped Page Card UI) -->
    <ModalDialog :show="showMutationTimelineModal" title="Timeline & Histori Mutasi Aset" maxWidth="640px" @close="showMutationTimelineModal = false">
      <div v-if="selectedAssetForTimeline" class="ios-timeline-modal-body">
        
        <!-- Header Info Card (iOS Style) -->
        <div class="ios-asset-info-card">
          <div class="ios-aic-header">
            <span class="ios-code-badge">{{ selectedAssetForTimeline.asset_code }}</span>
            <StatusBadge :status="selectedAssetForTimeline.status" />
          </div>
          <h3 class="ios-asset-name">{{ selectedAssetForTimeline.asset_name }}</h3>
          <div class="ios-registration-location">
            <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M20 10c0 6-8 12-8 12s-8-6-8-12a8 8 0 0 1 16 0Z"/><circle cx="12" cy="10" r="3"/></svg>
            <span>Lokasi Registrasi Awal: <strong>{{ selectedAssetForTimeline.registration_location || selectedAssetForTimeline.location }}</strong></span>
          </div>
        </div>

        <!-- Single Consolidated Page Box Card for Timeline Items -->
        <div class="ios-timeline-page-card">
          <div class="ios-tl-header">
            <div style="display: flex; align-items: center; gap: 8px;">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
              <h4>Histori Mutasi Lokasi</h4>
            </div>
            <button v-if="canMutate" class="icon-btn mut-btn" @click="handleNavigateToMutasi(selectedAssetForTimeline)" title="Lakukan Mutasi Aset" style="display: inline-flex; align-items: center; gap: 4px;">
              <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/></svg>
              <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="opacity: 0.85;"><path d="M16 3h5v5"/><path d="M4 20L21 3"/><path d="M21 16v5h-5"/><path d="M15 15l6 6"/><path d="M4 4l5 5"/></svg>
              <span style="font-size: 0.78rem;">Mutasi Aset</span>
            </button>
          </div>
          
          <div v-if="isTimelineLoading" class="ios-tl-loading">
            Memuat timeline mutasi...
          </div>
          
          <div v-else class="ios-tl-items-wrapper">
            <div v-for="(log, idx) in assetTimelineLogs" :key="log.id || idx" class="ios-tl-card-item">
              <!-- Item Header: Step Badge & Time -->
              <div class="ios-tl-item-top">
                <span class="ios-step-chip">Mutasi #{{ idx + 1 }}</span>
                <span class="ios-tl-time">
                  <svg xmlns="http://www.w3.org/2000/svg" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
                  {{ formatDate(log.moved_at) }}
                </span>
              </div>

              <!-- Item Flow (Stacked Vertically for Perfect Inside Fit) -->
              <div class="ios-tl-stacked-body">
                <div class="ios-loc-stack-row">
                  <div class="ios-loc-box prev-box">
                    <span class="lbl">Dari:</span>
                    <span class="val">{{ log.previous_location || 'Registrasi Awal' }}</span>
                  </div>
                  <div class="ios-arrow-connector">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#2563eb" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="5" x2="12" y2="19"/><polyline points="19 12 12 19 5 12"/></svg>
                  </div>
                  <div class="ios-loc-box new-box">
                    <span class="lbl">Ke:</span>
                    <span class="val">📍 {{ log.new_location }}</span>
                  </div>
                </div>

                <!-- PIC / Penanggung Jawab Pill Badge (Inside Card Container) -->
                <div class="ios-pic-info-pill">
                  <svg xmlns="http://www.w3.org/2000/svg" width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 21v-2a4 4 0 0 0-4-4H9a4 4 0 0 0-4 4v2"/><circle cx="12" cy="7" r="4"/></svg>
                  <span>Penanggung Jawab (PIC): <strong>{{ log.pic || 'Engineering' }}</strong></span>
                </div>

                <!-- Alasan Mutasi Note Box -->
                <div class="ios-reason-box" v-if="log.reason">
                  <p><strong>Alasan / Catatan:</strong> {{ log.reason }}</p>
                </div>
              </div>
            </div>

            <div v-if="assetTimelineLogs.length === 0" class="ios-empty-tl">
              Belum ada riwayat mutasi perpindahan lokasi untuk aset ini.
            </div>
          </div>
        </div>

      </div>
    </ModalDialog>

    <!-- Laporan Aset Modal -->
    <ModalDialog :show="showReportModal" title="Laporan Inventaris Aset" maxWidth="960px" @close="showReportModal = false">
      <div class="monthly-report-printable" id="printableAssetReport">
        <div class="report-header">
          <h2>LAPORAN INVENTARIS ASET</h2>
          <p class="report-sub">Sistem AsetKu — Dicetak: {{ reportDate }}</p>
          <hr class="report-divider" />
        </div>

        <div class="report-summary-boxes">
          <div class="rbox">
            <span>Total Aset</span>
            <strong>{{ assets.length }} Unit</strong>
          </div>
          <div class="rbox success">
            <span>Status Active</span>
            <strong>{{ countAssetStatus('Active') }} Unit</strong>
          </div>
          <div class="rbox warning">
            <span>Dalam Maintenance</span>
            <strong>{{ countAssetStatus('Maintenance') }} Unit</strong>
          </div>
          <div class="rbox danger">
            <span>Rusak / Retired</span>
            <strong>{{ countAssetStatus('Damaged') + countAssetStatus('Retired') }} Unit</strong>
          </div>
        </div>

        <div class="report-table-wrapper">
          <table class="report-table">
            <thead>
              <tr>
                <th>Kode Aset</th>
                <th>Nama Aset</th>
                <th>Kategori</th>
                <th>Lokasi Registrasi</th>
                <th>PIC</th>
                <th>Status</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="asset in assets" :key="asset.id">
                <td>{{ asset.asset_code }}</td>
                <td>{{ asset.asset_name }}</td>
                <td>{{ asset.category }}</td>
                <td>{{ asset.registration_location || asset.location }}</td>
                <td>{{ asset.pic }}</td>
                <td>{{ asset.status }}</td>
              </tr>
              <tr v-if="assets.length === 0">
                <td colspan="6" class="empty-state">Belum ada data aset terdaftar.</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <div class="report-actions no-print">
        <button class="excel-btn" @click="exportAssetToExcel">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z"/><polyline points="14 2 14 8 20 8"/><line x1="8" y1="13" x2="16" y2="13"/><line x1="8" y1="17" x2="16" y2="17"/></svg>
          <span>Export ke Excel (.xlsx)</span>
        </button>
        <button class="print-btn" @click="printAssetReport">
          <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="btn-icon"><polyline points="6 9 6 2 18 2 18 9"/><path d="M6 18H4a2 2 0 0 1-2-2v-5a2 2 0 0 1 2-2h16a2 2 0 0 1 2 2v5a2 2 0 0 1-2 2h-2"/><rect width="12" height="8" x="6" y="14"/></svg>
          <span>Cetak Dokumen Laporan (PDF / Print)</span>
        </button>
      </div>
    </ModalDialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import QRCode from 'qrcode'
import StatusBadge from '../components/StatusBadge.vue'
import ModalDialog from '../components/ModalDialog.vue'
import api from '../api'

const userRole = ref(sessionStorage.getItem('user_role') || 'external')
const canCreateAsset = computed(() => ['hod', 'management', 'supervisor', 'admin'].includes(userRole.value.toLowerCase()))
const canDeleteAsset = computed(() => ['hod', 'management', 'supervisor', 'admin'].includes(userRole.value.toLowerCase()))
const canMutate = computed(() => ['hod', 'management', 'supervisor', 'admin'].includes(userRole.value.toLowerCase()))

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

const assets = ref([])
const isLoading = ref(false)

const showReportModal = ref(false)
const reportDate = ref(new Date().toLocaleDateString('id-ID', { day: '2-digit', month: 'long', year: 'numeric' }))

function countAssetStatus(status) {
  return assets.value.filter(a => a.status === status).length
}

import { exportToExcel, triggerPrint } from '../utils/exportUtils'

function exportAssetToExcel() {
  const fileName = `Laporan_Inventaris_Aset_${reportDate.value.replace(/\s+/g, '_')}.xls`
  const headers = ['Kode Aset', 'Nama Aset', 'Kategori', 'Lokasi Registrasi', 'PIC', 'Status']
  const rows = assets.value.map(a => [
    a.asset_code || '',
    a.asset_name || '',
    a.category || '',
    a.registration_location || a.location || '',
    a.pic || '',
    a.status || ''
  ])
  exportToExcel(fileName, headers, rows)
}

function printAssetReport() {
  triggerPrint()
}

const searchQuery = ref('')
const filterStatus = ref('')
const sortBy = ref('id-desc')

const showAssetModal = ref(false)
const isEditMode = ref(false)
const formAsset = ref({ id: 0, asset_code: '', asset_name: '', category: '', location: '', pic: '', status: 'Active', document_url: '' })

const showMutModal = ref(false)
const selectedAssetForMut = ref(null)
const mutNewLocation = ref('')
const mutPIC = ref('')
const mutReason = ref('')

const showQrModal = ref(false)
const selectedAssetForQr = ref(null)

const showMutationTimelineModal = ref(false)
const selectedAssetForTimeline = ref(null)
const assetTimelineLogs = ref([])
const isTimelineLoading = ref(false)
const newMutLocation = ref('')
const newMutPic = ref('')
const newMutReason = ref('')
const isSubmittingMut = ref(false)

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

function handleNavigateToMutasi(asset) {
  showMutationTimelineModal.value = false
  openMutationModal(asset)
}

async function openMutationTimelineModal(asset) {
  if (!asset) return
  selectedAssetForTimeline.value = asset
  showMutationTimelineModal.value = true
  isTimelineLoading.value = true

  const defaultLogs = [
    {
      id: 1,
      asset_code: asset.asset_code,
      previous_location: '-',
      new_location: asset.registration_location || asset.location,
      pic: asset.pic || 'Engineering',
      reason: 'Registrasi awal aset terdaftar di sistem',
      moved_at: asset.created_at || new Date().toISOString()
    }
  ]

  try {
    const res = await api.get(`/mutations/timeline?asset_code=${encodeURIComponent(asset.asset_code)}`)
    const logsData = res.data?.data || res.data
    if (Array.isArray(logsData) && logsData.length > 0) {
      assetTimelineLogs.value = logsData
    } else {
      assetTimelineLogs.value = defaultLogs
    }
  } catch (e) {
    console.error('Fetch asset mutation timeline error:', e)
    assetTimelineLogs.value = defaultLogs
  } finally {
    isTimelineLoading.value = false
  }
}

async function submitAddMutationNote() {
  if (!selectedAssetForTimeline.value || !newMutLocation.value) return
  isSubmittingMut.value = true
  try {
    const payload = {
      asset_code: selectedAssetForTimeline.value.asset_code,
      new_location: newMutLocation.value,
      pic: newMutPic.value || selectedAssetForTimeline.value.pic,
      reason: newMutReason.value
    }
    await api.post('/assets/mutate', payload)
    notify(`Mutasi lokasi aset ${selectedAssetForTimeline.value.asset_code} ke ${newMutLocation.value} berhasil dicatat!`, 'success')
    newMutLocation.value = ''
    newMutReason.value = ''

    const res = await api.get(`/assets/mutation-timeline?asset_code=${selectedAssetForTimeline.value.asset_code}`)
    const logsData = res.data?.data || res.data
    if (Array.isArray(logsData) && logsData.length > 0) {
      assetTimelineLogs.value = logsData
    }
    await fetchAssets()
  } catch (e) {
    console.error('Failed to record mutation:', e)
    notify(e.response?.data?.error || 'Gagal mencatat mutasi aset.', 'error')
  } finally {
    isSubmittingMut.value = false
  }
}

const displayedAssets = computed(() => {
  let list = [...assets.value]
  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(a => 
      (a.asset_name && a.asset_name.toLowerCase().includes(q)) || 
      (a.asset_code && a.asset_code.toLowerCase().includes(q)) || 
      (a.location && a.location.toLowerCase().includes(q)) ||
      (a.pic && a.pic.toLowerCase().includes(q))
    )
  }
  if (filterStatus.value) {
    list = list.filter(a => a.status === filterStatus.value)
  }
  if (sortBy.value === 'name-asc') {
    list.sort((a, b) => (a.asset_name || '').localeCompare(b.asset_name || ''))
  } else if (sortBy.value === 'location-asc') {
    list.sort((a, b) => (a.location || '').localeCompare(b.location || ''))
  } else {
    list.sort((a, b) => b.id - a.id)
  }
  return list
})

async function fetchAssets() {
  isLoading.value = true
  try {
    const res = await api.get('/assets')
    if (res.data?.data && Array.isArray(res.data.data)) {
      assets.value = res.data.data
    } else {
      assets.value = []
    }
  } catch (e) {
    console.error('Failed to fetch assets from backend:', e)
    assets.value = []
  } finally {
    isLoading.value = false
  }
}

function generateUniqueAssetCode() {
  let code = ''
  let isDuplicate = true
  let attempts = 0
  const existingCodes = new Set((assets.value || []).map(a => a.asset_code))
  while (isDuplicate && attempts < 1000) {
    attempts++
    const randNum = Math.floor(100 + Math.random() * 900)
    code = `AST-RM${randNum}-UNIT`
    isDuplicate = existingCodes.has(code)
  }
  return code
}

function openAddModal() {
  isEditMode.value = false
  formAsset.value = { 
    id: 0, 
    asset_code: generateUniqueAssetCode(), 
    asset_name: '', 
    category: '', 
    location: '', 
    pic: 'Front Office', 
    status: 'Active', 
    document_url: '' 
  }
  showAssetModal.value = true
}

function openEditModal(asset) {
  isEditMode.value = true
  formAsset.value = { 
    ...asset,
    pic: asset.pic || 'Front Office'
  }
  showAssetModal.value = true
}

async function saveAsset() {
  try {
    if (isEditMode.value) {
      await api.put(`/assets?id=${formAsset.value.id}`, formAsset.value)
      notify('Perubahan data aset berhasil disimpan!', 'success')
    } else {
      await api.post('/assets', formAsset.value)
      notify('Aset baru berhasil didaftarkan!', 'success')
    }
    showAssetModal.value = false
    await fetchAssets()
  } catch (e) {
    notify('Gagal menyimpan aset: ' + (e.response?.data?.message || e.message), 'error')
  }
}

async function deleteAsset(asset) {
  try {
    await api.post('/assets/delete', { asset_id: asset.id })
    assets.value = assets.value.filter(a => a.id !== asset.id)
    notify(`Aset "${asset.asset_name}" berhasil dihapus permanen!`, 'success')
    await fetchAssets()
  } catch (e) {
    notify('Gagal menghapus aset: ' + (e.response?.data?.message || e.message), 'error')
  }
}

function openMutationModal(asset) {
  selectedAssetForMut.value = asset
  mutNewLocation.value = ''
  mutPIC.value = asset.pic || 'Front Office'
  mutReason.value = ''
  showMutModal.value = true
}

async function submitMutation() {
  if (!selectedAssetForMut.value) return
  try {
    await api.post('/mutations', {
      asset_id: selectedAssetForMut.value.id,
      new_location: mutNewLocation.value,
      pic: mutPIC.value,
      new_pic: mutPIC.value,
      reason: mutReason.value
    })
    showMutModal.value = false
    notify(`Mutasi aset "${selectedAssetForMut.value.asset_name}" ke ${mutNewLocation.value} berhasil dicatat!`, 'success')
    await fetchAssets()
  } catch (e) {
    notify('Gagal mencatat mutasi: ' + (e.response?.data?.message || e.message), 'error')
  }
}

const generatedQrUrl = ref('')

async function openQrPrint(asset) {
  selectedAssetForQr.value = asset
  try {
    generatedQrUrl.value = await QRCode.toDataURL(asset.asset_code || 'AST-UNKNOWN', {
      width: 260,
      margin: 2,
      color: {
        dark: '#0f172a',
        light: '#ffffff'
      }
    })
  } catch (err) {
    console.error('QRCode.toDataURL error:', err)
    generatedQrUrl.value = `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(asset.asset_code)}`
  }
  showQrModal.value = true
}

function printQrCard() {
  window.print()
}

async function downloadQrStickerPng() {
  if (!selectedAssetForQr.value) return
  const asset = selectedAssetForQr.value

  const canvas = document.createElement('canvas')
  canvas.width = 420
  canvas.height = 560
  const ctx = canvas.getContext('2d')

  // Background card
  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, canvas.width, canvas.height)
  
  // Outer Border
  ctx.strokeStyle = '#0f172a'
  ctx.lineWidth = 6
  ctx.strokeRect(12, 12, canvas.width - 24, canvas.height - 24)

  // Header Banner
  ctx.fillStyle = '#0f172a'
  ctx.fillRect(12, 12, canvas.width - 24, 70)

  ctx.fillStyle = '#f59e0b'
  ctx.font = 'bold 22px sans-serif'
  ctx.textAlign = 'center'
  ctx.fillText('ASETKU ASSET STICKER', canvas.width / 2, 56)

  // Load and Draw QR Code Image
  const qrImg = new Image()
  qrImg.crossOrigin = 'Anonymous'
  qrImg.src = generatedQrUrl.value || `https://api.qrserver.com/v1/create-qr-code/?size=220x220&data=${encodeURIComponent(asset.asset_code)}`
  
  await new Promise((resolve) => {
    qrImg.onload = resolve
    qrImg.onerror = resolve
  })

  ctx.drawImage(qrImg, (canvas.width - 220) / 2, 105, 220, 220)

  // Asset Details Text
  ctx.textAlign = 'center'

  // Code
  ctx.fillStyle = '#d97706'
  ctx.font = 'bold 22px monospace'
  ctx.fillText(asset.asset_code || '', canvas.width / 2, 360)

  // Name
  ctx.fillStyle = '#0f172a'
  ctx.font = 'bold 20px sans-serif'
  const nameText = (asset.asset_name || '').length > 26 ? asset.asset_name.substring(0, 26) + '...' : asset.asset_name
  ctx.fillText(nameText, canvas.width / 2, 400)

  // Location & PIC
  ctx.fillStyle = '#475569'
  ctx.font = '16px sans-serif'
  ctx.fillText(`📍 ${asset.location || 'Ruangan / Area'}`, canvas.width / 2, 440)
  ctx.fillText(`👤 PIC: ${asset.pic || 'Engineering'}`, canvas.width / 2, 475)

  // Footer Subtext
  ctx.fillStyle = '#94a3b8'
  ctx.font = '12px sans-serif'
  ctx.fillText('PROPERTY OF ASSET MANAGEMENT SYSTEM', canvas.width / 2, 520)

  // Trigger PNG download
  const imageURI = canvas.toDataURL('image/png')
  const link = document.createElement('a')
  link.download = `Stiker_QR_${asset.asset_code}.png`
  link.href = imageURI
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  notify(`Stiker QR ${asset.asset_code} berhasil diunduh sebagai gambar PNG!`, 'success')
}

function viewDetail(asset) {
  alert(`Detail Aset Hotel:\n\nKode: ${asset.asset_code}\nNama: ${asset.asset_name}\nLokasi: ${asset.location}\nPIC: ${asset.pic}\nStatus: ${asset.status}`)
}

onMounted(() => {
  fetchAssets()
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

.header-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  height: 37.6px;
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
  box-sizing: border-box !important;
  white-space: nowrap !important;
  margin: 0 !important;
}

.primary-btn:hover {
  background: #0062cc !important;
  border-color: #0062cc !important;
  transform: translateY(-1px);
}

.card-panel {
  background: #ffffff;
  border-radius: 20px;
  padding: 24px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.05);
  border: 1px solid #e2e8f0;
}

.toolbar-grid {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.search-input {
  flex: 1;
  min-width: 260px;
  padding: 10px 16px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
}

.filter-select, .sort-select {
  padding: 10px 14px;
  border: 1px solid #cbd5e1;
  border-radius: 12px;
  background: white;
}

.table-responsive {
  overflow-x: auto;
}

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

.code-badge {
  background: #e2e8f0;
  color: #0f172a;
  padding: 3px 8px;
  border-radius: 6px;
  font-family: monospace;
  font-weight: 700;
}

.asset-name {
  font-weight: 600;
  color: #2563eb;
  cursor: pointer;
}

.reserved-tag {
  margin-left: 6px;
  background: #ffedd5;
  color: #c2410c;
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 4px;
  font-weight: 700;
}

.doc-link {
  color: #0284c7;
  text-decoration: none;
  font-weight: 600;
}

.no-doc { color: #94a3b8; }

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

.log-btn {
  background: #eff6ff;
  color: #2563eb;
  border-color: #bfdbfe;
}

.log-btn:hover {
  background: #dbeafe;
  border-color: #93c5fd;
}

.qr-btn {
  background: #f1f5f9;
  color: #0f172a;
  border-color: #cbd5e1;
}

.qr-btn:hover {
  background: #e2e8f0;
  border-color: #94a3b8;
}

.mut-btn {
  background: #fff7ed;
  color: #c2410c;
  border-color: #ffedd5;
}

.mut-btn:hover {
  background: #ffedd5;
  border-color: #fed7aa;
}

.edit-btn {
  background: #f8fafc;
  color: #475569;
  border-color: #cbd5e1;
  padding: 0 8px;
}

.edit-btn:hover {
  background: #f1f5f9;
  color: #0f172a;
}

.delete-btn {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
  padding: 0 8px;
}

.delete-btn:hover {
  background: #fee2e2;
  border-color: #fca5a5;
}

.primary-btn:hover {
  background: #1e293b;
  border-color: #1e293b;
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

.modal-form input, .modal-form select {
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

.modal-form input:focus, .modal-form select:focus {
  border-color: #2563eb;
  box-shadow: 0 0 0 3px rgba(37, 99, 235, 0.15);
}

.submit-modal-btn {
  background: #0f172a;
  color: #ffffff;
  border: 1px solid #0f172a;
  padding: 13px;
  border-radius: 4px !important;
  font-size: 0.95rem;
  font-weight: 700;
  cursor: pointer;
  margin-top: 8px;
  width: 100%;
  transition: all 0.15s ease;
}

.submit-modal-btn:hover {
  background: #1e293b;
}

.warning-btn { background: #d97706; }

.current-info {
  background: #f8fafc;
  padding: 12px;
  border-radius: 10px;
  margin-bottom: 16px;
  font-size: 0.9rem;
}

.qr-print-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
}

.printable-badge {
  border: 2px solid #0f172a;
  border-radius: 16px;
  padding: 20px;
  display: flex;
  align-items: center;
  gap: 20px;
  background: #ffffff;
  width: 100%;
  max-width: 380px;
}

.qr-code-box {
  width: 110px;
  height: 110px;
  background: #ffffff;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid #cbd5e1;
  padding: 4px;
}

.real-qr-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.qr-meta h3 { margin: 0; font-size: 1.1rem; color: #0f172a; }
.q-code { margin: 4px 0; font-family: monospace; font-weight: 800; color: #2563eb; }
.q-name { margin: 2px 0; font-weight: 600; font-size: 0.9rem; }
.q-loc { margin: 4px 0 0; color: #64748b; font-size: 0.85rem; }

.qr-instruction {
  font-size: 0.85rem;
  color: #64748b;
  text-align: center;
}

.qr-btn-group {
  display: flex;
  gap: 10px;
  width: 100%;
  justify-content: center;
}

.secondary-btn {
  background: #f1f5f9;
  color: #0f172a;
  border: 1px solid #cbd5e1;
}

/* === Report Modal CSS === */
.monthly-report-printable { padding: 4px 0 16px; }
.report-header { text-align: center; margin-bottom: 20px; }
.report-header h2 { font-size: 1.1rem; font-weight: 800; color: #0f172a; margin: 0 0 4px; letter-spacing: -0.01em; }
.report-sub { color: #64748b; font-size: 0.85rem; margin: 0; }
.report-divider { border: none; border-top: 2px solid #e2e8f0; margin: 16px 0 0; }
.report-summary-boxes { display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-bottom: 20px; }
.rbox { background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 12px 14px; }
.rbox span { display: block; font-size: 0.75rem; color: #64748b; font-weight: 600; margin-bottom: 4px; }
.rbox strong { font-size: 1.1rem; font-weight: 800; color: #0f172a; }
.rbox.success strong { color: #16a34a; }
.rbox.warning strong { color: #d97706; }
.rbox.danger strong { color: #dc2626; }
.report-table-wrapper { border: 1px solid #e2e8f0; border-radius: 10px; overflow: hidden; }
.report-table { width: 100%; border-collapse: collapse; font-size: 0.82rem; }
.report-table th { background: #0f172a; color: #fff; padding: 9px 12px; text-align: left; font-weight: 700; white-space: nowrap; }
.report-table td { padding: 8px 12px; border-bottom: 1px solid #f1f5f9; color: #334155; }
.report-table tbody tr:last-child td { border-bottom: none; }
.report-table tbody tr:hover td { background: #f8fafc; }
.empty-state { text-align: center; color: #94a3b8; font-style: italic; padding: 20px !important; }
.report-actions { display: flex; gap: 12px; padding-top: 16px; justify-content: flex-end; border-top: 1px solid #e2e8f0; margin-top: 16px; }
.excel-btn, .print-btn {
  display: inline-flex; align-items: center; justify-content: center; gap: 8px;
  height: 42px; padding: 0 18px; border-radius: 10px; font-size: 0.88rem;
  font-weight: 700; cursor: pointer; border: none; transition: all 0.15s ease;
  line-height: 1; white-space: nowrap;
}
.excel-btn { background: #16a34a; color: #fff; }
.excel-btn:hover { background: #15803d; }
.print-btn { background: #007aff; color: #fff; }
.print-btn:hover { background: #0062cc; }
.btn-icon { flex-shrink: 0; }
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
@media print {
  .no-print { display: none !important; }
}

/* === Desktop vs Mobile Display Toggle === */
.mobile-only {
  display: none !important;
}

/* === Mobile Responsive CSS (Android & iOS) === */
@media (max-width: 640px) {
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
  .header-actions { width: 100%; display: flex; flex-wrap: wrap; gap: 8px; }
  .header-actions .primary-btn { flex: 1; min-width: 130px; justify-content: center; height: 40px !important; font-size: 0.82rem !important; }
  .toolbar-grid { flex-direction: column; gap: 10px; }
  .search-input, .select-input, .filter-select, .sort-select { width: 100%; }
  .card-panel { padding: 16px !important; border-radius: 14px !important; }
  .report-summary-boxes { grid-template-columns: repeat(2, 1fr); gap: 8px; }
  .report-actions { flex-direction: column; gap: 8px; }
  .excel-btn, .print-btn { width: 100%; justify-content: center; }

  .mobile-asset-list {
    width: 100%;
  }

  .mobile-asset-card {
    background: #ffffff;
    border: 1px solid #e2e8f0;
    border-radius: 14px !important;
    padding: 14px 16px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.04);
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .mac-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f1f5f9;
  }

  .mac-body {
    display: flex;
    flex-direction: column;
    gap: 8px;
    cursor: pointer;
  }

  .mac-title-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
  }

  .mac-title {
    margin: 0;
    font-size: 0.92rem;
    font-weight: 800;
    color: #0f172a;
    line-height: 1.3;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .mac-pills-row {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
  }

  .mac-pill {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    font-size: 0.76rem;
    font-weight: 700;
    padding: 4px 9px;
    border-radius: 6px;
    background: #f8fafc;
    border: 1px solid #e2e8f0;
    color: #475569;
    line-height: 1;
  }

  .mac-pill span {
    max-width: 110px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .mac-pill svg {
    flex-shrink: 0;
  }

  .category-pill { background: #eff6ff; border-color: #bfdbfe; color: #1e40af; }
  .location-pill { background: #f0fdf4; border-color: #bbf7d0; color: #166534; }
  .pic-pill { background: #fff7ed; border-color: #fed7aa; color: #c2410c; }

  .mac-mutation-info {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    font-size: 0.76rem;
    color: #64748b;
    background: #f8fafc;
    padding: 6px 10px;
    border-radius: 8px;
  }

  .mac-time {
    display: inline-flex;
    align-items: center;
  }

  .mac-actions-bar {
    display: flex;
    align-items: center;
    gap: 6px;
    flex-wrap: wrap;
    padding-top: 8px;
    border-top: 1px solid #f1f5f9;
  }

  .mac-actions-bar .icon-btn {
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

/* === iOS Style Grouped Card Mutation Timeline === */
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

.ios-loc-stack-row {
  display: flex;
  flex-direction: column;
  gap: 4px;
  background: #f8fafc;
  border: 1px solid #f1f5f9;
  border-radius: 10px;
  padding: 10px 12px;
}

.ios-loc-box {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.82rem;
  line-height: 1.3;
}

.ios-loc-box .lbl {
  font-weight: 700;
  color: #64748b;
  min-width: 36px;
}

.ios-loc-box .val {
  font-weight: 700;
  color: #0f172a;
  word-break: break-word;
}

.ios-loc-box.new-box .val {
  color: #15803d;
}

.ios-arrow-connector {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2px 0;
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

.ios-empty-tl {
  text-align: center;
  color: #64748b;
  font-size: 0.85rem;
  padding: 16px;
  background: #ffffff;
  border-radius: 10px;
  border: 1px dashed #cbd5e1;
}

.ios-add-mutation-card {
  background: #ffffff;
  border: 1px solid #e2e8f0;
  border-radius: 14px;
  padding: 14px 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
}

.ios-amc-title {
  margin: 0 0 12px;
  font-size: 0.92rem;
  font-weight: 800;
  color: #0f172a;
}

.ios-amc-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
}

@media (max-width: 640px) {
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
    min-width: 580px;
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
