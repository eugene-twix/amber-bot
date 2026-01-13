// internal/api/handlers/private.go
package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/eugene-twix/amber-bot/internal/api/middleware"
	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/gin-gonic/gin"
)

// === TEAMS ===

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

func (h *Handler) CreateTeam(c *gin.Context) {
	var req CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	user := middleware.GetUser(c)

	team := &domain.Team{
		Name:      req.Name,
		CreatedBy: user.TelegramID,
	}

	if err := h.teamRepo.Create(c.Request.Context(), team); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         team.ID,
		"name":       team.Name,
		"created_at": team.CreatedAt.Format(time.RFC3339),
		"version":    team.Version,
	})
}

type UpdateTeamRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=100"`
	Version int    `json:"version" binding:"required,min=1"`
}

func (h *Handler) UpdateTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var req UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// Get current team
	team, err := h.teamRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team_not_found"})
		return
	}

	// Check version (optimistic locking)
	if team.Version != req.Version {
		c.JSON(http.StatusConflict, gin.H{"error": "version_conflict", "current_version": team.Version})
		return
	}

	user := middleware.GetUser(c)
	now := time.Now()

	team.Name = req.Name
	team.UpdatedAt = &now
	team.UpdatedBy = &user.TelegramID
	team.Version = req.Version + 1

	// TODO: Implement Update method in repository
	c.JSON(http.StatusOK, gin.H{
		"id":         team.ID,
		"name":       team.Name,
		"updated_at": team.UpdatedAt.Format(time.RFC3339),
		"version":    team.Version,
	})
}

type DeleteRequest struct {
	Version int `json:"version" binding:"required,min=1"`
}

func (h *Handler) DeleteTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// Get current team
	team, err := h.teamRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team_not_found"})
		return
	}

	// Check version
	if team.Version != req.Version {
		c.JSON(http.StatusConflict, gin.H{"error": "version_conflict", "current_version": team.Version})
		return
	}

	if err := h.teamRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// === MEMBERS ===

type CreateMemberRequest struct {
	Name string `json:"name" binding:"required,min=1,max=100"`
}

func (h *Handler) CreateMember(c *gin.Context) {
	teamIDStr := c.Param("id")
	teamID, err := strconv.ParseInt(teamIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_team_id"})
		return
	}

	var req CreateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// Check team exists
	_, err = h.teamRepo.GetByID(c.Request.Context(), teamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team_not_found"})
		return
	}

	user := middleware.GetUser(c)

	member := &domain.Member{
		Name:      req.Name,
		TeamID:    teamID,
		CreatedBy: user.TelegramID,
	}

	if err := h.memberRepo.Create(c.Request.Context(), member); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":        member.ID,
		"name":      member.Name,
		"team_id":   member.TeamID,
		"joined_at": member.JoinedAt.Format(time.RFC3339),
		"version":   member.Version,
	})
}

type UpdateMemberRequest struct {
	Name    string `json:"name" binding:"required,min=1,max=100"`
	Version int    `json:"version" binding:"required,min=1"`
}

func (h *Handler) UpdateMember(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_member_id"})
		return
	}

	var req UpdateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// TODO: Implement GetByID and Update for member
	c.JSON(http.StatusOK, gin.H{
		"id":      memberID,
		"name":    req.Name,
		"version": req.Version + 1,
	})
}

func (h *Handler) DeleteMember(c *gin.Context) {
	memberIDStr := c.Param("member_id")
	memberID, err := strconv.ParseInt(memberIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_member_id"})
		return
	}

	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	if err := h.memberRepo.Delete(c.Request.Context(), memberID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// === TOURNAMENTS ===

type CreateTournamentRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=200"`
	Date     string `json:"date" binding:"required"` // Format: 2006-01-02
	Location string `json:"location" binding:"max=200"`
}

func (h *Handler) CreateTournament(c *gin.Context) {
	var req CreateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_date_format"})
		return
	}

	user := middleware.GetUser(c)

	tournament := &domain.Tournament{
		Name:      req.Name,
		Date:      date,
		Location:  req.Location,
		CreatedBy: user.TelegramID,
	}

	if err := h.tournamentRepo.Create(c.Request.Context(), tournament); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         tournament.ID,
		"name":       tournament.Name,
		"date":       tournament.Date.Format("2006-01-02"),
		"location":   tournament.Location,
		"created_at": tournament.CreatedAt.Format(time.RFC3339),
		"version":    tournament.Version,
	})
}

type UpdateTournamentRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=200"`
	Date     string `json:"date" binding:"required"`
	Location string `json:"location" binding:"max=200"`
	Version  int    `json:"version" binding:"required,min=1"`
}

func (h *Handler) UpdateTournament(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var req UpdateTournamentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	tournament, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament_not_found"})
		return
	}

	if tournament.Version != req.Version {
		c.JSON(http.StatusConflict, gin.H{"error": "version_conflict", "current_version": tournament.Version})
		return
	}

	// TODO: Implement Update
	c.JSON(http.StatusOK, gin.H{
		"id":      tournament.ID,
		"name":    req.Name,
		"version": req.Version + 1,
	})
}

func (h *Handler) DeleteTournament(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	tournament, err := h.tournamentRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament_not_found"})
		return
	}

	if tournament.Version != req.Version {
		c.JSON(http.StatusConflict, gin.H{"error": "version_conflict", "current_version": tournament.Version})
		return
	}

	// TODO: Implement Delete in tournament repo
	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

// === RESULTS ===

type CreateResultRequest struct {
	TeamID int64 `json:"team_id" binding:"required"`
	Place  int   `json:"place" binding:"required,min=1,max=1000"`
}

func (h *Handler) CreateResult(c *gin.Context) {
	tournamentIDStr := c.Param("id")
	tournamentID, err := strconv.ParseInt(tournamentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_tournament_id"})
		return
	}

	var req CreateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// Verify tournament exists
	_, err = h.tournamentRepo.GetByID(c.Request.Context(), tournamentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tournament_not_found"})
		return
	}

	// Verify team exists
	_, err = h.teamRepo.GetByID(c.Request.Context(), req.TeamID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "team_not_found"})
		return
	}

	user := middleware.GetUser(c)

	result := &domain.Result{
		TournamentID: tournamentID,
		TeamID:       req.TeamID,
		Place:        req.Place,
		RecordedBy:   user.TelegramID,
	}

	if err := h.resultRepo.Create(c.Request.Context(), result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	// Invalidate rating cache
	_ = h.cache.Delete(c.Request.Context(), ratingCacheKey)

	c.JSON(http.StatusCreated, gin.H{
		"id":            result.ID,
		"tournament_id": result.TournamentID,
		"team_id":       result.TeamID,
		"place":         result.Place,
		"recorded_at":   result.RecordedAt.Format(time.RFC3339),
		"version":       result.Version,
	})
}

type UpdateResultRequest struct {
	Place   int `json:"place" binding:"required,min=1,max=1000"`
	Version int `json:"version" binding:"required,min=1"`
}

func (h *Handler) UpdateResult(c *gin.Context) {
	resultIDStr := c.Param("result_id")
	resultID, err := strconv.ParseInt(resultIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_result_id"})
		return
	}

	var req UpdateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// TODO: Implement GetByID and Update for result
	// Invalidate rating cache
	_ = h.cache.Delete(c.Request.Context(), ratingCacheKey)

	c.JSON(http.StatusOK, gin.H{
		"id":      resultID,
		"place":   req.Place,
		"version": req.Version + 1,
	})
}

func (h *Handler) DeleteResult(c *gin.Context) {
	resultIDStr := c.Param("result_id")
	resultID, err := strconv.ParseInt(resultIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_result_id"})
		return
	}

	var req DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// TODO: Implement GetByID for version check and Delete
	// Invalidate rating cache
	_ = h.cache.Delete(c.Request.Context(), ratingCacheKey)

	c.JSON(http.StatusOK, gin.H{"deleted": true, "id": resultID})
}

// === USERS (Admin only) ===

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.userRepo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	items := make([]gin.H, 0, len(users))
	for _, u := range users {
		items = append(items, gin.H{
			"telegram_id": u.TelegramID,
			"username":    u.Username,
			"role":        u.Role,
			"created_at":  u.CreatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, NewListResponse(items, 50, 0, len(items)))
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=viewer organizer admin"`
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	telegramIDStr := c.Param("telegram_id")
	telegramID, err := strconv.ParseInt(telegramIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_telegram_id"})
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "details": err.Error()})
		return
	}

	// Check user exists
	user, err := h.userRepo.GetByTelegramID(c.Request.Context(), telegramID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user_not_found"})
		return
	}

	// Cannot change own role
	currentUser := middleware.GetUser(c)
	if currentUser.TelegramID == telegramID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot_change_own_role"})
		return
	}

	if err := h.userRepo.UpdateRole(c.Request.Context(), telegramID, domain.Role(req.Role)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"telegram_id": user.TelegramID,
		"username":    user.Username,
		"role":        req.Role,
	})
}
