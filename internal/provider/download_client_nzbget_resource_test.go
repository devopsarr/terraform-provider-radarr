package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientNzbgetResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientNzbgetResourceConfig("error", "nzbget") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientNzbgetResourceConfig("resourceNzbgetTest", "nzbget"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_nzbget.test", "host", "nzbget"),
					resource.TestCheckResourceAttr("radarr_download_client_nzbget.test", "url_base", "/nzbget/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_nzbget.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientNzbgetResourceConfig("error", "nzbget") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientNzbgetResourceConfig("resourceNzbgetTest", "nzbget-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_nzbget.test", "host", "nzbget-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_nzbget.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientNzbgetResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_nzbget" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/nzbget/"
		port = 9091
	}`, name, host)
}
