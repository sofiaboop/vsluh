# Деплой

## Схема

```
git push → GitHub Actions (npm run generate) → S3 → CloudFront → Cloudflare → домен
```

| Ветка  | Контур | Бакет            | Дистрибуция     | Домен                      |
| ------ | ------ | ---------------- | --------------- | -------------------------- |
| `dev`  | dev    | `dev.vsluh.club` | `EK7DCB0N5ONL0` | `dev.vsluh.club`           |
| `main` | prod   | `vsluh.club`     | `E2PDKIW0GX4CYJ` | `vsluh.club`, `www.vsluh.club` |

Бакеты закрыты от интернета целиком: публичный доступ заблокирован, читать их
может только CloudFront через Origin Access Control. Прямой S3-URL не работает —
и это правильно.

Cloudflare проксирует трафик и терминирует TLS своим сертификатом (режим `Full`),
до CloudFront идёт отдельное HTTPS-соединение с сертификатом ACM из `us-east-1`.
Записи сайта обязаны быть оранжевыми: серая выдаст наружу домен дистрибуции.

Бэкенд (`manifest.yaml`, `vsluh-backend/`) катится в AWS-аккаунт `661609779996`
через [deploy-backend.yml](../.github/workflows/deploy-backend.yml).

## AWS

| Что                     | Значение                                            |
| ----------------------- | --------------------------------------------------- |
| Аккаунт                 | `661609779996` (management-аккаунт `o-wnp0oanxbp`)  |
| Регион                  | `eu-central-1`                                      |
| Доступ человека         | SSO-роль `AdministratorAccess` через access portal  |
| Локальный ключ          | IAM-пользователь `vsluh-deploy`, лежит в `.envrc`   |
| Ключ CI                 | IAM-пользователь `vsluh-ci`, в GitHub Secrets       |
| Terraform state         | `s3://vsluh-tfstate-661609779996`, блокировки — DynamoDB `vsluh-tflock` |
| Бюджет-алерт            | `vsluh-monthly`, 20 USD/мес                         |

Оба IAM-пользователя ограничены политиками с тем же именем: полный доступ к сервисам
проекта (Lambda, API Gateway, DynamoDB, S3, CloudFront, CloudFormation, Logs),
из IAM — только управление ролями лямбд. Админки IAM и ротации ключей у них нет.

### Terraform

Гоняется локально, стейт — удалённый:

```hcl
terraform {
  backend "s3" {
    bucket         = "vsluh-tfstate-661609779996"
    key            = "<component>/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
}
```

У бакета включены версионирование, SSE-AES256 и полный block public access.

### Бэкенд через CI

| Ветка  | Окружение | Стек         | API                                                      |
| ------ | --------- | ------------ | -------------------------------------------------------- |
| `dev`  | `dev`     | `vsluh-dev`  | `https://l9rmwlm921.execute-api.eu-central-1.amazonaws.com/dev`  |
| `main` | `prod`    | `vsluh-prod` | `https://8wqkgr6sa0.execute-api.eu-central-1.amazonaws.com/prod` |

Кастомные домены `api.vsluh.club` / `api-dev.vsluh.club` из
[variables.yml](../variables/default/variables.yml) пока не настроены — API отвечает
по адресу, который выдал API Gateway. Проверка живости: `GET /echo`.

Workflow поднимает Go-лямбды (`sam build`) и катит стек (`sam deploy --resolve-s3`).
Параметры стека собираются из [variables.yml](../variables/default/variables.yml)
по префиксу `USER_BACKEND_{DEV|PROD}_`; секретные значения — одноимённые GitHub Secrets.
Ручной запуск с выбором окружения — через `workflow_dispatch`.

