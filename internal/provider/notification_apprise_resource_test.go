package provider

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNotificationAppriseResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccNotificationAppriseResourceConfig("resourceAppriseTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccNotificationAppriseResourceConfig("resourceAppriseTest", "token123"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_notification_apprise.test", "auth_password", "token123"),
					resource.TestCheckResourceAttrSet("radarr_notification_apprise.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccNotificationAppriseResourceConfig("resourceAppriseTest", "token123") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccNotificationAppriseResourceConfig("resourceAppriseTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_notification_apprise.test", "auth_password", "token234"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "radarr_notification_apprise.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccNotificationAppriseResourceConfig(name, token string) string {
	return fmt.Sprintf(`
	resource "radarr_notification_apprise" "test" {
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

		notification_type = 1
		server_url = "https://apprise.go"
		auth_username = "User"
		auth_password = "%s"
		field_tags = ["warning","skull"]
	}`, name, token)
}

// TestAccNotificationAppriseResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccNotificationAppriseResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccNotificationAppriseResourceConfig("resourceAppriseTest", "token234"),
				Check: testAccCheckResourceDisappears("radarr_notification_apprise.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.NotificationAPI.DeleteNotification(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccNotificationAppriseResourceConfig("resourceAppriseTest", "token234"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_notification_apprise.test", "id"),
				),
			},
		},
	})
}
