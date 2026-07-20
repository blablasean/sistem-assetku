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

	return c.db.Create(&newAsset).Error
}

func (c *AssetController) GenerateQRCode(assetCode string, callerRole string) (string, error) {
	if err := validateHodRole(callerRole); err != nil {
		return "", err
	}
	return "QR_CODE_PLACEHOLDER:" + assetCode, nil
}

func (c *AssetController) GetAssetDetail(assetID int) (models.Asset, error) {
	var asset models.Asset
	if err := c.db.First(&asset, assetID).Error; err != nil {
		return models.Asset{}, err
	}
	return asset, nil
}

func (c *AssetController) SearchAndFilterAssets(query string) ([]models.Asset, error) {
	var assets []models.Asset
	search := "%" + query + "%"
	if err := c.db.Where("asset_name LIKE ? OR asset_code LIKE ? OR status LIKE ?", search, search, search).Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func validateHodRole(callerRole string) error {
	if callerRole != "hod" {
		return errors.New("akses ditolak")
	}
	return nil
}
