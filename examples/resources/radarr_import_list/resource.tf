resource "radarr_import_list" "example" {
  enabled              = false
  enable_auto          = true
  search_on_add        = false
  monitor              = "movieOnly"
  minimum_availability = "tba"
  list_type            = "program"
  root_folder_path     = radarr_root_folder.example.path
  quality_profile_id   = radarr_quality_profile.example.id
  name                 = "Example"
  implementation       = "RadarrImport"
  config_contract      = "RadarrSettings"
  tags                 = [1, 2]

  tag_ids     = [1, 2]
  profile_ids = [1]
  base_url    = "http://127.0.0.1:8686"
  api_key     = "APIKey"
}