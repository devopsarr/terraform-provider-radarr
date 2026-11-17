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

func TestAccMetadataEmbyResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMetadataEmbyResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMetadataEmbyResourceConfig("embyResourceTest", "false"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_emby.test", "movie_metadata", "false"),
					resource.TestCheckResourceAttrSet("radarr_metadata_emby.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMetadataEmbyResourceConfig("error", "false") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMetadataEmbyResourceConfig("embyResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_metadata_emby.test", "movie_metadata", "true"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_metadata_emby.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMetadataEmbyResourceConfig(name, metadata string) string {
	return fmt.Sprintf(`
	resource "radarr_metadata_emby" "test" {
		enable = false
		name = "%s"
		movie_metadata = %s
	}`, name, metadata)
}

// TestAccMetadataEmbyResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccMetadataEmbyResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccMetadataEmbyResourceConfig("embyResourceTest", "true"),
				Check: testAccCheckResourceDisappears("radarr_metadata_emby.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.MetadataAPI.DeleteMetadata(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccMetadataEmbyResourceConfig("embyResourceTest", "true"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_metadata_emby.test", "id"),
				),
			},
		},
	})
}
