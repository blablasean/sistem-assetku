package controllers

import (
	"errors"
	"strings"
	"time"
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
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat membuat jadwal Preventive Maintenance")
	}
	return c.db.Create(&schedule).Error
}

func (c *MaintenanceController) EditPMSchedule(pmID int, updated models.PreventiveMaintenance, callerRole string) error {
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat mengubah jadwal Preventive Maintenance")
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
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menghapus jadwal Preventive Maintenance")
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

func (c *MaintenanceController) SubmitPMChecklist(pmID int, targetDate string, checklist string, callerRole string) error {
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menyelesaikan checklist maintenance")
	}
	var schedule models.PreventiveMaintenance
	if err := c.db.First(&schedule, pmID).Error; err != nil {
		return err
	}

	if checklist != "" {
		schedule.ChecklistData = checklist
	}

	cleanDate := strings.TrimSpace(targetDate)
	if strings.Contains(cleanDate, "T") {
		cleanDate = strings.Split(cleanDate, "T")[0]
	}
	if cleanDate == "" {
		cleanDate = time.Now().Format("2006-01-02")
	}

	if schedule.CompletedDates == "" {
		schedule.CompletedDates = cleanDate
	} else {
		existing := strings.Split(schedule.CompletedDates, ",")
		alreadyExists := false
		for _, d := range existing {
			if strings.TrimSpace(d) == cleanDate {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			schedule.CompletedDates = schedule.CompletedDates + "," + cleanDate
		}
	}

	schedule.Status = "Completed"
	if err := c.db.Save(&schedule).Error; err != nil {
		return err
	}

	actMsg := "Perawatan Berkala (" + schedule.ScheduleType + ") Tanggal " + cleanDate + ": " + schedule.ChecklistData

	history := models.MaintenanceHistory{
		AssetID:     schedule.AssetID,
		ActionTaken: actMsg,
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
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat mengubah riwayat maintenance")
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
	if !canManageAssets(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menghapus riwayat maintenance")
	}
	res := c.db.Where("id = ?", historyID).Delete(&models.MaintenanceHistory{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("riwayat maintenance tidak ditemukan di database")
	}
	return nil
}
