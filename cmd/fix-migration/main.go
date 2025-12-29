package main

import (
	"flag"
	"log"
	"prodyo-backend/cmd/internal/config"
	"prodyo-backend/cmd/internal/migrations"
)

func main() {
	forceVersion := flag.Int("force", 0, "Force migration to specific version")
	flag.Parse()

	cfg := config.Load()

	if *forceVersion > 0 {
		log.Printf("Forcing migration version to %d...\n", *forceVersion)
		if err := migrations.ForceVersion(cfg.DSN(), "cmd/migrations", *forceVersion); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		log.Println("Migration version forced successfully")
	} else {
		// Check current migration status
		version, dirty, err := migrations.CheckMigrationStatus(cfg.DSN(), "cmd/migrations")
		if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}
		log.Printf("Current migration version: %d, Dirty: %v\n", version, dirty)
		log.Println("Use -force=<version> to force a specific version")
	}
}
