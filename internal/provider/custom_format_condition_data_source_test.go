package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCustomFormatConditionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_custom_format_condition.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_custom_format_condition.test", "name", "Surround Sound"),
					resource.TestCheckResourceAttr("radarr_custom_format.test", "specifications.0.implementation", "ReleaseTitleSpecification")),
			},
		},
	})
}

const testAccCustomFormatConditionDataSourceConfig = `
data  "radarr_custom_format_condition" "test" {
	name = "Surround Sound"
	implementation = "ReleaseTitleSpecification"
	negate = false
	required = false
	value = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
}

data  "radarr_custom_format_condition" "test1" {
	name = "Size"
	implementation = "SizeSpecification"
	negate = false
	required = false
	min = 0
	max = 100
}

resource "radarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDS"
	
	specifications = [data.radarr_custom_format_condition.test,data.radarr_custom_format_condition.test1]	
}`
