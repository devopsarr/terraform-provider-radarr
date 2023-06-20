package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRestrictionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRestrictionDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccRestrictionDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find restriction"),
			},
			// Read testing
			{
				Config: testAccRestrictionResourceConfig("datatest1", "datatest2") + testAccRestrictionDataSourceConfig("radarr_restriction.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_restriction.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_restriction.test", "ignored", "datatest1")),
			},
		},
	})
}

func testAccRestrictionDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "radarr_restriction" "test" {
		id = %s
	}
	`, id)
}
