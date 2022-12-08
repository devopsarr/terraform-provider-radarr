package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientHadoukenResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_hadouken.test", "host", "hadouken"),
					resource.TestCheckResourceAttr("radarr_download_client_hadouken.test", "url_base", "/hadouken/"),
					resource.TestCheckResourceAttrSet("radarr_download_client_hadouken.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientHadoukenResourceConfig("resourceHadoukenTest", "hadouken-host"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_hadouken.test", "host", "hadouken-host"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_download_client_hadouken.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientHadoukenResourceConfig(name, host string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_hadouken" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "%s"
		url_base = "/hadouken/"
		port = 9091
		category = "radarr-tv"
		username = "username"
		password = "password"
	}`, name, host)
}
