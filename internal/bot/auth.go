// internal/bot/auth.go
package bot

import (
	"context"

	"github.com/eugene-twix/amber-bot/internal/domain"
	tele "gopkg.in/telebot.v3"
)

func (b *Bot) authMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		ctx := context.Background()
		sender := c.Sender()

		user, err := b.userRepo.GetOrCreate(ctx, sender.ID, sender.Username)
		if err != nil {
			return c.Send("Ошибка авторизации")
		}

		// Check if user is in admin list from config
		for _, adminID := range b.cfg.AdminIDs {
			if sender.ID == adminID && user.Role != domain.RoleAdmin {
				_ = b.userRepo.UpdateRole(ctx, sender.ID, domain.RoleAdmin)
				user.Role = domain.RoleAdmin
			}
		}

		c.Set(string(userKey), user)
		return next(c)
	}
}

func (b *Bot) getUser(c tele.Context) *domain.User {
	if u, ok := c.Get(string(userKey)).(*domain.User); ok {
		return u
	}
	return nil
}

func (b *Bot) requireOrganizer(c tele.Context) bool {
	user := b.getUser(c)
	if user == nil || !user.CanManage() {
		c.Send("У вас нет прав для этой команды")
		return false
	}
	return true
}

func (b *Bot) requireAdmin(c tele.Context) bool {
	user := b.getUser(c)
	if user == nil || !user.IsAdmin() {
		c.Send("Только администратор может выполнить эту команду")
		return false
	}
	return true
}
