import { looksLikeInitData, setInitData } from '~/shared/lib/telegramAuth'

declare global {
  interface Window {
    Telegram?: { WebApp?: { initData?: string; ready: () => void; expand: () => void } }
  }
}

export default defineNuxtPlugin(() => {
  const webApp = window.Telegram?.WebApp
  if (!webApp) return

  webApp.ready()
  webApp.expand()

  const raw = webApp.initData
  if (looksLikeInitData(raw)) {
    setInitData(raw)
  }
})
