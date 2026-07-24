package models

import "time"

type Mutation struct {
	ID               int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetID          int       `gorm:"not null;type:int" json:"asset_id"`
	PreviousLocation string    `gorm:"default:''" json:"previous_location"`
	NewLocation      string    `gorm:"not null" json:"new_location"`
	PIC              string    `gorm:"default:''" json:"new_pic"`
	Reason           string    `gorm:"default:''" json:"reason"`
	MutationDate     time.Time `gorm:"not null" json:"mutation_date"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}
