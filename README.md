# Task Manager

Monorepo for the Task Manager app.

## Repo layout

- `backend/`: Go HTTP API
- `frontend/`: React + Vite Cloudflare Pages app

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

## Frontend deploy (Cloudflare Pages)

From `frontend/`:

```bash
export CLOUDFLARE_PAGES_PROJECT=taskmanager-staging
export VITE_API_URL=https://taskmanager-api-staging.fly.dev/api
npm ci
npm run deploy:staging
```

`deploy:staging` deploys `frontend/dist` to Cloudflare Pages branch `main` and defaults to account `353a69adadf77371f073ef6bb0a66f4c`.

## Environment variables

See `.env.example` for the full list.
