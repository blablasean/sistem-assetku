package models

import "time"

type PreventiveMaintenance struct {
	ID            int       `gorm:"primaryKey;autoIncrement;type:int"`
	AssetID       int       `gorm:"not null;type:int"`
	ScheduleType  string    `gorm:"not null"`
	NextRun       time.Time `gorm:"not null"`
	ChecklistData string    `gorm:"type:text"`
	Status        string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

type MaintenanceHistory struct {
	ID          int       `gorm:"primaryKey;autoIncrement;type:int"`
	AssetID     int       `gorm:"not null;type:int"`
	ActionTaken string    `gorm:"not null"`
	Cost        int       `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
