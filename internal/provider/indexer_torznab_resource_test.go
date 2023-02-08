package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerTorznabResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerTorznabResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerTorznabResourceConfig("torznabResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_torznab.test", "minimum_seeders", "1"),
					resource.TestCheckResourceAttrSet("radarr_indexer_torznab.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerTorznabResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerTorznabResourceConfig("torznabResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_torznab.test", "minimum_seeders", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_torznab.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerTorznabResourceConfig(name string, seeders int) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_torznab" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "https://feed.animetosho.org"
		api_path = "/nabapi"
		minimum_seeders = %d
		categories = [2000,2010]
		remove_year = true
	}`, name, seeders)
}
