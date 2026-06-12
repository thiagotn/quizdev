# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# GoRound

Plataforma de quiz gamificado para aprender Go. Tagline: *Survive the code. Round by round.*

## Quick Start

**Backend + Database (Docker):**
```bash
docker-compose up -d
```
- API runs on `http://localhost:8080`
- Database: PostgreSQL 16 on port 5432 (user/pass: `quizdev`/`quizdev`)
- Migrations run automatically on startup

**Frontend (local Node.js):**
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```
- Runs on `http://localhost:5173`
- Set `VITE_API_URL=http://localhost:8080` in `.env`

## Stack
- **Frontend**: SvelteKit 2 + Tailwind CSS v4 + Shiki (syntax highlighting)
- **Backend**: Go 1.22 + Chi router + pgx
- **DB**: PostgreSQL 16
- **Auth**: JWT (RS256) + bcrypt
- **Infra**: Docker Compose

## Architecture

### Backend (Go)
Organized by domain in `backend/internal/`:
- **`auth/`** — Registration, login, JWT middleware, user context
- **`quiz/`** — Questions, sessions, answer submission, results
- **`score/`** — Scoring logic, leaderboard queries
- **`db/`** — Connection pooling, migration runner, embedded SQL files

Entry point: `backend/cmd/api/main.go`
- Initializes Chi router with CORS, logging, request ID middleware
- Mounts auth routes (public), quiz/score routes (protected by JWT middleware)
- Runs database migrations on startup using `go:embed`

**API design:**
- Handlers take `*pgxpool.Pool` as dependency for DB access
- All responses are JSON; errors use standard HTTP status codes
- Protected routes check JWT in `Authorization: Bearer <token>` header via `auth.Middleware`

### Frontend (SvelteKit)
- **`src/lib/api/`** — Client functions for API calls (typed fetch wrappers)
- **`src/lib/stores/`** — Reactive stores for auth state and quiz session data
- **`src/routes/`** — Page components; `[level]` and `[id]` are dynamic segments
- Components use Shiki for code syntax highlighting

## Development Commands

### Backend
```bash
# Build the Docker image
docker build -t goround-api ./backend

# Run tests (if present)
cd backend && go test ./...

# Format code
cd backend && go fmt ./...

# Run linter (if configured)
cd backend && golangci-lint run ./...

# Access database from CLI (when container is running)
docker exec -it goround-db psql -U quizdev -d quizdev
```

### Frontend
```bash
# Type-check (Svelte + TypeScript)
npm run check

# Lint + format
npm run lint

# Build for production
npm run build

# Run production preview locally
npm run preview
```

### Database
**Add a migration:**
1. Create `backend/internal/db/migrations/NNN_description.sql` (zero-padded 3-digit number)
2. Migration runs automatically on next `docker-compose up` or container restart
3. Filename is recorded in `schema_migrations` table to prevent re-running

Migrations are embedded in the binary and executed in filename order by `db.RunMigrations()`.

## Design System

### Naming System
- Iniciante → Round I
- Intermediário → Round II
- Avançado → Round III

### Color Palette
- **Arena Red**: #E63946 (primary accent, CTAs)
- **Void Black**: #111111 (background)
- **Surface**: #1A1A1A (cards/components)
- **Frost White**: #F1FAEE (primary text)
- **Go Blue**: #00ACD7 (Go identity)
- **Code Green**: #4ADE80 (correct answers)
- **Streak Amber**: #F59E0B (streaks)

## Project Conventions

- **Migrations**: Numbered `001_`, `002_`, `003_` in `backend/internal/db/migrations/`
- **Components**: Stored in `frontend/src/lib/components/`
- **Stores**: In `frontend/src/lib/stores/` (Svelte writable/readable stores)
- **API client**: Typed fetch wrappers in `frontend/src/lib/api/`
- **Error handling**: Backend returns standard HTTP status codes; frontend checks response status

## Current Phase: Phase 4
- [ ] OAuth2 social login (Google)
- [ ] Admin panel for managing questions (`/admin`)
- [ ] PWA support (manifest.json + service worker)
- [ ] Multi-language support (TypeScript, Rust, etc.)
