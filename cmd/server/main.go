package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Rohin-Gillgallon/GoJobTracker/internal/auth"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/config"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/database"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/handlers"
	"github.com/Rohin-Gillgallon/GoJobTracker/internal/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Load()

	db := database.New(cfg.DatabaseURL)
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)

	authHandler := handlers.NewAuthHandler(userRepo, cfg.JWTSecret)
	jobHandler := handlers.NewJobHandler(jobRepo)

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	router.Post("/auth/register", authHandler.Register)
	router.Post("/auth/login", authHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(auth.Middleware(cfg.JWTSecret))
		r.Get("/jobs", jobHandler.GetAllJobs)
		r.Post("/jobs", jobHandler.CreateJob)
		r.Get("/jobs/{id}", jobHandler.GetJobByID)
		r.Put("/jobs/{id}", jobHandler.UpdateJob)
		r.Delete("/jobs/{id}", jobHandler.DeleteJob)
	})

	address := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", address)

	if err := http.ListenAndServe(address, router); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
