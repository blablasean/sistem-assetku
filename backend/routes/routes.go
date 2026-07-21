package routes

import (
	"encoding/json"
	"net/http"
	"sistem-asetku-backend/controllers"
	"sistem-asetku-backend/middlewares"
	"sistem-asetku-backend/models"
	"sistem-asetku-backend/utils"
	"strconv"

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
	mux.HandleFunc("GET /", handleHome)
	mux.HandleFunc("POST /auth/login", handleLogin(authCtrl))
	mux.HandleFunc("POST /auth/register", handleRegister(db))

	// Protected endpoints wrapper
	authMux := http.NewServeMux()

	// Asset endpoints
	authMux.HandleFunc("GET /assets", handleGetAssets(assetCtrl))
	authMux.HandleFunc("POST /assets", handleCreateAsset(assetCtrl))
	authMux.HandleFunc("GET /assets/{id}", handleGetAssetDetail(assetCtrl))
	authMux.HandleFunc("GET /assets/{id}/history", handleGetAssetHistory(mutationCtrl))

	// Mutation endpoints
	authMux.HandleFunc("POST /mutations", handleCreateMutation(mutationCtrl))
	authMux.HandleFunc("GET /mutations/{assetID}", handleGetMutationHistory(mutationCtrl))

	// Work Order endpoints
	authMux.HandleFunc("POST /workorders", handleCreateWorkOrder(workOrderCtrl))
	authMux.HandleFunc("POST /workorders/{id}/assign", handleAssignWorker(workOrderCtrl))
	authMux.HandleFunc("PUT /workorders/{id}/status", handleUpdateWorkOrderStatus(workOrderCtrl))
	authMux.HandleFunc("GET /workorders/{id}/status", handleGetWorkOrderStatus(workOrderCtrl))

	// Maintenance endpoints
	authMux.HandleFunc("POST /maintenance/schedule", handleCreatePMSchedule(maintenanceCtrl))
	authMux.HandleFunc("POST /maintenance/{id}/checklist", handleSubmitPMChecklist(maintenanceCtrl))
	authMux.HandleFunc("GET /maintenance/{assetID}/history", handleGetMaintenanceHistory(maintenanceCtrl))

	// Apply auth middleware to protected routes
	mux.Handle("/assets", middlewares.AuthMiddleware(authMux))
	mux.Handle("/mutations", middlewares.AuthMiddleware(authMux))
	mux.Handle("/workorders", middlewares.AuthMiddleware(authMux))
	mux.Handle("/maintenance", middlewares.AuthMiddleware(authMux))

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
		if q == "" {
			utils.SendError(w, http.StatusBadRequest, "Missing search query", "query parameter 'q' required")
			return
		}

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
		if claims == nil || claims.Role != "hod" {
			utils.SendError(w, http.StatusForbidden, "Access denied", "Only HOD can create assets")
			return
		}

		var payload models.Asset
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := assetCtrl.RegistrasiAsset(payload, claims.Role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to create asset", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Asset created successfully", payload)
	}
}

func handleGetAssetDetail(assetCtrl *controllers.AssetController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
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
		idStr := r.PathValue("id")
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
		if claims == nil {
			utils.SendError(w, http.StatusUnauthorized, "Unauthorized", "Authentication required")
			return
		}

		var payload models.Mutation
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := mutationCtrl.CreateMutation(payload, claims.Role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to create mutation", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Mutation created successfully", payload)
	}
}

func handleGetMutationHistory(mutationCtrl *controllers.MutationController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("assetID")
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
		if claims == nil || claims.Role != "hod" {
			utils.SendError(w, http.StatusForbidden, "Access denied", "Only HOD can create work orders")
			return
		}

		var payload models.WorkOrder
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.CreateWorkOrder(payload, claims.Role); err != nil {
			utils.SendError(w, http.StatusForbidden, "Failed to create work order", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusCreated, "Work order created successfully", payload)
	}
}

func handleAssignWorker(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.Role != "hod" {
			utils.SendError(w, http.StatusForbidden, "Access denied", "Only HOD can assign workers")
			return
		}

		idStr := r.PathValue("id")
		woID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid work order ID", err.Error())
			return
		}

		var payload struct {
			EngineerID int `json:"engineer_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.AssignWorker(woID, payload.EngineerID, claims.Role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to assign worker", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Worker assigned successfully", map[string]int{
			"work_order_id": woID,
			"engineer_id":   payload.EngineerID,
		})
	}
}

func handleUpdateWorkOrderStatus(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middlewares.GetClaimsFromContext(r)
		if claims == nil || claims.Role != "hod" {
			utils.SendError(w, http.StatusForbidden, "Access denied", "Only HOD can update work order status")
			return
		}

		idStr := r.PathValue("id")
		woID, err := strconv.Atoi(idStr)
		if err != nil {
			utils.SendError(w, http.StatusBadRequest, "Invalid work order ID", err.Error())
			return
		}

		var payload struct {
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.SendError(w, http.StatusBadRequest, "Failed to decode request", err.Error())
			return
		}

		if err := workOrderCtrl.UpdateWOStatus(woID, payload.Status, claims.Role); err != nil {
			utils.SendError(w, http.StatusInternalServerError, "Failed to update status", err.Error())
			return
		}

		utils.SendSuccess(w, http.StatusOK, "Work order status updated successfully", map[string]interface{}{
			"work_order_id": woID,
			"status":        payload.Status,
		})
	}
}

func handleGetWorkOrderStatus(workOrderCtrl *controllers.WorkOrderController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
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
		idStr := r.PathValue("assetID")
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
