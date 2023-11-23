package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
	"net/http"
	"sync"

	"example.com/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	threadMax   = 10
	apiEndpoint = "https://mocha.lozi.vn/v6.1/search/eateries/near-by?cityId=50&limit=24&superCategoryId=1&lat=10.7765194&lng=106.700987&page="
)

var (
	totalPages = 1
)

// Address struct to represent the address field
type Address struct {
	Street   string `json:"street"`
	District string `json:"district"`
	City     string `json:"city"`
	Full     string `json:"full"`
}

// OperatingTime struct to represent the operatingTime field
type OperatingTime struct {
	Start   string `json:"start"`
	Finish  string `json:"finish"`
	Weekday string `json:"weekday"`
}

// Promotion struct to represent the promotions field
type Promotion struct {
	ID                             int    `json:"id"`
	Code                           string `json:"code"`
	PromotionType                  string `json:"promotionType"`
	Value                          int    `json:"value"`
	MaxDiscount                    int    `json:"maxDiscount"`
	ActiveFrom                     string `json:"activeFrom"`
	ActiveTo                       string `json:"activeTo"`
	IsFirstTimeCode                bool   `json:"isFirstTimeCode"`
	IsFirstTimeOrder               bool   `json:"isFirstTimeOrder"`
	PerUserActiveTimes             int    `json:"perUserActiveTimes"`
	PerUserDailyActiveTimes        int    `json:"perUserDailyActiveTimes"`
	ClientSupportOnly              bool   `json:"clientSupportOnly"`
	IsLimitMinimumOrderValue       bool   `json:"isLimitMinimumOrderValue"`
	MinimumOrderValue              int    `json:"minimumOrderValue"`
	IsLimitPaymentType             bool   `json:"isLimitPaymentType"`
	IsEateryHavingCommissionOnly   bool   `json:"IsEateryHavingCommissionOnly"`
	IsGroupPromotion               bool   `json:"isGroupPromotion"`
	MinimumUserToApply             int    `json:"minimumUserToApply"`
	DailyActiveFrom                int    `json:"dailyActiveFrom"`
	DailyActiveTo                  int    `json:"dailyActiveTo"`
	IsHaveUnsupportedDishes        bool   `json:"isHaveUnsupportedDishes"`
	PromotionUsagePercentage       int    `json:"promotionUsagePercentage"`
	PromotionUsagePercentageEnable bool   `json:"promotionUsagePercentageEnable"`
}

// ApiResponseData struct to represent the Data field in the API response
type ApiResponseData struct {
	ID              int             `json:"id"`
	Name            string          `json:"name"`
	Avatar          string          `json:"avatar"`
	Phone           string          `json:"phone"`
	CountryCode     string          `json:"countryCode"`
	Slug            string          `json:"slug"`
	Address         Address         `json:"address"`
	Rating          float64         `json:"rating"`
	Username        string          `json:"username"`
	Lat             float64         `json:"lat"`
	Long            float64         `json:"long"`
	OperatingTime   []OperatingTime `json:"operatingTime"`
	OperatingStatus struct {
		IsOpening              bool `json:"isOpening"`
		IsOpening24h           bool `json:"isOpening24h"`
		MinutesUntilNextStatus int  `json:"minutesUntilNextStatus"`
	} `json:"operatingStatus"`
	PromotedAt             string      `json:"promotedAt"`
	Promotions             []Promotion `json:"promotions"`
	IsLoshipPartner        bool        `json:"isLoshipPartner"`
	IsHonored              bool        `json:"isHonored"`
	Quote                  string      `json:"quote"`
	IsActive               bool        `json:"isActive"`
	IsCheckedIn            bool        `json:"isCheckedIn"`
	Closed                 bool        `json:"closed"`
	RecommendedRatio       float64     `json:"recommendedRatio"`
	RecommendedEnable      bool        `json:"recommendedEnable"`
	Distance               int         `json:"distance"`
	IsPurchasedSupplyItems bool        `json:"isPurchasedSupplyItems"`
	IsSponsored            bool        `json:"isSponsored"`
	FreeShippingMilestone  int         `json:"freeShippingMilestone"`
}

