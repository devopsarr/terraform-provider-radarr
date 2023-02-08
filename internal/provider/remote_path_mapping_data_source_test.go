package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccRemotePathMappingDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccRemotePathMappingDataSourceConfig("radarr_remote_path_mapping.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_remote_path_mapping.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_remote_path_mapping.test", "host", "transmission")),
			},
			// Not found testing
			{
				Config:      testAccRemotePathMappingDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find remote_path_mapping"),
			},
		},
	})
}

func testAccRemotePathMappingDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	resource "radarr_remote_path_mapping" "test" {
		host = "transmission"
		remote_path = "/datatest/"
		local_path = "/config/"
	}
	
	data "radarr_remote_path_mapping" "test" {
		id = %s
	}
	`, id)
}
