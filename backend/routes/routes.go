package routes

import (
	"dr-mario-backend/config"
	"dr-mario-backend/handlers"
	"dr-mario-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AppConfig.CORS.AllowedOrigins
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "Dr. Mario Retinal Imaging API",
			"version": "1.0.0",
		})
	})

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			profile := protected.Group("/profile")
			{
				profile.GET("", handlers.GetProfile)
				profile.PUT("", handlers.UpdateProfile)
			}

			// Patient routes
			patients := protected.Group("/patients")
			{
				patients.GET("/profile", handlers.GetPatientProfile)
				patients.PUT("/profile", handlers.UpdatePatientProfile)
				patients.GET("", middleware.RoleMiddleware("doctor", "admin"), handlers.GetPatients)
				patients.GET("/:id", handlers.GetPatient)
				patients.GET("/:id/images", handlers.GetPatientImages)
			}

			// Image routes
			images := protected.Group("/images")
			{
				images.POST("/upload", handlers.UploadImage)
				images.POST("/detect", handlers.DetectDR)
				images.GET("", handlers.GetImages)
				images.GET("/:id", handlers.GetImage)
				images.GET("/:id/file", handlers.ServeImage)
			}

			// Appointment routes
			appointments := protected.Group("/appointments")
			{
				appointments.POST("", handlers.CreateAppointment)
				appointments.GET("", handlers.GetAppointments)
				appointments.GET("/:id", handlers.GetAppointment)
				appointments.PUT("/:id", handlers.UpdateAppointment)
				appointments.DELETE("/:id", handlers.CancelAppointment)
			}

			// Doctor routes (for future expansion)
			doctors := protected.Group("/doctors")
			{
				doctors.GET("", handlers.GetDoctors)
				doctors.GET("/:id", handlers.GetDoctor)
			}

			// Analytics routes
			analytics := protected.Group("/analytics")
			{
				analytics.GET("/stats", middleware.RoleMiddleware("doctor", "admin"), handlers.GetAnalytics)
			}
		}
	}

	return router
}
