# Task Manager Frontend

React + TypeScript + Vite frontend for Task Manager.

## Environment

- `VITE_API_URL` (preferred): backend base URL. If it does not end with `/api`, the app appends it automatically.
- `VITE_API_BASE_URL` (legacy alias): accepted for compatibility with older deploy configs.
- Default (unset): `http://localhost:8080/api`

## Local development

```bash
npm install
npm run dev
```

## Cloudflare Pages deployment

Build output is `dist/` and SPA fallback is configured via `public/_redirects`.

```bash
# Example backend URL:
export VITE_API_URL="https://taskmanager-api.fly.dev/api"

npm run build:pages
npx wrangler pages deploy dist --project-name taskmanager-frontend
```

If your Pages project uses a different name, replace `taskmanager-frontend`.
