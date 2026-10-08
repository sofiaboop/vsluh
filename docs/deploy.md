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

## DNSSEC

Если `dig A vsluh.club @1.1.1.1` отдаёт **SERVFAIL**, а `dig +cd` (без валидации) отвечает
нормально — значит сломана цепочка DNSSEC: в зоне `.club` лежит DS-запись, а зона на
новых NS не подписана теми же ключами.

Так было при переезде с Namecheap на Cloudflare: DS от старого провайдера остался,
Cloudflare-зона была неподписана → все валидирующие резолверы (1.1.1.1, 8.8.8.8) отдавали
SERVFAIL, а домен выглядел «не прописавшимся».

Лечение — согласовать DS с реальным подписантом:

```bash
# включить подпись на стороне Cloudflare и забрать DS
curl -s -X PATCH -H "Authorization: Bearer $CLOUDFLARE_API_TOKEN" \
  -H "Content-Type: application/json" \
  "https://api.cloudflare.com/client/v4/zones/$ZONE/dnssec" --data '{"status":"active"}'
```

Полученные `key_tag`, `algorithm`, `digest_type`, `digest` вносятся в раздел DNSSEC у
регистратора **вместо** старой DS. Альтернатива — просто удалить DNSSEC у регистратора;
домен заработает, но останется без подписи.

Проверка, что цепочка сошлась:

```bash
dig DS vsluh.club @a.nic.club +short     # key tag должен совпасть с Cloudflare
dig +short A vsluh.club @1.1.1.1         # должен вернуть IP, а не пустоту
```

