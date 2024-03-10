provider "radarr" {
  url     = "http://example.radarr.tv:8989"
  api_key = "APIkey-example"
  extra_headers = [
    {
      name  = "exampleName"
      value = "exanpleValue"
    }
  ]
}
