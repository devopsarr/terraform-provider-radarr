package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const indexerConfigDataSourceName = "indexer_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerConfigDataSource{}

func NewIndexerConfigDataSource() datasource.DataSource {
	return &IndexerConfigDataSource{}
}

// IndexerConfigDataSource defines the indexer config implementation.
type IndexerConfigDataSource struct {
	client *radarr.APIClient
	auth   context.Context
}

func (d *IndexerConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerConfigDataSourceName
}

func (d *IndexerConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->\n[Indexer Config](../resources/indexer_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"maximum_size": schema.Int64Attribute{
				MarkdownDescription: "Maximum size.",
				Computed:            true,
			},
			"minimum_age": schema.Int64Attribute{
				MarkdownDescription: "Minimum age.",
				Computed:            true,
			},
			"retention": schema.Int64Attribute{
				MarkdownDescription: "Retention.",
				Computed:            true,
			},
			"rss_sync_interval": schema.Int64Attribute{
				MarkdownDescription: "RSS sync interval.",
				Computed:            true,
			},
			"availability_delay": schema.Int64Attribute{
				MarkdownDescription: "Availability delay.",
				Computed:            true,
			},
			"whitelisted_hardcoded_subs": schema.StringAttribute{
				MarkdownDescription: "Whitelisted hardconded subs.",
				Computed:            true,
			},
			"prefer_indexer_flags": schema.BoolAttribute{
				MarkdownDescription: "Prefer indexer flags.",
				Computed:            true,
			},
			"allow_hardcoded_subs": schema.BoolAttribute{
				MarkdownDescription: "Allow hardcoded subs.",
				Computed:            true,
			},
		},
	}
}

func (d *IndexerConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *IndexerConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, _, err := d.client.IndexerConfigAPI.GetIndexerConfig(d.auth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, indexerConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerConfigDataSourceName)

	status := IndexerConfig{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}
