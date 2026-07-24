package controllers

import (
	"errors"
	"sistem-asetku-backend/models"
	"time"

	"gorm.io/gorm"
)

type MutationController struct {
	db *gorm.DB
}

func NewMutationController(db *gorm.DB) *MutationController {
	return &MutationController{db: db}
}

func (c *MutationController) CreateMutation(data models.Mutation, callerRole string) error {
	// Only admin and hod can create mutations
	if !canMutate(callerRole) {
		return errors.New("akses ditolak: hanya Admin atau HOD yang dapat membuat mutasi lokasi aset")
	}

	// Auto-fill mutation date if missing
	if data.MutationDate.IsZero() {
		data.MutationDate = time.Now()
	}

	// Fetch current asset location to fill previous_location automatically
	if data.PreviousLocation == "" {
		var asset models.Asset
		if err := c.db.First(&asset, data.AssetID).Error; err == nil {
			data.PreviousLocation = asset.Location
		}
	}

	// Save mutation record
	if err := c.db.Create(&data).Error; err != nil {
		return err
	}

	// Update the asset's current location to new location
	if err := c.db.Model(&models.Asset{}).Where("id = ?", data.AssetID).Update("location", data.NewLocation).Error; err != nil {
		return err
	}

	return nil
}

func (c *MutationController) GetLocationHistory(assetID int) ([]models.Mutation, error) {
	var history []models.Mutation
	if err := c.db.Where("asset_id = ?", assetID).Order("mutation_date DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}
