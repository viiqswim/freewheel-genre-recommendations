// internal/ds/ds_client.go
package ds

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// DSClient handles communication with the DS service
type DSClient struct {
	ServiceURL string
	Client     *http.Client
}

// NewDSClient initializes and returns a new DSClient
func NewDSClient(serviceURL string) *DSClient {
	return &DSClient{
		ServiceURL: serviceURL,
		Client:     &http.Client{},
	}
}

// AssetInfo represents the structure of each asset
type AssetInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Add other relevant fields if necessary
}

// GenrePrediction represents the structure of genre predictions
type GenrePrediction struct {
	Genres []string `json:"genres"`
}

// PredictGenres sends asset information to the DS service and returns genre predictions
func (ds *DSClient) PredictGenres(asset AssetInfo) ([]string, error) {
	data, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	resp, err := ds.Client.Post(ds.ServiceURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Printf("DS service returned status %d: %s", resp.StatusCode, string(bodyBytes))
		return nil, err
	}

	var prediction GenrePrediction
	err = json.NewDecoder(resp.Body).Decode(&prediction)
	if err != nil {
		return nil, err
	}

	return prediction.Genres, nil
}
