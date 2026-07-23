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

	schedule.ChecklistData = checklist
	schedule.Status = "Completed"
	return c.db.Save(&schedule).Error
}

func (c *MaintenanceController) GetMaintenanceHistory(assetID int) ([]models.MaintenanceHistory, error) {
	var history []models.MaintenanceHistory
	if err := c.db.Where("asset_id = ?", assetID).Order("created_at DESC").Find(&history).Error; err != nil {
		return nil, err
	}
	return history, nil
}
