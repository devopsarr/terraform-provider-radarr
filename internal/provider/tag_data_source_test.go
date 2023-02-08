package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccTagDataSourceConfig("radarr_tag.test.label"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_tag.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_tag.test", "label", "tag_datasource"),
				),
			},
			// Not found testing
			{
				Config:      testAccTagDataSourceConfig("\"error\""),
				ExpectError: regexp.MustCompile("Unable to find tag"),
			},
		},
	})
}

func testAccTagDataSourceConfig(label string) string {
	return fmt.Sprintf(`
	resource "radarr_tag" "test" {
		label = "tag_datasource"
	}
	
	data "radarr_tag" "test" {
		label = %s
	}
	`, label)
}
