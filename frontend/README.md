# Task Manager Frontend

React + TypeScript + Vite app deployed to Cloudflare Pages.

## Local development

```bash
npm ci
npm run dev
```

The API base URL defaults to `http://localhost:8080/api` unless `VITE_API_URL` is set.

## Build

```bash
npm run build
```

For staging builds, `build:staging` sets a default hosted API URL and still allows override:

```bash
VITE_API_URL=https://your-api-host/api npm run build:staging
```

## Deploy to Cloudflare Pages (staging)

Set the Pages project name, then run the deploy script:

```bash
export CLOUDFLARE_PAGES_PROJECT=taskmanager-staging
npm run deploy:staging
```

Notes:
- `deploy:staging` always targets branch `main`.
- `CLOUDFLARE_ACCOUNT_ID` defaults to `353a69adadf77371f073ef6bb0a66f4c`.
- Override `VITE_API_URL` to point at the desired backend before deploying.
