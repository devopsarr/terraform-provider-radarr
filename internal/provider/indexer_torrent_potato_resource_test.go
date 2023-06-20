package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerTorrentPotatoResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerTorrentPotatoResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerTorrentPotatoResourceConfig("torrentPotatoResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_torrent_potato.test", "minimum_seeders", "1"),
					resource.TestCheckResourceAttrSet("radarr_indexer_torrent_potato.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerTorrentPotatoResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerTorrentPotatoResourceConfig("torrentPotatoResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_torrent_potato.test", "minimum_seeders", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_torrent_potato.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerTorrentPotatoResourceConfig(name string, seeders int) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_torrent_potato" "test" {
		enable_automatic_search = false
		name = "%s"
		base_url = "http://127.0.0.1"
		user = "testUser"
		passkey = "pass"
		minimum_seeders = %d
		seed_time = 10
		seed_ratio = 0.5
	}`, name, seeders)
}
