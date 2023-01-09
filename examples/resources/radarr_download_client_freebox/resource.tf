resource "radarr_download_client_freebox" "example" {
  enable    = true
  priority  = 1
  name      = "Example"
  host      = "mafreebox.freebox.fr"
  api_url   = "/api/v1/"
  port      = 443
  app_id    = "freebox"
  app_token = "Token123"
}