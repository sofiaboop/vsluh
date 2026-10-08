# Инфраструктура ВСЛУХ

Terraform создаёт статический хостинг для лендинга и мини-аппа:
приватный S3 + CloudFront (OAC) + ACM-сертификат + DNS-записи в Cloudflare.

## Требования

- `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY` — в окружении
- `CLOUDFLARE_API_TOKEN` — токен с правами `Zone:DNS:Edit` на зону `vsluh.club`

Удобнее всего через `direnv`: скопировать `.envrc.example` в `.envrc` и заполнить.

## Запуск

```bash
cd infra
cp terraform.tfvars.example terraform.tfvars   # вписать cloudflare_zone_id
terraform init
terraform plan
terraform apply
```

Валидация ACM идёт через DNS и занимает несколько минут — `apply` будет ждать.

После `apply` перенести выводы в GitHub Environment:

```bash
terraform output
```

| Вывод                     | Переменная GitHub                      |
| ------------------------- | -------------------------------------- |
| `site_bucket`             | `S3_BUCKET`                            |
| `site_distribution_id`    | `CLOUDFRONT_DISTRIBUTION_ID`           |
| `miniapp_bucket`          | `MINIAPP_S3_BUCKET`                    |
| `miniapp_distribution_id` | `MINIAPP_CLOUDFRONT_DISTRIBUTION_ID`   |

## Два окружения

State у окружений раздельный — используем workspaces:

```bash
terraform workspace new prod && terraform apply -var environment=prod
terraform workspace new dev  && terraform apply -var environment=dev
```

> Бэкенд state пока локальный. Перед работой в команде переключить на S3-бэкенд
> с DynamoDB-локом.
