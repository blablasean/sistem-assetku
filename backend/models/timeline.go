package models

import "time"

type Timeline struct {
	ID          int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	WorkOrderID int       `gorm:"not null;type:int;index" json:"work_order_id"`
	Status      string    `gorm:"not null" json:"status"`
	ActionTaken string    `gorm:"default:''" json:"action_taken"`
	Cost        int       `gorm:"default:0" json:"cost"`
	UpdatedBy   string    `gorm:"default:''" json:"updated_by"`
	UserRole    string    `gorm:"default:''" json:"user_role"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
