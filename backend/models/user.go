package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey;autoIncrement;type:int" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password,omitempty"`
	Name      string    `gorm:"not null" json:"name"`
	Role      string    `gorm:"not null" json:"role"`
	Avatar    string    `gorm:"type:longtext" json:"avatar"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
