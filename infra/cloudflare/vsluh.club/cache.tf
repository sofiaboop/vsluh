# ─────────────────────────────────────────────────────────────────────────────
# Кеш на уровне Cloudflare.
#
# По умолчанию Cloudflare кеширует только статику по расширению, а `text/html`
# не кеширует вовсе — каким бы ни был Cache-Control. Правило ниже включает кеш
# для всего, срок берётся из заголовков origin'а: их расставляет выкат (ассеты
# immutable, html — нулевой max-age и суточный s-maxage).
#
# ВЫКЛЮЧЕНО: токену Cloudflare в `.envrc` не хватает прав. Нужно добавить ему
# две области на зону `vsluh.club`:
#   • Zone → Cache Rules → Edit   — для правила ниже
#   • Zone → Cache Purge → Purge  — для сброса кеша из CI после выката
# После этого раскомментировать блок и выполнить `terraform apply`.
#
# Пока правило выключено, HTML кеширует CloudFront (s-maxage), а Cloudflare
# пропускает его насквозь. Ассеты с immutable кешируются на всех уровнях и так.
# ─────────────────────────────────────────────────────────────────────────────

# resource "cloudflare_ruleset" "cache" {
#   zone_id = local.zone_id
#   name    = "cache"
#   kind    = "zone"
#   phase   = "http_request_cache_settings"
#
#   rules {
#     ref         = "cache_everything_respect_origin"
#     description = "Кешировать всё, срок — по заголовкам origin"
#     expression  = "true"
#     action      = "set_cache_settings"
#
#     action_parameters {
#       cache = true
#
#       edge_ttl {
#         mode = "respect_origin"
#       }
#
#       browser_ttl {
#         mode = "respect_origin"
#       }
#     }
#   }
# }
