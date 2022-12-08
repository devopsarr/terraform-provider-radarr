resource "radarr_download_client_utorrent" "example" {
  enable         = true
  priority       = 1
  name           = "Example"
  host           = "utorrent"
  url_base       = "/utorrent/"
  port           = 9091
  movie_category = "tv-radarr"
}