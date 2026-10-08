# Деплой

## Схема

```
git push main  →  GitHub Actions (npm run generate)  →  GitHub Pages  →  vsluh.club
                                                            ↑
                                                   Cloudflare DNS (зона vsluh.club)
```

AWS не используется — аккаунта нет. `manifest.yaml` и `vsluh-backend/` лежат заготовкой.

## Константы

| Что                   | Значение                                            |
| --------------------- | --------------------------------------------------- |
| Репозиторий           | `sofiaboop/vsluh`, ветка `main`                     |
| Cloudflare Zone ID    | `a4530887122f3fb6b17e7287a9589d73`                  |
| Cloudflare Account ID | `782efe566a56ca624924a8de1050a95f`                  |
| NS для регистратора   | `ivan.ns.cloudflare.com`, `paris.ns.cloudflare.com` |

## DNS

Апекс `vsluh.club` указывает на GitHub Pages четырьмя A-записями:

```
185.199.108.153
185.199.109.153
185.199.110.153
185.199.111.153
```

`www.vsluh.club` — `CNAME` на `sofiaboop.github.io`.

Записи **не проксируются** Cloudflare (`proxied: false`): GitHub сам выпускает сертификат
Let's Encrypt на апекс, а оранжевое облако мешает его верификации.

Кастомный домен зафиксирован файлом [vsluh-frontend/web/public/CNAME](../vsluh-frontend/web/public/CNAME) —
он попадает в сборку, и GitHub Pages не теряет домен при редеплое.

## Выкат вручную

```bash
cd vsluh-frontend/web
NUXT_PUBLIC_SITE_URL=https://vsluh.club npm run generate
npx serve .output/public          # локальная проверка
```

Дальше достаточно пуша в `main` — workflow соберёт и выложит сам.

## Диагностика

```bash
set -a && source .envrc && set +a
ZONE=a4530887122f3fb6b17e7287a9589d73

# статус зоны (active / pending)
curl -s -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  "https://api.cloudflare.com/client/v4/zones/$ZONE" | python3 -m json.tool | head -20

# что реально отдаёт DNS
dig +short NS vsluh.club @1.1.1.1
dig +short A vsluh.club @1.1.1.1

# состояние GitHub Pages
curl -s -H "Authorization: Bearer $GITHUB_TOKEN" \
  https://api.github.com/repos/sofiaboop/vsluh/pages | python3 -m json.tool
```

Зона в статусе `pending` означает, что NS у регистратора ещё не переключены на Cloudflare —
домен не будет резолвиться нигде, сколько бы записей ни стояло внутри зоны.