Секреты репозитория: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`,
`USER_BACKEND_PROD_TELEGRAM_BOT_TOKEN`, `USER_BACKEND_DEV_DASHSCOPE_API_KEY`,
`USER_BACKEND_PROD_DASHSCOPE_API_KEY`.
Переменные: `AWS_REGION`, `AWS_ACCOUNT_ID`.

> Для `dev`-окружения секрет `USER_BACKEND_DEV_TELEGRAM_BOT_TOKEN` ещё не заведён —
> нужен токен бота `vsluh_dev_bot`. Без него стек поднимется с пустым токеном
> и проверка `initData` в dev работать не будет.

## Распознавание речи — Alibaba Model Studio

| Что      | Значение                                        |
| -------- | ----------------------------------------------- |
| Аккаунт  | `5673492054411126`, RAM-пользователь `vsluh`    |
| Регион   | Сингапур (`ap-southeast-1`), контур `-intl`     |
| Эндпоинт | `https://dashscope-intl.aliyuncs.com`           |
| Модель   | `qwen-audio-3.0-asr-flash`                      |
| Ключ     | `DASHSCOPE_API_KEY` в `.envrc`, в CI — `USER_BACKEND_{DEV\|PROD}_DASHSCOPE_API_KEY` |

Регионы Model Studio изолированы: ключ сингапурского контура в Пекине отдаёт 401.
Ключ workspace-scoped (`sk-ws-…`), открыт доступ только к включённым в workspace моделям —
остальные отвечают `AccessDenied.Unpurchased`.

Вызов — нативный мультимодальный эндпоинт, **не** OpenAI-совместимый:

```http
POST /api/v1/services/aigc/multimodal-generation/generation
```

```json
{
  "model": "qwen-audio-3.0-asr-flash",
  "input": { "messages": [{ "role": "user", "content": [{ "audio": "<url или data:audio/wav;base64,…>" }] }] },
  "parameters": { "format": "wav", "asr_options": { "language": "ru", "enable_lid": true, "enable_itn": true } }
}
```

`parameters.format` обязателен — без него приходит `UNSUPPORTED_FORMAT: format is empty`,
и он не читается из data-URI. Расшифровка лежит в `output.text`, потайминговая
разбивка — в `output.sentence.words`, израсходованные секунды — в `usage.duration`.
Русский проверен: распознаёт, расставляет пунктуацию.


## Инфраструктура как код — что где

| Слой                                     | Чем описан                              |
| ---------------------------------------- | --------------------------------------- |
| Cloudflare: зона и TLS-настройки         | Terraform, `infra/cloudflare/vsluh.club` |
| S3 + CloudFront + ACM + запись сайта     | Terraform, `infra/cloudfront/<домен>`   |
| AWS: лямбды, API Gateway, DynamoDB       | SAM, `manifest.yaml`, катится из CI      |
| GitHub: репозиторий и секреты            | руками через API (провайдер не заведён)  |

Раскладка каталогов повторяет `aimour`: `infra/<провайдер>/<ресурс>/` — по одному
стейту на компонент, `backend "s3"` с ключом `<провайдер>/<ресурс>/terraform.tfstate`.
Общая часть статического хостинга вынесена в модуль [infra/modules/static-site](../infra/modules/static-site).

### Статический хостинг

```bash
set -a && source .envrc && set +a
cd infra/cloudfront/dev.vsluh.club     # или infra/cloudfront/vsluh.club
terraform init && terraform plan       # «No changes» — норма
```

Модуль поднимает закрытый бакет, сертификат ACM в `us-east-1` с DNS-валидацией,
дистрибуцию с OAC и политикой заголовков безопасности и саму запись в Cloudflare.
Запись сайта живёт в том же стейте, что и дистрибуция, иначе они разъезжаются.

HSTS стоит на `max-age=300` — ступень лестницы. Поднимать только постепенно:
ошибку в большом `max-age` у уже зашедших посетителей откатить нельзя.

### Cloudflare

```bash
set -a && source .envrc && set +a       # CLOUDFLARE_API_TOKEN + ключ vsluh-deploy
cd infra/cloudflare/vsluh.club
terraform init
terraform plan                          # «No changes» — норма
```

Zone id не хардкодится, резолвится по имени через `data.cloudflare_zones`.
Здесь остались только настройки зоны (TLS и перф): записи сайта уехали в стейты
`infra/cloudfront/<домен>`. Почта Namecheap (`MX`, SPF) и DS-запись DNSSEC намеренно
вне стейта — их заводил регистратор, Terraform неуправляемые записи не трогает.

Ресурс, который уже существует в облаке, добавляем только через `import`: сначала
блок `import`, потом `plan` до нуля изменений, и только затем `apply`.

## Константы

