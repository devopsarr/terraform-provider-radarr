package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataRoksboxResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataRoksboxResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_roksbox.test", "movie_metadata", "false"),
					resource.TestCheckResourceAttrSet("radarr_metadata_roksbox.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataRoksboxResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataRoksboxResourceConfig("roksboxResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_roksbox.test", "movie_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata_roksbox.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataRoksboxResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata_roksbox" "test" {
		enable = false
		name = "%s"
		movie_metadata = %s
		movie_images = true
	}`, name, metadata)
}
