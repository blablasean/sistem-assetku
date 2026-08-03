# 🚀 Panduan Deployment Sistem AssetKu ke Server (Production Guide)

Dokumen ini menjelaskan langkah demi langkah untuk melakukan **Deployment Production** pada aplikasi **Sistem AssetKu** agar dapat diakses oleh seluruh pengguna di jaringan lokal (LAN Kantor) maupun secara Online via Internet (Cloud VPS).

---

## 📌 Pilihan 1: Deploy di Server Lokal / PC Kantor (LAN Windows Server)
*Opsi paling praktis untuk penggunaan internal kantor tanpa biaya sewa cloud.*

### Langkah 1: Build Frontend Vue
1. Buka terminal di folder `frontend`:
   ```bash
   cd frontend
   npm run build
   ```
2. Hasil build berupa folder `dist` siap digunakan.

### Langkah 2: Build Backend Go
1. Buka terminal di folder `backend`:
   ```bash
   cd backend
   go build -o server.exe main.go
   ```
2. Pindahkan atau pastikan folder `dist` dari frontend diletakkan berdampingan dengan `server.exe` atau berada di folder induk `../frontend/dist`.

### Langkah 3: Jalankan Server
1. Jalankan `server.exe`.
2. Backend secara otomatis akan melayani **API Backend + Tampilan Web Frontend** sekaligus pada port `8080`.
3. Akses aplikasi dari komputer manapun di jaringan LAN kantor melalui browser:
   ```text
   http://IP_SERVER_KANTOR:8080
   ```
   *(Contoh: http://192.168.1.100:8080)*

---

## ☁️ Pilihan 2: Deploy di Cloud Server (Docker / Linux VPS)
*Opsi untuk deploy ke VPS seperti Biznet, CloudKilat, Niagahoster, AWS, DigitalOcean, atau Railway.*

Proyek ini telah dilengkapi dengan **Dockerfile Multi-Stage** & **docker-compose.yml**.

### Langkah Deploy dengan Docker (1-Click):
1. Clone repositori ke server Linux VPS Anda:
   ```bash
   git clone https://github.com/blablasean/sistem-assetku.git
   cd sistem-assetku
   ```
2. Jalankan Docker Compose:
   ```bash
   docker-compose up -d --build
   ```
3. Docker akan secara otomatis:
   - Menjalankan container database MySQL 8.0 (`assetku_db`)
   - Membangun frontend & backend ke dalam 1 container cepat (`assetku_app`)
   - Menjalankan aplikasi secara otomatis 24/7 di port `8080`.

---

## 🛠️ Konfigurasi Environment Variable (`.env`)

Untuk produksi, pastikan file `.env` di folder `backend` dikonfigurasi sebagai berikut:

```env
SERVER_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASS=password_mysql_server
DB_NAME=db_sistemasetku
JWT_SECRET=sistemasetku_secret_production_key_2026
```

---

## 🔒 Akun Superadmin Default (Produksi)
Setelah dipasang di server baru, gunakan kredensial bawaan berikut untuk login pertama kali:
- **URL**: `http://IP_SERVER:8080`
- **Username**: `admin`
- **Password**: `admin123`
