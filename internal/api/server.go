// internal/api/server.go
package api

import (
	"context"
	"fmt"
	"time"

	"github.com/eugene-twix/amber-bot/internal/api/handlers"
	"github.com/eugene-twix/amber-bot/internal/api/middleware"
	"github.com/eugene-twix/amber-bot/internal/cache"
	"github.com/eugene-twix/amber-bot/internal/repository"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Port         int
	BotToken     string
	FrontendPath string // Path to frontend/dist
}

type Server struct {
	config  Config
	engine  *gin.Engine
	handler *handlers.Handler
}

func NewServer(cfg Config, repos *Repositories, cache *cache.Cache) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// Global middleware
	engine.Use(gin.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.Timeout(15 * time.Second))
	engine.Use(middleware.CORS())

	// Create handler with all dependencies
	h := handlers.NewHandler(repos.User, repos.Team, repos.Member, repos.Tournament, repos.Result, cache)

	// Auth middleware
	authMW := middleware.NewAuthMiddleware(cfg.BotToken, repos.User, cache)
	rateLimitMW := middleware.NewRateLimitMiddleware(cache)

	s := &Server{
		config:  cfg,
		engine:  engine,
		handler: h,
	}

	s.setupRoutes(authMW, rateLimitMW)

	return s
}

func (s *Server) setupRoutes(authMW *middleware.AuthMiddleware, rateLimitMW *middleware.RateLimitMiddleware) {
	api := s.engine.Group("/api/v1")

	// Public routes (Viewer+)
	public := api.Group("/public")
	public.Use(authMW.Authenticate())
	public.Use(rateLimitMW.LimitRead())
	{
		public.GET("/me", s.handler.GetMe)
		public.GET("/teams", s.handler.ListTeams)
		public.GET("/teams/:id", s.handler.GetTeam)
		public.GET("/teams/:id/members", s.handler.ListTeamMembers)
		public.GET("/teams/:id/results", s.handler.ListTeamResults)
		public.GET("/tournaments", s.handler.ListTournaments)
		public.GET("/tournaments/:id", s.handler.GetTournament)
		public.GET("/tournaments/:id/results", s.handler.ListTournamentResults)
		public.GET("/rating", s.handler.GetRating)
	}

	// Private routes (Organizer/Admin)
	private := api.Group("/private")
	private.Use(authMW.Authenticate())
	private.Use(authMW.RequireOrganizer())
	{
		// Teams
		private.POST("/teams", rateLimitMW.LimitWrite(), s.handler.CreateTeam)
		private.PATCH("/teams/:id", rateLimitMW.LimitWrite(), s.handler.UpdateTeam)
		private.DELETE("/teams/:id", rateLimitMW.LimitWrite(), s.handler.DeleteTeam)

		// Members
		private.POST("/teams/:id/members", rateLimitMW.LimitWrite(), s.handler.CreateMember)
		private.PATCH("/teams/:team_id/members/:member_id", rateLimitMW.LimitWrite(), s.handler.UpdateMember)
		private.DELETE("/teams/:team_id/members/:member_id", rateLimitMW.LimitWrite(), s.handler.DeleteMember)

		// Tournaments
		private.POST("/tournaments", rateLimitMW.LimitWrite(), s.handler.CreateTournament)
		private.PATCH("/tournaments/:id", rateLimitMW.LimitWrite(), s.handler.UpdateTournament)
		private.DELETE("/tournaments/:id", rateLimitMW.LimitWrite(), s.handler.DeleteTournament)

		// Results
		private.POST("/tournaments/:id/results", rateLimitMW.LimitWrite(), s.handler.CreateResult)
		private.PATCH("/tournaments/:id/results/:result_id", rateLimitMW.LimitWrite(), s.handler.UpdateResult)
		private.DELETE("/tournaments/:id/results/:result_id", rateLimitMW.LimitWrite(), s.handler.DeleteResult)

		// Admin only routes
		admin := private.Group("")
		admin.Use(authMW.RequireAdmin())
		{
			admin.GET("/users", rateLimitMW.LimitRead(), s.handler.ListUsers)
			admin.PUT("/users/:telegram_id/role", rateLimitMW.LimitWrite(), s.handler.UpdateUserRole)
		}
	}

	// Serve frontend static files
	if s.config.FrontendPath != "" {
		s.engine.Static("/assets", s.config.FrontendPath+"/assets")
		s.engine.StaticFile("/favicon.ico", s.config.FrontendPath+"/favicon.ico")

		// SPA fallback - serve index.html for all non-API routes
		s.engine.NoRoute(func(c *gin.Context) {
			c.File(s.config.FrontendPath + "/index.html")
		})
	}
}

func (s *Server) Run() error {
	addr := fmt.Sprintf(":%d", s.config.Port)
	return s.engine.Run(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Gin doesn't have built-in graceful shutdown, but we can implement it
	// For now, just return nil
	return nil
}

// Repositories holds all repository interfaces
type Repositories struct {
	User       repository.UserRepository
	Team       repository.TeamRepository
	Member     repository.MemberRepository
	Tournament repository.TournamentRepository
	Result     repository.ResultRepository
}
