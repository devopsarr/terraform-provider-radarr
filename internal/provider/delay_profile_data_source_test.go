package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDelayProfileDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDelayProfileDataSourceConfig(1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_delay_profile.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_delay_profile.test", "enable_usenet", "true")),
			},
			// Not found testing
			{
				Config:      testAccDelayProfileDataSourceConfig(999),
				ExpectError: regexp.MustCompile("Unable to find delay_profile"),
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
