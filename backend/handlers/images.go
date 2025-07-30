package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"dr-mario-backend/config"
	"dr-mario-backend/middleware"
	"dr-mario-backend/services"
	"dr-mario-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageUploadRequest struct {
	ImageType string `form:"image_type" binding:"required,oneof=left_eye right_eye"`
	Notes     string `form:"notes"`
}

type DetectionRequest struct {
	ImageID uuid.UUID `json:"image_id" binding:"required"`
}

// UploadImage handles retinal image upload
func UploadImage(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(config.AppConfig.Upload.MaxFileSize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	var req ImageUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}
	defer file.Close()

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := false
	for _, allowedExt := range config.AppConfig.Upload.AllowedExtensions {
		if "."+allowedExt == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(config.AppConfig.Upload.UploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	filepath := filepath.Join(config.AppConfig.Upload.UploadDir, filename)

	// Create file
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Get patient ID
	var patient *storage.Patient
	if user.Role == "patient" {
		patient, err = storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Patient profile not found"})
			return
		}
	} else {
		// If doctor is uploading, get patient ID from request
		patientIDStr := c.PostForm("patient_id")
		if patientIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Patient ID required for doctor uploads"})
			return
		}
		patientID, err := uuid.Parse(patientIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
			return
		}
		patient, err = storage.GlobalStorage.GetPatientByID(patientID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
			return
		}
	}

	// Create image record
	image := &storage.RetinalImage{
		PatientID:  patient.ID,
		FileName:   header.Filename,
		FilePath:   filepath,
		FileSize:   header.Size,
		ImageType:  req.ImageType,
		UploadDate: time.Now(),
		Notes:      req.Notes,
		Status:     "pending",
	}

	if user.Role == "doctor" {
		doctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Doctor profile not found"})
			return
		}
		image.DoctorID = doctor.ID
	}

	if err := storage.GlobalStorage.CreateImage(image); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image record"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Image uploaded successfully",
		"image":   image,
	})
}

// DetectDR performs AI detection on uploaded images
func DetectDR(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req DetectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get image
	image, err := storage.GlobalStorage.GetImageByID(req.ImageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		patient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if image.PatientID != patient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Check if image exists
	if _, err := os.Stat(image.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image file not found"})
		return
	}

	// Perform AI detection
	startTime := time.Now()
	result, err := services.DetectDiabeticRetinopathy(image.FilePath)
	processingTime := time.Since(startTime).Seconds()

	if err != nil {
		// Update image status to error
		image.Status = "error"
		storage.GlobalStorage.UpdateImage(image)

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Detection failed: " + err.Error()})
		return
	}

	// Update image status
	image.Status = "processed"
	storage.GlobalStorage.UpdateImage(image)

	// Create detection result
	detectionResult := &storage.DetectionResult{
		ImageID:           image.ID,
		HasDR:             result.HasDR,
		DRStage:           result.DRStage,
		Confidence:        result.Confidence,
		HasMacularEdema:   result.HasMacularEdema,
		HasHemorrhages:    result.HasHemorrhages,
		HasExudates:       result.HasExudates,
		HasMicroaneurysms: result.HasMicroaneurysms,
		AnalysisDate:      time.Now(),
		ProcessingTime:    processingTime,
		ModelVersion:      "v1.0.0",
	}

	if user.Role == "doctor" {
		doctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err == nil {
			detectionResult.DoctorID = doctor.ID
		}
	}

	if err := storage.GlobalStorage.CreateDetectionResult(detectionResult); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save detection result"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detection completed successfully",
		"result":  detectionResult,
	})
}

// GetImages returns images for the current user
func GetImages(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var images []*storage.RetinalImage
	if user.Role == "patient" {
		patient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Patient profile not found"})
			return
		}
		images, err = storage.GlobalStorage.GetImagesByPatientID(patient.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
			return
		}
	} else if user.Role == "doctor" {
		doctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Doctor profile not found"})
			return
		}
		// For doctors, we'll return all images they have access to
		// In a real system, you might want to implement proper access control
		allPatients, err := storage.GlobalStorage.GetAllPatients()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
			return
		}
		for _, patient := range allPatients {
			patientImages, err := storage.GlobalStorage.GetImagesByPatientID(patient.ID)
			if err == nil {
				images = append(images, patientImages...)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}

// GetImage returns a specific image with its detection results
func GetImage(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	imageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	image, err := storage.GlobalStorage.GetImageByID(imageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		patient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if image.PatientID != patient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Get detection results
	results, err := storage.GlobalStorage.GetDetectionResultsByImageID(imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch detection results"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"image":   image,
		"results": results,
	})
}

// ServeImage serves the image file
func ServeImage(c *gin.Context) {
	imageID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image ID"})
		return
	}

	image, err := storage.GlobalStorage.GetImageByID(imageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
		return
	}

	// Check if file exists
	if _, err := os.Stat(image.FilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Image file not found"})
		return
	}

	c.File(image.FilePath)
}
