package models

import "time"

type Mutation struct {
	ID               int       `gorm:"primaryKey;autoIncrement;type:int"`
	AssetID          int       `gorm:"not null;type:int"`
	PreviousLocation string    `gorm:"not null"`
	NewLocation      string    `gorm:"not null"`
	PICID            int       `gorm:"not null;type:int"`
	MutationDate     time.Time `gorm:"not null"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
}
