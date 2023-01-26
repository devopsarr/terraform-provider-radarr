resource "radarr_import_list_couch_potato" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  monitor              = "none"
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  link                 = "http://localhost"
  api_key              = "APIKey"
  port                 = 5050
  only_active          = true
}