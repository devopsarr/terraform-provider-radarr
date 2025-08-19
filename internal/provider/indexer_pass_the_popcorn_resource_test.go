package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerPassThePopcornResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerPassThePopcornResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerPassThePopcornResourceConfig("passThePopcornResourceTest", 1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_pass_the_popcorn.test", "minimum_seeders", "1"),
					resource.TestCheckResourceAttrSet("radarr_indexer_pass_the_popcorn.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerPassThePopcornResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerPassThePopcornResourceConfig("passThePopcornResourceTest", 2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_pass_the_popcorn.test", "minimum_seeders", "2"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_pass_the_popcorn.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerPassThePopcornResourceConfig(name string, seeders int) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_pass_the_popcorn" "test" {
		enable_automatic_search = false
		priority = 1
		name = "%s"
		base_url = "https://passthepopcorn.me"
		api_user = "test"
		minimum_seeders = %d
	}`, name, seeders)
}
