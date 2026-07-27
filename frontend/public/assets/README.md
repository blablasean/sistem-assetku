# 🖼️ Panduan Manajemen Asset Gambar & Media - AsetKu Hotel

Folder ini (`frontend/public/assets/`) disiapkan khusus untuk menyimpan seluruh file gambar fisik, logo resmi hotel, ikon, serta media aset terdaftar agar tampilan aplikasi **AsetKu Hotel** tampak 100% profesional, bersih, dan konsisten (bebas dari penggunaan emoji/stiker sementara).

---

## 📁 Daftar File Gambar Yang Diperlukan & Penamaannya

Silakan masukkan file gambar Anda ke dalam folder ini dengan **nama file persis** seperti yang tercantum pada tabel di bawah ini:

| Nama File | Format | Ukuran Disarankan | Deskripsi & Kegunaan |
| :--- | :--- | :--- | :--- |
| **`logo.png`** | PNG (Transparan) | 512 × 512 px | **Logo Resmi Hotel / Perusahaan**. Digunakan pada Header Navbar atas, Halaman Login, dan Cetak Sticker QR Code Aset. |
| **`logo_icon.png`** | PNG (Transparan) | 128 × 128 px | **Logo Ikon Ringkas**. Digunakan pada sudut tombol kecil & favicon tab browser. |
| **`asset_placeholder.png`** | PNG / JPG | 800 × 600 px | **Foto Fallback Aset**. Digunakan secara otomatis jika suatu aset terdaftar belum diunggah foto fisiknya. |
| **`hotel_banner.jpg`** | JPG / WEBP | 1920 × 1080 px | **Banner Foto Fisik Hotel**. Digunakan sebagai gambar latar belakang opsional pada halaman Login & Dashboard. |
| **`qr_scanner_frame.png`** | PNG (Transparan) | 600 × 600 px | **Bingkai Pemindai QR Code**. Digunakan sebagai overlay panduan kamera pemindai QR Code di HP. |

---

## 💡 Petunjuk Penggunaan Di Aplikasi

* Setelah Anda menambahkan file gambar (misalnya `logo.png`), file tersebut secara otomatis dapat diakses oleh sistem di URL `/assets/logo.png`.
* Sistem frontend secara otomatis akan mendeteksi dan menampilkan `logo.png` resmi Anda pada Header Navbar, Modal Cetak QR Code, dan Halaman Login menggantikan stiker emoji sebelumnya.

---

## 📐 Rekomendasi Format & Kualitas File

1. **Format PNG Transparan**: Untuk `logo.png` dan `logo_icon.png`, sangat disarankan menggunakan background transparan (PNG-24) agar menyatu sempurna dengan latar belakang putih maupun gelap.
2. **Kompresi Gambar**: Gunakan alat kompresi gambar gratis seperti [TinyPNG](https://tinypng.com/) agar ukuran file tetap ringan (< 300 KB) sehingga loading aplikasi super cepat.
