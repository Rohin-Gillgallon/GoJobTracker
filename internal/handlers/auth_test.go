package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/database"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/handlers"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *repository.UserRepository {
	db := database.New("postgres://postgres:postgres@localhost:5433/jobtracker_test?sslmode=disable")
	t.Cleanup(func() { db.Close() })
	return repository.NewUserRepository(db)
}

func TestRegister_Success(t *testing.T) {
	userRepo := setupTestDB(t)
	handler := handlers.NewAuthHandler(userRepo, "test-secret")

	email := fmt.Sprintf("register_%d@example.com", time.Now().UnixNano())
	body, _ := json.Marshal(models.RegisterRequest{
		Email:    email,
		Password: "password123",
	})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp models.AuthResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestRegister_MissingFields(t *testing.T) {
	userRepo := setupTestDB(t)
	handler := handlers.NewAuthHandler(userRepo, "test-secret")

	body, _ := json.Marshal(models.RegisterRequest{Email: ""})

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Register(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin_Success(t *testing.T) {
	userRepo := setupTestDB(t)
	handler := handlers.NewAuthHandler(userRepo, "test-secret")

	email := fmt.Sprintf("login_%d@example.com", time.Now().UnixNano())
	regBody, _ := json.Marshal(models.RegisterRequest{
		Email:    email,
		Password: "password123",
	})

	regReq := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regW := httptest.NewRecorder()
	handler.Register(regW, regReq)
	require.Equal(t, http.StatusCreated, regW.Code)

	body, _ := json.Marshal(models.LoginRequest{
		Email:    email,
		Password: "password123",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.AuthResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotEmpty(t, resp.AccessToken)
}

func TestLogin_WrongPassword(t *testing.T) {
	userRepo := setupTestDB(t)
	handler := handlers.NewAuthHandler(userRepo, "test-secret")

	body, _ := json.Marshal(models.LoginRequest{
		Email:    "login@example.com",
		Password: "wrongpassword",
	})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
