// internal/api/middleware/auth.go
package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/eugene-twix/amber-bot/internal/cache"
	"github.com/eugene-twix/amber-bot/internal/domain"
	"github.com/eugene-twix/amber-bot/internal/repository"
	"github.com/gin-gonic/gin"
)

const (
	// TTL for auth_date validation
	ReadTTL  = 1 * time.Hour
	WriteTTL = 10 * time.Minute

	// Context keys
	ContextKeyUser     = "user"
	ContextKeyInitData = "initData"
)

var (
	ErrMissingAuth     = errors.New("missing authorization header")
	ErrInvalidAuth     = errors.New("invalid authorization format")
	ErrInvalidInitData = errors.New("invalid init data")
	ErrExpiredAuth     = errors.New("auth_date expired")
	ErrInvalidHash     = errors.New("invalid hash")
	ErrReplayAttack    = errors.New("replay attack detected")
)

// InitData represents parsed Telegram Web App init data
type InitData struct {
	QueryID  string
	User     *InitDataUser
	AuthDate time.Time
	Hash     string
	Raw      string
}

type InitDataUser struct {
	ID           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type AuthMiddleware struct {
	botToken  string
	userRepo  repository.UserRepository
	cache     *cache.Cache
	secretKey []byte
	devMode   bool
	devUserID int64
}

func NewAuthMiddleware(botToken string, userRepo repository.UserRepository, cache *cache.Cache) *AuthMiddleware {
	// Compute secret key: HMAC_SHA256(bot_token, "WebAppData")
	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secretKey := h.Sum(nil)

	return &AuthMiddleware{
		botToken:  botToken,
		userRepo:  userRepo,
		cache:     cache,
		secretKey: secretKey,
		devMode:   false,
		devUserID: 0,
	}
}

// EnableDevMode enables development mode with a mock user
func (m *AuthMiddleware) EnableDevMode(userID int64) {
	m.devMode = true
	m.devUserID = userID
}

// Authenticate validates Telegram initData and loads user
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Dev mode: bypass authentication
		if m.devMode {
			user, err := m.userRepo.GetOrCreate(c.Request.Context(), m.devUserID, "dev_user")
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
				return
			}
			c.Set(ContextKeyUser, user)
			c.Next()
			return
		}

		// Get Authorization header: "TMA <initData>"
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing_authorization"})
			return
		}

		if !strings.HasPrefix(authHeader, "TMA ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_authorization_format"})
			return
		}

		rawInitData := strings.TrimPrefix(authHeader, "TMA ")

		// Parse and validate initData
		initData, err := m.parseInitData(rawInitData)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Determine TTL based on request method
		ttl := ReadTTL
		if c.Request.Method != http.MethodGet {
			ttl = WriteTTL
		}

		// Validate auth_date
		if time.Since(initData.AuthDate) > ttl {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth_expired"})
			return
		}

		// Validate hash
		if !m.validateHash(rawInitData, initData.Hash) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid_hash"})
			return
		}

		// Check replay attack for mutations
		if c.Request.Method != http.MethodGet && initData.QueryID != "" {
			if err := m.checkReplay(c.Request.Context(), initData.QueryID, c.Request.Method, c.Request.URL.Path); err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "replay_detected"})
				return
			}
		}

		// Get or create user
		user, err := m.userRepo.GetOrCreate(c.Request.Context(), initData.User.ID, initData.User.Username)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal_error"})
			return
		}

		// Store in context
		c.Set(ContextKeyUser, user)
		c.Set(ContextKeyInitData, initData)

		c.Next()
	}
}

// RequireOrganizer checks if user has organizer or admin role
func (m *AuthMiddleware) RequireOrganizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if !user.CanManage() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

