package main

import (
	"fmt"
	"log"
	"net/http"

	"cheapeats-api/internal/config"
	"cheapeats-api/internal/database"
	"cheapeats-api/internal/handlers"
	"cheapeats-api/internal/services"

	_ "cheapeats-api/docs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title CheapEats API
// @version 1.0
// @description API for fetching and tracking restaurant prices
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@cheapeats.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @schemes http https

func main() {
	cfg := config.LoadConfig()

	dbConfig := database.DBConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	if err := database.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	apiClient := services.NewRestaurantAPIClient(cfg.API.GooglePlacesAPIKey)
	priceFetcher := services.NewPriceFetcher(apiClient)
	restaurantHandler := handlers.NewRestaurantHandler(priceFetcher)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy"}`))
		})

		r.Route("/restaurants", func(r chi.Router) {
			r.Get("/", restaurantHandler.GetAllRestaurants)
			r.Get("/search", restaurantHandler.SearchNearby)
			r.Get("/{id}", restaurantHandler.GetRestaurant)
			r.Get("/{id}/menu", restaurantHandler.GetMenuItems)
		})

		r.Route("/menu-items", func(r chi.Router) {
			r.Get("/{itemId}", restaurantHandler.GetMenuItem)
			r.Get("/{itemId}/price-history", restaurantHandler.GetPriceHistory)
		})
	})

	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", serverAddr)
	
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}