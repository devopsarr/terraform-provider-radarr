package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientTransmissionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientTransmissionResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientTransmissionResourceConfig("resourceTransmissionTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_transmission.test", "enable", "false"),
					resource.TestCheckResourceAttr("radarr_download_client_transmission.test", "url_base", "/transmission/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_transmission.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientTransmissionResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientTransmissionResourceConfig("resourceTransmissionTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_transmission.test", "enable", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_transmission.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientTransmissionResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_transmission" "test" {
		enable = %s
		priority = 1
		name = "%s"
		host = "transmission"
		url_base = "/transmission/"
		port = 9091
	}`, enable, name)
}
