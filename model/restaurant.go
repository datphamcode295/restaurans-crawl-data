package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Create a GORM model for the restaurants table
type Restaurant struct {
	gorm.Model
	ID           uuid.UUID `gorm:"type:primary_key;default:uuid_generate_v4()"`
	Name         string
	Avatar       string
	Phone        string
	Slug         string
	Street       string
	District     string
	City         string
	FullAddress  string
	Lat          float64
	Long         float64
	IsOpening24h bool `gorm:"column:is_opening_24h"`
	ExternalId   string
}
