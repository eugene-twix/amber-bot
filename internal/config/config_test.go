// internal/config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("TELEGRAM_TOKEN", "test_token")
	os.Setenv("DATABASE_URL", "postgres://test")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("ADMIN_IDS", "123,456")
	defer func() {
		os.Unsetenv("TELEGRAM_TOKEN")
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("ADMIN_IDS")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.TelegramToken != "test_token" {
		t.Errorf("expected test_token, got %s", cfg.TelegramToken)
	}
	if len(cfg.AdminIDs) != 2 {
		t.Errorf("expected 2 admin IDs, got %d", len(cfg.AdminIDs))
	}
}
