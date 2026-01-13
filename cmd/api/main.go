// cmd/api/main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eugene-twix/amber-bot/internal/api"
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
	defer c.Close()

	// Initialize repositories
	repos := &api.Repositories{
		User:       bunrepo.NewUserRepo(db),
		Team:       bunrepo.NewTeamRepo(db),
		Member:     bunrepo.NewMemberRepo(db),
		Tournament: bunrepo.NewTournamentRepo(db),
		Result:     bunrepo.NewResultRepo(db),
	}

	// Create API server
	server := api.NewServer(api.Config{
		Port:         cfg.APIPort,
		BotToken:     cfg.TelegramToken,
		FrontendPath: cfg.FrontendPath,
		DevMode:      cfg.DevMode,
		DevUserID:    cfg.DevUserID,
	}, repos, c)

	if cfg.DevMode {
		log.Printf("WARNING: Running in DEV MODE with user ID %d", cfg.DevUserID)
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down API server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
	}()

	log.Printf("Starting API server on port %d", cfg.APIPort)
	if err := server.Run(); err != nil {
		log.Fatalf("API server error: %v", err)
	}
}
