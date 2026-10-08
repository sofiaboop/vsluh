# ВСЛУХ — vsluh.club

Сервис практики спонтанной речи. Монорепозиторий.

## Структура

```
vsluh/
├── .github/                  # Правила (L1/L2) + GitHub Actions
├── docs/                     # L3: архитектура, деплой, дизайн
├── variables/default/        # Переменные CI
├── manifest.yaml             # AWS SAM — заготовка, не деплоится
├── vsluh-backend/            # Go-лямбды — заготовка, в работу не взята
├── vsluh-frontend/web/       # Nuxt 4 — лендинг vsluh.club
└── vsluh-miniapp/web/        # Nuxt 4 — Telegram Mini App (заготовка)
```

## Стадия

MVP: только визуальный лендинг. Бэкенда и AWS нет.
Хостинг — GitHub Pages, DNS — Cloudflare.

## Локально

```bash
nvm use                                   # Node 22
cd vsluh-frontend/web && npm ci && npm run dev
```

Секреты: `set -a && source .envrc && set +a`.

## Деплой

Пуш в `main` → [deploy-pages.yml](.github/workflows/deploy-pages.yml) → `vsluh.club`.
Подробности — [docs/deploy.md](docs/deploy.md).
