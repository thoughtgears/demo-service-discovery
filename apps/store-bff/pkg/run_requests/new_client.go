package run_requests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/api/idtoken"

	"github.com/thoughtgears/demo-service-discovery/apps/store-bff/pkg/cfg"
)

type Client struct {
	http        *http.Client
	accessToken string
	BackendURL  string
}

func NewClient(config *cfg.Config, service string) (*Client, error) {
	var client Client

	client.http = &http.Client{
		Timeout: 10 * time.Second,
	}
	if config.BackendURL != "" {
		client.BackendURL = config.BackendURL
	}

	log.Info().Msgf("Requires token: %v", requiresToken(config.Service))

	if config.BackendURL == "" {
		requestURL := fmt.Sprintf("%s/services/%s?environment=%s", config.DiscoveryURL, service, config.Environment)
		req, _ := http.NewRequest(http.MethodGet, requestURL, nil)

		if requiresToken(config.Service) {
			token, err := getToken(config.DiscoveryURL)
			if err != nil {
				return nil, fmt.Errorf("error getting token for discovery: %v", err)
			}
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		}

		resp, err := client.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error doing request to service discovery: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error getting service from discovery endpoint: %v", err)
		}
		defer resp.Body.Close()

		var ServiceDiscovery struct {
			Environment string `json:"environment"`
			URL         string `json:"url"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&ServiceDiscovery); err != nil {
			return nil, fmt.Errorf("error decoding service discovery response: %v", err)
		}

		if ServiceDiscovery.Environment != config.Environment {
			return nil, fmt.Errorf("wrong environment url: %s", ServiceDiscovery.Environment)
		}

		client.BackendURL = ServiceDiscovery.URL
	}

	if requiresToken(config.Service) {
		token, err := getToken(client.BackendURL)
		if err != nil {
			return nil, fmt.Errorf("error getting token for backend: %v", err)
		}

		client.accessToken = token
	}

	return &client, nil
}

func requiresToken(service string) bool {
	return service != "local"
}

func getToken(audience string) (string, error) {
	token, err := idtoken.NewTokenSource(context.TODO(), audience)
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
