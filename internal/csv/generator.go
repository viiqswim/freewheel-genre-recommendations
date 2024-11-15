// internal/csv/generator.go
package csv

import (
	"encoding/csv"
	"io"
	"strings"
)

// AggregatedData represents the combined asset information and genre predictions
type AggregatedData struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Genres []string `json:"genres"`
}

// GenerateCSV creates a CSV from aggregated data and writes it to the provided writer
func GenerateCSV(data []AggregatedData, writer io.Writer) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	header := []string{"ID", "Title", "Genres"}
	if err := csvWriter.Write(header); err != nil {
		return err
	}

	// Write records
	for _, record := range data {
		genres := strings.Join(record.Genres, "|") // Use "|" as a separator for genres
		csvRecord := []string{record.ID, record.Title, genres}
		if err := csvWriter.Write(csvRecord); err != nil {
			return err
		}
	}

	return nil
}
