// internal/migrations/migrations.go
package migrations

import (
	"context"
	"embed"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

//go:embed *.sql
var sqlMigrations embed.FS

func newMigrator(db *bun.DB) (*migrate.Migrator, error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(sqlMigrations); err != nil {
		return nil, fmt.Errorf("discover migrations: %w", err)
	}
	return migrate.NewMigrator(db, migrations), nil
}

// Up applies all pending migrations.
func Up(ctx context.Context, db *bun.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return err
	}

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	if !group.IsZero() {
		fmt.Printf("Migrated to %s\n", group)
	}
	return nil
}

// Down rolls back the last migration group.
func Down(ctx context.Context, db *bun.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return err
	}

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}

	group, err := migrator.Rollback(ctx)
	if err != nil {
		return fmt.Errorf("rollback: %w", err)
	}

	if group.IsZero() {
		fmt.Println("Nothing to rollback")
	} else {
		fmt.Printf("Rolled back %s\n", group)
	}
	return nil
}

// Status prints migration status.
func Status(ctx context.Context, db *bun.DB) error {
	migrator, err := newMigrator(db)
	if err != nil {
		return err
	}

	if err := migrator.Init(ctx); err != nil {
		return fmt.Errorf("init migrator: %w", err)
	}

	ms, err := migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return fmt.Errorf("get status: %w", err)
	}

	fmt.Println("Migrations:")
	for _, m := range ms {
		fmt.Printf("  %s: %s\n", m.Name, m.String())
	}
	return nil
}
