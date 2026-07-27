-- ========================================================
-- Database schema for db_sistemasetku (Updated Latest Version)
-- ========================================================
-- Safe recreate script for db_sistemasetku
-- Drops tables in dependency-safe order, then recreates them with
-- explicit column types and named foreign-key constraints matching GORM models.

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS work_order_logs;
DROP TABLE IF EXISTS work_orders;
DROP TABLE IF EXISTS asset_mutation_timelines;
DROP TABLE IF EXISTS mutations;
DROP TABLE IF EXISTS preventive_maintenances;
DROP TABLE IF EXISTS maintenance_histories;
DROP TABLE IF EXISTS activity_logs;
DROP TABLE IF EXISTS spare_parts;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS users;

SET FOREIGN_KEY_CHECKS = 1;

-- 1. MASTER DATA TABLES

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(30) NOT NULL,
    avatar LONGTEXT,
    active_token TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE assets (
    id INT NOT NULL AUTO_INCREMENT,
    asset_code VARCHAR(50) NOT NULL UNIQUE,
    asset_name VARCHAR(100) NOT NULL,
    category VARCHAR(100) DEFAULT 'General',
    location VARCHAR(100) DEFAULT 'Main Area',
    registration_location VARCHAR(100) DEFAULT 'Main Area',
    pic VARCHAR(100) DEFAULT 'Engineering',
    status VARCHAR(30) NOT NULL DEFAULT 'Active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE spare_parts (
    id INT NOT NULL AUTO_INCREMENT,
    part_name VARCHAR(100) NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    min_stock INT NOT NULL DEFAULT 5,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. TRANSACTIONS & BUSINESS TABLES

CREATE TABLE mutations (
    id INT NOT NULL AUTO_INCREMENT,
    asset_id INT NOT NULL,
    previous_location VARCHAR(100) NOT NULL,
    new_location VARCHAR(100) NOT NULL,
    pic_id INT NOT NULL,
    mutation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_mutations_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    CONSTRAINT fk_mutations_pic FOREIGN KEY (pic_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE asset_mutation_timelines (
    id INT NOT NULL AUTO_INCREMENT,
    asset_id INT NOT NULL,
    previous_location VARCHAR(100) NOT NULL,
    new_location VARCHAR(100) NOT NULL,
    pic VARCHAR(100) NOT NULL,
    mutation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE work_orders (
    id INT NOT NULL AUTO_INCREMENT,
    wo_code VARCHAR(50),
    asset_id INT NOT NULL,
    location VARCHAR(100),
    description TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'Pending',
    priority VARCHAR(30) DEFAULT 'Normal',
    reporter VARCHAR(100) DEFAULT 'Staff',
    department VARCHAR(100) DEFAULT 'General',
    requester_id INT NOT NULL,
    engineer_id INT DEFAULT NULL,
    cost INT DEFAULT 0,
    action_taken TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_work_orders_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    CONSTRAINT fk_work_orders_requester FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE work_order_logs (
    id INT NOT NULL AUTO_INCREMENT,
    work_order_id INT NOT NULL,
    status VARCHAR(30) NOT NULL,
    action_taken TEXT,
    cost INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE preventive_maintenances (
    id INT NOT NULL AUTO_INCREMENT,
    asset_id INT NOT NULL,
    schedule_type VARCHAR(30) NOT NULL,
    next_run DATETIME NOT NULL,
    checklist_data TEXT,
    completed_dates TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'Active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_preventive_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. HISTORY & AUDIT TABLES

CREATE TABLE maintenance_histories (
    id INT NOT NULL AUTO_INCREMENT,
    asset_id INT NOT NULL,
    action_taken TEXT NOT NULL,
    cost INT NOT NULL DEFAULT 0,
    completion_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_maint_hist_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE activity_logs (
    id INT NOT NULL AUTO_INCREMENT,
    user_id INT NOT NULL,
    action VARCHAR(255) NOT NULL,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
