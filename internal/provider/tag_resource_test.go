package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "eng"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_tag.test", "label", "eng"),
					resource.TestCheckResourceAttrSet("radarr_tag.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "1080p"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_tag.test", "label", "1080p"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccTagResourceConfig(name, label string) string {
	return fmt.Sprintf(`
		resource "radarr_tag" "%s" {
  			label = "%s"
		}
	`, name, label)
}
