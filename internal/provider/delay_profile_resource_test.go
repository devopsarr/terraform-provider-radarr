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

func TestAccDelayProfileResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccDelayProfileResourceConfig("usenet", "0") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("usenet", "radarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_delay_profile.test", "preferred_protocol", "usenet"),
					resource.TestCheckResourceAttrSet("radarr_delay_profile.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccDelayProfileResourceConfig("usenet", "0") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("torrent", "radarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_delay_profile.test", "preferred_protocol", "torrent"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_delay_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccDelayProfileResourceConfig(protocol, tag string) string {
	return fmt.Sprintf(`
	resource "radarr_delay_profile" "test" {
		enable_usenet = true
		enable_torrent = true
		bypass_if_highest_quality = true
		order = 100
		usenet_delay = 0
		torrent_delay = 0
		preferred_protocol= "%s"
		tags = [%s]
	}`, protocol, tag)
}

// TestAccDelayProfileResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccDelayProfileResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("torrent", "radarr_tag.test.id"),
				Check: testAccCheckResourceDisappears("radarr_delay_profile.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.DelayProfileAPI.DeleteDelayProfile(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccTagResourceConfig("test", "delay-profile-resource") + testAccDelayProfileResourceConfig("torrent", "radarr_tag.test.id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_delay_profile.test", "id"),
				),
			},
		},
	})
}
