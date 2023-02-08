package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImportListConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListConfigResourceConfig("removeAndDelete") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccImportListConfigResourceConfig("removeAndDelete"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_config.test", "sync_level", "removeAndDelete"),
					resource.TestCheckResourceAttrSet("radarr_import_list_config.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListConfigResourceConfig("removeAndDelete") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListConfigResourceConfig("logOnly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_config.test", "sync_level", "logOnly"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListConfigResourceConfig(level string) string {
	return fmt.Sprintf(`
	resource "radarr_import_list_config" "test" {
		sync_interval = 24
		sync_level = "%s"
	}`, level)
}
