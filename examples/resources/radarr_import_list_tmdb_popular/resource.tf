resource "radarr_import_list_tmdb_popular" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  monitor              = "none"
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  tmdb_list_type       = 2
  min_vote_average     = "5"
  min_votes            = "1"
  tmdb_certification   = "PG-13"
  language_code        = 2
}