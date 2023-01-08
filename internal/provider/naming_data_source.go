package provider

import (
	"context"
	"fmt"

	"github.com/devopsarr/radarr-go/radarr"
	"github.com/devopsarr/terraform-provider-radarr/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const namingDataSourceName = "naming"

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &NamingDataSource{}

func NewNamingDataSource() datasource.DataSource {
	return &NamingDataSource{}
}

// NamingDataSource defines the naming implementation.
type NamingDataSource struct {
	client *radarr.APIClient
}

func (d *NamingDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + namingDataSourceName
}

func (d *NamingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the delay server.
		MarkdownDescription: "<!-- subcategory:Media Management -->[Naming](../resources/naming).",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				MarkdownDescription: "Delay Profile ID.",
				Computed:            true,
			},
			"include_quality": schema.BoolAttribute{
				MarkdownDescription: "Include quality in file name.",
				Computed:            true,
			},
			"rename_movies": schema.BoolAttribute{
				MarkdownDescription: "Radarr will use the existing file name if false.",
				Computed:            true,
			},
			"replace_illegal_characters": schema.BoolAttribute{
				MarkdownDescription: "Replace illegal characters. They will be removed if false.",
				Computed:            true,
			},
			"replace_spaces": schema.BoolAttribute{
				MarkdownDescription: "Replace spaces.",
				Computed:            true,
			},
			"colon_replacement_format": schema.StringAttribute{
				MarkdownDescription: "Change how Radarr handles colon replacement. Valid values are: 'delete', 'dash', 'spaceDash', and 'spaceDashSpace'.",
				Computed:            true,
			},
			"movie_folder_format": schema.StringAttribute{
				MarkdownDescription: "Movie folder format.",
				Computed:            true,
			},
			"standard_movie_format": schema.StringAttribute{
				MarkdownDescription: "Standard movie formatss.",
				Computed:            true,
			},
		},
	}
}

func (d *NamingDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*radarr.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			tools.UnexpectedDataSourceConfigureType,
			fmt.Sprintf("Expected *radarr.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *NamingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get naming current value
	response, _, err := d.client.NamingConfigApi.GetConfigNaming(ctx).Execute()
	if err != nil {
		resp.Diagnostics.AddError(tools.ClientError, fmt.Sprintf("Unable to read %s, got error: %s", namingDataSourceName, err))

		return
	}

	tflog.Trace(ctx, "read "+namingDataSourceName)

	state := Naming{}
	state.write(response)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}
