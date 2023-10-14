package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDownloadClientFreeboxResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDownloadClientFreeboxResourceConfig("error", "Token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccDownloadClientFreeboxResourceConfig("resourceFreeboxTest", "Token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_freebox.test", "app_token", "Token123"),
					resource.TestCheckResourceAttrSet("radarr_download_client_freebox.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDownloadClientFreeboxResourceConfig("error", "Token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccDownloadClientFreeboxResourceConfig("resourceFreeboxTest", "Token321"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_download_client_freebox.test", "app_token", "Token321"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "radarr_download_client_freebox.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"app_token"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDownloadClientFreeboxResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client_freebox" "test" {
		enable = false
		priority = 1
		name = "%s"
		host = "mafreebox.freebox.fr"
		api_url = "/api/v1/"
		port = 443
		app_id = "test"
		app_token = "%s"
	}`, name, token)
}
