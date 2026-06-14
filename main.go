// @title           IMS - Institution Management Service API
// @version         1.0
// @description     API for managing institution member information.
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"log"

	_ "github.com/vickyaruldoss/ims/docs"
	"github.com/vickyaruldoss/ims/router"
)

func main() {
	r := router.SetupRouter()
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
