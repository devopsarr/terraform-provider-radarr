---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_notification_slack Resource - terraform-provider-radarr"
subcategory: "Notifications"
description: |-
  Notification Slack resource.
  For more information refer to Notification https://wiki.servarr.com/radarr/settings#connect and Slack https://wiki.servarr.com/radarr/supported#slack.
---

# radarr_notification_slack (Resource)

<!-- subcategory:Notifications -->
Notification Slack resource.
For more information refer to [Notification](https://wiki.servarr.com/radarr/settings#connect) and [Slack](https://wiki.servarr.com/radarr/supported#slack).

## Example Usage

```terraform
resource "radarr_notification_slack" "example" {
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

  web_hook_url = "http://my.slack.com/test"
  username     = "user"
  channel      = "example-channel"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) NotificationSlack name.
- `on_movie_delete` (Boolean) On movie delete flag.
- `username` (String) Username.
- `web_hook_url` (String) URL.

### Optional

- `channel` (String) Channel.
- `icon` (String) Icon.
- `include_health_warnings` (Boolean) Include health warnings.
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
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_notification_slack.example 1
```
