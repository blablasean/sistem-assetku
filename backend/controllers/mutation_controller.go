package controllers

import (
	"errors"
	"strconv"
	"strings"
	"time"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"

	"gorm.io/gorm"
)

type MutationController struct {
	db *gorm.DB
}

func NewMutationController(db *gorm.DB) *MutationController {
	return &MutationController{db: db}
}

type MutationInput struct {
	AssetID          int       `json:"asset_id"`
	PreviousLocation string    `json:"previous_location"`
	NewLocation      string    `json:"new_location"`
	PIC              string    `json:"pic"`
	Reason           string    `json:"reason"`
	MutationDate     time.Time `json:"mutation_date"`
}

func (c *MutationController) CreateMutation(data MutationInput, callerRole string) error {
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

	// Save entry to dedicated asset_mutation_timelines table
	tl := models.AssetMutationTimeline{
		AssetCode:        asset.AssetCode,
		PreviousLocation: data.PreviousLocation,
		NewLocation:      data.NewLocation,
		PIC:              data.PIC,
		Reason:           data.Reason,
		MovedAt:          data.MutationDate,
	}
	if err := c.db.Create(&tl).Error; err != nil {
		return err
	}

	// Record in activity log
	utils.RecordActivity(c.db, "MUTASI_ASET", data.PIC,
		"Mutasi aset "+asset.AssetCode+" ("+asset.AssetName+"): "+data.PreviousLocation+" → "+data.NewLocation+". Alasan: "+data.Reason,
		asset.AssetCode)

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

	// Save entry to dedicated asset_mutation_timelines table
	tl := models.AssetMutationTimeline{
		AssetCode:        assetCode,
		PreviousLocation: prevLoc,
		NewLocation:      newLocation,
		PIC:              pic,
		Reason:           reason,
		MovedAt:          now,
	}
	if err := c.db.Create(&tl).Error; err != nil {
		return err
	}

	// Record in activity log
	utils.RecordActivity(c.db, "MUTASI_ASET", pic,
		"Mutasi aset "+assetCode+": "+prevLoc+" → "+newLocation+". Alasan: "+reason,
		assetCode)

	return nil
}

func (c *MutationController) GetLocationHistory(assetID int) ([]models.AssetMutationTimeline, error) {
	var asset models.Asset
	if err := c.db.First(&asset, assetID).Error; err != nil {
		return nil, err
	}
	return c.GetAssetMutationTimeline(asset.AssetCode)
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
	return list, nil
}

// GetMutationByID retrieves a single asset mutation timeline by its integer ID
func (c *MutationController) GetMutationByID(id int) (models.AssetMutationTimeline, error) {
	var mut models.AssetMutationTimeline
	if err := c.db.First(&mut, id).Error; err != nil {
		return models.AssetMutationTimeline{}, err
	}
	return mut, nil
}

// DeleteMutation removes an asset mutation timeline entry
func (c *MutationController) DeleteMutation(id int, callerRole string) error {
	r := strings.ToLower(strings.TrimSpace(callerRole))
	if r == "external" || r == "" {
		return errors.New("akses ditolak: anda harus login terlebih dahulu")
	}
	res := c.db.Where("id = ?", id).Delete(&models.AssetMutationTimeline{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("data mutasi tidak ditemukan di database")
	}

	// Record in activity log
	utils.RecordActivity(c.db, "MUTASI_ASET", callerRole,
		"Menghapus record mutasi aset #"+strconv.Itoa(id),
		strconv.Itoa(id))

	return nil
}

// EditMutation updates an asset mutation timeline entry
func (c *MutationController) EditMutation(id int, previousLocation string, newLocation string, pic string, reason string, callerRole string) error {
	r := strings.ToLower(strings.TrimSpace(callerRole))
	if r == "external" || r == "" {
		return errors.New("akses ditolak: anda harus login terlebih dahulu")
	}
	var mut models.AssetMutationTimeline
	if err := c.db.First(&mut, id).Error; err != nil {
		return errors.New("data mutasi aset tidak ditemukan")
	}
	mut.PreviousLocation = previousLocation
	mut.NewLocation = newLocation
	mut.PIC = pic
	mut.Reason = reason
	return c.db.Save(&mut).Error
}

