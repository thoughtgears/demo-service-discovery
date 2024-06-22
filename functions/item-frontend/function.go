package item_frontend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

type Config struct {
	DiscoveryURL string `envconfig:"DISCOVERY_URL" required:"true"`
	Environment  string `envconfig:"ENVIRONMENT" default:"dev"`
}

type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Currency string `json:"currency"`
}

var config Config

func init() {
	envconfig.MustProcess("", &config)
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.LevelFieldName = "severity"
	functions.HTTP("app", app)
}

func app(w http.ResponseWriter, r *http.Request) {
	backendURL, err := getUrl("item-api", config.Environment)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting backend URL: %v", err), http.StatusInternalServerError)
		return
	}

	log.Info().Msgf("Calling backend API: %s", backendURL)

	// Read count query parameter
	count := r.URL.Query().Get("count")
	if count == "" {
		count = "10"
	}

	resp, err := http.Get(fmt.Sprintf("%s:8080/?count=%s", backendURL, count))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling backend API: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("Backend API returned status code %d", resp.StatusCode), http.StatusInternalServerError)
		return
	}

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding backend response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
	}
}

func getUrl(name, env string) (string, error) {
	url := fmt.Sprintf("%s/services/%s?environment=%s", config.DiscoveryURL, name, env)
	req, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error calling discovery service: %v", err)
	}
	if req.StatusCode != http.StatusOK {
		return "", fmt.Errorf("discovery service returned status code %d", req.StatusCode)
	}

	defer req.Body.Close()

	var resp struct {
		Environment string `json:"environment"`
		Name        string `json:"name"`
		URL         string `json:"url"`
	}
	if err := json.NewDecoder(req.Body).Decode(&resp); err != nil {
		return "", fmt.Errorf("error decoding discovery response: %v", err)
	}

	return resp.URL, nil
}
