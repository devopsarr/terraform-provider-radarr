resource "radarr_indexer_newznab" "example" {
  enable_automatic_search = true
  name                    = "Test"
  base_url                = "https://lolo.sickbeard.com"
  api_path                = "/api"
  categories              = [8000, 5000]
  tags                    = [1, 2]
}