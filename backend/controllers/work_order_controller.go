package controllers

import (
	"errors"
	"strconv"
	"time"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"

	"gorm.io/gorm"
)

type WorkOrderController struct {
	db *gorm.DB
}

func NewWorkOrderController(db *gorm.DB) *WorkOrderController {
	return &WorkOrderController{db: db}
}

// logWO is a helper to insert a row into work_order_logs and optionally activity_logs
func (c *WorkOrderController) logWO(woID int, status string, actionTaken string, cost int, updatedBy string, callerRole string) {
	wol := models.WorkOrderLog{
		WorkOrderID: woID,
		Status:      status,
		ActionTaken: actionTaken,
		Cost:        cost,
		UpdatedBy:   updatedBy,
		UserRole:    callerRole,
	}
	c.db.Create(&wol)
}

func (c *WorkOrderController) CreateWorkOrder(data *models.WorkOrder, callerRole string) error {
	if data.Priority == "" {
		data.Priority = "Medium"
	}
	data.Status = "Open"
	if err := c.db.Create(data).Error; err != nil {
		return err
	}

	// Log to work_order_logs
	c.logWO(data.ID, "Open", "Laporan kerusakan diajukan: "+data.Description, 0, data.RequestedBy, data.Department)

	// Record in activity_logs
	utils.RecordActivity(c.db, "WORK_ORDER", data.RequestedBy,
		"Membuat Work Order #"+strconv.Itoa(data.ID)+": "+data.Description,
		strconv.Itoa(data.ID))

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

	// Log to work_order_logs
	c.logWO(woID, "In Progress",
		"Penugasan Teknisi #"+strconv.Itoa(engineerID)+" untuk perbaikan lokasi "+workOrder.Location,
		workOrder.Cost, updatedBy, callerRole)

	// Record in activity_logs
	utils.RecordActivity(c.db, "WORK_ORDER", updatedBy,
		"Menugaskan Teknisi #"+strconv.Itoa(engineerID)+" ke Work Order #"+strconv.Itoa(woID)+" ("+workOrder.Location+")",
		strconv.Itoa(woID))

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

	// Log to work_order_logs
	c.logWO(woID, status, actionTaken, cost, updatedBy, callerRole)

	// Record in activity_logs
	logMsg := "Update status Work Order #" + strconv.Itoa(woID) + " → " + status
	if actionTaken != "" {
		logMsg += ": " + actionTaken
	}
	utils.RecordActivity(c.db, "WORK_ORDER", updatedBy, logMsg, strconv.Itoa(woID))

	return nil
}

func (c *WorkOrderController) CancelWorkOrder(woID int, callerRole string, updatedBy string) error {
	if !canCloseOrDeleteWO(callerRole) {
		return errors.New("akses ditolak: hanya Admin, HOD, atau Supervisor (Management) yang dapat membatalkan Work Order")
	}
	if err := c.db.Model(&models.WorkOrder{}).Where("id = ?", woID).Update("status", "Cancelled").Error; err != nil {
		return err
	}

	// Log to work_order_logs
	c.logWO(woID, "Cancelled", "Tiket Work Order dibatalkan oleh "+updatedBy, 0, updatedBy, callerRole)

	// Record in activity_logs
	utils.RecordActivity(c.db, "WORK_ORDER", updatedBy,
		"Membatalkan Work Order #"+strconv.Itoa(woID),
		strconv.Itoa(woID))

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

	// Log to work_order_logs
	c.logWO(woID, "Finish", "Work Order terverifikasi selesai", 0, updatedBy, callerRole)

	// Record in activity_logs
	utils.RecordActivity(c.db, "WORK_ORDER", updatedBy,
		"Menutup & memverifikasi selesai Work Order #"+strconv.Itoa(woID),
		strconv.Itoa(woID))

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
	// Cascade delete logs
	c.db.Where("work_order_id = ?", woID).Delete(&models.WorkOrderLog{})

	// Record in activity_logs
	utils.RecordActivity(c.db, "WORK_ORDER", "admin",
		"Menghapus Work Order #"+strconv.Itoa(woID),
		strconv.Itoa(woID))

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
	if err := c.db.Save(&workOrder).Error; err != nil {
		return err
	}

	// Record in activity_logs
	utils.RecordActivity(c.db, "WORK_ORDER", callerRole,
		"Mengedit detail Work Order #"+strconv.Itoa(woID),
		strconv.Itoa(woID))

	return nil
}

