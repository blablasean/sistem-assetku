package models

import "time"

type User struct {
    ID        int       `gorm:"primaryKey;autoIncrement"`
    Username  string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Name      string    `gorm:"not null"`
    Role      string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Asset struct {
    ID        int       `gorm:"primaryKey;autoIncrement"`
    AssetCode string    `gorm:"unique;not null"`
    AssetName string    `gorm:"not null"`
    Status    string    `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}

type WorkOrder struct {
    ID          int       `gorm:"primaryKey;autoIncrement"`
    AssetID     int       `gorm:"not null"`
    Description string    `gorm:"not null"`
    Status      string    `gorm:"not null"`
    RequesterID int       `gorm:"not null"`
    EngineerID  int       `gorm:""
    CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type Mutation struct {
    ID               int       `gorm:"primaryKey;autoIncrement"`
    AssetID          int       `gorm:"not null"`
    PreviousLocation string    `gorm:"not null"`
    NewLocation      string    `gorm:"not null"`
    PICID            int       `gorm:"not null"`
    MutationDate     time.Time `gorm:"not null"`
    CreatedAt        time.Time `gorm:"autoCreateTime"`
}
