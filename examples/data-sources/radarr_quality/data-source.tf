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