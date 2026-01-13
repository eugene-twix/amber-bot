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
	_, err := r.db.NewInsert().
		Model(res).
		On("CONFLICT (tournament_id, team_id) DO UPDATE").
		Set("place = EXCLUDED.place").
		Set("recorded_by = EXCLUDED.recorded_by").
		Set("recorded_at = CURRENT_TIMESTAMP").
		Returning("*").
		Exec(ctx)
	return err
}

func (r *ResultRepo) GetByTeamID(ctx context.Context, teamID int64) ([]*domain.Result, error) {
	var results []*domain.Result
	err := r.db.NewSelect().
		Model(&results).
		Relation("Tournament").
		Where("result.team_id = ?", teamID).
		Where("result.deleted_at IS NULL").
		Order("recorded_at DESC").
		Scan(ctx)
	return results, err
}

func (r *ResultRepo) GetByTournamentID(ctx context.Context, tournamentID int64) ([]*domain.Result, error) {
	var results []*domain.Result
	err := r.db.NewSelect().
		Model(&results).
		Relation("Team").
		Where("result.tournament_id = ?", tournamentID).
		Where("result.deleted_at IS NULL").
		Order("place ASC").
		Scan(ctx)
	return results, err
}

func (r *ResultRepo) GetByID(ctx context.Context, id int64) (*domain.Result, error) {
	result := new(domain.Result)
	err := r.db.NewSelect().Model(result).Where("id = ?", id).Where("deleted_at IS NULL").Scan(ctx)
	return result, err
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
		LEFT JOIN results r ON t.id = r.team_id AND r.deleted_at IS NULL
		WHERE t.deleted_at IS NULL
		GROUP BY t.id, t.name
		HAVING COUNT(r.id) > 0
		ORDER BY wins DESC, avg_place ASC
	`).Scan(ctx, &ratings)
	return ratings, err
}

func (r *ResultRepo) Update(ctx context.Context, res *domain.Result) error {
	_, err := r.db.NewUpdate().Model(res).WherePK().Returning("*").Exec(ctx)
	return err
}

func (r *ResultRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*domain.Result)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *ResultRepo) DeleteWithShift(ctx context.Context, id int64, tournamentID int64, deletedPlace int) error {
	return r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		// Delete the result
		_, err := tx.NewDelete().Model((*domain.Result)(nil)).Where("id = ?", id).Exec(ctx)
		if err != nil {
			return err
		}

		// Safety check: only shift if deletedPlace >= 1
		if deletedPlace < 1 {
			return nil
		}

		// Shift places down to fill the gap
		_, err = tx.NewUpdate().
			Model((*domain.Result)(nil)).
			Set("place = place - 1").
			Where("tournament_id = ?", tournamentID).
			Where("place > ?", deletedPlace).
			Where("place > 1"). // Ensure we never go below 1
			Exec(ctx)
		return err
	})
}
