---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_download_client_freebox Resource - terraform-provider-radarr"
subcategory: "Download Clients"
description: |-
  Download Client Freebox resource.
  For more information refer to Download Client https://wiki.servarr.com/radarr/settings#download-clients and Freebox https://wiki.servarr.com/radarr/supported#torrentfreeboxdownload.
---

# radarr_download_client_freebox (Resource)

<!-- subcategory:Download Clients -->
Download Client Freebox resource.
For more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients) and [Freebox](https://wiki.servarr.com/radarr/supported#torrentfreeboxdownload).

## Example Usage

```terraform
resource "radarr_download_client_freebox" "example" {
  enable    = true
  priority  = 1
  name      = "Example"
  host      = "mafreebox.freebox.fr"
  api_url   = "/api/v1/"
  port      = 443
  app_id    = "freebox"
  app_token = "Token123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_url` (String) API URL.
- `app_id` (String) App ID.
- `app_token` (String, Sensitive) App Token.
- `host` (String) host.
- `name` (String) Download Client name.
- `port` (Number) Port.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `category` (String) category.
- `destination_directory` (String) Movie directory.
- `enable` (Boolean) Enable flag.
- `older_priority` (Number) Older Movie priority. `0` Last, `1` First.
- `priority` (Number) Priority.
- `recent_priority` (Number) Recent Movie priority. `0` Last, `1` First.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `tags` (Set of Number) List of associated tags.
- `use_ssl` (Boolean) Use SSL flag.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_download_client_freebox.example 1
```
