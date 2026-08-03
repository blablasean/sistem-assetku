-- ===================================================
-- SQL SETUP INITIAL DATABASE SISTEM ASETKU HOTEL
-- Database: db_sistemasetku (Bersih / Siap Operasional)
-- ===================================================

CREATE DATABASE IF NOT EXISTS `db_sistemasetku`;
USE `db_sistemasetku`;

-- 1. AKUN SUPERADMIN UTAMA (Default Login)
-- Username: admin | Password: admin123 | Role: admin
INSERT INTO `users` (`username`, `password`, `name`, `role`, `created_at`) VALUES
('admin', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Administrator Utama Hotel', 'admin', NOW())
ON DUPLICATE KEY UPDATE `password` = VALUES(`password`), `name` = VALUES(`name`), `role` = VALUES(`role`);

