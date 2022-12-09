resource "radarr_notification_ntfy" "example" {
  on_grab                          = false
  on_download                      = true
  on_upgrade                       = true
  on_movie_added                   = false
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true
  on_health_issue                  = false
  on_application_update            = false

  include_health_warnings = false
  name                    = "Example"

  priority   = 1
  server_url = "https://ntfy.sh"
  username   = "User"
  password   = "Pass"
  topics     = ["Topic1234", "Topic4321"]
  field_tags = ["warning", "skull"]
}