package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerOmgwtfnzbsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerOmgwtfnzbsResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerOmgwtfnzbsResourceConfig("omgwtfnzbsResourceTest", 30),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_omgwtfnzbs.test", "delay", "30"),
					resource.TestCheckResourceAttrSet("radarr_indexer_omgwtfnzbs.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerOmgwtfnzbsResourceConfig("error", 1) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerOmgwtfnzbsResourceConfig("omgwtfnzbsResourceTest", 60),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_omgwtfnzbs.test", "delay", "60"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_omgwtfnzbs.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerOmgwtfnzbsResourceConfig(name string, delay int) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_omgwtfnzbs" "test" {
		enable_automatic_search = false
		name = "%s"
		username = "Username"
		api_key = "API_Key"
		delay = %d
	}`, name, delay)
}
