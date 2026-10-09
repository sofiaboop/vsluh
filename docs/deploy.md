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

AWS не используется — аккаунта нет. `manifest.yaml` и `vsluh-backend/` лежат заготовкой.

## Константы

| Что                   | Значение                                            |
| --------------------- | --------------------------------------------------- |
| Репозиторий           | `sofiaboop/vsluh`, ветка `main`                     |
| Cloudflare Zone ID    | `a4530887122f3fb6b17e7287a9589d73`                  |
| Cloudflare Account ID | `782efe566a56ca624924a8de1050a95f`                  |
| NS у регистратора     | `ivan.ns.cloudflare.com`, `paris.ns.cloudflare.com` |
| DS-запись (DNSSEC)    | key tag `2371`, alg `13`, digest type `2`           |

## DNS

Все записи **проксируются** Cloudflare (`proxied: true`):

```
A      vsluh.club      → 185.199.108.153 / .109.153 / .110.153 / .111.153
CNAME  www.vsluh.club  → sofiaboop.github.io
```

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
