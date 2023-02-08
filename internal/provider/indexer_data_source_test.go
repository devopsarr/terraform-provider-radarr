package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIndexerDataSourceConfig("radarr_indexer.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_indexer.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_indexer.test", "protocol", "usenet")),
			},
			// Not found testing
			{
				Config:      testAccIndexerDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find indexer"),
			},
		},
	})
}

func testAccIndexerDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "radarr_indexer" "test" {
		enable_automatic_search = false
		name = "indexerdata"
		implementation = "Newznab"
		protocol = "usenet"
		config_contract = "NewznabSettings"
		base_url = "https://lolo.sickbeard.com"
		api_path = "/api"
		categories = [5030, 5040]
	}
	
	data "radarr_indexer" "test" {
		name = %s
	}
	`, name)
}
