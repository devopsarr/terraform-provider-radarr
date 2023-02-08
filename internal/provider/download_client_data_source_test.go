package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDownloadClientDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccDownloadClientDataSourceConfig("radarr_download_client.test.name"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_download_client.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_download_client.test", "protocol", "torrent")),
			},
			// Not found testing
			{
				Config:      testAccDownloadClientDataSourceConfig("\"Error\""),
				ExpectError: regexp.MustCompile("Unable to find download_client"),
			},
		},
	})
}

func testAccDownloadClientDataSourceConfig(name string) string {
	return fmt.Sprintf(`
	resource "radarr_download_client" "test" {
		enable = false
		priority = 1
		name = "dataTest"
		implementation = "Transmission"
		protocol = "torrent"
		config_contract = "TransmissionSettings"
		host = "transmission"
		url_base = "/transmission/"
		port = 9091
	}
	
	data "radarr_download_client" "test" {
		name = %s
	}
	`, name)
}
