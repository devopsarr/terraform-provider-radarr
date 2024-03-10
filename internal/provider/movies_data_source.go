package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const moviesDataSourceName = "movies"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MoviesDataSource{}

func NewMoviesDataSource() datasource.DataSource {
	return &MoviesDataSource{}
}

// MoviesDataSource defines the movies implementation.
type MoviesDataSource struct {
	client *radarr.APIClient
	auth   context.Context
}

// Movies describes the movies data model.
type Movies struct {
	Movies types.Set    `tfsdk:"movies"`
	ID     types.String `tfsdk:"id"`
}

func (d *MoviesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + moviesDataSourceName
}

func (d *MoviesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "<!-- subcategory:Movies -->\nList all available [Movies](../resources/movie).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"movies": schema.SetNestedAttribute{
				MarkdownDescription: "Movie list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:            true,
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
				},
			},
		},
	}
}

func (d *MoviesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *MoviesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get movies current value
	response, _, err := d.client.MovieAPI.ListMovie(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, moviesDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+moviesDataSourceName)
	// Map response body to resource schema attribute
	movies := make([]Movie, len(response))
	for i, m := range response {
		movies[i].write(ctx, &m, &resp.Diagnostics)
	}

	movieList, diags := types.SetValueFrom(ctx, Movie{}.getType(), movies)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, Movies{Movies: movieList, ID: types.StringValue(strconv.Itoa(len(response)))})...)
}
