// migrations/20260112000001_init.go
package migrations

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		// Users table
		_, err := db.ExecContext(ctx, `
			CREATE TABLE users (
				telegram_id BIGINT PRIMARY KEY,
				username VARCHAR(255),
				role VARCHAR(20) NOT NULL DEFAULT 'viewer',
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			)
		`)
		if err != nil {
			return err
		}

		// Teams table
		_, err = db.ExecContext(ctx, `
			CREATE TABLE teams (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL UNIQUE,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				created_by BIGINT REFERENCES users(telegram_id)
			)
		`)
		if err != nil {
			return err
		}

		// Members table
		_, err = db.ExecContext(ctx, `
			CREATE TABLE members (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
				joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				UNIQUE(name, team_id)
			)
		`)
		if err != nil {
			return err
		}

		// Tournaments table
		_, err = db.ExecContext(ctx, `
			CREATE TABLE tournaments (
				id BIGSERIAL PRIMARY KEY,
				name VARCHAR(255) NOT NULL,
				date DATE NOT NULL,
				location VARCHAR(255),
				created_by BIGINT REFERENCES users(telegram_id),
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			)
		`)
		if err != nil {
			return err
		}

		// Results table
		_, err = db.ExecContext(ctx, `
			CREATE TABLE results (
				id BIGSERIAL PRIMARY KEY,
				team_id BIGINT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
				tournament_id BIGINT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,
				place INT NOT NULL CHECK (place > 0),
				recorded_by BIGINT REFERENCES users(telegram_id),
				recorded_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				UNIQUE(team_id, tournament_id)
			)
		`)
		if err != nil {
			return err
		}

		// Indexes
		_, err = db.ExecContext(ctx, `
			CREATE INDEX idx_results_tournament ON results(tournament_id);
			CREATE INDEX idx_results_team ON results(team_id);
			CREATE INDEX idx_members_team ON members(team_id)
		`)
		return err
	}, func(ctx context.Context, db *bun.DB) error {
		// Down migration
		_, err := db.ExecContext(ctx, `
			DROP TABLE IF EXISTS results;
			DROP TABLE IF EXISTS members;
			DROP TABLE IF EXISTS tournaments;
			DROP TABLE IF EXISTS teams;
			DROP TABLE IF EXISTS users
		`)
		return err
	})
}
