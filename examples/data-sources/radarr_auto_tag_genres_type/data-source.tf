data "radarr_auto_tag_condition_genres" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = "horror comedy"
}

resource "radarr_custom_format" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.radarr_auto_tag_condition_genres.example]
}