package models

type SparePart struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	PartName string `gorm:"not null"`
	Stock    int    `gorm:"not null;default:0"`
	MinStock int    `gorm:"not null;default:5"`
}
