package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMoviesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccMoviesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.radarr_movies.test", "movies.*", map[string]string{"title": "Gladiator"}),
				),
			},
		},
	})
}

const testAccMoviesDataSourceConfig = `
resource "radarr_movie" "test" {
	monitored = false
	title = "Gladiator"
	path = "/config/Gladiator_2000"
	quality_profile_id = 1
	tmdb_id = 98

	minimum_availability = "inCinemas"
}

data "radarr_movies" "test" {
	depends_on = [radarr_movie.test]
}
`
