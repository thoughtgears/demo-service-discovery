package item_api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Currency string `json:"currency" default:"USD"`
}

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.LevelFieldName = "severity"
	functions.HTTP("app", app)
}

func app(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := 30
	queryCount := r.URL.Query().Get("count")
	if queryCount != "" {
		if parsedCount, err := strconv.Atoi(queryCount); err == nil {
			count = parsedCount
		}
	}

	items := getItems(count)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(items); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (i Item) MarshalJSON() ([]byte, error) {
	type Alias Item
	return json.Marshal(&struct {
		Price string `json:"price"`
		*Alias
	}{
		Price: i.Price,
		Alias: (*Alias)(&i),
	})
}

func getRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func getRandomElement(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func generateRandomName(prefixes, suffixes []string) string {
	return strings.Join([]string{getRandomElement(prefixes), getRandomElement(suffixes)}, " ")
}

func getItems(n int) []Item {
	prefix := []string{"Amazing", "Big", "Small", "Fancy", "Shiny", "Cool", "Fast", "Sleek"}
	suffix := []string{"Tool", "Stool", "Gadget", "Widget", "Device", "Contraption", "Item", "Thing"}
	itemTypes := []string{"Electronics", "Furniture", "Toys", "Clothing", "Books", "Tools", "Sports", "Groceries"}
	rand.New(rand.NewSource(time.Now().UnixNano()))

	items := make([]Item, n)
	for i := 0; i < n; i++ {
		priceFloat := getRandomFloat(10, 100000)
		items[i] = Item{
			ID:       uuid.New().String(),
			Name:     generateRandomName(prefix, suffix),
			Type:     getRandomElement(itemTypes),
			Price:    fmt.Sprintf("%.2f", priceFloat),
			Currency: "USD",
		}
	}

	return items
}
