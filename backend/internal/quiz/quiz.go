package quiz

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/goround/api/internal/auth"
)

// ─── Models ───────────────────────────────────────────────────────────────────

type Option struct {
	ID          string `json:"id"`
	CodeSnippet string `json:"code_snippet"`
	IsCorrect   bool   `json:"is_correct,omitempty"` // hidden on GET, shown on POST answer
	DisplayOrder int   `json:"display_order"`
}

type Question struct {
	ID          string   `json:"id"`
	Language    string   `json:"language"`
	Level       string   `json:"level"`
	Title       string   `json:"title"`
	Explanation string   `json:"explanation,omitempty"` // only shown after answer
	Options     []Option `json:"options"`
}

type Session struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	Language       string     `json:"language"`
	Level          string     `json:"level"`
	TotalQuestions int        `json:"total_questions"`
	CorrectAnswers int        `json:"correct_answers"`
	Score          int        `json:"score"`
	StartedAt      time.Time  `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	Questions      []Question `json:"questions,omitempty"`
}

type AnswerRequest struct {
	QuestionID    string `json:"question_id"`
	OptionID      string `json:"option_id"`
	TimeRemaining int    `json:"time_remaining"` // seconds left on timer (0 = timed out)
}

type AnswerResult struct {
	IsCorrect   bool     `json:"is_correct"`
	Explanation string   `json:"explanation"`
	CorrectOption Option `json:"correct_option"`
	PointsEarned int     `json:"points_earned"`
	CurrentStreak int    `json:"current_streak"`
	TotalScore   int     `json:"total_score"`
}

// ─── Handler ──────────────────────────────────────────────────────────────────

type Handler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// GET /questions?language=go&level=beginner&limit=10
func (h *Handler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	language := r.URL.Query().Get("language")
	level := r.URL.Query().Get("level")
	if language == "" {
		language = "go"
	}
	if level == "" || !isValidLevel(level) {
		writeError(w, "invalid level: must be beginner, intermediate or advanced", http.StatusBadRequest)
		return
	}

	rows, err := h.db.Query(r.Context(),
		`SELECT id, language, level, title
		 FROM questions
		 WHERE language = $1 AND level = $2
		 ORDER BY RANDOM()
		 LIMIT 10`,
		language, level,
	)
	if err != nil {
		writeError(w, "failed to fetch questions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	questions := []Question{}
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.Language, &q.Level, &q.Title); err != nil {
			continue
		}
		questions = append(questions, q)
	}

	writeJSON(w, questions, http.StatusOK)
}

// GET /questions/:id
func (h *Handler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var q Question
	err := h.db.QueryRow(r.Context(),
		`SELECT id, language, level, title FROM questions WHERE id = $1`, id,
	).Scan(&q.ID, &q.Language, &q.Level, &q.Title)
	if err != nil {
		writeError(w, "question not found", http.StatusNotFound)
		return
	}

	opts, err := h.fetchOptions(r, q.ID, false)
	if err != nil {
		writeError(w, "failed to fetch options", http.StatusInternalServerError)
		return
	}
	q.Options = opts

	writeJSON(w, q, http.StatusOK)
}

// POST /sessions
func (h *Handler) CreateSession(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var body struct {
		Language string `json:"language"`
		Level    string `json:"level"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}
	if body.Language == "" {
		body.Language = "go"
	}
	if !isValidLevel(body.Level) {
		writeError(w, "invalid level", http.StatusBadRequest)
		return
	}

	// Fetch 10 random questions
	rows, err := h.db.Query(r.Context(),
		`SELECT id, language, level, title
		 FROM questions
		 WHERE language = $1 AND level = $2
		 ORDER BY RANDOM()
		 LIMIT 10`,
		body.Language, body.Level,
	)
	if err != nil {
		writeError(w, "failed to fetch questions", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	questions := []Question{}
	for rows.Next() {
		var q Question
		if err := rows.Scan(&q.ID, &q.Language, &q.Level, &q.Title); err != nil {
			continue
		}
		opts, _ := h.fetchOptions(r, q.ID, false)
		q.Options = opts
		questions = append(questions, q)
	}

	if len(questions) == 0 {
		writeError(w, "no questions available for this level", http.StatusNotFound)
		return
	}

	var session Session
	err = h.db.QueryRow(r.Context(),
		`INSERT INTO quiz_sessions (user_id, language, level, total_questions)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, language, level, total_questions, correct_answers, score, started_at`,
		claims.UserID, body.Language, body.Level, len(questions),
	).Scan(&session.ID, &session.UserID, &session.Language, &session.Level,
		&session.TotalQuestions, &session.CorrectAnswers, &session.Score, &session.StartedAt)
	if err != nil {
		writeError(w, "failed to create session", http.StatusInternalServerError)
		return
	}

	session.Questions = questions
	writeJSON(w, session, http.StatusCreated)
}

// POST /sessions/:id/answer
func (h *Handler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := chi.URLParam(r, "id")

	var req AnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Verify session belongs to user
	var userID string
	err := h.db.QueryRow(r.Context(),
		`SELECT user_id FROM quiz_sessions WHERE id = $1`, sessionID,
	).Scan(&userID)
	if err != nil || userID != claims.UserID {
		writeError(w, "session not found", http.StatusNotFound)
		return
	}

	// Check if option is correct
	var isCorrect bool
	var explanation string
	err = h.db.QueryRow(r.Context(),
		`SELECT qo.is_correct, q.explanation
		 FROM question_options qo
		 JOIN questions q ON q.id = qo.question_id
		 WHERE qo.id = $1 AND qo.question_id = $2`,
		req.OptionID, req.QuestionID,
	).Scan(&isCorrect, &explanation)
	if err != nil {
		writeError(w, "invalid option or question", http.StatusBadRequest)
		return
	}

	// Count current streak
	var streak int
	h.db.QueryRow(r.Context(),
		`SELECT COUNT(*) FROM (
			SELECT is_correct FROM user_answers
			WHERE session_id = $1
			ORDER BY answered_at DESC
		) sub WHERE is_correct = true`,
		sessionID,
	).Scan(&streak)
	if !isCorrect {
		streak = 0
	} else {
		streak++
	}

	// Calculate points
	points := 0
	if isCorrect {
		points = 10
		// Streak bonus
		if streak >= 5 {
			points += 15
		} else if streak >= 3 {
			points += 5
		}
		// Speed bonus: answered with >10s remaining
		if req.TimeRemaining > 10 {
			points += 2
		}
	}

	// Save answer
	h.db.Exec(r.Context(),
		`INSERT INTO user_answers (session_id, question_id, option_id, is_correct, streak_at_answer)
		 VALUES ($1, $2, $3, $4, $5)`,
		sessionID, req.QuestionID, req.OptionID, isCorrect, streak,
	)

	// Update session score
	var totalScore int
	h.db.QueryRow(r.Context(),
		`UPDATE quiz_sessions
		 SET score = score + $1,
		     correct_answers = correct_answers + $2
		 WHERE id = $3
		 RETURNING score`,
		points, boolToInt(isCorrect), sessionID,
	).Scan(&totalScore)

	// Fetch correct option
	correctOption, _ := h.fetchCorrectOption(r, req.QuestionID)

	result := AnswerResult{
		IsCorrect:     isCorrect,
		Explanation:   explanation,
		CorrectOption: correctOption,
		PointsEarned:  points,
		CurrentStreak: streak,
		TotalScore:    totalScore,
	}

	writeJSON(w, result, http.StatusOK)
}

// GET /sessions/:id/result
func (h *Handler) GetResult(w http.ResponseWriter, r *http.Request) {
	claims, ok := auth.ClaimsFromContext(r.Context())
	if !ok {
		writeError(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID := chi.URLParam(r, "id")

	var session Session
	err := h.db.QueryRow(r.Context(),
		`UPDATE quiz_sessions
		 SET finished_at = COALESCE(finished_at, NOW())
		 WHERE id = $1 AND user_id = $2
		 RETURNING id, user_id, language, level, total_questions, correct_answers, score, started_at, finished_at`,
		sessionID, claims.UserID,
	).Scan(&session.ID, &session.UserID, &session.Language, &session.Level,
		&session.TotalQuestions, &session.CorrectAnswers, &session.Score,
		&session.StartedAt, &session.FinishedAt)
	if err != nil {
		writeError(w, "session not found", http.StatusNotFound)
		return
	}

	// Upsert aggregated score — also track best streak and best accuracy
	accuracy := 0
	if session.TotalQuestions > 0 {
		accuracy = (session.CorrectAnswers * 100) / session.TotalQuestions
	}
	// Get max streak from this session
	var bestStreak int
	h.db.QueryRow(r.Context(),
		`SELECT COALESCE(MAX(streak_at_answer), 0) FROM user_answers WHERE session_id = $1`,
		sessionID,
	).Scan(&bestStreak)

	h.db.Exec(r.Context(),
		`INSERT INTO scores (user_id, language, level, total_points, sessions_played, best_streak, best_accuracy)
		 VALUES ($1, $2, $3, $4, 1, $5, $6)
		 ON CONFLICT (user_id, language, level)
		 DO UPDATE SET
		   total_points    = scores.total_points + EXCLUDED.total_points,
		   sessions_played = scores.sessions_played + 1,
		   best_streak     = GREATEST(scores.best_streak, EXCLUDED.best_streak),
		   best_accuracy   = GREATEST(scores.best_accuracy, EXCLUDED.best_accuracy),
		   updated_at      = NOW()`,
		claims.UserID, session.Language, session.Level, session.Score, bestStreak, accuracy,
	)

	writeJSON(w, session, http.StatusOK)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func (h *Handler) fetchOptions(r *http.Request, questionID string, showCorrect bool) ([]Option, error) {
	rows, err := h.db.Query(r.Context(),
		`SELECT id, code_snippet, is_correct, display_order
		 FROM question_options
		 WHERE question_id = $1
		 ORDER BY display_order`,
		questionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	opts := []Option{}
	for rows.Next() {
		var o Option
		rows.Scan(&o.ID, &o.CodeSnippet, &o.IsCorrect, &o.DisplayOrder)
		if !showCorrect {
			o.IsCorrect = false
		}
		opts = append(opts, o)
	}
	return opts, nil
}

func (h *Handler) fetchCorrectOption(r *http.Request, questionID string) (Option, error) {
	var o Option
	err := h.db.QueryRow(r.Context(),
		`SELECT id, code_snippet, is_correct, display_order
		 FROM question_options
		 WHERE question_id = $1 AND is_correct = true
		 LIMIT 1`,
		questionID,
	).Scan(&o.ID, &o.CodeSnippet, &o.IsCorrect, &o.DisplayOrder)
	return o, err
}

func isValidLevel(level string) bool {
	return level == "beginner" || level == "intermediate" || level == "advanced"
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{"error": msg}, status)
}
