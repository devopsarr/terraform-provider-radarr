package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const importListConfigDataSourceName = "import_list_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ImportListConfigDataSource{}

func NewImportListConfigDataSource() datasource.DataSource {
	return &ImportListConfigDataSource{}
}

// ImportListConfigDataSource defines the import list config implementation.
type ImportListConfigDataSource struct {
	client *radarr.APIClient
}

func (d *ImportListConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + importListConfigDataSourceName
}

func (d *ImportListConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Import Lists -->[Import List Config](../resources/import_list_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Import List Config ID.",
				Computed:            true,
			},
			"sync_interval": schema.Int64Attribute{
				MarkdownDescription: "List Update Interval.",
				Computed:            true,
			},
			"sync_level": schema.StringAttribute{
				MarkdownDescription: "Clean library level.",
				Computed:            true,
			},
		},
	}
}

func (d *ImportListConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *ImportListConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.ImportListConfigApi.GetImportListConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, importListConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+importListConfigDataSourceName)

	config := ImportListConfig{}
	config.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, config)...)
}
