package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIndexerIptorrentsResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerIptorrentsResourceConfig("error", "https://iptorrents.org") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.org"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_iptorrents.test", "base_url", "https://iptorrents.org"),
					resource.TestCheckResourceAttrSet("radarr_indexer_iptorrents.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerIptorrentsResourceConfig("error", "https://iptorrents.org") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccIndexerIptorrentsResourceConfig("iptorrentsResourceTest", "https://iptorrents.net"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer_iptorrents.test", "base_url", "https://iptorrents.net"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_indexer_iptorrents.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccIndexerIptorrentsResourceConfig(name, url string) string {
	return fmt.Sprintf(`
	resource "radarr_indexer_iptorrents" "test" {
		enable_rss = false
		name = "%s"
		base_url = "%s"
		minimum_seeders = 1
	}`, name, url)
}
