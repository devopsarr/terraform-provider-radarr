resource "radarr_metadata" "example" {
  enable          = true
  name            = "Example"
  implementation  = "MediaBrowserMetadata"
  config_contract = "MediaBrowserMetadataSettings"
  movie_metadata  = true
  tags            = [1, 2]
}