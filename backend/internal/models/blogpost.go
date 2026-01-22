package models

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BlogPost struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
}

type BlogPostRepository struct {
	db *pgxpool.Pool
}

func NewBlogPostRepository(db *pgxpool.Pool) *BlogPostRepository {
	return &BlogPostRepository{db: db}
}

func (r *BlogPostRepository) Create(ctx context.Context, post *BlogPost) error {
	query := `
		INSERT INTO blogposts (title, content, country, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(ctx, query, post.Title, post.Content, post.Country, time.Now()).Scan(&post.ID)
	if err != nil {
		return fmt.Errorf("error creating blog post: %w", err)
	}

	return nil
}

func (r *BlogPostRepository) GetRandom(ctx context.Context, limit int, offset int) ([]BlogPost, error) {
	query := `
		SELECT id, title, content, country, created_at
		FROM blogposts
		ORDER BY RANDOM()
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching random posts: %w", err)
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Country, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *BlogPostRepository) GetAll(ctx context.Context, limit int, offset int) ([]BlogPost, error) {
	query := `
		SELECT id, title, content, country, created_at
		FROM blogposts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error fetching posts: %w", err)
	}
	defer rows.Close()

	var posts []BlogPost
	for rows.Next() {
		var post BlogPost
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Country, &post.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, post)
	}

	return posts, nil
}

func (r *BlogPostRepository) Count(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM blogposts`
	err := r.db.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting posts: %w", err)
	}
	return count, nil
}

func (r *BlogPostRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM blogposts WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("post not found")
	}

	return nil
}
