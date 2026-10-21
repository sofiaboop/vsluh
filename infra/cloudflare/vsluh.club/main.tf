# ─────────────────────────────────────────────────────────────────────────────
# Cloudflare zone — vsluh.club (Free plan)
#
# Зона обслуживает два контура, оба живут на GitHub Pages и оба проксируются
# Cloudflare (оранжевое облако), который и терминирует TLS:
#   prod  vsluh.club / www.vsluh.club  -> sofiaboop/vsluh        (Pages, ветка main)
#   test  test.vsluh.club              -> sofiaboop/vsluh-test   (Pages, ветка gh-pages)
#
# Всё, что было заведено руками через API, импортировано в imports.tf —
# первый `terraform plan` обязан показывать ноль изменений.
#
# НЕ управляется здесь:
#   MX/SPF Namecheap  — почта заведена регистратором, в IaC не тянем (records.tf);
#   DNSSEC            — DS-запись согласована с реестром `.club` вручную, трогать
#                       её терраформом опасно (docs/incidents/2026-10-09-dnssec-tls.md);
#   GitHub Pages      — настройки репозиториев, не Cloudflare.
#
# Auth: провайдер читает CLOUDFLARE_API_TOKEN из .envrc (direnv).
#       Запускать из корня репозитория: set -a && source .envrc && set +a
# ─────────────────────────────────────────────────────────────────────────────

terraform {
  required_version = ">= 1.5"

  required_providers {
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "4.52.0"
    }
  }

  backend "s3" {
    bucket         = "vsluh-tfstate-661609779996"
    key            = "cloudflare/vsluh.club/terraform.tfstate"
    region         = "eu-central-1"
    dynamodb_table = "vsluh-tflock"
    encrypt        = true
  }
}

provider "cloudflare" {
  # api_token берётся из переменной окружения CLOUDFLARE_API_TOKEN (direnv / .envrc)
}
