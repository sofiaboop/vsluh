# Деплой

## Схема

```
git push main → GitHub Actions (npm run generate) → GitHub Pages
                                                         ↑
                                 Cloudflare (proxy, TLS) ┘
                                          ↑
                                     vsluh.club
```

Cloudflare проксирует трафик и терминирует TLS своим сертификатом.
До GitHub Pages запрос идёт отдельным соединением по HTTPS (режим `Full`).

AWS не используется для статики — лендинг живёт на GitHub Pages.
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

| Ветка  | Окружение | Стек         |
| ------ | --------- | ------------ |
| `dev`  | `dev`     | `vsluh-dev`  |
| `main` | `prod`    | `vsluh-prod` |

Workflow поднимает Go-лямбды (`sam build`) и катит стек (`sam deploy --resolve-s3`).
Параметры стека собираются из [variables.yml](../variables/default/variables.yml)
по префиксу `USER_BACKEND_{DEV|PROD}_`; секретные значения — одноимённые GitHub Secrets.
Ручной запуск с выбором окружения — через `workflow_dispatch`.

Секреты репозитория: `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`,
`USER_BACKEND_PROD_TELEGRAM_BOT_TOKEN`, `TEST_PAGES_TOKEN`.
Переменные: `AWS_REGION`, `AWS_ACCOUNT_ID`.

> Для `dev`-окружения секрет `USER_BACKEND_DEV_TELEGRAM_BOT_TOKEN` ещё не заведён —
> нужен токен бота `vsluh_dev_bot`. Без него стек поднимется с пустым токеном
> и проверка `initData` в dev работать не будет.


## Инфраструктура как код — чего нет

Terraform в проекте **нет**. Cloudflare (зона, DNS, SSL) и GitHub (репозитории, Pages,
секреты, Actions) настроены вручную через API — командами из этого файла. Источник правды
по инфраструктуре — таблицы ниже, а не стейт. Менять руками, сверяясь с `dig` и API.

## Константы

| Что                   | Значение                                            |
| --------------------- | --------------------------------------------------- |
| Репозиторий           | `sofiaboop/vsluh`, ветка `main`                     |
| Cloudflare Zone ID    | `a4530887122f3fb6b17e7287a9589d73`                  |
| Cloudflare Account ID | `782efe566a56ca624924a8de1050a95f`                  |
| NS у регистратора     | `ivan.ns.cloudflare.com`, `paris.ns.cloudflare.com` |
| DS-запись (DNSSEC)    | key tag `2371`, alg `13`, digest type `2`           |

## DNS

Записи сайта **проксируются** Cloudflare (`proxied: true`):

```
A      vsluh.club       → 185.199.108.153 / .109.153 / .110.153 / .111.153
CNAME  www.vsluh.club   → sofiaboop.github.io
CNAME  test.vsluh.club  → sofiaboop.github.io
```

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

Пуш в `main` → [deploy-pages.yml](../.github/workflows/deploy-pages.yml) собирает
`vsluh-frontend/web` и публикует в Pages. Кастомный домен задан в настройках Pages,
файла `CNAME` в репозитории нет намеренно: в артефакте он переустанавливает домен
при каждом деплое и сбрасывает выпуск сертификата.

Локальная проверка сборки:

```bash
cd vsluh-frontend/web && npm ci && npm run generate
npx serve .output/public
```

## Тестовый контур

```
sofiaboop/vsluh, ветка test
        │  push (paths: vsluh-frontend/**)
        ▼
GitHub Actions: deploy-test.yml
   npm ci → npm run generate → .output/public
        │  peaceiris/actions-gh-pages, personal_token = TEST_PAGES_TOKEN
        ▼
sofiaboop/vsluh-test, ветка gh-pages   ← только артефакт, исходников нет
        │  Pages, source = branch gh-pages, домен из файла CNAME
        ▼
   GitHub Pages (sofiaboop.github.io)
        ▲
        │  origin по HTTPS, SSL mode full
   Cloudflare (proxy + TLS)
        ▲
   test.vsluh.club
```

Почему отдельный репозиторий: GitHub Pages отдаёт **один сайт на репозиторий** с одним
кастомным доменом. `sofiaboop/vsluh` уже занят продом на `vsluh.club`, поэтому тест
живёт в `sofiaboop/vsluh-test`. Исходники там не дублируются — только собранная статика.

| Что               | Значение                                                       |
| ----------------- | -------------------------------------------------------------- |
| Ветка исходников  | `test` в `sofiaboop/vsluh`                                       |
| Репозиторий Pages | `sofiaboop/vsluh-test`, ветка `gh-pages`                         |
| Pages source      | `branch gh-pages`, path `/`, build type `legacy`                 |
| Домен             | `test.vsluh.club`, задан файлом `CNAME` в артефакте              |
| DNS               | `CNAME test → sofiaboop.github.io`, `proxied: true`              |
| Воркфлоу          | [deploy-test.yml](../.github/workflows/deploy-test.yml)          |
| Триггер           | push в `test` по путям `vsluh-frontend/**`, плюс `workflow_dispatch` |
| Секрет            | `TEST_PAGES_TOKEN` в `sofiaboop/vsluh` — PAT для пуша в соседний репозиторий |

Отличия сборки от прода:

- `NUXT_PUBLIC_SITE_URL=https://test.vsluh.club`;
- `robots.txt` перезаписывается на `Disallow: /`, чтобы тест не попал в поиск;
- домен задан файлом `CNAME` внутри артефакта. Для прода это запрещено (артефакт
  Actions переустанавливает домен и сбрасывает сертификат), а для деплоя из ветки —
  единственный способ: без `CNAME` в `gh-pages` Pages сбросит кастомный домен.

TLS тот же, что на проде: сертификат Cloudflare (Universal SSL покрывает `*.vsluh.club`),
до GitHub Pages — отдельное HTTPS-соединение в режиме `full`. Собственный сертификат
GitHub для `test.vsluh.club` не нужен и `https_enforced` остаётся `false` — редирект
на HTTPS делает Cloudflare.

### Как катим

```bash
# 1. правка едет в test
git checkout test
git commit -am "..."
git push origin test

# 2. смотрим https://test.vsluh.club

# 3. только после проверки — в прод
git checkout main && git merge --ff-only test && git push origin main
```

Подтянуть в `test` то, что уже в `main` (например после хотфикса):

```bash
git checkout test && git merge --ff-only main && git push origin test
```

Откат теста — реверт коммита в `test` и пуш; `gh-pages` перезаписывается целиком
каждым деплоем, отдельного отката артефакта не нужно.

### Проверка контура

```bash
set -a && source .envrc && set +a

# что отдаёт тест
curl -sI https://test.vsluh.club/ | head -3
curl -s https://test.vsluh.club/robots.txt

# последние запуски деплоя
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  "https://api.github.com/repos/sofiaboop/vsluh/actions/workflows/deploy-test.yml/runs?per_page=3" \
  | python3 -c 'import json,sys;[print(r["head_branch"],r["status"],r["conclusion"],r["html_url"]) for r in json.load(sys.stdin)["workflow_runs"]]'

# настройки Pages тестового репозитория
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" https://api.github.com/repos/sofiaboop/vsluh-test/pages
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
