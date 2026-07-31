# 🔒🚀 Panduan Deploy HTTPS Sistem AssetKu (IP 172.17.10.109 Tanpa Port)

Dokumen ini berisi panduan resmi untuk menjalankan **Sistem AssetKu Hotel** di protokol **HTTPS** secara langsung menggunakan alamat IP (`https://172.17.10.109`) tanpa menggunakan nomor port di URL.

---

## 🌟 Mengapa Menggunakan HTTPS pada IP?

1. **Akses Kamera Live 100% Bebas Izin**: Browser modern (Chrome/Safari) mewajibkan koneksi HTTPS untuk izin akses streaming kamera live (`getUserMedia`). Dengan HTTPS, pemindaian QR code via kamera HP/laptop berjalan 100% mulus.
2. **URL Bersih & Profesional**: Pengguna cukup mengetik `https://172.17.10.109` tanpa perlu mengingat port `:3000` atau `:8080`.
3. **Enkripsi Keamanan SSL/TLS**: Seluruh data kredensial login, aset, dan tiket Work Order terenkripsi aman saat melalui jaringan Wi-Fi LAN/Server.

---

## 📁 Struktur Paket Deploy (`/deploy`)

```text
sistem-assetku/
├── deploy/
│   ├── ssl/
│   │   └── generate_ssl.ps1     # Script pembuat sertifikat SSL Self-Signed IP
│   ├── nginx/
│   │   └── nginx.conf           # Konfigurasi Nginx Production HTTPS (Port 443)
│   ├── caddy/
│   │   └── Caddyfile            # Alternative 1-File Caddy HTTPS Web Server
│   └── DEPLOYMENT_GUIDE.md      # Panduan Deployment Resmi
```

---

## 🛠️ Langkah-Langkah Deployment (Windows Server / Local Server)

### 1. Buat Sertifikat SSL untuk IP `172.17.10.109`
1. Buka PowerShell sebagai **Administrator**.
2. Jalankan script pembuat sertifikat SSL:
   ```powershell
   cd c:\Users\Sean\Documents\KERJA\PROJECT\sistem-assetku\deploy\ssl
   .\generate_ssl.ps1
   ```
3. Script akan menghasilkan file `server.crt` dan `server.key` di folder `deploy/ssl/`.

---

### 2. Jalankan Backend Go Executable
1. Pastikan backend telah di-build menjadi file `server.exe`:
   ```powershell
   cd backend
   go build -o server.exe main.go
   ```
2. Jalankan `server.exe`:
   ```powershell
   .\server.exe
   ```
   *Backend akan aktif di `http://127.0.0.1:8080` di latar belakang.*

---

### 3. Jalankan Frontend Build
1. Pastikan frontend telah di-compile ke folder `dist/`:
   ```powershell
   cd frontend
   npm run build
   ```

---

### 4. Menjalankan Web Server HTTPS (Pilih Salah Satu Opsi)

#### 🔹 OPSI A: Menggunakan Nginx (Sangat Direkomendasikan)
1. Unduh **Nginx for Windows** dari [nginx.org](https://nginx.org/en/download.html) dan ekstrak ke `C:\nginx`.
2. Salin isi file `deploy/nginx/nginx.conf` ke `C:\nginx\conf\nginx.conf`.
3. Salin file `server.crt` & `server.key` dari `deploy/ssl/` ke `C:\nginx\ssl\`.
4. Jalankan Nginx dari Command Prompt / PowerShell Administrator:
   ```cmd
   cd C:\nginx
   start nginx
   ```
5. Akses aplikasi di browser via **`https://172.17.10.109`**!

#### 🔹 OPSI B: Menggunakan Caddy (Zero-Config 1-File Server)
1. Unduh **Caddy** dari [caddyserver.com](https://caddyserver.com/download) (1 file `caddy.exe`).
2. Jalankan Caddy dari folder `deploy/caddy`:
   ```cmd
   caddy run --config Caddyfile
   ```
3. Caddy akan otomatis memproses SSL dan reverse proxy ke port 8080. Akses via **`https://172.17.10.109`**!

---

## 🎯 Verifikasi Deployment Selesai

- Buka browser di HP/Laptop: `https://172.17.10.109`
- Cek URL bar: Protokol **HTTPS** aktif tanpa port tambahan.
- Klik tombol **"Buka Kamera"** pada modal QR Scanner: Kamera live HP/Laptop langsung aktif secara mulus!
