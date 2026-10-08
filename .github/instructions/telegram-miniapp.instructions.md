---
applyTo: "vsluh-miniapp/**"
---

# Telegram Mini App

- SDK — `vue-tg`, плюс официальный скрипт `https://telegram.org/js/telegram-web-app.js`
  подключается через `app.head.script` в `nuxt.config.ts`.
- `initData` снимается в `plugins/telegram.client.ts` и кладётся в модуль
  `shared/lib/telegramAuth.ts` (не в Pinia — нужен до инициализации сторов).
- В запросы `initData` уходит **Base64-кодированным** в заголовке `X-Telegram-Init-Data`,
  иначе API Gateway ломает сырую строку.
- Бэкенд обязан проверять HMAC-SHA256 подпись через `TELEGRAM_BOT_TOKEN`.
  Доверять `initDataUnsafe` нельзя никогда.
- Локальная разработка без Telegram: если `initData` пустой — работаем в гостевом режиме,
  запросы к защищённым ручкам не шлём.
