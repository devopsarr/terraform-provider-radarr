package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccImportListExclusionDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccImportListExclusionDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccImportListExclusionDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find import_list_exclusion"),
			},
			// Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 987) + testAccImportListExclusionDataSourceConfig("radarr_import_list_exclusion.test.tmdb_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_import_list_exclusion.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_import_list_exclusion.test", "title", "Test"),
				),
			},
		},
	})
}

func testAccImportListExclusionDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "radarr_import_list_exclusion" "test" {
		tmdb_id = %s
	}
	`, id)
}
