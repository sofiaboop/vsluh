# ─────────────────────────────────────────────────────────────────────────────
# Импорт того, что было заведено руками через API до появления Terraform.
#
# Id записей сняты живьём 21 окт 2026:
#   A     vsluh.club      185.199.108.153  d4aa16aa436f7e65b1cd8762a86bb109
#   A     vsluh.club      185.199.109.153  ee8c828b613fee1ef95594b081395fe9
#   A     vsluh.club      185.199.110.153  30cf531f1aea354aaa8f73f9ef6cdfa8
#   A     vsluh.club      185.199.111.153  c33231ad256c9395dcab98d790c38bfc
#   CNAME www.vsluh.club                   c2c9e1c203065b5dde30f6a7feacd7f8
#   CNAME test.vsluh.club                  aa459a1b71d60a33cb0ade4a39f70229
#
# ⚠️ Перед apply: `terraform plan` обязан показать 6 импортов и ноль изменений.
#    Любой diff по DNS — это расхождение кода с реальностью, а не повод жать apply.
# ─────────────────────────────────────────────────────────────────────────────

import {
  to = cloudflare_record.apex["185.199.108.153"]
  id = "${local.zone_id}/d4aa16aa436f7e65b1cd8762a86bb109"
}

import {
  to = cloudflare_record.apex["185.199.109.153"]
  id = "${local.zone_id}/ee8c828b613fee1ef95594b081395fe9"
}

import {
  to = cloudflare_record.apex["185.199.110.153"]
  id = "${local.zone_id}/30cf531f1aea354aaa8f73f9ef6cdfa8"
}

import {
  to = cloudflare_record.apex["185.199.111.153"]
  id = "${local.zone_id}/c33231ad256c9395dcab98d790c38bfc"
}

import {
  to = cloudflare_record.www
  id = "${local.zone_id}/c2c9e1c203065b5dde30f6a7feacd7f8"
}

import {
  to = cloudflare_record.test
  id = "${local.zone_id}/aa459a1b71d60a33cb0ade4a39f70229"
}

# cloudflare_zone_settings_override импорта не поддерживает: ресурс не создаёт
# сущность, а перезаписывает настройки зоны. Первый apply просто приводит их
# к значениям из settings.tf — они уже такие, поведение зоны не меняется.
