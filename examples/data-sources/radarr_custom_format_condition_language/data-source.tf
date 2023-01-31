data "radarr_custom_format_condition_language" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = "31"
}

resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.radarr_custom_format_condition_language.example]
}