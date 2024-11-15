// internal/config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the configuration variables
type Config struct {
	AWSRegion          string
	S3Bucket           string
	AssetInfoKey       string
	RecommendationsKey string
	DSServiceURL       string
}

// LoadConfig reads configuration from .env file or environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables.")
	}

	config := &Config{
		AWSRegion:          getEnv("AWS_REGION", "us-east-1"),
		S3Bucket:           getEnv("S3_BUCKET", "genre-recommendations"),
		AssetInfoKey:       getEnv("ASSET_INFO_KEY", "assets/asset_info.json"),
		RecommendationsKey: getEnv("RECOMMENDATIONS_KEY", "recommendations/"),
		DSServiceURL:       getEnv("DS_SERVICE_URL", "http://localhost:8080/predict"),
	}

	return config
}

// getEnv retrieves the value of the environment variable named by the key.
// It returns the fallback value if the variable is not present.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
