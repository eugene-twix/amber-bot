// internal/domain/team.go
package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Team struct {
	bun.BaseModel `bun:"table:teams"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,unique,notnull"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
	CreatedBy int64     `bun:"created_by"`

	// Metadata for updates
	UpdatedAt *time.Time `bun:"updated_at"`
	UpdatedBy *int64     `bun:"updated_by"`

	// Soft delete
	DeletedAt *time.Time `bun:"deleted_at,soft_delete"`
	DeletedBy *int64     `bun:"deleted_by"`

	// Optimistic locking
	Version int `bun:"version,default:1"`
}
