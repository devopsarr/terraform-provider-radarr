package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	auth   context.Context
}

func (d *ImportListExclusionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListExclusionDataSourceName
}

func (d *ImportListExclusionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "<!-- subcategory:Import Lists -->\nSingle [ImportListExclusion](../resources/import_list_exclusion).",
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
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *ImportListExclusionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ImportListExclusion

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get importListExclusions current value
	response, _, err := d.client.ImportExclusionsAPI.ListExclusions(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListExclusionDataSourceName, err))

		return
	}

	data.find(data.TMDBID.ValueInt64(), response, &resp.Diagnostics)
	tflog.Trace(ctx, "read "+importListExclusionDataSourceName)
	// Map response body to resource schema attribute
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *ImportListExclusion) find(tmdbID int64, importListExclusions []radarr.ImportExclusionsResource, diags *diag.Diagnostics) {
	for _, t := range importListExclusions {
		if t.GetTmdbId() == int32(tmdbID) {
			i.write(&t)

			return
		}
	}

	diags.AddError(helpers.DataSourceError, helpers.ParseNotFoundError(importListExclusionDataSourceName, "tmdb_id", strconv.Itoa(int(tmdbID))))
}
