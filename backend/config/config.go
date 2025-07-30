package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	JWT    JWTConfig
	Upload UploadConfig
	AI     AIConfig
	CORS   CORSConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type JWTConfig struct {
	Secret string
	Expiry string
}

type UploadConfig struct {
	MaxFileSize       int64
	UploadDir         string
	AllowedExtensions []string
}

type AIConfig struct {
	ModelPath           string
	ConfidenceThreshold float64
}

type CORSConfig struct {
	AllowedOrigins []string
}

var AppConfig Config

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		// If .env file doesn't exist, continue with system environment variables
	}

	AppConfig = Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "default-secret-key"),
			Expiry: getEnv("JWT_EXPIRY", "24h"),
		},
		Upload: UploadConfig{
			MaxFileSize:       getEnvAsInt64("MAX_FILE_SIZE", 10485760), // 10MB
			UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
			AllowedExtensions: getEnvAsSlice("ALLOWED_EXTENSIONS", []string{"jpg", "jpeg", "png", "tiff", "bmp"}),
		},
		AI: AIConfig{
			ModelPath:           getEnv("MODEL_PATH", "./models/dr_detection_model"),
			ConfidenceThreshold: getEnvAsFloat("CONFIDENCE_THRESHOLD", 0.7),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3000", "http://localhost:5173"}),
		},
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return defaultValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values
		// In production, you might want more sophisticated parsing
		return []string{value}
	}
	return defaultValue
}
