// internal/bot/handlers_admin.go
package bot

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/eugene-twix/amber-bot/internal/fsm"
	tele "gopkg.in/telebot.v3"
)

// /grant - выдать роль пользователю
func (b *Bot) handleGrant(c tele.Context) error {
	if !b.requireAdmin(c) {
		return nil
	}

	ctx := context.Background()
	_ = b.fsm.Set(ctx, c.Sender().ID, fsm.StateGrantUser, fsm.Data{})
	return c.Send("Введите Telegram ID пользователя:")
}

func (b *Bot) processGrantUser(c tele.Context, state *fsm.UserState) error {
	ctx := context.Background()
	idStr := strings.TrimSpace(c.Text())

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Send("Неверный ID. Введите числовой Telegram ID:")
	}

	// Check user exists
	_, err = b.userRepo.GetByTelegramID(ctx, userID)
	if err != nil {
		return c.Send("Пользователь не найден. Он должен сначала написать боту /start")
	}

	_ = b.fsm.Update(ctx, c.Sender().ID, fsm.StateGrantRole, "user_id", userID)

	buttons := [][]tele.InlineButton{
		{
			{Text: "Viewer", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleViewer)},
			{Text: "Organizer", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleOrganizer)},
		},
		{
			{Text: "Admin", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleAdmin)},
		},
	}

	return c.Send("Выберите роль:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}
