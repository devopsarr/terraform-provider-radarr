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

func TestAccImportListCouchPotatoResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListCouchPotatoResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListCouchPotatoResourceConfig("resourceCouchPotatoTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_couch_potato.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_couch_potato.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListCouchPotatoResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListCouchPotatoResourceConfig("resourceCouchPotatoTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_couch_potato.test", "monitor", "movieOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "radarr_import_list_couch_potato.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListCouchPotatoResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_couch_potato" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		monitor = "%s"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		link = "http://localhost"
		api_key = "APIKey"
		port = 5050
		only_active = true
	}`, monitor, name)
}

// TestAccImportListCouchPotatoResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccImportListCouchPotatoResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccImportListCouchPotatoResourceConfig("resourceCouchPotatoTest", "movieOnly"),
				Check: testAccCheckResourceDisappears("radarr_import_list_couch_potato.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.ImportListAPI.DeleteImportList(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccImportListCouchPotatoResourceConfig("resourceCouchPotatoTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_import_list_couch_potato.test", "id"),
				),
			},
		},
	})
}