| Что                   | Значение                                            |
| --------------------- | --------------------------------------------------- |
| Репозиторий           | `sofiaboop/vsluh`, ветка `main`                     |
| Cloudflare Zone ID    | `a4530887122f3fb6b17e7287a9589d73`                  |
| Cloudflare Account ID | `782efe566a56ca624924a8de1050a95f`                  |
| NS у регистратора     | `ivan.ns.cloudflare.com`, `paris.ns.cloudflare.com` |
| DS-запись (DNSSEC)    | key tag `2371`, alg `13`, digest type `2`           |

## DNS

Записи сайта создаёт модуль статического хостинга
([infra/modules/static-site/dns.tf](../infra/modules/static-site/dns.tf)) и **проксирует**
их Cloudflare (`proxied: true`):

```
CNAME  vsluh.club       → dkdwy3qlpoh8y.cloudfront.net
CNAME  www.vsluh.club   → dkdwy3qlpoh8y.cloudfront.net
CNAME  dev.vsluh.club   → d14ad7v88t1xys.cloudfront.net
```

Апекс задан CNAME — Cloudflare умеет CNAME flattening и отдаёт на него свои адреса.
Рядом лежат серые `_acme-challenge`-подобные записи валидации ACM: их прокси сломает,
поэтому они намеренно `proxied: false`.

Почта Namecheap (`MX` на `eforward*.registrar-servers.com` и `TXT` со SPF) не проксируется
и к деплою отношения не имеет — не трогать.

Настройки зоны:

| Настройка          | Значение | Зачем                                   |
| ------------------ | -------- | --------------------------------------- |
| SSL mode           | `full`   | до origin тоже по HTTPS                 |
| Universal SSL      | on       | сертификат Let's Encrypt от Cloudflare  |
| Always Use HTTPS   | on       | редирект `http` → `https`               |
| Minimum TLS        | `1.2`    | отсекаем устаревшие протоколы           |

> `Always Use HTTPS` включать **только после** выпуска сертификата. Иначе
> посетители уедут на нерабочий протокол и сайт ляжет целиком.

## Выкат

Пуш в `dev` или `main` → [deploy-frontend.yml](../.github/workflows/deploy-frontend.yml):
`npm ci && npm run generate`, затем `aws s3 sync` в бакет контура и инвалидация CloudFront.
Контур выбирается по ветке, переменные сборки берутся из `variables/default/variables.yml`.

Кеш раскладывается тремя проходами — иначе клиент либо залипнет на старой сборке,
либо будет каждый раз тянуть мегабайты ассетов заново:

| Что                     | `Cache-Control`                        |
| ----------------------- | -------------------------------------- |
| `_nuxt/**` и прочие ассеты | `public, max-age=31536000, immutable` |
| `*.html`, `*payload.json`  | `no-cache, must-revalidate`           |
| `robots.txt`               | `public, max-age=300, must-revalidate` |

Имена ассетов хешированные, поэтому `--delete` в `sync` безопасен. После заливки
воркфлоу находит дистрибуцию по алиасу домена, запускает инвалидацию `/*`
и ждёт её завершения.

На тестовом контуре `robots.txt` перезаписывается на `Disallow: /`, чтобы
`dev.vsluh.club` не попал в поиск.

Локальная проверка сборки:

```bash
cd vsluh-frontend/web && npm ci && npm run generate
npx serve .output/public
```

## Тестовый контур

```
sofiaboop/vsluh, ветка dev
        │  push
        ▼
GitHub Actions: deploy-frontend.yml
   npm ci → npm run generate → .output/public
        │  aws s3 sync + cloudfront create-invalidation
        ▼
   s3://dev.vsluh.club  (приватный, читает только CloudFront через OAC)
        ▼
   CloudFront EK7DCB0N5ONL0
        ▲
        │  origin по HTTPS, SSL mode full, сертификат ACM
   Cloudflare (proxy + TLS)
        ▲
   dev.vsluh.club
```

Контур устроен ровно как прод — тот же модуль Terraform, тот же воркфлоу,
отличаются только домен, бакет и дистрибуция. Отдельных репозиториев и артефактных
веток больше нет: GitHub Pages выпилен целиком.

