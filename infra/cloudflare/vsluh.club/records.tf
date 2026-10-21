# DNS сайта. Управляем только тем, что ведёт на GitHub Pages.
#
# 🔴 `proxied = true` здесь не вкусовщина: TLS для домена выпускает и терминирует
# Cloudflare. Серая запись оставит хост на сертификате GitHub, которого для
# кастомного домена за прокси нет — сайт ляжет целиком.
#
# Почта регистратора (MX eforward*.registrar-servers.com и TXT со SPF) намеренно
# не импортирована: её заводил Namecheap, она к деплою отношения не имеет.
# Terraform не удаляет неуправляемые записи — они просто вне стейта.

resource "cloudflare_record" "apex" {
  for_each = toset(var.pages_ipv4)

  zone_id = local.zone_id
  name    = var.zone_name
  type    = "A"
  content = each.value
  proxied = true
  ttl     = 1 # 1 = automatic, обязательно при proxied
}

resource "cloudflare_record" "www" {
  zone_id = local.zone_id
  name    = "www"
  type    = "CNAME"
  content = var.pages_host
  proxied = true
  ttl     = 1
}

resource "cloudflare_record" "test" {
  zone_id = local.zone_id
  name    = trimsuffix(var.test_host, ".${var.zone_name}")
  type    = "CNAME"
  content = var.pages_host
  proxied = true
  ttl     = 1
  comment = "test env -> GitHub Pages sofiaboop/vsluh-test"
}
