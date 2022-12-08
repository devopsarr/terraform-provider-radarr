resource "radarr_indexer_pass_the_popcorn" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://passthepopcorn.me"
  api_user                = "User"
  api_key                 = "Key"
  minimum_seeders         = 1
}
