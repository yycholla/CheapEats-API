package services

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"cheapeats-api/internal/database"
	"cheapeats-api/internal/models"
)

type PriceFetcher struct {
	apiClient *RestaurantAPIClient
}

func NewPriceFetcher(apiClient *RestaurantAPIClient) *PriceFetcher {
	return &PriceFetcher{
		apiClient: apiClient,
	}
}

func (pf *PriceFetcher) FetchAndSaveRestaurants(ctx context.Context, lat, lng float64, radius int) error {
	searchResp, err := pf.apiClient.SearchRestaurantsByLocation(ctx, lat, lng, radius)
	if err != nil {
		return fmt.Errorf("failed to search restaurants: %w", err)
	}

	db := database.GetDB()

	for _, place := range searchResp.Results {
		restaurant := models.Restaurant{
			ExternalID:  place.PlaceID,
			Name:        place.Name,
			Address:     place.Address,
			Latitude:    place.Geometry.Location.Lat,
			Longitude:   place.Geometry.Location.Lng,
			Rating:      place.Rating,
			PriceRange:  pf.convertPriceLevel(place.PriceLevel),
			CuisineType: pf.extractCuisineType(place.Types),
		}

		addressParts := strings.Split(place.Address, ", ")
		if len(addressParts) >= 3 {
			restaurant.City = addressParts[len(addressParts)-3]
			stateZip := addressParts[len(addressParts)-2]
			stateParts := strings.Split(stateZip, " ")
			if len(stateParts) >= 2 {
				restaurant.State = stateParts[0]
				restaurant.ZipCode = stateParts[1]
			}
			restaurant.Country = addressParts[len(addressParts)-1]
		}

		var existingRestaurant models.Restaurant
		result := db.Where("external_id = ?", restaurant.ExternalID).First(&existingRestaurant)
		
		if result.Error != nil {
			if err := db.Create(&restaurant).Error; err != nil {
				fmt.Printf("Failed to create restaurant %s: %v\n", restaurant.Name, err)
				continue
			}
			existingRestaurant = restaurant
		} else {
			if err := db.Model(&existingRestaurant).Updates(&restaurant).Error; err != nil {
				fmt.Printf("Failed to update restaurant %s: %v\n", restaurant.Name, err)
				continue
			}
		}

		detailsResp, err := pf.apiClient.GetRestaurantDetails(ctx, place.PlaceID)
		if err != nil {
			fmt.Printf("Failed to get details for restaurant %s: %v\n", place.Name, err)
			continue
		}

		if detailsResp.Result.PhoneNumber != "" {
			existingRestaurant.Phone = detailsResp.Result.PhoneNumber
		}
		if detailsResp.Result.Website != "" {
			existingRestaurant.Website = detailsResp.Result.Website
		}
		
		db.Save(&existingRestaurant)

		pf.generateSampleMenuItems(existingRestaurant.ID, place.PriceLevel)

		scrapedData := models.ScrapedData{
			Source:       "google_places",
			RestaurantID: &existingRestaurant.ID,
			RawData: models.JSONB{
				"search_result": place,
				"details":       detailsResp.Result,
			},
			ScrapedAt: time.Now(),
		}
		db.Create(&scrapedData)

		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (pf *PriceFetcher) generateSampleMenuItems(restaurantID uint, priceLevel int) {
	db := database.GetDB()
	
	basePrice := 10.0
	if priceLevel > 0 {
		basePrice = float64(priceLevel) * 15.0
	}

	categories := []string{"Appetizers", "Main Course", "Desserts", "Beverages"}
	
	menuItems := []models.MenuItem{
		{
			RestaurantID: restaurantID,
			Name:         "Signature Appetizer",
			Description:  "Chef's special starter",
			Category:     categories[0],
			Price:        basePrice * 0.7 + rand.Float64()*5,
			Currency:     "USD",
			IsAvailable:  true,
		},
		{
			RestaurantID: restaurantID,
			Name:         "House Special Main",
			Description:  "Most popular main dish",
			Category:     categories[1],
			Price:        basePrice + rand.Float64()*10,
			Currency:     "USD",
			IsAvailable:  true,
		},
		{
			RestaurantID: restaurantID,
			Name:         "Daily Special",
			Description:  "Today's featured dish",
			Category:     categories[1],
			Price:        basePrice * 1.2 + rand.Float64()*8,
			Currency:     "USD",
			IsAvailable:  true,
		},
		{
			RestaurantID: restaurantID,
			Name:         "Classic Burger",
			Description:  "Traditional burger with fries",
			Category:     categories[1],
			Price:        basePrice * 0.9 + rand.Float64()*5,
			Currency:     "USD",
			IsAvailable:  true,
		},
		{
			RestaurantID: restaurantID,
			Name:         "Dessert of the Day",
			Description:  "Sweet treat to end your meal",
			Category:     categories[2],
			Price:        basePrice * 0.5 + rand.Float64()*3,
			Currency:     "USD",
			IsAvailable:  true,
		},
		{
			RestaurantID: restaurantID,
			Name:         "Soft Drink",
			Description:  "Various sodas available",
			Category:     categories[3],
			Price:        3.50 + rand.Float64()*2,
			Currency:     "USD",
			IsAvailable:  true,
		},
	}

	for _, item := range menuItems {
		var existingItem models.MenuItem
		result := db.Where("restaurant_id = ? AND name = ?", item.RestaurantID, item.Name).First(&existingItem)
		
		if result.Error != nil {
			db.Create(&item)
			
			priceHistory := models.PriceHistory{
				MenuItemID: item.ID,
				Price:      item.Price,
				RecordedAt: time.Now(),
			}
			db.Create(&priceHistory)
		} else {
			if existingItem.Price != item.Price {
				priceHistory := models.PriceHistory{
					MenuItemID: existingItem.ID,
					Price:      item.Price,
					RecordedAt: time.Now(),
				}
				db.Create(&priceHistory)
				
				db.Model(&existingItem).Update("price", item.Price)
			}
		}
	}
}

func (pf *PriceFetcher) convertPriceLevel(level int) string {
	switch level {
	case 0:
		return "Free"
	case 1:
		return "$"
	case 2:
		return "$$"
	case 3:
		return "$$$"
	case 4:
		return "$$$$"
	default:
		return "N/A"
	}
}

func (pf *PriceFetcher) extractCuisineType(types []string) string {
	cuisineTypes := map[string]string{
		"chinese_restaurant":  "Chinese",
		"italian_restaurant":  "Italian",
		"mexican_restaurant":  "Mexican",
		"japanese_restaurant": "Japanese",
		"indian_restaurant":   "Indian",
		"thai_restaurant":     "Thai",
		"french_restaurant":   "French",
		"pizza":              "Pizza",
		"burger":             "Burger",
		"seafood":            "Seafood",
		"vegetarian":         "Vegetarian",
		"cafe":               "Cafe",
		"bakery":             "Bakery",
		"bar":                "Bar",
	}

	for _, t := range types {
		if cuisine, ok := cuisineTypes[t]; ok {
			return cuisine
		}
	}

	if contains(types, "restaurant") {
		return "General"
	}

	return "Other"
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}