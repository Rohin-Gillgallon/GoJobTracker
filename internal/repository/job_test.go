package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/database"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJobRepo(t *testing.T) (*repository.JobRepository, string) {
	db := database.New("postgres://postgres:postgres@localhost:5433/jobtracker_test?sslmode=disable")
	t.Cleanup(func() { db.Close() })

	userRepo := repository.NewUserRepository(db)
	user := &models.User{
		Email:    fmt.Sprintf("jobrepo_%d@test.com", time.Now().UnixNano()),
		Password: "hashedpassword",
	}
	require.NoError(t, userRepo.Create(user))

	return repository.NewJobRepository(db), user.ID
}

func createTestJob(t *testing.T, repo *repository.JobRepository, userID string) *models.Job {
	now := time.Now()
	job := &models.Job{
		UserID:      userID,
		Company:     "Test Corp",
		Role:        "Engineer",
		Status:      models.StatusApplied,
		Notes:       "Test notes",
		AppliedDate: &now,
	}
	require.NoError(t, repo.Create(job))
	return job
}

func TestJobRepo_Create(t *testing.T) {
	repo, userID := setupJobRepo(t)

	now := time.Now()
	job := &models.Job{
		UserID:      userID,
		Company:     "Acme",
		Role:        "Developer",
		Status:      models.StatusApplied,
		AppliedDate: &now,
	}

	err := repo.Create(job)
	require.NoError(t, err)
	assert.NotEmpty(t, job.ID)
}

func TestJobRepo_GetByID(t *testing.T) {
	repo, userID := setupJobRepo(t)
	job := createTestJob(t, repo, userID)

	found, err := repo.GetByID(job.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, job.ID, found.ID)
	assert.Equal(t, "Test Corp", found.Company)
}

func TestJobRepo_GetByID_NotFound(t *testing.T) {
	repo, userID := setupJobRepo(t)

	_, err := repo.GetByID("00000000-0000-0000-0000-000000000000", userID)
	assert.Error(t, err)
}

func TestJobRepo_GetAll(t *testing.T) {
	repo, userID := setupJobRepo(t)
	createTestJob(t, repo, userID)
	createTestJob(t, repo, userID)

	jobs, err := repo.GetAll(userID, models.JobsFilter{Page: 1, Limit: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(jobs), 2)
}

func TestJobRepo_GetAll_FilterByStatus(t *testing.T) {
	repo, userID := setupJobRepo(t)
	createTestJob(t, repo, userID)

	jobs, err := repo.GetAll(userID, models.JobsFilter{Status: models.StatusApplied, Page: 1, Limit: 10})
	require.NoError(t, err)
	for _, j := range jobs {
		assert.Equal(t, models.StatusApplied, j.Status)
	}
}

func TestJobRepo_Update(t *testing.T) {
	repo, userID := setupJobRepo(t)
	job := createTestJob(t, repo, userID)

	job.Company = "Updated Corp"
	job.Status = models.StatusInterview

	err := repo.Update(job)
	require.NoError(t, err)

	found, err := repo.GetByID(job.ID, userID)
	require.NoError(t, err)
	assert.Equal(t, "Updated Corp", found.Company)
	assert.Equal(t, models.StatusInterview, found.Status)
}

func TestJobRepo_Delete(t *testing.T) {
	repo, userID := setupJobRepo(t)
	job := createTestJob(t, repo, userID)

	err := repo.Delete(job.ID, userID)
	require.NoError(t, err)

	_, err = repo.GetByID(job.ID, userID)
	assert.Error(t, err)
}
