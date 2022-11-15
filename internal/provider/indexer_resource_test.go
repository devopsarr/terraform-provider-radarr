package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "25"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer.test", "priority", "25"),
					resource.TestCheckResourceAttr("radarr_indexer.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("radarr_indexer.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "30"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer.test", "priority", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerResourceConfig(name, aSearch string) string {
	return fmt.Sprintf(`
	resource "radarr_indexer" "test" {
		priority = %s
		name = "%s"
		implementation = "Newznab"
		protocol = "usenet"
    	config_contract = "NewznabSettings"
		base_url = "https://lolo.sickbeard.com"
		api_path = "/api"
		categories = [8000, 5000]
	}`, aSearch, name)
}
