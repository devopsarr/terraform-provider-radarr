package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListResourceConfig("error", "true") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListResourceConfig("importListResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list.test", "enable_auto", "true"),
					resource.TestCheckResourceAttrSet("radarr_import_list.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListResourceConfig("error", "true") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListResourceConfig("importListResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list.test", "enable_auto", "false"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list" "test" {
		enabled = false
		enable_auto = "%s"
		search_on_add = false
		list_type = "program"
		root_folder_path = "/config"
		monitor = "movieOnly"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		implementation = "RadarrImport"
    	config_contract = "RadarrSettings"
		base_url = "http://127.0.0.1:7878"
		api_key = "testAPIKey"
		tag_ids = [1,2]
		profile_ids = [1]
		tags = []
	}`, monitor, name)
}
