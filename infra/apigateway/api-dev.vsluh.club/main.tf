# ─────────────────────────────────────────────────────────────────────────────
# Публичный хост API тестового контура — api-dev.vsluh.club
#
#   api-dev.vsluh.club -> Cloudflare (прокси и TLS) -> API Gateway (REGIONAL)
#                      -> стейдж dev стека vsluh-dev (SAM)
#
# Сам REST API и лямбды описаны в manifest.yaml и живут в CloudFormation —
# здесь только доменное имя, сертификат и DNS.
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
    key            = "apigateway/api-dev.vsluh.club/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

provider "cloudflare" {
  # api_token берётся из переменной окружения CLOUDFLARE_API_TOKEN (direnv / .envrc)
}

# Сертификат региональный: у REGIONAL-домена API Gateway он обязан лежать
# в том же регионе, что и сам API, в отличие от CloudFront с его us-east-1.
resource "aws_acm_certificate" "api" {
  domain_name       = var.domain
  validation_method = "DNS"

  tags = {
    Product     = var.project
    Environment = var.environment
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_acm_certificate_validation" "api" {
  certificate_arn         = aws_acm_certificate.api.arn
  validation_record_fqdns = [cloudflare_record.cert_validation.hostname]
}

resource "aws_api_gateway_domain_name" "api" {
  domain_name              = var.domain
  regional_certificate_arn = aws_acm_certificate_validation.api.certificate_arn
  security_policy          = "TLS_1_2"

  endpoint_configuration {
    types = ["REGIONAL"]
  }

  tags = {
    Product     = var.project
    Environment = var.environment
  }
}

# Пустой base_path: /v1/... запросы уходят в стейдж как есть.
resource "aws_api_gateway_base_path_mapping" "api" {
  api_id      = data.aws_api_gateway_rest_api.api.id
  stage_name  = var.environment
  domain_name = aws_api_gateway_domain_name.api.domain_name
}
