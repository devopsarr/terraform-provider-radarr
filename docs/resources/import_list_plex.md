---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_import_list_plex Resource - terraform-provider-radarr"
subcategory: "Import Lists"
description: |-
  Import List Plex resource.
  For more information refer to Import List https://wiki.servarr.com/radarr/settings#import-lists and Plex https://wiki.servarr.com/radarr/supported#pleximport.
---

# radarr_import_list_plex (Resource)

<!-- subcategory:Import Lists -->
Import List Plex resource.
For more information refer to [Import List](https://wiki.servarr.com/radarr/settings#import-lists) and [Plex](https://wiki.servarr.com/radarr/supported#pleximport).

## Example Usage

```terraform
resource "radarr_import_list_plex" "example" {
  enabled              = true
  enable_auto          = false
  search_on_add        = false
  root_folder_path     = "/config"
  monitor              = "none"
  minimum_availability = "tba"
  quality_profile_id   = 1
  name                 = "Example"
  access_token         = "YourToken"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_token` (String, Sensitive) Access token.
- `minimum_availability` (String) Minimum availability.
- `monitor` (String) Should monitor.
- `name` (String) Import List name.
- `quality_profile_id` (Number) Quality profile ID.
- `root_folder_path` (String) Root folder path.

### Optional

- `enable_auto` (Boolean) Enable automatic add flag.
- `enabled` (Boolean) Enabled flag.
- `list_order` (Number) List order.
- `search_on_add` (Boolean) Search on add flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Import List ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_import_list_plex.example 1
```
