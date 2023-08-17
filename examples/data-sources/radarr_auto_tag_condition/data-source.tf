data "radarr_auto_tag_condition" "example" {
  name           = "Example"
  implementation = "RootFolderSpecification"
  negate         = false
  required       = false
  value          = "/movies"
}

resource "radarr_auto_tag" "example" {
  remove_tags_automatically = false
  name                      = "Example"

  tags = [1, 2]

  specifications = [data.radarr_auto_tag.example]
}