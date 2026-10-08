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

func TestAccRootFolderResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized create
			{
				Config:      testAccRootFolderResourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccRootFolderResourceConfig("/config/asp"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_root_folder.test", "path", "/config/asp"),
					resource.TestCheckResourceAttrSet("radarr_root_folder.test", "id"),
				),
			},
			// Unauthorized read
			{
				Config:      testAccRootFolderResourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccRootFolderResourceConfig("/config/logs"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_root_folder.test", "path", "/config/logs"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_root_folder.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccRootFolderResourceConfig(path string) string {
	return fmt.Sprintf(`
		resource "radarr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}

// TestAccRootFolderResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccRootFolderResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccRootFolderResourceConfig("/config/logs"),
				Check: testAccCheckResourceDisappears("radarr_root_folder.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.RootFolderAPI.DeleteRootFolder(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccRootFolderResourceConfig("/config/logs"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_root_folder.test", "id"),
				),
			},
		},
	})
}
