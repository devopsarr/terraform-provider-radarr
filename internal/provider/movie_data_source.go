package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const movieDataSourceName = "movie"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MovieDataSource{}

func NewMovieDataSource() datasource.DataSource {
	return &MovieDataSource{}
}

// MovieDataSource defines the movie implementation.
type MovieDataSource struct {
	client *radarr.APIClient
}

func (d *MovieDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + movieDataSourceName
}

func (d *MovieDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Movies -->Single [Movie](../resources/movie).",
		Attributes: map[string]schema.Attribute{
			"monitored": schema.BoolAttribute{
				MarkdownDescription: "Monitored flag.",
				Computed:            true,
			},
			"is_available": schema.BoolAttribute{
				MarkdownDescription: "Availability flag.",
				Computed:            true,
			},
			"quality_profile_id": schema.Int64Attribute{
				MarkdownDescription: "Quality profile ID.",
				Computed:            true,
			},
			"tmdb_id": schema.Int64Attribute{
				MarkdownDescription: "TMDB ID.",
				Required:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Movie ID.",
				Computed:            true,
			},
			"year": schema.Int64Attribute{
				MarkdownDescription: "Year.",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Movie title.",
				Computed:            true,
			},
			"path": schema.StringAttribute{
				MarkdownDescription: "Full movie path.",
				Computed:            true,
			},
			"minimum_availability": schema.StringAttribute{
				MarkdownDescription: "Minimum availability.\nAllowed values: 'tba', 'announced', 'inCinemas', 'released', 'deleted'.",
				Computed:            true,
			},
			"original_title": schema.StringAttribute{
				MarkdownDescription: "Movie original title.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Movie status.",
				Computed:            true,
			},
			"overview": schema.StringAttribute{
				MarkdownDescription: "Overview.",
				Computed:            true,
			},
			"website": schema.StringAttribute{
				MarkdownDescription: "Website.",
				Computed:            true,
			},
			"imdb_id": schema.StringAttribute{
				MarkdownDescription: "IMDB ID.",
				Computed:            true,
			},
			"youtube_trailer_id": schema.StringAttribute{
				MarkdownDescription: "Youtube trailer ID.",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				MarkdownDescription: "List of associated tags.",
				Computed:            true,
				ElementType:         types.Int64Type,
			},
			"genres": schema.SetAttribute{
				MarkdownDescription: "List genres.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"original_language": schema.SingleNestedAttribute{
				MarkdownDescription: "Origina language.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.Int64Attribute{
						MarkdownDescription: "ID.",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Name.",
						Computed:            true,
					},
				},
			},
		},
	}
}

func (d *MovieDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MovieDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var movie *Movie

	resp.Diagnostics.Append(req.Config.Get(ctx, &movie)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get movies current value
	response, _, err := d.client.MovieApi.ListMovie(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, movieDataSourceName, err))

		return
	}

	value, err := findMovie(movie.TMDBID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", movieDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+movieDataSourceName)
	movie.write(ctx, value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &movie)...)
}

func findMovie(ID int64, movies []*radarr.MovieResource) (*radarr.MovieResource, error) {
	for _, t := range movies {
		if t.GetTmdbId() == int32(ID) {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(movieDataSourceName, "TMDB ID", strconv.Itoa(int(ID)))
}
