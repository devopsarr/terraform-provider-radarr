package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientUtorrentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientUtorrentResourceConfig("error", "utorrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientUtorrentResourceConfig("resourceUtorrentTest", "utorrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_utorrent.test", "host", "utorrent"),
					resource.TestCheckResourceAttr("radarr_download_client_utorrent.test", "url_base", "/utorrent/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_utorrent.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientUtorrentResourceConfig("error", "utorrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientUtorrentResourceConfig("resourceUtorrentTest", "utorrent-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_utorrent.test", "host", "utorrent-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_utorrent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientUtorrentResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_utorrent" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/utorrent/"
		port = 9091
		movie_category = "tv-radarr"
	}`, name, host)
}
