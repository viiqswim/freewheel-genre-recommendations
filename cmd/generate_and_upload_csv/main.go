package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"genre_recommendation/internal/aws"
	"genre_recommendation/internal/config"
)

type AggregatedData struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

func main() {
	ctx := context.Background()
	err := Handler(ctx, nil)
	if err != nil {
		log.Fatalf("Handler failed: %v", err)
	}
}

func Handler(ctx context.Context, event interface{}) error {
	cfg := config.LoadConfig()
	s3Client := aws.NewS3Client(cfg.AWSRegion)

	// Retrieve aggregated data from S3
	body, err := s3Client.GetObject(cfg.S3Bucket, "aggregated_data.json")
	if err != nil {
		log.Printf("Error retrieving aggregated data: %v", err)
		return err
	}
	defer body.Close()

	var aggregated []AggregatedData
	decoder := json.NewDecoder(body)
	err = decoder.Decode(&aggregated)
	if err != nil {
		log.Printf("Error parsing aggregated data: %v", err)
		return err
	}

	// Generate CSV data
	csvBuffer := &bytes.Buffer{}
	writer := csv.NewWriter(csvBuffer)

	// Write header
	err = writer.Write([]string{"ID", "Title", "Genres"})
	if err != nil {
		log.Printf("Error writing CSV header: %v", err)
		return err
	}

	// Write records
	for _, record := range aggregated {
		genres := strings.Join(record.Genres, "|")
		err = writer.Write([]string{record.ID, record.Title, genres})
		if err != nil {
			log.Printf("Error writing CSV record: %v", err)
			return err
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		log.Printf("Error flushing CSV writer: %v", err)
		return err
	}

	// Upload CSV to S3
	err = s3Client.PutObject(cfg.S3Bucket, cfg.RecommendationsKey+"recommendations.csv", NewReadSeekCloser(bytes.NewReader(csvBuffer.Bytes())))
	if err != nil {
		log.Printf("Error uploading CSV to S3: %v", err)
		return err
	}

	fmt.Println("Recommendations CSV generated and uploaded successfully.")
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
