package storage

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// In-memory storage with thread-safe operations
type Storage struct {
	users            map[uuid.UUID]*User
	patients         map[uuid.UUID]*Patient
	doctors          map[uuid.UUID]*Doctor
	images           map[uuid.UUID]*RetinalImage
	detectionResults map[uuid.UUID]*DetectionResult
	appointments     map[uuid.UUID]*Appointment
	userByEmail      map[string]*User
	mu               sync.RWMutex
}

// User represents the base user model
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Patient represents a patient in the system
type Patient struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	User             *User     `json:"user"`
	DateOfBirth      time.Time `json:"date_of_birth"`
	Gender           string    `json:"gender"`
	Address          string    `json:"address"`
	EmergencyContact string    `json:"emergency_contact"`
	MedicalHistory   string    `json:"medical_history"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Doctor represents a doctor in the system
type Doctor struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	User           *User     `json:"user"`
	License        string    `json:"license"`
	Specialization string    `json:"specialization"`
	Experience     int       `json:"experience"`
	Hospital       string    `json:"hospital"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// RetinalImage represents uploaded retinal images
type RetinalImage struct {
	ID         uuid.UUID `json:"id"`
	PatientID  uuid.UUID `json:"patient_id"`
	Patient    *Patient  `json:"patient"`
	DoctorID   uuid.UUID `json:"doctor_id"`
	Doctor     *Doctor   `json:"doctor"`
	FileName   string    `json:"file_name"`
	FilePath   string    `json:"file_path"`
	FileSize   int64     `json:"file_size"`
	ImageType  string    `json:"image_type"`
	UploadDate time.Time `json:"upload_date"`
	Notes      string    `json:"notes"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// DetectionResult represents AI detection results
type DetectionResult struct {
	ID                uuid.UUID     `json:"id"`
	ImageID           uuid.UUID     `json:"image_id"`
	Image             *RetinalImage `json:"image"`
	DoctorID          uuid.UUID     `json:"doctor_id"`
	Doctor            *Doctor       `json:"doctor"`
	HasDR             bool          `json:"has_dr"`
	DRStage           string        `json:"dr_stage"`
	Confidence        float64       `json:"confidence"`
	HasMacularEdema   bool          `json:"has_macular_edema"`
	HasHemorrhages    bool          `json:"has_hemorrhages"`
	HasExudates       bool          `json:"has_exudates"`
	HasMicroaneurysms bool          `json:"has_microaneurysms"`
	AnalysisDate      time.Time     `json:"analysis_date"`
	ProcessingTime    float64       `json:"processing_time"`
	ModelVersion      string        `json:"model_version"`
	ReviewedBy        uuid.UUID     `json:"reviewed_by"`
	ReviewDate        time.Time     `json:"review_date"`
	ReviewNotes       string        `json:"review_notes"`
	IsConfirmed       bool          `json:"is_confirmed"`
	CreatedAt         time.Time     `json:"created_at"`
	UpdatedAt         time.Time     `json:"updated_at"`
}

// Appointment represents patient appointments
type Appointment struct {
	ID              uuid.UUID `json:"id"`
	PatientID       uuid.UUID `json:"patient_id"`
	Patient         *Patient  `json:"patient"`
	DoctorID        uuid.UUID `json:"doctor_id"`
	Doctor          *Doctor   `json:"doctor"`
	AppointmentDate time.Time `json:"appointment_date"`
	Duration        int       `json:"duration"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

var GlobalStorage *Storage

func init() {
	GlobalStorage = &Storage{
		users:            make(map[uuid.UUID]*User),
		patients:         make(map[uuid.UUID]*Patient),
		doctors:          make(map[uuid.UUID]*Doctor),
		images:           make(map[uuid.UUID]*RetinalImage),
		detectionResults: make(map[uuid.UUID]*DetectionResult),
		appointments:     make(map[uuid.UUID]*Appointment),
		userByEmail:      make(map[string]*User),
	}
}

// User operations
func (s *Storage) CreateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	s.users[user.ID] = user
	s.userByEmail[user.Email] = user
	return nil
}

