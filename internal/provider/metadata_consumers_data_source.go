package provider

import (
	"context"
	"strconv"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const metadataConsumersDataSourceName = "metadata_consumers"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MetadataConsumersDataSource{}

func NewMetadataConsumersDataSource() datasource.DataSource {
	return &MetadataConsumersDataSource{}
}

// MetadataConsumersDataSource defines the metadataConsumers implementation.
type MetadataConsumersDataSource struct {
	client *radarr.APIClient
}

// MetadataConsumers describes the metadataConsumers data model.
type MetadataConsumers struct {
	MetadataConsumers types.Set    `tfsdk:"metadata_consumers"`
	ID                types.String `tfsdk:"id"`
}

func (d *MetadataConsumersDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + metadataConsumersDataSourceName
}

func (d *MetadataConsumersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Metadata -->List all available [Metadata Consumers](../resources/metadata).",
		Attributes: map[string]schema.Attribute{
			// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
			"id": schema.StringAttribute{
				Computed: true,
			},
			"metadata_consumers": schema.SetNestedAttribute{
				MarkdownDescription: "MetadataConsumer list.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							MarkdownDescription: "Enable flag.",
							Computed:            true,
						},
						"config_contract": schema.StringAttribute{
							MarkdownDescription: "Metadata configuration template.",
							Computed:            true,
						},
						"implementation": schema.StringAttribute{
							MarkdownDescription: "Metadata implementation name.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Metadata name.",
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							MarkdownDescription: "List of associated tags.",
							Computed:            true,
							ElementType:         types.Int64Type,
						},
						"id": schema.Int64Attribute{
							MarkdownDescription: "Metadata ID.",
							Computed:            true,
						},
						// Field values
						"add_collection_name": schema.BoolAttribute{
							MarkdownDescription: "Add collection name flag.",
							Computed:            true,
						},
						"use_movie_nfo": schema.BoolAttribute{
							MarkdownDescription: "Use movie nfo flag.",
							Computed:            true,
						},
						"movie_images": schema.BoolAttribute{
							MarkdownDescription: "Movie images flag.",
							Computed:            true,
						},
						"movie_metadata": schema.BoolAttribute{
							MarkdownDescription: "Movie metafata flag.",
							Computed:            true,
						},
						"movie_metadata_url": schema.BoolAttribute{
							MarkdownDescription: "Movie metadata URL flag.",
							Computed:            true,
						},
						"movie_metadata_language": schema.Int64Attribute{
							MarkdownDescription: "Movie metadata language.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *MetadataConsumersDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if client := helpers.DataSourceConfigure(ctx, req, resp); client != nil {
		d.client = client
	}
}

func (d *MetadataConsumersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *MetadataConsumers

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Get metadataConsumers current value
	response, _, err := d.client.MetadataApi.ListMetadata(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(helpers.ClientError, helpers.ParseClientError(helpers.List, metadataConsumersDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+metadataConsumersDataSourceName)
	// Map response body to resource schema attribute
	profiles := make([]Metadata, len(response))
	for i, p := range response {
		profiles[i].write(ctx, p)
	}

	tfsdk.ValueFrom(ctx, profiles, data.MetadataConsumers.Type(context.Background()), &data.MetadataConsumers)
	// TODO: remove ID once framework support tests without ID https://www.terraform.io/plugin/framework/acctests#implement-id-attribute
	data.ID = types.StringValue(strconv.Itoa(len(response)))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
