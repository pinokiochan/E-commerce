package myrouter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"order-service/internal/adapter/http/myrouter/invdto"
	"order-service/internal/models"
)

type InventoryRouter struct {
	url string
}

func NewInventoryRouter(baseURL string) (*InventoryRouter, error) {
	// Validate the base URL
	_, err := url.ParseRequestURI(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	// Ensure the URL ends with a slash
	if baseURL[len(baseURL)-1] != '/' {
		baseURL += "/"
	}

	return &InventoryRouter{
		url: baseURL + "products/",
	}, nil
}

func (r *InventoryRouter) GetById(id int64) (models.Inventory, error) {
	// Construct the full URL
	fullURL := r.url + fmt.Sprintf("%d", id)

	// Make the HTTP GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return models.Inventory{}, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		return models.Inventory{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Parse the response
	var response invdto.InventoryResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return models.Inventory{}, fmt.Errorf("failed to parse response: %v", err)
	}

	return invdto.ToInventoryModel(response), nil
}

// Sends http PATCH request to change "availability" of product with Header "X-Expected-Version"
func (r *InventoryRouter) Substruct(id, newAvailability int64, version int32) error {
	fullURL := r.url + fmt.Sprintf("%d", id)

	// Create request body
	requestBody := map[string]int64{
		"available": newAvailability,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create new request
	req, err := http.NewRequest(http.MethodPatch, fullURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Expected-Version", fmt.Sprintf("%d", version))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
