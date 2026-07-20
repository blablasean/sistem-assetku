package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"sistem-asetku-backend/models"
)

var DB *gorm.DB

// InitDatabase initializes a GORM MySQL connection using environment variables.
// If no env vars are set, it defaults to localhost:3306 and database `db_barang`.
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, dbname)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("gorm open: %w", err)
	}

	DB = db

	// Auto-migrate models present in the project
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Asset{},
		&models.WorkOrder{},
		&models.Mutation{},
		&models.PreventiveMaintenance{},
		&models.MaintenanceHistory{},
		&models.SparePart{},
		&models.ActivityLog{},
	); err != nil {
		// do not fail app startup just because auto-migrate had issues on an existing DB
		log.Printf("warning: auto migrate encountered an error: %v", err)
	}

	return nil
}
