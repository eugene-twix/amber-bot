// internal/config/config.go
package config

import (
	"strconv"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string  `env:"TELEGRAM_TOKEN,required"`
	DatabaseURL   string  `env:"DATABASE_URL,required"`
	RedisURL      string  `env:"REDIS_URL,required"`
	AdminIDsRaw   string  `env:"ADMIN_IDS" envDefault:""`
	AdminIDs      []int64 `env:"-"`

	// API server config
	APIPort      int    `env:"API_PORT" envDefault:"8080"`
	FrontendPath string `env:"FRONTEND_PATH" envDefault:"./frontend/dist"`

	// Mini App URL (for bot button)
	MiniAppURL string `env:"MINI_APP_URL" envDefault:""`

	// Dev mode (bypasses Telegram auth)
	DevMode      bool  `env:"DEV_MODE" envDefault:"false"`
	DevUserID    int64 `env:"DEV_USER_ID" envDefault:"123456789"`
}

func Load() (*Config, error) {
	// Load public env first, then secrets (secrets override)
	_ = godotenv.Load(".env")
	_ = godotenv.Load(".secret.env")

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	cfg.AdminIDs = parseAdminIDs(cfg.AdminIDsRaw)
	return cfg, nil
}

func parseAdminIDs(raw string) []int64 {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	ids := make([]int64, 0, len(parts))
	for _, p := range parts {
		if id, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	return ids
}
