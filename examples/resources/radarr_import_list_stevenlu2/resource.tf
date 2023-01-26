resource "radarr_import_list_stevenlu2" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  monitor              = "none"
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  source               = 0
  min_score            = 5
}