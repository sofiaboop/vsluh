---
applyTo: "vsluh-frontend/**"
---

# Nuxt 4 / Vue 3

Один проект обслуживает и лендинг, и Telegram Mini App — разделять их не нужно.
Структура — Feature-Sliced Design.

```
app/
├── pages/        # роуты, тонкие — делегируют в features
├── widgets/      # составные блоки экрана (TopBar, BackdropGlow)
├── features/     # topic-picker, wheel, timer, finish — UI в <name>/ui/
├── entities/     # доменные сущности
├── shared/       # api/, lib/, constants/, data/
├── composables/  # usePractice, useTelegram
└── layouts/
```

## Правила

- Только `<script setup lang="ts">`. Пропсы и эмиты типизированы дженериками.
- Импорты — абсолютные через `~/`. Кросс-слойные относительные импорты запрещены.
- Стили экранов лежат в `app/assets/css/tailwind.css` и сняты 1:1 с макета.
  Трогать значения без сверки с эталоном нельзя (см. L1, правило 1).
- Тексты интерфейса — в `shared/constants/`, банк вопросов — в `shared/data/`.
- Состояние практики — composable `usePractice` на `useState` (SSR-safe,
  общий между компонентами). Таймеры и `requestAnimationFrame` держим на уровне
  модуля: экран-инициатор размонтируется раньше, чем они отработают.
- Доступность: интерактив — `<button>`/`<a>`, с `aria-pressed` / `aria-label`.

## Telegram Mini App

- Среда определяется в рантайме (`shared/lib/telegram.ts`): `window.TelegramWebviewProxy`
  либо `tgWebApp` в хэше. SDK грузится динамически **только внутри Telegram**,
  чтобы не замедлять лендинг.
- `initData` снимается в `plugins/telegram.client.ts` и кладётся в модуль
  `shared/lib/telegramAuth.ts` — он нужен до инициализации любых сторов.
- В запросы `initData` уходит **Base64-кодированным** в заголовке `X-Telegram-Init-Data`,
  иначе API Gateway ломает сырую строку.
- Бэкенд обязан проверять HMAC-SHA256 подпись через `TELEGRAM_BOT_TOKEN`.
  Доверять `initDataUnsafe` нельзя никогда.
- Отличия интерфейса внутри Telegram — через класс `is-telegram` на `<html>` (ставит плагин)
  и `useTelegram().inTelegram` в логике, а не отдельной сборкой. Класс на `<html>` надёжнее
  реактивного бинда: плагин отрабатывает после гидрации, и Vue уже не патчит статичный `class`.
