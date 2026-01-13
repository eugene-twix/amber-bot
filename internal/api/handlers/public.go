// internal/api/handlers/public.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/eugene-twix/amber-bot/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

const (
	ratingCacheKey = "api:rating"
	ratingCacheTTL = 5 * time.Minute
)

// GetMe returns current user info
func (h *Handler) GetMe(c *gin.Context) {
	user := middleware.GetUser(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"telegram_id": user.TelegramID,
		"username":    user.Username,
		"role":        user.Role,
		"created_at":  user.CreatedAt.Format(time.RFC3339),
	})
}

// ListTeams returns list of teams
func (h *Handler) ListTeams(c *gin.Context) {
	teams, err := h.teamRepo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	// Transform to response format
	items := make([]gin.H, 0, len(teams))
	for _, t := range teams {
		items = append(items, gin.H{
			"id":         t.ID,
			"name":       t.Name,
			"created_at": t.CreatedAt.Format(time.RFC3339),
			"version":    t.Version,
		})
	}

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}

// GetTeam returns team details
func (h *Handler) GetTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	team, err := h.teamRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team_not_found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         team.ID,
		"name":       team.Name,
		"created_at": team.CreatedAt.Format(time.RFC3339),
		"created_by": team.CreatedBy,
		"version":    team.Version,
	})
}

// ListTeamMembers returns team members
func (h *Handler) ListTeamMembers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	members, err := h.memberRepo.GetByTeamID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	items := make([]gin.H, 0, len(members))
	for _, m := range members {
		items = append(items, gin.H{
			"id":        m.ID,
			"name":      m.Name,
			"team_id":   m.TeamID,
			"joined_at": m.JoinedAt.Format(time.RFC3339),
			"version":   m.Version,
		})
	}

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}

// ListTeamResults returns team's tournament results
func (h *Handler) ListTeamResults(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	results, err := h.resultRepo.GetByTeamID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	items := make([]gin.H, 0, len(results))
	for _, r := range results {
		item := gin.H{
			"id":            r.ID,
			"team_id":       r.TeamID,
			"tournament_id": r.TournamentID,
			"place":         r.Place,
			"recorded_at":   r.RecordedAt.Format(time.RFC3339),
			"version":       r.Version,
		}
		if r.Tournament != nil {
			item["tournament_name"] = r.Tournament.Name
			item["tournament_date"] = r.Tournament.Date.Format("2006-01-02")
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}

// ListTournaments returns list of tournaments
func (h *Handler) ListTournaments(c *gin.Context) {
	tournaments, err := h.tournamentRepo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	items := make([]gin.H, 0, len(tournaments))
	for _, t := range tournaments {
		items = append(items, gin.H{
			"id":         t.ID,
			"name":       t.Name,
			"date":       t.Date.Format("2006-01-02"),
			"location":   t.Location,
			"created_at": t.CreatedAt.Format(time.RFC3339),
			"version":    t.Version,
		})
	}

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}

// GetTournament returns tournament details
func (h *Handler) GetTournament(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	tournament, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament_not_found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         tournament.ID,
		"name":       tournament.Name,
		"date":       tournament.Date.Format("2006-01-02"),
		"location":   tournament.Location,
		"created_at": tournament.CreatedAt.Format(time.RFC3339),
		"created_by": tournament.CreatedBy,
		"version":    tournament.Version,
	})
}

// ListTournamentResults returns tournament results
func (h *Handler) ListTournamentResults(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	results, err := h.resultRepo.GetByTournamentID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	items := make([]gin.H, 0, len(results))
	for _, r := range results {
		item := gin.H{
			"id":            r.ID,
			"team_id":       r.TeamID,
			"tournament_id": r.TournamentID,
			"place":         r.Place,
			"recorded_at":   r.RecordedAt.Format(time.RFC3339),
			"version":       r.Version,
		}
		if r.Team != nil {
			item["team_name"] = r.Team.Name
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}

// GetRating returns team ratings (cached)
func (h *Handler) GetRating(c *gin.Context) {
	// Try cache first
	var cachedRating []gin.H
	if err := h.cache.Get(c.Request.Context(), ratingCacheKey, &cachedRating); err == nil {
		c.JSON(http.StatusOK, NewListResponse(cachedRating, 50, 0, len(cachedRating)))
		return
	}

	// Get from DB
	ratings, err := h.resultRepo.GetTeamRating(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	items := make([]gin.H, 0, len(ratings))
	for _, r := range ratings {
		items = append(items, gin.H{
			"team_id":     r.TeamID,
			"team_name":   r.TeamName,
			"top_places":  r.Wins,
			"total_games": r.TotalGames,
			"avg_place":   r.AvgPlace,
		})
	}

	// Cache result
	_ = h.cache.Set(c.Request.Context(), ratingCacheKey, items, ratingCacheTTL)

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}
