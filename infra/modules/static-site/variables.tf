variable "project" {
  description = "Имя продукта, идёт в теги"
  type        = string
  default     = "vsluh"
}

variable "environment" {
  description = "Окружение: dev или prod"
  type        = string
}

variable "domain" {
  description = "Основной хост раздачи, он же CN сертификата"
  type        = string
}

variable "extra_domains" {
  description = "Дополнительные хосты той же раздачи (например www) — уходят в SAN сертификата и в aliases"
  type        = list(string)
  default     = []
}

variable "zone_name" {
  description = "Имя зоны Cloudflare, в которой живут хосты"
  type        = string
  default     = "vsluh.club"
}

variable "bucket_name" {
  description = "Имя S3-бакета со статикой"
  type        = string
}

variable "price_class" {
  description = "Ценовой класс CloudFront; PriceClass_100 — США и Европа, нам хватает: трафик всё равно приходит через Cloudflare"
  type        = string
  default     = "PriceClass_100"
}

variable "hsts_max_age" {
  description = "HSTS. Поднимаем лесенкой 300 -> 3600 -> 86400 -> 31536000 и только убедившись, что все хосты живут по HTTPS"
  type        = number
  default     = 300
}
