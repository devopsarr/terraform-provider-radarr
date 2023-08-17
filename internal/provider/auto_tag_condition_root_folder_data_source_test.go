package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAutoTagConditionRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccAutoTagConditionRootFolderDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_auto_tag_condition_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_auto_tag_condition_root_folder.test", "name", "Test"),
					resource.TestCheckResourceAttr("radarr_auto_tag.test", "specifications.0.value", "/test")),
			},
		},
	})
}

const testAccAutoTagConditionRootFolderDataSourceConfig = `
resource "radarr_tag" "test" {
	label = "atconditionfolder"
}

data  "radarr_auto_tag_condition_root_folder" "test" {
	name = "Test"
	negate = false
	required = false
	value = "/test"
}

resource "radarr_auto_tag" "test" {
	remove_tags_automatically = false
	name = "TestWithDSRootFolder"

	tags = [radarr_tag.test.id]
	
	specifications = [data.radarr_auto_tag_condition_root_folder.test]	
}`
