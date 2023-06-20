package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatConditionQualityModifierDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccCustomFormatConditionQualityModifierDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_custom_format_condition_quality_modifier.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_custom_format_condition_quality_modifier.test", "name", "REMUX"),
					resource.TestCheckResourceAttr("radarr_custom_format.test", "specifications.0.value", "5")),
			},
		},
	})
}

const testAccCustomFormatConditionQualityModifierDataSourceConfig = `
data  "radarr_custom_format_condition_quality_modifier" "test" {
	name = "REMUX"
	negate = false
	required = false
	value = "5"
}

resource "radarr_custom_format" "test" {
	include_custom_format_when_renaming = false
	name = "TestWithDSQualityModifier"
	
	specifications = [data.radarr_custom_format_condition_quality_modifier.test]	
}`
