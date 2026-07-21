package models

import "time"

type ActivityLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement;type:int"`
	UserID    int       `gorm:"not null;type:int"`
	Action    string    `gorm:"not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}
