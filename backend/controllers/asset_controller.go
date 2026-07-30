package controllers

import (
	"errors"

	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
	"strconv"

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
	if err := c.db.Create(&newAsset).Error; err != nil {
		return err
	}

	// Record in activity log
	utils.RecordActivity(c.db, "ASET", callerRole,
		"Mendaftarkan aset baru: "+newAsset.AssetName+" ("+newAsset.AssetCode+") di "+newAsset.Location,
		newAsset.AssetCode)

	return nil
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
	if err := c.db.Save(&existing).Error; err != nil {
		return err
	}

	// Record in activity log
	utils.RecordActivity(c.db, "ASET", callerRole,
		"Mengedit data aset: "+existing.AssetName+" ("+existing.AssetCode+") → Status: "+existing.Status+", Lokasi: "+existing.Location,
		existing.AssetCode)

	return nil
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
	if err := c.db.Save(&asset).Error; err != nil {
		return err
	}

	// Record in activity log
	action := "Mereservasi"
	if !isReserved {
		action = "Membatalkan reservasi"
	}
	utils.RecordActivity(c.db, "ASET", callerRole,
		action+" aset: "+asset.AssetName+" ("+asset.AssetCode+")",
		asset.AssetCode)

	return nil
}

func (c *AssetController) DeleteAsset(assetID int, callerRole string) error {
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menghapus aset")
	}
	var existing models.Asset
	if err := c.db.First(&existing, assetID).Error; err != nil {
		return errors.New("aset tidak ditemukan di database")
	}
	res := c.db.Where("id = ?", assetID).Delete(&models.Asset{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("aset tidak ditemukan di database")
	}

	// Record in activity log
	utils.RecordActivity(c.db, "ASET", callerRole,
		"Menghapus aset: "+existing.AssetName+" ("+existing.AssetCode+")",
		existing.AssetCode)

	return nil
}

// ChangeAssetStatus updates just the status field of an asset and records the change
func (c *AssetController) ChangeAssetStatus(assetID int, newStatus string, callerRole string) error {
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat mengubah status aset")
	}
	var asset models.Asset
	if err := c.db.First(&asset, assetID).Error; err != nil {
		return errors.New("aset tidak ditemukan")
	}
	oldStatus := asset.Status
	asset.Status = newStatus
	if err := c.db.Save(&asset).Error; err != nil {
		return err
	}

	// Record in activity log
	utils.RecordActivity(c.db, "ASET", callerRole,
		"Mengubah status aset "+asset.AssetName+" ("+asset.AssetCode+"): "+oldStatus+" → "+newStatus,
		strconv.Itoa(assetID))

	return nil
}
