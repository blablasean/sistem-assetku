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
	if callerRole != "hod" {
		return errors.New("akses ditolak")
	}
	return c.db.Create(&schedule).Error
}

func (c *MaintenanceController) SubmitPMChecklist(pmID int, checklist string, callerRole string) error {
	if callerRole != "engineer" && callerRole != "hod" {
		return errors.New("akses ditolak")
	}

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
