package handlers

import (
	"net/http"

	"dr-mario-backend/middleware"
	"dr-mario-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DoctorProfileRequest struct {
	License        string `json:"license" binding:"required"`
	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	Hospital       string `json:"hospital"`
}

// GetDoctors returns all doctors
func GetDoctors(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	doctors, err := storage.GlobalStorage.GetAllDoctors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch doctors"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctors": doctors})
}

// GetDoctor returns a specific doctor by ID
func GetDoctor(c *gin.Context) {
	doctorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	doctor, err := storage.GlobalStorage.GetDoctorByID(doctorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctor": doctor})
}

// GetDoctorProfile returns the current doctor's profile
func GetDoctorProfile(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	doctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctor": doctor})
}

// UpdateDoctorProfile updates the current doctor's profile
func UpdateDoctorProfile(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if user.Role != "doctor" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req DoctorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor profile not found"})
		return
	}

	// Update fields
	if req.License != "" {
		doctor.License = req.License
	}
	if req.Specialization != "" {
		doctor.Specialization = req.Specialization
	}
	if req.Experience > 0 {
		doctor.Experience = req.Experience
	}
	if req.Hospital != "" {
		doctor.Hospital = req.Hospital
	}

	if err := storage.GlobalStorage.UpdateDoctor(doctor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"doctor": doctor})
}
