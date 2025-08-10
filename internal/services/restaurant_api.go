package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type RestaurantAPIClient struct {
	httpClient *http.Client
	apiKey     string
}

func NewRestaurantAPIClient(apiKey string) *RestaurantAPIClient {
	return &RestaurantAPIClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

type PlaceSearchResponse struct {
	Results []PlaceResult `json:"results"`
	Status  string        `json:"status"`
}

type PlaceResult struct {
	PlaceID      string   `json:"place_id"`
	Name         string   `json:"name"`
	Address      string   `json:"formatted_address"`
	Geometry     Geometry `json:"geometry"`
	Rating       float32  `json:"rating"`
	PriceLevel   int      `json:"price_level"`
	Types        []string `json:"types"`
	OpeningHours *struct {
		OpenNow bool `json:"open_now"`
	} `json:"opening_hours,omitempty"`
}

type Geometry struct {
	Location struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"location"`
}

type PlaceDetailsResponse struct {
	Result PlaceDetails `json:"result"`
	Status string       `json:"status"`
}

type PlaceDetails struct {
	PlaceID         string   `json:"place_id"`
	Name            string   `json:"name"`
	FormattedAddress string  `json:"formatted_address"`
	PhoneNumber     string   `json:"formatted_phone_number"`
	Website         string   `json:"website"`
	Rating          float32  `json:"rating"`
	PriceLevel      int      `json:"price_level"`
	Types           []string `json:"types"`
	Geometry        Geometry `json:"geometry"`
}

func (c *RestaurantAPIClient) SearchRestaurantsByLocation(ctx context.Context, lat, lng float64, radius int) (*PlaceSearchResponse, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	
	params := url.Values{}
	params.Add("location", fmt.Sprintf("%f,%f", lat, lng))
	params.Add("radius", fmt.Sprintf("%d", radius))
	params.Add("type", "restaurant")
	params.Add("key", c.apiKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PlaceSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "OK" && result.Status != "ZERO_RESULTS" {
		return nil, fmt.Errorf("API returned error status: %s", result.Status)
	}

	return &result, nil
}

func (c *RestaurantAPIClient) GetRestaurantDetails(ctx context.Context, placeID string) (*PlaceDetailsResponse, error) {
	baseURL := "https://maps.googleapis.com/maps/api/place/details/json"
	
	params := url.Values{}
	params.Add("place_id", placeID)
	params.Add("fields", "place_id,name,formatted_address,formatted_phone_number,website,rating,price_level,types,geometry")
	params.Add("key", c.apiKey)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result PlaceDetailsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if result.Status != "OK" {
		return nil, fmt.Errorf("API returned error status: %s", result.Status)
	}

	return &result, nil
}