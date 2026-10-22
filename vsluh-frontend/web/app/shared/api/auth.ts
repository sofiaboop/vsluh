import { apiFetch } from '~/shared/api/apiInstance'

/** Профиль пользователя. На границе API snake_case превращается в camelCase (L1, правило 3). */
export interface AuthUser {
  id: string
  email: string
  name: string
  createdAt: number
  /** Остаток разборов речи: один разбор стоит один кредит. */
  credits: number
}

export interface AuthSession {
  token: string
  user: AuthUser
}

interface UserPayload {
  id: string
  email: string
  name: string
  created_at: number
  credits: number
}

interface SessionPayload {
  token: string
  user: UserPayload
}

function toUser(payload: UserPayload): AuthUser {
  return {
    id: payload.id,
    email: payload.email,
    name: payload.name,
    createdAt: payload.created_at,
    credits: payload.credits,
  }
}

export async function registerUser(input: {
  name: string
  email: string
  password: string
}): Promise<AuthSession> {
  const payload = await apiFetch<SessionPayload>('/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify({ name: input.name, email: input.email, password: input.password }),
  })
  return { token: payload.token, user: toUser(payload.user) }
}

export async function loginUser(input: { email: string; password: string }): Promise<AuthSession> {
  const payload = await apiFetch<SessionPayload>('/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email: input.email, password: input.password }),
  })
  return { token: payload.token, user: toUser(payload.user) }
}

export async function fetchMe(): Promise<AuthUser> {
  const payload = await apiFetch<{ user: UserPayload }>('/v1/auth/me')
  return toUser(payload.user)
}

/** Списывает кредит за разбор и возвращает остаток. 402 — кредиты кончились. */
export async function spendCredit(): Promise<number> {
  const payload = await apiFetch<{ credits: number }>('/v1/credits/spend', { method: 'POST' })
  return payload.credits
}
