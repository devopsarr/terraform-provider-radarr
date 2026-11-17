package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/devopsarr/radarr-go/radarr"
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

// TestAccDownloadClientTransmissionResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccDownloadClientTransmissionResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccDownloadClientTransmissionResourceConfig("resourceTransmissionTest", "true"),
				Check: testAccCheckResourceDisappears("radarr_download_client_transmission.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.DownloadClientAPI.DeleteDownloadClient(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccDownloadClientTransmissionResourceConfig("resourceTransmissionTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_download_client_transmission.test", "id"),
				),
			},
		},
	})
}
