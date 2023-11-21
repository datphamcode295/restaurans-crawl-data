package model

import (
	"time"

	"gorm.io/gorm"
)

// Create a GORM model for the restaurants table
type Restaurant struct {
	gorm.Model
	Name                   string
	Avatar                 string
	Phone                  string
	CountryCode            string
	Slug                   string
	Street                 string
	District               string
	City                   string
	FullAddress            string
	Rating                 float64
	Username               string
	Lat                    float64
	Long                   float64
	IsOpening              bool
	IsOpening24h           bool `gorm:"column:is_opening_24h"`
	MinutesUntilNextStatus int
	PromotedAt             time.Time
	IsLoshipPartner        bool
	IsHonored              bool
	Quote                  string
	IsActive               bool
	IsCheckedIn            bool
	Closed                 bool
	RecommendedRatio       float64
	RecommendedEnable      bool
	Distance               int
	IsPurchasedSupplyItems bool
	IsSponsored            bool
	FreeShippingMilestone  int
}
