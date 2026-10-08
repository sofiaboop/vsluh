# Деплой

## Окружения

| Ветка  | Окружение | Лендинг             | Mini App                | API                   |
| ------ | --------- | ------------------- | ----------------------- | --------------------- |
| `dev`  | dev       | `dev.vsluh.club`    | `app-dev.vsluh.club`    | `api-dev.vsluh.club`  |
| `main` | prod      | `vsluh.club`        | `app.vsluh.club`        | `api.vsluh.club`      |

Фича-ветки деплой не триггерят — только CI.

## Что нужно один раз

1. **AWS**: аккаунт, IAM-пользователь для деплоя (S3, CloudFront, CloudFormation, Lambda,
   API Gateway, DynamoDB), ACM-сертификат на `*.vsluh.club` в `us-east-1` (требование CloudFront).
2. **Cloudflare**: зона `vsluh.club`, API-токен с правами `Zone:DNS:Edit` для этой зоны.
3. **Terraform**: `cd infra && terraform init && terraform apply` — создаёт S3-бакеты,
   CloudFront-дистрибутивы и DNS-записи.
4. **GitHub**: окружения `dev` и `prod`, в каждом переменные и секреты (ниже).

## Переменные окружений GitHub

Variables (не секреты):

| Имя                                   | Пример                          |
| ------------------------------------- | ------------------------------- |
| `AWS_REGION`                          | `eu-central-1`                  |
| `SITE_URL`                            | `https://vsluh.club`            |
| `API_URL`                             | `https://api.vsluh.club`        |
| `MINIAPP_URL`                         | `https://app.vsluh.club`        |
| `S3_BUCKET`                           | `vsluh-frontend-prod`           |
| `CLOUDFRONT_DISTRIBUTION_ID`          | `E...`                          |
| `MINIAPP_S3_BUCKET`                   | `vsluh-miniapp-prod`            |
| `MINIAPP_CLOUDFRONT_DISTRIBUTION_ID`  | `E...`                          |
| `BOT_NAME`                            | `vsluh_bot`                     |
| `DELETION_PROTECTION_ENABLED`         | `true` на проде, `false` на dev |
| `THROTTLING_RATE_LIMIT`               | `50`                            |
| `THROTTLING_BURST_LIMIT`              | `100`                           |

Secrets:

| Имя                     | Что                                         |
| ----------------------- | ------------------------------------------- |
| `AWS_ACCESS_KEY_ID`     | ключ деплой-пользователя                    |
| `AWS_SECRET_ACCESS_KEY` | секрет деплой-пользователя                  |
| `TELEGRAM_BOT_TOKEN`    | токен бота, нужен для проверки `initData`   |

> Долгоживущие AWS-ключи — временное решение. Как только дойдут руки, переключить
> `configure-aws-credentials` на OIDC (`role-to-assume`) и ключи удалить.

## Выкат

Пуш в `dev` или `main` запускает нужные workflow по изменённым путям:

- `vsluh-frontend/**` → `deploy-frontend.yml`
- `vsluh-miniapp/**` → `deploy-miniapp.yml`
- `vsluh-backend/**`, `manifest.yaml` → `deploy-backend.yml`

Кэш: `_nuxt/*` заливается с `max-age=31536000, immutable`, остальное — с
`no-cache, must-revalidate`, затем инвалидация CloudFront по HTML.

## Ручной выкат статики

```bash
cd vsluh-frontend/web
NUXT_PUBLIC_SITE_URL=https://vsluh.club npm run generate
aws s3 sync .output/public s3://vsluh-frontend-prod --delete
aws cloudfront create-invalidation --distribution-id E... --paths "/" "/index.html"
```

## Бэкенд локально

```bash
cd vsluh-backend && make test-all
sam build --template ../manifest.yaml && sam local start-api
```
