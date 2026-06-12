-- 001_initial_schema.sql

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Users
CREATE TABLE users (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email       TEXT UNIQUE NOT NULL,
    username    TEXT UNIQUE NOT NULL,
    password_hash TEXT,                      -- nullable for future OAuth users
    provider    TEXT NOT NULL DEFAULT 'local', -- 'local' | 'google' | 'github'
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Questions
CREATE TABLE questions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    language    TEXT NOT NULL DEFAULT 'go',
    level       TEXT NOT NULL CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    title       TEXT NOT NULL,
    explanation TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Answer options (code snippets)
CREATE TABLE question_options (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    code_snippet TEXT NOT NULL,
    is_correct  BOOLEAN NOT NULL DEFAULT FALSE,
    display_order INT NOT NULL DEFAULT 0
);

-- Quiz sessions
CREATE TABLE quiz_sessions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language    TEXT NOT NULL DEFAULT 'go',
    level       TEXT NOT NULL CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    total_questions INT NOT NULL DEFAULT 0,
    correct_answers INT NOT NULL DEFAULT 0,
    score       INT NOT NULL DEFAULT 0,
    started_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ
);

-- Individual answers within a session
CREATE TABLE user_answers (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id  UUID NOT NULL REFERENCES quiz_sessions(id) ON DELETE CASCADE,
    question_id UUID NOT NULL REFERENCES questions(id),
    option_id   UUID NOT NULL REFERENCES question_options(id),
    is_correct  BOOLEAN NOT NULL,
    streak_at_answer INT NOT NULL DEFAULT 0,
    answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Aggregated scores per user/language/level
CREATE TABLE scores (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    language     TEXT NOT NULL DEFAULT 'go',
    level        TEXT NOT NULL,
    total_points INT NOT NULL DEFAULT 0,
    best_streak  INT NOT NULL DEFAULT 0,
    sessions_played INT NOT NULL DEFAULT 0,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, language, level)
);

-- Indexes
CREATE INDEX idx_questions_language_level ON questions(language, level);
CREATE INDEX idx_quiz_sessions_user_id ON quiz_sessions(user_id);
CREATE INDEX idx_scores_user_id ON scores(user_id);
CREATE INDEX idx_scores_total_points ON scores(total_points DESC);
