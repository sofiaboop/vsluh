# Запись доменной валидации ACM — служебная, прокси её сломает.
resource "cloudflare_record" "cert_validation" {
  zone_id = local.zone_id
  name = trimsuffix(one([
    for o in aws_acm_certificate.api.domain_validation_options :
    o.resource_record_name if o.domain_name == var.domain
  ]), ".")
  type = one([
    for o in aws_acm_certificate.api.domain_validation_options :
    o.resource_record_type if o.domain_name == var.domain
  ])
  content = trimsuffix(one([
    for o in aws_acm_certificate.api.domain_validation_options :
    o.resource_record_value if o.domain_name == var.domain
  ]), ".")
  proxied         = false
  ttl             = 60
  allow_overwrite = true
  comment         = "ACM validation для ${var.domain}"
}

# 🔴 proxied = true: снаружи TLS терминирует Cloudflare своим сертификатом зоны,
# до API Gateway он идёт по HTTPS и упирается в региональный сертификат ACM.
resource "cloudflare_record" "api" {
  zone_id = local.zone_id
  name    = local.host
  type    = "CNAME"
  content = aws_api_gateway_domain_name.api.regional_domain_name
  proxied = true
  ttl     = 1
  comment = "${var.domain} -> API Gateway (${var.environment})"
}
