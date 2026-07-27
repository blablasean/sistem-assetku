# 🏨🏢🛠️ Sistem AsetKu Hotel — Asset & Work Order Management System

**Sistem AsetKu** adalah sistem informasi manajemen aset hotel, perawatan rutin (*Preventive Maintenance*), pelaporan kerusakan (*Work Order*), mutasi lokasi aset, dan log aktivitas operasional real-time berbasis web yang dirancang khusus untuk memenuhi standar efisiensi industri perhotelan dan tim *Engineering*.

Aplikasi ini mengusung desain **Apple iOS Modern Aesthetic** dengan antarmuka yang bersih, sangat responsif di perangkat Smartphone (Android/iOS) dan Laptop/Desktop, serta dilengkapi proteksi keamanan tingkat tinggi seperti **Single Active Session per Account** (Anti Double-Login).

---

## 🌟 Fitur Utama Sistem

### 1. 🔐 Autentikasi & Keamanan Sesi Tingkat Tinggi
* **Single Active Session per Account (Anti Double-Login)**: Mencegah akun yang sama digunakan secara bersamaan oleh lebih dari satu perangkat. Apabila akun sedang aktif di Perangkat A, upaya login di Perangkat B akan ditolak secara otomatis oleh server backend dengan notifikasi yang presisi.
* **Terisolasi Per-Tab (`sessionStorage`)**: Menyimpan token autentikasi secara terisolasi per-tab browser. Menyalin URL ke tab baru atau window lain tidak akan merusak atau mengeluarkan sesi tab utama.
* **Pengalihan Sesi Instan**: Apabila sesi tidak valid atau kedaluwarsa, aplikasi seketika mengarahkan pengguna secara bersih ke layar Login UI tanpa menampilkan teks mentah JSON.

### 2. 📱 Pemindai & Pembuat QR Code Real-Time
* **Kamera Live & Upload Gambar QR**: Pemindaian otomatis via kamera WebCam live maupun **Upload Foto Stiker QR** dengan dekode piksel JavaScript `jsQR` (0% CPU lag).
* **Cetak & Download Stiker QR (PNG)**: Menggenerasi stiker QR Code aset secara murni di peramban (*client-side canvas*) yang dapat langsung diunduh sebagai gambar PNG atau dicetak.

### 3. 📦 Manajemen Inventaris Aset & Mutasi Lokasi
* **Pencatatan Aset Terpadu**: Mengelola data aset lengkap dengan Kode Aset, Nama, Kategori, Lokasi Registrasi, PIC, dan Status (*Active, Maintenance, Damaged, Disposed*).
* **Riwayat Mutasi Lokasi**: Pencatatan resmi perpindahan fisik aset dari lokasi awal ke lokasi baru beserta penanggung jawab (PIC).

### 4. 🛠️ Pelaporan Kerusakan (Work Order) & Timeline Progress
* **Pelaporan Mudah & Fleksibel**: Pengajuan tiket kerusakan dengan tingkat prioritas (*Low, Normal, High, Emergency*) dan pencatatan lokasi kamar/ruangan.
* **Timeline Progres Perbaikan**: Pencatatan kronologi tindakan teknisi (*Action Taken*), biaya perbaikan (*Cost*), dan update status tiket hingga tahap **Finish / Closed**.

### 5. 📅 Kalender Preventive Maintenance (PM)
* **Jadwal Rutin Otomatis**: Manajemen jadwal perawatan berkala (*Daily, Weekly, Monthly, Quarterly, Yearly*) dengan tampilan kalender interaktif dan penyelesaian *Checklist*.

### 6. 📄 Ekspor Laporan & Cetak PDF Terpadu (`exportUtils.js`)
* **Ekspor Excel (.xls)** & **Cetak Laporan PDF**: Modul utilitas independen `exportUtils.js` yang memungkinkan pengunduhan data aset, work order, dan maintenance dalam format Excel serta pembuatan dokumen siap cetak.

### 7. 📱 Desain UI Responsive (Android & iOS)
* **Layout Adaptif**: Antarmuka disesuaikan khusus untuk layar smartphone dengan *touch-friendly target*, tombol elastis bebas terpotong, dan *smooth horizontal swipe*.

