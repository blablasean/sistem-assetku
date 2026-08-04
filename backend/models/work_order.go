package models

import "time"

type WorkOrder struct {
	ID          int        `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetID     int        `gorm:"not null;type:int" json:"asset_id"`
	Category    string     `gorm:"default:'HVAC / AC'" json:"category"`
	Location    string     `gorm:"default:''" json:"location"`
	Priority    string     `gorm:"default:'Medium'" json:"priority"` // Low, Medium, High, Emergency
	Description string     `gorm:"not null" json:"description"`
	Status      string     `gorm:"not null;default:'Open'" json:"status"` // Open, In Progress, Under Review, Completed, Closed, Cancelled
	RequesterID int        `gorm:"not null;type:int" json:"requester_id"`
	RequestedBy string     `gorm:"default:''" json:"requested_by"`
	Department  string     `gorm:"default:''" json:"department"`
	EngineerID  int        `gorm:"default:0;type:int" json:"engineer_id"`
	ActionTaken      string     `gorm:"default:''" json:"action_taken"`
	AlasanPembatalan string     `gorm:"column:alasan_pembatalan;type:text;default:''" json:"alasan_pembatalan"`
	Cost             int        `gorm:"default:0" json:"cost"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
	ClosedAt         *time.Time `json:"closed_at"`
}

type WorkOrderLog struct {
	ID          int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	WorkOrderID int       `gorm:"not null;type:int;index" json:"work_order_id"`
	Status      string    `gorm:"not null" json:"status"`
	ActionTaken string    `gorm:"default:''" json:"action_taken"`
	Cost        int       `gorm:"default:0" json:"cost"`
	UpdatedBy   string    `gorm:"default:''" json:"updated_by"`
	UserRole    string    `gorm:"default:''" json:"user_role"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
