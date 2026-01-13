// internal/bot/handlers.go
package bot

import (
	"context"
	"fmt"
	"html"
	"log"
	"strings"

	"github.com/eugene-twix/amber-bot/internal/fsm"
	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleStart(c tele.Context) error {
	user := b.getUser(c)

	// Сначала убираем Reply Keyboard (отдельным сообщением)
	_ = c.Send("👋", &tele.ReplyMarkup{RemoveKeyboard: true})

	msg := fmt.Sprintf(`Привет, %s!

Это бот для учёта результатов квизов и настолок.

Ваша роль: %s`, c.Sender().FirstName, user.Role)

	// Если Mini App URL настроен — добавляем inline кнопку
	if b.miniAppURL != "" {
		msg += "\n\nНажми кнопку ниже, чтобы открыть приложение 👇"

		kb := &tele.ReplyMarkup{}
		webAppBtn := kb.WebApp("🚀 Открыть приложение", &tele.WebApp{URL: b.miniAppURL})
		kb.Inline(kb.Row(webAppBtn))

		return c.Send(msg, kb)
	}

	return c.Send(msg)
}

func (b *Bot) handleText(c tele.Context) error {
	ctx := context.Background()
	text := strings.TrimSpace(c.Text())

	// Check if user is in FSM state
	state, err := b.fsm.Get(ctx, c.Sender().ID)
	if err == nil && state.State != fsm.StateNone {
		// Check for cancel button
		if text == BtnCancel {
			return b.handleCancel(c)
		}
		// Route to FSM handler
		return b.handleFSMText(c, state)
	}

	// Not in FSM — check for main menu buttons
	switch text {
	case BtnTeams:
		return b.handleTeams(c)
	case BtnRating:
		return b.handleRating(c)
	case BtnNewTeam:
		return b.handleNewTeam(c)
	case BtnAddMember:
		return b.handleAddMember(c)
	case BtnNewTournament:
		return b.handleNewTournament(c)
	case BtnResult:
		return b.handleResult(c)
	case BtnGrant:
		return b.handleGrant(c)
	default:
		user := b.getUser(c)
		return c.Send("Используйте кнопки для работы с ботом.", MainMenu(user.Role))
	}
}

func (b *Bot) handleFSMText(c tele.Context, state *fsm.UserState) error {
	switch state.State {
	case fsm.StateNewTeamName:
		return b.processNewTeamName(c, state)
	case fsm.StateNewTeamMoreMember:
		return b.processNewTeamMemberName(c, state)
	case fsm.StateAddMemberName:
		return b.processAddMemberName(c, state)
	case fsm.StateNewTournamentName:
		return b.processNewTournamentName(c, state)
	case fsm.StateNewTournamentDate:
		return b.processNewTournamentDate(c, state)
	case fsm.StateNewTournamentLocation:
		return b.processNewTournamentLocation(c, state)
	case fsm.StateResultPlace:
		return b.processResultPlace(c, state)
	case fsm.StateGrantUser:
		return b.processGrantUser(c, state)
	default:
		user := b.getUser(c)
		return c.Send("Неизвестное состояние. Попробуйте снова.", MainMenu(user.Role))
	}
}

func (b *Bot) handleTeams(c tele.Context) error {
	return b.showTeamsPage(c, 0, false)
}

func (b *Bot) showTeamsPage(c tele.Context, page int, edit bool) error {
	ctx := context.Background()
	teams, err := b.teamRepo.List(ctx)
	if err != nil {
		log.Printf("ERROR: failed to list teams: %v", err)
		return c.Send("Ошибка получения списка команд")
	}

	if len(teams) == 0 {
		return c.Send("Команд пока нет")
	}

	// Convert to paginated items
	items := make([]PaginatedItem, len(teams))
	for i, t := range teams {
		items[i] = PaginatedItem{
			Text: t.Name,
			Data: fmt.Sprintf("team_info:%d", t.ID),
		}
	}

	kb := PaginatedKeyboard("team_page", items, page)

	if edit {
		return c.Edit("Выберите команду:", kb)
	}
	return c.Send("Выберите команду:", kb)
}

func (b *Bot) handleRating(c tele.Context) error {
	return b.showRatingPage(c, 0, false)
}

func (b *Bot) showRatingPage(c tele.Context, page int, edit bool) error {
	ctx := context.Background()
	ratings, err := b.resultRepo.GetTeamRating(ctx)
	if err != nil {
		log.Printf("ERROR: failed to get team rating: %v", err)
		return c.Send("Ошибка получения рейтинга")
	}

	if len(ratings) == 0 {
		return c.Send("Рейтинг пуст — нет результатов")
	}

	const ratingPageSize = 5
	totalPages := (len(ratings) + ratingPageSize - 1) / ratingPageSize

	if page < 0 {
		page = 0
	}
	if page >= totalPages {
		page = totalPages - 1
	}

	start := page * ratingPageSize
	end := start + ratingPageSize
	if end > len(ratings) {
		end = len(ratings)
	}

	var sb strings.Builder
	sb.WriteString("<b>🏆 Рейтинг команд</b>\n")
	sb.WriteString(Separator + "\n\n")
	for i := start; i < end; i++ {
		r := ratings[i]
		medal := "    "
		switch i {
		case 0:
			medal = "🥇 "
		case 1:
			medal = "🥈 "
		case 2:
			medal = "🥉 "
		}
		sb.WriteString(fmt.Sprintf("%s<b>%d. %s</b>\n", medal, i+1, html.EscapeString(r.TeamName)))
		sb.WriteString(fmt.Sprintf("    Побед: <code>%d</code> | Игр: <code>%d</code> | Ср: <code>%.1f</code>\n\n", r.Wins, r.TotalGames, r.AvgPlace))
	}

	// Навигация
	var navRow []tele.InlineButton
	if totalPages > 1 {
		if page > 0 {
			navRow = append(navRow, tele.InlineButton{
				Text: "◀️",
				Data: fmt.Sprintf("rating_page:%d", page-1),
			})
		}
		navRow = append(navRow, tele.InlineButton{
			Text: fmt.Sprintf("%d/%d", page+1, totalPages),
			Data: "noop",
		})
		if page < totalPages-1 {
			navRow = append(navRow, tele.InlineButton{
				Text: "▶️",
				Data: fmt.Sprintf("rating_page:%d", page+1),
			})
		}
	}

	var kb *tele.ReplyMarkup
	if len(navRow) > 0 {
		kb = &tele.ReplyMarkup{InlineKeyboard: [][]tele.InlineButton{navRow}}
	}

	if edit {
		if kb != nil {
			return c.Edit(sb.String(), kb, tele.ModeHTML)
		}
		return c.Edit(sb.String(), tele.ModeHTML)
	}
	if kb != nil {
		return c.Send(sb.String(), kb, tele.ModeHTML)
	}
	return c.Send(sb.String(), tele.ModeHTML)
}

func (b *Bot) handleCancel(c tele.Context) error {
	ctx := context.Background()
	user := b.getUser(c)
	_ = b.fsm.Clear(ctx, c.Sender().ID)
	return c.Send("Действие отменено", MainMenu(user.Role))
}
