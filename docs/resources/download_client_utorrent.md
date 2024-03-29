---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_download_client_utorrent Resource - terraform-provider-radarr"
subcategory: "Download Clients"
description: |-
  Download Client uTorrent resource.
  For more information refer to Download Client https://wiki.servarr.com/radarr/settings#download-clients and uTorrent https://wiki.servarr.com/radarr/supported#utorrent.
---

# radarr_download_client_utorrent (Resource)

<!-- subcategory:Download Clients -->
Download Client uTorrent resource.
For more information refer to [Download Client](https://wiki.servarr.com/radarr/settings#download-clients) and [uTorrent](https://wiki.servarr.com/radarr/supported#utorrent).

## Example Usage

```terraform
resource "radarr_download_client_utorrent" "example" {
  enable         = true
  priority       = 1
  name           = "Example"
  host           = "utorrent"
  url_base       = "/utorrent/"
  port           = 9091
  movie_category = "tv-radarr"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `intial_state` (Number) Initial state, with Stop support. `0` Start, `1` ForceStart, `2` Pause, `3` Stop.
- `movie_category` (String) Movie category.
- `movie_imported_category` (String) Movie imported category.
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
terraform import radarr_download_client_utorrent.example 1
```
