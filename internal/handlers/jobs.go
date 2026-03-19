package handlers

import (
	"database/sql"
	"encoding/json"

	"net/http"
	"strconv"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/auth"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/models"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	jobRepo *repository.JobRepository
}

func NewJobHandler(jobRepo *repository.JobRepository) *JobHandler {
	return &JobHandler{jobRepo: jobRepo}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	var req models.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Company == "" || req.Role == "" {
		http.Error(w, "missing required fields: company or role", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		req.Status = models.StatusApplied
	}

	job := &models.Job{
		UserID:      userID,
		Company:     req.Company,
		Role:        req.Role,
		Status:      req.Status,
		Notes:       req.Notes,
		AppliedDate: req.AppliedDate,
	}

	if err := h.jobRepo.Create(job); err != nil {
		http.Error(w, "failed to create job in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)

	filter := models.JobsFilter{
		Status: models.JobStatus(r.URL.Query().Get("status")),
		Page:   1,
		Limit:  10,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if limit := r.URL.Query().Get("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filter.Limit = l
		}
	}

	jobs, err := h.jobRepo.GetAll(userID, filter)
	if err != nil {
		http.Error(w, "failed to fetch jobs from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

func (h *JobHandler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	id := r.URL.Query().Get("id")
	chi.URLParam(r, "id")

	job, err := h.jobRepo.GetByID(id, userID)
	if err != nil {
		http.Error(w, "failed to fetch job from database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	id := chi.URLParam(r, "id")

	job, err := h.jobRepo.GetByID(id, userID)
	if err != nil {
		http.Error(w, "failed to fetch job for update", http.StatusInternalServerError)
		return
	}

	var req models.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid update request body", http.StatusBadRequest)
		return
	}

	job.Company = req.Company
	job.Role = req.Role
	job.Status = req.Status
	job.Notes = req.Notes
	job.AppliedDate = req.AppliedDate

	if err := h.jobRepo.Update(job); err != nil {
		http.Error(w, "failed to update job in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(job)
}

func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserID(r)
	id := chi.URLParam(r, "id")

	if err := h.jobRepo.Delete(id, userID); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
