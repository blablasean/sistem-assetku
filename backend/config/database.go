package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
)

var DB *gorm.DB

// InitDatabase initializes a GORM MySQL connection using environment variables.
// If no env vars are set, it defaults to localhost:3306 and database `db_sistemasetku`.
func InitDatabase() error {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "127.0.0.1"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "3306"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "root"
	}
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "db_sistemasetku"
	}

	// connection timeout (e.g., "5s") for the MySQL dial
	connectTimeoutStr := os.Getenv("DB_CONNECT_TIMEOUT")
	if connectTimeoutStr == "" {
		connectTimeoutStr = "5s"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&timeout=%s", user, pass, host, port, dbname, connectTimeoutStr)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db

	// Auto-migrate all models
	if err := AutoMigrateModels(); err != nil {
		return fmt.Errorf("failed to auto-migrate models: %w", err)
	}

	// Ensure avatar and session tracking columns exist on users table
	db.Exec("ALTER TABLE users MODIFY COLUMN avatar LONGTEXT NULL;")
	db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS is_logged_in TINYINT(1) NOT NULL DEFAULT 0;")
	db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS active_token TEXT NULL;")
	db.Exec("ALTER TABLE users ADD COLUMN IF NOT EXISTS last_seen_at DATETIME NULL;")

	// Reset all active session flags on server restart so users can log in cleanly after restart
	db.Exec("UPDATE users SET is_logged_in = 0, active_token = '', last_seen_at = NULL;")

	log.Println("✓ Database models auto-migrated successfully & active sessions reset")

	// Seed initial default users and hotel assets if database is empty
	if err := SeedInitialData(db); err != nil {
		log.Printf("⚠️ Warning seeding database: %v\n", err)
	}

	return nil
}

// SeedInitialData seeds default users and initial sample data for hotel assets and work orders
func SeedInitialData(db *gorm.DB) error {
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil || hashedPassword == "" {
		hashedPassword = "$2a$10$aXZN4vu2Nt7mOJR.a5rpSuj1sKaPpVM75B0.YEG5QC/8gJQmGcAwu"
	}

	// Initial Default Users — only created if username does NOT exist yet.
	// Existing accounts are NEVER modified to preserve custom passwords.
	users := []models.User{
		{Username: "admin", Password: hashedPassword, Name: "Administrator Utama Hotel", Role: "admin"},
		{Username: "hod_eng", Password: hashedPassword, Name: "Pak Alex (HOD Engineer)", Role: "hod"},
		{Username: "spv_eng", Password: hashedPassword, Name: "Pak Hendra (Supervisor)", Role: "management"},
		{Username: "teknisi_budi", Password: hashedPassword, Name: "Budi Santoso (Teknisi)", Role: "engineer"},
		{Username: "staff_frontdesk", Password: hashedPassword, Name: "Rina (Staff Front Office)", Role: "dept_frontoffice"},
	}
	for _, u := range users {
		var existing models.User
		if err := db.Where("username = ?", u.Username).First(&existing).Error; err != nil {
			// Only create if the user doesn't exist yet
			db.Create(&u)
		}
		// DO NOT update existing users — preserve their custom passwords & data
	}

	log.Println("✓ Initial default users seeded (existing accounts unchanged)")
	return nil
}

// AutoMigrateModels automatically migrates active database models
func AutoMigrateModels() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Asset{},
		&models.AssetMutationTimeline{},
		&models.WorkOrder{},
		&models.WorkOrderLog{},
		&models.PreventiveMaintenance{},
		&models.MaintenanceHistory{},
		&models.ActivityLog{},
	)
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
