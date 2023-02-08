package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataKodiResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataKodiResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_kodi.test", "movie_metadata", "false"),
					resource.TestCheckResourceAttrSet("radarr_metadata_kodi.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataKodiResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataKodiResourceConfig("kodiResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_kodi.test", "movie_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata_kodi.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataKodiResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata_kodi" "test" {
		enable = false
		name = "%s"
		movie_metadata = %s
		movie_images = true
		movie_metadata_language = -2
		movie_metadata_url = false
		use_movie_nfo = true
		add_collection_name = false
	}`, name, metadata)
}
