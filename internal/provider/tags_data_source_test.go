package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTagsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create a tag to have a value to check
			{
				Config: testAccTagResourceConfig("test-1", "1080p") + testAccTagResourceConfig("test-2", "4k"),
			},
			// Read testing
			{
				Config: testAccTagsDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.radarr_tags.test", "tags.*", map[string]string{"label": "1080p"}),
				),
			},
		},
	})
}

const testAccTagsDataSourceConfig = `
data "radarr_tags" "test" {
}
`
