package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIndexerResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccIndexerResourceConfig("error", "25") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccIndexerResourceConfig("resourceTest", "25"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_indexer.test", "priority", "25"),
					resource.TestCheckResourceAttr("radarr_indexer.test", "base_url", "https://lolo.sickbeard.com"),
					resource.TestCheckResourceAttrSet("radarr_indexer.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccIndexerResourceConfig("error", "25") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
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
			{
				ResourceName:            "radarr_indexer.test_sensitive",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"passkey"},
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
	}

	resource "radarr_indexer" "test_sensitive" {
		enable_automatic_search = false
		name = "%sWithSensitive"
		base_url = "https://filelist.io"
		username = "test"
		passkey = "Pass"
		categories = [21,23,27]
		minimum_seeders = 1
		implementation = "FileList"
		protocol = "torrent"
    	config_contract = "FileListSettings"
	}
	`, aSearch, name, name)
}
