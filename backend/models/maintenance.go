package models

import "time"

type PreventiveMaintenance struct {
	ID             int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetID        int       `gorm:"not null;type:int" json:"asset_id"`
	ScheduleType   string    `gorm:"not null" json:"schedule_type"`
	NextRun        time.Time `gorm:"not null" json:"next_run"`
	ChecklistData  string    `gorm:"type:text" json:"checklist_data"`
	CompletedDates string    `gorm:"type:text" json:"completed_dates"`
	Status         string    `gorm:"not null;default:'Active'" json:"status"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type MaintenanceHistory struct {
	ID          int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetID     int       `gorm:"not null;type:int" json:"asset_id"`
	ActionTaken string    `gorm:"not null" json:"action_taken"`
	Cost        int       `gorm:"not null" json:"cost"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
