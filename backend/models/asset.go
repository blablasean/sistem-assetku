package models

import "time"

type Asset struct {
	ID        int       `gorm:"primaryKey;autoIncrement;type:int"`
	AssetCode string    `gorm:"unique;not null"`
	AssetName string    `gorm:"not null"`
	Status    string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
