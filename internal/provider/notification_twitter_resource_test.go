package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNotificationTwitterResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationTwitterResourceConfig("Error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationTwitterResourceConfig("resourceTwitterTest", "me"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_notification_twitter.test", "mention", "me"),
					resource.TestCheckResourceAttrSet("radarr_notification_twitter.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationTwitterResourceConfig("Error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationTwitterResourceConfig("resourceTwitterTest", "myself"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_notification_twitter.test", "mention", "myself"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_notification_twitter.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationTwitterResourceConfig(name, mention string) string {
	return fmt.Sprintf(`
	resource "radarr_notification_twitter" "test" {
		on_grab                            = false
		on_download                        = false
		on_upgrade                         = false
		on_movie_added                     = false
		on_movie_delete                    = false
		on_movie_file_delete               = false
		on_movie_file_delete_for_upgrade   = false
		on_health_issue                    = false
		on_application_update              = false
	  
		include_health_warnings = false
		name                    = "%s"
	  
		access_token = "Token"
		access_token_secret = "TokenSecret"
		consumer_key = "Key"
		consumer_secret = "Secret"
		mention = "%s"
	}`, name, mention)
}
