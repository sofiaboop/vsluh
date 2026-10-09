export interface TelegramWebApp {
  initData?: string
  ready: () => void
  expand: () => void
  colorScheme?: 'light' | 'dark'
  setHeaderColor?: (color: string) => void
  setBackgroundColor?: (color: string) => void
}

declare global {
  interface Window {
    Telegram?: { WebApp?: TelegramWebApp }
    TelegramWebviewProxy?: unknown
  }
}

/**
 * Признак запуска внутри Telegram. Проверяется до загрузки SDK:
 * Telegram передаёт параметры в хэше и подкладывает свой прокси-объект.
 */
export function isTelegramEnv(): boolean {
  if (typeof window === 'undefined') return false
  return Boolean(window.TelegramWebviewProxy) || window.location.hash.includes('tgWebApp')
}

const SDK_URL = 'https://telegram.org/js/telegram-web-app.js'

/** Живём ли внутри Telegram. Выставляется плагином после инициализации SDK. */
export const inTelegram = ref(false)

/** Подгружает официальный SDK — только внутри Telegram, чтобы не тормозить лендинг. */
export function loadTelegramSdk(): Promise<TelegramWebApp | null> {
  if (typeof window === 'undefined') return Promise.resolve(null)
  if (window.Telegram?.WebApp) return Promise.resolve(window.Telegram.WebApp)

  return new Promise((resolve) => {
    const script = document.createElement('script')
    script.src = SDK_URL
    script.async = true
    script.onload = () => resolve(window.Telegram?.WebApp ?? null)
    script.onerror = () => resolve(null)
    document.head.appendChild(script)
  })
}
