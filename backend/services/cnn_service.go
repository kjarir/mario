package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"dr-mario-backend/config"
)

// CNNScanResult represents the result from CNN analysis
type CNNScanResult struct {
	Success        bool    `json:"success"`
	HasDR          bool    `json:"has_dr"`
	DRStage        string  `json:"dr_stage"`
	Confidence     float64 `json:"confidence"`
	Severity       string  `json:"severity"`
	RiskLevel      string  `json:"risk_level"`
	Recommendation string  `json:"recommendation"`

	// Detailed Analysis
	MacularEdema       bool `json:"macular_edema"`
	Hemorrhages        bool `json:"hemorrhages"`
	Exudates           bool `json:"exudates"`
	Microaneurysms     bool `json:"microaneurysms"`
	Neovascularization bool `json:"neovascularization"`

	// Quantitative Measurements
	LesionCount      int     `json:"lesion_count"`
	LesionArea       float64 `json:"lesion_area_percentage"`
	VesselTortuosity float64 `json:"vessel_tortuosity"`

	// Processing Information
	ProcessingTime float64 `json:"processing_time"`
	ModelVersion   string  `json:"model_version"`
	AnalysisDate   string  `json:"analysis_date"`

	// Error Information
	Error string `json:"error,omitempty"`
}

// CNNService handles communication with the CNN model
type CNNService struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewCNNService creates a new CNN service instance
func NewCNNService() *CNNService {
	return &CNNService{
		baseURL: getCNNBaseURL(),
		apiKey:  getCNNAPIKey(),
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // 60 seconds timeout for CNN processing
		},
	}
}

// ScanImageWithCNN sends an image to the CNN for complex analysis
func (c *CNNService) ScanImageWithCNN(imagePath string) (*CNNScanResult, error) {
	startTime := time.Now()

	// Validate image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file not found: %s", imagePath)
	}

	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image file
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("image", filepath.Base(imagePath))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy image data: %v", err)
	}

	// Add additional parameters
	writer.WriteField("api_key", c.apiKey)
	writer.WriteField("model_version", "v2.1.0")
	writer.WriteField("analysis_type", "comprehensive")
	writer.WriteField("confidence_threshold", "0.7")

	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/scan", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", "DrMario-Backend/1.0")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to CNN: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return &CNNScanResult{
			Success: false,
			Error:   fmt.Sprintf("CNN API returned status %d: %s", resp.StatusCode, string(respBody)),
		}, nil
	}

	// Parse JSON response
	var result CNNScanResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CNN response: %v", err)
	}

	// Calculate processing time
	result.ProcessingTime = time.Since(startTime).Seconds()
	result.AnalysisDate = time.Now().Format(time.RFC3339)

	return &result, nil
}

// ScanImageBytes scans image data directly from bytes
func (c *CNNService) ScanImageBytes(imageData []byte, filename string) (*CNNScanResult, error) {
	startTime := time.Now()

	// Prepare multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add image data
	part, err := writer.CreateFormFile("image", filename)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}

	_, err = part.Write(imageData)
	if err != nil {
		return nil, fmt.Errorf("failed to write image data: %v", err)
	}

	// Add additional parameters
	writer.WriteField("api_key", c.apiKey)
	writer.WriteField("model_version", "v2.1.0")
	writer.WriteField("analysis_type", "comprehensive")
	writer.WriteField("confidence_threshold", "0.7")

	writer.Close()

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/scan", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", "DrMario-Backend/1.0")

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to CNN: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return &CNNScanResult{
			Success: false,
			Error:   fmt.Sprintf("CNN API returned status %d: %s", resp.StatusCode, string(respBody)),
		}, nil
	}

	// Parse JSON response
	var result CNNScanResult
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse CNN response: %v", err)
	}

	// Calculate processing time
	result.ProcessingTime = time.Since(startTime).Seconds()
	result.AnalysisDate = time.Now().Format(time.RFC3339)

	return &result, nil
}

// ValidateImage validates if the image is suitable for CNN analysis
func (c *CNNService) ValidateImage(imagePath string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	// Decode image to get dimensions
	img, format, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// Check supported formats
	if format != "jpeg" && format != "png" {
		return fmt.Errorf("unsupported image format: %s. Only JPEG and PNG are supported", format)
	}

	// Check minimum resolution
	if img.Width < 512 || img.Height < 512 {
		return fmt.Errorf("image resolution too low: %dx%d. Minimum required: 512x512", img.Width, img.Height)
	}

	// Check maximum resolution
	if img.Width > 4096 || img.Height > 4096 {
		return fmt.Errorf("image resolution too high: %dx%d. Maximum allowed: 4096x4096", img.Width, img.Height)
	}

	return nil
}

// PreprocessImage prepares the image for CNN analysis
func (c *CNNService) PreprocessImage(imagePath string) (string, error) {
	// Create preprocessing directory
	preprocessDir := filepath.Join(config.AppConfig.Upload.UploadDir, "preprocessed")
	if err := os.MkdirAll(preprocessDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create preprocess directory: %v", err)
	}

	// Generate output filename
	ext := filepath.Ext(imagePath)
	baseName := filepath.Base(imagePath)
	outputPath := filepath.Join(preprocessDir, "preprocessed_"+baseName)

	// Open source image
	srcFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open source image: %v", err)
	}
	defer srcFile.Close()

	// Decode image
	var img image.Image
	ext = filepath.Ext(imagePath)
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(srcFile)
	case ".png":
		img, err = png.Decode(srcFile)
	default:
		return "", fmt.Errorf("unsupported image format: %s", ext)
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	// Encode with high quality
	var encodeErr error
	switch ext {
	case ".jpg", ".jpeg":
		encodeErr = jpeg.Encode(outFile, img, &jpeg.Options{Quality: 95})
	case ".png":
		encodeErr = png.Encode(outFile, img)
	}

	if encodeErr != nil {
		return "", fmt.Errorf("failed to encode preprocessed image: %v", encodeErr)
	}

	return outputPath, nil
}

// GetCNNHealth checks if the CNN service is available
func (c *CNNService) GetCNNHealth() error {
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("CNN health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CNN service unhealthy, status: %d", resp.StatusCode)
	}

	return nil
}

// Helper functions to get configuration
func getCNNBaseURL() string {
	baseURL := os.Getenv("CNN_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8000/api/v1/cnn" // Default CNN service URL
	}
	return baseURL
}

func getCNNAPIKey() string {
	apiKey := os.Getenv("CNN_API_KEY")
	if apiKey == "" {
		apiKey = "default-cnn-api-key" // Default for development
	}
	return apiKey
}
