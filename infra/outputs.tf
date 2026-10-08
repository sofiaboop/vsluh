output "site_domain" {
  value = local.site_domain
}

output "site_bucket" {
  value = module.site.bucket_name
}

output "site_distribution_id" {
  description = "Положить в GitHub Environment как CLOUDFRONT_DISTRIBUTION_ID"
  value       = module.site.distribution_id
}

output "miniapp_domain" {
  value = local.miniapp_domain
}

output "miniapp_bucket" {
  value = module.miniapp.bucket_name
}

output "miniapp_distribution_id" {
  description = "Положить в GitHub Environment как MINIAPP_CLOUDFRONT_DISTRIBUTION_ID"
  value       = module.miniapp.distribution_id
}
