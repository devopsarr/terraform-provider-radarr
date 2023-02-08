package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLanguageDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLanguageDataSourceConfig("English"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_language.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_language.test", "name_lower", "english"),
				),
			},
			// Not found testing
			{
				Config:      testAccLanguageDataSourceConfig("Error"),
				ExpectError: regexp.MustCompile("Unable to find language"),
			},
		},
	})
}

func testAccLanguageDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	data "radarr_language" "test" {
		name = "%s"
	}
	`, name)
}
