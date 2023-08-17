package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccTagResourceConfig("test", "autotag") + testAccAutoTagResourceConfig("test", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "autotag") + testAccAutoTagResourceConfig("Test", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_auto_tag.test", "remove_tags_automatically", "false"),
					resource.TestCheckResourceAttrSet("radarr_auto_tag.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccTagResourceConfig("test", "autotag") + testAccAutoTagResourceConfig("test", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
				Destroy:     true,
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "autotag") + testAccAutoTagResourceConfig("Test", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_auto_tag.test", "remove_tags_automatically", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_auto_tag.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccAutoTagResourceConfig(name, flag string) string {
	return fmt.Sprintf(`
		resource "radarr_auto_tag" "test" {
  			name = "%s"
			remove_tags_automatically = %s
			tags = [radarr_tag.test.id]

			specifications = [
				{
					name = "folder"
					implementation = "RootFolderSpecification"
					negate = true
					required = false
					value = "/config"
				},
				{
					name = "year"
					implementation = "YearSpecification"
					negate = false
					required = false
					min = 1900
					max = 1910
				},
				{
					name = "genre"
					implementation = "GenreSpecification"
					negate = false
					required = false
					value = "horror comedy"
				},
			]
		}
	`, name, flag)
}
