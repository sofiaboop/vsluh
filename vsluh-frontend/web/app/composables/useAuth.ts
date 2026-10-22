import { ApiError } from '~/shared/api/apiInstance'
import { fetchMe, loginUser, registerUser, type AuthUser } from '~/shared/api/auth'
import { clearAuthToken, getAuthToken, setAuthToken } from '~/shared/lib/authToken'

export type AuthMode = 'login' | 'register'

export interface AuthFormPayload {
  name?: string
  email: string
  password: string
}

/** Ключ строки в i18n (`auth.errors.*`), а не готовый текст: язык выбирается в рантайме. */
function errorKey(error: unknown): string {
  if (!(error instanceof ApiError)) return 'auth.errors.network'
  switch (error.status) {
    case 401:
      return 'auth.errors.invalidCredentials'
    case 409:
      return 'auth.errors.emailTaken'
    case 422:
      return 'auth.errors.invalidInput'
    case 503:
      return 'auth.errors.unavailable'
    default:
      return 'auth.errors.network'
  }
}

export function useAuth() {
  const user = useState<AuthUser | null>('auth-user', () => null)
  const dialogMode = useState<AuthMode | null>('auth-dialog', () => null)
  const pending = useState('auth-pending', () => false)
  const errorKeyRef = useState<string | null>('auth-error', () => null)
  // Сколько кредитов подарили при регистрации. Число запоминаем отдельно от баланса:
  // поздравление должно говорить про подарок, а не про текущий остаток.
  const giftCredits = useState<number | null>('auth-gift', () => null)

  const isAuthenticated = computed(() => user.value !== null)
  const userName = computed(() => user.value?.name ?? '')
  const credits = computed(() => user.value?.credits ?? 0)

  function openAuth(mode: AuthMode) {
    errorKeyRef.value = null
    dialogMode.value = mode
  }

  function closeAuth() {
    errorKeyRef.value = null
    dialogMode.value = null
  }

  /** Возвращает true, если вход удался: диалогу этого хватает, чтобы закрыться. */
  async function submitAuth(payload: AuthFormPayload): Promise<boolean> {
    if (pending.value) return false

    pending.value = true
    errorKeyRef.value = null
    try {
      const session =
        dialogMode.value === 'register'
          ? await registerUser({
              name: payload.name?.trim() ?? '',
              email: payload.email.trim(),
              password: payload.password,
            })
          : await loginUser({ email: payload.email.trim(), password: payload.password })

      setAuthToken(session.token)
      const registered = dialogMode.value === 'register'
      user.value = session.user
      dialogMode.value = null
      giftCredits.value = registered && session.user.credits > 0 ? session.user.credits : null
      return true
    } catch (error) {
      errorKeyRef.value = errorKey(error)
      return false
    } finally {
      pending.value = false
    }
  }

  /**
   * Восстанавливает сессию по сохранённому токену. Токен может быть просрочен
   * или подписан прежним секретом — тогда просто забываем его.
   */
  async function restoreSession(): Promise<void> {
    if (!getAuthToken() || user.value) return
    try {
      user.value = await fetchMe()
    } catch (error) {
      if (error instanceof ApiError && error.status === 401) {
        clearAuthToken()
      }
      user.value = null
    }
  }

  /** Остаток после списания кредита: профиль перечитывать ради одного числа не нужно. */
  function setCredits(value: number) {
    if (user.value) user.value = { ...user.value, credits: value }
  }

  function closeGift() {
    giftCredits.value = null
  }

  function logout() {
    clearAuthToken()
    user.value = null
    dialogMode.value = null
    errorKeyRef.value = null
    giftCredits.value = null
  }

  return {
    user,
    isAuthenticated,
    userName,
    credits,
    giftCredits,
    dialogMode,
    pending,
    errorKey: errorKeyRef,
    openAuth,
    closeAuth,
    submitAuth,
    setCredits,
    closeGift,
    restoreSession,
    logout,
  }
}
