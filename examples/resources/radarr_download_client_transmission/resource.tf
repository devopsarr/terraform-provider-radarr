resource "radarr_download_client_transmission" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "transmission"
  url_base = "/transmission/"
  port     = 9091
}