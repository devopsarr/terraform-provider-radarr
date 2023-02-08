package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerRarbgResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerRarbgResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerRarbgResourceConfig("rarbgResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_rarbg.test", "enable_automatic_search", "false"),
					resource.TestCheckResourceAttr("radarr_indexer_rarbg.test", "base_url", "https://torrentapi.org"),
					resource.TestCheckResourceAttrSet("radarr_indexer_rarbg.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerRarbgResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerRarbgResourceConfig("rarbgResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_rarbg.test", "enable_automatic_search", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_rarbg.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerRarbgResourceConfig(name, aSearch string) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_rarbg" "test" {
		enable_automatic_search = %s
		name = "%s"
		base_url = "https://torrentapi.org"
		ranked_only = "false"
		minimum_seeders = 1
	}`, aSearch, name)
}
