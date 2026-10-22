# ─────────────────────────────────────────────────────────────────────────────
# Права CI на выкат фронта.
#
# `vsluh-ci` катит бэкенд через SAM и статику через `aws s3 sync`. Политика
# бэкенда заведена отдельно; здесь — ровно то, что нужно фронту: запись в два
# бакета и инвалидация их дистрибуций.
#
# Привязка к конкретным бакетам, а не `*`: ключ лежит в GitHub Actions, и его
# радиус поражения должен быть минимальным.
# ─────────────────────────────────────────────────────────────────────────────

terraform {
  required_version = ">= 1.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }

  backend "s3" {
    bucket         = "vsluh-tfstate-661609779996"
    key            = "aws/ci-frontend/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

variable "ci_user" {
  description = "IAM-пользователь, под которым работает GitHub Actions"
  type        = string
  default     = "vsluh-ci"
}

variable "buckets" {
  description = "Бакеты статики по контурам"
  type        = list(string)
  default     = ["vsluh.club", "dev.vsluh.club"]
}

data "aws_iam_policy_document" "frontend_deploy" {
  statement {
    sid       = "ListSiteBuckets"
    actions   = ["s3:ListBucket"]
    resources = [for b in var.buckets : "arn:aws:s3:::${b}"]
  }

  statement {
    sid = "WriteSiteObjects"
    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:DeleteObject",
    ]
    resources = [for b in var.buckets : "arn:aws:s3:::${b}/*"]
  }

  # Дистрибуцию CI ищет по алиасу домена, поэтому нужен список.
  statement {
    sid = "InvalidateDistributions"
    actions = [
      "cloudfront:ListDistributions",
      "cloudfront:CreateInvalidation",
      "cloudfront:GetInvalidation",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_policy" "frontend_deploy" {
  name        = "vsluh-ci-frontend-deploy"
  description = "Выкат статики лендинга в S3 и сброс кеша CloudFront"
  policy      = data.aws_iam_policy_document.frontend_deploy.json
}

resource "aws_iam_user_policy_attachment" "frontend_deploy" {
  user       = var.ci_user
  policy_arn = aws_iam_policy.frontend_deploy.arn
}
