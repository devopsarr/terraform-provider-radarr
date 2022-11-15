resource "radarr_indexer_config" "example" {
  maximum_size               = 0
  minimum_age                = 0
  retention                  = 0
  rss_sync_interval          = 25
  availability_delay         = 0
  whitelisted_hardcoded_subs = ""
  prefer_indexer_flags       = false
  allow_hardcoded_subs       = false
}