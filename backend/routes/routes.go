package routes

import (
	"encoding/json"
	"net/http"
	"sistem-asetku-backend/controllers"
	"sistem-asetku-backend/middlewares"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// RegisterRoutes registers all API routes
func RegisterRoutes(
	db *gorm.DB,
	authCtrl *controllers.AuthController,
	assetCtrl *controllers.AssetController,
	mutationCtrl *controllers.MutationController,
	workOrderCtrl *controllers.WorkOrderController,
	maintenanceCtrl *controllers.MaintenanceController,
) *http.ServeMux {
	mux := http.NewServeMux()
	// Public endpoints
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/auth/login", handleLogin(authCtrl))
	mux.HandleFunc("/auth/register", handleRegister(db))

	// Protected endpoints (single handler per path; switch by method)
	mux.Handle("/assets", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetAssets(assetCtrl)(w, r)
		case http.MethodPost:
			handleCreateAsset(assetCtrl)(w, r)
		case http.MethodPut:
			handleEditAsset(assetCtrl)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	mux.Handle("/assets/code", middlewares.AuthMiddleware(handleGetAssetByCode(assetCtrl)))
	mux.Handle("/assets/reserve", middlewares.AuthMiddleware(handleReserveAsset(assetCtrl)))

	// asset detail and history - use prefix "/assets/"
	mux.Handle("/assets/", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tail := strings.TrimPrefix(r.URL.Path, "/assets/")
		if tail == "" {
			http.NotFound(w, r)
			return
		}
		if strings.HasSuffix(tail, "/history") && r.Method == http.MethodGet {
			idStr := strings.TrimSuffix(tail, "/history")
			q := r.URL.Query()
			q.Set("id", idStr)
			r.URL.RawQuery = q.Encode()
			handleGetAssetHistory(mutationCtrl)(w, r)
			return
		}
		if r.Method == http.MethodGet {
			q := r.URL.Query()
			q.Set("id", tail)
			r.URL.RawQuery = q.Encode()
			handleGetAssetDetail(assetCtrl)(w, r)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})))

	// Mutations
	mux.Handle("/mutations", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateMutation(mutationCtrl)(w, r)
		case http.MethodGet:
			handleGetMutationHistory(mutationCtrl)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Workorders
	mux.Handle("/workorders/assign", middlewares.AuthMiddleware(handleAssignWorker(workOrderCtrl)))
	mux.Handle("/workorders/status", middlewares.AuthMiddleware(handleUpdateWorkOrderStatus(workOrderCtrl)))
	mux.Handle("/workorders/close", middlewares.AuthMiddleware(handleCloseWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders/cancel", middlewares.AuthMiddleware(handleCancelWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateWorkOrder(workOrderCtrl)(w, r)
		case http.MethodGet:
			handleGetWorkOrders(workOrderCtrl)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Maintenance
	mux.Handle("/maintenance/schedule", middlewares.AuthMiddleware(handleCreatePMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tail := strings.TrimPrefix(r.URL.Path, "/maintenance/")
		if strings.HasSuffix(tail, "/checklist") && r.Method == http.MethodPost {
			idStr := strings.TrimSuffix(tail, "/checklist")
			q := r.URL.Query()
			q.Set("id", idStr)
			r.URL.RawQuery = q.Encode()
			handleSubmitPMChecklist(maintenanceCtrl)(w, r)
			return
		}
		if strings.HasSuffix(tail, "/history") && r.Method == http.MethodGet {
			idStr := strings.TrimSuffix(tail, "/history")
			q := r.URL.Query()
			q.Set("assetID", idStr)
			r.URL.RawQuery = q.Encode()
			handleGetMaintenanceHistory(maintenanceCtrl)(w, r)
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})))

	// Activity logs
	mux.Handle("/activitylogs", middlewares.AuthMiddleware(handleGetActivityLogs(db)))

	return mux
}

// Handlers

func handleHome(w http.ResponseWriter, r *http.Request) {
	utils.SendSuccess(w, http.StatusOK, "Sistem AsetKu backend running", map[string]string{
		"version": "1.0.0",
		"status":  "active",
	})
}

func handleLogin(authCtrl *controllers.AuthController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		token, err := authCtrl.Login(payload.Username, payload.Password)
		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, "Login failed", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Login successful", map[string]string{
			"token": token,
		})
	}
}

