package score

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/goround/api/internal/auth"
)

// ─── Models ───────────────────────────────────────────────────────────────────

type Score struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Language       string    `json:"language"`
	Level          string    `json:"level"`
	TotalPoints    int       `json:"total_points"`
	BestStreak     int       `json:"best_streak"`
	BestAccuracy   int       `json:"best_accuracy"`
	SessionsPlayed int       `json:"sessions_played"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ProfileUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SessionSummary struct {
	ID             string     `json:"id"`
	Level          string     `json:"level"`
	Score          int        `json:"score"`
	CorrectAnswers int        `json:"correct_answers"`
	TotalQuestions int        `json:"total_questions"`
	Accuracy       int        `json:"accuracy"`
	FinishedAt     *time.Time `json:"finished_at"`
}

type ProfileResponse struct {
	User     ProfileUser      `json:"user"`
	Scores   []Score          `json:"scores"`
	History  []SessionSummary `json:"history"`
	GlobalRank int            `json:"global_rank"`
	TotalPoints int           `json:"total_points"`
}

type LeaderboardEntry struct {
	Rank           int    `json:"rank"`
	Username       string `json:"username"`
	TotalPoints    int    `json:"total_points"`
	SessionsPlayed int    `json:"sessions_played"`
	BestStreak     int    `json:"best_streak"`
	IsCurrentUser  bool   `json:"is_current_user"`
}

// ─── Handler ──────────────────────────────────────────────────────────────────

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// GET /scores/me — full profile with history and rank
func (h *Handler) MyScores(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var user ProfileUser
	h.db.QueryRow(r.Context(),
		`SELECT id, username, email FROM users WHERE id = $1`, claims.UserID,
	).Scan(&user.ID, &user.Username, &user.Email)

	// Per-level scores
	rows, err := h.db.Query(r.Context(),
		`SELECT id, user_id, language, level, total_points, best_streak,
		        COALESCE(best_accuracy, 0), sessions_played, updated_at
		 FROM scores WHERE user_id = $1 ORDER BY level`,
		claims.UserID,
	)
	if err != nil {
		writeError(w, "failed to fetch scores", http.StatusInternalServerError)
		return
	}
	scores := []Score{}
	for rows.Next() {
		var s Score
		rows.Scan(&s.ID, &s.UserID, &s.Language, &s.Level, &s.TotalPoints,
			&s.BestStreak, &s.BestAccuracy, &s.SessionsPlayed, &s.UpdatedAt)
		scores = append(scores, s)
	}
	rows.Close()

	// Recent session history (last 20 finished)
	histRows, err := h.db.Query(r.Context(),
		`SELECT id, level, score, correct_answers, total_questions, finished_at
		 FROM quiz_sessions
		 WHERE user_id = $1 AND finished_at IS NOT NULL
		 ORDER BY finished_at DESC
		 LIMIT 20`,
		claims.UserID,
	)
	history := []SessionSummary{}
	if err == nil {
		for histRows.Next() {
			var s SessionSummary
			histRows.Scan(&s.ID, &s.Level, &s.Score, &s.CorrectAnswers, &s.TotalQuestions, &s.FinishedAt)
			if s.TotalQuestions > 0 {
				s.Accuracy = (s.CorrectAnswers * 100) / s.TotalQuestions
			}
			history = append(history, s)
		}
		histRows.Close()
	}

	// Global rank and total points from view
	var globalRank, totalPoints int
	h.db.QueryRow(r.Context(),
		`SELECT rank, total_points FROM leaderboard_global WHERE user_id = $1`,
		claims.UserID,
	).Scan(&globalRank, &totalPoints)

	writeJSON(w, ProfileResponse{
		User:        user,
		Scores:      scores,
		History:     history,
		GlobalRank:  globalRank,
		TotalPoints: totalPoints,
	}, http.StatusOK)
}

// GET /scores/leaderboard?language=go&level=beginner
func (h *Handler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	claims, _ := auth.ClaimsFromContext(r.Context())
	level := r.URL.Query().Get("level")

	var rows interface{ Next() bool; Scan(...any) error; Close() }
	var err error

	if level != "" {
		// Per-level leaderboard
		rows, err = h.db.Query(r.Context(),
			`SELECT
				RANK() OVER (ORDER BY s.total_points DESC) AS rank,
				u.username,
				s.total_points,
				s.sessions_played,
				s.best_streak,
				u.id
			 FROM scores s
			 JOIN users u ON u.id = s.user_id
			 WHERE s.language = 'go' AND s.level = $1
			 ORDER BY s.total_points DESC
			 LIMIT 20`,
			level,
		)
	} else {
		// Global leaderboard
		rows, err = h.db.Query(r.Context(),
			`SELECT rank, username, total_points, sessions_played, best_streak, user_id
			 FROM leaderboard_global
			 ORDER BY rank
			 LIMIT 20`,
		)
	}
	if err != nil {
		writeError(w, "failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	entries := []LeaderboardEntry{}
	for rows.Next() {
		var e LeaderboardEntry
		var userID string
		rows.Scan(&e.Rank, &e.Username, &e.TotalPoints, &e.SessionsPlayed, &e.BestStreak, &userID)
		e.IsCurrentUser = claims != nil && userID == claims.UserID
		entries = append(entries, e)
	}

	writeJSON(w, entries, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{"error": msg}, status)
}
