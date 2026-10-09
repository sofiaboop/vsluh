import { inTelegram, isTelegramEnv, loadTelegramSdk } from '~/shared/lib/telegram'
import { looksLikeInitData, setInitData } from '~/shared/lib/telegramAuth'

export default defineNuxtPlugin({
  name: 'telegram',
  // не блокируем гидрацию: SDK внешний и может отвечать долго
  parallel: true,
  async setup() {
    if (!isTelegramEnv()) return

    const webApp = await loadTelegramSdk()
    if (!webApp) return

    inTelegram.value = true
    // класс на <html>: не зависит от того, успела ли пройти гидрация
    document.documentElement.classList.add('is-telegram')
    webApp.ready()
    webApp.expand()
    webApp.setHeaderColor?.('#f3f6fa')
    webApp.setBackgroundColor?.('#f3f6fa')

    if (looksLikeInitData(webApp.initData)) {
      setInitData(webApp.initData)
    }
  },
})
