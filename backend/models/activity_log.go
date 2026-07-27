package models

import "time"

type ActivityLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	Category  string    `gorm:"type:varchar(50);default:'GENERAL'" json:"category"`
	UserID    int       `gorm:"type:int;default:0" json:"user_id"`
	Actor     string    `gorm:"type:varchar(255);default:''" json:"actor"`
	Action    string    `gorm:"not null;type:text" json:"action"`
	EntityID  string    `gorm:"type:varchar(255);default:''" json:"entity_id"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}
