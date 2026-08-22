resource "radarr_notification_emby" "example" {
  on_grab                          = false
  on_download                      = true
  on_upgrade                       = true
  on_rename                        = false
  on_movie_added                   = false
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true
  on_health_issue                  = false
  on_application_update            = false

  include_health_warnings = false
  name                    = "Example"

  host    = "emby.lcl"
  port    = 8096
  api_key = "API_Key"

  # optional path mapping, for when Radarr and Emby/Jellyfin see the library at
  # different paths. only applied when update_library is enabled.
  update_library = true
  map_from       = "/movies"
  map_to         = "/media/movies"
}