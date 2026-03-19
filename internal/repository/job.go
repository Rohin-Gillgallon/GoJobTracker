package repository

import (
	"database/sql"
	"fmt"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/jmoiron/sqlx"
)

type JobRepository struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) Create(job *models.Job) error {
	query := `
		INSERT INTO jobs (user_id, company, role, status, notes, applied_date)
		VALUES (:user_id, :company, :role, :status, :notes, :applied_date)
		RETURNING id, created_at, updated_at`

	rows, err := r.db.NamedQuery(query, job)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.StructScan(job)
	}
	return nil
}

func (r *JobRepository) GetByID(id, userID string) (*models.Job, error) {
	job := &models.Job{}
	query := `SELECT * FROM jobs WHERE id = $1 AND user_id = $2`

	if err := r.db.Get(job, query, id, userID); err != nil {
		return nil, err
	}
	return job, nil
}

func (r *JobRepository) GetAll(userID string, filter models.JobsFilter) ([]models.Job, error) {
	jobs := []models.Job{}

	query := `SELECT * FROM jobs WHERE user_id = $1`
	args := []interface{}{userID}
	argCount := 2

	if filter.Status != "" {
		query += fmt.Sprintf(" AND status = $%d", argCount)
		args = append(args, filter.Status)
		argCount++
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, filter.Limit, (filter.Page-1)*filter.Limit)

	if err := r.db.Select(&jobs, query, args...); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) Update(job *models.Job) error {
	query := `
		UPDATE jobs
		SET company = :company,
		    role = :role,
		    status = :status,
		    notes = :notes,
		    applied_date = :applied_date,
		    updated_at = NOW()
		WHERE id = :id AND user_id = :user_id
		RETURNING updated_at`

	rows, err := r.db.NamedQuery(query, job)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return rows.Scan(&job.UpdatedAt)
	}
	return nil
}

func (r *JobRepository) Delete(id, userID string) error {
	query := `DELETE FROM jobs WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
