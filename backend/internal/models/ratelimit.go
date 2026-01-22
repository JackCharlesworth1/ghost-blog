package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RateLimit struct {
	IPAddress   string    `json:"ip_address"`
	PostCount   int       `json:"post_count"`
	WindowStart time.Time `json:"window_start"`
}

type RateLimitRepository struct {
	db *pgxpool.Pool
}

func NewRateLimitRepository(db *pgxpool.Pool) *RateLimitRepository {
	return &RateLimitRepository{db: db}
}

func (r *RateLimitRepository) CheckAndIncrement(ctx context.Context, ipAddress string, maxPosts int, windowDuration time.Duration) (bool, error) {
	// Start a transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get current rate limit for this IP
	var rateLimit RateLimit
	query := `SELECT ip_address, post_count, window_start FROM ratelimits WHERE ip_address = $1`
	err = tx.QueryRow(ctx, query, ipAddress).Scan(&rateLimit.IPAddress, &rateLimit.PostCount, &rateLimit.WindowStart)

	now := time.Now()

	if err != nil {
		// No existing record, create one
		insertQuery := `INSERT INTO ratelimits (ip_address, post_count, window_start) VALUES ($1, 1, $2)`
		_, err = tx.Exec(ctx, insertQuery, ipAddress, now)
		if err != nil {
			return false, fmt.Errorf("error creating rate limit record: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("error committing transaction: %w", err)
		}
		return true, nil
	}

	// Check if window has expired
	if now.Sub(rateLimit.WindowStart) > windowDuration {
		// Reset the window
		updateQuery := `UPDATE ratelimits SET post_count = 1, window_start = $1 WHERE ip_address = $2`
		_, err = tx.Exec(ctx, updateQuery, now, ipAddress)
		if err != nil {
			return false, fmt.Errorf("error resetting rate limit: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return false, fmt.Errorf("error committing transaction: %w", err)
		}
		return true, nil
	}

	// Check if under limit
	if rateLimit.PostCount >= maxPosts {
		return false, nil
	}

	// Increment counter
	updateQuery := `UPDATE ratelimits SET post_count = post_count + 1 WHERE ip_address = $1`
	_, err = tx.Exec(ctx, updateQuery, ipAddress)
	if err != nil {
		return false, fmt.Errorf("error incrementing rate limit: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return false, fmt.Errorf("error committing transaction: %w", err)
	}

	return true, nil
}

func (r *RateLimitRepository) CleanupExpired(ctx context.Context, windowDuration time.Duration) error {
	query := `DELETE FROM ratelimits WHERE window_start < $1`
	cutoff := time.Now().Add(-windowDuration)

	_, err := r.db.Exec(ctx, query, cutoff)
	if err != nil {
		return fmt.Errorf("error cleaning up expired rate limits: %w", err)
	}

	return nil
}
