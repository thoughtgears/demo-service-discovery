package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/thoughtgears/demo-service-discovery/apps/frontend/run_request"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/appengine"
)

type PageData struct {
	Items []Item
}

var decoder = schema.NewDecoder()

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
	r := mux.NewRouter()
	r.HandleFunc("/", handleIndex).Methods("GET")
	r.HandleFunc("/items", handleItems).Methods("POST")
	http.Handle("/", r)

	appengine.Main()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	tmpl.Execute(w, nil)
}

func handleItems(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	log.Info().Msg("Fetching items from backend")

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
	tmpl := template.Must(template.ParseFiles("static/index.html"))
	data := PageData{Items: responseItems}
	tmpl.Execute(w, data)
}
