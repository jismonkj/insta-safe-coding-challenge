package main

import (
	"github.com/gin-gonic/gin"
)

// startGinRouter setups gin router.
func startGinRouter() *gin.Engine {
	router := gin.Default()

	// API endpoints
	// Manages transactions
	{
		router.POST("/transactions", storeTransaction)
		router.DELETE("/transactions", removeTransactions)
	}

	// Manage user location
	{
		router.POST("/location", setLocation)
		router.DELETE("/location", resetLocation)
	}

	// Getting statistics
	{
		router.GET("/statistics", loadTransactionStatistics)
	}

	return router
}

func main() {

	ginEngine := startGinRouter()
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	ginEngine.Run()
}
