import { getApiBaseUrl } from './env'
import type { ApiError } from './types'

export type RequestOptions = {
  method?: 'GET' | 'POST' | 'PATCH' | 'DELETE'
  path: string
  query?: Record<string, string | number | boolean | undefined>
  body?: unknown
  token?: string | null
  signal?: AbortSignal
}

function toQueryString(query: NonNullable<RequestOptions['query']>) {
  const params = new URLSearchParams()
  for (const [k, v] of Object.entries(query)) {
    if (v === undefined) continue
    params.set(k, String(v))
  }
  const s = params.toString()
  return s ? `?${s}` : ''
}

export async function apiRequest<T>(opts: RequestOptions): Promise<T> {
  const baseUrl = getApiBaseUrl()
  const qs = opts.query ? toQueryString(opts.query) : ''
  const url = `${baseUrl}${opts.path}${qs}`

  const headers: Record<string, string> = {
    Accept: 'application/json',
  }
  if (opts.body !== undefined) headers['Content-Type'] = 'application/json'
  if (opts.token) headers.Authorization = `Bearer ${opts.token}`

  const res = await fetch(url, {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body === undefined ? undefined : JSON.stringify(opts.body),
    signal: opts.signal,
  })

  const contentType = res.headers.get('content-type') ?? ''
  const isJson = contentType.includes('application/json')
  const payload = isJson ? await res.json().catch(() => undefined) : undefined

  if (!res.ok) {
    const fallback: ApiError = { message: `Request failed (${res.status})` }
    const err: ApiError =
      payload && typeof payload === 'object' && 'message' in payload
        ? (payload as ApiError)
        : fallback
    throw new Error(err.message)
  }

  return payload as T
}

