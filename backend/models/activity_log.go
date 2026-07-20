package models

import "time"

type ActivityLog struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	UserID    int       `gorm:"not null"`
	Action    string    `gorm:"not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}
