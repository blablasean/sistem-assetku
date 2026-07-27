package utils

import (
	"sistem-asetku-backend/models"
	"gorm.io/gorm"
)

// RecordActivity records a unified activity log entry into the activity_logs table
func RecordActivity(db *gorm.DB, category string, actor string, action string, entityID string) {
	if db == nil {
		return
	}
	log := models.ActivityLog{
		Category: category,
		Actor:    actor,
		Action:   action,
		EntityID: entityID,
	}
	_ = db.Create(&log)
}
