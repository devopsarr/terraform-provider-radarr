package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMoviesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMovieResourceConfig("Error", "error", 0) + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Read testing
			{
				Config: testAccMovieResourceConfig("Gladiator", "Gladiator_2000", 98) + testAccMoviesDataSourceConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.radarr_movies.test", "movies.*", map[string]string{"title": "Gladiator"}),
				),
			},
		},
	})
}

const testAccMoviesDataSourceConfig = `
data "radarr_movies" "test" {
	depends_on = [radarr_movie.test]
}
`
