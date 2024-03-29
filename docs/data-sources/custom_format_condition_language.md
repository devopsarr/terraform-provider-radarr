---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_custom_format_condition_language Data Source - terraform-provider-radarr"
subcategory: "Profiles"
description: |-
  Custom Format Condition Language data source.
  For more information refer to Custom Format Conditions https://wiki.servarr.com/radarr/settings#conditions.
---

# radarr_custom_format_condition_language (Data Source)

<!-- subcategory:Profiles -->
 Custom Format Condition Language data source.
For more information refer to [Custom Format Conditions](https://wiki.servarr.com/radarr/settings#conditions).

## Example Usage

```terraform
data "radarr_custom_format_condition_language" "example" {
  name     = "Example"
  negate   = false
  required = false
  value    = "31"
}

resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = false
  name                                = "Example"

  specifications = [data.radarr_custom_format_condition_language.example]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Specification name.
- `negate` (Boolean) Negate flag.
- `required` (Boolean) Computed flag.
- `value` (String) Language ID.

### Read-Only

- `id` (Number) Custom format condition language ID.
- `implementation` (String) Implementation.
