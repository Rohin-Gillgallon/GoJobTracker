package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/auth"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
	secret   string
}

func NewAuthHandler(userRepo *repository.UserRepository, secret string) *AuthHandler {
	return &AuthHandler{userRepo: userRepo, secret: secret}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := h.userRepo.Create(user); err != nil {
		log.Printf("Register error: %v", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	tokens, err := auth.GenerateTokenPair(user.ID, h.secret)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(models.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	tokens, err := auth.GenerateTokenPair(user.ID, h.secret)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(models.AuthResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
