// internal/repository/bun/member.go
package bunrepo

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/uptrace/bun"
)

type MemberRepo struct {
	db *bun.DB
}

func NewMemberRepo(db *bun.DB) *MemberRepo {
	return &MemberRepo{db: db}
}

func (r *MemberRepo) Create(ctx context.Context, member *domain.Member) error {
	_, err := r.db.NewInsert().Model(member).Returning("*").Exec(ctx)
	return err
}

func (r *MemberRepo) GetByTeamID(ctx context.Context, teamID int64) ([]*domain.Member, error) {
	var members []*domain.Member
	err := r.db.NewSelect().
		Model(&members).
		Where("team_id = ?", teamID).
		Order("name ASC").
		Scan(ctx)
	return members, err
}

func (r *MemberRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.NewDelete().Model((*domain.Member)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}
