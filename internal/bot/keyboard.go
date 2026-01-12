// internal/bot/keyboard.go
package bot

import (
	"fmt"

	"github.com/eugene-twix/amber-bot/internal/domain"
	tele "gopkg.in/telebot.v3"
)

// Button text constants
const (
	BtnTeams         = "📋 Команды"
	BtnRating        = "🏆 Рейтинг"
	BtnNewTeam       = "➕ Новая команда"
	BtnAddMember     = "👤 Добавить игрока"
	BtnNewTournament = "🎯 Новый турнир"
	BtnResult        = "🏅 Записать место"
	BtnGrant         = "👑 Права"
	BtnCancel        = "❌ Отмена"
)

// PageSize - количество элементов на странице
const PageSize = 5

// Separator - разделитель для сообщений
const Separator = "─────────────────"

// PaginatedItem - элемент для пагинированного списка
type PaginatedItem struct {
	Text string
	Data string
}

// PaginatedKeyboard создаёт inline-клавиатуру с пагинацией
// action - префикс для callback навигации (например "team_page", "grant_page")
// items - все элементы списка
// page - текущая страница (0-based)
func PaginatedKeyboard(action string, items []PaginatedItem, page int) *tele.ReplyMarkup {
	totalPages := (len(items) + PageSize - 1) / PageSize
	if totalPages == 0 {
		totalPages = 1
	}

	// Корректируем страницу если выходит за границы
	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}

	// Вычисляем индексы для текущей страницы
	start := page * PageSize
	end := start + PageSize
	if end > len(items) {
		end = len(items)
	}

	var rows [][]tele.InlineButton

	// Добавляем элементы текущей страницы
	for _, item := range items[start:end] {
		rows = append(rows, []tele.InlineButton{
			{Text: item.Text, Data: item.Data},
		})
	}

	// Добавляем навигацию если нужна
	if totalPages > 1 {
		var navRow []tele.InlineButton

		if page > 0 {
			navRow = append(navRow, tele.InlineButton{
				Text: "◀️ Назад",
				Data: fmt.Sprintf("%s:%d", action, page-1),
			})
		}

		navRow = append(navRow, tele.InlineButton{
			Text: fmt.Sprintf("%d/%d", page+1, totalPages),
			Data: "noop",
		})

		if page < totalPages-1 {
			navRow = append(navRow, tele.InlineButton{
				Text: "Вперёд ▶️",
				Data: fmt.Sprintf("%s:%d", action, page+1),
			})
		}

		rows = append(rows, navRow)
	}

	return &tele.ReplyMarkup{InlineKeyboard: rows}
}

// MainMenu returns Reply Keyboard based on user role
func MainMenu(role domain.Role) *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}

	btnTeams := menu.Text(BtnTeams)
	btnRating := menu.Text(BtnRating)
	btnNewTeam := menu.Text(BtnNewTeam)
	btnAddMember := menu.Text(BtnAddMember)
	btnNewTournament := menu.Text(BtnNewTournament)
	btnResult := menu.Text(BtnResult)
	btnGrant := menu.Text(BtnGrant)

	switch role {
	case domain.RoleViewer:
		menu.Reply(
			menu.Row(btnTeams, btnRating),
		)
	case domain.RoleOrganizer:
		menu.Reply(
			menu.Row(btnTeams, btnRating),
			menu.Row(btnNewTeam, btnAddMember),
			menu.Row(btnNewTournament, btnResult),
		)
	case domain.RoleAdmin:
		menu.Reply(
			menu.Row(btnTeams, btnRating),
			menu.Row(btnNewTeam, btnAddMember),
			menu.Row(btnNewTournament, btnResult),
			menu.Row(btnGrant),
		)
	}

	return menu
}

// CancelMenu returns Reply Keyboard with cancel button (for FSM dialogs)
func CancelMenu() *tele.ReplyMarkup {
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	menu.Reply(
		menu.Row(menu.Text(BtnCancel)),
	)
	return menu
}
