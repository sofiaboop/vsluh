# ─────────────────────────────────────────────────────────────────────────────
# Статика тестового контура — dev.vsluh.club
#
#   ветка dev -> GitHub Actions (npm run generate) -> s3://dev.vsluh.club
#             -> CloudFront -> Cloudflare -> dev.vsluh.club
#
# Контур закрыт от индексации: CI перезаписывает robots.txt на Disallow.
# ─────────────────────────────────────────────────────────────────────────────

terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.52.0"
    }
  }

  backend "s3" {
    bucket         = "vsluh-tfstate-661609779996"
    key            = "cloudfront/dev.vsluh.club/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

# Сертификаты CloudFront живут только в us-east-1.
provider "aws" {
  alias  = "us_east_1"
  region = "us-east-1"
}

provider "cloudflare" {
  # api_token берётся из переменной окружения CLOUDFLARE_API_TOKEN (direnv / .envrc)
}

module "site" {
  source = "../../modules/static-site"

  providers = {
    aws           = aws
    aws.us_east_1 = aws.us_east_1
  }

  environment = "dev"
  domain      = "dev.vsluh.club"
  bucket_name = "dev.vsluh.club"
}
