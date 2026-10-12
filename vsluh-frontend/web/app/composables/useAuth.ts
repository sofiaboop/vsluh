/**
 * Состояние входа. Бэкенда пока нет, поэтому форма ничего не отправляет —
 * экраны готовы, подключение к API будет отдельным шагом.
 */
export type AuthMode = 'login' | 'register'

export function useAuth() {
  const isAuthenticated = useState('auth-ok', () => false)
  const userName = useState('auth-name', () => '')
  const dialogMode = useState<AuthMode | null>('auth-dialog', () => null)

  function openAuth(mode: AuthMode) {
    dialogMode.value = mode
  }

  function closeAuth() {
    dialogMode.value = null
  }

  function submitAuth(payload: { name?: string; email: string }) {
    userName.value = payload.name?.trim() || payload.email.split('@')[0] || ''
    isAuthenticated.value = true
    dialogMode.value = null
  }

  function logout() {
    isAuthenticated.value = false
    userName.value = ''
  }

  return { isAuthenticated, userName, dialogMode, openAuth, closeAuth, submitAuth, logout }
}
