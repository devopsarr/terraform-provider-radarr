package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDelayProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccDelayProfileDataSourceConfig(999) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccDelayProfileDataSourceConfig(999),
				ExpectError: regexp.MustCompile("Unable to find delay_profile"),
			},
			// Read testing
			{
				Config: testAccDelayProfileDataSourceConfig(1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_delay_profile.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_delay_profile.test", "enable_usenet", "true")),
			},
		},
	})
}

func testAccDelayProfileDataSourceConfig(id int) string {
	return fmt.Sprintf(`
	data "radarr_delay_profile" "test" {
		id = %d
	}
	`, id)
}
