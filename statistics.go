package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// StatLocation struct for unmarshalling location input.
type StatisticsLocation struct {
	City string `json:"city"`
}

// TransactionStatistics struct for holding transaction statistics.
type TransactionStatistics struct {
	// “sum”:””,
	// “avg”:””,
	// “max”:””,
	// “min”:””,
	// “count”:””
	Sum     float32 `json:"sum"` // Total transaction amount
	Average float32 `json:"avg"`
	Maximum float32 `json:"max"`
	Minimum float32 `json:"min"`
	Count   float32 `json:"count"` // Number of transactions
}

// TransactionStoreMap stores transaction statistics on each second.
// Store the transaction stat in the map against transaction time.
// This for getting transaction stat. summary in the last minute with O(1) time complexity when calling the statistics API.
var TransactionStoreMap = make(map[time.Time][]TransactionStatistics)
//
TransactionStoreMap

// Stores location from where statistics can be accessed.
var StatAllowedLocation StatisticsLocation

// SetLocation updates the location restriction for gettins statistics.
func SetLocation(ginCtx *gin.Context) {
	// Stores location input from API request.
	var location StatisticsLocation

	reqBody, err := ginCtx.GetRawData()
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to parse request")
		return
	}
	log.Printf("location = %v\n", string(reqBody))

	err = json.Unmarshal(reqBody, &location)
	if err != nil {
		log.Printf("err = %v\n", err)
		ginCtx.String(http.StatusBadRequest, "Failed to parse request")
		return
	}

	// Update the location.
	if location.City != "" {
		StatAllowedLocation = location
	}

	ginCtx.String(http.StatusOK, "Location updated")
}

// ResetLocation resets the location restriction for getting statistics.
func ResetLocation(ginCtx *gin.Context) {
	StatAllowedLocation = StatisticsLocation{}

	ginCtx.String(http.StatusOK, "Location reset")
}

// LoadTransactionStatistics generates statistics based on the transactions in the last 60 seconds.
func LoadTransactionStatistics(ginCtx *gin.Context) {

	// Parse location from the request.
	requestLocation := ginCtx.Query("location")

	if requestLocation == "" {
		log.Println("Location not provided")
		ginCtx.String(http.StatusUnauthorized, "Location not provided")
	}

	if requestLocation != StatAllowedLocation.City {
		log.Println("Unauthorized location")
		ginCtx.String(http.StatusUnauthorized, "Unauthorized location")
	}

}

// TransactionStatisticsUpdator udpates transaction statistics for the given second.
func TransactionStatisticsUpdator(transactionTime *time.Time, transactionAmount float32) {

	log.Printf("transaction time = %s", transactionTime.Format(time.RFC1123))

	if _, ok := TransactionStoreMap[transactionAmount]; ok {

		// Update existing statistics
		existingStatis
	}
}