func (s *Storage) GetUserByID(id uuid.UUID) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, ErrNotFound
	}
	return user, nil
}

func (s *Storage) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.userByEmail[email]
	if !exists {
		return nil, ErrNotFound
	}
	return user, nil
}

func (s *Storage) UpdateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	s.userByEmail[user.Email] = user
	return nil
}

// Patient operations
func (s *Storage) CreatePatient(patient *Patient) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	patient.ID = uuid.New()
	patient.CreatedAt = time.Now()
	patient.UpdatedAt = time.Now()

	s.patients[patient.ID] = patient
	return nil
}

func (s *Storage) GetPatientByUserID(userID uuid.UUID) (*Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, patient := range s.patients {
		if patient.UserID == userID {
			// Load user data
			if user, exists := s.users[patient.UserID]; exists {
				patient.User = user
			}
			return patient, nil
		}
	}
	return nil, ErrNotFound
}

func (s *Storage) GetPatientByID(id uuid.UUID) (*Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	patient, exists := s.patients[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Load user data
	if user, exists := s.users[patient.UserID]; exists {
		patient.User = user
	}

	return patient, nil
}

func (s *Storage) GetAllPatients() ([]*Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var patients []*Patient
	for _, patient := range s.patients {
		// Load user data
		if user, exists := s.users[patient.UserID]; exists {
			patient.User = user
		}
		patients = append(patients, patient)
	}
	return patients, nil
}

func (s *Storage) UpdatePatient(patient *Patient) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	patient.UpdatedAt = time.Now()
	s.patients[patient.ID] = patient
	return nil
}

// Doctor operations
func (s *Storage) CreateDoctor(doctor *Doctor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doctor.ID = uuid.New()
	doctor.CreatedAt = time.Now()
	doctor.UpdatedAt = time.Now()

	s.doctors[doctor.ID] = doctor
	return nil
}

func (s *Storage) GetDoctorByUserID(userID uuid.UUID) (*Doctor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, doctor := range s.doctors {
		if doctor.UserID == userID {
			// Load user data
			if user, exists := s.users[doctor.UserID]; exists {
				doctor.User = user
			}
			return doctor, nil
		}
	}
	return nil, ErrNotFound
}

func (s *Storage) GetDoctorByID(id uuid.UUID) (*Doctor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	doctor, exists := s.doctors[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Load user data
	if user, exists := s.users[doctor.UserID]; exists {
		doctor.User = user
	}

	return doctor, nil
}

func (s *Storage) GetAllDoctors() ([]*Doctor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var doctors []*Doctor
	for _, doctor := range s.doctors {
		// Load user data
		if user, exists := s.users[doctor.UserID]; exists {
			doctor.User = user
		}
		doctors = append(doctors, doctor)
	}
	return doctors, nil
}

func (s *Storage) UpdateDoctor(doctor *Doctor) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	doctor.UpdatedAt = time.Now()
	s.doctors[doctor.ID] = doctor
	return nil
}

// Image operations
func (s *Storage) CreateImage(image *RetinalImage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	image.ID = uuid.New()
	image.CreatedAt = time.Now()
	image.UpdatedAt = time.Now()

	s.images[image.ID] = image
	return nil
}

