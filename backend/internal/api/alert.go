package api

import (
	"golem/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterAlertRoutes(r *gin.Engine) {
	r.GET("/alerts", getAlerts)
}

func getAlerts(c *gin.Context) {
	alerts := services.GetAlerts()
	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}
