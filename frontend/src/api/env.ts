export function getApiBaseUrl(): string {
  const raw = import.meta.env.VITE_API_URL as string | undefined
  if (raw && raw.trim().length > 0) return raw.replace(/\/+$/, '')
  return 'http://localhost:8080/api'
}