| Что              | Значение                                                        |
| ---------------- | --------------------------------------------------------------- |
| Ветка исходников | `dev` в `sofiaboop/vsluh`                                         |
| Бакет            | `dev.vsluh.club`                                                  |
| Дистрибуция      | `EK7DCB0N5ONL0`                                                   |
| DNS              | `CNAME dev → d14ad7v88t1xys.cloudfront.net`, `proxied: true`      |
| Terraform        | [infra/cloudfront/dev.vsluh.club](../infra/cloudfront/dev.vsluh.club) |
| Воркфлоу         | [deploy-frontend.yml](../.github/workflows/deploy-frontend.yml)   |
| Триггер          | push в `dev`, плюс `workflow_dispatch`                            |

Отличия сборки от прода:

- `NUXT_PUBLIC_SITE_URL=https://dev.vsluh.club`;
- `robots.txt` перезаписывается на `Disallow: /`, чтобы тест не попал в поиск.

TLS тот же, что на проде: Universal SSL от Cloudflare покрывает `*.vsluh.club`,
до CloudFront — отдельное HTTPS-соединение с сертификатом ACM.

### Как катим

```bash
# 1. правка едет в dev
git checkout dev
git commit -am "..."
git push origin dev

# 2. смотрим https://dev.vsluh.club

# 3. только после проверки — в прод
git checkout main && git merge --ff-only dev && git push origin main
```

Подтянуть в `dev` то, что уже в `main` (например после хотфикса):

```bash
git checkout dev && git merge --ff-only main && git push origin dev
```

Откат — реверт коммита в ветке и пуш: бакет перезаписывается целиком каждым
деплоем (`sync --delete`), отдельного отката артефакта не нужно.

### Проверка контура

```bash
set -a && source .envrc && set +a

# что отдаёт тест
curl -sI https://dev.vsluh.club/ | head -3
curl -s https://dev.vsluh.club/robots.txt

# последние запуски деплоя
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  "https://api.github.com/repos/sofiaboop/vsluh/actions/workflows/deploy-frontend.yml/runs?per_page=3" \
  | python3 -c 'import json,sys;[print(r["head_branch"],r["status"],r["conclusion"],r["html_url"]) for r in json.load(sys.stdin)["workflow_runs"]]'

# что реально лежит в бакете и какой у объектов кеш
aws s3 ls s3://dev.vsluh.club --recursive | head
aws s3api head-object --bucket dev.vsluh.club --key index.html --query CacheControl
```

## Диагностика

```bash
set -a && source .envrc && set +a
ZONE=a4530887122f3fb6b17e7287a9589d73

# сертификат, который реально отдаётся
echo | openssl s_client -connect vsluh.club:443 -servername vsluh.club 2>/dev/null \
  | openssl x509 -noout -subject -issuer -dates

# статус пакета сертификатов Cloudflare
curl -s -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  "https://api.cloudflare.com/client/v4/zones/$ZONE/ssl/certificate_packs?status=all"

# что мешает валидации
curl -s -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  "https://api.cloudflare.com/client/v4/zones/$ZONE/ssl/verification"
```

## Грабли, на которые уже наступили

**SERVFAIL при резолве.** В реестре `.club` висела DS-запись от прежнего подписанта
(key tag `30543`), а зона на Cloudflare была неподписана — цепочка DNSSEC рвалась,
и валидирующие резолверы (1.1.1.1, 8.8.8.8) отдавали SERVFAIL. Лечится согласованием
DS с реальным подписантом. Признак: `dig +cd` отвечает, обычный `dig` — нет.

**Панель регистратора врёт.** В Namecheap DNSSEC показывался выключенным, хотя DS
в реестре присутствовала. Проверять только через `dig DS vsluh.club @a.nic.club`
и `whois` (реестр и регистратор отвечают разными строками).

**Выпуск сертификата нельзя торопить.** Каждый перезапуск ACME засчитывается
Let's Encrypt как неудачная валидация; после пяти за час домен блокируется.
Если статус `pending_validation` или `bad_authz` — ждать, не дёргая.

**Сертификат выпустится только после починки DNS.** И Cloudflare, и GitHub
проверяют владение доменом через DNS. Пока домен не резолвится, выпуск будет
падать бесконечно, сколько ни перезапускай.
