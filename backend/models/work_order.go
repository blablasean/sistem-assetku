package models

import "time"

type WorkOrder struct {
	ID          int       `gorm:"primaryKey;autoIncrement"`
	AssetID     int       `gorm:"not null"`
	Description string    `gorm:"not null"`
	Status      string    `gorm:"not null"`
	RequesterID int       `gorm:"not null"`
	EngineerID  int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
