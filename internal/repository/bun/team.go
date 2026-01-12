// internal/repository/bun/team.go
package bunrepo

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/uptrace/bun"
)

type TeamRepo struct {
	db *bun.DB
}

func NewTeamRepo(db *bun.DB) *TeamRepo {
	return &TeamRepo{db: db}
}

func (r *TeamRepo) Create(ctx context.Context, team *domain.Team) error {
	_, err := r.db.NewInsert().Model(team).Returning("*").Exec(ctx)
	return err
}

func (r *TeamRepo) GetByID(ctx context.Context, id int64) (*domain.Team, error) {
	team := new(domain.Team)
	err := r.db.NewSelect().Model(team).Where("id = ?", id).Scan(ctx)
	return team, err
}

func (r *TeamRepo) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	team := new(domain.Team)
	err := r.db.NewSelect().Model(team).Where("name = ?", name).Scan(ctx)
	return team, err
}

func (r *TeamRepo) List(ctx context.Context) ([]*domain.Team, error) {
	var teams []*domain.Team
	err := r.db.NewSelect().Model(&teams).Order("name ASC").Scan(ctx)
	return teams, err
}

func (r *TeamRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*domain.Team)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
