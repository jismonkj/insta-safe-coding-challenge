package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// For testing non-exported functions
var SetLocation = setLocation
var ResetLocation = resetLocation
var LoadTransactionStatistics = loadTransactionStatistics


// TestSetLocation tests updating statistics access location.
func TestSetLocation(t *testing.T) {

	// Starts the API server for sending requests.
	ginEngine := startGinRouter()

	var location StatisticsLocation = StatisticsLocation{
		City : "bangalore",
	}

	locationJson, err := json.Marshal(location)
	if err != nil {
		t.Errorf("Failed to parse transaction data")
	}

	bodyReader := bytes.NewReader(locationJson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/location", bodyReader)
	ginEngine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Errorf("status %d", w.Code)
	}
}

// TestResetLocation tests updating statistics access location.
func TestResetLocation(t *testing.T) {

	// Starts the API server for sending requests.
	ginEngine := startGinRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/location", nil)
	ginEngine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Errorf("status %d", w.Code)
	}
}

// TestLoadTransactionStatistics tests getting statstics.
func TestLoadTransactionStatistics(t *testing.T) {
	// Starts the API server for sending requests.
	ginEngine := startGinRouter()

	// Create list of transactions for getting statistics.
	createTransactions(t, ginEngine)

	// Set location access on statistics
	var location StatisticsLocation = StatisticsLocation{
		City : "bangalore",
	}

	locationJson, err := json.Marshal(location)
	if err != nil {
		t.Errorf("Failed to parse transaction data")
	}

	bodyReader := bytes.NewReader(locationJson)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/location", bodyReader)
	ginEngine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Errorf("status %d", w.Code)
	}

	req, _ = http.NewRequest(http.MethodGet, "/statistics?location=bangalore", nil)
	ginEngine.ServeHTTP(w, req)

	if http.StatusOK != w.Code {
		t.Errorf("status %d", w.Code)
	}

	fmt.Printf("client: response body: %s\n", w.Body.String())
}
