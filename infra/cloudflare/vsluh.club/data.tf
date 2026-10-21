# Zone id резолвим по имени — хардкода id в коде нет (конвенция репозитория).
data "cloudflare_zones" "this" {
  filter {
    name = var.zone_name
  }
}

locals {
  zone_id = data.cloudflare_zones.this.zones[0].id
}
