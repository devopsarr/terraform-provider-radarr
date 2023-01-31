data "radarr_custom_format_condition_indexer_flag" "example" {
  name     = "AHD_UserRelease"
  negate   = false
  required = false
  value    = "1024"
}

resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.radarr_custom_format_condition_indexer_flag.example]
}