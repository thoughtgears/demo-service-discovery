package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/thoughtgears/demo-service-discovery/apps/frontend/run_request"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/appengine"
)

type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Currency string `json:"currency" default:"USD"`
}

func init() {
	zerolog.LevelFieldName = "severity"
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = time.RFC3339Nano
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/items", handleItems)

	if os.Getenv("GAE_SERVICE") == "local" {
		// Running locally
		port := "8080"
		if p := os.Getenv("PORT"); p != "" {
			port = p
		}
		log.Info().Msgf("Starting server on port %s", port)
		log.Error().Err(http.ListenAndServe(":"+port, nil)).Msg("Error starting server")
	}

	appengine.Main()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	client, err := run_request.NewClient("store-bff")
	if err != nil {
		log.Error().Err(err).Msg("Error creating client")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	countStr := r.FormValue("count")
	count, err := strconv.Atoi(countStr)
	if err != nil || count <= 0 {
		count = 20
	}

	path := fmt.Sprintf("/items?count=%d", count)
	resp, err := client.Do(ctx, http.MethodGet, path)
	if err != nil {
		log.Error().Err(err).Msg("Error fetching items from backend")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().Msgf("Unexpected status code: %d", resp.StatusCode)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var responseItems []Item
	if err := json.NewDecoder(resp.Body).Decode(&responseItems); err != nil {
		log.Error().Err(err).Msg("Error decoding response")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
