package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationNtfyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationNtfyResourceConfig("error", "key1", "testtoken123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key1", "testtoken123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_notification_ntfy.test", "password", "key1"),
					resource.TestCheckResourceAttr("radarr_notification_ntfy.test", "access_token", "testtoken123"),
					resource.TestCheckResourceAttrSet("radarr_notification_ntfy.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationNtfyResourceConfig("error", "key1", "testtoken123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationNtfyResourceConfig("resourceNtfyTest", "key2", "testtoken234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_notification_ntfy.test", "password", "key2"),
					resource.TestCheckResourceAttr("radarr_notification_ntfy.test", "access_token", "testtoken234"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "radarr_notification_ntfy.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "access_token"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationNtfyResourceConfig(name, key, accessToken string) string {
	return fmt.Sprintf(`
	resource "radarr_notification_ntfy" "test" {
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

		priority = 1
		server_url = "https://ntfy.sh"
		username = "User"
		password = "%s"
		access_token = "%s"
		topics = ["Topic1234","Topic4321"]
		field_tags = ["warning","skull"]
	}`, name, key, accessToken)
}
