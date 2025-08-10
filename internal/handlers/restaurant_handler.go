package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cheapeats-api/internal/database"
	"cheapeats-api/internal/models"
	"cheapeats-api/internal/services"

	"github.com/go-chi/chi/v5"
)

type RestaurantHandler struct {
	priceFetcher *services.PriceFetcher
}

func NewRestaurantHandler(priceFetcher *services.PriceFetcher) *RestaurantHandler {
	return &RestaurantHandler{
		priceFetcher: priceFetcher,
	}
}

// GetAllRestaurants godoc
// @Summary List all restaurants
// @Description Get a list of all restaurants with optional filters
// @Tags restaurants
// @Accept json
// @Produce json
// @Param city query string false "Filter by city"
// @Param cuisine query string false "Filter by cuisine type"
// @Param price_range query string false "Filter by price range (e.g., $, $$, $$$, $$$$)"
// @Success 200 {array} models.Restaurant
// @Failure 500 {object} map[string]string
// @Router /restaurants [get]
func (h *RestaurantHandler) GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	db := database.GetDB()
	
	var restaurants []models.Restaurant
	
	query := db.Model(&models.Restaurant{})
	
	if city := r.URL.Query().Get("city"); city != "" {
		query = query.Where("city = ?", city)
	}
	
	if cuisine := r.URL.Query().Get("cuisine"); cuisine != "" {
		query = query.Where("cuisine_type = ?", cuisine)
	}
	
	if priceRange := r.URL.Query().Get("price_range"); priceRange != "" {
		query = query.Where("price_range = ?", priceRange)
	}
	
	if err := query.Find(&restaurants).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch restaurants")
		return
	}
	
	respondWithJSON(w, http.StatusOK, restaurants)
}

// GetRestaurant godoc
// @Summary Get restaurant by ID
// @Description Get detailed information about a specific restaurant including menu items
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Success 200 {object} models.Restaurant
// @Failure 404 {object} map[string]string
// @Router /restaurants/{id} [get]
func (h *RestaurantHandler) GetRestaurant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	db := database.GetDB()
	var restaurant models.Restaurant
	
	if err := db.Preload("MenuItems").First(&restaurant, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, restaurant)
}

// SearchNearby godoc
// @Summary Search nearby restaurants
// @Description Search for restaurants within a radius of given coordinates
// @Tags restaurants
// @Accept json
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param radius query int false "Search radius in meters (default: 1000)"
// @Success 200 {array} models.Restaurant
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /restaurants/search [get]
func (h *RestaurantHandler) SearchNearby(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lngStr := r.URL.Query().Get("lng")
	radiusStr := r.URL.Query().Get("radius")
	
	if latStr == "" || lngStr == "" {
		respondWithError(w, http.StatusBadRequest, "Latitude and longitude are required")
		return
	}
	
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid latitude")
		return
	}
	
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid longitude")
		return
	}
	
	radius := 1000
	if radiusStr != "" {
		radius, err = strconv.Atoi(radiusStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid radius")
			return
		}
	}
	
	if err := h.priceFetcher.FetchAndSaveRestaurants(r.Context(), lat, lng, radius); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch restaurants")
		return
	}
	
	db := database.GetDB()
	var restaurants []models.Restaurant
	
	earthRadius := 6371.0
	query := db.Raw(`
		SELECT * FROM restaurants 
		WHERE (
			? * acos(
				cos(radians(?)) * 
				cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * 
				sin(radians(latitude))
			)
		) <= ?
		ORDER BY (
			? * acos(
				cos(radians(?)) * 
				cos(radians(latitude)) * 
				cos(radians(longitude) - radians(?)) + 
				sin(radians(?)) * 
				sin(radians(latitude))
			)
		)`,
		earthRadius, lat, lng, lat, float64(radius)/1000.0,
		earthRadius, lat, lng, lat,
	)
	
	if err := query.Scan(&restaurants).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch nearby restaurants")
		return
	}
	
	respondWithJSON(w, http.StatusOK, restaurants)
}

// GetMenuItems godoc
// @Summary Get restaurant menu items
// @Description Get all menu items for a specific restaurant
// @Tags restaurants
// @Accept json
// @Produce json
// @Param id path int true "Restaurant ID"
// @Param category query string false "Filter by category"
// @Param max_price query number false "Maximum price filter"
// @Success 200 {array} models.MenuItem
// @Failure 500 {object} map[string]string
// @Router /restaurants/{id}/menu [get]
func (h *RestaurantHandler) GetMenuItems(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "id")
	
	db := database.GetDB()
	var menuItems []models.MenuItem
	
	query := db.Where("restaurant_id = ?", restaurantID)
	
	if category := r.URL.Query().Get("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	
	if maxPrice := r.URL.Query().Get("max_price"); maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)
		if err == nil {
			query = query.Where("price <= ?", price)
		}
	}
	
	if err := query.Find(&menuItems).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch menu items")
		return
	}
	
	respondWithJSON(w, http.StatusOK, menuItems)
}

// GetMenuItem godoc
// @Summary Get menu item by ID
// @Description Get detailed information about a specific menu item including price history
// @Tags menu-items
// @Accept json
// @Produce json
// @Param itemId path int true "Menu Item ID"
// @Success 200 {object} models.MenuItem
// @Failure 404 {object} map[string]string
// @Router /menu-items/{itemId} [get]
func (h *RestaurantHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "itemId")
	
	db := database.GetDB()
	var menuItem models.MenuItem
	
	if err := db.Preload("PriceHistory").First(&menuItem, id).Error; err != nil {
		respondWithError(w, http.StatusNotFound, "Menu item not found")
		return
	}
	
	respondWithJSON(w, http.StatusOK, menuItem)
}

// GetPriceHistory godoc
// @Summary Get price history for menu item
// @Description Get the price history of a specific menu item
// @Tags menu-items
// @Accept json
// @Produce json
// @Param itemId path int true "Menu Item ID"
// @Success 200 {array} models.PriceHistory
// @Failure 500 {object} map[string]string
// @Router /menu-items/{itemId}/price-history [get]
func (h *RestaurantHandler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "itemId")
	
	db := database.GetDB()
	var priceHistory []models.PriceHistory
	
	if err := db.Where("menu_item_id = ?", itemID).
		Order("recorded_at DESC").
		Find(&priceHistory).Error; err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch price history")
		return
	}
	
	respondWithJSON(w, http.StatusOK, priceHistory)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}