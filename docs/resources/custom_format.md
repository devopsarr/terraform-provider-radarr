---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_custom_format Resource - terraform-provider-radarr"
subcategory: "Profiles"
description: |-
  Custom Format resource.
  For more information refer to Custom Format https://wiki.servarr.com/radarr/settings#custom-formats.
---

# radarr_custom_format (Resource)

<!-- subcategory:Profiles -->Custom Format resource.
For more information refer to [Custom Format](https://wiki.servarr.com/radarr/settings#custom-formats).

## Example Usage

```terraform
resource "radarr_custom_format" "example" {
  include_custom_format_when_renaming = true
  name                                = "Example"

  specifications = [
    {
      name           = "Surround Sound"
      implementation = "ReleaseTitleSpecification"
      negate         = false
      required       = false
      value          = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
    },
    {
      name           = "Arabic"
      implementation = "LanguageSpecification"
      negate         = false
      required       = false
      value          = "31"
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Custom Format name.
- `specifications` (Attributes Set) Specifications. (see [below for nested schema](#nestedatt--specifications))

### Optional

- `include_custom_format_when_renaming` (Boolean) Include custom format when renaming flag.

### Read-Only

- `id` (Number) Custom Format ID.

<a id="nestedatt--specifications"></a>
### Nested Schema for `specifications`

Optional:

- `implementation` (String) Implementation.
- `max` (Number) Max.
- `min` (Number) Min.
- `name` (String) Specification name.
- `negate` (Boolean) Negate flag.
- `required` (Boolean) Required flag.
- `value` (String) Value.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import radarr_custom_format.example 1
```