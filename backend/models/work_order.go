package models

import "time"

type WorkOrder struct {
	ID          int        `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetID     int        `gorm:"not null;type:int" json:"asset_id"`
	Location    string     `gorm:"default:''" json:"location"`
	Priority    string     `gorm:"default:'Medium'" json:"priority"` // Low, Medium, High, Emergency
	Description string     `gorm:"not null" json:"description"`
	Status      string     `gorm:"not null;default:'Open'" json:"status"` // Open, In Progress, Under Review, Completed, Closed, Cancelled
	RequesterID int        `gorm:"not null;type:int" json:"requester_id"`
	EngineerID  int        `gorm:"default:0;type:int" json:"engineer_id"`
	ActionTaken string     `gorm:"default:''" json:"action_taken"`
	Cost        int        `gorm:"default:0" json:"cost"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ClosedAt    *time.Time `json:"closed_at"`
}
