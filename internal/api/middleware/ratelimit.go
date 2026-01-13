// internal/api/middleware/ratelimit.go
package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/eugene-twix/amber-bot/internal/cache"
	"github.com/gin-gonic/gin"
)

const (
	ReadRateLimit  = 100 // requests per minute for GET
	WriteRateLimit = 20  // requests per minute for POST/PATCH/DELETE
	RateLimitTTL   = 1 * time.Minute
)

type RateLimitMiddleware struct {
	cache *cache.Cache
}

func NewRateLimitMiddleware(cache *cache.Cache) *RateLimitMiddleware {
	return &RateLimitMiddleware{cache: cache}
}

// LimitRead applies rate limit for GET requests (100 req/min)
func (m *RateLimitMiddleware) LimitRead() gin.HandlerFunc {
	return m.limit(ReadRateLimit, "read")
}

// LimitWrite applies rate limit for mutation requests (20 req/min)
func (m *RateLimitMiddleware) LimitWrite() gin.HandlerFunc {
	return m.limit(WriteRateLimit, "write")
}

func (m *RateLimitMiddleware) limit(maxRequests int, prefix string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := GetUser(c)
		if user == nil {
			c.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:%s:%d", prefix, user.TelegramID)

		count, err := m.increment(c.Request.Context(), key)
		if err != nil {
			// Fail open on Redis errors
			c.Next()
			return
		}

		if count > int64(maxRequests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "rate_limit_exceeded",
				"retry_after": 60,
			})
			return
		}

		c.Next()
	}
}

func (m *RateLimitMiddleware) increment(ctx context.Context, key string) (int64, error) {
	count, err := m.cache.Incr(ctx, key)
	if err != nil {
		return 0, err
	}

	// Set TTL on first increment
	if count == 1 {
		_ = m.cache.Expire(ctx, key, RateLimitTTL)
	}

	return count, nil
}
