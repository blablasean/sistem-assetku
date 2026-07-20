package controllers

import (
	"errors"
	"sistem-asetku-backend/models"

	"gorm.io/gorm"
)

type MutationController struct {
	db *gorm.DB
}

func NewMutationController(db *gorm.DB) *MutationController {
	return &MutationController{db: db}
}

func (c *MutationController) CreateMutation(data models.Mutation, callerRole string) error {
	if callerRole != "hod" {
		return errors.New("akses ditolak")
	}
	return c.db.Create(&data).Error
}

func (c *MutationController) GetLocationHistory(assetID int) ([]models.Mutation, error) {
	var history []models.Mutation
	if err := c.db.Where("asset_id = ?", assetID).Order("mutation_date DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}
