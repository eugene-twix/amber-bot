// internal/bot/handlers.go
package bot

import (
	"context"
	"fmt"
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (b *Bot) handleStart(c tele.Context) error {
	user := b.getUser(c)
	msg := fmt.Sprintf(`Привет, %s!

Я бот для учёта результатов квизов и настолок.

Команды:
/teams — список команд
/team <название> — статистика команды
/rating — общий рейтинг

Ваша роль: %s`, c.Sender().FirstName, user.Role)

	if user.CanManage() {
		msg += `

Команды организатора:
/newteam — создать команду
/addmember — добавить участника
/newtournament — создать турнир
/result — записать результат`
	}

	if user.IsAdmin() {
		msg += `

Команды админа:
/grant — выдать роль`
	}

	return c.Send(msg)
}

func (b *Bot) handleTeams(c tele.Context) error {
	ctx := context.Background()
	teams, err := b.teamRepo.List(ctx)
	if err != nil {
		return c.Send("Ошибка получения списка команд")
	}

	if len(teams) == 0 {
		return c.Send("Команд пока нет")
	}

	var sb strings.Builder
	sb.WriteString("Команды:\n\n")
	for _, t := range teams {
		sb.WriteString(fmt.Sprintf("• %s\n", t.Name))
	}

	return c.Send(sb.String())
}

func (b *Bot) handleTeam(c tele.Context) error {
	ctx := context.Background()
	args := c.Args()
	if len(args) == 0 {
		return c.Send("Использование: /team <название команды>")
	}

	teamName := strings.Join(args, " ")
	team, err := b.teamRepo.GetByName(ctx, teamName)
	if err != nil {
		return c.Send(fmt.Sprintf("Команда '%s' не найдена", teamName))
	}

	members, _ := b.memberRepo.GetByTeamID(ctx, team.ID)
	results, _ := b.resultRepo.GetByTeamID(ctx, team.ID)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Команда: %s\n\n", team.Name))

	sb.WriteString("Участники:\n")
	if len(members) == 0 {
		sb.WriteString("  (нет участников)\n")
	} else {
		for _, m := range members {
			sb.WriteString(fmt.Sprintf("  • %s\n", m.Name))
		}
	}

	sb.WriteString("\nРезультаты:\n")
	if len(results) == 0 {
		sb.WriteString("  (нет результатов)\n")
	} else {
		for _, r := range results {
			tournament, _ := b.tournRepo.GetByID(ctx, r.TournamentID)
			if tournament != nil {
				sb.WriteString(fmt.Sprintf("  • %s — %d место\n", tournament.Name, r.Place))
			}
		}
	}

	return c.Send(sb.String())
}

func (b *Bot) handleRating(c tele.Context) error {
	ctx := context.Background()
	ratings, err := b.resultRepo.GetTeamRating(ctx)
	if err != nil {
		return c.Send("Ошибка получения рейтинга")
	}

	if len(ratings) == 0 {
		return c.Send("Рейтинг пуст — нет результатов")
	}

	var sb strings.Builder
	sb.WriteString("Рейтинг команд:\n\n")
	for i, r := range ratings {
		sb.WriteString(fmt.Sprintf("%d. %s — %d побед, %d игр, ср. место: %.1f\n",
			i+1, r.TeamName, r.Wins, r.TotalGames, r.AvgPlace))
	}

	return c.Send(sb.String())
}

func (b *Bot) handleCancel(c tele.Context) error {
	ctx := context.Background()
	_ = b.fsm.Clear(ctx, c.Sender().ID)
	return c.Send("Действие отменено")
}
