package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagConditionMonitoredDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionMonitoredDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_auto_tag_condition_monitored.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_auto_tag_condition_monitored.test", "name", "Test")),
			},
		},
	})
}

const testAccAutoTagConditionMonitoredDataSourceConfig = `
resource "radarr_tag" "test" {
	label = "atconditionmonitored"
}

data  "radarr_auto_tag_condition_monitored" "test" {
	name = "Test"
	negate = false
	required = false
}

resource "radarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSMonitored"

	tags = [radarr_tag.test.id]
	
	specifications = [data.radarr_auto_tag_condition_monitored.test]	
}`
