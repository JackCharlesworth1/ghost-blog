package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jimmyrecce/ghost-blog/internal/config"
	"github.com/jimmyrecce/ghost-blog/internal/database"
	"github.com/jimmyrecce/ghost-blog/internal/handlers"
	"github.com/jimmyrecce/ghost-blog/internal/middleware"
	"github.com/jimmyrecce/ghost-blog/internal/models"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	postRepo := models.NewBlogPostRepository(db.Pool)
	rateLimitRepo := models.NewRateLimitRepository(db.Pool)

	// Initialize handlers
	postHandler := handlers.NewPostHandler(postRepo)
	adminHandler := handlers.NewAdminHandler(postRepo)

	// Initialize middleware
	rateLimiter := middleware.NewRateLimiter(rateLimitRepo)
	adminAuth := middleware.NewAdminAuth(cfg.AdminPassword)

	// Setup router
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))
	r.Use(middleware.CORS)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Public routes
	r.Route("/api/posts", func(r chi.Router) {
		r.Get("/", postHandler.GetRandom)
		r.With(rateLimiter.Middleware).Post("/", postHandler.Create)
	})

	// Admin routes
	r.Route("/api/admin/posts", func(r chi.Router) {
		r.Use(adminAuth.Middleware)
		r.Get("/", adminHandler.GetAll)
		r.Delete("/{id}", adminHandler.Delete)
	})

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
