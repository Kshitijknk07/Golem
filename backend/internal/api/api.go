package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up the routes for the API
func SetupRoutes(r *gin.Engine) {
	r.GET("/metrics", getMetrics) // Endpoint to fetch metrics
}

// getMetrics handles the /metrics endpoint and returns mock metrics
func getMetrics(c *gin.Context) {
	// Example metrics (can be replaced with actual data collection)
	metrics := map[string]interface{}{
		"cpu_usage": 75.5,
		"mem_usage": 60.2,
	}

	// Return the metrics in a JSON response
	c.JSON(http.StatusOK, gin.H{
		"metrics": metrics,
	})
}
