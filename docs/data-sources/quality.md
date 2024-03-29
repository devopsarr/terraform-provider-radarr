---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "radarr_quality Data Source - terraform-provider-radarr"
subcategory: "Profiles"
description: |-
  Single Quality.
---

# radarr_quality (Data Source)

<!-- subcategory:Profiles -->
Single Quality.

## Example Usage

```terraform
data "radarr_quality" "bluray" {
  name = "Bluray-2160p"
}

data "radarr_quality" "webdl" {
  name = "WEBDL-2160p"
}

data "radarr_quality" "webrip" {
  name = "WEBRip-2160p"
}

resource "radarr_quality_profile" "Example" {
  name            = "Example"
  upgrade_allowed = true
  cutoff          = 2000

  language = data.radarr_language.test

  quality_groups = [
    {
      id   = 2000
      name = "WEB 2160p"
      qualities = [
        data.radarr_quality.webdl,
        data.radarr_quality.webrip,
      ]
    },
    {
      qualities = [data.radarr_quality.bluray]
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Quality Name.

### Read-Only

- `id` (Number) Quality  ID.
- `resolution` (Number) Quality Resolution.
- `source` (String) Quality source.
