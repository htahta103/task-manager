export function getApiBaseUrl(): string {
  const rawUrl = import.meta.env.VITE_API_URL as string | undefined
  const rawBaseUrl = import.meta.env.VITE_API_BASE_URL as string | undefined
  const raw = rawUrl ?? rawBaseUrl
  if (raw && raw.trim().length > 0) {
    const trimmed = raw.replace(/\/+$/, '')
    return trimmed.endsWith('/api') ? trimmed.slice(0, -4) : trimmed
  }
  return 'http://localhost:8080'
}

