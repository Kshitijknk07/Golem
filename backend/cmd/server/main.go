package main

import (
	"golem/internal/api"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api.RegisterMetricsRoutes(r)
	api.RegisterAlertRoutes(r)

	r.Use(api.JWTMiddleware())

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
