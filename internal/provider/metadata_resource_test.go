package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata.test", "movie_metadata", "true"),
					resource.TestCheckResourceAttrSet("radarr_metadata.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMetadataResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata.test", "movie_metadata", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata" "test" {
		enable = true
		name = "%s"
		implementation = "MediaBrowserMetadata"
    	config_contract = "MediaBrowserMetadataSettings"
		movie_metadata = %s
	}`, name, metadata)
}
