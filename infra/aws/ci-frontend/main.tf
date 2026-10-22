# ─────────────────────────────────────────────────────────────────────────────
# Права CI на сброс кеша CloudFront.
#
# Доступ к бакетам статики выдан политиками самих бакетов (модуль static-site) —
# так радиус поражения ключа из GitHub Actions ограничен ровно ими. У CloudFront
# ресурсных политик нет, поэтому инвалидация требует политики пользователя.
#
# Применяется воркфлоу `infra-ci-permissions.yml`, а не локально: у `vsluh-deploy`
# нет ни одного действия IAM.
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

data "aws_iam_policy_document" "frontend_deploy" {
  # Дистрибуцию CI ищет по алиасу домена, поэтому нужен и список.
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
  description = "Сброс кеша CloudFront при выкате статики"
  policy      = data.aws_iam_policy_document.frontend_deploy.json
}

resource "aws_iam_user_policy_attachment" "frontend_deploy" {
  user       = var.ci_user
  policy_arn = aws_iam_policy.frontend_deploy.arn
}