// RequireAdmin checks if user has admin role
func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if !user.IsAdmin() {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

// parseInitData parses and extracts data from initData string
func (m *AuthMiddleware) parseInitData(raw string) (*InitData, error) {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return nil, ErrInvalidInitData
	}

	result := &InitData{Raw: raw}

	// Extract hash
	result.Hash = values.Get("hash")
	if result.Hash == "" {
		return nil, ErrInvalidInitData
	}

	// Extract query_id
	result.QueryID = values.Get("query_id")

	// Extract auth_date
	authDateStr := values.Get("auth_date")
	if authDateStr == "" {
		return nil, ErrInvalidInitData
	}
	authDateInt, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidInitData
	}
	result.AuthDate = time.Unix(authDateInt, 0)

	// Extract user
	userStr := values.Get("user")
	if userStr == "" {
		return nil, ErrInvalidInitData
	}

	user, err := parseUserJSON(userStr)
	if err != nil {
		return nil, ErrInvalidInitData
	}
	result.User = user

	return result, nil
}

// parseUserJSON parses user JSON from initData
func parseUserJSON(s string) (*InitDataUser, error) {
	// Simple JSON parsing for user object
	// Format: {"id":123,"first_name":"John",...}
	user := &InitDataUser{}

	// Remove braces
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	// Parse key-value pairs
	for _, pair := range splitJSONPairs(s) {
		parts := strings.SplitN(pair, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.Trim(parts[0], `"`)
		value := strings.Trim(parts[1], `"`)

		switch key {
		case "id":
			id, _ := strconv.ParseInt(value, 10, 64)
			user.ID = id
		case "first_name":
			user.FirstName = value
		case "last_name":
			user.LastName = value
		case "username":
			user.Username = value
		case "language_code":
			user.LanguageCode = value
		}
	}

	if user.ID == 0 {
		return nil, errors.New("missing user id")
	}

	return user, nil
}

// splitJSONPairs splits JSON object into key-value pairs
func splitJSONPairs(s string) []string {
	var pairs []string
	var current strings.Builder
	depth := 0
	inString := false

	for _, r := range s {
		switch r {
		case '"':
			inString = !inString
			current.WriteRune(r)
		case '{', '[':
			depth++
			current.WriteRune(r)
		case '}', ']':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 && !inString {
				pairs = append(pairs, current.String())
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		pairs = append(pairs, current.String())
	}

	return pairs
}

// validateHash validates initData hash using bot token
func (m *AuthMiddleware) validateHash(raw, hash string) bool {
	// Parse query params
	values, err := url.ParseQuery(raw)
	if err != nil {
		return false
	}

	// Remove hash from params
	values.Del("hash")

	// Sort params alphabetically
	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Build data-check-string
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, values.Get(k)))
	}
	dataCheckString := strings.Join(parts, "\n")

	// Compute HMAC
	h := hmac.New(sha256.New, m.secretKey)
	h.Write([]byte(dataCheckString))
	computed := hex.EncodeToString(h.Sum(nil))

	return hmac.Equal([]byte(computed), []byte(hash))
}

// checkReplay checks if exact request was already made (replay attack)
// Uses query_id + method + path to allow multiple different requests in same session
func (m *AuthMiddleware) checkReplay(ctx context.Context, queryID, method, path string) error {
	key := fmt.Sprintf("replay:%s:%s:%s", queryID, method, path)

	// Try to set with NX (only if not exists)
	set, err := m.cache.SetNX(ctx, key, "1", WriteTTL)
	if err != nil {
		// If Redis error, allow request (fail open)
		return nil
	}

	if !set {
		return ErrReplayAttack
	}

	return nil
}

// GetUser extracts user from gin context
func GetUser(c *gin.Context) *domain.User {
	user, exists := c.Get(ContextKeyUser)
	if !exists {
		return nil
	}
	return user.(*domain.User)
}

// GetInitData extracts initData from gin context
func GetInitData(c *gin.Context) *InitData {
	data, exists := c.Get(ContextKeyInitData)
	if !exists {
		return nil
	}
	return data.(*InitData)
}
