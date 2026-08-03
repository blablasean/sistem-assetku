# 📖 Panduan Pengguna & Pengoperasian — Sistem AssetKu Hotel

Selamat datang di **Sistem AssetKu Hotel**! Dokumen ini berisi panduan lengkap penggunaan sistem bagi seluruh pengguna (Admin, HOD, Supervisor, Teknisi, dan Staff Departemen Operasional).

---

## 📋 Daftar Isi
1. [Pengenalan Sistem](#1-pengenalan-sistem)
2. [Hak Akses & Akun Bawaan (Default Login)](#2-hak-akses--akun-bawaan-default-login)
3. [Panduan Penggunaan Fitur](#3-panduan-penggunaan-fitur)
   - [A. Login & Keamanan Sesi](#a-login--keamanan-sesi)
   - [B. Dashboard Operasional](#b-dashboard-operasional)
   - [C. Pemindai QR Code Aset (Scan QR)](#c-pemindai-qr-code-aset-scan-qr)
   - [D. Tiket Kerusakan (Work Order)](#d-tiket-kerusakan-work-order)
   - [E. Manajemen Aset & Cetak Stiker QR](#e-manajemen-aset--cetak-stiker-qr)
   - [F. Pemeliharaan Rutin (Preventive Maintenance)](#f-pemeliharaan-rutin-preventive-maintenance)
   - [G. Log Aktivitas Sistem (Activity Log)](#g-log-aktivitas-sistem-activity-log)
   - [H. Manajemen User & Foto Profil](#h-manajemen-user--foto-profil)
4. [Panduan Menjalankan Server & Deployment HTTPS](#4-panduan-menjalankan-server--deployment-https)

---

## 1. Pengenalan Sistem

**Sistem AssetKu Hotel** adalah aplikasi manajemen aset dan tiket kerusakan (*Work Order*) berbasis web yang dirancang khusus untuk memenuhi kebutuhan operasional perhotelan.

### ✨ Fitur Utama:
- **PWA Mobile-First Interface**: Tampilan smartphone yang responsif dengan *iOS Navigation Dock* dan tombol cepat *Scan QR*.
- **Pemindai QR Code Real-Time**: Memindai stiker QR pada fisik aset menggunakan kamera HP tanpa perlu menginstal aplikasi tambahan.
- **Single Active Session**: Mencegah login ganda (*concurrent login*) dari dua perangkat berbeda demi keamanan.
- **Manajemen Work Order**: Pelaporan kerusakan kamar/fasilitas, penugasan teknisi, dan pembaruan status real-time.
- **Preventive Maintenance**: Penjadwalan perawatan harian, mingguan, bulanan, dan tahunan beserta checklist inspeksi.
- **Cetak Stiker QR**: Fitur ekspor dan cetak stiker QR code aset secara langsung dari peramban.

---

## 2. Hak Akses & Akun Bawaan (Default Login)

Sistem membagi pengguna ke dalam beberapa tingkatan peran (*Role*):

| Peran (Role) | Hak Akses Utama |
| :--- | :--- |
| **Administrator** | Akses penuh ke seluruh fitur (User Management, Asset, WO, Maintenance, Log). |
| **HOD Engineer** | Mengelola jadwal maintenance, aset, work order, dan konfirmasi checklist. |
| **Supervisor Engineer** | Mengelola work order, menetapkan teknisi, dan memantau pemeliharaan. |
| **Staff Engineer (Teknisi)** | Menangani tiket perbaikan work order dan melaksanakan inspeksi fisik. |
| **Staff Operasional (Departemen)** | Melaporkan kerusakan fasilitas (Work Order) dan memantau status perbaikan. |

### 🔑 Akun Superadmin Bawaan System:

> [!NOTE]
> Untuk keperluan awal setup produksi, sistem menyediakan 1 akun **Administrator Utama** bawaan:
> - Username: **`admin`**
> - Password: **`admin123`**
> 
> Administrator dapat langsung login dan menambahkan akun pengguna baru (HOD, Supervisor, Teknisi, & Staff) melalui menu **Manajemen Pengguna**.


---

## 3. Panduan Penggunaan Fitur

### A. Login & Keamanan Sesi
1. Buka aplikasi di peramban (misal: `https://172.17.10.109` atau `http://localhost:3000`).
2. Masukkan **Username** dan **Password** Anda, lalu klik **Login**.
3. *Catatan Keamanan*: Jika akun Anda sedang digunakan di perangkat/tab lain, login baru akan membatalkan sesi lama secara otomatis (*Single Active Session*).

---

### B. Dashboard Operasional
- Menampilkan ringkasan jumlah **Total Aset**, **Work Order Aktif**, dan **Jadwal Perawatan Rutin**.
- Di layar HP, tombol **Scan QR** dapat diakses 1-ketukan dari tombol lingkaran biru di tengah dock navigasi bawah.

---

### C. Pemindai QR Code Aset (Scan QR)
1. Klik tombol **Scan QR** (berwarna biru di bagian tengah dock bawah HP atau di navigasi atas).
2. Arahkan kamera HP Anda ke stiker QR Code yang terpasang pada aset hotel (misal: unit AC, Smart TV, Pompa Air).
3. Setelah QR terdeteksi (ditandai dengan nada bising *beep*), informasi detail aset beserta lokasi dan riwayatnya akan muncul otomatis di layar.
4. Klik **Laporkan Kerusakan Work Order** jika aset tersebut membutuhkan perbaikan.

---

### D. Tiket Kerusakan (Work Order)
#### 1. Membuat Tiket Kerusakan Baru:
- Klik tombol **+ Buat Work Order Baru** / **Laporkan Kerusakan**.
- Pilih **Lokasi Ruangan/Kamar**, **Aset yang Rusak**, **Tingkat Prioritas** (*Low*, *Medium*, *High*, *Emergency*), dan sertakan deskripsi detail kerusakan.
- Anda juga dapat mengunggah foto bukti kerusakan.

#### 2. Memproses & Menyelesaikan Tiket (Khusus Tim Engineering):
- Klik tombol **Ambil Pekerjaan** untuk menandai bahwa tiket sedang ditangani.
- Setelah perbaikan fisik selesai, ubah status tiket menjadi **Completed** dan isi rincian penanganan.

---

### E. Manajemen Aset & Cetak Stiker QR
*(Khusus Role Admin, HOD, Supervisor, & Engineer)*
1. Masuk ke halaman **Manajemen Aset**.
2. Klik **+ Tambah Aset Baru** untuk mendaftarkan barang/peralatan hotel baru.
3. Klik ikon **Cetak QR Code** pada baris aset untuk mencetak stiker QR yang dapat ditempelkan langsung pada fisik barang.

---

### F. Pemeliharaan Rutin (Preventive Maintenance)
*(Khusus Role Admin, HOD, Supervisor, & Engineer)*
1. Masuk ke halaman **Maintenance**.
2. Anda dapat melihat kalender jadwal perawatan berkala (**Daily**, **Weekly**, **Monthly**, **Yearly**).
3. Klik pada tanggal jadwal untuk melihat rincian poin inspeksi fisik.
4. Klik **Selesaikan Checklist Perawatan** saat inspeksi rutin telah dilakukan di lapangan.

---

### G. Log Aktivitas Sistem (Activity Log)
*(Khusus Role Admin, HOD, Supervisor)*
- Masuk ke menu **Activity Log** untuk melihat catatan audit (*audit trail*) mencakup waktu login, pembuatan work order, pembaruan aset, dan riwayat aktivitas seluruh pengguna.

---

### H. Manajemen User & Foto Profil
1. **Mengubah Foto Profil Pribadi**:
   - Klik **Foto Profil Avatar** di pojok kanan atas header HP / Laptop.
   - Unggah foto profil baru Anda dan klik **Simpan Foto Profil**.
2. **Pengelolaan Akun Pengguna** *(Khusus Admin)*:
   - Masuk ke menu **User Management**.
   - Admin dapat menambah akun pengguna baru, memperbarui role/jabatan, mereset password, atau menghapus akun.

---

## 4. Panduan Menjalankan Server & Deployment HTTPS

### 💻 Menjalankan di Lingkungan Pengembangan (Development):
1. **Jalankan Backend (Go)**:
   ```bash
   cd backend
   go run main.go
   ```
2. **Jalankan Frontend (Vue 3)**:
   ```bash
   cd frontend
   npm run dev
   ```

---

### 🌐 Deploy Resmi HTTPS di IP Server (`172.17.10.109` Tanpa Port):
1. **Generate Sertifikat SSL**:
   Jalankan script PowerShell di folder project:
   ```powershell
   .\deploy\ssl\generate_ssl.ps1
   ```
2. **Salin Konfigurasi Nginx**:
   Salin `deploy/nginx/nginx.conf` ke folder Nginx server Anda (`C:\nginx\conf\nginx.conf`).
3. **Jalankan Service**:
   - Jalankan Backend Go: `.\backend\server.exe`
   - Jalankan Nginx: `cd C:\nginx && start nginx`
4. Akses melalui smartphone: **`https://172.17.10.109`**
