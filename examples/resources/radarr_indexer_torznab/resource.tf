resource "radarr_indexer_torznab" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://feed.animetosho.org"
  api_path                = "/nabapi"
  categories              = [2000, 2010]
  remove_year             = true
  minimum_seeders         = 1
}
