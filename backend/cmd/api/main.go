package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	"github.com/goround/api/internal/auth"
	"github.com/goround/api/internal/db"
	"github.com/goround/api/internal/quiz"
	"github.com/goround/api/internal/score"
)

func main() {
	_ = godotenv.Load()

	ctx := context.Background()

	pool, err := db.Connect(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	authHandler := auth.NewHandler(pool)
	quizHandler := quiz.NewHandler(pool)
	scoreHandler := score.NewHandler(pool)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://quizdev.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Auth routes (public)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(auth.Middleware)

		r.Get("/auth/me", authHandler.Me)

		// Questions
		r.Get("/questions", quizHandler.ListQuestions)
		r.Get("/questions/{id}", quizHandler.GetQuestion)

		// Sessions
		r.Post("/sessions", quizHandler.CreateSession)
		r.Post("/sessions/{id}/answer", quizHandler.SubmitAnswer)
		r.Get("/sessions/{id}/result", quizHandler.GetResult)

		// Scores
		r.Get("/scores/me", scoreHandler.MyScores)
		r.Get("/scores/leaderboard", scoreHandler.Leaderboard)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
