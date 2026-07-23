package controllers

import (
	"errors"
	"time"
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
	if data.Priority == "" {
		data.Priority = "Medium"
	}
	data.Status = "Open"
	return c.db.Create(&data).Error
}

func (c *WorkOrderController) GetAllWorkOrders() ([]models.WorkOrder, error) {
	var list []models.WorkOrder
	if err := c.db.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (c *WorkOrderController) AssignWorker(woID int, engineerID int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" {
		return errors.New("akses ditolak: hanya HOD atau Management yang dapat menugaskan teknisi")
	}
	var workOrder models.WorkOrder
	if err := c.db.First(&workOrder, woID).Error; err != nil {
		return err
	}
	workOrder.EngineerID = engineerID
	workOrder.Status = "In Progress"
	return c.db.Save(&workOrder).Error
}

func (c *WorkOrderController) UpdateWOStatus(woID int, status string, actionTaken string, cost int, callerRole string) error {
	var workOrder models.WorkOrder
	if err := c.db.First(&workOrder, woID).Error; err != nil {
		return err
	}
	workOrder.Status = status
	if actionTaken != "" {
		workOrder.ActionTaken = actionTaken
	}
	if cost > 0 {
		workOrder.Cost = cost
	}
	return c.db.Save(&workOrder).Error
}

func (c *WorkOrderController) CancelWorkOrder(woID int, callerRole string) error {
	// hod, management, engineer, and external roles can cancel open work orders
	return c.db.Model(&models.WorkOrder{}).Where("id = ?", woID).Update("status", "Cancelled").Error
}

func (c *WorkOrderController) CloseWorkOrder(woID int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat menyelesaikan Work Order")
	}
	now := time.Now()
	return c.db.Model(&models.WorkOrder{}).Where("id = ?", woID).Updates(map[string]interface{}{
		"status":    "Finish",
		"closed_at": &now,
	}).Error
}

func (c *WorkOrderController) DeleteWorkOrder(woID int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat menghapus Work Order")
	}
	res := c.db.Where("id = ?", woID).Delete(&models.WorkOrder{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("work order tidak ditemukan di database")
	}
	return nil
}

func (c *WorkOrderController) GetWorkOrderStatus(woID int) string {
	var workOrder models.WorkOrder
	if err := c.db.First(&workOrder, woID).Error; err != nil {
		return "unknown"
	}
	return workOrder.Status
}

func (c *WorkOrderController) EditWorkOrderDetail(woID int, description string, actionTaken string, cost int, callerRole string) error {
	if callerRole != "hod" && callerRole != "management" && callerRole != "admin" {
		return errors.New("akses ditolak: hanya HOD, Management, atau Admin yang dapat merubah detail Work Order")
	}
	var workOrder models.WorkOrder
	if err := c.db.First(&workOrder, woID).Error; err != nil {
		return err
	}
	if description != "" {
		workOrder.Description = description
	}
	if actionTaken != "" {
		workOrder.ActionTaken = actionTaken
	}
	if cost >= 0 {
		workOrder.Cost = cost
	}
	return c.db.Save(&workOrder).Error
}
