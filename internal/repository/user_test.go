package repository_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/database"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserRepo(t *testing.T) *repository.UserRepository {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5433/jobtracker_test?sslmode=disable"
	}
	db := database.New(dbURL)
	t.Cleanup(func() { db.Close() })
	return repository.NewUserRepository(db)
}

func TestUserRepo_Create(t *testing.T) {
	repo := setupUserRepo(t)

	user := &models.User{
		Email:    fmt.Sprintf("repo_%d@test.com", time.Now().UnixNano()),
		Password: "hashedpassword",
	}

	err := repo.Create(user)
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
	assert.NotZero(t, user.CreatedAt)
}

func TestUserRepo_GetByEmail(t *testing.T) {
	repo := setupUserRepo(t)

	email := fmt.Sprintf("getbyemail_%d@test.com", time.Now().UnixNano())
	user := &models.User{Email: email, Password: "hashedpassword"}
	require.NoError(t, repo.Create(user))

	found, err := repo.GetByEmail(email)
	require.NoError(t, err)
	assert.Equal(t, email, found.Email)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepo_GetByEmail_NotFound(t *testing.T) {
	repo := setupUserRepo(t)

	_, err := repo.GetByEmail("nonexistent@test.com")
	assert.Error(t, err)
}

func TestUserRepo_GetByID(t *testing.T) {
	repo := setupUserRepo(t)

	user := &models.User{
		Email:    fmt.Sprintf("getbyid_%d@test.com", time.Now().UnixNano()),
		Password: "hashedpassword",
	}
	require.NoError(t, repo.Create(user))

	found, err := repo.GetByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	repo := setupUserRepo(t)

	_, err := repo.GetByID("00000000-0000-0000-0000-000000000000")
	assert.Error(t, err)
}
