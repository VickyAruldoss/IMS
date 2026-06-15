// @title           IMS - Institution Management Service API
// @version         1.0
// @description     API for managing institution member information.
// @host            localhost:8080
// @BasePath        /api/v1
package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/vickyaruldoss/ims/config"
	_ "github.com/vickyaruldoss/ims/docs"
	"github.com/vickyaruldoss/ims/router"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("connected to PostgreSQL")

	r := router.SetupRouter(db)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
