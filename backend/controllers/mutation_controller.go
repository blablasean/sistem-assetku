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
	if !canMutate(callerRole) {
		return errors.New("akses ditolak: hanya Admin atau HOD yang dapat membuat mutasi lokasi aset")
	}

	now := time.Now()
	if data.MutationDate.IsZero() {
		data.MutationDate = now
	}

	var asset models.Asset
	if err := c.db.First(&asset, data.AssetID).Error; err != nil {
		return errors.New("aset tidak ditemukan di database")
	}

	// Ensure registration location is permanently preserved
	if asset.RegistrationLocation == "" {
		asset.RegistrationLocation = asset.Location
	}

	data.PreviousLocation = asset.Location

	// Update asset current location and last_moved_at timestamp
	asset.Location = data.NewLocation
	asset.LastMovedAt = &now
	if err := c.db.Save(&asset).Error; err != nil {
		return err
	}

	// Save legacy mutation record
	if err := c.db.Create(&data).Error; err != nil {
		return err
	}

	// Save entry to dedicated asset_mutation_timelines table
	tl := models.AssetMutationTimeline{
		AssetCode:        asset.AssetCode,
		PreviousLocation: data.PreviousLocation,
		NewLocation:      data.NewLocation,
		PIC:              data.PIC,
		Reason:           data.Reason,
		MovedAt:          data.MutationDate,
	}
	c.db.Create(&tl)

	return nil
}

func (c *MutationController) CreateMutationByCode(assetCode string, newLocation string, pic string, reason string, callerRole string) error {
	if !canMutate(callerRole) {
		return errors.New("akses ditolak: hanya Admin atau HOD yang dapat membuat mutasi lokasi aset")
	}

	var asset models.Asset
	if err := c.db.Where("asset_code = ?", assetCode).First(&asset).Error; err != nil {
		return errors.New("aset dengan kode " + assetCode + " tidak ditemukan")
	}

	now := time.Now()
	if asset.RegistrationLocation == "" {
		asset.RegistrationLocation = asset.Location
	}

	prevLoc := asset.Location
	asset.Location = newLocation
	asset.LastMovedAt = &now
	if err := c.db.Save(&asset).Error; err != nil {
		return err
	}

	// Save legacy mutation record
	mut := models.Mutation{
		AssetID:          asset.ID,
		PreviousLocation: prevLoc,
		NewLocation:      newLocation,
		PIC:              pic,
		Reason:           reason,
		MutationDate:     now,
	}
	c.db.Create(&mut)

	// Save entry to dedicated asset_mutation_timelines table
	tl := models.AssetMutationTimeline{
		AssetCode:        assetCode,
		PreviousLocation: prevLoc,
		NewLocation:      newLocation,
		PIC:              pic,
		Reason:           reason,
		MovedAt:          now,
	}
	return c.db.Create(&tl).Error
}

func (c *MutationController) GetLocationHistory(assetID int) ([]models.Mutation, error) {
	var history []models.Mutation
	c.db.Where("asset_id = ?", assetID).Order("mutation_date DESC").Find(&history)
	return history, nil
}

func (c *MutationController) GetAssetMutationTimeline(assetCode string) ([]models.AssetMutationTimeline, error) {
	var timelines []models.AssetMutationTimeline
	c.db.Where("asset_code = ?", assetCode).Order("moved_at asc").Find(&timelines)

	// Fallback & Auto-populate: If no records exist yet in timelines table for this asset, construct initial registration item!
	if len(timelines) == 0 {
		var asset models.Asset
		if err := c.db.Where("asset_code = ?", assetCode).First(&asset).Error; err == nil {
			regLoc := asset.RegistrationLocation
			if regLoc == "" {
				regLoc = asset.Location
			}
			t1 := models.AssetMutationTimeline{
				AssetCode:        asset.AssetCode,
				PreviousLocation: "-",
				NewLocation:      regLoc,
				PIC:              asset.PIC,
				Reason:           "Registrasi awal aset terdaftar di sistem",
				MovedAt:          asset.CreatedAt,
			}
			c.db.Create(&t1)
			timelines = append(timelines, t1)

			if asset.Location != regLoc {
				movedTime := asset.CreatedAt.Add(24 * time.Hour)
				if asset.LastMovedAt != nil {
					movedTime = *asset.LastMovedAt
				}
				t2 := models.AssetMutationTimeline{
					AssetCode:        asset.AssetCode,
					PreviousLocation: regLoc,
					NewLocation:      asset.Location,
					PIC:              asset.PIC,
					Reason:           "Mutasi posisi aset ke lokasi " + asset.Location,
					MovedAt:          movedTime,
				}
				c.db.Create(&t2)
				timelines = append(timelines, t2)
			}
		}
	}
	return timelines, nil
}

func (c *MutationController) GetAllAssetMutationTimelines() ([]models.AssetMutationTimeline, error) {
	var list []models.AssetMutationTimeline
	c.db.Order("moved_at desc").Limit(100).Find(&list)

	// If empty, auto-populate all existing assets
	if len(list) == 0 {
		var assets []models.Asset
		c.db.Find(&assets)
		for _, a := range assets {
			c.GetAssetMutationTimeline(a.AssetCode)
		}
		c.db.Order("moved_at desc").Limit(100).Find(&list)
	}
	return list, nil
}
