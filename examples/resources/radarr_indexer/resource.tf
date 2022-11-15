resource "radarr_indexer" "example" {
  enable_automatic_search = true
  name                    = "Test"
  implementation          = "Newznab"
  protocol                = "usenet"
  config_contract         = "NewznabSettings"
  base_url                = "https://lolo.sickbeard.com"
  api_path                = "/api"
  categories              = [8000, 5000]
  tags                    = [1, 2]
}