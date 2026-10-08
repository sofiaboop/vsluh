---
applyTo: "manifest.yaml,infra/**,variables/**,.github/workflows/**"
---

# Инфраструктура

## SAM (`manifest.yaml`)

- Runtime `provided.al2023`, `Architectures: [arm64]`, `Handler: bootstrap`.
- Имена ресурсов — через `${CI_PROJECT_NAME}-${ENVIRONMENT}-...`, не хардкод.
- Политики — только SAM-полиси (`DynamoDBCrudPolicy`), ручные IAM-роли не пишем.
- Имена таблиц/очередей генерирует CloudFormation — читать их обратно через `aws cloudformation
  describe-stacks`, а не угадывать.

## Статика (лендинг и мини-апп)

S3 + CloudFront. Кэш-политика:

| Путь                        | Cache-Control                               |
| --------------------------- | ------------------------------------------- |
| `/_nuxt/**` (хэш в имени)   | `public, max-age=31536000, immutable`       |
| `index.html`, `200.html`    | `no-cache, must-revalidate`                 |
| `/robots.txt`, `/sitemap*`  | `public, max-age=300`                       |

После `aws s3 sync` обязательна инвалидация CloudFront по `/index.html` и `/`.

## Переменные

Формат: `USER_{BACKEND|FRONTEND|MINIAPP}_{DEV|PROD}_{KEY}`.
В `variables/default/variables.yml` — только несекретные значения; секреты живут в
GitHub Actions Secrets под тем же именем.
