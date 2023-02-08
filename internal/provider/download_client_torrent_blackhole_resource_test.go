package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientTorrentBlackholeResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientTorrentBlackholeResourceConfig("error", ".torrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientTorrentBlackholeResourceConfig("resourceTorrentBlackholeTest", ".torrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_torrent_blackhole.test", "magnet_file_extension", ".torrent"),
					resource.TestCheckResourceAttr("radarr_download_client_torrent_blackhole.test", "watch_folder", "/config/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_torrent_blackhole.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientTorrentBlackholeResourceConfig("error", ".torrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientTorrentBlackholeResourceConfig("resourceTorrentBlackholeTest", ".magnet"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_torrent_blackhole.test", "magnet_file_extension", ".magnet"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_torrent_blackhole.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientTorrentBlackholeResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_torrent_blackhole" "test" {
		enable = false
		priority = 1
		name = "%s"
		magnet_file_extension = "%s"
		watch_folder = "/config/"
		torrent_folder = "/config/"
	}`, name, host)
}
