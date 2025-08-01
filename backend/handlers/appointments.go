package handlers

import (
	"net/http"
	"time"

	"dr-mario-backend/middleware"
	"dr-mario-backend/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AppointmentRequest struct {
	PatientID       uuid.UUID `json:"patient_id" binding:"required"`
	DoctorID        uuid.UUID `json:"doctor_id" binding:"required"`
	AppointmentDate string    `json:"appointment_date" binding:"required"`
	Duration        int       `json:"duration" binding:"required,min=15,max=120"`
	Notes           string    `json:"notes"`
}

type UpdateAppointmentRequest struct {
	AppointmentDate string `json:"appointment_date"`
	Duration        int    `json:"duration"`
	Status          string `json:"status"`
	Notes           string `json:"notes"`
}

// CreateAppointment creates a new appointment
func CreateAppointment(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var req AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse appointment date
	appointmentDate, err := time.Parse("2006-01-02T15:04:05Z", req.AppointmentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Check if appointment is in the future
	if appointmentDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment date must be in the future"})
		return
	}

	// Verify patient exists
	_, err = storage.GlobalStorage.GetPatientByID(req.PatientID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Patient not found"})
		return
	}

	// Verify doctor exists
	_, err = storage.GlobalStorage.GetDoctorByID(req.DoctorID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Doctor not found"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		currentPatient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if currentPatient.ID != req.PatientID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Can only create appointments for yourself"})
			return
		}
	} else if user.Role == "doctor" {
		currentDoctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if currentDoctor.ID != req.DoctorID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Can only create appointments for yourself"})
			return
		}
	}

	// Check for scheduling conflicts
	doctorAppointments, err := storage.GlobalStorage.GetAppointmentsByDoctorID(req.DoctorID)
	if err == nil {
		for _, existingAppointment := range doctorAppointments {
			if existingAppointment.AppointmentDate.Equal(appointmentDate) && existingAppointment.Status != "cancelled" {
				c.JSON(http.StatusConflict, gin.H{"error": "Doctor has another appointment at this time"})
				return
			}
		}
	}

	// Create appointment
	appointment := &storage.Appointment{
		PatientID:       req.PatientID,
		DoctorID:        req.DoctorID,
		AppointmentDate: appointmentDate,
		Duration:        req.Duration,
		Status:          "scheduled",
		Notes:           req.Notes,
	}

	if err := storage.GlobalStorage.CreateAppointment(appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Appointment created successfully",
		"appointment": appointment,
	})
}

// GetAppointments returns appointments for the current user
func GetAppointments(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var appointments []*storage.Appointment
	if user.Role == "patient" {
		patient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Patient profile not found"})
			return
		}
		appointments, err = storage.GlobalStorage.GetAppointmentsByPatientID(patient.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
			return
		}
	} else if user.Role == "doctor" {
		doctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Doctor profile not found"})
			return
		}
		appointments, err = storage.GlobalStorage.GetAppointmentsByDoctorID(doctor.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"appointments": appointments})
}

// GetAppointment returns a specific appointment
func GetAppointment(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	appointmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := storage.GlobalStorage.GetAppointmentByID(appointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		currentPatient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if appointment.PatientID != currentPatient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	} else if user.Role == "doctor" {
		currentDoctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if appointment.DoctorID != currentDoctor.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"appointment": appointment})
}

// UpdateAppointment updates an appointment
func UpdateAppointment(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	appointmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	var req UpdateAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	appointment, err := storage.GlobalStorage.GetAppointmentByID(appointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		currentPatient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if appointment.PatientID != currentPatient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	} else if user.Role == "doctor" {
		currentDoctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if appointment.DoctorID != currentDoctor.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Update fields
	if req.AppointmentDate != "" {
		appointmentDate, err := time.Parse("2006-01-02T15:04:05Z", req.AppointmentDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		if appointmentDate.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment date must be in the future"})
			return
		}
		appointment.AppointmentDate = appointmentDate
	}

	if req.Duration > 0 {
		appointment.Duration = req.Duration
	}

	if req.Status != "" {
		appointment.Status = req.Status
	}

	if req.Notes != "" {
		appointment.Notes = req.Notes
	}

	if err := storage.GlobalStorage.UpdateAppointment(appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Appointment updated successfully",
		"appointment": appointment,
	})
}

// CancelAppointment cancels an appointment
func CancelAppointment(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	appointmentID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid appointment ID"})
		return
	}

	appointment, err := storage.GlobalStorage.GetAppointmentByID(appointmentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
		return
	}

	// Check permissions
	if user.Role == "patient" {
		currentPatient, err := storage.GlobalStorage.GetPatientByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if appointment.PatientID != currentPatient.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	} else if user.Role == "doctor" {
		currentDoctor, err := storage.GlobalStorage.GetDoctorByUserID(user.ID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		if appointment.DoctorID != currentDoctor.ID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Check if appointment can be cancelled
	if appointment.Status == "cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Appointment is already cancelled"})
		return
	}

	if appointment.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot cancel completed appointment"})
		return
	}

	appointment.Status = "cancelled"
	if err := storage.GlobalStorage.UpdateAppointment(appointment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel appointment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Appointment cancelled successfully",
	})
}
