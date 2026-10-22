data "cloudflare_zones" "this" {
  filter {
    name = var.zone_name
  }
}

# Сам REST API создаёт SAM (manifest.yaml), id не хардкодим.
data "aws_api_gateway_rest_api" "api" {
  name = "${var.project}-${var.environment}-api"
}

locals {
  zone_id = data.cloudflare_zones.this.zones[0].id
  host    = trimsuffix(var.domain, ".${var.zone_name}")
}
