package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListStevenluResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListStevenluResourceConfig("resourceStevenluTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_stevenlu.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_stevenlu.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListStevenluResourceConfig("resourceStevenluTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_stevenlu.test", "monitor", "movieOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_stevenlu.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListStevenluResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_stevenlu" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		monitor = "%s"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		link = "https://s3.amazonaws.com/popular-movies/movies.json"
	}`, monitor, name)
}
