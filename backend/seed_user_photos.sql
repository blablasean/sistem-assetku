-- ============================================================
-- SQL Script: Tabel & Fitur Foto Profil Pengguna (User Photos)
-- Database: MySQL / MariaDB (Database: db_sistemasetku atau sistem_asetku)
-- ============================================================

USE `sistem_asetku`;

-- 1. Tabel Khusus Penyimpanan Foto Profil (`user_photos`)
CREATE TABLE IF NOT EXISTS `user_photos` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `user_id` INT NOT NULL UNIQUE,
  `avatar_data` LONGTEXT NOT NULL,
  `file_name` VARCHAR(255) DEFAULT 'profile_avatar.jpg',
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_user_photos_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Tambah Kolom 'avatar' pada Tabel 'users' Utama (GORM Sync)
ALTER TABLE `users` ADD COLUMN `avatar` LONGTEXT NULL AFTER `role`;

-- 3. Tampilkan Verifikasi Struktur Tabel
DESCRIBE `users`;
DESCRIBE `user_photos`;
