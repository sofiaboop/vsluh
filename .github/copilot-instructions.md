# ВСЛУХ — правила проекта (L1)

Монорепозиторий сервиса практики спонтанной речи. Домен — `vsluh.club`.

## Слои знаний

- **L1 (этот файл)** — нерушимые правила и маршрутизация.
- **L2** — `.github/instructions/*.instructions.md`, подключаются автоматически по `applyTo`.
- **L3** — `docs/`, вход через [docs/README.md](../docs/README.md).

## Карта репозитория

| Путь                | Что это                                    |
| ------------------- | ------------------------------------------ |
| `vsluh-frontend/web`| Nuxt 4 — лендинг и веб-приложение          |
| `vsluh-miniapp/web` | Nuxt 4 — Telegram Mini App                 |
| `vsluh-backend`     | Go-лямбды (`lambdas-go/{common,echo,funcs}`)|
| `manifest.yaml`     | AWS SAM: Lambda, API Gateway, DynamoDB     |
| `variables/default` | Переменные CI, **секретов тут нет**        |
| `infra/`            | Terraform: S3, CloudFront, ACM, Cloudflare |

## Нерушимые правила

1. **Никаких секретов в git.** Значения живут в GitHub Actions Secrets и `.envrc` (в `.gitignore`).
   В `variables/default/variables.yml` — только несекретные значения и пустые плейсхолдеры.
2. **Валидация с двух сторон.** Любое бизнес-ограничение проверяется и на фронте (UX: блокируем
   кнопку, показываем подсказку), и на бэке (целостность: 400/409/422).
3. **Именование.** В коде — `camelCase` (TS) / `camelCase` (Go-внутренние поля экспортируются как
   `PascalCase`). В JSON API и в DynamoDB — `snake_case`. Конвертация — только на границе
   (`shared/api/*`, Go-теги `json:"..."`).
4. **Фронт и мини-апп повторяют одну структуру** (Feature-Sliced Design). Отличие мини-аппа —
   авторизация через Telegram `initData`.
5. **Lambda-рантайм** — `provided.al2023`, архитектура `arm64`, хендлер `bootstrap`.
6. **DynamoDB** — `PAY_PER_REQUEST`. На проде `DeletionProtectionEnabled: true`.
7. **Ветки → окружения.** `dev` → dev-стенд, `main` → прод. Фича-ветки деплой не триггерят.
8. **Не коммить сгенерированное**: `.nuxt/`, `.output/`, `.aws-sam/`, `bootstrap`.

## Язык

Интерфейс и контент — русский. Код, имена файлов, коммиты — английский.
