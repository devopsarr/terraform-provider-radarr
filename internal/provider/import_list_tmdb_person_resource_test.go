package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListTMDBPersonResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListTMDBPersonResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListTMDBPersonResourceConfig("resourceTMDPersonTest", "none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_tmdb_person.test", "monitor", "none"),
					resource.TestCheckResourceAttrSet("radarr_import_list_tmdb_person.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListTMDBPersonResourceConfig("error", "none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListTMDBPersonResourceConfig("resourceTMDPersonTest", "movieOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_tmdb_person.test", "monitor", "movieOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_tmdb_person.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListTMDBPersonResourceConfig(name, monitor string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_tmdb_person" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		root_folder_path = "/config"
		monitor = "%s"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "%s"
		person_id = "11842"
		cast = true
		cast_director = true
		cast_producer = true
		cast_sound = true
	}`, monitor, name)
}
