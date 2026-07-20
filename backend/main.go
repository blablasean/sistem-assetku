package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"sistem-asetku-backend/config"
	controllers "sistem-asetku-backend/controllers"
	"sistem-asetku-backend/models"

	"github.com/joho/godotenv"
)

func main() {
	// load .env if present
	_ = godotenv.Load()

	if err := config.InitDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	log.Println("database connected")

	// instantiate controllers
	authCtrl := controllers.NewAuthController(config.DB)
	mutationCtrl := controllers.NewMutationController(config.DB)
	assetCtrl := controllers.NewAssetController(config.DB)
	workOrderCtrl := controllers.NewWorkOrderController(config.DB)
	maintenanceCtrl := controllers.NewMaintenanceController(config.DB)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Sistem AsetKu backend running")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		token, err := authCtrl.Login(payload.Username, payload.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		resp := map[string]string{"token": token}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	// Asset endpoints
	http.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var payload models.Asset
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				http.Error(w, "invalid body", http.StatusBadRequest)
				return
			}
			callerRole := r.Header.Get("X-User-Role")
			if callerRole == "" {
				callerRole = "guest"
			}
			if err := assetCtrl.RegistrasiAsset(payload, callerRole); err != nil {
				http.Error(w, err.Error(), http.StatusForbidden)
				return
			}
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintln(w, "asset created")
			return
		case http.MethodGet:
			// support ?id= or ?q=
			q := r.URL.Query().Get("q")
			id := r.URL.Query().Get("id")
			if id != "" {
				// fetching by numeric id not implemented in this demo endpoint
				http.Error(w, "fetch by id not implemented in this endpoint; use query param 'q' to search", http.StatusNotImplemented)
				return
			}
			if q != "" {
				assets, err := assetCtrl.SearchAndFilterAssets(q)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(assets)
				return
			}
			http.Error(w, "missing parameters", http.StatusBadRequest)
			return
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})

	// WorkOrder endpoints
	http.HandleFunc("/workorders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload models.WorkOrder
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		callerRole := r.Header.Get("X-User-Role")
		if callerRole == "" {
			callerRole = "guest"
		}
		if err := workOrderCtrl.CreateWorkOrder(payload, callerRole); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "workorder created")
	})

	http.HandleFunc("/workorders/assign", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			WOId       int `json:"wo_id"`
			EngineerId int `json:"engineer_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		callerRole := r.Header.Get("X-User-Role")
		if callerRole == "" {
			callerRole = "guest"
		}
		if err := workOrderCtrl.AssignWorker(payload.WOId, payload.EngineerId, callerRole); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		fmt.Fprintln(w, "worker assigned")
	})

	http.HandleFunc("/workorders/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload struct {
			WOId   int    `json:"wo_id"`
			Status string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		callerRole := r.Header.Get("X-User-Role")
		if callerRole == "" {
			callerRole = "guest"
		}
		if err := workOrderCtrl.UpdateWOStatus(payload.WOId, payload.Status, callerRole); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		fmt.Fprintln(w, "status updated")
	})

	// Maintenance endpoints
	http.HandleFunc("/maintenance/pm", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload models.PreventiveMaintenance
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		callerRole := r.Header.Get("X-User-Role")
		if callerRole == "" {
			callerRole = "guest"
		}
		if err := maintenanceCtrl.CreatePMSchedule(payload, callerRole); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "pm schedule created")
	})

	http.HandleFunc("/maintenance/history", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		assetID := r.URL.Query().Get("asset_id")
		if assetID == "" {
			http.Error(w, "asset_id required", http.StatusBadRequest)
			return
		}
		// for brevity, not parsing to int here - controller expects int
		http.Error(w, "asset_id parsing not implemented in this demo endpoint", http.StatusNotImplemented)
	})

	http.HandleFunc("/mutation", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var payload models.Mutation
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		// caller role header (for demo purposes)
		callerRole := r.Header.Get("X-User-Role")
		if callerRole == "" {
			callerRole = "guest"
		}

		if err := mutationCtrl.CreateMutation(payload, callerRole); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "mutation created")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
