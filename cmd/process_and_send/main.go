package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"genre_recommendation/internal/aws"
	"genre_recommendation/internal/config"
	"genre_recommendation/internal/ds"
)

type AggregatedData struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

func main() {
	// Entry point for the Lambda function
	ctx := context.Background()
	err := Handler(ctx, nil)
	if err != nil {
		log.Fatalf("Handler failed: %v", err)
	}
}

func Handler(ctx context.Context, event interface{}) error {
	cfg := config.LoadConfig()
	s3Client := aws.NewS3Client(cfg.AWSRegion)

	// Retrieve asset information from S3
	body, err := s3Client.GetObject(cfg.S3Bucket, cfg.AssetInfoKey)
	if err != nil {
		log.Printf("Error retrieving asset info: %v", err)
		return err
	}
	defer body.Close()

	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("Error reading asset info: %v", err)
		return err
	}

	var assets []ds.AssetInfo
	err = json.Unmarshal(data, &assets)
	if err != nil {
		log.Printf("Error parsing asset info: %v", err)
		return err
	}

	dsClient := ds.NewDSClient(cfg.DSServiceURL)
	var aggregated []AggregatedData

	for _, asset := range assets {
		genres, err := dsClient.PredictGenres(asset)
		if err != nil {
			log.Printf("Error predicting genres for asset %s: %v", asset.ID, err)
			continue
		}
		aggregated = append(aggregated, AggregatedData{
			ID:     asset.ID,
			Title:  asset.Title,
			Genres: genres,
		})
	}

	// Serialize aggregated data and pass to next Lambda
	aggregatedData, err := json.Marshal(aggregated)
	if err != nil {
		log.Printf("Error serializing aggregated data: %v", err)
		return err
	}

	// Use a custom ReadSeekCloser
	fmt.Println("Uploading object to S3 from generated CSV...")
	err = s3Client.PutObject(cfg.S3Bucket, "aggregated_data.json", NewReadSeekCloser(bytes.NewReader(aggregatedData)))
	if err != nil {
		log.Printf("Error uploading aggregated data: %v", err)
		return err
	}

	fmt.Println("Processed assets and sent to DS successfully.")
	return nil
}

// NewReadSeekCloser wraps a bytes.Reader to implement io.ReadSeekCloser
func NewReadSeekCloser(reader *bytes.Reader) io.ReadSeekCloser {
	return &readSeekCloser{Reader: reader}
}

type readSeekCloser struct {
	*bytes.Reader
}

func (r *readSeekCloser) Close() error {
	return nil // No-op for bytes.Reader
}
