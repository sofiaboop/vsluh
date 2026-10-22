# ─────────────────────────────────────────────────────────────────────────────
# Статика прода — vsluh.club (+ www.vsluh.club)
#
#   ветка main -> GitHub Actions (npm run generate) -> s3://vsluh.club
#              -> CloudFront -> Cloudflare -> vsluh.club
#
# Apex в Cloudflare задан CNAME: зона умеет CNAME flattening, отдельные A-записи
# на адреса origin не нужны.
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
    key            = "cloudfront/vsluh.club/terraform.tfstate"
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

  environment   = "prod"
  domain        = "vsluh.club"
  extra_domains = ["www.vsluh.club"]
  bucket_name   = "vsluh.club"
}
