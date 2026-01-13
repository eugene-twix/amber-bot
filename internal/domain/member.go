// internal/domain/member.go
package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Member struct {
	bun.BaseModel `bun:"table:members"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	TeamID    int64     `bun:"team_id,notnull"`
	JoinedAt  time.Time `bun:"joined_at,default:current_timestamp"`
	CreatedBy int64     `bun:"created_by"`

	// Metadata for updates
	UpdatedAt *time.Time `bun:"updated_at"`
	UpdatedBy *int64     `bun:"updated_by"`

	// Soft delete
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	DeletedBy *int64     `bun:"deleted_by"`

	// Optimistic locking
	Version int `bun:"version,default:1"`

	// Relations
	Team *Team `bun:"rel:belongs-to,join:team_id=id"`
}
