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

func TestAccImportListExclusionResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccImportListExclusionResourceConfig("error", 123) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 123),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_exclusion.test", "tmdb_id", "123"),
					resource.TestCheckResourceAttrSet("radarr_import_list_exclusion.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccImportListExclusionResourceConfig("error", 123) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccImportListExclusionResourceConfig("test", 1234),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_import_list_exclusion.test", "tmdb_id", "1234"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_import_list_exclusion.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccImportListExclusionResourceConfig(name string, tvID int) string {
	return fmt.Sprintf(`
		resource "radarr_import_list_exclusion" "%s" {
  			title = "Test"
			tmdb_id = %d
			year = 1900
		}
	`, name, tvID)
}

// TestAccImportListExclusionResourceDisappears is not parallel: it deletes an object outside Terraform, and a parallel test could otherwise take over its freed ID.
func TestAccImportListExclusionResourceDisappears(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create, then delete outside Terraform: Read removes it from state instead of failing the plan
			{
				Config: testAccImportListExclusionResourceConfig("test", 1234),
				Check: testAccCheckResourceDisappears("radarr_import_list_exclusion.test", func(client *radarr.APIClient, id int32) (*http.Response, error) {
					return client.ImportListExclusionAPI.DeleteExclusions(context.TODO(), id).Execute()
				}),
				ExpectNonEmptyPlan: true,
			},
			// Create again after the deletion outside Terraform
			{
				Config: testAccImportListExclusionResourceConfig("test", 1234),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("radarr_import_list_exclusion.test", "id"),
				),
			},
		},
	})
}
