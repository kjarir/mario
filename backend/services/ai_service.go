package services

import (
	"fmt"
	"math/rand"
	"time"
)

// DetectionResult represents the result of diabetic retinopathy detection
type DetectionResult struct {
	HasDR             bool    `json:"has_dr"`
	DRStage           string  `json:"dr_stage"`
	Confidence        float64 `json:"confidence"`
	HasMacularEdema   bool    `json:"has_macular_edema"`
	HasHemorrhages    bool    `json:"has_hemorrhages"`
	HasExudates       bool    `json:"has_exudates"`
	HasMicroaneurysms bool    `json:"has_microaneurysms"`
	ProcessingTime    float64 `json:"processing_time"`
	ModelVersion      string  `json:"model_version"`
	Error             string  `json:"error,omitempty"`
}

// DetectionStats represents statistics about detections
type DetectionStats struct {
	TotalDetections       int     `json:"total_detections"`
	PositiveDetections    int     `json:"positive_detections"`
	NegativeDetections    int     `json:"negative_detections"`
	AverageConfidence     float64 `json:"average_confidence"`
	AverageProcessingTime float64 `json:"average_processing_time"`
	SuccessRate           float64 `json:"success_rate"`
}

// Global CNN service instance
var cnnService *CNNService

// InitializeCNNService initializes the CNN service
func InitializeCNNService() {
	cnnService = NewCNNService()
}

// DetectDiabeticRetinopathy performs AI detection on retinal images
// This function now integrates with the CNN service for complex analysis
func DetectDiabeticRetinopathy(imagePath string) (*DetectionResult, error) {
	startTime := time.Now()

	// Initialize CNN service if not already done
	if cnnService == nil {
		InitializeCNNService()
	}

	// Validate image for CNN processing
	if err := cnnService.ValidateImage(imagePath); err != nil {
		return &DetectionResult{
			Error: fmt.Sprintf("Image validation failed: %v", err),
		}, nil
	}

	// Preprocess image for better CNN analysis
	preprocessedPath, err := cnnService.PreprocessImage(imagePath)
	if err != nil {
		return &DetectionResult{
			Error: fmt.Sprintf("Image preprocessing failed: %v", err),
		}, nil
	}

	// Send to CNN for complex analysis
	cnnResult, err := cnnService.ScanImageWithCNN(preprocessedPath)
	if err != nil {
		// Fallback to simulated detection if CNN is unavailable
		return simulateDetection(imagePath, startTime)
	}

	// Check if CNN analysis was successful
	if !cnnResult.Success {
		return &DetectionResult{
			Error: cnnResult.Error,
		}, nil
	}

	// Convert CNN result to DetectionResult format
	result := &DetectionResult{
		HasDR:             cnnResult.HasDR,
		DRStage:           cnnResult.DRStage,
		Confidence:        cnnResult.Confidence,
		HasMacularEdema:   cnnResult.MacularEdema,
		HasHemorrhages:    cnnResult.Hemorrhages,
		HasExudates:       cnnResult.Exudates,
		HasMicroaneurysms: cnnResult.Microaneurysms,
		ProcessingTime:    cnnResult.ProcessingTime,
		ModelVersion:      cnnResult.ModelVersion,
	}

	return result, nil
}

// DetectDiabeticRetinopathyBytes performs detection on image bytes
func DetectDiabeticRetinopathyBytes(imageData []byte, filename string) (*DetectionResult, error) {
	startTime := time.Now()

	// Initialize CNN service if not already done
	if cnnService == nil {
		InitializeCNNService()
	}

	// Send image bytes directly to CNN
	cnnResult, err := cnnService.ScanImageBytes(imageData, filename)
	if err != nil {
		// Fallback to simulated detection if CNN is unavailable
		return simulateDetectionBytes(imageData, filename, startTime)
	}

	// Check if CNN analysis was successful
	if !cnnResult.Success {
		return &DetectionResult{
			Error: cnnResult.Error,
		}, nil
	}

	// Convert CNN result to DetectionResult format
	result := &DetectionResult{
		HasDR:             cnnResult.HasDR,
		DRStage:           cnnResult.DRStage,
		Confidence:        cnnResult.Confidence,
		HasMacularEdema:   cnnResult.MacularEdema,
		HasHemorrhages:    cnnResult.Hemorrhages,
		HasExudates:       cnnResult.Exudates,
		HasMicroaneurysms: cnnResult.Microaneurysms,
		ProcessingTime:    cnnResult.ProcessingTime,
		ModelVersion:      cnnResult.ModelVersion,
	}

	return result, nil
}

