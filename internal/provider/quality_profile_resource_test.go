package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccQualityProfileResource(t *testing.T) {
	// no parallel to avoid conflict with custom formats
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccQualityProfileResourceError + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("401 Unauthorized"),
			},
			// Create and Read testing
			{
				Config: testAccQualityProfileResourceConfig("example-4k"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_quality_profile.test", "name", "example-4k"),
					resource.TestCheckResourceAttrSet("radarr_quality_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccQualityProfileResourceError + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("401 Unauthorized"),
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

const testAccQualityProfileResourceError = `
resource "radarr_quality_profile" "test" {
	name            = "Error"
	upgrade_allowed = true
	cutoff          = 2000
	quality_groups = []
	language = {
		id = 0
	}
}
`

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

	data "radarr_language" "test" {
		name = "English"
	}

	data "radarr_quality" "bluray" {
		name = "Bluray-2160p"
	}

	data "radarr_quality" "webdl" {
		name = "WEBDL-2160p"
	}

	data "radarr_quality" "webrip" {
		name = "WEBRip-2160p"
	}

	resource "radarr_quality_profile" "test" {
		name            = "%s"
		upgrade_allowed = true
		cutoff          = 2000

		language = data.radarr_language.test

		quality_groups = [
			{
				id   = 2000
				name = "WEB 2160p"
				qualities = [
					data.radarr_quality.webdl,
					data.radarr_quality.webrip,
				]
			},
			{
				qualities = [data.radarr_quality.bluray]
			}
		]

		format_items = [
			{
				name   = radarr_custom_format.test.name
				format = radarr_custom_format.test.id
				score  = 10
			}
		]
	}`, name)
}
