package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientNzbvortexResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientNzbvortexResourceConfig("error", "nzbvortex") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientNzbvortexResourceConfig("resourceNzbvortexTest", "nzbvortex"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_nzbvortex.test", "host", "nzbvortex"),
					resource.TestCheckResourceAttr("radarr_download_client_nzbvortex.test", "url_base", "/nzbvortex/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_nzbvortex.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientNzbvortexResourceConfig("error", "nzbvortex") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientNzbvortexResourceConfig("resourceNzbvortexTest", "nzbvortex-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_nzbvortex.test", "host", "nzbvortex-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_nzbvortex.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientNzbvortexResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_nzbvortex" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/nzbvortex/"
		port = 4321
		api_key = "testAPIkey"
	}`, name, host)
}
