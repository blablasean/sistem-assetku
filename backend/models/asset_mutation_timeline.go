package models

import "time"

type AssetMutationTimeline struct {
	ID               int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetCode        string    `gorm:"not null;type:varchar(255);index" json:"asset_code"`
	PreviousLocation string    `gorm:"default:''" json:"previous_location"`
	NewLocation      string    `gorm:"default:''" json:"new_location"`
	PIC              string    `gorm:"default:''" json:"pic"`
	Reason           string    `gorm:"type:text" json:"reason"`
	MovedAt          time.Time `gorm:"autoCreateTime" json:"moved_at"`
}
