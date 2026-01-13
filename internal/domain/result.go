// internal/domain/result.go
package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Result struct {
	bun.BaseModel `bun:"table:results"`

	ID           int64     `bun:"id,pk,autoincrement"`
	TeamID       int64     `bun:"team_id,notnull"`
	TournamentID int64     `bun:"tournament_id,notnull"`
	Place        int       `bun:"place,notnull"`
	RecordedBy   int64     `bun:"recorded_by"`
	RecordedAt   time.Time `bun:"recorded_at,default:current_timestamp"`

	// Metadata for updates
	UpdatedAt *time.Time `bun:"updated_at"`
	UpdatedBy *int64     `bun:"updated_by"`

	// Soft delete
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	DeletedBy *int64     `bun:"deleted_by"`

	// Optimistic locking
	Version int `bun:"version,default:1"`

	// Relations
	Team       *Team       `bun:"rel:belongs-to,join:team_id=id"`
	Tournament *Tournament `bun:"rel:belongs-to,join:tournament_id=id"`
}
