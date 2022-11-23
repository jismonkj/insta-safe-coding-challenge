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
		router.POST("/transactions", StoreTransaction)
		router.DELETE("/transactions", RemoveTransactions)
	}

	// Manage user location
	{
		router.POST("/location", SetLocation)
		router.DELETE("/location", ResetLocation)
	}

	return router
}

func main() {

	ginEngine := startGinRouter()
	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	ginEngine.Run()
}
