package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_quality_profile.test", "name", "example-4k"),
					resource.TestCheckResourceAttrSet("radarr_quality_profile.test", "id"),
				),
			},
			// Update and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-HD"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_quality_profile.test", "name", "example-HD"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_quality_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccQualityProfileResourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "radarr_custom_format" "test" {
		include_custom_format_when_renaming = false
		name = "QualityFormatTest"
		
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
	}

	data "radarr_custom_formats" "test" {
		depends_on = [radarr_custom_format.test]
	}

	resource "radarr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 1003

		language = {
			id   = 1
			name = "English"
		}

		quality_groups = [
			{
				id   = 1003
				name = "WEB 2160p"
				qualities = [
					{
						id         = 18
						name       = "WEBDL-2160p"
						source     = "webdl"
						resolution = 2160
					},
					{
						id         = 17
						name       = "WEBRip-2160p"
						source     = "webrip"
						resolution = 2160
					}
				]
			},
			{
				qualities = [
					{
						id = 19
						name = "Bluray-2160p"
						source = "bluray"
						resolution = 2160
					}
				]
			}
		]

		format_items = [
			for format in data.radarr_custom_formats.test.custom_formats :
			{
				name   = format.name
				format = format.id
				score  = 0
			}
		]
	}`, name)
}
