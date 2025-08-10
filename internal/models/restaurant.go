package models

import (
	"time"

	"gorm.io/gorm"
)

type Restaurant struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ExternalID  string         `gorm:"uniqueIndex;size:255" json:"external_id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Address     string         `json:"address"`
	City        string         `gorm:"size:100" json:"city"`
	State       string         `gorm:"size:50" json:"state"`
	ZipCode     string         `gorm:"size:20" json:"zip_code"`
	Country     string         `gorm:"size:100" json:"country"`
	Latitude    float64        `gorm:"index:idx_location" json:"latitude"`
	Longitude   float64        `gorm:"index:idx_location" json:"longitude"`
	CuisineType string         `gorm:"size:100" json:"cuisine_type"`
	Phone       string         `gorm:"size:50" json:"phone"`
	Website     string         `gorm:"size:255" json:"website"`
	Rating      float32        `json:"rating"`
	PriceRange  string         `gorm:"size:10" json:"price_range"`
	MenuItems   []MenuItem     `gorm:"foreignKey:RestaurantID" json:"menu_items,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}