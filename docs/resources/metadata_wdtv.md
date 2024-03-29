---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_metadata_wdtv Resource - terraform-provider-radarr"
subcategory: "Metadata"
description: |-
  Metadata Wdtv resource.
  For more information refer to Metadata https://wiki.servarr.com/radarr/settings#metadata and WDTV https://wiki.servarr.com/radarr/supported#wdtvmetadata.
---

# radarr_metadata_wdtv (Resource)

<!-- subcategory:Metadata -->
Metadata Wdtv resource.
For more information refer to [Metadata](https://wiki.servarr.com/radarr/settings#metadata) and [WDTV](https://wiki.servarr.com/radarr/supported#wdtvmetadata).

## Example Usage

```terraform
resource "radarr_metadata_wdtv" "example" {
  enable         = true
  name           = "Example"
  movie_metadata = true
  movie_images   = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `movie_images` (Boolean) Movie images flag.
- `movie_metadata` (Boolean) Movie metadata flag.
- `name` (String) Metadata name.

### Optional

- `enable` (Boolean) Enable flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Metadata ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_metadata_wdtv.example 1
```
