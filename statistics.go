package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"sync"

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
var TransactionStoreMap = make(map[time.Time]TransactionStatistics)
// Mutex for handling concurrent accesss on TransactionStoreMap
var TransactionStoreMapMutex sync.Mutex

// Stores location from where statistics can be accessed.
var StatAllowedLocation StatisticsLocation

// SetLocation updates the location restriction for gettins statistics.
func setLocation(ginCtx *gin.Context) {
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
func resetLocation(ginCtx *gin.Context) {
	StatAllowedLocation = StatisticsLocation{}

	ginCtx.String(http.StatusOK, "Location reset")
}

// LoadTransactionStatistics generates statistics based on the transactions in the last 60 seconds.
func loadTransactionStatistics(ginCtx *gin.Context) {

	// Parse location from the request.
	requestLocation := ginCtx.Query("location")

	if requestLocation == "" {
		log.Println("Location not provided")
		ginCtx.String(http.StatusUnauthorized, "Location not provided")
		return
	}

	if requestLocation != StatAllowedLocation.City {
		log.Println("Unauthorized location")
		ginCtx.String(http.StatusUnauthorized, "Unauthorized location")
		return
	}

	// Iterate the transaction statistics from last 60 seconds.
	endTime := time.Now().UTC().Round(time.Millisecond * 1000)
	startTime := endTime.Add(time.Second * -60)
	log.Printf("start time = %v, end time = %v", startTime, endTime)

	// Stores statistics summary from last 60s
	var statisticsSummary TransactionStatistics

	// Iterate through last 60 seconds.
	// This iteration will always be 60, no matter how many transactions are there in the last 60s.
	for mapKey := startTime; (mapKey.Before(endTime) || mapKey.Equal(endTime)); mapKey = mapKey.Add(time.Second * 1) {

		// Get the statistics for the second.
		TransactionStoreMapMutex.Lock()
		existingStatistics, ok := TransactionStoreMap[mapKey]
		TransactionStoreMapMutex.Unlock()

		if ok {
		log.Printf("transaction time = %v", mapKey)
			// Add statistics on the second to the summary
			statisticsSummary.Sum = statisticsSummary.Sum + existingStatistics.Sum
			statisticsSummary.Count = statisticsSummary.Count + existingStatistics.Count
			statisticsSummary.Average = statisticsSummary.Sum / statisticsSummary.Count

			if statisticsSummary.Minimum == 0 ||  statisticsSummary.Minimum > existingStatistics.Minimum {
				statisticsSummary.Minimum = existingStatistics.Minimum
			} 

			if statisticsSummary.Maximum == 0 ||  statisticsSummary.Maximum < existingStatistics.Maximum {
				statisticsSummary.Maximum = existingStatistics.Maximum
			} 

		}
	}

	// Return the response.
	ginCtx.JSON(http.StatusOK, statisticsSummary)
}

// TransactionStatisticsUpdator udpates transaction statistics for the given second.
func transactionStatisticsUpdator(transactionTime time.Time, transactionAmount float32) {

	log.Printf("transaction time = %s", transactionTime.Format(time.RFC1123))

	TransactionStoreMapMutex.Lock()
	existingStatistics, ok := TransactionStoreMap[transactionTime]
	TransactionStoreMapMutex.Unlock()
	
	if ok {
		// Update existing statistics
		sum := existingStatistics.Sum + transactionAmount
		count := existingStatistics.Count + 1
		min := existingStatistics.Minimum

		// Update minimum 
		if transactionAmount != 0 && transactionAmount < min {
			min = transactionAmount
		}

		max := existingStatistics.Maximum
		// Update minimum 
		if transactionAmount != 0 && transactionAmount < min {
			max = transactionAmount
		}

		existingStatistics = TransactionStatistics{
			Sum: sum,
			Average: sum/count,
			Maximum: max,
			Minimum: min,
			Count: count,
		}

		// Update the global store map.
		TransactionStoreMapMutex.Lock()
		TransactionStoreMap[transactionTime] = existingStatistics
		TransactionStoreMapMutex.Unlock()
	} else {
		newStatistics := TransactionStatistics{
			Sum: transactionAmount,
			Average: transactionAmount,
			Maximum: transactionAmount,
			Minimum: transactionAmount,
			Count: 1,
		}

		// Update the global store map.
		TransactionStoreMapMutex.Lock()
		TransactionStoreMap[transactionTime] = newStatistics
		TransactionStoreMapMutex.Unlock()
	}

	log.Println("TransactionStatistics updated")
	log.Printf("data = %v", TransactionStoreMap)
}
