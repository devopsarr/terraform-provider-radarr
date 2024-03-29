---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_import_list_config Resource - terraform-provider-radarr"
subcategory: "Import Lists"
description: |-
  Import List Config resource.
  For more information refer to Import List https://wiki.servarr.com/radarr/settings#completed-download-handling documentation.
---

# radarr_import_list_config (Resource)

<!-- subcategory:Import Lists -->
Import List Config resource.
For more information refer to [Import List](https://wiki.servarr.com/radarr/settings#completed-download-handling) documentation.

## Example Usage

```terraform
resource "radarr_import_list_config" "example" {
  sync_level = "logOnly"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `sync_level` (String) Clean library level.

### Read-Only

- `id` (Number) Import List Config ID.

## Import

Import is supported using the following syntax:

```shell
# import does not need parameters
terraform import radarr_import_list_config.example ""
```
