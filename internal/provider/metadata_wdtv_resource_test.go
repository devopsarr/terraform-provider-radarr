package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataWdtvResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataWdtvResourceConfig("wdtvResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_wdtv.test", "movie_metadata", "false"),
					resource.TestCheckResourceAttrSet("radarr_metadata_wdtv.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMetadataWdtvResourceConfig("wdtvResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_wdtv.test", "movie_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata_wdtv.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataWdtvResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata_wdtv" "test" {
		enable = false
		name = "%s"
		movie_metadata = %s
		movie_images = true
	}`, name, metadata)
}
