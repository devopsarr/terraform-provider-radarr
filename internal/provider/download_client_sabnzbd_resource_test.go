package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientSabnzbdResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientSabnzbdResourceConfig("resourceSabnzbdTest", "sabnzbd"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_sabnzbd.test", "host", "sabnzbd"),
					resource.TestCheckResourceAttr("radarr_download_client_sabnzbd.test", "url_base", "/sabnzbd/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_sabnzbd.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientSabnzbdResourceConfig("resourceSabnzbdTest", "sabnzbd-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_sabnzbd.test", "host", "sabnzbd-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_sabnzbd.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientSabnzbdResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_sabnzbd" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/sabnzbd/"
		port = 8080
		api_key = "testAPIkey"
	}`, name, host)
}
