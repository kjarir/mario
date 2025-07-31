package routes

import (
	"dr-mario-backend/handlers"
	"dr-mario-backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Dr. Mario Backend is running",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Protected routes
		protected := v1.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User profile
			protected.GET("/profile", handlers.GetProfile)
			protected.PUT("/profile", handlers.UpdateProfile)

			// Patient routes
			patients := protected.Group("/patients")
			{
				patients.GET("/profile", handlers.GetPatientProfile)
				patients.PUT("/profile", handlers.UpdatePatientProfile)
				patients.GET("/", middleware.RoleMiddleware("doctor", "admin"), handlers.GetPatients)
				patients.GET("/:id", handlers.GetPatient)
				patients.GET("/:id/images", handlers.GetPatientImages)
			}

			// Doctor routes
			doctors := protected.Group("/doctors")
			{
				doctors.GET("/", handlers.GetDoctors)
				doctors.GET("/:id", handlers.GetDoctor)
				doctors.GET("/profile", middleware.RoleMiddleware("doctor"), handlers.GetDoctorProfile)
				doctors.PUT("/profile", middleware.RoleMiddleware("doctor"), handlers.UpdateDoctorProfile)
			}

			// Image routes
			images := protected.Group("/images")
			{
				images.POST("/upload", handlers.UploadImage)
				images.POST("/detect", handlers.DetectDR)
				images.POST("/scan-cnn", handlers.ScanWithCNN) // New CNN scanning endpoint
				images.GET("/", handlers.GetImages)
				images.GET("/:id", handlers.GetImage)
				images.GET("/:id/file", handlers.ServeImage)
			}

			// Appointment routes
			appointments := protected.Group("/appointments")
			{
				appointments.POST("/", handlers.CreateAppointment)
				appointments.GET("/", handlers.GetAppointments)
				appointments.GET("/:id", handlers.GetAppointment)
				appointments.PUT("/:id", handlers.UpdateAppointment)
				appointments.DELETE("/:id", handlers.CancelAppointment)
			}

			// Analytics routes (doctors and admins only)
			analytics := protected.Group("/analytics")
			analytics.Use(middleware.RoleMiddleware("doctor", "admin"))
			{
				analytics.GET("/stats", handlers.GetAnalytics)
				analytics.GET("/patient/:id", handlers.GetPatientAnalytics)
				analytics.GET("/doctor/:id", handlers.GetDoctorAnalytics)
			}

			// CNN service routes
			cnn := protected.Group("/cnn")
			{
				cnn.GET("/health", handlers.GetCNNHealth) // CNN health check
			}
		}
	}

	return router
}
