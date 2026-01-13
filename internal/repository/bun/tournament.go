// internal/repository/bun/tournament.go
package bunrepo

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/uptrace/bun"
)

type TournamentRepo struct {
	db *bun.DB
}

func NewTournamentRepo(db *bun.DB) *TournamentRepo {
	return &TournamentRepo{db: db}
}

func (r *TournamentRepo) Create(ctx context.Context, t *domain.Tournament) error {
	_, err := r.db.NewInsert().Model(t).Returning("*").Exec(ctx)
	return err
}

func (r *TournamentRepo) GetByID(ctx context.Context, id int64) (*domain.Tournament, error) {
	t := new(domain.Tournament)
	err := r.db.NewSelect().Model(t).Where("id = ?", id).Where("deleted_at IS NULL").Scan(ctx)
	return t, err
}

func (r *TournamentRepo) List(ctx context.Context) ([]*domain.Tournament, error) {
	var tournaments []*domain.Tournament
	err := r.db.NewSelect().Model(&tournaments).Where("deleted_at IS NULL").Order("date DESC").Scan(ctx)
	return tournaments, err
}

func (r *TournamentRepo) ListRecent(ctx context.Context, limit int) ([]*domain.Tournament, error) {
	var tournaments []*domain.Tournament
	err := r.db.NewSelect().
		Model(&tournaments).
		Where("deleted_at IS NULL").
		Order("date DESC").
		Limit(limit).
		Scan(ctx)
	return tournaments, err
}

func (r *TournamentRepo) Update(ctx context.Context, t *domain.Tournament) error {
	_, err := r.db.NewUpdate().Model(t).WherePK().Returning("*").Exec(ctx)
	return err
}

func (r *TournamentRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*domain.Tournament)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
