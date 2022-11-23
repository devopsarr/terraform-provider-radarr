package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_custom_format.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_custom_format.test", "include_custom_format_when_renaming", "false")),
			},
		},
	})
}

const testAccCustomFormatDataSourceConfig = `
resource "radarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "dataTest"
	
	specifications = [
		{
			name = "Surround Sound"
			implementation = "ReleaseTitleSpecification"
			negate = false
			required = false
			value = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
		},
		{
			name = "Arabic"
			implementation = "LanguageSpecification"
			negate = false
			required = false
			value = "31"
		}
	]	
}

data "radarr_custom_format" "test" {
	name = radarr_custom_format.test.name
}
`
