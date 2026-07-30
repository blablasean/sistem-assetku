package main

import (
	"log"
	"net/http"
	"os"

	"sistem-asetku-backend/config"
	"sistem-asetku-backend/controllers"
	"sistem-asetku-backend/middlewares"
	"sistem-asetku-backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if present
	_ = godotenv.Load()

	// Initialize database
	if err := config.InitDatabase(); err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	log.Println("✓ Database connected successfully")

	// Get database connection
	db := config.GetDB()

	// Instantiate controllers (Single Responsibility Architecture)
	authCtrl := controllers.NewAuthController(db)
	assetCtrl := controllers.NewAssetController(db)
	mutationCtrl := controllers.NewMutationController(db)
	workOrderCtrl := controllers.NewWorkOrderController(db)
	maintenanceCtrl := controllers.NewMaintenanceController(db)
	userCtrl := controllers.NewUserController(db)

	// Register routes
	mux := routes.RegisterRoutes(
		db,
		authCtrl,
		assetCtrl,
		mutationCtrl,
		workOrderCtrl,
		maintenanceCtrl,
		userCtrl,
	)

	// Wrap routes with CORS middleware
	handler := middlewares.CORSMiddleware(mux)

	// Get server port
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("✓ Server starting on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
