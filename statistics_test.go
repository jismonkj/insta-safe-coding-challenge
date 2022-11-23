package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestSetLocation tests updating statistics access location.
func TestSetLocation(t *testing.T) {

	// Starts the API server for sending requests.
	ginEngine := startGinRouter()

	var location StatisticsLocation = StatisticsLocation{
		City : "bangalore"
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
