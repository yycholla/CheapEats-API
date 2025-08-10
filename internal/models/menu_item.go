package models

import (
	"time"

	"gorm.io/gorm"
)

type MenuItem struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	RestaurantID uint           `gorm:"index" json:"restaurant_id"`
	Restaurant   *Restaurant    `json:"restaurant,omitempty"`
	Name         string         `gorm:"not null;size:255" json:"name"`
	Description  string         `json:"description"`
	Category     string         `gorm:"size:100" json:"category"`
	Price        float64        `gorm:"not null;index" json:"price"`
	Currency     string         `gorm:"default:USD;size:10" json:"currency"`
	IsAvailable  bool           `gorm:"default:true" json:"is_available"`
	PriceHistory []PriceHistory `gorm:"foreignKey:MenuItemID" json:"price_history,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}