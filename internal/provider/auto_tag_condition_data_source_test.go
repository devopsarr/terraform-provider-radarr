package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagConditionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_auto_tag_condition.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_auto_tag_condition.test", "name", "year"),
					resource.TestCheckResourceAttr("radarr_auto_tag.test", "specifications.0.implementation", "YearSpecification")),
			},
		},
	})
}

const testAccAutoTagConditionDataSourceConfig = `
resource "radarr_tag" "test" {
	label = "atcondition"
}

data  "radarr_auto_tag_condition" "test" {
	name = "year"
	implementation = "YearSpecification"
	min = 1900
	max = 1920
	required = false
	negate = true
}

resource "radarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDS"
	tags = [radarr_tag.test.id]
	
	specifications = [data.radarr_auto_tag_condition.test]	
}`
