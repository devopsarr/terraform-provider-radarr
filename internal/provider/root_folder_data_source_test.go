package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRootFolderDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccRootFolderDataSourceConfig("/error") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccRootFolderDataSourceConfig("/error"),
				ExpectError: regexp.MustCompile("Unable to find root_folder"),
			},
			// Read testing
			{
				PreConfig: rootFolderDSInit,
				Config:    testAccRootFolderDataSourceConfig("/config"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_root_folder.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_root_folder.test", "path", "/config")),
			},
		},
	})
}

func testAccRootFolderDataSourceConfig(path string) string {
	return fmt.Sprintf(`
	data "radarr_root_folder" "test" {
  			path = "%s"
		}
	`, path)
}

func rootFolderDSInit() {
	// ensure a /config root path is configured
	client := testAccAPIClient()
	folder := radarr.NewRootFolderResource()
	folder.SetPath("/config")
	_, _, _ = client.RootFolderAPI.CreateRootFolder(context.TODO()).RootFolderResource(*folder).Execute()
}
