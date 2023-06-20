package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientRtorrentResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientRtorrentResourceConfig("error", "rtorrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientRtorrentResourceConfig("resourceRtorrentTest", "rtorrent"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_rtorrent.test", "host", "rtorrent"),
					resource.TestCheckResourceAttr("radarr_download_client_rtorrent.test", "url_base", "/rtorrent/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_rtorrent.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientRtorrentResourceConfig("error", "rtorrent") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientRtorrentResourceConfig("resourceRtorrentTest", "rtorrent-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_rtorrent.test", "host", "rtorrent-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_rtorrent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientRtorrentResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_rtorrent" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/rtorrent/"
		port = 9091
	}`, name, host)
}
