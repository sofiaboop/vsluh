/**
 * Токен сессии нужен раньше, чем поднимаются composables (его читает apiFetch),
 * поэтому он живёт в модуле, а localStorage — лишь способ пережить перезагрузку.
 */
const STORAGE_KEY = 'vsluh_token'

let token = ''

export function getAuthToken(): string {
  if (!token && import.meta.client) {
    token = localStorage.getItem(STORAGE_KEY) ?? ''
  }
  return token
}

export function setAuthToken(value: string): void {
  token = value
  if (!import.meta.client) return
  if (value) {
    localStorage.setItem(STORAGE_KEY, value)
  } else {
    localStorage.removeItem(STORAGE_KEY)
  }
}

export function clearAuthToken(): void {
  setAuthToken('')
}
