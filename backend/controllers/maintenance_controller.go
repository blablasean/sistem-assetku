package controllers

import (
	"errors"
	"sistem-asetku-backend/models"

	"gorm.io/gorm"
)

type MaintenanceController struct {
	db *gorm.DB
}

func NewMaintenanceController(db *gorm.DB) *MaintenanceController {
	return &MaintenanceController{db: db}
}

func (c *MaintenanceController) CreatePMSchedule(schedule models.PreventiveMaintenance, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat membuat jadwal PM")
	}
	return c.db.Create(&schedule).Error
}

func (c *MaintenanceController) EditPMSchedule(pmID int, updated models.PreventiveMaintenance, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat mengubah jadwal PM")
	}
	var schedule models.PreventiveMaintenance
	if err := c.db.First(&schedule, pmID).Error; err != nil {
		return err
	}
	if updated.ScheduleType != "" {
		schedule.ScheduleType = updated.ScheduleType
	}
	if updated.ChecklistData != "" {
		schedule.ChecklistData = updated.ChecklistData
	}
	if !updated.NextRun.IsZero() {
		schedule.NextRun = updated.NextRun
	}
	return c.db.Save(&schedule).Error
}

func (c *MaintenanceController) DeletePMSchedule(pmID int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat menghapus jadwal PM")
	}
	res := c.db.Where("id = ?", pmID).Delete(&models.PreventiveMaintenance{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("jadwal PM tidak ditemukan di database")
	}
	return nil
}

func (c *MaintenanceController) SubmitPMChecklist(pmID int, checklist string, callerRole string) error {
	var schedule models.PreventiveMaintenance
	if err := c.db.First(&schedule, pmID).Error; err != nil {
		return err
	}

	if checklist != "" {
		schedule.ChecklistData = checklist
	}
	schedule.Status = "Completed"
	if err := c.db.Save(&schedule).Error; err != nil {
		return err
	}

	history := models.MaintenanceHistory{
		AssetID:     schedule.AssetID,
		ActionTaken: "Perawatan Berkala (" + schedule.ScheduleType + "): " + schedule.ChecklistData,
		Cost:        50000,
	}
	return c.db.Create(&history).Error
}

func (c *MaintenanceController) GetMaintenanceHistory(assetID int) ([]models.MaintenanceHistory, error) {
	var history []models.MaintenanceHistory
	if err := c.db.Where("asset_id = ?", assetID).Order("created_at DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}

func (c *MaintenanceController) GetAllPMSchedules() ([]models.PreventiveMaintenance, error) {
	var schedules []models.PreventiveMaintenance
	if err := c.db.Order("id desc").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (c *MaintenanceController) EditMaintenanceHistory(historyID int, actionTaken string, cost int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat merubah riwayat maintenance")
	}
	var mh models.MaintenanceHistory
	if err := c.db.First(&mh, historyID).Error; err != nil {
		return err
	}
	if actionTaken != "" {
		mh.ActionTaken = actionTaken
	}
	if cost >= 0 {
		mh.Cost = cost
	}
	return c.db.Save(&mh).Error
}

func (c *MaintenanceController) DeleteMaintenanceHistory(historyID int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat menghapus riwayat maintenance")
	}
	res := c.db.Where("id = ?", historyID).Delete(&models.MaintenanceHistory{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("riwayat maintenance tidak ditemukan")
	}
	return nil
}
