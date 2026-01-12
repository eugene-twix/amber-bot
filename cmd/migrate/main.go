// cmd/migrate/main.go
package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/eugene-twix/amber-bot/internal/config"
	"github.com/eugene-twix/amber-bot/migrations"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DatabaseURL)))
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	migrator := migrate.NewMigrator(db, migrations.Migrations)

	ctx := context.Background()

	if err := migrator.Init(ctx); err != nil {
		log.Fatalf("Failed to init migrator: %v", err)
	}

	cmd := "up"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	switch cmd {
	case "up":
		group, err := migrator.Migrate(ctx)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		if group.IsZero() {
			log.Println("No new migrations to run")
		} else {
			log.Printf("Migrated to %s\n", group)
		}
	case "down":
		group, err := migrator.Rollback(ctx)
		if err != nil {
			log.Fatalf("Rollback failed: %v", err)
		}
		if group.IsZero() {
			log.Println("Nothing to rollback")
		} else {
			log.Printf("Rolled back %s\n", group)
		}
	case "status":
		ms, err := migrator.MigrationsWithStatus(ctx)
		if err != nil {
			log.Fatalf("Failed to get status: %v", err)
		}
		log.Println("Migrations:")
		for _, m := range ms {
			log.Printf("  %s: %s\n", m.Name, m.String())
		}
	default:
		log.Fatalf("Unknown command: %s (use: up, down, status)", cmd)
	}
}
