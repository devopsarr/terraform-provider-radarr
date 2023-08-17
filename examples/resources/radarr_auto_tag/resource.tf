resource "radarr_auto_tag" "example" {
  name                      = "Example"
  remove_tags_automatically = true
  tags                      = [1, 2]

  specifications = [
    {
      name           = "folder"
      implementation = "RootFolderSpecification"
      negate         = true
      required       = false
      value          = "/series"
    },
    {
      name           = "year"
      implementation = "YearSpecification"
      negate         = true
      required       = false
      min            = 1900
      max            = 1910
    },
    {
      name           = "genre"
      implementation = "GenreSpecification"
      negate         = false
      required       = false
      value          = "horror comedy"
    },
  ]
}