import { apiFetch } from '~/shared/api/apiInstance'

/** Профиль пользователя. На границе API snake_case превращается в camelCase (L1, правило 3). */
export interface AuthUser {
  id: string
  email: string
  name: string
  createdAt: number
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
}

interface SessionPayload {
  token: string
  user: UserPayload
}

function toUser(payload: UserPayload): AuthUser {
  return { id: payload.id, email: payload.email, name: payload.name, createdAt: payload.created_at }
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
