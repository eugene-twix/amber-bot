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
}
