package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccMovieDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Unauthorized
			{
				Config:      testAccMovieDataSourceConfig("999") + testUnauthorizedProvider,
				ExpectError: regexp.MustCompile("Client Error"),
			},
			// Not found testing
			{
				Config:      testAccMovieDataSourceConfig("999"),
				ExpectError: regexp.MustCompile("Unable to find movie"),
			},
			// Read testing
			{
				Config: testAccMovieResourceConfig("Pulp Fiction", "Pulp_Fiction_1994", 680) + testAccMovieDataSourceConfig("radarr_movie.test.tmdb_id"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.radarr_movie.test", "id"),
					resource.TestCheckResourceAttr("data.radarr_movie.test", "title", "Pulp Fiction"),
				),
			},
		},
	})
}

func testAccMovieDataSourceConfig(id string) string {
	return fmt.Sprintf(`
	data "radarr_movie" "test" {
		tmdb_id = %s
	}
	`, id)
}
