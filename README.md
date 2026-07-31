# 🏨🏢🛠️ Sistem AsetKu Hotel — Asset & Work Order Management System

**Sistem AsetKu** adalah sistem informasi manajemen aset hotel, perawatan rutin (*Preventive Maintenance*), pelaporan kerusakan (*Work Order*), mutasi lokasi aset, dan log aktivitas operasional real-time berbasis web yang dirancang khusus untuk memenuhi standar efisiensi industri perhotelan dan tim *Engineering*.

Aplikasi ini mengusung desain **Apple iOS Modern Aesthetic** dengan antarmuka yang bersih, sangat responsif di perangkat Smartphone (Android/iOS) dan Laptop/Desktop, serta dilengkapi proteksi keamanan tingkat tinggi seperti **Single Active Session per Account dengan Real-Time Heartbeat Architecture**.

---

## 📚 Dokumentasi Terkait & Panduan Lengkap
- 📖 **[USER_GUIDE.md](file:///c:/Users/Sean/Documents/KERJA/PROJECT/sistem-assetku/USER_GUIDE.md)**: Panduan Pengguna & Pengoperasian Seluruh Fitur Sistem (Penggunaan Kamera HP, Hak Akses Akun, Work Order, PM, Cetak Stiker QR).
- 🛠️ **[SETUP.md](file:///c:/Users/Sean/Documents/KERJA/PROJECT/sistem-assetku/SETUP.md)**: Panduan Instalasi Lokal & Konfigurasi Database MySQL.
- 🌐 **[DEPLOYMENT_GUIDE.md](file:///c:/Users/Sean/Documents/KERJA/PROJECT/sistem-assetku/deploy/DEPLOYMENT_GUIDE.md)**: Panduan Deploy HTTPS Server Production (Port 443 Nginx/Caddy & IP `172.17.10.109`).

---

## 🌟 Fitur Utama Sistem

### 1. 🔐 Autentikasi & Keamanan Sesi Real-Time (Anti Lockout & Multi-Device Protection)
* **Single Active Session per Account (Anti Double-Login)**: Mencegah akun yang sama digunakan secara bersamaan oleh lebih dari satu perangkat/browser. Apabila akun sedang aktif di Perangkat A, upaya login di Perangkat B akan ditolak secara otomatis oleh server backend.
* **Real-Time Heartbeat Architecture (Setiap 12 Detik)**: Frontend secara berkala mengirimkan sinyal keleluasaan ke backend untuk menjaga timestamp `last_seen_at` tetap aktif selama tab terbuka.
* **Instant Tab Release (`pagehide` + `sendBeacon`)**: Saat tab atau peramban ditutup, aplikasi otomatis mengirimkan sinyal pelepasan sesi ke server agar akun bebas dan siap digunakan kembali.
* **Account Lockout Protection (30 Detik Toleransi)**: Apabila peramban ditutup secara mendadak (misal mati listrik), server secara otomatis melepaskan kunci akun dalam 30 detik tanpa pernah mengunci pengguna.
* **Terisolasi Per-Tab (`sessionStorage`)**: Menyimpan token autentikasi secara terisolasi per-tab browser.

### 2. 📋 Log Aktivitas Sistem & Proteksi Khusus Admin
* **Restriksi Khusus Admin**: Log Aktivitas Sistem (Audit Trail) hanya dapat diakses & dilihat oleh role `admin` (baik dari sisi frontend UI maupun proteksi API backend).
* **Filter Rentang Tanggal**: Dilengkapi filter tanggal presisi dengan tombol pintas cepat (*Hari Ini, 7 Hari, 30 Hari, Bulan Ini*) serta pilihan tanggal custom.
* **Tampilan Ringkas & Pagination**: Tabel dengan limit bawaan 5 baris, scrollable container (`max-height: 340px`), dan pilihan per-halaman (5, 10, 25, 50, Semua).

### 3. 📱 Pemindai & Pembuat QR Code Real-Time
* **Kamera Live & File Upload QR**: Pemindaian otomatis via kamera WebCam live maupun **Upload Foto Stiker QR** dengan dekode piksel JavaScript `jsQR` (0% CPU lag).
* **Cetak & Download Stiker QR (PNG)**: Menggenerasi stiker QR Code aset secara murni di peramban (*client-side canvas*) yang dapat langsung diunduh sebagai gambar PNG atau dicetak.

### 4. 📦 Manajemen Inventaris Aset & Mutasi Lokasi
* **Pencatatan Aset Terpadu**: Mengelola data aset lengkap dengan Kode Aset, Nama, Kategori, Lokasi Registrasi, PIC, dan Status (*Active, Maintenance, Damaged, Disposed*).
* **Riwayat Mutasi Lokasi**: Pencatatan resmi perpindahan fisik aset (`asset_mutation_timelines`) dari lokasi awal ke lokasi baru beserta PIC.

### 5. 🛠️ Pelaporan Kerusakan (Work Order) & Timeline Progress
* **Pelaporan Mudah & Fleksibel**: Pengajuan tiket kerusakan dengan tingkat prioritas (*Low, Normal, High, Emergency*) dan pencatatan lokasi kamar/ruangan.
* **Timeline Progres Perbaikan (`work_order_logs`)**: Pencatatan kronologi tindakan teknisi (*Action Taken*), biaya perbaikan (*Cost*), dan update status tiket hingga tahap **Finish / Closed**.

### 6. 📅 Kalender Preventive Maintenance (PM)
* **Jadwal Rutin Otomatis**: Manajemen jadwal perawatan berkala (*Daily, Weekly, Monthly, Quarterly, Yearly*) dengan tampilan kalender interaktif dan penyelesaian *Checklist*.

### 7. 📄 Ekspor Laporan & Cetak PDF Terpadu (`exportUtils.js`)
* **Ekspor Excel (.xls)** & **Cetak Laporan PDF**: Modul utilitas independen `exportUtils.js` yang memungkinkan pengunduhan data aset, work order, dan maintenance dalam format Excel serta pembuatan dokumen siap cetak.

---

## 👥 Matriks Hak Akses (Role-Based Access Control)

| Role / Jabatan | Hak Akses Utama |
| :--- | :--- |
| **👑 Administrator (`admin`)** | Akses penuh atas seluruh aset, Work Order, penugasan, jadwal PM, user management, dan audit log sistem. |
| **⭐ HOD Engineer (`hod`)** | Registrasi & Edit Aset, Penugasan Teknisi, Penutupan Tiket (*Finish*), Hapus WO/Aset/PM, Cetak Laporan & Export Excel. |
| **👔 Supervisor (`management`)** | Monitoring Aset, Penugasan Teknisi, Update Status, Mutasi Lokasi, Hapus Data, Cetak Laporan. |
| **🛠️ Staff Engineer (`engineer`)** | Update progres pengerjaan tiket, isi tindakan perbaikan, biaya, dan penyelesaian *Checklist PM*. |
| **🛎️ Staff Hotel (`external` / Dept)** | Pengajuan Tiket Kerusakan Aset berdasarkan Lokasi/Kamar dan pemantauan status perbaikan real-time. |

---

## 🛠️ Arsitektur & Tech Stack

* **Frontend**: Vue 3 (Composition API) + Vite + Vanilla CSS iOS Design System + `jsQR`
* **Backend**: Go (Golang) + `net/http` + GORM ORM + JWT Authentication + Bcrypt
* **Database**: MySQL / MariaDB (`db_sistemasetku`)
* **Keamanan**: Single Active Session dengan 30s Heartbeat Timeout, Role-Based Access Control (RBAC), Interceptor 401/403.

---

## 🚀 Panduan Jalankan Sistem & Deploy Production

### A. Pengujian Lokal / Development
```bash
# 1. Jalankan Backend Go (Port 8080)
cd backend
go run main.go

# 2. Jalankan Frontend Vue Dev Server (Port 3000)
cd frontend
npm run dev -- --host
```

---

### B. Deployment Server Production (Rekomendasi Resmi)

#### 1. Build & Run Backend (Golang Binary)
```bash
cd backend

# Compile menjadi file binary executable tunggal
go build -o server.exe main.go

# Jalankan server executable
.\server.exe
```
*Atau gunakan **NSSM** (Windows Service) atau **Systemd/PM2** (Linux) agar `server.exe` menyala otomatis saat server restart.*

#### 2. Build & Serve Frontend (Production Build)
```bash
cd frontend

# Kompilasi aset produksi ke folder dist/
npm run build

# Jalankan server produksi
npx serve -s dist -l 3000
```
*Modul `api.js` secara otomatis memetakan permintaan API dari port 3000 langsung ke port backend `8080` tanpa perlu konfigurasi proxy manual.*

---

## 📂 Struktur Repositori

```text
sistem-assetku/
├── backend/                  # RESTful API Go (Golang)
│   ├── config/               # Database Connection & Auto-Migration
│   ├── controllers/          # Business Logic (Auth, Asset, WO, Mutation, PM, User)
│   ├── db/                   # Updated Schema SQL (schema.sql)
│   ├── middlewares/          # JWT & Single Active Session Middleware
│   ├── models/               # Struct GORM (User, Asset, WorkOrder, PM, ActivityLog)
│   ├── routes/               # HTTP Route Handler Definitions
│   ├── server.exe            # Binary Executable Production Server
│   └── main.go               # Entrypoint Backend Go
│
├── frontend/                 # Single Page Application (Vue 3 + Vite)
│   ├── dist/                 # Production Build Assets Output
│   ├── src/
│   │   ├── components/       # HeaderNavbar, QrScannerModal, StatusBadge, ModalDialog
│   │   ├── utils/            # exportUtils.js, auth.js
│   │   ├── views/            # Dashboard, AssetManagement, WorkOrder, Maintenance, ActivityLog, UserManagement
│   │   ├── router/           # Vue Router Navigation Guards
│   │   ├── App.vue           # Root Component with Heartbeat & PageHide Tab Release
│   │   └── api.js            # Axios Interceptors & Dynamic BaseURL Port Auto-Mapping
│   └── vite.config.js        # Vite Config with HTML Refresh Bypass Proxy
│
└── README.md                 # Dokumentasi Resmi Proyek
```
