---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_download_client_vuze Resource - terraform-provider-radarr"
subcategory: "Download Clients"
description: |-
  Download Client Vuze resource.
  For more information refer to Download Client https://wiki.servarr.com/radarr/settings#download-clients and Vuze https://wiki.servarr.com/radarr/supported#vuze.
---

# radarr_download_client_vuze (Resource)

<!-- subcategory:Download Clients -->
Download Client Vuze resource.
For more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients) and [Vuze](https://wiki.servarr.com/radarr/supported#vuze).

## Example Usage

```terraform
resource "radarr_download_client_vuze" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "vuze"
  url_base = "/vuze/"
  port     = 9091
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `movie_category` (String) Movie category.
- `movie_directory` (String) Movie directory.
- `older_movie_priority` (Number) Older Movie priority. `0` Last, `1` First.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `recent_movie_priority` (Number) Recent Movie priority. `0` Last, `1` First.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `tags` (Set of Number) List of associated tags.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_download_client_vuze.example 1
```
