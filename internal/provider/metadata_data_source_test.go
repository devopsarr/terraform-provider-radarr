package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMetadataDataSourceConfig("radarr_metadata.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_metadata.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_metadata.test", "movie_metadata", "false")),
			},
			// Not found testing
			{
				Config:      testAccMetadataDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find metadata"),
			},
		},
	})
}

func testAccMetadataDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata" "test" {
		enable = true
		name = "metadataData"
		implementation = "MediaBrowserMetadata"
		config_contract = "MediaBrowserMetadataSettings"
		movie_metadata = false
	}
	
	data "radarr_metadata" "test" {
		name = %s
	}
	`, name)
}
