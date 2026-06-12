# GoRound 🐹

Plataforma de quiz para aprender programação com questões de código real, syntax highlight e sistema de pontuação gamificado.

## Stack

| Camada | Tecnologia |
|---|---|
| Frontend | SvelteKit 2 + Tailwind CSS v4 |
| Syntax Highlight | Shiki |
| Backend | Go 1.22 + Chi |
| Banco de dados | PostgreSQL 16 |
| Auth | JWT + bcrypt |
| Infra | Docker + Docker Compose |

## Início rápido

### Pré-requisitos

- Docker e Docker Compose
- Go 1.22+ (para desenvolvimento local do backend)
- Node.js 20+ (para desenvolvimento local do frontend)

### 1. Subir o backend + banco

```bash
docker-compose up -d
```

O banco sobe, as migrations rodam automaticamente e o seed de questões é aplicado.

API disponível em: `http://localhost:8080`

### 2. Subir o frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

Frontend disponível em: `http://localhost:5173`

## Estrutura do projeto

```
quiz-dev/
├── docker-compose.yml
├── backend/
│   ├── cmd/api/main.go          # Entrypoint, roteamento Chi
│   ├── internal/
│   │   ├── auth/auth.go         # JWT, registro, login, middleware
│   │   ├── quiz/quiz.go         # Questões, sessões, respostas
│   │   ├── score/score.go       # Pontuação, leaderboard
│   │   └── db/
│   │       ├── db.go            # Conexão + runner de migrations
│   │       └── migrations/
│   │           ├── 001_initial_schema.sql
│   │           └── 002_seed_go_questions.sql
│   └── Dockerfile
└── frontend/
    └── src/
        ├── lib/
        │   ├── api/client.ts    # Funções de fetch tipadas
        │   ├── stores/
        │   │   ├── auth.ts      # Store de autenticação (persistido)
        │   │   └── quiz.ts      # Estado da sessão de quiz
        │   └── components/
        │       └── CodeBlock.svelte  # Syntax highlight com Shiki
        └── routes/
            ├── +page.svelte          # Landing / seleção de nível
            ├── login/+page.svelte
            ├── register/+page.svelte
            ├── quiz/[level]/+page.svelte  # Tela principal do quiz
            ├── results/[id]/+page.svelte  # Resultado da sessão
            └── profile/+page.svelte       # Perfil e pontuação
```

## API Reference

### Auth (públicas)

```
POST /auth/register  { email, username, password }
POST /auth/login     { email, password }
```

### Quiz (requer Bearer token)

```
GET  /auth/me
GET  /questions?language=go&level=beginner
POST /sessions                    { language, level }
POST /sessions/:id/answer         { question_id, option_id }
GET  /sessions/:id/result
GET  /scores/me
GET  /scores/leaderboard?language=go&level=beginner
```

## Sistema de pontuação

| Evento | Pontos |
|---|---|
| Resposta correta | +10 |
| Streak ≥ 3 acertos | +5 bônus |
| Streak ≥ 5 acertos | +15 bônus |
| Resposta errada | streak zerado |

## Roadmap

- [x] **Fase 1** — Fundação: Docker, Go API, auth JWT, SvelteKit, Shiki, seed de questões
- [ ] **Fase 2** — Core do Quiz: fluxo completo, timer por questão, feedback visual
- [ ] **Fase 3** — Gamificação: progresso detalhado, animações, leaderboard
- [ ] **Fase 4** — Expansão: login social (Google), outras linguagens, admin panel, PWA

## Variáveis de ambiente

### Backend (via Docker Compose)
```
DATABASE_URL=postgres://quizdev:quizdev@db:5432/quizdev?sslmode=disable
JWT_SECRET=change-me-in-production
PORT=8080
```

### Frontend
```
VITE_API_URL=http://localhost:8080
```
