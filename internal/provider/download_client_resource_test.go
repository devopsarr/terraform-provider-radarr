package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client.test", "enable", "false"),
					resource.TestCheckResourceAttr("radarr_download_client.test", "url_base", "/transmission/"),
					resource.TestCheckResourceAttrSet("radarr_download_client.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client.test", "enable", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client" "test" {
		enable = %s
		priority = 1
		name = "%s"
		implementation = "Transmission"
		protocol = "torrent"
    	config_contract = "TransmissionSettings"
		host = "transmission"
		url_base = "/transmission/"
		port = 9091
	}`, enable, name)
}
