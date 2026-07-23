-- ===================================================
-- SQL SEED DATABASE UNTUK SISTEM ASETKU HOTEL
-- Impor SQL ini di phpMyAdmin (Database: db_sistemasetku)
-- ===================================================

CREATE DATABASE IF NOT EXISTS `db_sistemasetku`;
USE `db_sistemasetku`;

-- 1. TABEL USERS (4 Kategori Role User)
-- Password default untuk semua user di bawah adalah: password
INSERT INTO `users` (`id`, `username`, `password`, `name`, `role`, `created_at`) VALUES
(1, 'mgr_eng', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Budi (Supervisor)', 'management', NOW()),
(2, 'hod_eng', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Pak Alex (HOD Engineer)', 'hod', NOW()),
(3, 'tech_deni', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Deni (Staff Engineer)', 'engineer', NOW()),
(4, 'staff_frontdesk', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Rina (Staff Front Desk)', 'external', NOW()),
(5, 'admin', '$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu', 'Administrator Hotel', 'management', NOW())
ON DUPLICATE KEY UPDATE `password` = VALUES(`password`), `name` = VALUES(`name`), `role` = VALUES(`role`);

-- 2. TABEL ASSETS (Aset Hotel Sampel)
INSERT INTO `assets` (`id`, `asset_code`, `asset_name`, `category`, `location`, `pic`, `status`, `document_url`, `is_reserved`, `created_at`) VALUES
(1, 'AST-RM301-AC', 'AC Split Daikin 1.5 PK', 'HVAC / AC', 'Kamar 301', 'Deni (Tech)', 'Active', '', 0, NOW()),
(2, 'AST-RM102-TV', 'Smart TV LG 43 Inch', 'Elektronik & TV', 'Kamar 102', 'Front Desk Team', 'Active', '', 0, NOW()),
(3, 'AST-KCH-CHILLER', 'Chiller Dapur Utama', 'Kitchen Equipment', 'Kitchen Dapur', 'Kitchen Chef', 'Maintenance', '', 0, NOW()),
(4, 'AST-GEN-01', 'Generator Unit Cummins 500kVA', 'Mesin & Generator', 'Power House', 'Engineering Supervisor', 'Active', '', 1, NOW()),
(5, 'AST-LBY-SOFA', 'Set Sofa Premium Leather', 'Mebel & Furniture', 'Lobby Lounge', 'Housekeeping', 'Active', '', 0, NOW())
ON DUPLICATE KEY UPDATE `asset_name`=`asset_name`;

-- 3. TABEL WORK_ORDERS (Tiket Laporan Kerusakan & Perbaikan)
INSERT INTO `work_orders` (`id`, `asset_id`, `location`, `priority`, `description`, `status`, `requester_id`, `engineer_id`, `action_taken`, `cost`, `created_at`) VALUES
(1, 1, 'Kamar 301', 'Emergency', 'AC Kamar 301 bocor air dan tidak dingin', 'In Progress', 4, 3, 'Pengecekan pipa freon dan pembersihan drainase', 150000, NOW()),
(2, 3, 'Kitchen Dapur', 'High', 'Chiller Dapur Utama suhu naik ke -5°C', 'Open', 4, 0, '', 0, NOW()),
(3, 2, 'Kamar 102', 'Medium', 'Smart TV HDMI port tidak terdeteksi', 'Closed', 4, 3, 'Ganti modul port HDMI TV', 250000, NOW())
ON DUPLICATE KEY UPDATE `description`=`description`;
