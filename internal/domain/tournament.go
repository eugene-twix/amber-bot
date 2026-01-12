// internal/domain/tournament.go
package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Tournament struct {
	bun.BaseModel `bun:"table:tournaments"`

	ID        int64     `bun:"id,pk,autoincrement"`
	Name      string    `bun:"name,notnull"`
	Date      time.Time `bun:"date,notnull"`
	Location  string    `bun:"location"`
	CreatedBy int64     `bun:"created_by"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
}
