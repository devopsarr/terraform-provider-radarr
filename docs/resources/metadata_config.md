---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_metadata_config Resource - terraform-provider-radarr"
subcategory: "Metadata"
description: |-
  Metadata Config resource.
  For more information refer to Metadata https://wiki.servarr.com/radarr/settings#options documentation.
---

# radarr_metadata_config (Resource)

<!-- subcategory:Metadata -->Metadata Config resource.
For more information refer to [Metadata](https://wiki.servarr.com/radarr/settings#options) documentation.

## Example Usage

```terraform
resource "radarr_metadata_config" "example" {
  certification_country = "us"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `certification_country` (String) Certification Country.

### Read-Only

- `id` (Number) Metadata Config ID.

## Import

Import is supported using the following syntax:

```shell
# import does not need parameters
terraform import radarr_metadata_config.example
```