-- 003_session_history.sql
-- Add best_streak tracking to scores, add rank/position view

-- Update scores to also track best streak per session
ALTER TABLE scores ADD COLUMN IF NOT EXISTS best_accuracy INT NOT NULL DEFAULT 0;

-- Leaderboard view: aggregate total points across all levels per user
CREATE OR REPLACE VIEW leaderboard_global AS
SELECT
    u.id         AS user_id,
    u.username,
    SUM(s.total_points)    AS total_points,
    SUM(s.sessions_played) AS sessions_played,
    MAX(s.best_streak)     AS best_streak,
    RANK() OVER (ORDER BY SUM(s.total_points) DESC) AS rank
FROM scores s
JOIN users u ON u.id = s.user_id
GROUP BY u.id, u.username;

-- Index for faster session history queries
CREATE INDEX IF NOT EXISTS idx_quiz_sessions_user_finished
    ON quiz_sessions(user_id, finished_at DESC NULLS LAST);