func (s *Storage) GetImageByID(id uuid.UUID) (*RetinalImage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	image, exists := s.images[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Load related data
	if patient, exists := s.patients[image.PatientID]; exists {
		image.Patient = patient
		if user, exists := s.users[patient.UserID]; exists {
			image.Patient.User = user
		}
	}

	if image.DoctorID != uuid.Nil {
		if doctor, exists := s.doctors[image.DoctorID]; exists {
			image.Doctor = doctor
			if user, exists := s.users[doctor.UserID]; exists {
				image.Doctor.User = user
			}
		}
	}

	return image, nil
}

func (s *Storage) GetImagesByPatientID(patientID uuid.UUID) ([]*RetinalImage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var images []*RetinalImage
	for _, image := range s.images {
		if image.PatientID == patientID {
			// Load related data
			if patient, exists := s.patients[image.PatientID]; exists {
				image.Patient = patient
				if user, exists := s.users[patient.UserID]; exists {
					image.Patient.User = user
				}
			}

			if image.DoctorID != uuid.Nil {
				if doctor, exists := s.doctors[image.DoctorID]; exists {
					image.Doctor = doctor
					if user, exists := s.users[doctor.UserID]; exists {
						image.Doctor.User = user
					}
				}
			}

			images = append(images, image)
		}
	}
	return images, nil
}

func (s *Storage) UpdateImage(image *RetinalImage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	image.UpdatedAt = time.Now()
	s.images[image.ID] = image
	return nil
}

// Detection Result operations
func (s *Storage) CreateDetectionResult(result *DetectionResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	result.ID = uuid.New()
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()

	s.detectionResults[result.ID] = result
	return nil
}

func (s *Storage) GetDetectionResultsByImageID(imageID uuid.UUID) ([]*DetectionResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*DetectionResult
	for _, result := range s.detectionResults {
		if result.ImageID == imageID {
			// Load related data
			if image, exists := s.images[result.ImageID]; exists {
				result.Image = image
			}

			if result.DoctorID != uuid.Nil {
				if doctor, exists := s.doctors[result.DoctorID]; exists {
					result.Doctor = doctor
				}
			}

			results = append(results, result)
		}
	}
	return results, nil
}

// Appointment operations
func (s *Storage) CreateAppointment(appointment *Appointment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	appointment.ID = uuid.New()
	appointment.CreatedAt = time.Now()
	appointment.UpdatedAt = time.Now()

	s.appointments[appointment.ID] = appointment
	return nil
}

func (s *Storage) GetAppointmentByID(id uuid.UUID) (*Appointment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	appointment, exists := s.appointments[id]
	if !exists {
		return nil, ErrNotFound
	}

	// Load related data
	if patient, exists := s.patients[appointment.PatientID]; exists {
		appointment.Patient = patient
		if user, exists := s.users[patient.UserID]; exists {
			appointment.Patient.User = user
		}
	}

	if doctor, exists := s.doctors[appointment.DoctorID]; exists {
		appointment.Doctor = doctor
		if user, exists := s.users[doctor.UserID]; exists {
			appointment.Doctor.User = user
		}
	}

	return appointment, nil
}

func (s *Storage) GetAppointmentsByPatientID(patientID uuid.UUID) ([]*Appointment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var appointments []*Appointment
	for _, appointment := range s.appointments {
		if appointment.PatientID == patientID {
			// Load related data
			if patient, exists := s.patients[appointment.PatientID]; exists {
				appointment.Patient = patient
				if user, exists := s.users[patient.UserID]; exists {
					appointment.Patient.User = user
				}
			}

			if doctor, exists := s.doctors[appointment.DoctorID]; exists {
				appointment.Doctor = doctor
				if user, exists := s.users[doctor.UserID]; exists {
					appointment.Doctor.User = user
				}
			}

			appointments = append(appointments, appointment)
		}
	}
	return appointments, nil
}

func (s *Storage) GetAppointmentsByDoctorID(doctorID uuid.UUID) ([]*Appointment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var appointments []*Appointment
	for _, appointment := range s.appointments {
		if appointment.DoctorID == doctorID {
			// Load related data
			if patient, exists := s.patients[appointment.PatientID]; exists {
				appointment.Patient = patient
				if user, exists := s.users[patient.UserID]; exists {
					appointment.Patient.User = user
				}
			}

			if doctor, exists := s.doctors[appointment.DoctorID]; exists {
				appointment.Doctor = doctor
				if user, exists := s.users[doctor.UserID]; exists {
					appointment.Doctor.User = user
				}
			}

			appointments = append(appointments, appointment)
		}
	}
	return appointments, nil
}

func (s *Storage) UpdateAppointment(appointment *Appointment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	appointment.UpdatedAt = time.Now()
	s.appointments[appointment.ID] = appointment
	return nil
}

// Statistics
func (s *Storage) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return map[string]interface{}{
		"total_patients":     len(s.patients),
		"total_doctors":      len(s.doctors),
		"total_images":       len(s.images),
		"total_appointments": len(s.appointments),
		"total_detections":   len(s.detectionResults),
	}
}
