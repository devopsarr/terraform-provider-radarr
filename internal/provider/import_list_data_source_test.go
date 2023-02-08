package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccImportListDataSourceConfig("radarr_import_list.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_import_list.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_import_list.test", "monitor", "movieOnly")),
			},
			// Not found testing
			{
				Config:      testAccImportListDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find import_list"),
			},
		},
	})
}

func testAccImportListDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list" "test" {
		enabled = false
		enable_auto = false
		search_on_add = false
		list_type = "program"
		root_folder_path = "/config"
		monitor = "movieOnly"
		minimum_availability = "tba"
		quality_profile_id = 1
		name = "importListDataTest"
		implementation = "RadarrImport"
		config_contract = "RadarrSettings"
		base_url = "http://127.0.0.1:7878"
		api_key = "testAPIKey"
	}
	
	data "radarr_import_list" "test" {
		name = %s
	}
	`, name)
}
