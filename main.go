package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"example.com/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	threadMax   = 10
	totalPages  = 200
	apiEndpoint = "https://mocha.lozi.vn/v6.1/search/eateries/near-by?cityId=50&limit=50&superCategoryId=1&lat=10.7765194&lng=106.700987&page="
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

func fetchData(page int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Construct the API URL for the given page
	apiURL := apiEndpoint + fmt.Sprint(page)

	// Make the HTTP GET request
	response, err := http.Get(apiURL)
	if err != nil {
		fmt.Printf("Error making the request for page %d: %v\n", page, err)
		return
	}
	defer response.Body.Close()

	// Decode the JSON response
	var apiResponse ApiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		fmt.Printf("Error decoding JSON for page %d: %v\n", page, err)
		return
	}

	// Process the results (you can customize this part based on your needs)
	fmt.Printf("Data for page %d:\n", page)
	for _, data := range apiResponse.Data {
		fmt.Printf("  ID: %d\n", data.ID)
		err = saveToDB(data)
		if err != nil {
			fmt.Printf("Error saving to DB for page %d: %v\n", page, err)
			return
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

	var wg sync.WaitGroup

	for page := 1; page <= totalPages; page++ {
		// Limit the number of concurrent threads
		if page%threadMax == 0 {
			wg.Wait()
		}

		wg.Add(1)
		go fetchData(page, &wg)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

func saveToDB(data ApiResponseData) error {
	restaurant := model.Restaurant{
		Name:                   data.Name,
		Avatar:                 data.Avatar,
		Phone:                  data.Phone,
		CountryCode:            data.CountryCode,
		Slug:                   data.Slug,
		Street:                 data.Address.Street,
		District:               data.Address.District,
		City:                   data.Address.City,
		FullAddress:            data.Address.Full,
		Rating:                 data.Rating,
		Username:               data.Username,
		Lat:                    data.Lat,
		Long:                   data.Long,
		IsOpening:              data.OperatingStatus.IsOpening,
		IsOpening24h:           data.OperatingStatus.IsOpening24h,
		MinutesUntilNextStatus: data.OperatingStatus.MinutesUntilNextStatus,
		IsLoshipPartner:        data.IsLoshipPartner,
		IsHonored:              data.IsHonored,
		Quote:                  data.Quote,
		IsActive:               data.IsActive,
		IsCheckedIn:            data.IsCheckedIn,
		Closed:                 data.Closed,
		RecommendedRatio:       data.RecommendedRatio,
		RecommendedEnable:      data.RecommendedEnable,
		Distance:               data.Distance,
		IsPurchasedSupplyItems: data.IsPurchasedSupplyItems,
		IsSponsored:            data.IsSponsored,
		FreeShippingMilestone:  data.FreeShippingMilestone,
	}

	return db.Create(&restaurant).Error
}
