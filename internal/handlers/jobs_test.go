package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/auth"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/database"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/handlers"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJobHandler(t *testing.T) (*handlers.JobHandler, string) {
	db := database.New("postgres://postgres:postgres@localhost:5433/jobtracker_test?sslmode=disable")
	t.Cleanup(func() { db.Close() })

	jobRepo := repository.NewJobRepository(db)
	userRepo := repository.NewUserRepository(db)

	email := fmt.Sprintf("jobtest_%d@example.com", time.Now().UnixNano())
	user := &models.User{
		Email:    email,
		Password: "hashedpassword",
	}
	err := userRepo.Create(user)
	require.NoError(t, err)
	require.NotEmpty(t, user.ID)

	return handlers.NewJobHandler(jobRepo), user.ID
}

func TestCreateJob_Success(t *testing.T) {
	handler, userID := setupJobHandler(t)

	now := time.Now()
	body, _ := json.Marshal(models.CreateJobRequest{
		Company:     "Acme Corp",
		Role:        "Software Engineer",
		Status:      models.StatusApplied,
		AppliedDate: &now,
	})

	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateJob(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var job models.Job
	require.NoError(t, json.NewDecoder(w.Body).Decode(&job))
	assert.Equal(t, "Acme Corp", job.Company)
	assert.Equal(t, models.StatusApplied, job.Status)
}

func TestCreateJob_MissingFields(t *testing.T) {
	handler, userID := setupJobHandler(t)

	body, _ := json.Marshal(models.CreateJobRequest{})

	req := httptest.NewRequest(http.MethodPost, "/jobs", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.CreateJob(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetAll_Success(t *testing.T) {
	handler, userID := setupJobHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.GetAllJobs(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteJob_NotFound(t *testing.T) {
	handler, userID := setupJobHandler(t)

	req := httptest.NewRequest(http.MethodDelete, "/jobs/00000000-0000-0000-0000-000000000000", nil)
	ctx := context.WithValue(req.Context(), auth.UserIDKey, userID)
	req = req.WithContext(ctx)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "00000000-0000-0000-0000-000000000000")
	ctx = context.WithValue(req.Context(), chi.RouteCtxKey, rctx)
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()
	handler.DeleteJob(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
