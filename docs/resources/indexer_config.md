---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_indexer_config Resource - terraform-provider-radarr"
subcategory: "Indexers"
description: |-
  Indexer Config resource.
  For more information refer to Indexer https://wiki.servarr.com/radarr/settings#options documentation.
---

# radarr_indexer_config (Resource)

<!-- subcategory:Indexers -->Indexer Config resource.
For more information refer to [Indexer](https://wiki.servarr.com/radarr/settings#options) documentation.

## Example Usage

```terraform
resource "radarr_indexer_config" "example" {
  maximum_size               = 0
  minimum_age                = 0
  retention                  = 0
  rss_sync_interval          = 25
  availability_delay         = 0
  whitelisted_hardcoded_subs = ""
  prefer_indexer_flags       = false
  allow_hardcoded_subs       = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `allow_hardcoded_subs` (Boolean) Allow hardcoded subs.
- `availability_delay` (Number) Availability delay.
- `maximum_size` (Number) Maximum size.
- `minimum_age` (Number) Minimum age.
- `prefer_indexer_flags` (Boolean) Prefer indexer flags.
- `retention` (Number) Retention.
- `rss_sync_interval` (Number) RSS sync interval.
- `whitelisted_hardcoded_subs` (String) Whitelisted hardconded subs.

### Read-Only

- `id` (Number) Indexer Config ID.

## Import

Import is supported using the following syntax:

```shell
# import does not need parameters
terraform import radarr_indexer_config.example
```