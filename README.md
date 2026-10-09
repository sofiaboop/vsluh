# ВСЛУХ — vsluh.club

Тренажёр спонтанной речи: выбираешь темы, барабан выдаёт вопрос,
30 секунд на подготовку и минута, чтобы говорить вслух.

## Структура

```
vsluh/
├── .github/              # Правила (L1/L2) + GitHub Actions
├── docs/                 # L3: архитектура, деплой, дизайн
├── variables/default/    # Переменные CI
├── manifest.yaml         # AWS SAM — заготовка, не деплоится
├── vsluh-backend/        # Go-лямбды — заготовка, в работу не взята
└── vsluh-frontend/web/   # Nuxt 4 — лендинг и Telegram Mini App одной сборкой
```

Лендинг и мини-апп — **один проект**. Telegram определяется в рантайме,
его SDK грузится только внутри мессенджера.

## Стадия

MVP: фронт целиком рабочий (выбор тем → барабан → таймеры → финиш).
Бэкенда и AWS пока нет, хостинг — GitHub Pages, DNS — Cloudflare.

## Локально

```bash
nvm use                                   # Node 22
cd vsluh-frontend/web && npm ci && npm run dev
```

Секреты: `set -a && source .envrc && set +a`.

## Деплой

Пуш в `main` → [deploy-pages.yml](.github/workflows/deploy-pages.yml) → `vsluh.club`.
Подробности — [docs/deploy.md](docs/deploy.md).
