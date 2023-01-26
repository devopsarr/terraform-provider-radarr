package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListRadarrResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListRadarrResourceConfig("resourceRadarrTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_radarr.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_radarr.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListRadarrResourceConfig("resourceRadarrTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_radarr.test", "monitor", "movieOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_radarr.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListRadarrResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_radarr" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		monitor = "%s"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		base_url = "http://127.0.0.1:7878"
		api_key = "testAPIKey"
		tag_ids = [1,2]
		profile_ids = [1]
	}`, monitor, name)
}
