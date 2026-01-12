// internal/bot/bot.go
package bot

import (
	"log"
	"time"

	"github.com/eugene-twix/amber-bot/internal/cache"
	"github.com/eugene-twix/amber-bot/internal/config"
	"github.com/eugene-twix/amber-bot/internal/fsm"
	bunrepo "github.com/eugene-twix/amber-bot/internal/repository/bun"
	"github.com/uptrace/bun"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	tg         *tele.Bot
	cfg        *config.Config
	db         *bun.DB
	cache      *cache.Cache
	fsm        *fsm.Manager
	userRepo   *bunrepo.UserRepo
	teamRepo   *bunrepo.TeamRepo
	memberRepo *bunrepo.MemberRepo
	tournRepo  *bunrepo.TournamentRepo
	resultRepo *bunrepo.ResultRepo
}

func New(cfg *config.Config, db *bun.DB, cache *cache.Cache) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.TelegramToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	tg, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		tg:         tg,
		cfg:        cfg,
		db:         db,
		cache:      cache,
		fsm:        fsm.NewManager(cache),
		userRepo:   bunrepo.NewUserRepo(db),
		teamRepo:   bunrepo.NewTeamRepo(db),
		memberRepo: bunrepo.NewMemberRepo(db),
		tournRepo:  bunrepo.NewTournamentRepo(db),
		resultRepo: bunrepo.NewResultRepo(db),
	}

	b.registerHandlers()
	return b, nil
}

func (b *Bot) registerHandlers() {
	// Middleware
	b.tg.Use(b.authMiddleware)

	// Public commands
	b.tg.Handle("/start", b.handleStart)
	b.tg.Handle("/teams", b.handleTeams)
	b.tg.Handle("/team", b.handleTeam)
	b.tg.Handle("/rating", b.handleRating)
	b.tg.Handle("/cancel", b.handleCancel)

	// Organizer commands
	b.tg.Handle("/newteam", b.handleNewTeam)
	b.tg.Handle("/addmember", b.handleAddMember)
	b.tg.Handle("/newtournament", b.handleNewTournament)
	b.tg.Handle("/result", b.handleResult)

	// Admin commands
	b.tg.Handle("/grant", b.handleGrant)

	// Callbacks
	b.tg.Handle(tele.OnCallback, b.handleCallback)

	// Text messages (for FSM)
	b.tg.Handle(tele.OnText, b.handleText)
}

func (b *Bot) Start() {
	log.Println("Bot started")
	b.tg.Start()
}

func (b *Bot) Stop() {
	b.tg.Stop()
}

// Context key for user
type ctxKey string

const userKey ctxKey = "user"

// Stub handlers - will be implemented in later tasks

func (b *Bot) handleStart(c tele.Context) error {
	return nil
}

func (b *Bot) handleTeams(c tele.Context) error {
	return nil
}

func (b *Bot) handleTeam(c tele.Context) error {
	return nil
}

func (b *Bot) handleRating(c tele.Context) error {
	return nil
}

func (b *Bot) handleCancel(c tele.Context) error {
	return nil
}

func (b *Bot) handleNewTeam(c tele.Context) error {
	return nil
}

func (b *Bot) handleAddMember(c tele.Context) error {
	return nil
}

func (b *Bot) handleNewTournament(c tele.Context) error {
	return nil
}

func (b *Bot) handleResult(c tele.Context) error {
	return nil
}

func (b *Bot) handleGrant(c tele.Context) error {
	return nil
}

func (b *Bot) handleCallback(c tele.Context) error {
	return nil
}

func (b *Bot) handleText(c tele.Context) error {
	return nil
}
