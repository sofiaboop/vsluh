---
applyTo: "manifest.yaml,infra/**,variables/**,.github/workflows/**"
---

# Инфраструктура

## Terraform (`infra/**`)

Всё, кроме ресурсов SAM, описывается Terraform. Руками в панелях и через API —
только разведка и диагностика, не изменения.

- Раскладка как в `aimour`: `infra/<провайдер>/<ресурс>/`, один стейт на компонент.
- Файлы делим по смыслу: `main.tf` (terraform + provider + backend), `data.tf`,
  `variables.tf`, `imports.tf`, дальше по теме — `records.tf`, `settings.tf`.
- Стейт только удалённый:

  ```hcl
  backend "s3" {
    bucket         = "vsluh-tfstate-661609779996"
    key            = "<провайдер>/<ресурс>/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
  ```

- Id зон и аккаунтов не хардкодим — резолвим через `data`.
- Существующий ресурс сначала `import`, затем `plan` до нуля изменений,
  и только потом `apply`. Любой неожиданный diff — повод разбираться, а не применять.
- `terraform` гоняется локально под ключом `vsluh-deploy` из `.envrc`, не из CI.
- `.terraform.lock.hcl` коммитим, `.terraform/` и `*.tfstate` — нет.

## SAM (`manifest.yaml`)

- Runtime `provided.al2023`, `Architectures: [arm64]`, `Handler: bootstrap`.
- Имена ресурсов — через `${CI_PROJECT_NAME}-${ENVIRONMENT}-...`, не хардкод.
- Политики — только SAM-полиси (`DynamoDBCrudPolicy`), ручные IAM-роли не пишем.
- Имена таблиц/очередей генерирует CloudFormation — читать их обратно через `aws cloudformation
  describe-stacks`, а не угадывать.

## Статика (лендинг и мини-апп)

S3 + CloudFront + Cloudflare. Инфраструктура — модуль `infra/modules/static-site`,
по корневому конфигу на домен: `infra/cloudfront/vsluh.club`, `infra/cloudfront/dev.vsluh.club`.
Бакет приватный, читает только CloudFront через OAC; запись в Cloudflare создаёт
тот же модуль, что и дистрибуцию, и всегда `proxied = true`.

Выкат — `deploy-frontend.yml`. Кэш-политика:

| Путь                        | Cache-Control                               |
| --------------------------- | ------------------------------------------- |
| `/_nuxt/**` (хэш в имени)   | `public, max-age=31536000, immutable`       |
| `*.html`, `*payload.json`   | `no-cache, must-revalidate`                 |
| `/robots.txt`, `/sitemap*`  | `public, max-age=300, must-revalidate`      |

После `aws s3 sync` обязательна инвалидация CloudFront с ожиданием завершения.

## Переменные

Формат: `USER_{BACKEND|FRONTEND|MINIAPP}_{DEV|PROD}_{KEY}`.
В `variables/default/variables.yml` — только несекретные значения; секреты живут в
GitHub Actions Secrets под тем же именем.
