package main

import (
	"log"
	"os"

	"dr-mario-backend/config"
	"dr-mario-backend/routes"
	"dr-mario-backend/services"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize CNN service
	services.InitializeCNNService()
	log.Println("🔬 CNN Service initialized")

	// Setup router
	router := routes.SetupRouter()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Dr. Mario Backend Server starting on port %s", port)
	log.Printf("🔬 Retinal Imaging Detection API Ready!")
	log.Printf("🤖 CNN Integration: Active")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
