resource "radarr_indexer_filelist" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://filelist.io"
  username                = "User"
  passkey                 = "PassKey"
  minimum_seeders         = 1
  categories              = [4, 6, 1]
  required_flags          = [1, 4]
}
