# sistem-assetku
Sistem manajemen aset internal, pelacakan tiket perbaikan (Work Order), dan pemantauan utilitas berbasis local server (on-premise).

# Sistem AsetKu 🏢🛠️

**Sistem AsetKu** adalah aplikasi manajemen aset, perawatan (*maintenance*), dan pemantauan utilitas internal perusahaan[cite: 2]. Aplikasi ini dirancang untuk dijalankan pada *local server* (on-premise) guna mempermudah tim *Engineering* dan manajemen dalam melacak siklus hidup aset, mengelola tiket perbaikan, serta memantau konsumsi energi.

---

## 🚀 Fitur Utama

Sistem ini dibagi menjadi beberapa modul utama berdasarkan kebutuhan operasional[cite: 2]:

### 📦 1. Asset Management
* **Registrasi Aset:** Pencatatan aset baru ke dalam sistem secara terstruktur[cite: 2].
* **Generate & Scan QR Code:** Sistem otomatis membuat QR Code unik untuk setiap aset guna mempermudah akses informasi dan pelaporan kerusakan[cite: 2].
* **Mutasi & Tracking Lokasi:** Mencatat perpindahan aset beserta PIC yang bertanggung jawab[cite: 2].
* **Asset Status Management:** Pemantauan kondisi aset secara real-time (*Active, Maintenance, Damaged, Retired*)[cite: 2].

### 🔧 2. Maintenance & Work Order Management
* **Work Order / Repair Request:** Pembuatan tiket perbaikan oleh departemen terkait atau setelah melakukan scan QR Code aset[cite: 2].
* **Work Order Status Tracking:** Pemantauan progres perbaikan teknisi (*Open, In Progress, Completed*)[cite: 2].
* **Preventive Maintenance Schedule & Checklist:** Penjadwalan perawatan rutin otomatis untuk mencegah kerusakan fatal pada aset[cite: 2].
* **Maintenance History:** Penyimpanan riwayat perbaikan untuk analisis performa aset[cite: 2].

### 📊 3. Utility Monitoring & Reporting
* **Monitoring Air PDAM & Listrik:** Pencatatan dan pemantauan konsumsi utilitas harian/bulanan.
* **Dashboard Asset Monitoring:** Visualisasi data statistik jumlah aset, kondisi terkini, dan laporan performa tim engineering[cite: 2].

### 🔐 4. Sistem Pendukung
* **User Management (Role & Permission):** Pengaturan hak akses multi-user (Management Engineering, Staff Engineer, Head of Department, & External User)[cite: 1, 2].
* **Activity Log:** Pencatatan riwayat aktivitas pengguna di dalam sistem untuk kebutuhan audit[cite: 1].

---

## 🛠️ Tech Stack (Teknologi yang Digunakan)

Aplikasi ini menggunakan pendekatan *Separation of Concerns* dengan memisahkan Frontend dan Backend (Monorepo):

* **Frontend:** Vue.js 3 (Vite) + Tailwind CSS / UI Framework
* **Backend:** Go (Golang) + Gin/Fiber Framework
* **Database:** PostgreSQL / MySQL
* **Deployment Target:** Windows Server (Local Network/On-Premise)

---

## 📂 Struktur Folder Proyek

```text
sistem-asetku/
├── backend/             # RESTful API menggunakan Go (Golang)
│   ├── config/          # Koneksi database & environment
│   ├── controllers/     # Logika bisnis & handling request
│   ├── models/          # Skema tabel database (GORM Structs)
│   ├── routes/          # Endpoint API URL
│   └── main.go          # Main entrypoint aplikasi backend
│
└── frontend/            # Single Page Application (SPA) menggunakan Vue.js
    ├── src/
    │   ├── components/  # Komponen reusable (Scanner, Buttons, Forms)
    │   ├── views/       # Halaman utama (Dashboard, Asset, Work Order)
    │   └── main.js      # Main entrypoint aplikasi frontend
    └── vite.config.js
