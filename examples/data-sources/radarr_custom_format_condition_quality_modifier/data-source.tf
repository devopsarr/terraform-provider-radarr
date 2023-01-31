data "radarr_custom_format_condition_quality_modifier" "example" {
  name     = "REMUX"
  negate   = false
  required = false
  value    = "5"
}

resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.radarr_custom_format_condition_quality_modifier.example]
}