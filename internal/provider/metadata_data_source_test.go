package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMetadataDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMetadataDataSourceConfig("\"Error\"") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccMetadataDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find metadata"),
			},
			// Read testing
			{
				Config: testAccMetadataResourceConfig("metadataData", "false") + testAccMetadataDataSourceConfig("radarr_metadata.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_metadata.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_metadata.test", "movie_metadata", "false")),
			},
		},
	})
}

func testAccMetadataDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "radarr_metadata" "test" {
		name = %s
	}
	`, name)
}
