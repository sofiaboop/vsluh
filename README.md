# ВСЛУХ — vsluh.club

Сервис практики спонтанной речи. Монорепозиторий.

## Структура

```
vsluh/
├── .github/                  # Copilot-инструкции (L1/L2) + GitHub Actions
├── docs/                     # L3: архитектура, плейбуки, деплой
├── variables/default/        # Переменные окружения для CI (без секретов)
├── manifest.yaml             # AWS SAM: Lambda + API Gateway + DynamoDB
├── vsluh-backend/            # Go-лямбды (go.work, lambdas-go/*)
├── vsluh-frontend/web/       # Nuxt 4 — публичный лендинг/веб-приложение (vsluh.club)
└── vsluh-miniapp/web/        # Nuxt 4 — Telegram Mini App (app.vsluh.club)
```

## Домены

| Домен             | Назначение                     | Хостинг                  |
| ----------------- | ------------------------------ | ------------------------ |
| `vsluh.club`      | Лендинг / веб-приложение       | S3 + CloudFront (Cloudflare DNS) |
| `app.vsluh.club`  | Telegram Mini App              | S3 + CloudFront          |
| `api.vsluh.club`  | REST API (Go Lambda)           | API Gateway + ACM        |

## Быстрый старт

```bash
nvm use                       # Node 22
cd vsluh-frontend/web && npm ci && npm run dev
```

Бэкенд:

```bash
cd vsluh-backend && make deps && make test-all
```

## Документация

Точка входа — [docs/README.md](docs/README.md).
