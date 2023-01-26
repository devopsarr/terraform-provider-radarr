resource "radarr_metadata_kodi" "example" {
  enable         = true
  name           = "Example"
  movie_metadata = true
  movie_images   = true
}