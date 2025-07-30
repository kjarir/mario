package main

import (
	"log"
	"os"

	"dr-mario-backend/config"
	"dr-mario-backend/routes"
)

func main() {
	// Load environment variables
	if err := config.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Setup router
	router := routes.SetupRouter()

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Dr. Mario Backend Server starting on port %s", port)
	log.Printf("🔬 Retinal Imaging Detection API Ready!")

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
