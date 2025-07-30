package handlers

import (
	"net/http"

	"dr-mario-backend/services"
	"dr-mario-backend/storage"

	"github.com/gin-gonic/gin"
)

// GetAnalytics returns system analytics and statistics
func GetAnalytics(c *gin.Context) {
	// Get detection statistics
	detectionStats := services.GetDetectionStats()

	// Get storage statistics
	storageStats := storage.GlobalStorage.GetStats()

	analytics := gin.H{
		"system_stats":    storageStats,
		"detection_stats": detectionStats,
	}

	c.JSON(http.StatusOK, analytics)
}

// GetPatientAnalytics returns analytics for a specific patient
func GetPatientAnalytics(c *gin.Context) {
	// This would be implemented to show patient-specific analytics
	// For now, returning a placeholder
	c.JSON(http.StatusOK, gin.H{
		"message": "Patient analytics endpoint - to be implemented",
	})
}

// GetDoctorAnalytics returns analytics for a specific doctor
func GetDoctorAnalytics(c *gin.Context) {
	// This would be implemented to show doctor-specific analytics
	// For now, returning a placeholder
	c.JSON(http.StatusOK, gin.H{
		"message": "Doctor analytics endpoint - to be implemented",
	})
}
