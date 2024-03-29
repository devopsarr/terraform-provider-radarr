---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_notification_kodi Resource - terraform-provider-radarr"
subcategory: "Notifications"
description: |-
  Notification Kodi resource.
  For more information refer to Notification https://wiki.servarr.com/radarr/settings#connect and Kodi https://wiki.servarr.com/radarr/supported#xbmc.
---

# radarr_notification_kodi (Resource)

<!-- subcategory:Notifications -->
Notification Kodi resource.
For more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Kodi](https://wiki.servarr.com/radarr/supported#xbmc).

## Example Usage

```terraform
resource "radarr_notification_kodi" "example" {
  on_grab                          = false
  on_download                      = false
  on_upgrade                       = false
  on_rename                        = false
  on_movie_added                   = false
  on_movie_delete                  = false
  on_movie_file_delete             = false
  on_movie_file_delete_for_upgrade = true
  on_health_issue                  = false
  on_application_update            = false

  include_health_warnings = false
  name                    = "Example"

  host     = "http://kodi.com"
  port     = 8080
  username = "User"
  password = "MyPass"
  notify   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `host` (String) Host.
- `name` (String) NotificationKodi name.
- `on_movie_delete` (Boolean) On movie delete flag.
- `port` (Number) Port.

### Optional

- `always_update` (Boolean) Always update flag.
- `clean_library` (Boolean) Clean library flag.
- `display_time` (Number) Display time.
- `include_health_warnings` (Boolean) Include health warnings.
- `notify` (Boolean) Notification flag.
- `on_application_update` (Boolean) On application update flag.
- `on_download` (Boolean) On download flag.
- `on_grab` (Boolean) On grab flag.
- `on_health_issue` (Boolean) On health issue flag.
- `on_health_restored` (Boolean) On health restored flag.
- `on_manual_interaction_required` (Boolean) On manual interaction required flag.
- `on_movie_added` (Boolean) On movie added flag.
- `on_movie_file_delete` (Boolean) On movie file delete flag.
- `on_movie_file_delete_for_upgrade` (Boolean) On movie file delete for upgrade flag.
- `on_rename` (Boolean) On rename flag.
- `on_upgrade` (Boolean) On upgrade flag.
- `password` (String, Sensitive) Password.
- `tags` (Set of Number) List of associated tags.
- `update_library` (Boolean) Update library flag.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_notification_kodi.example 1
```
