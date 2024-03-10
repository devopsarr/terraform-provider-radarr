package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mitchellh/hashstructure/v2"
)

const (
	customFormatConditionIndexerFlagDataSourceName = "custom_format_condition_indexer_flag"
	customFormatConditionIndexerFlagImplementation = "IndexerFlagSpecification"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &CustomFormatConditionIndexerFlagDataSource{}

func NewCustomFormatConditionIndexerFlagDataSource() datasource.DataSource {
	return &CustomFormatConditionIndexerFlagDataSource{}
}

// CustomFormatConditionIndexerFlagDataSource defines the custom format condition indexer flag implementation.
type CustomFormatConditionIndexerFlagDataSource struct {
	client *radarr.APIClient
	auth   context.Context
}

func (d *CustomFormatConditionIndexerFlagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + customFormatConditionIndexerFlagDataSourceName
}

func (d *CustomFormatConditionIndexerFlagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Profiles -->\n Custom Format Condition Indexer Flag data source.\nFor more information refer to [Custom Format Conditions](https://wiki.servarr.com/radarr/settings#conditions) and [Indexer Flag](https://wiki.servarr.com/radarr/settings#indexer-flags).",
		Attributes: map[string]schema.Attribute{
			"negate": schema.BoolAttribute{
				MarkdownDescription: "Negate flag.",
				Required:            true,
			},
			"required": schema.BoolAttribute{
				MarkdownDescription: "Computed flag.",
				Required:            true,
			},
			"implementation": schema.StringAttribute{
				MarkdownDescription: "Implementation.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Specification name.",
				Required:            true,
			},
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.Int64Attribute{
				MarkdownDescription: "Custom format condition indexer flag ID.",
				Computed:            true,
			},
			// Field values
			"value": schema.StringAttribute{
				MarkdownDescription: "Indexer flag ID. `1` G Freeleech, `2` G Halfleech, `4` G DoubleUpload, `8` PTP Golden, `16` PTP Approved, `32` HDB Internal, `64` AHD Internal, `128` G Scene, `256` G Freeleech75, `512` G Freeleech25, `1024` AHD UserRelease.",
				Required:            true,
			},
		},
	}
}

func (d *CustomFormatConditionIndexerFlagDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if auth, client := dataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
		d.auth = auth
	}
}

func (d *CustomFormatConditionIndexerFlagDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CustomFormatConditionValue

	hash, err := hashstructure.Hash(&data, hashstructure.FormatV2, nil)
	if err != nil {
		resp.Diagnostics.AddError(helpers.DataSourceError, helpers.ParseClientError(helpers.Create, customFormatConditionIndexerFlagDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+customFormatConditionIndexerFlagDataSourceName)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("implementation"), customFormatConditionIndexerFlagImplementation)...)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), int64(hash))...)
}
