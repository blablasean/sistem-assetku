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
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat mendaftarkan aset baru")
	}
	return c.db.Create(&newAsset).Error
}

func (c *AssetController) EditAsset(assetID int, updated models.Asset, callerRole string) error {
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat mengubah data aset")
	}
	var existing models.Asset
	if err := c.db.First(&existing, assetID).Error; err != nil {
		return err
	}
	existing.AssetName = updated.AssetName
	existing.Category = updated.Category
	existing.Location = updated.Location
	existing.PIC = updated.PIC
	existing.Status = updated.Status
	existing.DocumentURL = updated.DocumentURL
	return c.db.Save(&existing).Error
}

func (c *AssetController) GenerateQRCode(assetCode string, callerRole string) (string, error) {
	return "QR_CODE:" + assetCode, nil
}

func (c *AssetController) GetAssetDetail(assetID int) (models.Asset, error) {
	var asset models.Asset
	if err := c.db.First(&asset, assetID).Error; err != nil {
		return models.Asset{}, err
	}
	return asset, nil
}

func (c *AssetController) GetAssetByCode(assetCode string) (models.Asset, error) {
	var asset models.Asset
	if err := c.db.Where("asset_code = ?", assetCode).First(&asset).Error; err != nil {
		return models.Asset{}, err
	}
	return asset, nil
}

func (c *AssetController) SearchAndFilterAssets(query string) ([]models.Asset, error) {
	var assets []models.Asset
	if query == "" {
		if err := c.db.Order("id desc").Find(&assets).Error; err != nil {
			return nil, err
		}
		return assets, nil
	}
	search := "%" + query + "%"
	if err := c.db.Where("asset_name LIKE ? OR asset_code LIKE ? OR status LIKE ? OR location LIKE ? OR category LIKE ?", search, search, search, search, search).Order("id desc").Find(&assets).Error; err != nil {
		return nil, err
	}
	return assets, nil
}

func (c *AssetController) ReserveAsset(assetID int, isReserved bool, callerRole string) error {
	var asset models.Asset
	if err := c.db.First(&asset, assetID).Error; err != nil {
		return err
	}
	asset.IsReserved = isReserved
	if isReserved {
		asset.Status = "Reserved"
	} else if asset.Status == "Reserved" {
		asset.Status = "Active"
	}
	return c.db.Save(&asset).Error
}

func (c *AssetController) DeleteAsset(assetID int, callerRole string) error {
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menghapus aset")
	}
	res := c.db.Where("id = ?", assetID).Delete(&models.Asset{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("aset tidak ditemukan di database")
	}
	return nil
}
