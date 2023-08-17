data "radarr_auto_tag_condition_year" "example" {
  name     = "Example"
  negate   = false
  required = false
  min      = 1900
  max      = 2000
}

resource "radarr_custom_format" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.radarr_auto_tag_condition_year.example]
}