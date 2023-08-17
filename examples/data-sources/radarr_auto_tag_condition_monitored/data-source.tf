data "radarr_auto_tag_condition_monitored" "example" {
  name     = "Example"
  negate   = false
  required = false
}

resource "radarr_custom_format" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.radarr_auto_tag_condition_monitored.example]
}