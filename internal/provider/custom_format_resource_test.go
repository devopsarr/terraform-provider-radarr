package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFormatResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccCustomFormatResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccCustomFormatResourceConfig("resourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_custom_format.test", "include_custom_format_when_renaming", "false"),
					resource.TestCheckResourceAttrSet("radarr_custom_format.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccCustomFormatResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccCustomFormatResourceConfig("resourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_custom_format.test", "include_custom_format_when_renaming", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_custom_format.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCustomFormatResourceConfig(name, enable string) string {
	return fmt.Sprintf(`
	resource "radarr_custom_format" "test" {
		include_custom_format_when_renaming = %s
		name = "%s"
		
		specifications = [
			{
				name = "Surround Sound"
				implementation = "ReleaseTitleSpecification"
				negate = false
				required = false
				value = "DTS.?(HD|ES|X(?!\\D))|TRUEHD|ATMOS|DD(\\+|P).?([5-9])|EAC3.?([5-9])"
			},
			{
				name = "Arabic"
				implementation = "LanguageSpecification"
				negate = false
				required = false
				value = "31"
			},
			{
				name = "Size"
				implementation = "SizeSpecification"
				negate = false
				required = false
				min = 0
				max = 100
			}
		]
	}`, enable, name)
}
