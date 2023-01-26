package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListTMDBCompanyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTMDBCompanyResourceConfig("resourceTMDBCompanyTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_tmdb_company.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_tmdb_company.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccImportListTMDBCompanyResourceConfig("resourceTMDBCompanyTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_tmdb_company.test", "monitor", "movieOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_tmdb_company.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTMDBCompanyResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_tmdb_company" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		monitor = "%s"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		company_id = "11842"
	}`, monitor, name)
}
