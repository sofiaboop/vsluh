import { inTelegram } from '~/shared/lib/telegram'

/** Запущено ли приложение внутри Telegram Mini App. */
export function useTelegram() {
  return { inTelegram }
}
