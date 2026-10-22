import { encodeInitData, getInitData } from '~/shared/lib/telegramAuth'
import { getAuthToken } from '~/shared/lib/authToken'

export class ApiError extends Error {
  constructor(
    readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

/**
 * Запрос к API ВСЛУХ. initData уходит Base64-кодированным в X-Telegram-Init-Data,
 * токен сессии — в Authorization.
 */
export async function apiFetch<T>(path: string, options: RequestInit = {}): Promise<T> {
  const { apiUrl } = useRuntimeConfig().public

  const headers = new Headers(options.headers)
  headers.set('Content-Type', 'application/json')

  const initData = getInitData()
  if (initData) {
    headers.set('X-Telegram-Init-Data', encodeInitData(initData))
  }

  const token = getAuthToken()
  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`${apiUrl}${path}`, { ...options, headers })
  if (!response.ok) {
    const detail = await response.text()
    throw new ApiError(response.status, detail || response.statusText)
  }

  return (await response.json()) as T
}
