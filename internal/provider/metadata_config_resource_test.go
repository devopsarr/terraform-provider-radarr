package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMetadataConfigResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccMetadataConfigResourceConfig("us"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_config.test", "certification_country", "us"),
					resource.TestCheckResourceAttrSet("radarr_metadata_config.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccMetadataConfigResourceConfig("it"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_config.test", "certification_country", "it"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata_config.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataConfigResourceConfig(country string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata_config" "test" {
		certification_country = "%s"
	}`, country)
}
