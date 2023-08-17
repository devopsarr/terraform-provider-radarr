package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagConditionGenresDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionGenresDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_auto_tag_condition_genres.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_auto_tag_condition_genres.test", "name", "Test"),
					resource.TestCheckResourceAttr("radarr_auto_tag.test", "specifications.0.value", "horror comedy")),
			},
		},
	})
}

const testAccAutoTagConditionGenresDataSourceConfig = `
resource "radarr_tag" "test" {
	label = "atconditiongenre"
}

data  "radarr_auto_tag_condition_genres" "test" {
	name = "Test"
	negate = false
	required = false
	value = "horror comedy"
}

resource "radarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSGenres"

	tags = [radarr_tag.test.id]
	
	specifications = [data.radarr_auto_tag_condition_genres.test]	
}`
