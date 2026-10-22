# Сертификат для CloudFront обязан лежать в us-east-1 — отсюда алиас провайдера.
# Валидация доменная, записи кладём в Cloudflare (dns.tf).

resource "aws_acm_certificate" "site" {
  provider = aws.us_east_1

  domain_name               = var.domain
  subject_alternative_names = var.extra_domains
  validation_method         = "DNS"

  tags = {
    Product     = var.project
    Environment = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "site" {
  provider = aws.us_east_1

  certificate_arn         = aws_acm_certificate.site.arn
  validation_record_fqdns = [for r in cloudflare_record.cert_validation : r.hostname]
}
