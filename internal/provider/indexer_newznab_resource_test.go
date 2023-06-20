package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerNewznabResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerNewznabResourceConfig("error", "25") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "25"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_newznab.test", "priority", "25"),
					resource.TestCheckResourceAttr("radarr_indexer_newznab.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("radarr_indexer_newznab.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerNewznabResourceConfig("error", "25") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerNewznabResourceConfig("newzabResourceTest", "30"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_newznab.test", "priority", "30"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_newznab.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerNewznabResourceConfig(name, aSearch string) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_newznab" "test" {
		priority = %s
		name = "%s"
		base_url = "https://lolo.sickbeard.com"
		api_path = "/api"
		categories = [5030, 5040]
	}`, aSearch, name)
}
