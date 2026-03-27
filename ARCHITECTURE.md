# Task Manager — Architecture (v1)

This document is produced by `tm-mol-k4md` (“Design system architecture”).

## Goals

- Single-user task manager with **REST API**, **web UI**, and a **local CLI** client.
- Simple domain model focused on `Task`.
- Consistent JSON API with predictable error handling.
- Deploy frontend to **Cloudflare Pages** and backend to **Fly.io** (Go + Postgres).

Non-goals (v1): multi-user auth, collaboration, tags, subtasks, ordering.

## Component diagram

```
                    +---------------------------+
                    |   Cloudflare Pages        |
                    |   React (Vite) SPA        |
                    +-------------+-------------+
                                  |
                                  | HTTPS (JSON)
                                  v
+-------------------+   HTTPS   +---------------------------+      TCP
| Local CLI client  +----------> | Fly.io App: task-manager  +------------+
| (Node or Go)      |           | Go HTTP API               |            |
+-------------------+           | - validation + errors      |            |
                                | - CORS                     |            |
                                | - persistence (sqlc/pgx)   |            |
                                +-------------+--------------+            |
                                              |                           |
                                              v                           |
                                +---------------------------+             |
                                | Fly Postgres (or managed) | <-----------+
                                | tasks table + enums        |
                                +---------------------------+
```

## Domain model

### Entity: `Task`

- **id**: UUID (server-generated)
- **title**: string, required, max 255
- **description**: string, optional
- **status**: enum `pending | in_progress | done`
- **priority**: enum `low | medium | high`
- **due_date**: date, optional
- **created_at**: timestamptz, server-managed
- **updated_at**: timestamptz, server-managed (auto-updated)

No relationships in v1 (no tags/subtasks).

## API surface overview

Base path: `/api`

- `GET /health`
- `GET /api/tasks?status=&priority=&search=`
- `POST /api/tasks`
- `GET /api/tasks/{id}`
- `PATCH /api/tasks/{id}`
- `DELETE /api/tasks/{id}`
- `DELETE /api/tasks/clear/done`

Canonical shapes are defined in `api-spec.yaml`.

### Auth model

- **v1 has no user authentication** (single-tenant app scope).
- All API endpoints are currently unauthenticated and rely on deployment-level network controls.
- Future migration path: add bearer token auth (JWT or opaque session token) without changing resource paths.

### Error model (high level)

- **400**: invalid UUID, validation errors, malformed JSON
- **404**: task not found
- **500**: unexpected server errors (don’t leak internals)

All responses are JSON with `Content-Type: application/json`.

## Technology choices (rationale)

The PRD draft references Node/Express + Render + Supabase; however this rig’s operational baseline is:
- **Backend**: Go (simple, fast, strong stdlib HTTP, easy Fly deploy)
- **DB**: PostgreSQL (Fly Postgres in staging/prod)
- **Frontend**: React + Vite + Tailwind (Cloudflare Pages friendly)
- **CLI**: Node (portable) or Go (single binary). Either way, it calls the deployed REST API.

Backend implementation guidance:
- Use `net/http` + a small router (or stdlib patterns) and a dedicated package for validation + error mapping.
- Persist with `pgx` (or `database/sql`) and generate queries with `sqlc` (or small hand-written queries initially).

## Deployment topology

### Staging / production

- **Frontend**: Cloudflare Pages project
  - Build: `npm run build` → `dist/`
  - Runtime config: `VITE_API_URL` points to the Fly backend public URL
- **Backend**: Fly.io app (`PORT=8080`)
  - Env:
    - `DATABASE_URL` (Postgres)
    - `CORS_ORIGIN` (Cloudflare Pages URL + local dev origin)
    - `PORT` (usually 8080)
- **Database**: Fly Postgres 17 (or compatible managed Postgres)

### Local development

- Backend: `go test ./...` then `go run ./...` (or `air`)
- DB: local Postgres (Docker) or `fly proxy` to staging DB for manual testing
- Frontend: Vite dev server uses `VITE_API_URL=http://localhost:8080/api`

## Key flows (mapped to stories)

- **List tasks**: filters by `status`, `priority`, and `search` (title substring)
- **Create**: validates `title`; returns `201` with full task JSON
- **Edit**: partial update via PATCH; reflects immediately in UI
- **Delete**: confirmation in UI; API returns `{message,id}`
- **Clear done**: API returns deleted count
- **CLI**: `list`, `add`, `done <id>`; handles network errors and invalid IDs gracefully

