package handlers

import (
	"net/http"
	"time"

	"dr-mario-backend/middleware"
	"dr-mario-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PatientProfileRequest struct {
	DateOfBirth      string `json:"date_of_birth"`
	Gender           string `json:"gender"`
	Address          string `json:"address"`
	EmergencyContact string `json:"emergency_contact"`
	MedicalHistory   string `json:"medical_history"`
}

// GetPatientProfile returns the current patient's profile
func GetPatientProfile(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "patient" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	patient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patient": patient})
}

// UpdatePatientProfile updates the current patient's profile
func UpdatePatientProfile(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "patient" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req PatientProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	patient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient profile not found"})
		return
	}

	// Parse date of birth
	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		patient.DateOfBirth = dob
	}

	if req.Gender != "" {
		patient.Gender = req.Gender
	}
	if req.Address != "" {
		patient.Address = req.Address
	}
	if req.EmergencyContact != "" {
		patient.EmergencyContact = req.EmergencyContact
	}
	if req.MedicalHistory != "" {
		patient.MedicalHistory = req.MedicalHistory
	}

	if err := storage.GlobalStorage.UpdatePatient(patient); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patient": patient})
}

// GetPatients returns all patients (for doctors and admins)
func GetPatients(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "doctor" && user.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	patients, err := storage.GlobalStorage.GetAllPatients()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch patients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patients": patients})
}

// GetPatient returns a specific patient by ID
func GetPatient(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		currentPatient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if currentPatient.ID != patientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	patient, err := storage.GlobalStorage.GetPatientByID(patientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"patient": patient})
}

// GetPatientImages returns all images for a specific patient
func GetPatientImages(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	patientID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid patient ID"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		currentPatient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if currentPatient.ID != patientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	images, err := storage.GlobalStorage.GetImagesByPatientID(patientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch images"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"images": images})
}
