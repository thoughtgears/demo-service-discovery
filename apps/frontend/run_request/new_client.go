package run_request

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/api/idtoken"
)

type Client struct {
	http        *http.Client
	accessToken string
	BackendURL  string
}

func NewClient(service string) (*Client, error) {
	client := &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	discoveryURL := os.Getenv("DISCOVERY_URL")
	if discoveryURL == "" {
		return nil, fmt.Errorf("DISCOVERY_URL not set")
	}

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		return nil, errors.New("ENVIRONMENT not set")
	}

	client.BackendURL = os.Getenv("BACKEND_URL")
	gaeService := os.Getenv("GAE_SERVICE")

	log.Info().Msgf("service: %s, discoveryURL: %s, environment: %s, backendURL: %s, gaeService: %s", service, discoveryURL, environment, client.BackendURL, gaeService)

	if client.BackendURL == "" {
		if requiresToken(gaeService) {
			token, err := getToken(discoveryURL)
			if err != nil {
				return nil, fmt.Errorf("error getting token for discovery: %v", err)
			}
			client.accessToken = token
		}

		requestURL := fmt.Sprintf("%s/services/%s?environment=%s", discoveryURL, service, environment)
		req, _ := http.NewRequest(http.MethodGet, requestURL, nil)
		if requiresToken(gaeService) {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.accessToken))
		}

		resp, err := client.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error doing request to service discovery: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error getting service from discovery endpoint: %v", resp.Status)
		}

		var ServiceDiscovery struct {
			Environment string `json:"environment"`
			URL         string `json:"url"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&ServiceDiscovery); err != nil {
			return nil, fmt.Errorf("error decoding service discovery response: %v", err)
		}

		if ServiceDiscovery.Environment != environment {
			return nil, fmt.Errorf("wrong environment url: %s", ServiceDiscovery.Environment)
		}

		client.BackendURL = ServiceDiscovery.URL
	}

	if requiresToken(gaeService) {
		token, err := getToken(client.BackendURL)
		if err != nil {
			return nil, fmt.Errorf("error getting token for backend: %v", err)
		}
		client.accessToken = token
	}

	return client, nil
}

func requiresToken(service string) bool {
	return service != ""
}

func getToken(audience string) (string, error) {
	token, err := idtoken.NewTokenSource(context.Background(), audience)
	if err != nil {
		return "", fmt.Errorf("idtoken.NewTokenSource: %w", err)
	}

	t, err := token.Token()
	if err != nil {
		return "", fmt.Errorf("token.Token: %w", err)
	}

	return t.AccessToken, nil
}

func (c *Client) Do(ctx context.Context, method, path string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.BackendURL, path)
	req, _ := http.NewRequest(method, url, nil)
	req.WithContext(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	return c.http.Do(req)
}
