package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataEmbyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataEmbyResourceConfig("embyResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_emby.test", "movie_metadata", "false"),
					resource.TestCheckResourceAttrSet("radarr_metadata_emby.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMetadataEmbyResourceConfig("embyResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_emby.test", "movie_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata_emby.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataEmbyResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata_emby" "test" {
		enable = false
		name = "%s"
		movie_metadata = %s
	}`, name, metadata)
}
