package db

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    string `json:"price"`
	Currency string `json:"currency" default:"USD"`
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

func GetItems(n int) []Item {
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
