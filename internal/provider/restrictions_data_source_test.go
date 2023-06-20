package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRestrictionsDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRestrictionsDataSourceConfig + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create a resource to check
			{
				Config: testAccRestrictionResourceConfig("testDataSource", "testDataSource2"),
			},
			// Read testing
			{
				Config: testAccRestrictionsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.radarr_restrictions.test", "restrictions.*", map[string]string{"ignored": "testDataSource"}),
				),
			},
		},
	})
}

const testAccRestrictionsDataSourceConfig = `
data "radarr_restrictions" "test" {
}
`
