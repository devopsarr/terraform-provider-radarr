package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTraktListResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTraktListResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTraktListResourceConfig("resourceTraktListTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_trakt_list.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_trakt_list.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTraktListResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTraktListResourceConfig("resourceTraktListTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_trakt_list.test", "monitor", "movieOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_trakt_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTraktListResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_trakt_list" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		monitor = "%s"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		auth_user = "User1"
		access_token = "Token"
		username = "User2"
		listname = "test"
		limit = 100
	}`, monitor, name)
}
