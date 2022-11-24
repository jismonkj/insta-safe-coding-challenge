package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Transaction struct used to unmarshall transaction details recieved through API request
type Transaction struct {
	Amount         string `json:"amount"`
	Timestamp      string `json:"timestamp"` // Expected format: ISO 8601 format YYYY-MM-DDThh:mm:ss.sssZ in the UTC timezone (RFC3339)
	ExpectedStatus int    // Used for testing (Expected response status on different transactions)
}

// storeTransaction stores a valid transaction in the global map.
//
// Transaction (example):
// {
//		“amount”:”100.25”
//		“timestamp”:"2021-07-17T09:59:51.312Z"
// }
//
// Responses:
// 201 – in case of success (http.StatusCreated)
// 204 – if the transaction is older than 60 seconds (http.StatusNoContent)
// 400 – if the JSON is invalid	(http.StatusBadRequest)
// 422 – if any of the fields are not parsable or the transaction date is in the future (http.StatusUnprocessableEntity)
//
func storeTransaction(ginCtx *gin.Context) {
	reqBody, err := ginCtx.GetRawData()
	if err != nil {
		ginCtx.String(http.StatusInternalServerError, "Failed to parse request")
		return
	}
	log.Printf("transaction = %v\n", string(reqBody))

	// Unmarshall the request body
	var transaction Transaction

	err = json.Unmarshal(reqBody, &transaction)
	if err != nil {
		log.Printf("err = %v\n", err)
		ginCtx.String(http.StatusBadRequest, "Failed to parse request")
		return
	}

	// Verify the transaction timestamp.
	transactionTime, err := time.Parse(time.RFC3339, transaction.Timestamp)
	if err != nil {
		ginCtx.String(http.StatusUnprocessableEntity, "Failed to parse request")
		return
	}
	log.Printf("transactionTime = %v\n", transactionTime)

	// Validate the transaction timestamp.
	// # A valid timestamp should not be a future time.
	if transactionTime.After(time.Now().UTC()) {
		ginCtx.String(http.StatusUnprocessableEntity, "Time is in the future")
		return
	}

	// # A valid timestamp should not be older than 60 seconds and not a future time.
	transactionEndTime := time.Now().UTC().Add(time.Second * time.Duration(-60))
	if transactionTime.Before(transactionEndTime) {
		ginCtx.String(http.StatusNoContent, "Time is older than 60s")
		return
	}

	// Convert amount to float32
	value, err := strconv.ParseFloat(transaction.Amount, 32)
	if err != nil {
		ginCtx.String(http.StatusUnprocessableEntity, "Failed to parse amount")
		return
	}
	amount := float32(value)

	// All validations are complete!!.
	log.Println("Transaction is valid")

	// Update the statistics.
	transactionStatisticsUpdator(transactionTime, amount)
	ginCtx.String(http.StatusCreated, "Transaction updated")
}

// RemoveTransactions removes stored transactions from the global map.
func removeTransactions(ginCtx *gin.Context) {
	TransactionStoreMap = make(map[time.Time]TransactionStatistics)
	log.Println("Cleared transaction map")
	ginCtx.String(http.StatusNoContent, "Transactions cleared")
}
