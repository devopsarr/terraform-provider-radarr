package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/terraform-provider-sonarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golift.io/starr/radarr"
)

const indexerConfigDataSourceName = "indexer_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &IndexerConfigDataSource{}

func NewIndexerConfigDataSource() datasource.DataSource {
	return &IndexerConfigDataSource{}
}

// IndexerConfigDataSource defines the indexer config implementation.
type IndexerConfigDataSource struct {
	client *radarr.Radarr
}

func (d *IndexerConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + indexerConfigDataSourceName
}

func (d *IndexerConfigDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Indexers -->[Indexer Config](../resources/indexer_config).",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"maximum_size": {
				MarkdownDescription: "Maximum size.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"minimum_age": {
				MarkdownDescription: "Minimum age.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"retention": {
				MarkdownDescription: "Retention.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"rss_sync_interval": {
				MarkdownDescription: "RSS sync interval.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"availability_delay": {
				MarkdownDescription: "Availability delay.",
				Computed:            true,
				Type:                types.Int64Type,
			},
			"whitelisted_hardcoded_subs": {
				MarkdownDescription: "Whitelisted hardconded subs.",
				Computed:            true,
				Type:                types.StringType,
			},
			"prefer_indexer_flags": {
				MarkdownDescription: "Prefer indexer flags.",
				Computed:            true,
				Type:                types.BoolType,
			},
			"allow_hardcoded_subs": {
				MarkdownDescription: "Allow hardcoded subs.",
				Computed:            true,
				Type:                types.BoolType,
			},
		},
	}, nil
}

func (d *IndexerConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.Radarr)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *radarr.Radarr, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *IndexerConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get indexer config current value
	response, err := d.client.GetIndexerConfigContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", indexerConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+indexerConfigDataSourceName)

	status := IndexerConfig{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}