package controllers

import (
	"errors"
	"sistem-asetku-backend/models"

	"gorm.io/gorm"
)

type WorkOrderController struct {
	db *gorm.DB
}

func NewWorkOrderController(db *gorm.DB) *WorkOrderController {
	return &WorkOrderController{db: db}
}

func (c *WorkOrderController) CreateWorkOrder(data models.WorkOrder, callerRole string) error {
	if callerRole != "hod" {
		return errors.New("akses ditolak")
	}
	data.Status = "Open"
	return c.db.Create(&data).Error
}

func (c *WorkOrderController) AssignWorker(woID int, engineerID int, callerRole string) error {
	if err := validateHodRole(callerRole); err != nil {
		return err
	}
	return c.assignWorkerToOrder(woID, engineerID)
}

func (c *WorkOrderController) assignWorkerToOrder(woID int, engineerID int) error {
	tx := c.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var workOrder models.WorkOrder
	if err := tx.First(&workOrder, woID).Error; err != nil {
		tx.Rollback()
		return err
	}

	workOrder.EngineerID = engineerID
	workOrder.Status = "In Progress"

	if err := tx.Save(&workOrder).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (c *WorkOrderController) UpdateWOStatus(woID int, status string, callerRole string) error {
	if err := validateHodRole(callerRole); err != nil {
		return err
	}

	return c.db.Model(&models.WorkOrder{}).Where("id = ?", woID).Update("status", status).Error
}

func (c *WorkOrderController) GetWorkOrderStatus(woID int) string {
	var workOrder models.WorkOrder
	if err := c.db.First(&workOrder, woID).Error; err != nil {
		return "unknown"
	}
	return workOrder.Status
}
