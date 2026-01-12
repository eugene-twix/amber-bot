// internal/domain/member.go
package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Member struct {
	bun.BaseModel `bun:"table:members"`

	ID       int64     `bun:"id,pk,autoincrement"`
	Name     string    `bun:"name,notnull"`
	TeamID   int64     `bun:"team_id,notnull"`
	JoinedAt time.Time `bun:"joined_at,default:current_timestamp"`

	Team *Team `bun:"rel:belongs-to,join:team_id=id"`
}
