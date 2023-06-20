package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientPneumaticResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientPneumaticResourceConfig("error", "/config/") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientPneumaticResourceConfig("resourcePneumaticTest", "/config/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_pneumatic.test", "nzb_folder", "/config/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_pneumatic.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientPneumaticResourceConfig("error", "/config/") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientPneumaticResourceConfig("resourcePneumaticTest", "/config/logs/"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_pneumatic.test", "nzb_folder", "/config/logs/"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_pneumatic.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientPneumaticResourceConfig(name, folder string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_pneumatic" "test" {
		enable = false
		priority = 1
		name = "%s"
		nzb_folder = "%s"
		strm_folder = "/config/"
	}`, name, folder)
}
