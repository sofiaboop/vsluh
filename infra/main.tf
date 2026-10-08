locals {
  is_prod = var.environment == "prod"

  site_domain    = local.is_prod ? var.root_domain : "dev.${var.root_domain}"
  miniapp_domain = local.is_prod ? "app.${var.root_domain}" : "app-dev.${var.root_domain}"

  tags = {
    Product     = var.project_name
    Environment = var.environment
    ManagedBy   = "terraform"
  }
}

# Один сертификат на оба хоста окружения
resource "aws_acm_certificate" "site" {
  provider                  = aws.us_east_1
  domain_name               = local.site_domain
  subject_alternative_names = [local.miniapp_domain]
  validation_method         = "DNS"
  tags                      = local.tags

  lifecycle {
    create_before_destroy = true
  }
}

resource "cloudflare_record" "acm_validation" {
  for_each = {
    for dvo in aws_acm_certificate.site.domain_validation_options : dvo.domain_name => {
      name  = dvo.resource_record_name
      value = dvo.resource_record_value
      type  = dvo.resource_record_type
    }
  }

  zone_id = var.cloudflare_zone_id
  name    = trimsuffix(each.value.name, ".")
  content = trimsuffix(each.value.value, ".")
  type    = each.value.type
  ttl     = 60
  proxied = false
}

resource "aws_acm_certificate_validation" "site" {
  provider                = aws.us_east_1
  certificate_arn         = aws_acm_certificate.site.arn
  validation_record_fqdns = [for r in cloudflare_record.acm_validation : r.hostname]
}

module "site" {
  source = "./modules/static-site"

  bucket_name         = "${var.project_name}-frontend-${var.environment}"
  domain              = local.site_domain
  acm_certificate_arn = aws_acm_certificate_validation.site.certificate_arn
  cloudflare_zone_id  = var.cloudflare_zone_id
  tags                = local.tags
}

module "miniapp" {
  source = "./modules/static-site"

  bucket_name         = "${var.project_name}-miniapp-${var.environment}"
  domain              = local.miniapp_domain
  acm_certificate_arn = aws_acm_certificate_validation.site.certificate_arn
  cloudflare_zone_id  = var.cloudflare_zone_id
  tags                = local.tags
}
