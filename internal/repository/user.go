package repository

import (
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password)
		VALUES (:email, :password)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQuery(query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.StructScan(user)
	}
	return nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT * FROM users WHERE email = $1`

	if err := r.db.Get(user, query, email); err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT * FROM users WHERE id = $1`

	if err := r.db.Get(user, query, id); err != nil {
		return nil, err
	}
	return user, nil
}
