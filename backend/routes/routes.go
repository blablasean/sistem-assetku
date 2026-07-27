package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sistem-asetku-backend/controllers"
	"sistem-asetku-backend/middlewares"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"

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
	mux.Handle("/auth/profile", middlewares.AuthMiddleware(handleUpdateProfile(db)))

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
	mux.Handle("/assets/delete", middlewares.AuthMiddleware(handleDeleteAsset(assetCtrl)))

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

	mux.Handle("/assets/mutation-timeline", middlewares.AuthMiddleware(handleGetAssetMutationTimeline(mutationCtrl)))
	mux.Handle("/assets/mutate", middlewares.AuthMiddleware(handleMutateAssetByCode(mutationCtrl)))

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
	mux.Handle("/workorders/edit", middlewares.AuthMiddleware(handleEditWorkOrderDetail(workOrderCtrl)))
	mux.Handle("/workorders/close", middlewares.AuthMiddleware(handleCloseWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders/cancel", middlewares.AuthMiddleware(handleCancelWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders/delete", middlewares.AuthMiddleware(handleDeleteWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders/logs", middlewares.AuthMiddleware(handleGetWorkOrderTimeline(workOrderCtrl)))
	mux.Handle("/workorders/timeline/add", middlewares.AuthMiddleware(handleAddTimelineEntry(workOrderCtrl)))
	mux.Handle("/workorders/timeline", middlewares.AuthMiddleware(handleGetWorkOrderTimeline(workOrderCtrl)))
	mux.Handle("/workorders", middlewares.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateWorkOrder(db, workOrderCtrl)(w, r)
		case http.MethodGet:
			handleGetWorkOrders(workOrderCtrl)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Maintenance
	mux.Handle("/maintenance/schedules", middlewares.AuthMiddleware(handleGetAllPMSchedules(maintenanceCtrl)))
	mux.Handle("/maintenance/schedule", middlewares.AuthMiddleware(handleCreatePMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/edit", middlewares.AuthMiddleware(handleEditPMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/delete", middlewares.AuthMiddleware(handleDeletePMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/history/edit", middlewares.AuthMiddleware(handleEditMaintenanceHistory(maintenanceCtrl)))
	mux.Handle("/maintenance/history/delete", middlewares.AuthMiddleware(handleDeleteMaintenanceHistory(maintenanceCtrl)))
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
	mux.Handle("/activitylogs", middlewares.AuthMiddleware(handleGetActivityLogs(db, workOrderCtrl, mutationCtrl)))

	// User Management endpoints (Admin)
	mux.Handle("/users/create", middlewares.AuthMiddleware(handleCreateUser(db)))
	mux.Handle("/users/edit", middlewares.AuthMiddleware(handleEditUser(db)))
	mux.Handle("/users/delete", middlewares.AuthMiddleware(handleDeleteUser(db)))
	mux.Handle("/users", middlewares.AuthMiddleware(handleGetAllUsers(db)))

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

		token, user, err := authCtrl.Login(payload.Username, payload.Password)
		if err != nil {
			utils.SendError(w, http.StatusUnauthorized, "Login gagal: Username atau password salah", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Login successful", map[string]interface{}{
			"token":    token,
			"role":     user.Role,
			"name":     user.Name,
			"username": user.Username,
			"id":       user.ID,
			"avatar":   user.Avatar,
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

func handleDeleteAsset(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			AssetID int `json:"asset_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := assetCtrl.DeleteAsset(payload.AssetID, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to delete asset", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset deleted successfully", payload)
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
		role := "external"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload models.Mutation
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := mutationCtrl.CreateMutation(payload, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to create mutation", err.Error())
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

func handleGetAssetMutationTimeline(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assetCode := r.URL.Query().Get("asset_code")
		if assetCode == "" {
			utils.SendError(w, http.StatusBadRequest, "Asset code required", "asset_code query parameter missing")
			return
		}

		timelines, err := mutationCtrl.GetAssetMutationTimeline(assetCode)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch mutation timeline", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset mutation timeline retrieved successfully", timelines)
	}
}

func handleMutateAssetByCode(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			AssetCode   string `json:"asset_code"`
			NewLocation string `json:"new_location"`
			PIC         string `json:"pic"`
			Reason      string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := mutationCtrl.CreateMutationByCode(payload.AssetCode, payload.NewLocation, payload.PIC, payload.Reason, role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to record mutation", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Mutasi aset berhasil dicatat", payload)
	}
}

func handleCreateWorkOrder(db *gorm.DB, workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "external"
		userID := 1
		username := "user_hotel"
		if claims != nil {
			if claims.Role != "" {
				role = claims.Role
			}
			if claims.UserID > 0 {
				userID = claims.UserID
			}
			if claims.Username != "" {
				username = claims.Username
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

		// Always fetch the exact user record from DB to guarantee accurate username & department
		if userID > 0 {
			var user models.User
			if err := db.First(&user, userID).Error; err == nil {
				if payload.RequestedBy == "" || payload.RequestedBy == "user_hotel" {
					payload.RequestedBy = user.Username
				}
				if payload.Department == "" || payload.Department == "external" {
					payload.Department = user.Role
				}
				role = user.Role
			}
		}

		if payload.RequestedBy == "" {
			payload.RequestedBy = username
		}
		if payload.Department == "" {
			payload.Department = role
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
		username := "User"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}
		if claims != nil && claims.Username != "" {
			username = claims.Username
		}

		var payload struct {
			WOID       int `json:"wo_id"`
			EngineerID int `json:"engineer_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.AssignWorker(payload.WOID, payload.EngineerID, role, username); err != nil {
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
		username := "User"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}
		if claims != nil && claims.Username != "" {
			username = claims.Username
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

		if err := workOrderCtrl.UpdateWOStatus(payload.WOID, payload.Status, payload.ActionTaken, payload.Cost, role, username); err != nil {
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
		username := "User"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}
		if claims != nil && claims.Username != "" {
			username = claims.Username
		}

		var payload struct {
			WOID int `json:"wo_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.CancelWorkOrder(payload.WOID, role, username); err != nil {
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
		username := "User"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}
		if claims != nil && claims.Username != "" {
			username = claims.Username
		}

		var payload struct {
			WOID int `json:"wo_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.CloseWorkOrder(payload.WOID, role, username); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to close work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order closed successfully", payload)
	}
}

func handleDeleteWorkOrder(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
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

		if err := workOrderCtrl.DeleteWorkOrder(payload.WOID, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to delete work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order deleted successfully", payload)
	}
}

func handleGetWorkOrderTimeline(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("wo_id")
		woID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid work order ID", err.Error())
			return
		}

		timelines, err := workOrderCtrl.GetWorkOrderTimeline(woID)
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch work order timeline", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order timeline retrieved successfully", timelines)
	}
}

func handleAddTimelineEntry(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.SendError(w, http.StatusMethodNotAllowed, "Method not allowed", "POST required")
			return
		}
		claims := middlewares.GetClaimsFromContext(r)
		role := "engineer"
		username := "User"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}
		if claims != nil && claims.Username != "" {
			username = claims.Username
		}

		var payload struct {
			WOID        int    `json:"work_order_id"`
			Status      string `json:"status"`
			ActionTaken string `json:"action_taken"`
			Cost        int    `json:"cost"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.AddTimelineEntry(payload.WOID, payload.Status, payload.ActionTaken, payload.Cost, role, username); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to add timeline entry", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Catatan timeline berhasil ditambahkan", payload)
	}
}

func handleGetActivityLogs(db *gorm.DB, workOrderCtrl *controllers.WorkOrderController, mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Fetch finished / completed Work Orders
		var finishedWOs []models.WorkOrder
		db.Where("status IN ('Finish','Completed','Closed')").Order("id desc").Limit(50).Find(&finishedWOs)

		// Fetch all Work Order Timelines
		timelines, _ := workOrderCtrl.GetAllTimelines()

		// Fetch all Asset Mutation Timelines
		assetMutationTimelines, _ := mutationCtrl.GetAllAssetMutationTimelines()

		// Fetch all maintenance history entries
		var maintenanceHistory []models.MaintenanceHistory
		db.Order("id desc").Limit(50).Find(&maintenanceHistory)

		// Fetch all mutation history
		var mutations []models.Mutation
		db.Order("mutation_date desc").Limit(100).Find(&mutations)

		result := map[string]interface{}{
			"work_orders":              finishedWOs,
			"timelines":                timelines,
			"asset_mutation_timelines": assetMutationTimelines,
			"maintenance_history":      maintenanceHistory,
			"mutations":                mutations,
		}
		utils.SendSuccess(w, http.StatusOK, "Activity logs retrieved successfully", result)
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
		role := "admin"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var input struct {
			AssetID       interface{} `json:"asset_id"`
			ScheduleType  string      `json:"schedule_type"`
			NextRun       string      `json:"next_run"`
			ChecklistData string      `json:"checklist_data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		assetID := 0
		switch v := input.AssetID.(type) {
		case float64:
			assetID = int(v)
		case int:
			assetID = v
		case string:
			assetID, _ = strconv.Atoi(v)
		}

		nextRunTime := time.Now()
		if input.NextRun != "" {
			if t, err := time.Parse("2006-01-02", input.NextRun); err == nil {
				nextRunTime = t
			} else if t, err := time.Parse(time.RFC3339, input.NextRun); err == nil {
				nextRunTime = t
			}
		}

		payload := models.PreventiveMaintenance{
			AssetID:       assetID,
			ScheduleType:  input.ScheduleType,
			NextRun:       nextRunTime,
			ChecklistData: input.ChecklistData,
			Status:        "Active",
		}

		if err := maintenanceCtrl.CreatePMSchedule(payload, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to create PM schedule", err.Error())
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
		if idStr == "" {
			idStr = r.URL.Query().Get("id")
		}
		pmID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid PM ID", err.Error())
			return
		}

		var payload struct {
			TargetDate string `json:"target_date"`
			Checklist  string `json:"checklist"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := maintenanceCtrl.SubmitPMChecklist(pmID, payload.TargetDate, payload.Checklist, claims.Role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to submit checklist", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "PM checklist submitted successfully", map[string]interface{}{
			"pm_id":       pmID,
			"target_date": payload.TargetDate,
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

func handleEditPMSchedule(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			PMID          int    `json:"pm_id"`
			ScheduleType  string `json:"schedule_type"`
			ChecklistData string `json:"checklist_data"`
			NextRun       string `json:"next_run"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		var parsedNextRun time.Time
		if payload.NextRun != "" {
			parsedNextRun, _ = time.Parse("2006-01-02", payload.NextRun)
		}

		updated := models.PreventiveMaintenance{
			ScheduleType:  payload.ScheduleType,
			ChecklistData: payload.ChecklistData,
			NextRun:       parsedNextRun,
		}

		if err := maintenanceCtrl.EditPMSchedule(payload.PMID, updated, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to edit PM schedule", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "PM schedule updated successfully", payload)
	}
}

func handleDeletePMSchedule(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			PMID int `json:"pm_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := maintenanceCtrl.DeletePMSchedule(payload.PMID, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to delete PM schedule", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "PM schedule deleted successfully", payload)
	}
}

func handleGetAllPMSchedules(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		schedules, err := maintenanceCtrl.GetAllPMSchedules()
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch PM schedules", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "PM schedules retrieved successfully", schedules)
	}
}

func handleEditWorkOrderDetail(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			WOID        int    `json:"wo_id"`
			Description string `json:"description"`
			ActionTaken string `json:"action_taken"`
			Cost        int    `json:"cost"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.EditWorkOrderDetail(payload.WOID, payload.Description, payload.ActionTaken, payload.Cost, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to edit Work Order detail", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work Order detail updated successfully", payload)
	}
}

func handleEditMaintenanceHistory(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			HistoryID   int    `json:"history_id"`
			ActionTaken string `json:"action_taken"`
			Cost        int    `json:"cost"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := maintenanceCtrl.EditMaintenanceHistory(payload.HistoryID, payload.ActionTaken, payload.Cost, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to edit Maintenance History", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Maintenance History updated successfully", payload)
	}
}

func handleDeleteMaintenanceHistory(maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "hod"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			HistoryID int `json:"history_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := maintenanceCtrl.DeleteMaintenanceHistory(payload.HistoryID, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to delete Maintenance History", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Maintenance History deleted successfully", payload)
	}
}

func handleGetAllUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var users []models.User
		if err := db.Order("id desc").Find(&users).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch users", err.Error())
			return
		}
		for i := range users {
			users[i].Password = ""
		}
		utils.SendSuccess(w, http.StatusOK, "Users retrieved successfully", users)
	}
}

func handleCreateUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.Role != "admin" {
			utils.SendError(w, http.StatusForbidden, "Akses ditolak: Hanya Admin yang dapat menambah pengguna baru", "Access denied")
			return
		}

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

		if strings.TrimSpace(payload.Username) == "" || strings.TrimSpace(payload.Password) == "" || strings.TrimSpace(payload.Name) == "" {
			utils.SendError(w, http.StatusBadRequest, "Data pengguna tidak lengkap", "Username, Password, and Name required")
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
			utils.SendError(w, http.StatusBadRequest, "Gagal membuat pengguna (username mungkin sudah digunakan)", err.Error())
			return
		}

		user.Password = ""
		utils.SendSuccess(w, http.StatusCreated, "Pengguna baru berhasil dibuat!", user)
	}
}

func handleEditUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.Role != "admin" {
			utils.SendError(w, http.StatusForbidden, "Akses ditolak: Hanya Admin yang dapat mengubah data pengguna", "Access denied")
			return
		}

		var payload struct {
			UserID   int    `json:"user_id"`
			Username string `json:"username"`
			Password string `json:"password"`
			Name     string `json:"name"`
			Role     string `json:"role"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		var user models.User
		if err := db.First(&user, payload.UserID).Error; err != nil {
			utils.SendError(w, http.StatusNotFound, "Pengguna tidak ditemukan", err.Error())
			return
		}

		user.Name = payload.Name
		user.Role = payload.Role
		if strings.TrimSpace(payload.Username) != "" {
			user.Username = payload.Username
		}
		if strings.TrimSpace(payload.Password) != "" {
			hashedPassword, err := utils.HashPassword(payload.Password)
			if err == nil {
				user.Password = hashedPassword
			}
		}

		if err := db.Save(&user).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Gagal menyimpan perubahan pengguna", err.Error())
			return
		}

		user.Password = ""
		utils.SendSuccess(w, http.StatusOK, "Data pengguna berhasil diperbarui!", user)
	}
}

func handleDeleteUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.Role != "admin" {
			utils.SendError(w, http.StatusForbidden, "Akses ditolak: Hanya Admin yang dapat menghapus pengguna", "Access denied")
			return
		}

		var payload struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if payload.UserID == claims.UserID {
			utils.SendError(w, http.StatusBadRequest, "Anda tidak dapat menghapus akun Anda sendiri yang sedang aktif!", "Self deletion prohibited")
			return
		}

		if err := db.Where("id = ?", payload.UserID).Delete(&models.User{}).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Gagal menghapus pengguna", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Pengguna berhasil dihapus permanen!", payload)
	}
}

func handleUpdateProfile(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.UserID == 0 {
			utils.SendError(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
			return
		}

		var payload struct {
			Avatar string `json:"avatar"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			utils.SendError(w, http.StatusNotFound, "Pengguna tidak ditemukan", err.Error())
			return
		}

		user.Avatar = payload.Avatar

		if err := db.Save(&user).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Gagal memperbarui foto profil", err.Error())
			return
		}

		user.Password = ""
		utils.SendSuccess(w, http.StatusOK, "Foto profil Anda berhasil diperbarui!", user)
	}
}
