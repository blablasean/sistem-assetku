package models

import "time"

type Mutation struct {
	ID               int       `gorm:"primaryKey;autoIncrement"`
	AssetID          int       `gorm:"not null"`
	PreviousLocation string    `gorm:"not null"`
	NewLocation      string    `gorm:"not null"`
	PICID            int       `gorm:"not null"`
	MutationDate     time.Time `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
