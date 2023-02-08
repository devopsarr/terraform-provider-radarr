package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMediaManagementResource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized Create
			{
				Config:      testAccMediaManagementResourceConfig("none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Create and Read testing
			{
				Config: testAccMediaManagementResourceConfig("none"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_media_management.test", "file_date", "none"),
					resource.TestCheckResourceAttrSet("radarr_media_management.test", "id"),
				),
			},
			// Unauthorized Read
			{
				Config:      testAccMediaManagementResourceConfig("none") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Update and Read testing
			{
				Config: testAccMediaManagementResourceConfig("cinemas"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("radarr_media_management.test", "file_date", "cinemas"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "radarr_media_management.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccMediaManagementResourceConfig(date string) string {
	return fmt.Sprintf(`
	resource "radarr_media_management" "test" {
		auto_unmonitor_previously_downloaded_movies = false
		recycle_bin = ""
		recycle_bin_cleanup_days = 7
		download_propers_and_repacks = "doNotPrefer"
		create_empty_movie_folders = false
		delete_empty_folders = false
		file_date = "%s"
		rescan_after_refresh = "afterManual"
		auto_rename_folders = false
		paths_default_static = false
		set_permissions_linux = false
		chmod_folder = 755
		chown_group = ""
		skip_free_space_check_when_importing = false
		minimum_free_space_when_importing = 100
		copy_using_hardlinks = true
		import_extra_files = true
		extra_file_extensions = "srt"
		enable_media_info = true
	}`, date)
}
