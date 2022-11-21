package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"golift.io/starr"
	"golift.io/starr/radarr"
)

func TestAccQualityProfilesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				PreConfig: qualityprofilesDSInit,
				Config:    testAccQualityProfilesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.radarr_quality_profiles.test", "quality_profiles.*", map[string]string{"name": "Any"}),
				),
			},
		},
	})
}

const testAccQualityProfilesDataSourceConfig = `
data "radarr_quality_profiles" "test" {
}
`

func qualityprofilesDSInit() {
	// keep only first two profiles to avoid longer tests
	client := *radarr.New(starr.New(os.Getenv("RADARR_API_KEY"), os.Getenv("RADARR_URL"), 0))
	for i := 3; i < 7; i++ {
		_ = client.DeleteQualityProfile(int64(i))
	}
}