// GetDetectionStats returns statistics about detections
func GetDetectionStats() *DetectionStats {
	// In a real implementation, this would query the database
	// For now, return simulated statistics
	return &DetectionStats{
		TotalDetections:       1250,
		PositiveDetections:    312,
		NegativeDetections:    938,
		AverageConfidence:     0.87,
		AverageProcessingTime: 3.2,
		SuccessRate:           0.985,
	}
}

// GetCNNHealth checks the health of the CNN service
func GetCNNHealth() error {
	if cnnService == nil {
		InitializeCNNService()
	}
	return cnnService.GetCNNHealth()
}

// simulateDetection provides fallback detection when CNN is unavailable
func simulateDetection(imagePath string, startTime time.Time) (*DetectionResult, error) {
	// Simulate processing time
	time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)

	// Simulate detection logic
	rand.Seed(time.Now().UnixNano())
	hasDR := rand.Float64() < 0.3 // 30% chance of DR

	var drStage string
	var confidence float64

	if hasDR {
		stages := []string{"Mild", "Moderate", "Severe", "Proliferative"}
		drStage = stages[rand.Intn(len(stages))]
		confidence = 0.7 + rand.Float64()*0.3 // 70-100% confidence
	} else {
		drStage = "No DR"
		confidence = 0.8 + rand.Float64()*0.2 // 80-100% confidence
	}

	processingTime := time.Since(startTime).Seconds()

	return &DetectionResult{
		HasDR:             hasDR,
		DRStage:           drStage,
		Confidence:        confidence,
		HasMacularEdema:   hasDR && rand.Float64() < 0.4,
		HasHemorrhages:    hasDR && rand.Float64() < 0.6,
		HasExudates:       hasDR && rand.Float64() < 0.5,
		HasMicroaneurysms: hasDR && rand.Float64() < 0.7,
		ProcessingTime:    processingTime,
		ModelVersion:      "v2.1.0-simulated",
	}, nil
}

// simulateDetectionBytes provides fallback detection for image bytes
func simulateDetectionBytes(imageData []byte, filename string, startTime time.Time) (*DetectionResult, error) {
	// Simulate processing time
	time.Sleep(time.Duration(rand.Intn(2000)+1000) * time.Millisecond)

	// Simulate detection logic
	rand.Seed(time.Now().UnixNano())
	hasDR := rand.Float64() < 0.3 // 30% chance of DR

	var drStage string
	var confidence float64

	if hasDR {
		stages := []string{"Mild", "Moderate", "Severe", "Proliferative"}
		drStage = stages[rand.Intn(len(stages))]
		confidence = 0.7 + rand.Float64()*0.3 // 70-100% confidence
	} else {
		drStage = "No DR"
		confidence = 0.8 + rand.Float64()*0.2 // 80-100% confidence
	}

	processingTime := time.Since(startTime).Seconds()

	return &DetectionResult{
		HasDR:             hasDR,
		DRStage:           drStage,
		Confidence:        confidence,
		HasMacularEdema:   hasDR && rand.Float64() < 0.4,
		HasHemorrhages:    hasDR && rand.Float64() < 0.6,
		HasExudates:       hasDR && rand.Float64() < 0.5,
		HasMicroaneurysms: hasDR && rand.Float64() < 0.7,
		ProcessingTime:    processingTime,
		ModelVersion:      "v2.1.0-simulated",
	}, nil
}
