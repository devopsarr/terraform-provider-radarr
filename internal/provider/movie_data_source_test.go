package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMovieDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccMovieDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_movie.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_movie.test", "title", "Pulp Fiction"),
				),
			},
		},
	})
}

const testAccMovieDataSourceConfig = `
resource "radarr_movie" "test" {
	monitored = false
	title = "Pulp Fiction"
	path = "/config/Pulp_Fiction_2994"
	quality_profile_id = 1
	tmdb_id = 680
}

data "radarr_movie" "test" {
	tmdb_id = radarr_movie.test.tmdb_id
}
`
