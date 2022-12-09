resource "radarr_notification_plex" "example" {
  on_download                      = true
  on_upgrade                       = true
  on_rename                        = false
  on_movie_added                   = false
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true

  include_health_warnings = false
  name                    = "Example"

  host       = "plex.lcl"
  port       = 32400
  auth_token = "AuthTOKEN"
}