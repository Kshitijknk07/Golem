package main

import (
	"log"

	"golem/internal/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.SetupRoutes(r)

	if err := r.Run(":3000"); err != nil {
		log.Fatalf("could not start the server: %v", err)
	}
}
