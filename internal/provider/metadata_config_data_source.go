package provider

import (
	"context"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataConfigDataSourceName = "metadata_config"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataConfigDataSource{}

func NewMetadataConfigDataSource() datasource.DataSource {
	return &MetadataConfigDataSource{}
}

// MetadataConfigDataSource defines the metadata config implementation.
type MetadataConfigDataSource struct {
	client *radarr.APIClient
}

func (d *MetadataConfigDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataConfigDataSourceName
}

func (d *MetadataConfigDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Metadata -->[Metadata Config](../resources/metadata_config).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"certification_country": schema.StringAttribute{
				MarkdownDescription: "Certification Country.",
				Computed:            true,
			},
		},
	}
}

func (d *MetadataConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MetadataConfigDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get metadata config current value
	response, _, err := d.client.MetadataConfigApi.GetMetadataConfig(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.Read, metadataConfigDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataConfigDataSourceName)

	status := MetadataConfig{}
	status.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, status)...)
}
