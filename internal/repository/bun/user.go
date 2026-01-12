// internal/repository/bun/user.go
package bunrepo

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/uptrace/bun"
)

type UserRepo struct {
	db *bun.DB
}

func NewUserRepo(db *bun.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetOrCreate(ctx context.Context, telegramID int64, username string) (*domain.User, error) {
	user := &domain.User{
		TelegramID: telegramID,
		Username:   username,
		Role:       domain.RoleViewer,
	}

	_, err := r.db.NewInsert().
		Model(user).
		On("CONFLICT (telegram_id) DO UPDATE").
		Set("username = EXCLUDED.username").
		Returning("*").
		Exec(ctx)

	return user, err
}

func (r *UserRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	user := new(domain.User)
	err := r.db.NewSelect().
		Model(user).
		Where("telegram_id = ?", telegramID).
		Scan(ctx)
	return user, err
}

func (r *UserRepo) UpdateRole(ctx context.Context, telegramID int64, role domain.Role) error {
	_, err := r.db.NewUpdate().
		Model((*domain.User)(nil)).
		Set("role = ?", role).
		Where("telegram_id = ?", telegramID).
		Exec(ctx)
	return err
}
