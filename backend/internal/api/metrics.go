package api

import (
	"golem/internal/models"
	"golem/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterMetricsRoutes(r *gin.Engine) {
	r.GET("/metrics", getMetrics)
	r.POST("/metrics", postMetrics)
	r.PUT("/metrics/:id", updateMetric)
}

func getMetrics(c *gin.Context) {
	metrics, err := services.GetMetrics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve metrics"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"metrics": metrics})
}

func postMetrics(c *gin.Context) {
	var newMetric models.Metric
	if err := c.BindJSON(&newMetric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	err := services.SaveMetric(newMetric)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metric saved successfully"})
}

func updateMetric(c *gin.Context) {
	id := c.Param("id")
	var updatedMetric models.Metric
	if err := c.BindJSON(&updatedMetric); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	err := services.UpdateMetric(id, updatedMetric)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update metric"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Metric updated successfully"})
}
