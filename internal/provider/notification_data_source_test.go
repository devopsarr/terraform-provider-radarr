package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccNotificationDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_notification.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_notification.test", "path", "/scripts/test.sh")),
			},
		},
	})
}

const testAccNotificationDataSourceConfig = `
resource "radarr_notification" "test" {
	on_grab                            = false
	on_download                        = true
	on_upgrade                         = false
	on_rename                          = false
	on_movie_added                     = false
	on_movie_delete                    = false
	on_movie_file_delete               = false
	on_movie_file_delete_for_upgrade   = true
	on_health_issue                    = false
	on_application_update              = false
  
	include_health_warnings = false
	name                    = "notificationData"
  
	implementation  = "CustomScript"
	config_contract = "CustomScriptSettings"
  
	path = "/scripts/test.sh"
}

data "radarr_notification" "test" {
	name = radarr_notification.test.name
}
`
