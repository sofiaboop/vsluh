variable "zone_name" {
  description = "Имя зоны Cloudflare; zone id резолвится через data.cloudflare_zones"
  type        = string
  default     = "vsluh.club"
}

variable "pages_host" {
  description = "Хост GitHub Pages аккаунта — цель CNAME для www и test"
  type        = string
  default     = "sofiaboop.github.io"
}

variable "pages_ipv4" {
  description = "Анонсируемые GitHub Pages адреса для apex-записи (A на домене второго уровня, CNAME там невозможен)"
  type        = list(string)
  default = [
    "185.199.108.153",
    "185.199.109.153",
    "185.199.110.153",
    "185.199.111.153",
  ]
}

variable "test_host" {
  description = "Хост тестового контура; артефакт отдаёт репозиторий sofiaboop/vsluh-test"
  type        = string
  default     = "test.vsluh.club"
}
