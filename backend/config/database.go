package config

import (
	"fmt"
	"log"
	"os"
	"time"

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

	// Ensure avatar column is LONGTEXT to hold Base64 images
	db.Exec("ALTER TABLE users MODIFY COLUMN avatar LONGTEXT NULL;")

	log.Println("✓ Database models auto-migrated successfully")

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

	// Update ALL existing users in database so their password becomes 'admin123'
	db.Model(&models.User{}).Where("1 = 1").Update("password", hashedPassword)

	// Initial Default Users for All 5 System Roles (Password: admin123)
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
			db.Create(&u)
		} else {
			db.Model(&existing).Updates(map[string]interface{}{
				"password": hashedPassword,
				"name":     u.Name,
				"role":     u.Role,
			})
		}
	}

	// Check if assets exist
	var assetCount int64
	db.Model(&models.Asset{}).Count(&assetCount)
	if assetCount == 0 {
		assets := []models.Asset{
			{AssetCode: "AST-RM301-AC", AssetName: "AC Split Daikin 1.5 PK", Category: "HVAC / AC", Location: "Kamar 301", PIC: "Deni (Tech)", Status: "Active"},
			{AssetCode: "AST-RM102-TV", AssetName: "Smart TV LG 43 Inch", Category: "Elektronik & TV", Location: "Kamar 102", PIC: "Front Desk Team", Status: "Active"},
			{AssetCode: "AST-KCH-CHILLER", AssetName: "Chiller Dapur Utama", Category: "Kitchen Equipment", Location: "Kitchen Dapur", PIC: "Kitchen Chef", Status: "Maintenance"},
			{AssetCode: "AST-GEN-01", AssetName: "Generator Unit Cummins 500kVA", Category: "Mesin & Generator", Location: "Power House", PIC: "Engineering Supervisor", Status: "Active", IsReserved: true},
			{AssetCode: "AST-LBY-SOFA", AssetName: "Set Sofa Premium Leather", Category: "Mebel & Furniture", Location: "Lobby Lounge", PIC: "Housekeeping", Status: "Active"},
		}
		for _, a := range assets {
			db.Create(&a)
		}
	}

	// Check if work orders exist
	var woCount int64
	db.Model(&models.WorkOrder{}).Count(&woCount)
	if woCount == 0 {
		workOrders := []models.WorkOrder{
			{AssetID: 1, Category: "HVAC / AC", Location: "Kamar 301", Priority: "Emergency", Description: "AC Kamar 301 bocor air dan tidak dingin", Status: "In Progress", RequesterID: 4, EngineerID: 3},
			{AssetID: 3, Category: "Kitchen Equipment", Location: "Kitchen Dapur", Priority: "High", Description: "Chiller Dapur Utama suhu naik ke -5°C", Status: "Open", RequesterID: 4, EngineerID: 0},
			{AssetID: 2, Category: "Elektronik & TV", Location: "Kamar 102", Priority: "Medium", Description: "Smart TV HDMI port tidak terdeteksi", Status: "Closed", RequesterID: 4, EngineerID: 3},
		}
		for _, wo := range workOrders {
			db.Create(&wo)
		}
	}

	// Check if PM schedules exist
	var pmCount int64
	db.Model(&models.PreventiveMaintenance{}).Count(&pmCount)
	if pmCount == 0 {
		pmList := []models.PreventiveMaintenance{
			{
				AssetID:       1,
				ScheduleType:  "Monthly",
				NextRun:       time.Now().AddDate(0, 1, 0),
				ChecklistData: "1. Cek tekanan freon AC\n2. Bersihkan filter evaporator\n3. Cek drainase air kondensasi",
				Status:        "Active",
			},
			{
				AssetID:       4,
				ScheduleType:  "Weekly",
				NextRun:       time.Now().AddDate(0, 0, 7),
				ChecklistData: "1. Tes running generator 15 menit\n2. Cek level solar tangki harian\n3. Ukur tegangan aki starter",
				Status:        "Active",
			},
		}
		for _, pm := range pmList {
			db.Create(&pm)
		}
	}

	log.Println("✓ Initial default users and hotel sample data seeded & verified successfully")
	return nil
}

// AutoMigrateModels automatically migrates all database models
func AutoMigrateModels() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Asset{},
		&models.SparePart{},
		&models.Mutation{},
		&models.WorkOrder{},
		&models.PreventiveMaintenance{},
		&models.MaintenanceHistory{},
		&models.ActivityLog{},
	)
}

// GetDB returns the database connection
func GetDB() *gorm.DB {
	return DB
}
