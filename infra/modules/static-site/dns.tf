data "cloudflare_zones" "this" {
  filter {
    name = var.zone_name
  }
}

locals {
  zone_id = data.cloudflare_zones.this.zones[0].id
}

# Записи для доменной валидации ACM — по одной на каждый хост сертификата.
# Ключи for_each берём из списка хостов, а не из domain_validation_options:
# последние неизвестны на этапе plan и ломают планирование.
# Серые: это служебный CNAME, прокси его сломает.
resource "cloudflare_record" "cert_validation" {
  for_each = toset(local.aliases)

  zone_id = local.zone_id
  name = trimsuffix(one([
    for o in aws_acm_certificate.site.domain_validation_options :
    o.resource_record_name if o.domain_name == each.value
  ]), ".")
  type = one([
    for o in aws_acm_certificate.site.domain_validation_options :
    o.resource_record_type if o.domain_name == each.value
  ])
  content = trimsuffix(one([
    for o in aws_acm_certificate.site.domain_validation_options :
    o.resource_record_value if o.domain_name == each.value
  ]), ".")
  proxied         = false
  ttl             = 60
  allow_overwrite = true
  comment         = "ACM validation для ${each.value}"
}

# 🔴 proxied = true обязателен: TLS для посетителя терминирует Cloudflare, его же
# сертификат покрывает зону. Серая запись выдаст наружу домен CloudFront и уронит
# совпадение имени в сертификате.
resource "cloudflare_record" "site" {
  for_each = toset(local.aliases)

  zone_id = local.zone_id
  name    = each.value == var.zone_name ? var.zone_name : trimsuffix(each.value, ".${var.zone_name}")
  type    = "CNAME"
  content = aws_cloudfront_distribution.site.domain_name
  proxied = true
  ttl     = 1
  comment = "${each.value} -> CloudFront (${var.environment})"
}
