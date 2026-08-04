package routes

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
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
	userCtrl *controllers.UserController,
) *http.ServeMux {
	mux := http.NewServeMux()
	authMW := middlewares.AuthMiddlewareWithDB(db)

	// Public endpoints
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/auth/login", handleLogin(authCtrl))
	mux.HandleFunc("/auth/register", handleRegister(db))
	mux.Handle("/auth/me", authMW(handleGetMe(db)))
	mux.Handle("/auth/profile", authMW(handleUpdateProfile(db)))
	mux.Handle("/auth/logout", authMW(handleLogout(authCtrl)))

	// Assets
	mux.Handle("/assets/scan", authMW(handleGetAssetByCode(assetCtrl)))
	mux.Handle("/assets/detail", authMW(handleGetAssetDetail(assetCtrl)))
	mux.Handle("/assets/edit", authMW(handleEditAsset(assetCtrl)))
	mux.Handle("/assets/delete", authMW(handleDeleteAsset(assetCtrl)))
	mux.Handle("/assets/reserve", authMW(handleReserveAsset(assetCtrl)))
	mux.Handle("/assets", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateAsset(assetCtrl)(w, r)
		case http.MethodGet:
			handleGetAssets(assetCtrl)(w, r)
		case http.MethodPut:
			handleEditAsset(assetCtrl)(w, r)
		case http.MethodDelete:
			handleDeleteAsset(assetCtrl)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Activity Logs
	mux.Handle("/activitylogs/edit", authMW(handleEditActivityLog(db)))
	mux.Handle("/activitylogs/delete", authMW(handleDeleteActivityLog(db)))
	mux.Handle("/activitylogs", authMW(handleGetActivityLogs(db, workOrderCtrl, mutationCtrl)))

	// Mutations
	mux.Handle("/mutations/timeline/edit", authMW(handleEditMutationTimeline(mutationCtrl)))
	mux.Handle("/mutations/timeline/delete", authMW(handleDeleteMutationTimeline(mutationCtrl)))
	mux.Handle("/mutations/timeline/all", authMW(handleGetAllAssetMutationTimelines(mutationCtrl)))
	mux.Handle("/mutations/timeline", authMW(handleGetAssetMutationTimeline(mutationCtrl)))
	mux.Handle("/assets/mutation-timeline", authMW(handleGetAssetMutationTimeline(mutationCtrl)))
	mux.Handle("/mutations/code", authMW(handleMutateAssetByCode(db, mutationCtrl)))
	mux.Handle("/mutations", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleCreateMutation(db, mutationCtrl)(w, r)
		case http.MethodGet:
			handleGetMutationHistory(mutationCtrl)(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Workorders
	mux.Handle("/workorders/assign", authMW(handleAssignWorker(db, workOrderCtrl)))
	mux.Handle("/workorders/status", authMW(handleUpdateWorkOrderStatus(db, workOrderCtrl)))
	mux.Handle("/workorders/edit", authMW(handleEditWorkOrderDetail(workOrderCtrl)))
	mux.Handle("/workorders/close", authMW(handleCloseWorkOrder(db, workOrderCtrl)))
	mux.Handle("/workorders/cancel", authMW(handleCancelWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders/delete", authMW(handleDeleteWorkOrder(workOrderCtrl)))
	mux.Handle("/workorders/logs/edit", authMW(handleEditWorkOrderLog(workOrderCtrl)))
	mux.Handle("/workorders/logs/delete", authMW(handleDeleteWorkOrderLog(workOrderCtrl)))
	mux.Handle("/workorders/logs", authMW(handleGetWorkOrderTimeline(workOrderCtrl)))
	mux.Handle("/workorders/timeline/add", authMW(handleAddTimelineEntry(workOrderCtrl)))
	mux.Handle("/workorders/timeline", authMW(handleGetWorkOrderTimeline(workOrderCtrl)))
	mux.Handle("/workorders", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	mux.Handle("/maintenance/schedules", authMW(handleGetAllPMSchedules(db, maintenanceCtrl)))
	mux.Handle("/maintenance/schedule", authMW(handleCreatePMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/edit", authMW(handleEditPMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/delete", authMW(handleDeletePMSchedule(maintenanceCtrl)))
	mux.Handle("/maintenance/history/edit", authMW(handleEditMaintenanceHistory(maintenanceCtrl)))
	mux.Handle("/maintenance/history/delete", authMW(handleDeleteMaintenanceHistory(maintenanceCtrl)))
	mux.Handle("/maintenance/", authMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// User Management endpoints
	mux.Handle("/users/engineers", authMW(handleGetEngineers(userCtrl)))
	mux.Handle("/users/create", authMW(handleCreateUser(userCtrl)))
	mux.Handle("/users/edit", authMW(handleEditUser(userCtrl)))
	mux.Handle("/users/delete", authMW(handleDeleteUser(userCtrl)))
	mux.Handle("/users", authMW(handleGetAllUsers(userCtrl)))

	return mux
}

// Handlers

func handleHome(w http.ResponseWriter, r *http.Request) {
	// Try serving static frontend from ./dist or ../frontend/dist
	dirs := []string{"./dist", "../frontend/dist", "./public"}
	for _, dir := range dirs {
		if stat, err := os.Stat(dir); err == nil && stat.IsDir() {
			filePath := filepath.Join(dir, filepath.Clean(r.URL.Path))
			if fstat, err := os.Stat(filePath); err == nil && !fstat.IsDir() {
				http.ServeFile(w, r, filePath)
				return
			}
			// Fallback to index.html for SPA client-side routes
			indexPath := filepath.Join(dir, "index.html")
			if _, err := os.Stat(indexPath); err == nil {
				http.ServeFile(w, r, indexPath)
				return
			}
		}
	}
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
			status := http.StatusUnauthorized
			if err.Error() == "Akun sedang tersambung di perangkat lain" {
				status = http.StatusConflict
			}
			utils.SendError(w, status, err.Error(), err.Error())
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
		var payload models.Asset
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		assetID := 0
		if idStr != "" {
			assetID, _ = strconv.Atoi(idStr)
		}
		if assetID == 0 {
			assetID = payload.ID
		}

		if assetID == 0 {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", "ID parameter is missing")
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

		idStr := r.URL.Query().Get("id")
		var payload struct {
			AssetID int `json:"asset_id"`
		}
		json.NewDecoder(r.Body).Decode(&payload)

		assetID := 0
		if idStr != "" {
			assetID, _ = strconv.Atoi(idStr)
		}
		if assetID == 0 {
			assetID = payload.AssetID
		}

		if assetID == 0 {
			utils.SendError(w, http.StatusBadRequest, "Invalid asset ID", "ID parameter missing")
			return
		}

		if err := assetCtrl.DeleteAsset(assetID, role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to delete asset", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Asset deleted successfully", map[string]int{"asset_id": assetID})
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

func handleCreateMutation(db *gorm.DB, mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "external"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var raw struct {
			AssetID          int       `json:"asset_id"`
			PreviousLocation string    `json:"previous_location"`
			NewLocation      string    `json:"new_location"`
			PIC              string    `json:"pic"`
			NewPIC           string    `json:"new_pic"`
			Reason           string    `json:"reason"`
			MutationDate     time.Time `json:"mutation_date"`
		}
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		picVal := strings.TrimSpace(raw.PIC)
		if picVal == "" {
			picVal = strings.TrimSpace(raw.NewPIC)
		}
		if picVal == "" {
			picVal = "Engineering"
		}

		payload := controllers.MutationInput{
			AssetID:          raw.AssetID,
			PreviousLocation: raw.PreviousLocation,
			NewLocation:      raw.NewLocation,
			PIC:              picVal,
			Reason:           raw.Reason,
			MutationDate:     raw.MutationDate,
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

func handleGetAllAssetMutationTimelines(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timelines, err := mutationCtrl.GetAllAssetMutationTimelines()
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch mutation timelines", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "All asset mutation timelines retrieved successfully", timelines)
	}
}

func handleGetAssetMutationTimeline(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		assetCode := r.URL.Query().Get("asset_code")
		idStr := r.URL.Query().Get("assetID")

		if assetCode == "" && idStr != "" {
			assetID, _ := strconv.Atoi(idStr)
			history, err := mutationCtrl.GetLocationHistory(assetID)
			if err != nil {
				utils.SendError(w, http.StatusInternalServerError, "Failed to fetch mutation history", err.Error())
				return
			}
			utils.SendSuccess(w, http.StatusOK, "Asset mutation timeline retrieved successfully", history)
			return
		}

		if assetCode == "" {
			utils.SendError(w, http.StatusBadRequest, "Asset code parameter required", "asset_code query parameter missing")
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

func handleMutateAssetByCode(db *gorm.DB, mutationCtrl *controllers.MutationController) http.HandlerFunc {
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

		if err := workOrderCtrl.CreateWorkOrder(&payload, role); err != nil {
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

func handleAssignWorker(db *gorm.DB, workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
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

func handleUpdateWorkOrderStatus(db *gorm.DB, workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
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
		role := "external"
		username := "User"
		if claims != nil {
			if claims.Role != "" {
				role = claims.Role
			}
			if claims.Username != "" {
				username = claims.Username
			}
		}

		var payload struct {
			WOID   int    `json:"wo_id"`
			Reason string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.CancelWorkOrder(payload.WOID, role, username, payload.Reason); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to cancel work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order cancelled successfully", payload)
	}
}

func handleCloseWorkOrder(db *gorm.DB, workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
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
		// Fetch ALL Work Orders for full audit coverage
		var allWOs []models.WorkOrder
		db.Order("id desc").Limit(100).Find(&allWOs)

		// Fetch work_order_logs (primary WO audit trail)
		workOrderLogs, _ := workOrderCtrl.GetAllWorkOrderLogs()

		// Fetch all Asset Mutation Timelines
		assetMutationTimelines, _ := mutationCtrl.GetAllAssetMutationTimelines()

		// Fetch maintenance_histories
		var maintenanceHistory []models.MaintenanceHistory
		db.Order("id desc").Limit(100).Find(&maintenanceHistory)

		// Fetch structured activity_logs
		var activityLogs []models.ActivityLog
		db.Order("timestamp desc").Limit(200).Find(&activityLogs)

		result := map[string]interface{}{
			"work_orders":              allWOs,
			"work_order_logs":          workOrderLogs,
			"timelines":                workOrderLogs,
			"asset_mutation_timelines": assetMutationTimelines,
			"maintenance_history":      maintenanceHistory,
			"mutations":                assetMutationTimelines,
			"activity_logs":            activityLogs,
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

func handleGetAllPMSchedules(db *gorm.DB, maintenanceCtrl *controllers.MaintenanceController) http.HandlerFunc {
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

func handleGetAllUsers(userCtrl *controllers.UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := userCtrl.GetAllUsers()
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch users", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Users retrieved successfully", users)
	}
}

func handleGetEngineers(userCtrl *controllers.UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		engineers, err := userCtrl.GetEngineers()
		if err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to fetch engineers", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Engineers retrieved successfully", engineers)
	}
}

func handleCreateUser(userCtrl *controllers.UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "admin"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload models.User
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := userCtrl.CreateUser(&payload, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, err.Error(), err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Pengguna baru berhasil ditambahkan!", payload)
	}
}

func handleEditUser(userCtrl *controllers.UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)

		var payload struct {
			ID       int    `json:"id"`
			UserID   int    `json:"user_id"`
			Username string `json:"username"`
			Password string `json:"password"`
			Name     string `json:"name"`
			Role     string `json:"role"`
			Avatar   string `json:"avatar"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Gagal membaca data JSON", err.Error())
			return
		}

		targetID := payload.UserID
		if targetID <= 0 {
			targetID = payload.ID
		}

		if targetID <= 0 {
			utils.SendError(w, http.StatusBadRequest, "ID Pengguna tidak valid", "user_id or id is required")
			return
		}

		user := models.User{
			ID:       targetID,
			Username: strings.TrimSpace(payload.Username),
			Password: strings.TrimSpace(payload.Password),
			Name:     strings.TrimSpace(payload.Name),
			Role:     strings.TrimSpace(payload.Role),
			Avatar:   payload.Avatar,
		}

		if err := userCtrl.EditUser(&user, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, err.Error(), err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Data pengguna berhasil diperbarui!", user)
	}
}

func handleDeleteUser(userCtrl *controllers.UserController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		role := "admin"
		if claims != nil && claims.Role != "" {
			role = claims.Role
		}

		var payload struct {
			UserID int `json:"user_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if claims != nil && payload.UserID == claims.UserID {
			utils.SendError(w, http.StatusBadRequest, "Anda tidak dapat menghapus akun Anda sendiri yang sedang aktif!", "Self deletion prohibited")
			return
		}

		if err := userCtrl.DeleteUser(payload.UserID, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, err.Error(), err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Pengguna berhasil dihapus permanen!", payload)
	}
}

func handleGetMe(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.UserID == 0 {
			utils.SendError(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
			return
		}
		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			utils.SendError(w, http.StatusNotFound, "User not found", err.Error())
			return
		}
		user.Password = ""
		utils.SendSuccess(w, http.StatusOK, "User active profile", user)
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

func handleLogout(authCtrl *controllers.AuthController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims != nil && claims.UserID > 0 {
			_ = authCtrl.Logout(claims.UserID)
		}
		utils.SendSuccess(w, http.StatusOK, "Logout successful", nil)
	}
}

func getUserRoleFromRequest(r *http.Request) string {
	role := "external"
	claims := middlewares.GetClaimsFromContext(r)
	if claims != nil && claims.Role != "" {
		role = claims.Role
	}
	return role
}

func handleEditActivityLog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)
		if role != "admin" && role != "hod" {
			utils.SendError(w, http.StatusForbidden, "Forbidden", "Hanya Admin atau HOD yang dapat mengubah log aktivitas")
			return
		}
		var payload struct {
			ID       int    `json:"id"`
			Action   string `json:"action"`
			Category string `json:"category"`
			Actor    string `json:"actor"`
			EntityID string `json:"entity_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}
		var log models.ActivityLog
		if err := db.First(&log, payload.ID).Error; err != nil {
			utils.SendError(w, http.StatusNotFound, "Not Found", "Log aktivitas tidak ditemukan")
			return
		}
		log.Action = payload.Action
		log.Category = payload.Category
		log.Actor = payload.Actor
		log.EntityID = payload.EntityID
		if err := db.Save(&log).Error; err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Database Error", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Log aktivitas berhasil diperbarui", log)
	}
}

func handleDeleteActivityLog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)
		if role != "admin" {
			utils.SendError(w, http.StatusForbidden, "Forbidden", "Hanya Admin yang dapat menghapus log aktivitas sistem")
			return
		}
		var payload struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}
		res := db.Where("id = ?", payload.ID).Delete(&models.ActivityLog{})
		if res.Error != nil {
			utils.SendError(w, http.StatusInternalServerError, "Database Error", res.Error.Error())
			return
		}
		if res.RowsAffected == 0 {
			utils.SendError(w, http.StatusNotFound, "Not Found", "Log aktivitas tidak ditemukan")
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Log aktivitas berhasil dihapus", payload)
	}
}

func handleEditWorkOrderLog(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)
		var payload struct {
			LogID       int    `json:"log_id"`
			ActionTaken string `json:"action_taken"`
			Cost        int    `json:"cost"`
			Status      string `json:"status"`
			UpdatedBy   string `json:"updated_by"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}
		if err := workOrderCtrl.EditWorkOrderLog(payload.LogID, payload.ActionTaken, payload.Cost, payload.Status, payload.UpdatedBy, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to edit WO Log", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Log WO berhasil diperbarui", payload)
	}
}

func handleDeleteWorkOrderLog(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)
		var payload struct {
			LogID int `json:"log_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}
		if err := workOrderCtrl.DeleteWorkOrderLog(payload.LogID, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to delete WO Log", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Log WO berhasil dihapus", payload)
	}
}

func handleEditMutationTimeline(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)
		var payload struct {
			ID               int    `json:"id"`
			PreviousLocation string `json:"previous_location"`
			NewLocation      string `json:"new_location"`
			PIC              string `json:"pic"`
			Reason           string `json:"reason"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}
		if err := mutationCtrl.EditMutation(payload.ID, payload.PreviousLocation, payload.NewLocation, payload.PIC, payload.Reason, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to edit Mutation Timeline", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Mutasi aset berhasil diperbarui", payload)
	}
}

func handleDeleteMutationTimeline(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		role := getUserRoleFromRequest(r)
		var payload struct {
			ID int `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
			return
		}
		if err := mutationCtrl.DeleteMutation(payload.ID, role); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to delete Mutation Timeline", err.Error())
			return
		}
		utils.SendSuccess(w, http.StatusOK, "Mutasi aset berhasil dihapus", payload)
	}
}

