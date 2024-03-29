package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTMDBCompanyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTMDBCompanyResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTMDBCompanyResourceConfig("resourceTMDCompanyTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_tmdb_company.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_tmdb_company.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTMDBCompanyResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTMDBCompanyResourceConfig("resourceTMDCompanyTest", "movieOnly"),
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
