// internal/bot/handlers_admin.go
package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/eugene-twix/amber-bot/internal/fsm"
	tele "gopkg.in/telebot.v3"
)

// handleGrant - выдать роль пользователю
func (b *Bot) handleGrant(c tele.Context) error {
	if !b.requireAdmin(c) {
		return nil
	}
	return b.showGrantUsersPage(c, 0, false)
}

func (b *Bot) showGrantUsersPage(c tele.Context, page int, edit bool) error {
	if !b.requireAdmin(c) {
		return nil
	}

	ctx := context.Background()

	// Get list of users
	users, err := b.userRepo.List(ctx)
	if err != nil {
		log.Printf("ERROR: failed to list users: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.", MainMenu(domain.RoleAdmin))
	}

	// Build paginated items
	var items []PaginatedItem
	for _, u := range users {
		// Skip current admin
		if u.TelegramID == c.Sender().ID {
			continue
		}

		displayName := u.Username
		if displayName == "" {
			displayName = fmt.Sprintf("ID: %d", u.TelegramID)
		}

		// Add role icon
		roleIcon := "👀"
		switch u.Role {
		case domain.RoleOrganizer:
			roleIcon = "📝"
		case domain.RoleAdmin:
			roleIcon = "👑"
		}

		items = append(items, PaginatedItem{
			Text: fmt.Sprintf("%s %s", roleIcon, displayName),
			Data: fmt.Sprintf("grant_user:%d", u.TelegramID),
		})
	}

	// Add option to enter ID manually at the end
	items = append(items, PaginatedItem{
		Text: "✏️ Ввести ID вручную",
		Data: "grant_user:manual",
	})

	if len(items) == 1 {
		// Only "enter manually" option
		if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateGrantUser, fsm.Data{}); err != nil {
			log.Printf("ERROR: failed to set FSM state: %v", err)
			return c.Send("Ошибка сервиса. Попробуйте позже.")
		}
		return c.Send("Нет других пользователей. Введите Telegram ID:", CancelMenu())
	}

	kb := PaginatedKeyboard("grant_page", items, page)

	if edit {
		return c.Edit("Выберите пользователя:", kb)
	}
	return c.Send("Выберите пользователя:", kb)
}

// handleGrantUserCallback - обработка выбора пользователя из списка
func (b *Bot) handleGrantUserCallback(c tele.Context, payload string) error {
	if !b.requireAdmin(c) {
		return nil
	}

	ctx := context.Background()

	// Handle manual input option
	if payload == "manual" {
		if err := b.fsm.Set(ctx, c.Sender().ID, fsm.StateGrantUser, fsm.Data{}); err != nil {
			log.Printf("ERROR: failed to set FSM state: %v", err)
			return c.Send("Ошибка сервиса. Попробуйте позже.")
		}
		return c.Edit("Введите Telegram ID пользователя:")
	}

	// Parse user ID
	userID, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		log.Printf("ERROR: invalid user ID in callback: %s", payload)
		return c.Respond(&tele.CallbackResponse{Text: "Ошибка"})
	}

	// Get user to display name
	targetUser, err := b.userRepo.GetByTelegramID(ctx, userID)
	if err != nil {
		log.Printf("ERROR: failed to get user: %v", err)
		return c.Respond(&tele.CallbackResponse{Text: "Пользователь не найден"})
	}

	// Show role selection
	buttons := [][]tele.InlineButton{
		{
			{Text: "👀 Viewer", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleViewer)},
			{Text: "📝 Organizer", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleOrganizer)},
		},
		{
			{Text: "👑 Admin", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleAdmin)},
		},
	}

	displayName := targetUser.Username
	if displayName == "" {
		displayName = fmt.Sprintf("ID: %d", userID)
	}

	return c.Edit(fmt.Sprintf("Выберите роль для %s:", displayName), &tele.ReplyMarkup{InlineKeyboard: buttons})
}

func (b *Bot) processGrantUser(c tele.Context, _ *fsm.UserState) error {
	ctx := context.Background()

	// Verify state to prevent race condition
	if _, err := b.verifyState(ctx, c.Sender().ID, fsm.StateGrantUser); err != nil {
		log.Printf("ERROR: state verification failed: %v", err)
		user := b.getUser(c)
		return c.Send("Состояние изменилось. Попробуйте снова.", MainMenu(user.Role))
	}

	idStr := strings.TrimSpace(c.Text())

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Send("Неверный ID. Введите числовой Telegram ID:", CancelMenu())
	}

	// Check user exists
	_, err = b.userRepo.GetByTelegramID(ctx, userID)
	if err != nil {
		log.Printf("ERROR: failed to get user by telegram ID: %v", err)
		return c.Send("Пользователь не найден. Он должен сначала написать боту /start", CancelMenu())
	}

	if err := b.fsm.Update(ctx, c.Sender().ID, fsm.StateGrantRole, "user_id", userID); err != nil {
		log.Printf("ERROR: failed to update FSM state: %v", err)
		return c.Send("Ошибка сервиса. Попробуйте позже.")
	}

	buttons := [][]tele.InlineButton{
		{
			{Text: "👀 Viewer", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleViewer)},
			{Text: "📝 Organizer", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleOrganizer)},
		},
		{
			{Text: "👑 Admin", Data: fmt.Sprintf("grant_role:%d:%s", userID, domain.RoleAdmin)},
		},
	}

	return c.Send("Выберите роль:", &tele.ReplyMarkup{InlineKeyboard: buttons})
}