// GetWorkOrderLogs returns all work_order_logs for a specific WO, with auto-populate fallback
func (c *WorkOrderController) GetWorkOrderLogs(woID int) ([]models.WorkOrderLog, error) {
	var logs []models.WorkOrderLog
	c.db.Where("work_order_id = ?", woID).Order("created_at asc").Find(&logs)

	// Fallback: if no logs yet, auto-populate from WO data
	if len(logs) == 0 {
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

			l1 := models.WorkOrderLog{
				WorkOrderID: wo.ID,
				Status:      "Open",
				ActionTaken: "Laporan kerusakan diajukan: " + wo.Description,
				UpdatedBy:   reqBy,
				UserRole:    dept,
				CreatedAt:   wo.CreatedAt,
			}
			c.db.Create(&l1)
			logs = append(logs, l1)

			if wo.Status != "Open" {
				l2 := models.WorkOrderLog{
					WorkOrderID: wo.ID,
					Status:      "In Progress",
					ActionTaken: "Penugasan Teknisi untuk perbaikan lokasi " + wo.Location,
					UpdatedBy:   "HOD Engineer",
					UserRole:    "hod",
					CreatedAt:   wo.CreatedAt.Add(15 * time.Minute),
				}
				c.db.Create(&l2)
				logs = append(logs, l2)
			}

			if wo.Status == "Under Review" || wo.Status == "Finish" || wo.Status == "Completed" || wo.Status == "Closed" {
				act := wo.ActionTaken
				if act == "" {
					act = "Perbaikan unit selesai. Menunggu review."
				}
				l3 := models.WorkOrderLog{
					WorkOrderID: wo.ID,
					Status:      "Under Review",
					ActionTaken: act,
					Cost:        wo.Cost,
					UpdatedBy:   "Teknisi",
					UserRole:    "engineer",
					CreatedAt:   wo.CreatedAt.Add(45 * time.Minute),
				}
				c.db.Create(&l3)
				logs = append(logs, l3)
			}

			if wo.Status == "Finish" || wo.Status == "Completed" || wo.Status == "Closed" || wo.Status == "Cancelled" {
				fTime := wo.CreatedAt.Add(90 * time.Minute)
				if wo.ClosedAt != nil {
					fTime = *wo.ClosedAt
				}
				fAct := "Work Order diverifikasi selesai"
				if wo.Status == "Cancelled" {
					fAct = "Work Order dibatalkan"
				}
				l4 := models.WorkOrderLog{
					WorkOrderID: wo.ID,
					Status:      wo.Status,
					ActionTaken: fAct,
					Cost:        wo.Cost,
					UpdatedBy:   "Administrator",
					UserRole:    "admin",
					CreatedAt:   fTime,
				}
				c.db.Create(&l4)
				logs = append(logs, l4)
			}
		}
	}
	return logs, nil
}

// GetAllWorkOrderLogs returns all work_order_logs, auto-populating if empty
func (c *WorkOrderController) GetAllWorkOrderLogs() ([]models.WorkOrderLog, error) {
	var list []models.WorkOrderLog
	c.db.Order("created_at desc").Limit(200).Find(&list)

	if len(list) == 0 {
		var wos []models.WorkOrder
		c.db.Find(&wos)
		for _, wo := range wos {
			c.GetWorkOrderLogs(wo.ID)
		}
		c.db.Order("created_at desc").Limit(200).Find(&list)
	}
	return list, nil
}

// GetWorkOrderTimeline returns work_order_logs for a specific WO
func (c *WorkOrderController) GetWorkOrderTimeline(woID int) ([]models.WorkOrderLog, error) {
	return c.GetWorkOrderLogs(woID)
}

// GetAllTimelines returns all work_order_logs
func (c *WorkOrderController) GetAllTimelines() ([]models.WorkOrderLog, error) {
	return c.GetAllWorkOrderLogs()
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

	// Log to work_order_logs
	c.logWO(woID, wo.Status, actionTaken, cost, updatedBy, callerRole)

	// Record in activity_logs
	logMsg := "Menambah progress Work Order #" + strconv.Itoa(woID)
	if status != "" {
		logMsg += " → " + status
	}
	if actionTaken != "" {
		logMsg += ": " + actionTaken
	}
	utils.RecordActivity(c.db, "WORK_ORDER", updatedBy, logMsg, strconv.Itoa(woID))

	return nil
}
