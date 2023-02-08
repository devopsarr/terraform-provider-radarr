package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRestrictionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRestrictionDataSourceConfig("radarr_restriction.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_restriction.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_restriction.test", "ignored", "datatest1")),
			},
			// Not found testing
			{
				Config:      testAccRestrictionDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find restriction"),
			},
		},
	})
}

func testAccRestrictionDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	resource "radarr_restriction" "test" {
		ignored = "datatest1"
		required = "datatest2"
	}
	
	data "radarr_restriction" "test" {
		id = %s
	}
	`, id)
}
