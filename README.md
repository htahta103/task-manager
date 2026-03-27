# Task Manager

Monorepo for the Task Manager app.

## Repo layout

- `backend/`: Go HTTP API
- `frontend/`: (to be added by frontend crew) Cloudflare Pages app

## Prerequisites

- Go \(1.22+\)
- Docker + Docker Compose

## Local development

### Option A: Docker Compose (recommended)

```bash
cp .env.example .env
make dev
```

API should be available at `http://localhost:8080`.

### Option B: Run backend directly

```bash
cd backend
go test ./...
go run .
```

## Common commands

```bash
make build
make test
make lint
make typecheck
```

## Environment variables

See `.env.example` for the full list.
