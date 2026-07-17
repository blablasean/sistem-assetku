package controllers

import (
	"errors"

	"sistem-asetku-backend/models"

	"gorm.io/gorm"
)

type AssetController struct {
	db *gorm.DB
}

func NewAssetController(db *gorm.DB) *AssetController {
	return &AssetController{db: db}
}

func (c *AssetController) RegistrasiAsset(newAsset models.Asset, callerRole string) error {
	if err := validateHodRole(callerRole); err != nil {
		return err
	}

	return c.executeAssetRegistration(newAsset)
}

func validateHodRole(callerRole string) error {
	if callerRole != "hod" {
		return errors.New("Akses ditolak")
	}
	return nil
}

func (c *AssetController) executeAssetRegistration(asset models.Asset) error {
	return c.db.Create(&asset).Error
}
