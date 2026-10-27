package main

import (
	"log"

	"motoflow-backend/config"
	"motoflow-backend/routes"
	"motoflow-backend/seeders"
)

func main() {
	cfg := config.LoadConfig()

	log.Printf("[MotoFlow AHASS] Starting backend service on port %s (env: %s)...", cfg.Port, cfg.Environment)

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("[MotoFlow FATAL] Database connection failed: %v", err)
	}

	// Run realistic seed data if database is fresh
	if err := seeders.SeedAll(db); err != nil {
		log.Printf("[MotoFlow WARNING] Seeder encountered issue: %v", err)
	}
	seeders.SeedRetailAndShowroom(db)

	r := routes.SetupRouter(db)

	address := ":" + cfg.Port
	log.Printf("[MotoFlow READY] Listening at http://localhost%s", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("[MotoFlow FATAL] Failed to start HTTP server: %v", err)
	}
}

