-- ============================================================
-- SQL Script: Seed Admin & Default Users (Sistem AsetKu)
-- Database: MySQL / MariaDB (Database Name: sistem_asetku atau db_sistemasetku)
-- ============================================================

-- 1. Buat / Pastikan Tabel 'users' Tersedia
CREATE TABLE IF NOT EXISTS `users` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(100) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `name` VARCHAR(150) NOT NULL,
  `role` VARCHAR(50) NOT NULL DEFAULT 'external',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Insert / Update Akun Admin Utama
-- Username: admin | Password: admin123 | Role: admin
INSERT INTO `users` (`username`, `password`, `name`, `role`, `created_at`)
VALUES 
  ('admin', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Administrator Utama Hotel', 'admin', NOW())
ON DUPLICATE KEY UPDATE 
  `password` = '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu',
  `name` = 'Administrator Utama Hotel',
  `role` = 'admin';

-- 3. Insert / Update 4 Akun Default Peran Lainnya (Password untuk semua akun: admin123)
INSERT INTO `users` (`username`, `password`, `name`, `role`, `created_at`)
VALUES
  ('hod_eng', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Pak Alex (HOD Engineer)', 'hod', NOW()),
  ('spv_eng', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Pak Hendra (Supervisor)', 'management', NOW()),
  ('teknisi_budi', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Budi Santoso (Teknisi)', 'engineer', NOW()),
  ('staff_frontdesk', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Siti Rahma (Staff Hotel)', 'external', NOW())
ON DUPLICATE KEY UPDATE 
  `password` = '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu',
  `name` = VALUES(`name`);

-- 4. Verifikasi Data Pengguna Terdaftar
SELECT `id`, `username`, `name`, `role`, `created_at` FROM `users`;
