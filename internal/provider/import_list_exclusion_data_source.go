package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListExclusionDataSourceName = "import_list_exclusion"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListExclusionDataSource{}

func NewImportListExclusionDataSource() datasource.DataSource {
	return &ImportListExclusionDataSource{}
}

// ImportListExclusionDataSource defines the importListExclusion implementation.
type ImportListExclusionDataSource struct {
	client *radarr.APIClient
}

func (d *ImportListExclusionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListExclusionDataSourceName
}

func (d *ImportListExclusionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->Single [ImportListExclusion](../resources/import_list_exclusion).",
		Attributes: map[string]schema.Attribute{
			"tmdb_id": schema.Int64Attribute{
				MarkdownDescription: "Movie TMDB ID.",
				Required:            true,
			},
			"year": schema.Int64Attribute{
				MarkdownDescription: "Year.",
				Computed:            true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "Movie to be excluded.",
				Computed:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List Exclusion ID.",
				Computed:            true,
			},
		},
	}
}

func (d *ImportListExclusionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ImportListExclusionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var importListExclusion *ImportListExclusion

	resp.Diagnostics.Append(req.Config.Get(ctx, &importListExclusion)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListExclusions current value
	response, _, err := d.client.ImportExclusionsApi.ListExclusions(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionDataSourceName, err))

		return
	}

	value, err := findImportListExclusion(importListExclusion.TMDBID.ValueInt64(), response)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, fmt.Sprintf("Unable to find %s, got error: %s", importListExclusionDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListExclusionDataSourceName)
	importListExclusion.write(value)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &importListExclusion)...)
}

func findImportListExclusion(tvID int64, importListExclusions []*radarr.ImportExclusionsResource) (*radarr.ImportExclusionsResource, error) {
	for _, t := range importListExclusions {
		if t.GetTmdbId() == int32(tvID) {
			return t, nil
		}
	}

	return nil, helpers.ErrDataNotFoundError(importListExclusionDataSourceName, "tmdb_id", strconv.Itoa(int(tvID)))
}
