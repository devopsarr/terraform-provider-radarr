package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccQualityProfileDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_quality_profile.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_quality_profile.test", "cutoff", "1003")),
			},
		},
	})
}

const testAccQualityProfileDataSourceConfig = `
resource "radarr_quality_profile" "test" {
	name            = "qpdata"
	upgrade_allowed = true
	cutoff          = 1003

	language = {
		id   = 1
		name = "English"
	}

	quality_groups = [
		{
			id   = 1003
			name = "WEB 2160p"
			qualities = [
				{
					id         = 18
					name       = "WEBDL-2160p"
					source     = "webdl"
					resolution = 2160
				},
				{
					id         = 17
					name       = "WEBRip-2160p"
					source     = "webrip"
					resolution = 2160
				}
			]
		}
	]
}

data "radarr_quality_profile" "test" {
	name = radarr_quality_profile.test.name
}
`
