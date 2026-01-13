// internal/api/handlers/handler.go
package handlers

import (
	"github.com/eugene-twix/amber-bot/internal/cache"
	"github.com/eugene-twix/amber-bot/internal/repository"
)

// Handler holds all API handlers
type Handler struct {
	userRepo       repository.UserRepository
	teamRepo       repository.TeamRepository
	memberRepo     repository.MemberRepository
	tournamentRepo repository.TournamentRepository
	resultRepo     repository.ResultRepository
	cache          *cache.Cache
}

func NewHandler(
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
	memberRepo repository.MemberRepository,
	tournamentRepo repository.TournamentRepository,
	resultRepo repository.ResultRepository,
	cache *cache.Cache,
) *Handler {
	return &Handler{
		userRepo:       userRepo,
		teamRepo:       teamRepo,
		memberRepo:     memberRepo,
		tournamentRepo: tournamentRepo,
		resultRepo:     resultRepo,
		cache:          cache,
	}
}

// Pagination params
type PaginationParams struct {
	Limit  int `form:"limit,default=50"`
	Offset int `form:"offset,default=0"`
}

// SortParams for sorting
type SortParams struct {
	SortBy string `form:"sort_by"`
	Order  string `form:"order,default=asc"`
}

// ListResponse is standard response for list endpoints
type ListResponse struct {
	Items any       `json:"items"`
	Meta  *ListMeta `json:"meta"`
}

type ListMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func NewListResponse(items any, limit, offset, total int) *ListResponse {
	return &ListResponse{
		Items: items,
		Meta: &ListMeta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}
}