// ApiResponse struct to represent the entire API response
type ApiResponse struct {
	Data       []ApiResponseData `json:"Data"`
	Pagination struct {
		Total   int    `json:"total"`
		Page    int    `json:"page"`
		Limit   int    `json:"limit"`
		NextURL string `json:"nextUrl"`
	} `json:"pagination"`
}

func fetchData(page int) (ApiResponse, error) {
	apiURL := apiEndpoint + fmt.Sprint(page)

	// Make the HTTP GET request
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error making the request for page %d: %v\n", page, err)
		return ApiResponse{}, err
	}
	defer response.Body.Close()

	// Decode the JSON response
	var apiResponse ApiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		fmt.Printf("Error decoding JSON for page %d: %v\n", page, err)
		return ApiResponse{}, err
	}
	return apiResponse, nil
}
func saveData(page int, wg *sync.WaitGroup) {
	defer wg.Done()

	apiResponse, err := fetchData(page)
	if err != nil {
		fmt.Printf("Error fetching data for page %d: %v\n", page, err)
		return
	}

	for _, data := range apiResponse.Data {
		err = saveByExternalId(data)
		if err != nil {
			fmt.Printf("Error saving to DB for page %d: %v\n", page, err)
		}
	}

	fmt.Printf("Pagination for page %d:\n", page)
	// fmt.Printf("  Total: %d\n", apiResponse.Pagination.Total)
	// fmt.Printf("  Page: %d\n", apiResponse.Pagination.Page)
	// fmt.Printf("  Limit: %d\n", apiResponse.Pagination.Limit)
	// fmt.Printf("  NextURL: %s\n", apiResponse.Pagination.NextURL)
}

var db *gorm.DB

const (
	dbHost     = "localhost"
	dbPort     = 5436
	dbUser     = "postgres"
	dbPassword = "postgres"
	dbName     = "food-sharing-db"
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// Calculate the total number of pages
	apiResponse, err := fetchData(1)
	if err != nil {
		fmt.Printf("Error fetching data for page %d: %v\n", 1, err)
		return
	}
	totalPages = int(math.Ceil(float64(apiResponse.Pagination.Total / apiResponse.Pagination.Limit)))

	// get all page data concurrently
	var wg sync.WaitGroup

	for page := 1; page <= totalPages; page++ {
		// Limit the number of concurrent threads
		if page%threadMax == 0 {
			wg.Wait()
		}

		wg.Add(1)
		go saveData(page, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func saveByExternalId(data ApiResponseData) error {
	var restaurant model.Restaurant
	dbError := db.Where("external_id = ?", fmt.Sprint(data.ID)).First(&restaurant).Error
	if restaurant.ID != uuid.Nil || dbError == nil {
		restaurant.Name = data.Name
		restaurant.Avatar = data.Avatar
		restaurant.Phone = data.Phone
		restaurant.Slug = data.Slug
		restaurant.Street = data.Address.Street
		restaurant.District = data.Address.District
		restaurant.City = data.Address.City
		restaurant.FullAddress = data.Address.Full
		restaurant.Lat = data.Lat
		restaurant.Long = data.Long
		restaurant.IsOpening24h = data.OperatingStatus.IsOpening24h
		restaurant.ExternalId = fmt.Sprint(data.ID)
		return db.Save(&restaurant).Error
	}
	restaurant = model.Restaurant{
		Name:         data.Name,
		Avatar:       data.Avatar,
		Phone:        data.Phone,
		Slug:         data.Slug,
		Street:       data.Address.Street,
		District:     data.Address.District,
		City:         data.Address.City,
		FullAddress:  data.Address.Full,
		Lat:          data.Lat,
		Long:         data.Long,
		IsOpening24h: data.OperatingStatus.IsOpening24h,
		ExternalId:   fmt.Sprint(data.ID),
	}
	return db.Save(&restaurant).Error

}
