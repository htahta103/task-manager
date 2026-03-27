# Frontend (React + Vite)

This app is built for Cloudflare Pages.

## Local development

```bash
npm ci
VITE_API_URL=http://localhost:8080/api npm run dev
```

## Build

```bash
VITE_API_URL=https://taskmanager-api.fly.dev/api npm run build
```

## Cloudflare Pages deployment

The project is configured for Pages via `wrangler.toml`:

- Build output directory: `dist`
- Node runtime pinned in `.nvmrc` (Vite 8 requires Node 20+)

Deploy command (from `frontend/`):

```bash
CLOUDFLARE_ACCOUNT_ID=353a69adadf77371f073ef6bb0a66f4c \
npx wrangler pages deploy dist --project-name taskmanager-staging --branch main
```

Required build variable in Cloudflare Pages:

- `VITE_API_URL=https://taskmanager-api.fly.dev/api`
