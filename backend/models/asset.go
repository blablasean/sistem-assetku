package models

import "time"

type Asset struct {
	ID          int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	AssetCode   string    `gorm:"unique;not null" json:"asset_code"`
	AssetName   string    `gorm:"not null" json:"asset_name"`
	Category    string    `gorm:"default:'General'" json:"category"`
	Location    string    `gorm:"default:'Main Store'" json:"location"`
	PIC         string    `gorm:"default:''" json:"pic"`
	Status      string    `gorm:"not null;default:'Active'" json:"status"` // Active, Maintenance, Damaged, Retired, Reserved
	DocumentURL string    `gorm:"default:''" json:"document_url"`
	IsReserved  bool      `gorm:"default:false" json:"is_reserved"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
