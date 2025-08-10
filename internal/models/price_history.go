package models

import (
	"time"
)

type PriceHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MenuItemID uint      `gorm:"index" json:"menu_item_id"`
	MenuItem   *MenuItem `json:"menu_item,omitempty"`
	Price      float64   `gorm:"not null" json:"price"`
	RecordedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"recorded_at"`
}