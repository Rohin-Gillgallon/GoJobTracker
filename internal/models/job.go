package models

import "time"

type JobStatus string

const (
	StatusApplied   JobStatus = "applied"
	StatusInterview JobStatus = "interview"
	StatusOffer     JobStatus = "offer"
	StatusRejected  JobStatus = "rejected"
)

type Job struct {
	ID          string    `db:"id"           json:"id"`
	UserID      string    `db:"user_id"      json:"user_id"`
	Company     string    `db:"company"      json:"company"`
	Role        string    `db:"role"         json:"role"`
	Status      JobStatus `db:"status"       json:"status"`
	Notes       string    `db:"notes"        json:"notes"`
	AppliedDate time.Time `db:"applied_date" json:"applied_date"`
	CreatedAt   time.Time `db:"created_at"   json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"   json:"updated_at"`
}

type CreateJobRequest struct {
	Company     string    `json:"company"`
	Role        string    `json:"role"`
	Status      JobStatus `json:"status"`
	Notes       string    `json:"notes"`
	AppliedDate time.Time `json:"applied_date"`
}

type UpdateJobRequest struct {
	Company     string    `json:"company"`
	Role        string    `json:"role"`
	Status      JobStatus `json:"status"`
	Notes       string    `json:"notes"`
	AppliedDate time.Time `json:"applied_date"`
}

type JobsFilter struct {
	Status JobStatus `json:"status"`
	Page   int       `json:"page"`
	Limit  int       `json:"limit"`
}
