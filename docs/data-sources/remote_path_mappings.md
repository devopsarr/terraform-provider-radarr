---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_remote_path_mappings Data Source - terraform-provider-radarr"
subcategory: "Download Clients"
description: |-
  List all available Remote Path Mappings ../resources/remote_path_mapping.
---

# radarr_remote_path_mappings (Data Source)

<!-- subcategory:Download Clients -->
List all available [Remote Path Mappings](../resources/remote_path_mapping).

## Example Usage

```terraform
data "radarr_remote_path_mappings" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `remote_path_mappings` (Attributes Set) Remote Path Mapping list. (see [below for nested schema](#nestedatt--remote_path_mappings))

<a id="nestedatt--remote_path_mappings"></a>
### Nested Schema for `remote_path_mappings`

Read-Only:

- `host` (String) Download Client host.
- `id` (Number) RemotePathMapping ID.
- `local_path` (String) Local path.
- `remote_path` (String) Download Client remote path.
