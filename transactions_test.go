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

// TestStoreTransaction tests transaction creation functonality.
func TestStoreTransaction(t *testing.T) {

	// Starts the API server for sending requests.
	ginEngine := startGinRouter()

	// List of transactions with different expected outputs.
	currentTime := time.Now().UTC()
	var transactionCases []Transaction

	// 422 : Transaction timestamp is in future.
	transactionCases = append(transactionCases, Transaction{Amount: "100.25", Timestamp: currentTime.Add(time.Second * time.Duration(60)).Format(time.RFC3339), ExpectedStatus: 422})

	// 400 : Invalid Json.
	transactionCases = append(transactionCases, Transaction{ExpectedStatus: 400})

	// 204 : Transaction is older than 60s
	transactionCases = append(transactionCases, Transaction{Amount: "100.25", Timestamp: currentTime.Add(time.Second * time.Duration(-70)).Format(time.RFC3339), ExpectedStatus: 204})

	// 201 : Transaction created
	transactionCases = append(transactionCases, Transaction{Amount: "100.25", Timestamp: currentTime.Format(time.RFC3339), ExpectedStatus: http.StatusCreated})

	// Iterate through each Transaction cases and run sub tests.
	for _, transaction := range transactionCases {
		t.Run(fmt.Sprintf("Transaction Expected Status: %d", transaction.ExpectedStatus), func(t *testing.T) {

			var transactionJson []byte
			var err error

			if transaction.Amount != "" {
				transactionJson, err = json.Marshal(transaction)
				if err != nil {
					t.Errorf("Failed to parse transaction data")
				}
			}

			bodyReader := bytes.NewReader(transactionJson)

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/transactions", bodyReader)
			ginEngine.ServeHTTP(w, req)

			if transaction.ExpectedStatus != w.Code {
				t.Errorf("status %d", w.Code)
			}
		})
	}
}

// TestRemoveTransactions tests functionality for removing all transactions.
func TestRemoveTransactions(t *testing.T) {
	ginEngine := startGinRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodDelete, "/transactions", nil)
	ginEngine.ServeHTTP(w, req)

	if err != nil {
		t.Errorf("err = %v", err)
	}

	if http.StatusNoContent != w.Code {
		t.Errorf("status %d", w.Code)
	}
}
