resource "radarr_naming" "example" {
  include_quality            = false
  rename_movies              = true
  replace_illegal_characters = false
  replace_spaces             = false
  colon_replacement_format   = "%s"
  standard_movie_format      = "{Movie Title} ({Release Year}) {Quality Full}"
  movie_folder_format        = "{Movie Title} ({Release Year})"
}
