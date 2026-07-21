-- ========================================================
-- Database schema for db_sistemasetku (provided by user)
-- ========================================================
-- Safe recreate script for db_sistemasetku
-- Drops tables in dependency-safe order, then recreates them with
-- explicit column types and named foreign-key constraints.

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS work_orders;
DROP TABLE IF EXISTS mutations;
DROP TABLE IF EXISTS preventive_maintenances;
DROP TABLE IF EXISTS maintenance_histories;
DROP TABLE IF EXISTS activity_logs;
DROP TABLE IF EXISTS spare_parts;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS users;

SET FOREIGN_KEY_CHECKS = 1;

-- 1. MASTER DATA

CREATE TABLE users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(30) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE assets (
    id INT NOT NULL AUTO_INCREMENT,
    asset_code VARCHAR(50) NOT NULL UNIQUE,
    asset_name VARCHAR(100) NOT NULL,
    status VARCHAR(30) NOT NULL,
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

-- 2. TRANSACTIONS & BUSINESS TABLES (depend on users/assets)

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

CREATE TABLE work_orders (
    id INT NOT NULL AUTO_INCREMENT,
    asset_id INT NOT NULL,
    description TEXT NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'Pending',
    requester_id INT NOT NULL,
    engineer_id INT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_work_orders_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE ON UPDATE RESTRICT,
    CONSTRAINT fk_work_orders_requester FOREIGN KEY (requester_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT,
    CONSTRAINT fk_work_orders_engineer FOREIGN KEY (engineer_id) REFERENCES users(id) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE preventive_maintenances (
    id INT NOT NULL AUTO_INCREMENT,
    asset_id INT NOT NULL,
    schedule_type VARCHAR(30) NOT NULL,
    next_run DATETIME NOT NULL,
    checklist_data TEXT,
    status VARCHAR(30) NOT NULL DEFAULT 'Scheduled',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    CONSTRAINT fk_preventive_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. HISTORY & AUDIT

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
    PRIMARY KEY (id),
    CONSTRAINT fk_activity_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- End of recreate script
