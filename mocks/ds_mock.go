package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

// AssetInfo represents the incoming asset structure
type AssetInfo struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	// Add other relevant fields if necessary
}

// GenrePrediction represents the genre predictions structure
type GenrePrediction struct {
	Genres []string `json:"genres"`
}

// Handler for the /predict endpoint
func predictHandler(w http.ResponseWriter, r *http.Request) {
	var asset AssetInfo

	// Decode the incoming JSON request into the AssetInfo struct
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Generate mock genres based on asset ID or Title (for demonstration)
	genres := generateMockGenres(asset)

	// Create the GenrePrediction response
	prediction := GenrePrediction{
		Genres: genres,
	}

	// Set the response header to JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prediction)
}

// generateMockGenres creates a mock list of genres based on asset information
func generateMockGenres(asset AssetInfo) []string {
	// List of possible genres
	genresList := []string{"Comedy", "Drama", "Action", "Thriller", "Sci-Fi", "Adventure", "Documentary", "History", "Fantasy", "Horror"}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Shuffle the genres list
	rand.Shuffle(len(genresList), func(i, j int) {
		genresList[i], genresList[j] = genresList[j], genresList[i]
	})

	// Determine a random number of genres to return (between 1 and len(genresList))
	numGenres := rand.Intn(len(genresList)) + 1

	// Return a slice of the shuffled list with the random number of genres
	return genresList[:numGenres]
}

func main() {
	// Get the port from environment variables or default to 9090
	port := os.Getenv("DS_PORT")
	if port == "" {
		port = "9090"
	}

	http.HandleFunc("/predict", predictHandler)

	log.Printf("Mock DS service running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start mock DS server: %v", err)
	}
}
