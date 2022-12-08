resource "radarr_indexer_torrent_potato" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "http://127.0.0.1"
  user                    = "User"
  passkey                 = "Key"
  minimum_seeders         = 1
}
