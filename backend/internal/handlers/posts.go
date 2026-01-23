package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jimmyrecce/ghost-blog/internal/middleware"
	"github.com/jimmyrecce/ghost-blog/internal/models"
)

type PostHandler struct {
	repo *models.BlogPostRepository
}

func NewPostHandler(repo *models.BlogPostRepository) *PostHandler {
	return &PostHandler{repo: repo}
}

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate input
	if req.Title == "" || req.Content == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Title and content are required"})
		return
	}

	if len(req.Title) > 255 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Title must be less than 255 characters"})
		return
	}

	// Get client IP and country
	ipAddress := middleware.GetClientIP(r)
	country := middleware.GetCountryFromIP(ipAddress)
	log.Printf("Post submission - IP: %s, Country: %s", ipAddress, country)

	post := &models.BlogPost{
		Title:   req.Title,
		Content: req.Content,
		Country: country,
	}

	if err := h.repo.Create(r.Context(), post); err != nil {
		log.Printf("Error creating post: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create post"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (h *PostHandler) GetRandom(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 20 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	posts, err := h.repo.GetRandom(r.Context(), limit, offset)
	if err != nil {
		log.Printf("Error fetching random posts: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to fetch posts"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(posts)
}
