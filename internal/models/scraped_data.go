package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if data, ok := value.(string); ok {
		return json.Unmarshal([]byte(data), j)
	}
	if data, ok := value.([]byte); ok {
		return json.Unmarshal(data, j)
	}
	return nil
}

type ScrapedData struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	Source       string      `gorm:"not null;size:100;index" json:"source"`
	RestaurantID *uint       `json:"restaurant_id,omitempty"`
	Restaurant   *Restaurant `json:"restaurant,omitempty"`
	RawData      JSONB       `gorm:"type:jsonb" json:"raw_data"`
	ScrapedAt    time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"scraped_at"`
}