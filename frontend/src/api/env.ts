export function getApiBaseUrl(): string {
  const rawPrimary = import.meta.env.VITE_API_URL as string | undefined
  const rawLegacy = import.meta.env.VITE_API_BASE_URL as string | undefined
  const raw = rawPrimary ?? rawLegacy
  if (raw && raw.trim().length > 0) {
    const trimmed = raw.replace(/\/+$/, '')
    return trimmed.endsWith('/api') ? trimmed : `${trimmed}/api`
  }
  return 'http://localhost:8080/api'
}

