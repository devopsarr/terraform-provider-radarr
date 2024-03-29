---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_delay_profiles Data Source - terraform-provider-radarr"
subcategory: "Profiles"
description: |-
  List all available Delay Profiles ../resources/delay_profile.
---

# radarr_delay_profiles (Data Source)

<!-- subcategory:Profiles -->
List all available [Delay Profiles](../resources/delay_profile).

## Example Usage

```terraform
data "radarr_delay_profiles" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `delay_profiles` (Attributes Set) Delay Profile list. (see [below for nested schema](#nestedatt--delay_profiles))
- `id` (String) The ID of this resource.

<a id="nestedatt--delay_profiles"></a>
### Nested Schema for `delay_profiles`

Read-Only:

- `bypass_if_highest_quality` (Boolean) Bypass for highest quality Flag.
- `enable_torrent` (Boolean) Torrent allowed Flag.
- `enable_usenet` (Boolean) Usenet allowed Flag.
- `id` (Number) Delay Profile ID.
- `order` (Number) Order.
- `preferred_protocol` (String) Preferred protocol.
- `tags` (Set of Number) List of associated tags.
- `torrent_delay` (Number) Torrent Delay.
- `usenet_delay` (Number) Usenet delay.
