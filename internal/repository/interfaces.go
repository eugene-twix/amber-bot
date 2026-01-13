// internal/repository/interfaces.go
package repository

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
)

type UserRepository interface {
	GetOrCreate(ctx context.Context, telegramID int64, username string) (*domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	UpdateRole(ctx context.Context, telegramID int64, role domain.Role) error
	List(ctx context.Context) ([]*domain.User, error)
}

type TeamRepository interface {
	Create(ctx context.Context, team *domain.Team) error
	GetByID(ctx context.Context, id int64) (*domain.Team, error)
	GetByName(ctx context.Context, name string) (*domain.Team, error)
	List(ctx context.Context) ([]*domain.Team, error)
	Update(ctx context.Context, team *domain.Team) error
	Delete(ctx context.Context, id int64) error
}

type MemberRepository interface {
	Create(ctx context.Context, member *domain.Member) error
	GetByID(ctx context.Context, id int64) (*domain.Member, error)
	GetByTeamID(ctx context.Context, teamID int64) ([]*domain.Member, error)
	Update(ctx context.Context, member *domain.Member) error
	Delete(ctx context.Context, id int64) error
}

type TournamentRepository interface {
	Create(ctx context.Context, tournament *domain.Tournament) error
	GetByID(ctx context.Context, id int64) (*domain.Tournament, error)
	List(ctx context.Context) ([]*domain.Tournament, error)
	ListRecent(ctx context.Context, limit int) ([]*domain.Tournament, error)
	Update(ctx context.Context, tournament *domain.Tournament) error
	Delete(ctx context.Context, id int64) error
}

type ResultRepository interface {
	Create(ctx context.Context, result *domain.Result) error
	GetByID(ctx context.Context, id int64) (*domain.Result, error)
	GetByTeamID(ctx context.Context, teamID int64) ([]*domain.Result, error)
	GetByTournamentID(ctx context.Context, tournamentID int64) ([]*domain.Result, error)
	GetTeamRating(ctx context.Context) ([]TeamRating, error)
	Update(ctx context.Context, result *domain.Result) error
	Delete(ctx context.Context, id int64) error
	// DeleteWithShift deletes result and shifts higher places down in a transaction
	DeleteWithShift(ctx context.Context, id int64, tournamentID int64, deletedPlace int) error
}

type TeamRating struct {
	TeamID     int64
	TeamName   string
	Wins       int // 1 места
	TotalGames int
	AvgPlace   float64
}
