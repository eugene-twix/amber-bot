// internal/domain/user.go
package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Role string

const (
	RoleViewer    Role = "viewer"
	RoleOrganizer Role = "organizer"
	RoleAdmin     Role = "admin"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	TelegramID int64     `bun:"telegram_id,pk"`
	Username   string    `bun:"username"`
	Role       Role      `bun:"role,default:'viewer'"`
	CreatedAt  time.Time `bun:"created_at,default:current_timestamp"`
}

func (u *User) CanManage() bool {
	return u.Role == RoleOrganizer || u.Role == RoleAdmin
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
