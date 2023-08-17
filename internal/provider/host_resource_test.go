package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccHostResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccHostResourceConfig("Radarr") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccHostResourceConfig("Radarr"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_host.test", "port", "7878"),
					resource.TestCheckResourceAttrSet("radarr_host.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccHostResourceConfig("Radarr") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccHostResourceConfig("RadarrTest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_host.test", "port", "7878"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_host.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccHostResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "radarr_host" "test" {
		launch_browser = true
		port = 7878
		url_base = ""
		bind_address = "*"
		application_url =  ""
		instance_name = "%s"
		proxy = {
			enabled = false
		}
		ssl = {
			enabled = false
			certificate_validation = "enabled"
		}
		logging = {
			log_level = "info"
		}
		backup = {
			folder = "/backup"
			interval = 5
			retention = 10
		}
		authentication = {
			method = "none"
		}
		update = {
			mechanism = "docker"
			branch = "develop"
		}
	}`, name)
}
