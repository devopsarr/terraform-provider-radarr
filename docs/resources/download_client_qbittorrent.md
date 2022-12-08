---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_download_client_qbittorrent Resource - terraform-provider-radarr"
subcategory: "Download Clients"
description: |-
  Download Client qBittorrent resource.
  For more information refer to Download Client https://wiki.servarr.com/radarr/settings#download-clients and qBittorrent https://wiki.servarr.com/radarr/supported#qbittorrent.
---

# radarr_download_client_qbittorrent (Resource)

<!-- subcategory:Download Clients -->Download Client qBittorrent resource.
For more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients) and [qBittorrent](https://wiki.servarr.com/radarr/supported#qbittorrent).

## Example Usage

```terraform
resource "radarr_download_client_qbittorrent" "example" {
  enable         = true
  priority       = 1
  name           = "Example"
  host           = "qbittorrent"
  url_base       = "/qbittorrent/"
  port           = 9091
  movie_category = "tv-radarr"
  first_and_last = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `enable` (Boolean) Enable flag.
- `first_and_last` (Boolean) First and last flag.
- `host` (String) host.
- `initial_state` (Number) Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause.
- `movie_category` (String) TV category.
- `movie_directory` (String) TV directory.
- `movie_imported_category` (String) TV imported category.
- `older_movie_priority` (Number) Older TV priority. `0` Last, `1` First.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `recent_movie_priority` (Number) Recent TV priority. `0` Last, `1` First.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `sequential_order` (Boolean) Sequential order flag.
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
terraform import radarr_download_client_qbittorrent.example 1
```