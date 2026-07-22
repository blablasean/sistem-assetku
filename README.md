# Sistem AsetKu 🏨🏢🛠️

**Sistem AsetKu** adalah aplikasi manajemen aset hotel, perawatan rutin (*Preventive Maintenance*), pelaporan kerusakan (*Work Order*), dan pemantauan utilitas internal (listrik & air PDAM) berbasis *on-premise local server*.

Aplikasi ini dirancang khusus untuk efisiensi operasional hotel dan tim *Engineering*, mempermudah staf melapor kerusakan aset di lokasi/kamar, serta mencatat riwayat mutasi lokasi barang (seperti TV, AC, dan kursi yang kerap berpindah antar ruangan).

---

## 👥 4 Kategori User & Matriks Hak Akses (RBAC)

1. **👔 Management Engineer (Supervisor)**
   * Login, **Activity Log** (Audit Trail), **Dashboard Asset Monitoring**, **Asset Management** (Sort, Filter, Search, Detail, Mutasi, Reserve, Document, Location Tracking), dan **Work Order Management** (Priority, History, Cancel, Close, Assign Worker).

2. **🛠️ Staff Engineer (Teknisi Lapangan)**
   * Login, **Activity Log**, **Dashboard Asset Monitoring**, **Asset Management**, dan **Maintenance Management** (Reminder, Report, History, Preventive Maintenance Schedule).

3. **🛎️ External User (Staff Hotel / Departemen Lain)**
   * Login, **Activity Log**, **Work Order / Repair Request** (Pelaporan kerusakan berdasarkan Lokasi/Kamar & Priority), dan **Work Order Status Tracking** (Pemantauan status perbaikan real-time).

4. **🗂️ Head of Department (HOD) Engineer**
   * Login, **Profile Asset** (Registrasi & Edit Asset), **Generate QR Code**, **Asset Management**, **Maintenance Management**, dan **Preventive Maintenance Checklist**.

---

## 🔄 Alur Operasional Perbaikan Kerusakan & Mutasi

1. **Pelaporan Kerusakan**: Staff hotel melapor kerusakan aset di ruangan/kamar dengan menentukan **Lokasi** dan **Priority** (*Low, Medium, High, Emergency*). Bisa diajukan langsung atau melalui hasil pemindaian **QR Code** aset.
2. **Penugasan**: HOD / Management Engineer meninjau prioritas lalu menugaskan teknisi (*Assign Worker*).
3. **Pengerjaan**: Staff Engineer menerima notifikasi, mengerjakan perbaikan (*In Progress*), dan melapor setelah selesai.
4. **Review & Closed**: HOD/Management mereview pengerjaan lalu menutup tiket (*Close Work Order*).
5. **Pelacakan Mutasi Lokasi**: Pencatatan resmi perpindahan fisik aset dari ruang awal ke lokasi baru beserta penanggung jawab (PIC) sehingga riwayat lokasi tidak hilang.

---

## 🛠️ Tech Stack & Deployment Target

* **Frontend**: Vue 3 (Vite) + Component Architecture (Modern & Lightweight)
* **Backend**: Go (Golang) + `net/http` + GORM + JWT Auth + Bcrypt
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
│   │   ├── components/       # HeaderNavbar (Account Dropdown), QrScannerModal, StatusBadge, ModalDialog
│   │   ├── views/            # Dashboard, AssetManagement, WorkOrder, PreventiveMaintenance, ActivityLog, Utility
│   │   ├── router/           # Vue Router config
│   │   └── api.js            # Axios client with JWT interceptors
│   └── vite.config.js
│
├── docker-compose.yml        # Docker Multi-container orchestration
└── *.md                      # Dokumentasi Proyek (QUICK_START, SETUP, DEPLOYMENT, SECURITY, CHANGELOG)
```
