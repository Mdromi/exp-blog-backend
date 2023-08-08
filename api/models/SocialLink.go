package models

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// SocialLink represents social links for a user's profile
type SocialLink struct {
	gorm.Model
	Facebook  string `json:"facebook"`
	Twitter   string `json:"twitter"`
	Instagram string `json:"instagram"`
	// Add other social media fields as needed
}

// Implement Valuer interface to convert SocialLink to a JSON-encoded string when saving to the database
func (sl SocialLink) Value() (driver.Value, error) {
	return json.Marshal(sl)
}

// Implement Scanner interface to convert a JSON-encoded string from the database to a SocialLink object
func (sl *SocialLink) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), sl)
}
