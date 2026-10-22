# ─────────────────────────────────────────────────────────────────────────────
# static-site — S3 + CloudFront + ACM + запись в Cloudflare.
#
# Схема (как в aimour, projects/modules/cloudfront):
#   бакет приватный -> CloudFront читает его через OAC -> Cloudflare проксирует
#   домен на CloudFront и терминирует TLS своим сертификатом.
#
# Отличие от aimour: вместо legacy OAI используется OAC. OAI помечен AWS как
# устаревший и не умеет SigV4 для новых регионов; для нового кода брать его
# смысла нет, поведение то же — бакет остаётся закрытым от прямого доступа.
#
# Контент заливает CI (`aws s3 sync`), Terraform содержимое бакета не трогает.
# ─────────────────────────────────────────────────────────────────────────────

terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source                = "hashicorp/aws"
      version               = "~> 5.0"
      configuration_aliases = [aws.us_east_1]
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.52.0"
    }
  }
}
