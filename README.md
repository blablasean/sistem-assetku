# Sistem AsetKu 🏨🏢🛠️

**Sistem AsetKu** adalah aplikasi manajemen aset hotel, perawatan rutin (*Preventive Maintenance*), pelaporan kerusakan (*Work Order*), dan pencatatan mutasi lokasi barang berbasis *on-premise local server*.

Aplikasi ini dirancang khusus untuk efisiensi operasional hotel dan tim *Engineering*, mempermudah staf melapor kerusakan aset di lokasi/kamar, serta mencatat riwayat mutasi lokasi barang (seperti TV, AC, dan kursi yang kerap berpindah antar ruangan).

---

## 👥 5 Kategori User & Matriks Hak Akses (RBAC)

1. **👑 Administrator Hotel (`admin`)**
   * **Kontrol Penuh (Full System Control)** atas seluruh aset, Work Order, penugasan, jadwal PM, user management, dan audit log.

2. **⭐ Head of Department (`hod`) Engineer**
   * Pimpinan operasional engineering: Registrasi & Edit Aset, Penugasan Teknisi (*Assign Worker*), Penutupan Tiket (*Finish Work Order*), **Hapus Work Order**, **Hapus Aset**, **Hapus & Edit PM Schedule**, Edit/Hapus Activity Log, Cetak Laporan & Export Excel.

3. **👔 Management Engineer (`management` / Supervisor)**
   * Supervisor operasional: Monitoring Aset, Penugasan Teknisi, Update Status, Mutasi Lokasi, Hapus WO/Aset/PM/Log, Cetak Laporan & Export Excel.

4. **🛠️ Staff Engineer (`engineer` / Teknisi Lapangan)**
   * Pemeliharaan & Perbaikan: Update progres pengerjaan tiket, isi rincian tindakan perbaikan, biaya, dan penyelesaian *Checklist Preventive Maintenance*.

5. **🛎️ External User (`external` / Staff Hotel & Front Desk)**
   * Pelaporan & Pantau: Pembuatan Tiket Kerusakan Aset berdasarkan Lokasi/Kamar, Pembatalan Tiket status Open, dan pemantauan status perbaikan real-time.

---

## 🚀 Fitur Utama & Keunggulan Sistem

* **📱 Pemindai & Dekoder QR Code Real-Time (`jsQR`)**: Memindai stiker QR Code Aset via kamera WebCam live maupun **Upload Foto QR** dengan dekode piksel JavaScript `jsQR`, audio *beep* respon, dan optimasi delay (0% CPU lag).
* **🖨️ Generator QR Code Client-Side**: Menggenerasi gambar 2D QR Code aset secara murni di peramban (*client-side canvas*) bebas error dan siap cetak stiker.
* **📄 Integrated Single Button "Laporan & Export"**: Menyatukan tombol laporan di header halaman Work Order & Activity Log menjadi **`📄 Laporan & Export`** yang membuka pratinjau dokumen lengkap dengan opsi cetak PDF & export Excel (.xlsx).
* **📋 Activity Log Terintegrasi**: Menyajikan riwayat **Work Order Selesai (Finish)** dan **Maintenance Selesai** lengkap dengan fitur Edit, Hapus, pencarian serbaguna, dan ringkasan total biaya.
* **✅ Status Tiket "Finish"**: Tahap akhir penutupan tiket perbaikan secara resmi dilabeli sebagai **Finish**.
* **🗑️ Manajemen Penghapusan Data (WO, Aset, PM, Log)**: Penghapusan data permanen dari MySQL dengan *Optimistic UI Update* yang responsif.
* **👤 Custom UI Profile & Floating Toast**: Notifikasi respon aksi dan tampilan profil pengguna berbentuk UI Modal & Floating Toast modern tanpa `alert()` bawaan browser.
* **🔒 Keamanan Sesi Browser (`sessionStorage`)**: Sesi login otomatis terhapus saat tab atau browser ditutup untuk keamanan data hotel.

---

## 🔄 Alur Operasional Perbaikan Kerusakan & Mutasi

1. **Pelaporan Kerusakan**: Staff hotel melapor kerusakan aset di ruangan/kamar dengan menentukan **Lokasi** dan **Priority** (*Low, Medium, High, Emergency*). Bisa diajukan langsung atau melalui hasil pemindaian **QR Code** aset.
2. **Penugasan**: HOD / Management Engineer meninjau prioritas lalu menugaskan teknisi (*Assign Worker*).
3. **Pengerjaan**: Staff Engineer menerima penugasan, mengerjakan perbaikan (*In Progress*), dan memasukkan rincian biaya & tindakan perbaikan.
4. **Review & Finish**: HOD/Management/Admin mereview pengerjaan lalu menyelesaikan tiket (*Finish Work Order*).
5. **Pelacakan Mutasi Lokasi**: Pencatatan resmi perpindahan fisik aset dari ruang awal ke lokasi baru beserta penanggung jawab (PIC) sehingga riwayat lokasi tidak hilang.

---

## 🛠️ Tech Stack & Deployment Target

* **Frontend**: Vue 3 (Vite) + Vanilla CSS System + `jsQR` (Modern & Responsive UI)
* **Backend**: Go (Golang) + `net/http` + GORM ORM + JWT Auth + Bcrypt
* **Database**: MySQL / MariaDB (Auto-Migration GORM)
* **Deployment Target**: On-Premise Server (Windows Server / Linux Docker Container)

---

## 📂 Struktur Repositori

```text
sistem-assetku/
├── backend/                  # RESTful API Go (Golang)
│   ├── config/               # Database connection & Auto-migration
│   ├── controllers/          # Business logic (Asset, WO, Mutation, Maintenance, Auth)
│   ├── middlewares/          # Auth JWT & CORS middleware
│   ├── models/               # Struct GORM (User, Asset, WorkOrder, Mutation, Maintenance, Log)
│   ├── routes/               # Endpoint URL routers
│   └── main.go               # Main entrypoint backend
│
├── frontend/                 # Single Page Application (Vue.js 3 + Vite)
│   ├── src/
│   │   ├── components/       # HeaderNavbar, QrScannerModal, StatusBadge, ModalDialog
│   │   ├── views/            # Dashboard, AssetManagement, WorkOrder, PreventiveMaintenance, ActivityLog
│   │   ├── router/           # Vue Router config
│   │   └── api.js            # Axios client with JWT interceptors
│   └── vite.config.js
│
├── docker-compose.yml        # Docker Multi-container orchestration
└── *.md                      # Dokumentasi Proyek (QUICK_START, SETUP, DEPLOYMENT, SECURITY, CHANGELOG)
```
