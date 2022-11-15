package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRestrictionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccRestrictionResourceConfig("test1", "test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_restriction.test", "ignored", "test1"),
					resource.TestCheckResourceAttrSet("radarr_restriction.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccRestrictionResourceConfig("test3", "test2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_restriction.test", "ignored", "test3"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_restriction.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRestrictionResourceConfig(ignore, require string) string {
	return fmt.Sprintf(`
		resource "radarr_restriction" "test" {
  			ignored = "%s"
			required = "%s"
		}
	`, ignore, require)
}
