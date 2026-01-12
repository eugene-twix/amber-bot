// internal/repository/bun/result.go
package bunrepo

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/eugene-twix/amber-bot/internal/repository"
	"github.com/uptrace/bun"
)

type ResultRepo struct {
	db *bun.DB
}

func NewResultRepo(db *bun.DB) *ResultRepo {
	return &ResultRepo{db: db}
}

func (r *ResultRepo) Create(ctx context.Context, res *domain.Result) error {
	_, err := r.db.NewInsert().Model(res).Returning("*").Exec(ctx)
	return err
}

func (r *ResultRepo) GetByTeamID(ctx context.Context, teamID int64) ([]*domain.Result, error) {
	var results []*domain.Result
	err := r.db.NewSelect().
		Model(&results).
		Where("team_id = ?", teamID).
		Order("recorded_at DESC").
		Scan(ctx)
	return results, err
}

func (r *ResultRepo) GetByTournamentID(ctx context.Context, tournamentID int64) ([]*domain.Result, error) {
	var results []*domain.Result
	err := r.db.NewSelect().
		Model(&results).
		Where("tournament_id = ?", tournamentID).
		Order("place ASC").
		Scan(ctx)
	return results, err
}

func (r *ResultRepo) GetTeamRating(ctx context.Context) ([]repository.TeamRating, error) {
	var ratings []repository.TeamRating
	err := r.db.NewRaw(`
		SELECT
			t.id as team_id,
			t.name as team_name,
			COUNT(CASE WHEN r.place = 1 THEN 1 END) as wins,
			COUNT(r.id) as total_games,
			COALESCE(AVG(r.place), 0) as avg_place
		FROM teams t
		LEFT JOIN results r ON t.id = r.team_id
		GROUP BY t.id, t.name
		HAVING COUNT(r.id) > 0
		ORDER BY wins DESC, avg_place ASC
	`).Scan(ctx, &ratings)
	return ratings, err
}
