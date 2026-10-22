# ─────────────────────────────────────────────────────────────────────────────
# IAM-пользователи проекта.
#
#   vsluh-ci     — ключ в GitHub Actions: SAM-стек бэкенда и статика фронта
#   vsluh-deploy — локальный ключ в .envrc: Terraform и разведка
#
# Обе политики заводились руками до Terraform и импортированы как есть, после
# чего дополнены недостающим. Админский доступ — SSO-роль AdministratorAccess
# (профиль `vsluh-admin`), у самих ключей прав на IAM почти нет и это намеренно.
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
    key            = "aws/iam/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
}

provider "aws" {
  region = "eu-central-1"
}

variable "site_buckets" {
  description = "Бакеты статики по контурам"
  type        = list(string)
  default     = ["vsluh.club", "dev.vsluh.club"]
}
