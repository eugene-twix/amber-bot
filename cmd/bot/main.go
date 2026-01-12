// cmd/bot/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eugene-twix/amber-bot/internal/bot"
	"github.com/eugene-twix/amber-bot/internal/cache"
	"github.com/eugene-twix/amber-bot/internal/config"
	"github.com/eugene-twix/amber-bot/internal/migrations"
	bunrepo "github.com/eugene-twix/amber-bot/internal/repository/bun"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db := bunrepo.NewDB(cfg.DatabaseURL, false)
	defer db.Close()

	// Auto-migrate on startup
	if err := migrations.Up(context.Background(), db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	c, err := cache.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("Failed to connect to cache: %v", err)
	}

	b, err := bot.New(cfg, db, c)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down...")
		b.Stop()
	}()

	b.Start()
}
