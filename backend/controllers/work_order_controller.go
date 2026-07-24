package controllers

import (
	"errors"
	"strconv"
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
	if err := c.db.Create(&data).Error; err != nil {
		return err
	}

	// Insert into timelines table
	tl := models.Timeline{
		WorkOrderID: data.ID,
		Status:      "Open",
		ActionTaken: "Laporan kerusakan diajukan: " + data.Description,
		Cost:        0,
		UpdatedBy:   data.RequestedBy,
		UserRole:    data.Department,
	}
	c.db.Create(&tl)
	return nil
}

func (c *WorkOrderController) GetAllWorkOrders() ([]models.WorkOrder, error) {
	var list []models.WorkOrder
	if err := c.db.Order("id desc").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (c *WorkOrderController) AssignWorker(woID int, engineerID int, callerRole string, updatedBy string) error {
	if !canAssignEngineerRole(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menugaskan teknisi")
	}
	var workOrder models.WorkOrder
	if err := c.db.First(&workOrder, woID).Error; err != nil {
		return err
	}
	workOrder.EngineerID = engineerID
	workOrder.Status = "In Progress"
	if err := c.db.Save(&workOrder).Error; err != nil {
		return err
	}

	// Insert into timelines table
	tl := models.Timeline{
		WorkOrderID: woID,
		Status:      "In Progress",
		ActionTaken: "Penugasan Teknisi #" + strconv.Itoa(engineerID) + " untuk perbaikan lokasi " + workOrder.Location,
		Cost:        workOrder.Cost,
		UpdatedBy:   updatedBy,
		UserRole:    callerRole,
	}
	c.db.Create(&tl)
	return nil
}

func (c *WorkOrderController) UpdateWOStatus(woID int, status string, actionTaken string, cost int, callerRole string, updatedBy string) error {
	if !canUpdateProgressRole(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, Supervisor, atau Teknisi (Engineer) yang dapat mengupdate progres Work Order")
	}
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
	if status == "Finish" || status == "Completed" || status == "Closed" {
		now := time.Now()
		workOrder.ClosedAt = &now
	}
	if err := c.db.Save(&workOrder).Error; err != nil {
		return err
	}

	// Insert into timelines table
	tl := models.Timeline{
		WorkOrderID: woID,
		Status:      status,
		ActionTaken: actionTaken,
		Cost:        cost,
		UpdatedBy:   updatedBy,
		UserRole:    callerRole,
	}
	c.db.Create(&tl)
	return nil
}

func (c *WorkOrderController) CancelWorkOrder(woID int, callerRole string, updatedBy string) error {
	if !canCloseOrDeleteWO(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat membatalkan Work Order")
	}
	if err := c.db.Model(&models.WorkOrder{}).Where("id = ?", woID).Update("status", "Cancelled").Error; err != nil {
		return err
	}
	tl := models.Timeline{
		WorkOrderID: woID,
		Status:      "Cancelled",
		ActionTaken: "Tiket Work Order dibatalkan oleh " + updatedBy,
		UpdatedBy:   updatedBy,
		UserRole:    callerRole,
	}
	c.db.Create(&tl)
	return nil
}

func (c *WorkOrderController) CloseWorkOrder(woID int, callerRole string, updatedBy string) error {
	if !canCloseOrDeleteWO(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menyelesaikan Work Order")
	}
	now := time.Now()
	if err := c.db.Model(&models.WorkOrder{}).Where("id = ?", woID).Updates(map[string]interface{}{
		"status":    "Finish",
		"closed_at": &now,
	}).Error; err != nil {
		return err
	}

	tl := models.Timeline{
		WorkOrderID: woID,
		Status:      "Finish",
		ActionTaken: "Work Order terverifikasi selesai",
		UpdatedBy:   updatedBy,
		UserRole:    callerRole,
	}
	c.db.Create(&tl)
	return nil
}

func (c *WorkOrderController) DeleteWorkOrder(woID int, callerRole string) error {
	if !canCloseOrDeleteWO(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat menghapus Work Order")
	}
	res := c.db.Where("id = ?", woID).Delete(&models.WorkOrder{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("work order tidak ditemukan di database")
	}
	c.db.Where("work_order_id = ?", woID).Delete(&models.Timeline{})
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
	if !canCloseOrDeleteWO(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat mengubah detail Work Order")
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

func (c *WorkOrderController) GetWorkOrderTimeline(woID int) ([]models.Timeline, error) {
	var timelines []models.Timeline
	c.db.Where("work_order_id = ?", woID).Order("created_at asc").Find(&timelines)

	// Fallback & Auto-populate: If no records exist in timelines table for this WO, insert milestones!
	if len(timelines) == 0 {
		var wo models.WorkOrder
		if err := c.db.First(&wo, woID).Error; err == nil {
			reqBy := wo.RequestedBy
			if reqBy == "" {
				reqBy = "User Hotel"
			}
			dept := wo.Department
			if dept == "" {
				dept = "Hotel Staff"
			}

			// 1. Open
			t1 := models.Timeline{
				WorkOrderID: wo.ID,
				Status:      "Open",
				ActionTaken: "Laporan kerusakan diajukan: " + wo.Description,
				UpdatedBy:   reqBy,
				UserRole:    dept,
				CreatedAt:   wo.CreatedAt,
			}
			c.db.Create(&t1)
			timelines = append(timelines, t1)

			// 2. In Progress
			if wo.Status != "Open" {
				t2 := models.Timeline{
					WorkOrderID: wo.ID,
					Status:      "In Progress",
					ActionTaken: "Penugasan Teknisi untuk perbaikan lokasi " + wo.Location,
					UpdatedBy:   "HOD Engineer",
					UserRole:    "hod",
					CreatedAt:   wo.CreatedAt.Add(15 * time.Minute),
				}
				c.db.Create(&t2)
				timelines = append(timelines, t2)
			}

			// 3. Under Review / Finish
			if wo.Status == "Under Review" || wo.Status == "Finish" || wo.Status == "Completed" || wo.Status == "Closed" {
				act := wo.ActionTaken
				if act == "" {
					act = "Perbaikan unit selesai. Menunggu review."
				}
				t3 := models.Timeline{
					WorkOrderID: wo.ID,
					Status:      "Under Review",
					ActionTaken: act,
					Cost:        wo.Cost,
					UpdatedBy:   "Budi Santoso (Teknisi)",
					UserRole:    "engineer",
					CreatedAt:   wo.CreatedAt.Add(45 * time.Minute),
				}
				c.db.Create(&t3)
				timelines = append(timelines, t3)
			}

			// 4. Finish / Cancelled
			if wo.Status == "Finish" || wo.Status == "Completed" || wo.Status == "Closed" || wo.Status == "Cancelled" {
				fTime := wo.CreatedAt.Add(90 * time.Minute)
				if wo.ClosedAt != nil {
					fTime = *wo.ClosedAt
				}
				fAct := "Work Order diverifikasi selesai"
				if wo.Status == "Cancelled" {
					fAct = "Work Order dibatalkan"
				}
				t4 := models.Timeline{
					WorkOrderID: wo.ID,
					Status:      wo.Status,
					ActionTaken: fAct,
					Cost:        wo.Cost,
					UpdatedBy:   "Administrator",
					UserRole:    "admin",
					CreatedAt:   fTime,
				}
				c.db.Create(&t4)
				timelines = append(timelines, t4)
			}
		}
	}
	return timelines, nil
}

func (c *WorkOrderController) GetAllTimelines() ([]models.Timeline, error) {
	var list []models.Timeline
	c.db.Order("created_at desc").Limit(100).Find(&list)

	// If empty in DB, populate all existing work orders into timelines table
	if len(list) == 0 {
		var wos []models.WorkOrder
		c.db.Find(&wos)
		for _, wo := range wos {
			c.GetWorkOrderTimeline(wo.ID)
		}
		c.db.Order("created_at desc").Limit(100).Find(&list)
	}
	return list, nil
}

func (c *WorkOrderController) AddTimelineEntry(woID int, status string, actionTaken string, cost int, callerRole string, updatedBy string) error {
	if !canUpdateProgressRole(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, Supervisor, atau Teknisi (Engineer) yang dapat menambah catatan timeline Work Order")
	}
	var wo models.WorkOrder
	if err := c.db.First(&wo, woID).Error; err != nil {
		return errors.New("work order tidak ditemukan")
	}

	if status != "" {
		wo.Status = status
	}
	if actionTaken != "" {
		wo.ActionTaken = actionTaken
	}
	if cost > 0 {
		wo.Cost = cost
	}
	if status == "Finish" || status == "Completed" || status == "Closed" {
		now := time.Now()
		wo.ClosedAt = &now
	}
	if err := c.db.Save(&wo).Error; err != nil {
		return err
	}

	tl := models.Timeline{
		WorkOrderID: woID,
		Status:      wo.Status,
		ActionTaken: actionTaken,
		Cost:        cost,
		UpdatedBy:   updatedBy,
		UserRole:    callerRole,
	}
	return c.db.Create(&tl).Error
}
