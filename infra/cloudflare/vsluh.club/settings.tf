# Настройки зоны. Управляем только тем, что влияет на доступность и TLS,
# остальное (security_level, cache level, …) оставлено как есть.
#
#   ssl = "full"            — до origin тоже по HTTPS, но без проверки цепочки:
#                             GitHub Pages за прокси свой сертификат для домена
#                             не выпускает. "strict" здесь положит сайт.
#   always_use_https = "on" — включать только после выпуска сертификата, иначе
#                             редирект уводит на нерабочий протокол.
#   min_tls_version = "1.2" — отсекаем устаревшие протоколы.
#
# Значения совпадают с текущими в панели — plan после импорта должен быть пустым.

resource "cloudflare_zone_settings_override" "this" {
  zone_id = local.zone_id

  settings {
    ssl                      = "full"
    always_use_https         = "on"
    automatic_https_rewrites = "on"
    min_tls_version          = "1.2"
    tls_1_3                  = "on"
    brotli                   = "on"
    http3                    = "on"
    early_hints              = "off"
  }
}