func handleRegister(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Name     string `json:"name"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		hashedPassword, err := utils.HashPassword(payload.Password)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to hash password", err.Error())
			return
		}

		user := models.User{
			Username: payload.Username,
			Password: hashedPassword,
			Name:     payload.Name,
			Role:     payload.Role,
		}

		if err := db.Create(&user).Error; err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to create user", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "User registered successfully", user)
	}
}

func handleGetAssets(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		assets, err := assetCtrl.SearchAndFilterAssets(q)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch assets", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Assets retrieved successfully", assets)
	}
}

func handleCreateAsset(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload models.Asset
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := assetCtrl.RegistrasiAsset(payload, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to create asset", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Asset created successfully", payload)
	}
}

func handleEditAsset(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		idStr := r.URL.Query().Get("id")
		assetID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", err.Error())
			return
		}

		var payload models.Asset
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := assetCtrl.EditAsset(assetID, payload, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to edit asset", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset updated successfully", payload)
	}
}

func handleGetAssetByCode(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			utils.SendError(w, http.StatusBadRequest, "Asset code parameter required", "code is missing")
			return
		}
		asset, err := assetCtrl.GetAssetByCode(code)
		if err != nil {
			utils.SendError(w, http.StatusNotFound, "Asset not found", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Asset retrieved successfully", asset)
	}
}

func handleReserveAsset(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "management"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			AssetID    int  `json:"asset_id"`
			IsReserved bool `json:"is_reserved"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := assetCtrl.ReserveAsset(payload.AssetID, payload.IsReserved, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to update reservation", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset reservation updated", payload)
	}
}

func handleGetAssetDetail(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		assetID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", err.Error())
			return
		}

		asset, err := assetCtrl.GetAssetDetail(assetID)
		if err != nil {
			utils.SendError(w, http.StatusNotFound, "Asset not found", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset retrieved successfully", asset)
	}
}

func handleGetAssetHistory(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		assetID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", err.Error())
			return
		}

		history, err := mutationCtrl.GetLocationHistory(assetID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch history", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset history retrieved successfully", history)
	}
}

func handleCreateMutation(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload models.Mutation
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := mutationCtrl.CreateMutation(payload, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to create mutation", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Mutation created successfully", payload)
	}
}

func handleGetMutationHistory(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("assetID")
		assetID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", err.Error())
			return
		}

		history, err := mutationCtrl.GetLocationHistory(assetID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch mutation history", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Mutation history retrieved successfully", history)
	}
}

func handleCreateWorkOrder(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "external"
		userID := 1
		if claims != nil {
			if claims.Role != "" {
				role = claims.Role
			}
			if claims.UserID > 0 {
				userID = claims.UserID
			}
		}

		var payload models.WorkOrder
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if payload.RequesterID == 0 {
			payload.RequesterID = userID
		}

		if err := workOrderCtrl.CreateWorkOrder(payload, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to create work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Work order created successfully", payload)
	}
}

func handleGetWorkOrders(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := workOrderCtrl.GetAllWorkOrders()
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch work orders", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Work orders retrieved successfully", orders)
	}
}

func handleAssignWorker(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			WOID       int `json:"wo_id"`
			EngineerID int `json:"engineer_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.AssignWorker(payload.WOID, payload.EngineerID, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to assign worker", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Worker assigned successfully", payload)
	}
}

func handleUpdateWorkOrderStatus(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "engineer"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			WOID        int    `json:"wo_id"`
			Status      string `json:"status"`
			ActionTaken string `json:"action_taken"`
			Cost        int    `json:"cost"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.UpdateWOStatus(payload.WOID, payload.Status, payload.ActionTaken, payload.Cost, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to update status", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order status updated successfully", payload)
	}
}

func handleCancelWorkOrder(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			WOID int `json:"wo_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.CancelWorkOrder(payload.WOID, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to cancel work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order cancelled successfully", payload)
	}
}

func handleCloseWorkOrder(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			WOID int `json:"wo_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.CloseWorkOrder(payload.WOID, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to close work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order closed successfully", payload)
	}
}

func handleGetActivityLogs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logs []models.ActivityLog
		if err := db.Order("id desc").Limit(50).Find(&logs).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch activity logs", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Activity logs retrieved successfully", logs)
	}
}

func handleGetWorkOrderStatus(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		woID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid work order ID", err.Error())
			return
		}

		status := workOrderCtrl.GetWorkOrderStatus(woID)
		utils.SendSuccess(w, http.StatusOK, "Work order status retrieved successfully", map[string]string{
			"status": status,
		})
	}
}

func handleCreatePMSchedule(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.Role != "hod" {
			utils.SendError(w, http.StatusForbidden, "Access denied", "Only HOD can create PM schedules")
			return
		}

		var payload models.PreventiveMaintenance
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := maintenanceCtrl.CreatePMSchedule(payload, claims.Role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to create PM schedule", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "PM schedule created successfully", payload)
	}
}

func handleSubmitPMChecklist(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil {
			utils.SendError(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
			return
		}

		idStr := r.PathValue("id")
		pmID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid PM ID", err.Error())
			return
		}

		var payload struct {
			Checklist string `json:"checklist"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := maintenanceCtrl.SubmitPMChecklist(pmID, payload.Checklist, claims.Role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to submit checklist", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "PM checklist submitted successfully", map[string]int{
			"pm_id": pmID,
		})
	}
}

func handleGetMaintenanceHistory(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("assetID")
		assetID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", err.Error())
			return
		}

		history, err := maintenanceCtrl.GetMaintenanceHistory(assetID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch maintenance history", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Maintenance history retrieved successfully", history)
	}
}
