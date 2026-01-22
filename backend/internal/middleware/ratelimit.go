package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/jimmyrecce/ghost-blog/internal/models"
)

type RateLimiter struct {
	repo          *models.RateLimitRepository
	maxPosts      int
	windowDuration time.Duration
}

func NewRateLimiter(repo *models.RateLimitRepository) *RateLimiter {
	return &RateLimiter{
		repo:          repo,
		maxPosts:      5,
		windowDuration: time.Hour,
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := GetClientIP(r)

		allowed, err := rl.repo.CheckAndIncrement(context.Background(), ipAddress, rl.maxPosts, rl.windowDuration)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Rate limit exceeded. Maximum 5 posts per hour.", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
