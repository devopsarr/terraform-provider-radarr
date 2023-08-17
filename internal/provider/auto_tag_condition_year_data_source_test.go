package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagConditionYearDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionYearDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_auto_tag_condition_year.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_auto_tag_condition_year.test", "name", "Test"),
					resource.TestCheckResourceAttr("radarr_auto_tag.test", "specifications.0.min", "1900")),
			},
		},
	})
}

const testAccAutoTagConditionYearDataSourceConfig = `
resource "radarr_tag" "test" {
	label = "atconditiontype"
}

data  "radarr_auto_tag_condition_year" "test" {
	name = "Test"
	negate = false
	required = false
	min = 1900
	max = 1910
}

resource "radarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSYear"

	tags = [radarr_tag.test.id]
	
	specifications = [data.radarr_auto_tag_condition_year.test]	
}`