---

## 👥 Matriks Hak Akses (Role-Based Access Control)

| Role / Jabatan | Hak Akses Utama |
| :--- | :--- |
| **👑 Administrator (`admin`)** | Akses penuh atas seluruh aset, Work Order, penugasan, jadwal PM, user management, dan audit log. |
| **⭐ HOD Engineer (`hod`)** | Registrasi & Edit Aset, Penugasan Teknisi, Penutupan Tiket (*Finish*), Hapus WO/Aset/PM, Cetak Laporan & Export Excel. |
| **👔 Supervisor (`management`)** | Monitoring Aset, Penugasan Teknisi, Update Status, Mutasi Lokasi, Hapus Data, Cetak Laporan. |
| **🛠️ Staff Engineer (`engineer`)** | Update progres pengerjaan tiket, isi tindakan perbaikan, biaya, dan penyelesaian *Checklist PM*. |
| **🛎️ Staff Hotel (`external` / Dept)** | Pengajuan Tiket Kerusakan Aset berdasarkan Lokasi/Kamar dan pemantauan status perbaikan real-time. |

---

## 🛠️ Arsitektur & Tech Stack

* **Frontend**: Vue 3 (Composition API) + Vite + Vanilla CSS iOS Design System + `jsQR`
* **Backend**: Go (Golang) + `net/http` + GORM ORM + JWT Authentication + Bcrypt
* **Database**: MySQL / MariaDB (`db_sistemasetku`) & SQLite Support
* **Keamanan**: Single Active Session Token Validation, Role-Based Access Control (RBAC), Interceptor 401/403.

---

## 🚀 Panduan Instalasi & Jalankan Sistem

### 1. Prasyarat Sistem
* **Go (Golang)**: Versi 1.20 atau lebih baru
* **Node.js**: Versi 18 atau lebih baru
* **MySQL / MariaDB**: Server database aktif (Port 3306)

### 2. Jalankan Backend (Go API Server)
```bash
cd backend

# Salin file environment jika belum ada
copy .env.example .env

# Jalankan server Go (Port 8080)
go run main.go
```
*Server Backend akan berjalan otomatis di `http://localhost:8080` dan melakukan auto-migration basis data `db_sistemasetku`.*

### 3. Jalankan Frontend (Vue 3 Dev Server)
```bash
cd frontend

# Install dependensi
npm install

# Jalankan dev server lokal & LAN
npm run dev -- --host
```
*Frontend akan terbuka di `http://localhost:3000` (atau IP Lokal LAN seperti `http://192.168.x.x:3000` untuk akses dari HP).*

---

## 📂 Struktur Repositori

```text
sistem-assetku/
├── backend/                  # RESTful API Go (Golang)
│   ├── config/               # Database Connection & Auto-Migration
│   ├── controllers/          # Business Logic (Auth, Asset, WO, Mutation, PM)
│   ├── db/                   # Updated Schema SQL (schema.sql)
│   ├── middlewares/          # JWT & Single Active Session Middleware
│   ├── models/               # Struct GORM (User, Asset, WorkOrder, PM, Log)
│   ├── routes/               # HTTP Route Handler Definitions
│   └── main.go               # Entrypoint Backend Go
│
├── frontend/                 # Single Page Application (Vue 3 + Vite)
│   ├── src/
│   │   ├── components/       # HeaderNavbar, QrScannerModal, StatusBadge, ModalDialog
│   │   ├── utils/            # exportUtils.js, auth.js
│   │   ├── views/            # Dashboard, AssetManagement, WorkOrder, Maintenance, ActivityLog, UserManagement
│   │   ├── router/           # Vue Router Navigation Guards
│   │   └── api.js            # Axios Interceptors & Dynamic BaseURL
│   └── vite.config.js        # Vite Config with HTML Refresh Bypass Proxy
│
└── README.md                 # Dokumentasi Resmi Proyek
```

---

## 📄 Hak Cipta & Lisensi

© 2026 **Sistem AsetKu Hotel**. Hak Cipta Dilindungi Undang-Undang.  
Sistem Informasi Manajemen Aset & Operasional Perhotelan.

