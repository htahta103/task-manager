# Architecture Requirements: Task Manager (v1.0)

This document is produced by `tm-mol-qoto` (“Parse PRD into user stories”) to unblock the next workflow step: `tm-mol-k4md` (“Design system architecture”).

## Project Summary

Build a full-stack, single-user task manager:
- REST API backend (Node.js + Express)
- React web frontend (Vite + Tailwind)
- Local CLI client (Node.js) that calls the deployed REST API

Primary domain entity: `Task`.

## Core Domain Model

### `Task` entity
- `id` (uuid, primary key)
- `title` (string, required, max 255 chars)
- `description` (text, optional)
- `status` (enum): `pending | in_progress | done`
- `priority` (enum): `low | medium | high`
- `due_date` (date, optional)
- `created_at` (timestamptz)
- `updated_at` (timestamptz)

### Relationships
- Tasks are independent (no sub-tasks, no tags, no ordering).

## API Surface (from PRD)

Base URLs:
- Production: `https://task-manager-api.onrender.com/api`
- Local dev: `http://localhost:3000/api`

Endpoints:
- `GET /health` -> `{ status: "ok", timestamp: <iso> }`
- `GET /tasks?status=&priority=&search=` -> returns `{ data: [...], count: <n> }` (or equivalent shape per implementation)
- `POST /tasks` -> creates a task; validates `title` and returns `201`
- `GET /tasks/:id` -> `404` for missing, `400` for invalid UUID
- `PATCH /tasks/:id` -> partial update; `404` for missing, `400` for invalid UUID
- `DELETE /tasks/:id` -> deletes a task; returns `{ message, id }`
- `DELETE /tasks/clear/done` -> deletes all tasks with `status=done`; returns `{ message, deleted }`

Error + validation expectations:
- Missing/invalid payloads must return `400` with a descriptive error message (see PRD examples).
- Unknown IDs must return `404` with `{ error: "Task not found" }`.
- All endpoints must return `Content-Type: application/json`.

## User Flows (UI -> API mapping)

Web UI:
1. Load dashboard: fetch list of tasks (default filters show all).
2. Add task: submit modal -> POST `/tasks` -> insert/update in the list (optimistic update or refetch).
3. Edit task: modal -> PATCH `/tasks/:id` -> update row in place.
4. Delete task: confirm -> DELETE `/tasks/:id` -> remove from list.
5. Filters + search:
   - Status filter maps to `GET /tasks?status=...`
   - Priority filter maps to `GET /tasks?priority=...`
   - Search maps to `GET /tasks?search=...` (real-time, debounced client-side).
6. Deploy: frontend must build and deploy successfully to Cloudflare Pages (workspace `htahta103`).

CLI:
1. `task list` -> GET `/tasks` with optional filtering flags (implementation-defined mapping).
2. `task add` -> POST `/tasks` -> prints the created UUID.
3. `task done <id>` -> PATCH `/tasks/:id` (set `status=done`) -> prints `✓` confirmation.
4. All commands must handle network errors gracefully (no crashes) and print helpful validation errors for invalid/missing IDs.

## Non-Functional Requirements

- No authentication (explicitly out of scope).
- CORS must allow both:
  - Cloudflare Pages domain
  - Local dev origin (and any secondary configured domain, as needed)
- JSON consistency:
  - Always respond with JSON bodies
  - Always set `Content-Type: application/json`
- Robust error handling:
  - Client and server should handle invalid input without unhandled exceptions.

## Deployment Topology

Back-end:
- Render.com (Node.js + Express)
- Supabase Postgres database (free tier)
- Environment variables:
  - `SUPABASE_URL`
  - `SUPABASE_KEY`
  - `PORT` (Render-provided or explicitly set)
  - `ALLOWED_ORIGINS` (must be updated once the Cloudflare Pages URL is known)

Front-end:
- Cloudflare Pages (Vite build; output `dist`)
- Environment variables:
  - `VITE_API_URL=https://task-manager-api.onrender.com/api`
  - Preview/local: `VITE_API_URL=http://localhost:3000/api`

CLI:
- Local execution; reads `TASK_API_URL` from local `.env`

## Implementation Notes

- Implementation must use the PRD’s database schema (Supabase SQL):
  - `tasks` table with `uuid` primary key using `pgcrypto`’s `gen_random_uuid()`
  - enum types `task_status` and `task_priority`
  - `updated_at` auto-update trigger
- Prefer a thin API layer with centralized validation + error mapping.
- Prefer consistent field names across:
  - DB schema
  - API JSON payloads
  - Frontend client models
  - CLI output parsing/formatting

## Tracking Artifacts

- User stories convoy: `hq-cv-rpvqc`
- User story beads created under `tm-mol-qoto` (see convoy for the full list).

